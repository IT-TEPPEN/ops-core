package persistence

import (
	"context"
	"database/sql"
	"errors" // Import errors package
	"fmt"    // For error wrapping
	"time"

	"opscore/backend/domain/model"
	"opscore/backend/domain/repository"

	"github.com/jackc/pgx/v5"        // For checking specific errors
	"github.com/jackc/pgx/v5/pgconn" // Import pgconn for error handling
	"github.com/jackc/pgx/v5/pgxpool"
)

// PostgresRepository is a PostgreSQL implementation of the repository.Repository interface.
type PostgresRepository struct {
	db *pgxpool.Pool
}

// NewPostgresRepository creates a new PostgresRepository.
func NewPostgresRepository(db *pgxpool.Pool) repository.Repository {
	return &PostgresRepository{db: db}
}

// Save persists a repository in the PostgreSQL database.
func (r *PostgresRepository) Save(ctx context.Context, repo model.Repository) error {
	query := `
		INSERT INTO repositories (id, name, url, access_token, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6)
		ON CONFLICT (id) DO UPDATE SET
			name = EXCLUDED.name,
			url = EXCLUDED.url,
			access_token = EXCLUDED.access_token,
			updated_at = EXCLUDED.updated_at;
	`
	_, err := r.db.Exec(ctx, query, repo.ID(), repo.Name(), repo.URL(), repo.AccessToken(), repo.CreatedAt(), repo.UpdatedAt())
	if err != nil {
		// Check for unique constraint violation on URL if a separate constraint exists
		// var pgErr *pgconn.PgError
		// if errors.As(err, &pgErr) && pgErr.Code == "23505" { // 23505 is unique_violation
		//  // Check constraint name if needed to differentiate between ID and URL conflicts
		// 	return repository.ErrRepositoryAlreadyExists // Or a more specific error
		// }
		return fmt.Errorf("failed to save repository: %w", err)
	}
	return nil
}

// FindByURL retrieves a repository by its URL from the PostgreSQL database.
func (r *PostgresRepository) FindByURL(ctx context.Context, url string) (model.Repository, error) {
	query := `
		SELECT id, name, url, access_token, created_at, updated_at
		FROM repositories
		WHERE url = $1;
	`
	var id, name, repoURL string
	var accessToken sql.NullString // トークンは NULL の可能性があるため sql.NullString を使用
	var createdAt, updatedAt time.Time

	err := r.db.QueryRow(ctx, query, url).Scan(
		&id,
		&name,
		&repoURL,
		&accessToken,
		&createdAt,
		&updatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil // Not found, no error
		}
		return nil, fmt.Errorf("failed to find repository by URL: %w", err)
	}

	// NullString から通常の string へ変換
	tokenStr := ""
	if accessToken.Valid {
		tokenStr = accessToken.String
	}

	return model.ReconstructRepository(id, name, repoURL, tokenStr, createdAt, updatedAt), nil
}

// FindByID retrieves a repository by its ID from the PostgreSQL database.
func (r *PostgresRepository) FindByID(ctx context.Context, id string) (model.Repository, error) {
	query := `
		SELECT id, name, url, access_token, created_at, updated_at
		FROM repositories
		WHERE id = $1;
	`
	var repoID, name, url string
	var accessToken sql.NullString
	var createdAt, updatedAt time.Time

	err := r.db.QueryRow(ctx, query, id).Scan(
		&repoID,
		&name,
		&url,
		&accessToken,
		&createdAt,
		&updatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			// Return a specific error type if needed by use cases
			// return nil, repository.ErrNotFound
			return nil, nil // Not found, no error
		}
		return nil, fmt.Errorf("failed to find repository by ID: %w", err)
	}

	// NullString から通常の string へ変換
	tokenStr := ""
	if accessToken.Valid {
		tokenStr = accessToken.String
	}

	return model.ReconstructRepository(repoID, name, url, tokenStr, createdAt, updatedAt), nil
}

