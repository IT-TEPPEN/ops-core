package persistence

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"opscore/backend/internal/document/domain/entity"
	"opscore/backend/internal/document/domain/repository"
	"opscore/backend/internal/document/domain/value_object"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

// PostgresDocumentRepository provides a PostgreSQL implementation of DocumentRepository.
type PostgresDocumentRepository struct {
	db *pgxpool.Pool
}

// NewPostgresDocumentRepository creates a new repository instance.
func NewPostgresDocumentRepository(db *pgxpool.Pool) repository.DocumentRepository {
	return &PostgresDocumentRepository{db: db}
}

// Save creates or updates a document and its versions.
func (r *PostgresDocumentRepository) Save(ctx context.Context, document entity.Document) error {
	return r.persistDocument(ctx, document)
}

// Update updates an existing document and its versions.
func (r *PostgresDocumentRepository) Update(ctx context.Context, document entity.Document) error {
	return r.persistDocument(ctx, document)
}

// FindByID retrieves a document by its ID.
func (r *PostgresDocumentRepository) FindByID(ctx context.Context, id value_object.DocumentID) (entity.Document, error) {
	query := `
		SELECT id, repository_id, provider_repository_id, owner, repository, is_published, is_auto_update, access_scope, current_version_id, created_at, updated_at
		FROM documents
		WHERE id = $1;
	`

	row := r.db.QueryRow(ctx, query, id.String())
	docRec, err := scanDocumentRow(row)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return r.toDomainDocument(ctx, docRec)
}

// FindByRepositoryID retrieves all documents for a repository.
func (r *PostgresDocumentRepository) FindByRepositoryID(ctx context.Context, repoID value_object.RepositoryID) ([]entity.Document, error) {
	query := `
		SELECT id, repository_id, provider_repository_id, owner, repository, is_published, is_auto_update, access_scope, current_version_id, created_at, updated_at
		FROM documents
		WHERE repository_id = $1
		ORDER BY created_at DESC;
	`

	rows, err := r.db.Query(ctx, query, repoID.String())
	if err != nil {
		return nil, fmt.Errorf("failed to query documents: %w", err)
	}
	defer rows.Close()

	var docs []entity.Document
	for rows.Next() {
		docRec, scanErr := scanDocumentRow(rows)
		if scanErr != nil {
			return nil, scanErr
		}
		doc, convErr := r.toDomainDocument(ctx, docRec)
		if convErr != nil {
			return nil, convErr
		}
		docs = append(docs, doc)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating document rows: %w", err)
	}

	return docs, nil
}

// FindPublished retrieves published documents with optional filters.
func (r *PostgresDocumentRepository) FindPublished(ctx context.Context, filters ...repository.Filter) ([]entity.Document, error) {
	ctx = repository.ApplyFilters(ctx, filters...)
	parsed := repository.ParseFilters(ctx)

	query := `
		SELECT id, repository_id, provider_repository_id, owner, repository, is_published, is_auto_update, access_scope, current_version_id, created_at, updated_at
		FROM documents
		WHERE is_published = TRUE
	`

	args := []interface{}{}
	idx := 1

	if parsed.RepositoryID != "" {
		query += fmt.Sprintf(" AND repository_id = $%d", idx)
		args = append(args, parsed.RepositoryID)
		idx++
	}
	if parsed.ProviderRepositoryID != "" {
		query += fmt.Sprintf(" AND provider_repository_id = $%d", idx)
		args = append(args, parsed.ProviderRepositoryID)
		idx++
	}
	if parsed.Owner != "" {
		query += fmt.Sprintf(" AND owner = $%d", idx)
		args = append(args, parsed.Owner)
		idx++
	}
	if parsed.Repository != "" {
		query += fmt.Sprintf(" AND repository = $%d", idx)
		args = append(args, parsed.Repository)
		idx++
	}

	query += " ORDER BY created_at DESC"

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to query published documents: %w", err)
	}
	defer rows.Close()

	var docs []entity.Document
	for rows.Next() {
		docRec, scanErr := scanDocumentRow(rows)
		if scanErr != nil {
			return nil, scanErr
		}
		doc, convErr := r.toDomainDocument(ctx, docRec)
		if convErr != nil {
			return nil, convErr
		}
		docs = append(docs, doc)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating published document rows: %w", err)
	}

	return docs, nil
}

