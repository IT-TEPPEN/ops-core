package persistence

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"opscore/backend/internal/auth/domain"
)

// PostgresRefreshTokenRepository implements RefreshTokenRepository using PostgreSQL
type PostgresRefreshTokenRepository struct {
	db *pgxpool.Pool
}

// NewPostgresRefreshTokenRepository creates a new PostgresRefreshTokenRepository
func NewPostgresRefreshTokenRepository(db *pgxpool.Pool) *PostgresRefreshTokenRepository {
	return &PostgresRefreshTokenRepository{db: db}
}

// Create creates a new refresh token
func (r *PostgresRefreshTokenRepository) Create(ctx context.Context, token *domain.RefreshToken) error {
	query := `
		INSERT INTO refresh_tokens (id, user_id, token_hash, jti, expires_at, created_at, last_used_at, revoked_at, is_revoked, remember_me)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
	`

	_, err := r.db.Exec(ctx, query,
		token.ID,
		token.UserID,
		token.TokenHash,
		token.JTI,
		token.ExpiresAt,
		token.CreatedAt,
		token.LastUsedAt,
		token.RevokedAt,
		token.IsRevoked,
		token.RememberMe,
	)

	if err != nil {
		return fmt.Errorf("failed to create refresh token: %w", err)
	}

	return nil
}

// FindByTokenHash finds a refresh token by its hash
func (r *PostgresRefreshTokenRepository) FindByTokenHash(ctx context.Context, tokenHash string) (*domain.RefreshToken, error) {
	query := `
		SELECT id, user_id, token_hash, jti, expires_at, created_at, last_used_at, revoked_at, is_revoked, remember_me
		FROM refresh_tokens
		WHERE token_hash = $1
	`

	var token domain.RefreshToken
	err := r.db.QueryRow(ctx, query, tokenHash).Scan(
		&token.ID,
		&token.UserID,
		&token.TokenHash,
		&token.JTI,
		&token.ExpiresAt,
		&token.CreatedAt,
		&token.LastUsedAt,
		&token.RevokedAt,
		&token.IsRevoked,
		&token.RememberMe,
	)

	if err == pgx.ErrNoRows {
		return nil, fmt.Errorf("refresh token not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to query refresh token: %w", err)
	}

	return &token, nil
}

// FindByJTI finds a refresh token by its JWT ID
func (r *PostgresRefreshTokenRepository) FindByJTI(ctx context.Context, jti string) (*domain.RefreshToken, error) {
	query := `
		SELECT id, user_id, token_hash, jti, expires_at, created_at, last_used_at, revoked_at, is_revoked, remember_me
		FROM refresh_tokens
		WHERE jti = $1
	`

	var token domain.RefreshToken
	err := r.db.QueryRow(ctx, query, jti).Scan(
		&token.ID,
		&token.UserID,
		&token.TokenHash,
		&token.JTI,
		&token.ExpiresAt,
		&token.CreatedAt,
		&token.LastUsedAt,
		&token.RevokedAt,
		&token.IsRevoked,
		&token.RememberMe,
	)

	if err == pgx.ErrNoRows {
		return nil, fmt.Errorf("refresh token not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to query refresh token: %w", err)
	}

	return &token, nil
}

// Update updates a refresh token
func (r *PostgresRefreshTokenRepository) Update(ctx context.Context, token *domain.RefreshToken) error {
	query := `
		UPDATE refresh_tokens
		SET last_used_at = $1, revoked_at = $2, is_revoked = $3
		WHERE id = $4
	`

	_, err := r.db.Exec(ctx, query,
		token.LastUsedAt,
		token.RevokedAt,
		token.IsRevoked,
		token.ID,
	)

	if err != nil {
		return fmt.Errorf("failed to update refresh token: %w", err)
	}

	return nil
}

// RevokeByUserID revokes all tokens for a user
func (r *PostgresRefreshTokenRepository) RevokeByUserID(ctx context.Context, userID string) error {
	query := `
		UPDATE refresh_tokens
		SET is_revoked = true, revoked_at = NOW()
		WHERE user_id = $1 AND is_revoked = false
	`

	_, err := r.db.Exec(ctx, query, userID)
	if err != nil {
		return fmt.Errorf("failed to revoke refresh tokens: %w", err)
	}

	return nil
}

// DeleteExpired deletes expired tokens (cleanup job)
func (r *PostgresRefreshTokenRepository) DeleteExpired(ctx context.Context) error {
	query := `
		DELETE FROM refresh_tokens
		WHERE expires_at < NOW()
	`

	result, err := r.db.Exec(ctx, query)
	if err != nil {
		return fmt.Errorf("failed to delete expired refresh tokens: %w", err)
	}

	rowsAffected := result.RowsAffected()
	if rowsAffected > 0 {
		// Log for monitoring purposes
		fmt.Printf("Deleted %d expired refresh tokens\n", rowsAffected)
	}

	return nil
}