// FindAll retrieves all repositories from the PostgreSQL database.
func (r *PostgresRepository) FindAll(ctx context.Context) ([]model.Repository, error) {
	query := `
		SELECT id, name, url, access_token, created_at, updated_at
		FROM repositories
		ORDER BY created_at DESC;
	`

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to query repositories: %w", err)
	}
	defer rows.Close()

	var repositories []model.Repository
	for rows.Next() {
		var id, name, url string
		var accessToken sql.NullString
		var createdAt, updatedAt time.Time

		if err := rows.Scan(&id, &name, &url, &accessToken, &createdAt, &updatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan repository row: %w", err)
		}

		// NullString から通常の string へ変換
		tokenStr := ""
		if accessToken.Valid {
			tokenStr = accessToken.String
		}

		repo := model.ReconstructRepository(id, name, url, tokenStr, createdAt, updatedAt)
		repositories = append(repositories, repo)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over repository rows: %w", err)
	}

	return repositories, nil
}

// SaveManagedFiles saves the list of selected file paths for a repository.
// It first deletes existing entries for the repoID and then inserts the new ones.
func (r *PostgresRepository) SaveManagedFiles(ctx context.Context, repoID string, filePaths []string) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(ctx) // Rollback if anything fails

	// Delete existing managed files for this repository
	deleteQuery := `DELETE FROM managed_files WHERE repository_id = $1;`
	_, err = tx.Exec(ctx, deleteQuery, repoID)
	if err != nil {
		return fmt.Errorf("failed to delete existing managed files: %w", err)
	}

	// Insert new managed files if any paths are provided
	if len(filePaths) > 0 {
		insertQuery := `
			INSERT INTO managed_files (repository_id, file_path)
			VALUES ($1, $2);
		`
		// Use Batch for potentially better performance with many files
		batch := &pgx.Batch{}
		for _, filePath := range filePaths {
			batch.Queue(insertQuery, repoID, filePath)
		}

		results := tx.SendBatch(ctx, batch)
		// Check results for errors
		for i := 0; i < len(filePaths); i++ {
			_, err = results.Exec()
			if err != nil {
				// Check for foreign key violation (repository_id doesn't exist)
				var pgErr *pgconn.PgError
				if errors.As(err, &pgErr) && pgErr.Code == "23503" { // foreign_key_violation
					// Consider returning a domain-specific error like repository.ErrNotFound
					return fmt.Errorf("failed to insert managed file: repository with ID %s not found: %w", repoID, err)
				}
				return fmt.Errorf("failed to insert managed file '%s': %w", filePaths[i], err)
			}
		}
		err = results.Close()
		if err != nil {
			return fmt.Errorf("failed to close batch results: %w", err)
		}
	}

	// Commit the transaction
	err = tx.Commit(ctx)
	if err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

// GetManagedFiles retrieves the list of selected file paths for a repository.
func (r *PostgresRepository) GetManagedFiles(ctx context.Context, repoID string) ([]string, error) {
	query := `
		SELECT file_path
		FROM managed_files
		WHERE repository_id = $1
		ORDER BY file_path; -- Optional: order for consistency
	`
	rows, err := r.db.Query(ctx, query, repoID)
	if err != nil {
		// Check if the error is ErrNoRows, although Query usually doesn't return it directly
		if errors.Is(err, pgx.ErrNoRows) {
			return []string{}, nil // No files found is not an error, return empty slice
		}
		return nil, fmt.Errorf("failed to query managed files: %w", err)
	}
	defer rows.Close()

	var filePaths []string
	for rows.Next() {
		var filePath string
		if err := rows.Scan(&filePath); err != nil {
			return nil, fmt.Errorf("failed to scan managed file path: %w", err)
		}
		filePaths = append(filePaths, filePath)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over managed files rows: %w", err)
	}

	// If no rows were found, filePaths will be an empty slice, which is correct.
	return filePaths, nil
}

// UpdateAccessToken updates the access token for a repository.
func (r *PostgresRepository) UpdateAccessToken(ctx context.Context, repoID string, accessToken string) error {
	query := `
		UPDATE repositories
		SET access_token = $1, updated_at = $2
		WHERE id = $3;
	`
	now := time.Now()
	res, err := r.db.Exec(ctx, query, accessToken, now, repoID)
	if err != nil {
		return fmt.Errorf("failed to update repository access token: %w", err)
	}

	// Check if repository exists
	if res.RowsAffected() == 0 {
		return fmt.Errorf("repository with ID %s not found", repoID)
	}

	return nil
}
