package persistence

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"opscore/backend/internal/auth/domain"
)

// PostgresUserRepository implements UserRepository using PostgreSQL
type PostgresUserRepository struct {
	db *pgxpool.Pool
}

// NewPostgresUserRepository creates a new PostgresUserRepository
func NewPostgresUserRepository(db *pgxpool.Pool) *PostgresUserRepository {
	return &PostgresUserRepository{db: db}
}

// FindByEmail finds a user by their email address
func (r *PostgresUserRepository) FindByEmail(ctx context.Context, email string) (*domain.User, error) {
	query := `
		SELECT id, primary_email, display_name, picture_url, last_login_at, created_at, updated_at
		FROM users
		WHERE primary_email = $1
	`

	var user domain.User
	err := r.db.QueryRow(ctx, query, email).Scan(
		&user.ID,
		&user.PrimaryEmail,
		&user.DisplayName,
		&user.PictureURL,
		&user.LastLoginAt,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err == pgx.ErrNoRows {
		return nil, fmt.Errorf("user not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to query user: %w", err)
	}

	return &user, nil
}

// FindByID finds a user by their ID
func (r *PostgresUserRepository) FindByID(ctx context.Context, id string) (*domain.User, error) {
	query := `
		SELECT id, primary_email, display_name, picture_url, last_login_at, created_at, updated_at
		FROM users
		WHERE id = $1
	`

	var user domain.User
	err := r.db.QueryRow(ctx, query, id).Scan(
		&user.ID,
		&user.PrimaryEmail,
		&user.DisplayName,
		&user.PictureURL,
		&user.LastLoginAt,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err == pgx.ErrNoRows {
		return nil, fmt.Errorf("user not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to query user: %w", err)
	}

	return &user, nil
}

// Create creates a new user
func (r *PostgresUserRepository) Create(ctx context.Context, user *domain.User) error {
	query := `
		INSERT INTO users (id, primary_email, display_name, picture_url, last_login_at, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`

	_, err := r.db.Exec(ctx, query,
		user.ID,
		user.PrimaryEmail,
		user.DisplayName,
		user.PictureURL,
		user.LastLoginAt,
		user.CreatedAt,
		user.UpdatedAt,
	)

	if err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}

	return nil
}

// Update updates an existing user
func (r *PostgresUserRepository) Update(ctx context.Context, user *domain.User) error {
	query := `
		UPDATE users
		SET primary_email = $2, display_name = $3, picture_url = $4, last_login_at = $5, updated_at = $6
		WHERE id = $1
	`

	result, err := r.db.Exec(ctx, query,
		user.ID,
		user.PrimaryEmail,
		user.DisplayName,
		user.PictureURL,
		user.LastLoginAt,
		user.UpdatedAt,
	)

	if err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("user not found")
	}

	return nil
}
