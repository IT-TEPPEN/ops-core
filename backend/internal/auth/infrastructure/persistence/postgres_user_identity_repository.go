package persistence

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"opscore/backend/internal/auth/domain"
)

// PostgresUserIdentityRepository implements UserIdentityRepository using PostgreSQL
type PostgresUserIdentityRepository struct {
	db *pgxpool.Pool
}

// NewPostgresUserIdentityRepository creates a new PostgresUserIdentityRepository
func NewPostgresUserIdentityRepository(db *pgxpool.Pool) *PostgresUserIdentityRepository {
	return &PostgresUserIdentityRepository{db: db}
}

// FindByProviderIdentity finds a user identity by provider and provider user ID
func (r *PostgresUserIdentityRepository) FindByProviderIdentity(ctx context.Context, provider, providerUserID string) (*domain.UserIdentity, error) {
	query := `
		SELECT id, user_id, provider, provider_user_id, email, name, picture_url, created_at, updated_at, last_used_at
		FROM user_identities
		WHERE provider = $1 AND provider_user_id = $2
	`

	var identity domain.UserIdentity
	err := r.db.QueryRow(ctx, query, provider, providerUserID).Scan(
		&identity.ID,
		&identity.UserID,
		&identity.Provider,
		&identity.ProviderUserID,
		&identity.Email,
		&identity.Name,
		&identity.PictureURL,
		&identity.CreatedAt,
		&identity.UpdatedAt,
		&identity.LastUsedAt,
	)

	if err == pgx.ErrNoRows {
		return nil, fmt.Errorf("user identity not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to query user identity: %w", err)
	}

	return &identity, nil
}

// FindByUserID finds all identities for a user
func (r *PostgresUserIdentityRepository) FindByUserID(ctx context.Context, userID string) ([]*domain.UserIdentity, error) {
	query := `
		SELECT id, user_id, provider, provider_user_id, email, name, picture_url, created_at, updated_at, last_used_at
		FROM user_identities
		WHERE user_id = $1
		ORDER BY last_used_at DESC
	`

	rows, err := r.db.Query(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to query user identities: %w", err)
	}
	defer rows.Close()

	var identities []*domain.UserIdentity
	for rows.Next() {
		var identity domain.UserIdentity
		err := rows.Scan(
			&identity.ID,
			&identity.UserID,
			&identity.Provider,
			&identity.ProviderUserID,
			&identity.Email,
			&identity.Name,
			&identity.PictureURL,
			&identity.CreatedAt,
			&identity.UpdatedAt,
			&identity.LastUsedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan user identity: %w", err)
		}
		identities = append(identities, &identity)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating user identities: %w", err)
	}

	return identities, nil
}

// Create creates a new user identity
func (r *PostgresUserIdentityRepository) Create(ctx context.Context, identity *domain.UserIdentity) error {
	query := `
		INSERT INTO user_identities (id, user_id, provider, provider_user_id, email, name, picture_url, created_at, updated_at, last_used_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
	`

	_, err := r.db.Exec(ctx, query,
		identity.ID,
		identity.UserID,
		identity.Provider,
		identity.ProviderUserID,
		identity.Email,
		identity.Name,
		identity.PictureURL,
		identity.CreatedAt,
		identity.UpdatedAt,
		identity.LastUsedAt,
	)

	if err != nil {
		return fmt.Errorf("failed to create user identity: %w", err)
	}

	return nil
}

// Update updates an existing user identity
func (r *PostgresUserIdentityRepository) Update(ctx context.Context, identity *domain.UserIdentity) error {
	query := `
		UPDATE user_identities
		SET email = $1, name = $2, picture_url = $3, updated_at = $4, last_used_at = $5
		WHERE id = $6
	`

	result, err := r.db.Exec(ctx, query,
		identity.Email,
		identity.Name,
		identity.PictureURL,
		identity.UpdatedAt,
		identity.LastUsedAt,
		identity.ID,
	)

	if err != nil {
		return fmt.Errorf("failed to update user identity: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("user identity not found")
	}

	return nil
}

// Delete deletes a user identity
func (r *PostgresUserIdentityRepository) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM user_identities WHERE id = $1`

	result, err := r.db.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete user identity: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("user identity not found")
	}

	return nil
}
