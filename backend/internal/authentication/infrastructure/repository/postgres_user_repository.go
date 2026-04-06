package repository

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"

	"opscore/backend/internal/authentication/domain/entity"
	domain_err "opscore/backend/internal/authentication/domain/err"
	"opscore/backend/internal/authentication/domain/repository"
	"opscore/backend/internal/authentication/domain/value_object"
)

// PostgresUserRepository is the PostgreSQL implementation of UserRepository.
// It manages the User aggregate root and its child Identity entities.
// Changes are tracked via domain events from the User aggregate.
type PostgresUserRepository struct {
	db *pgxpool.Pool
}

// NewPostgresUserRepository creates a new PostgresUserRepository.
func NewPostgresUserRepository(db *pgxpool.Pool) repository.UserRepository {
	return &PostgresUserRepository{db: db}
}

func (r *PostgresUserRepository) Save(ctx context.Context, user *entity.User) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return domain_err.NewDataPersistFailureError("user: begin transaction").WithParent(err)
	}
	defer tx.Rollback(ctx)

	for _, event := range user.GetEvents() {
		if err := r.processEvent(ctx, tx, event); err != nil {
			return err
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return domain_err.NewDataPersistFailureError("user: commit transaction").WithParent(err)
	}

	user.ClearEvents()
	return nil
}

// processEvent handles a single domain event within a transaction.
func (r *PostgresUserRepository) processEvent(ctx context.Context, tx pgx.Tx, event entity.DomainEvent) error {
	switch e := event.(type) {
	case *entity.UserCreatedEvent:
		return r.insertUser(ctx, tx, e)
	case *entity.IdentityAddedEvent:
		return r.insertIdentity(ctx, tx, e)
	case *entity.IdentityUpdatedEvent:
		return r.updateIdentity(ctx, tx, e)
	case *entity.IdentityRemovedEvent:
		return r.deleteIdentity(ctx, tx, e)
	case *entity.PrimaryIdentityChangedEvent:
		return r.changePrimaryIdentity(ctx, tx, e)
	case *entity.UserLastLoginUpdatedEvent:
		return r.updateUserLastLogin(ctx, tx, e)
	default:
		return nil
	}
}

func (r *PostgresUserRepository) insertUser(ctx context.Context, tx pgx.Tx, e *entity.UserCreatedEvent) error {
	query := `
		INSERT INTO users (id, email, display_name, picture_url, created_at, updated_at, last_login_at)
		VALUES ($1, $2, $3, $4, NOW(), NOW(), NOW())`

	_, err := tx.Exec(ctx, query,
		e.UserID.String(), e.Email, e.DisplayName, nullableString(e.PictureURL),
	)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return domain_err.NewDataConflictError("users", pgErr.ConstraintName).WithParent(err)
		}
		return domain_err.NewDataPersistFailureError("user").WithParent(err)
	}
	return nil
}

func (r *PostgresUserRepository) insertIdentity(ctx context.Context, tx pgx.Tx, e *entity.IdentityAddedEvent) error {
	query := `
		INSERT INTO user_identities (id, user_id, provider, provider_user_id, email, name, picture_url, is_primary, created_at, updated_at, last_used_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, NOW(), NOW(), NOW())`

	_, err := tx.Exec(ctx, query,
		e.IdentityID.String(), e.UserID.String(), e.Provider, e.ProviderUserID,
		e.Email, nullableString(e.Name), nullableString(e.PictureURL), e.IsPrimary,
	)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return domain_err.NewDataConflictError("user_identities", pgErr.ConstraintName).WithParent(err)
		}
		return domain_err.NewDataPersistFailureError("identity").WithParent(err)
	}
	return nil
}

func (r *PostgresUserRepository) updateIdentity(ctx context.Context, tx pgx.Tx, e *entity.IdentityUpdatedEvent) error {
	query := `
		UPDATE user_identities
		SET email = $1, name = $2, picture_url = $3, updated_at = NOW()
		WHERE id = $4`

	_, err := tx.Exec(ctx, query,
		e.Email, nullableString(e.Name), nullableString(e.PictureURL), e.IdentityID.String(),
	)
	if err != nil {
		return domain_err.NewDataPersistFailureError("identity").WithParent(err)
	}
	return nil
}

func (r *PostgresUserRepository) deleteIdentity(ctx context.Context, tx pgx.Tx, e *entity.IdentityRemovedEvent) error {
	query := `DELETE FROM user_identities WHERE id = $1 AND user_id = $2`

	_, err := tx.Exec(ctx, query, e.IdentityID.String(), e.UserID.String())
	if err != nil {
		return domain_err.NewDataPersistFailureError("identity").WithParent(err)
	}
	return nil
}

func (r *PostgresUserRepository) changePrimaryIdentity(ctx context.Context, tx pgx.Tx, e *entity.PrimaryIdentityChangedEvent) error {
	// Unset old primary
	if !e.OldIdentityID.IsEmpty() {
		query := `UPDATE user_identities SET is_primary = FALSE, updated_at = NOW() WHERE id = $1`
		if _, err := tx.Exec(ctx, query, e.OldIdentityID.String()); err != nil {
			return domain_err.NewDataPersistFailureError("identity: unset primary").WithParent(err)
		}
	}

	// Set new primary
	query := `UPDATE user_identities SET is_primary = TRUE, updated_at = NOW() WHERE id = $1`
	if _, err := tx.Exec(ctx, query, e.NewIdentityID.String()); err != nil {
		return domain_err.NewDataPersistFailureError("identity: set primary").WithParent(err)
	}
	return nil
}