// Delete deletes a document by ID.
func (r *PostgresDocumentRepository) Delete(ctx context.Context, id value_object.DocumentID) error {
	_, err := r.db.Exec(ctx, "DELETE FROM documents WHERE id = $1;", id.String())
	if err != nil {
		return fmt.Errorf("failed to delete document: %w", err)
	}
	return nil
}

// SaveVersion saves a document version.
func (r *PostgresDocumentRepository) SaveVersion(ctx context.Context, version entity.DocumentVersion) error {
	return r.upsertVersion(ctx, nil, version)
}

// FindVersionsByDocumentID retrieves all versions for a document.
func (r *PostgresDocumentRepository) FindVersionsByDocumentID(ctx context.Context, docID value_object.DocumentID) ([]entity.DocumentVersion, error) {
	rows, err := r.db.Query(ctx, `
		SELECT id, document_id, version_number, file_path, commit_hash, title, doc_type, tags, variables, content, published_at, unpublished_at
		FROM document_versions
		WHERE document_id = $1
		ORDER BY version_number ASC;
	`, docID.String())
	if err != nil {
		return nil, fmt.Errorf("failed to query document versions: %w", err)
	}
	defer rows.Close()

	versions, _, err := r.scanVersions(ctx, rows, sql.NullString{})
	if err != nil {
		return nil, err
	}
	return versions, nil
}

// FindVersionByNumber retrieves a specific version by document ID and version number.
func (r *PostgresDocumentRepository) FindVersionByNumber(ctx context.Context, docID value_object.DocumentID, versionNumber value_object.VersionNumber) (entity.DocumentVersion, error) {
	var rec versionRecord
	var variablesBytes []byte
	var unpublishedAt sql.NullTime

	err := r.db.QueryRow(ctx, `
		SELECT id, document_id, version_number, file_path, commit_hash, title, doc_type, tags, variables, content, published_at, unpublished_at
		FROM document_versions
		WHERE document_id = $1 AND version_number = $2;
	`, docID.String(), versionNumber.Int()).Scan(
		&rec.ID,
		&rec.DocumentID,
		&rec.VersionNumber,
		&rec.FilePath,
		&rec.CommitHash,
		&rec.Title,
		&rec.DocType,
		&rec.Tags,
		&variablesBytes,
		&rec.Content,
		&rec.PublishedAt,
		&unpublishedAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to find document version: %w", err)
	}

	version, err := r.buildVersion(rec, variablesBytes, unpublishedAt, sql.NullString{})
	if err != nil {
		return nil, err
	}

	return version, nil
}

// scanDocumentRow reads a document row into a record struct.
func scanDocumentRow(row pgx.Row) (documentRecord, error) {
	var rec documentRecord
	err := row.Scan(
		&rec.ID,
		&rec.RepositoryID,
		&rec.ProviderRepositoryID,
		&rec.Owner,
		&rec.Repository,
		&rec.IsPublished,
		&rec.IsAutoUpdate,
		&rec.AccessScope,
		&rec.CurrentVersionID,
		&rec.CreatedAt,
		&rec.UpdatedAt,
	)
	if err != nil {
		return documentRecord{}, fmt.Errorf("failed to scan document row: %w", err)
	}
	return rec, nil
}

