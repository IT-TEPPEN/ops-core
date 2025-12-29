package persistence

import (
	"context"
	"encoding/json"
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
	// Encrypt access token
	encryptedAccessToken, err := r.encryptor.Encrypt(conn.AccessToken())
	if err != nil {
		return err
	}

	// Prepare metadata for encryption
	metadata := conn.Metadata()

	// Encrypt refresh token in metadata if present
	if metadata.RefreshTokenEncrypted != nil && *metadata.RefreshTokenEncrypted != "" {
		encrypted, err := r.encryptor.Encrypt(*metadata.RefreshTokenEncrypted)
		if err != nil {
			return err
		}
		metadata.RefreshTokenEncrypted = &encrypted
	}

	// Encrypt client credentials in metadata if present (for self-hosted GitLab)
	if metadata.ClientIDEncrypted != nil && *metadata.ClientIDEncrypted != "" {
		encrypted, err := r.encryptor.Encrypt(*metadata.ClientIDEncrypted)
		if err != nil {
			return err
		}
		metadata.ClientIDEncrypted = &encrypted
	}

	if metadata.ClientSecretEncrypted != nil && *metadata.ClientSecretEncrypted != "" {
		encrypted, err := r.encryptor.Encrypt(*metadata.ClientSecretEncrypted)
		if err != nil {
			return err
		}
		metadata.ClientSecretEncrypted = &encrypted
	}

	// Marshal metadata to JSONB
	metadataJSON, err := json.Marshal(metadata)
	if err != nil {
		return err
	}

	query := `
		INSERT INTO oauth_connections (
			id, user_id, provider, provider_host, provider_metadata,
			access_token_encrypted, created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		ON CONFLICT (user_id, provider, provider_host) DO UPDATE SET
			provider_metadata = EXCLUDED.provider_metadata,
			access_token_encrypted = EXCLUDED.access_token_encrypted,
			updated_at = EXCLUDED.updated_at
	`

	_, err = r.db.Exec(ctx, query,
		conn.ID(),
		conn.UserID(),
		string(conn.Provider()),
		conn.ProviderHost(),
		metadataJSON,
		encryptedAccessToken,
		conn.CreatedAt(),
		conn.UpdatedAt(),
	)

	return err
}

// FindByID finds an OAuth connection by ID
func (r *PostgresOAuthConnectionRepository) FindByID(ctx context.Context, id string) (*domain.OAuthConnection, error) {
	query := `
		SELECT id, user_id, provider, provider_host, provider_metadata,
			   access_token_encrypted, created_at, updated_at
		FROM oauth_connections
		WHERE id = $1
	`

	return r.scanConnection(ctx, query, id)
}

// FindByUserAndProvider finds an OAuth connection by user ID and provider
// Note: This will return the first match if multiple hosts exist for the same provider
func (r *PostgresOAuthConnectionRepository) FindByUserAndProvider(ctx context.Context, userID string, provider domain.Provider) (*domain.OAuthConnection, error) {
	query := `
		SELECT id, user_id, provider, provider_host, provider_metadata,
			   access_token_encrypted, created_at, updated_at
		FROM oauth_connections
		WHERE user_id = $1 AND provider = $2
		LIMIT 1
	`

	return r.scanConnection(ctx, query, userID, string(provider))
}

// FindByUserProviderAndHost finds an OAuth connection by user ID, provider, and host
func (r *PostgresOAuthConnectionRepository) FindByUserProviderAndHost(ctx context.Context, userID string, provider domain.Provider, providerHost string) (*domain.OAuthConnection, error) {
	query := `
		SELECT id, user_id, provider, provider_host, provider_metadata,
			   access_token_encrypted, created_at, updated_at
		FROM oauth_connections
		WHERE user_id = $1 AND provider = $2 AND provider_host = $3
	`

	return r.scanConnection(ctx, query, userID, string(provider), providerHost)
}

// FindByUser finds all OAuth connections for a user
func (r *PostgresOAuthConnectionRepository) FindByUser(ctx context.Context, userID string) ([]*domain.OAuthConnection, error) {
	query := `
		SELECT id, user_id, provider, provider_host, provider_metadata,
			   access_token_encrypted, created_at, updated_at
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
		id                   string
		userID               string
		provider             string
		providerHost         string
		providerMetadataJSON []byte
		accessTokenEncrypted string
		createdAt            time.Time
		updatedAt            time.Time
	)

	err := row.Scan(
		&id, &userID, &provider, &providerHost, &providerMetadataJSON,
		&accessTokenEncrypted, &createdAt, &updatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return r.buildConnection(
		id, userID, provider, providerHost, providerMetadataJSON,
		accessTokenEncrypted, createdAt, updatedAt,
	)
}

func (r *PostgresOAuthConnectionRepository) scanRow(rows pgx.Rows) (*domain.OAuthConnection, error) {
	var (
		id                   string
		userID               string
		provider             string
		providerHost         string
		providerMetadataJSON []byte
		accessTokenEncrypted string
		createdAt            time.Time
		updatedAt            time.Time
	)

	err := rows.Scan(
		&id, &userID, &provider, &providerHost, &providerMetadataJSON,
		&accessTokenEncrypted, &createdAt, &updatedAt,
	)
	if err != nil {
		return nil, err
	}

	return r.buildConnection(
		id, userID, provider, providerHost, providerMetadataJSON,
		accessTokenEncrypted, createdAt, updatedAt,
	)
}

func (r *PostgresOAuthConnectionRepository) buildConnection(
	id, userID, provider, providerHost string,
	providerMetadataJSON []byte,
	accessTokenEncrypted string,
	createdAt, updatedAt time.Time,
) (*domain.OAuthConnection, error) {
	// Unmarshal provider metadata
	var metadata domain.ProviderMetadata
	if err := json.Unmarshal(providerMetadataJSON, &metadata); err != nil {
		return nil, err
	}

	// Decrypt access token
	accessToken, err := r.encryptor.Decrypt(accessTokenEncrypted)
	if err != nil {
		return nil, err
	}

	// Decrypt refresh token in metadata if present
	if metadata.RefreshTokenEncrypted != nil && *metadata.RefreshTokenEncrypted != "" {
		decrypted, err := r.encryptor.Decrypt(*metadata.RefreshTokenEncrypted)
		if err != nil {
			return nil, err
		}
		metadata.RefreshTokenEncrypted = &decrypted
	}

	// Decrypt client credentials in metadata if present
	if metadata.ClientIDEncrypted != nil && *metadata.ClientIDEncrypted != "" {
		decrypted, err := r.encryptor.Decrypt(*metadata.ClientIDEncrypted)
		if err != nil {
			return nil, err
		}
		metadata.ClientIDEncrypted = &decrypted
	}

	if metadata.ClientSecretEncrypted != nil && *metadata.ClientSecretEncrypted != "" {
		decrypted, err := r.encryptor.Decrypt(*metadata.ClientSecretEncrypted)
		if err != nil {
			return nil, err
		}
		metadata.ClientSecretEncrypted = &decrypted
	}

	return domain.ReconstructOAuthConnection(
		id, userID, domain.Provider(provider), providerHost, metadata,
		accessToken, createdAt, updatedAt,
	), nil
}