func (r *PostgresUserRepository) updateUserLastLogin(ctx context.Context, tx pgx.Tx, e *entity.UserLastLoginUpdatedEvent) error {
	query := `UPDATE users SET last_login_at = NOW(), updated_at = NOW() WHERE id = $1`

	_, err := tx.Exec(ctx, query, e.UserID.String())
	if err != nil {
		return domain_err.NewDataPersistFailureError("user: last_login").WithParent(err)
	}
	return nil
}

func (r *PostgresUserRepository) FindByID(ctx context.Context, userID value_object.UserID) (*entity.User, error) {
	query := `
		SELECT id, email, display_name, picture_url, created_at, updated_at, last_login_at
		FROM users
		WHERE id = $1`

	user, err := r.scanUser(ctx, query, userID.String())
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (r *PostgresUserRepository) FindByIDWithIdentities(ctx context.Context, userID value_object.UserID) (*entity.User, error) {
	user, err := r.FindByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	identities, err := r.findIdentitiesByUserID(ctx, userID.String())
	if err != nil {
		return nil, err
	}

	user.LoadIdentities(identities)
	return user, nil
}

func (r *PostgresUserRepository) FindByProviderIdentity(ctx context.Context, provider, providerUserID string) (*entity.User, error) {
	// Find the user_id via user_identities, then load the full user with all identities.
	query := `
		SELECT u.id, u.email, u.display_name, u.picture_url, u.created_at, u.updated_at, u.last_login_at
		FROM users u
		INNER JOIN user_identities ui ON u.id = ui.user_id
		WHERE ui.provider = $1 AND ui.provider_user_id = $2`

	user, err := r.scanUser(ctx, query, provider, providerUserID)
	if err != nil {
		return nil, err
	}

	identities, err := r.findIdentitiesByUserID(ctx, user.ID().String())
	if err != nil {
		return nil, err
	}

	user.LoadIdentities(identities)
	return user, nil
}

func (r *PostgresUserRepository) Delete(ctx context.Context, userID value_object.UserID) error {
	query := `DELETE FROM users WHERE id = $1`

	result, err := r.db.Exec(ctx, query, userID.String())
	if err != nil {
		return domain_err.NewDataPersistFailureError("user").WithParent(err)
	}

	if result.RowsAffected() == 0 {
		return domain_err.NewNotFoundUserError(userID.String())
	}

	return nil
}

// scanUser executes a query and scans a single User row.
func (r *PostgresUserRepository) scanUser(ctx context.Context, query string, args ...interface{}) (*entity.User, error) {
	var (
		id          string
		email       string
		displayName string
		pictureURL  *string
		createdAt   time.Time
		updatedAt   time.Time
		lastLoginAt time.Time
	)

	err := r.db.QueryRow(ctx, query, args...).Scan(
		&id, &email, &displayName, &pictureURL, &createdAt, &updatedAt, &lastLoginAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain_err.NewNotFoundUserError(formatArgsForError(args))
		}
		return nil, domain_err.NewDataAccessFailureError("user").WithParent(err)
	}

	uid := value_object.ReconstructUserID(id)
	pic := ""
	if pictureURL != nil {
		pic = *pictureURL
	}
	return entity.ReconstructUser(uid, email, displayName, pic, createdAt, updatedAt, lastLoginAt), nil
}

// findIdentitiesByUserID loads all identities for a given user.
func (r *PostgresUserRepository) findIdentitiesByUserID(ctx context.Context, userID string) ([]*entity.Identity, error) {
	query := `
		SELECT id, user_id, provider, provider_user_id, email, name, picture_url, is_primary, created_at, updated_at, last_used_at
		FROM user_identities
		WHERE user_id = $1
		ORDER BY created_at ASC`

	rows, err := r.db.Query(ctx, query, userID)
	if err != nil {
		return nil, domain_err.NewDataAccessFailureError("identities").WithParent(err)
	}
	defer rows.Close()

	var identities []*entity.Identity
	for rows.Next() {
		identity, err := scanIdentityFromRow(rows)
		if err != nil {
			return nil, domain_err.NewDataAccessFailureError("identities").WithParent(err)
		}
		identities = append(identities, identity)
	}
	if err := rows.Err(); err != nil {
		return nil, domain_err.NewDataAccessFailureError("identities").WithParent(err)
	}

	return identities, nil
}

// scanIdentityFromRow scans a single Identity from a row.
func scanIdentityFromRow(rows pgx.Rows) (*entity.Identity, error) {
	var (
		id             string
		userID         string
		provider       string
		providerUserID string
		email          string
		name           *string
		pictureURL     *string
		isPrimary      bool
		createdAt      time.Time
		updatedAt      time.Time
		lastUsedAt     time.Time
	)

	err := rows.Scan(
		&id, &userID, &provider, &providerUserID, &email,
		&name, &pictureURL, &isPrimary, &createdAt, &updatedAt, &lastUsedAt,
	)
	if err != nil {
		return nil, err
	}

	identityID := value_object.ReconstructIdentityID(id)
	uid := value_object.ReconstructUserID(userID)
	nameStr := ""
	if name != nil {
		nameStr = *name
	}
	picStr := ""
	if pictureURL != nil {
		picStr = *pictureURL
	}

	return entity.ReconstructIdentity(
		identityID, uid, provider, providerUserID, email, nameStr, picStr,
		isPrimary, createdAt, updatedAt, lastUsedAt,
	), nil
}

// nullableString returns nil for empty strings (for nullable TEXT/VARCHAR columns).
func nullableString(s string) interface{} {
	if s == "" {
		return nil
	}
	return s
}

// formatArgsForError formats query arguments for error messages.
func formatArgsForError(args []interface{}) string {
	if len(args) > 0 {
		if s, ok := args[0].(string); ok {
			return s
		}
	}
	return "unknown"
}
