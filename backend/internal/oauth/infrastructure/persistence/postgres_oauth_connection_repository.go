package persistence

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"opscore/backend/internal/git_repository/infrastructure/encryption"
	"opscore/backend/internal/oauth/domain"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// PostgresOAuthConnectionRepository implements domain.OAuthConnectionRepository using PostgreSQL
type PostgresOAuthConnectionRepository struct {
	db        *pgxpool.Pool
	encryptor *encryption.Encryptor
}

// NewPostgresOAuthConnectionRepository creates a new PostgresOAuthConnectionRepository
func NewPostgresOAuthConnectionRepository(db *pgxpool.Pool, encryptor *encryption.Encryptor) *PostgresOAuthConnectionRepository {
	return &PostgresOAuthConnectionRepository{
		db:        db,
		encryptor: encryptor,
	}
}

// Save creates or updates an OAuth connection
func (r *PostgresOAuthConnectionRepository) Save(ctx context.Context, conn *domain.OAuthConnection) error {
	// Encrypt tokens
	encryptedAccessToken, err := r.encryptor.Encrypt(conn.AccessToken())
	if err != nil {
		return err
	}

	var encryptedRefreshToken *string
	if conn.RefreshToken() != "" {
		encrypted, err := r.encryptor.Encrypt(conn.RefreshToken())
		if err != nil {
			return err
		}
		encryptedRefreshToken = &encrypted
	}

	query := `
		INSERT INTO oauth_connections (
			id, user_id, provider, provider_user_id, provider_username,
			access_token_encrypted, refresh_token_encrypted, token_expires_at,
			scopes, created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
		ON CONFLICT (user_id, provider) DO UPDATE SET
			provider_user_id = EXCLUDED.provider_user_id,
			provider_username = EXCLUDED.provider_username,
			access_token_encrypted = EXCLUDED.access_token_encrypted,
			refresh_token_encrypted = EXCLUDED.refresh_token_encrypted,
			token_expires_at = EXCLUDED.token_expires_at,
			scopes = EXCLUDED.scopes,
			updated_at = EXCLUDED.updated_at
	`

	_, err = r.db.Exec(ctx, query,
		conn.ID(),
		conn.UserID(),
		string(conn.Provider()),
		conn.ProviderUserID(),
		conn.ProviderUsername(),
		encryptedAccessToken,
		encryptedRefreshToken,
		conn.TokenExpiresAt(),
		conn.Scopes(),
		conn.CreatedAt(),
		conn.UpdatedAt(),
	)

	return err
}

// FindByID finds an OAuth connection by ID
func (r *PostgresOAuthConnectionRepository) FindByID(ctx context.Context, id string) (*domain.OAuthConnection, error) {
	query := `
		SELECT id, user_id, provider, provider_user_id, provider_username,
			   access_token_encrypted, refresh_token_encrypted, token_expires_at,
			   scopes, created_at, updated_at
		FROM oauth_connections
		WHERE id = $1
	`

	return r.scanConnection(ctx, query, id)
}

// FindByUserAndProvider finds an OAuth connection by user ID and provider
func (r *PostgresOAuthConnectionRepository) FindByUserAndProvider(ctx context.Context, userID string, provider domain.Provider) (*domain.OAuthConnection, error) {
	query := `
		SELECT id, user_id, provider, provider_user_id, provider_username,
			   access_token_encrypted, refresh_token_encrypted, token_expires_at,
			   scopes, created_at, updated_at
		FROM oauth_connections
		WHERE user_id = $1 AND provider = $2
	`

	return r.scanConnection(ctx, query, userID, string(provider))
}