// persistDocument upserts the document and associated versions in a single transaction.
func (r *PostgresDocumentRepository) persistDocument(ctx context.Context, document entity.Document) error {
	origin := document.Origin()
	if origin == nil {
		return fmt.Errorf("document origin is required")
	}

	tx, err := r.db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	// Upsert document without current_version_id first to satisfy FK ordering.
	_, err = tx.Exec(ctx, `
		INSERT INTO documents (id, repository_id, provider_repository_id, owner, repository, is_published, is_auto_update, access_scope, current_version_id, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, NULL, $9, $10)
		ON CONFLICT (id) DO UPDATE SET
			repository_id = EXCLUDED.repository_id,
			provider_repository_id = EXCLUDED.provider_repository_id,
			owner = EXCLUDED.owner,
			repository = EXCLUDED.repository,
			is_published = EXCLUDED.is_published,
			is_auto_update = EXCLUDED.is_auto_update,
			access_scope = EXCLUDED.access_scope,
			updated_at = EXCLUDED.updated_at;
	`,
		document.ID().String(),
		document.RepositoryID().String(),
		origin.ProviderRepositoryID().String(),
		origin.Owner(),
		origin.Repository(),
		document.IsPublished(),
		document.IsAutoUpdate(),
		document.AccessScope().String(),
		document.CreatedAt(),
		document.UpdatedAt(),
	)
	if err != nil {
		return fmt.Errorf("failed to upsert document: %w", err)
	}

	// Upsert versions
	for _, v := range document.Versions() {
		if err := r.upsertVersion(ctx, tx, v); err != nil {
			return err
		}
	}

	// Update current_version_id if applicable
	var currentVersionID interface{}
	if cv := document.CurrentVersion(); cv != nil && cv.IsCurrentVersion() {
		currentVersionID = cv.ID().String()
	}

	_, err = tx.Exec(ctx, `
		UPDATE documents
		SET current_version_id = $1, is_published = $2, is_auto_update = $3, access_scope = $4, updated_at = $5
		WHERE id = $6;
	`,
		currentVersionID,
		document.IsPublished(),
		document.IsAutoUpdate(),
		document.AccessScope().String(),
		time.Now(),
		document.ID().String(),
	)
	if err != nil {
		return fmt.Errorf("failed to update document current version: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}
	return nil
}

// upsertVersion inserts or updates a version. If tx is nil, it operates without a transaction.
func (r *PostgresDocumentRepository) upsertVersion(ctx context.Context, tx pgx.Tx, version entity.DocumentVersion) error {
	tags := make([]string, 0, len(version.Tags()))
	for _, t := range version.Tags() {
		tags = append(tags, t.String())
	}

	varDefs := make([]variableDefinitionRecord, 0, len(version.Variables()))
	for _, v := range version.Variables() {
		varDefs = append(varDefs, variableDefinitionRecord{
			Name:         v.Name(),
			Label:        v.Label(),
			Description:  v.Description(),
			Type:         v.Type().String(),
			Required:     v.Required(),
			DefaultValue: v.DefaultValue(),
		})
	}

	var varBytes []byte
	if len(varDefs) > 0 {
		b, err := json.Marshal(varDefs)
		if err != nil {
			return fmt.Errorf("failed to marshal variables: %w", err)
		}
		varBytes = b
	}

	execer := execAdapter{db: r.db, tx: tx}
	_, err := execer.Exec(ctx, `
		INSERT INTO document_versions (id, document_id, version_number, file_path, commit_hash, title, doc_type, tags, variables, content, published_at, unpublished_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
		ON CONFLICT (id) DO UPDATE SET
			version_number = EXCLUDED.version_number,
			file_path = EXCLUDED.file_path,
			commit_hash = EXCLUDED.commit_hash,
			title = EXCLUDED.title,
			doc_type = EXCLUDED.doc_type,
			tags = EXCLUDED.tags,
			variables = EXCLUDED.variables,
			content = EXCLUDED.content,
			published_at = EXCLUDED.published_at,
			unpublished_at = EXCLUDED.unpublished_at;
	`,
		version.ID().String(),
		version.DocumentID().String(),
		version.VersionNumber().Int(),
		version.Source().FilePath().String(),
		version.Source().CommitHash().String(),
		version.Title(),
		version.Type().String(),
		tags,
		varBytes,
		version.Content(),
		version.PublishedAt(),
		nullableTime(version.UnpublishedAt()),
	)
	if err != nil {
		return fmt.Errorf("failed to upsert document version: %w", err)
	}

	return nil
}

// toDomainDocument reconstructs a domain entity from database records.
func (r *PostgresDocumentRepository) toDomainDocument(ctx context.Context, rec documentRecord) (entity.Document, error) {
	docID, err := value_object.NewDocumentID(rec.ID)
	if err != nil {
		return nil, fmt.Errorf("invalid document id: %w", err)
	}

	repoID, err := value_object.NewRepositoryID(rec.RepositoryID)
	if err != nil {
		return nil, fmt.Errorf("invalid repository id: %w", err)
	}

	origin, err := value_object.NewRepositoryOrigin(rec.ProviderRepositoryID, rec.Owner, rec.Repository)
	if err != nil {
		return nil, fmt.Errorf("invalid repository origin: %w", err)
	}

	accessScope, err := value_object.NewAccessScope(rec.AccessScope)
	if err != nil {
		return nil, fmt.Errorf("invalid access scope: %w", err)
	}

	versions, currentVersion, err := r.loadVersionsForDocument(ctx, docID.String(), rec.CurrentVersionID)
	if err != nil {
		return nil, err
	}

	doc := entity.ReconstructDocument(
		docID,
		repoID,
		&origin,
		rec.IsPublished,
		rec.IsAutoUpdate,
		accessScope,
		currentVersion,
		versions,
		rec.CreatedAt,
		rec.UpdatedAt,
	)
	return doc, nil
}

// loadVersionsForDocument fetches versions for a document and returns slice plus the current version if available.
func (r *PostgresDocumentRepository) loadVersionsForDocument(ctx context.Context, docID string, currentVersionID sql.NullString) ([]entity.DocumentVersion, entity.DocumentVersion, error) {
	rows, err := r.db.Query(ctx, `
		SELECT id, document_id, version_number, file_path, commit_hash, title, doc_type, tags, variables, content, published_at, unpublished_at
		FROM document_versions
		WHERE document_id = $1
		ORDER BY version_number ASC;
	`, docID)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to query versions: %w", err)
	}
	defer rows.Close()

	versions, current, err := r.scanVersions(ctx, rows, currentVersionID)
	if err != nil {
		return nil, nil, err
	}

	return versions, current, nil
}

// scanVersions converts rows into domain versions and determines the current version.
func (r *PostgresDocumentRepository) scanVersions(ctx context.Context, rows pgx.Rows, currentVersionID sql.NullString) ([]entity.DocumentVersion, entity.DocumentVersion, error) {
	var versions []entity.DocumentVersion
	var current entity.DocumentVersion

	for rows.Next() {
		var rec versionRecord
		var variablesBytes []byte
		var unpublishedAt sql.NullTime

		err := rows.Scan(
			&rec.ID,
			&rec.DocumentID,
			&rec.VersionNumber,
			&rec.FilePath,
			&rec.CommitHash,
			&rec.Title,
			&rec.DocType,
			&rec.Tags,
			&variablesBytes,
			&rec.Content,
			&rec.PublishedAt,
			&unpublishedAt,
		)
		if err != nil {
			if err == pgx.ErrNoRows {
				return versions, current, nil
			}
			return nil, nil, fmt.Errorf("failed to scan document version: %w", err)
		}

		version, err := r.buildVersion(rec, variablesBytes, unpublishedAt, currentVersionID)
		if err != nil {
			return nil, nil, err
		}

		versions = append(versions, version)
		if currentVersionID.Valid && currentVersionID.String == rec.ID {
			current = version
		}
	}

	if err := rows.Err(); err != nil {
		return nil, nil, fmt.Errorf("error iterating version rows: %w", err)
	}

	return versions, current, nil
}

// buildVersion constructs a DocumentVersion from scanned values.
func (r *PostgresDocumentRepository) buildVersion(rec versionRecord, variablesBytes []byte, unpublishedAt sql.NullTime, currentVersionID sql.NullString) (entity.DocumentVersion, error) {
	vars, convErr := toVariableDefinitions(variablesBytes)
	if convErr != nil {
		return nil, convErr
	}

	filePath, err := value_object.NewFilePath(rec.FilePath)
	if err != nil {
		return nil, fmt.Errorf("invalid file path: %w", err)
	}
	commitHash, err := value_object.NewCommitHash(rec.CommitHash)
	if err != nil {
		return nil, fmt.Errorf("invalid commit hash: %w", err)
	}
	source, err := value_object.NewDocumentSource(filePath, commitHash)
	if err != nil {
		return nil, fmt.Errorf("invalid document source: %w", err)
	}

	docType, err := value_object.NewDocumentType(rec.DocType)
	if err != nil {
		return nil, fmt.Errorf("invalid document type: %w", err)
	}

	tagValues := make([]value_object.Tag, 0, len(rec.Tags))
	for _, t := range rec.Tags {
		tag, tagErr := value_object.NewTag(t)
		if tagErr != nil {
			return nil, fmt.Errorf("invalid tag: %w", tagErr)
		}
		tagValues = append(tagValues, tag)
	}

	versionNumber, err := value_object.NewVersionNumber(rec.VersionNumber)
	if err != nil {
		return nil, fmt.Errorf("invalid version number: %w", err)
	}
	versionID, err := value_object.NewVersionID(rec.ID)
	if err != nil {
		return nil, fmt.Errorf("invalid version id: %w", err)
	}
	documentID, err := value_object.NewDocumentID(rec.DocumentID)
	if err != nil {
		return nil, fmt.Errorf("invalid document id on version: %w", err)
	}

	var unpublishedPtr *time.Time
	if unpublishedAt.Valid {
		unpublishedPtr = &unpublishedAt.Time
	}

	isCurrent := currentVersionID.Valid && currentVersionID.String == rec.ID

	version := entity.ReconstructDocumentVersion(
		versionID,
		documentID,
		versionNumber,
		source,
		rec.Title,
		docType,
		tagValues,
		vars,
		rec.Content,
		rec.PublishedAt,
		unpublishedPtr,
		isCurrent,
	)

	return version, nil
}

// toVariableDefinitions converts JSONB bytes into value objects.
func toVariableDefinitions(raw []byte) ([]value_object.VariableDefinition, error) {
	if len(raw) == 0 {
		return nil, nil
	}

	var records []variableDefinitionRecord
	if err := json.Unmarshal(raw, &records); err != nil {
		return nil, fmt.Errorf("failed to unmarshal variables: %w", err)
	}

	vars := make([]value_object.VariableDefinition, 0, len(records))
	for i, r := range records {
		t, err := value_object.NewVariableType(r.Type)
		if err != nil {
			return nil, fmt.Errorf("invalid variable type at index %d: %w", i, err)
		}
		v, err := value_object.NewVariableDefinition(r.Name, r.Label, r.Description, t, r.Required, r.DefaultValue)
		if err != nil {
			return nil, fmt.Errorf("invalid variable definition at index %d: %w", i, err)
		}
		vars = append(vars, v)
	}

	return vars, nil
}

// nullableTime returns the underlying time or nil.
func nullableTime(t *time.Time) interface{} {
	if t == nil {
		return nil
	}
	return *t
}

// execAdapter allows using either a transaction or the pool for Exec operations.
type execAdapter struct {
	db *pgxpool.Pool
	tx pgx.Tx
}

func (e execAdapter) Exec(ctx context.Context, sql string, args ...interface{}) (pgconn.CommandTag, error) {
	if e.tx != nil {
		return e.tx.Exec(ctx, sql, args...)
	}
	return e.db.Exec(ctx, sql, args...)
}

// documentRecord holds raw fields from the documents table.
type documentRecord struct {
	ID                   string
	RepositoryID         string
	ProviderRepositoryID string
	Owner                string
	Repository           string
	IsPublished          bool
	IsAutoUpdate         bool
	AccessScope          string
	CurrentVersionID     sql.NullString
	CreatedAt            time.Time
	UpdatedAt            time.Time
}

// versionRecord holds raw fields from the document_versions table.
type versionRecord struct {
	ID            string
	DocumentID    string
	VersionNumber int
	FilePath      string
	CommitHash    string
	Title         string
	DocType       string
	Tags          []string
	Content       string
	PublishedAt   time.Time
}

// variableDefinitionRecord mirrors the JSONB shape stored in document_versions.variables.
type variableDefinitionRecord struct {
	Name         string      `json:"name"`
	Label        string      `json:"label"`
	Description  string      `json:"description"`
	Type         string      `json:"type"`
	Required     bool        `json:"required"`
	DefaultValue interface{} `json:"defaultValue"`
}