// FindByUser finds all OAuth connections for a user
func (r *PostgresOAuthConnectionRepository) FindByUser(ctx context.Context, userID string) ([]*domain.OAuthConnection, error) {
	query := `
		SELECT id, user_id, provider, provider_user_id, provider_username,
			   access_token_encrypted, refresh_token_encrypted, token_expires_at,
			   scopes, created_at, updated_at
		FROM oauth_connections
		WHERE user_id = $1
	`

	rows, err := r.db.Query(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var connections []*domain.OAuthConnection
	for rows.Next() {
		conn, err := r.scanRow(rows)
		if err != nil {
			return nil, err
		}
		connections = append(connections, conn)
	}

	return connections, rows.Err()
}

// Delete deletes an OAuth connection
func (r *PostgresOAuthConnectionRepository) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM oauth_connections WHERE id = $1`
	_, err := r.db.Exec(ctx, query, id)
	return err
}

// DeleteByUserAndProvider deletes an OAuth connection by user ID and provider
func (r *PostgresOAuthConnectionRepository) DeleteByUserAndProvider(ctx context.Context, userID string, provider domain.Provider) error {
	query := `DELETE FROM oauth_connections WHERE user_id = $1 AND provider = $2`
	_, err := r.db.Exec(ctx, query, userID, string(provider))
	return err
}

func (r *PostgresOAuthConnectionRepository) scanConnection(ctx context.Context, query string, args ...interface{}) (*domain.OAuthConnection, error) {
	row := r.db.QueryRow(ctx, query, args...)

	var (
		id                    string
		userID                string
		provider              string
		providerUserID        string
		providerUsername      sql.NullString
		accessTokenEncrypted  string
		refreshTokenEncrypted sql.NullString
		tokenExpiresAt        sql.NullTime
		scopes                []string
		createdAt             time.Time
		updatedAt             time.Time
	)

	err := row.Scan(
		&id, &userID, &provider, &providerUserID, &providerUsername,
		&accessTokenEncrypted, &refreshTokenEncrypted, &tokenExpiresAt,
		&scopes, &createdAt, &updatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return r.buildConnection(
		id, userID, provider, providerUserID, providerUsername,
		accessTokenEncrypted, refreshTokenEncrypted, tokenExpiresAt,
		scopes, createdAt, updatedAt,
	)
}

func (r *PostgresOAuthConnectionRepository) scanRow(rows pgx.Rows) (*domain.OAuthConnection, error) {
	var (
		id                    string
		userID                string
		provider              string
		providerUserID        string
		providerUsername      sql.NullString
		accessTokenEncrypted  string
		refreshTokenEncrypted sql.NullString
		tokenExpiresAt        sql.NullTime
		scopes                []string
		createdAt             time.Time
		updatedAt             time.Time
	)

	err := rows.Scan(
		&id, &userID, &provider, &providerUserID, &providerUsername,
		&accessTokenEncrypted, &refreshTokenEncrypted, &tokenExpiresAt,
		&scopes, &createdAt, &updatedAt,
	)
	if err != nil {
		return nil, err
	}

	return r.buildConnection(
		id, userID, provider, providerUserID, providerUsername,
		accessTokenEncrypted, refreshTokenEncrypted, tokenExpiresAt,
		scopes, createdAt, updatedAt,
	)
}

func (r *PostgresOAuthConnectionRepository) buildConnection(
	id, userID, provider, providerUserID string,
	providerUsername sql.NullString,
	accessTokenEncrypted string,
	refreshTokenEncrypted sql.NullString,
	tokenExpiresAt sql.NullTime,
	scopes []string,
	createdAt, updatedAt time.Time,
) (*domain.OAuthConnection, error) {
	// Decrypt tokens
	accessToken, err := r.encryptor.Decrypt(accessTokenEncrypted)
	if err != nil {
		return nil, err
	}

	var refreshToken string
	if refreshTokenEncrypted.Valid {
		refreshToken, err = r.encryptor.Decrypt(refreshTokenEncrypted.String)
		if err != nil {
			return nil, err
		}
	}

	var expiresAt *time.Time
	if tokenExpiresAt.Valid {
		expiresAt = &tokenExpiresAt.Time
	}

	var username string
	if providerUsername.Valid {
		username = providerUsername.String
	}

	return domain.ReconstructOAuthConnection(
		id, userID, domain.Provider(provider), providerUserID, username,
		accessToken, refreshToken, expiresAt, scopes,
		createdAt, updatedAt,
	), nil
}
