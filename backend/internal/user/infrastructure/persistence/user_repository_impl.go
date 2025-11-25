package persistence

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"opscore/backend/internal/user/domain/entity"
	"opscore/backend/internal/user/domain/repository"
	"opscore/backend/internal/user/domain/value_object"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// UserRepositoryImpl is a PostgreSQL implementation of the UserRepository interface
type UserRepositoryImpl struct {
	db *pgxpool.Pool
}

// NewUserRepositoryImpl creates a new UserRepositoryImpl
func NewUserRepositoryImpl(db *pgxpool.Pool) repository.UserRepository {
	return &UserRepositoryImpl{db: db}
}

// Save persists a new user or updates an existing one
func (r *UserRepositoryImpl) Save(ctx context.Context, user entity.User) error {
	query := `
		INSERT INTO users (id, name, email, role, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6)
		ON CONFLICT (id) DO UPDATE SET
			name = EXCLUDED.name,
			email = EXCLUDED.email,
			role = EXCLUDED.role,
			updated_at = EXCLUDED.updated_at;
	`
	_, err := r.db.Exec(ctx, query,
		user.ID().String(),
		user.Name(),
		user.Email().String(),
		user.Role().String(),
		user.CreatedAt(),
		user.UpdatedAt(),
	)
	if err != nil {
		return fmt.Errorf("failed to save user: %w", err)
	}

	// Save user-group associations
	err = r.saveUserGroups(ctx, user.ID().String(), user.GroupIDs())
	if err != nil {
		return fmt.Errorf("failed to save user groups: %w", err)
	}

	return nil
}

// saveUserGroups saves the user-group associations
func (r *UserRepositoryImpl) saveUserGroups(ctx context.Context, userID string, groupIDs []value_object.GroupID) error {
	// Delete existing associations
	deleteQuery := `DELETE FROM user_groups WHERE user_id = $1;`
	_, err := r.db.Exec(ctx, deleteQuery, userID)
	if err != nil {
		return fmt.Errorf("failed to delete existing user groups: %w", err)
	}

	// Insert new associations
	if len(groupIDs) > 0 {
		insertQuery := `INSERT INTO user_groups (user_id, group_id) VALUES ($1, $2);`
		batch := &pgx.Batch{}
		for _, groupID := range groupIDs {
			batch.Queue(insertQuery, userID, groupID.String())
		}

		results := r.db.SendBatch(ctx, batch)
		defer results.Close()

		for range groupIDs {
			_, err := results.Exec()
			if err != nil {
				return fmt.Errorf("failed to insert user group: %w", err)
			}
		}
	}

	return nil
}

// FindByID retrieves a user by its ID
func (r *UserRepositoryImpl) FindByID(ctx context.Context, id value_object.UserID) (entity.User, error) {
	query := `
		SELECT id, name, email, role, created_at, updated_at
		FROM users
		WHERE id = $1;
	`

	var userID, name, email, role string
	var createdAt, updatedAt time.Time

	err := r.db.QueryRow(ctx, query, id.String()).Scan(
		&userID,
		&name,
		&email,
		&role,
		&createdAt,
		&updatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to find user by ID: %w", err)
	}

	// Get group IDs
	groupIDs, err := r.getUserGroupIDs(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user groups: %w", err)
	}

	return r.toDomainEntity(userID, name, email, role, groupIDs, createdAt, updatedAt)
}

// FindByEmail retrieves a user by its email
func (r *UserRepositoryImpl) FindByEmail(ctx context.Context, email value_object.Email) (entity.User, error) {
	query := `
		SELECT id, name, email, role, created_at, updated_at
		FROM users
		WHERE email = $1;
	`

	var userID, name, emailStr, role string
	var createdAt, updatedAt time.Time

	err := r.db.QueryRow(ctx, query, email.String()).Scan(
		&userID,
		&name,
		&emailStr,
		&role,
		&createdAt,
		&updatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to find user by email: %w", err)
	}

	// Get group IDs
	groupIDs, err := r.getUserGroupIDs(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user groups: %w", err)
	}

	return r.toDomainEntity(userID, name, emailStr, role, groupIDs, createdAt, updatedAt)
}

// FindAll retrieves all users
func (r *UserRepositoryImpl) FindAll(ctx context.Context) ([]entity.User, error) {
	query := `
		SELECT id, name, email, role, created_at, updated_at
		FROM users
		ORDER BY created_at DESC;
	`

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to query users: %w", err)
	}
	defer rows.Close()

	var users []entity.User
	for rows.Next() {
		var userID, name, email, role string
		var createdAt, updatedAt time.Time

		if err := rows.Scan(&userID, &name, &email, &role, &createdAt, &updatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan user row: %w", err)
		}

		// Get group IDs
		groupIDs, err := r.getUserGroupIDs(ctx, userID)
		if err != nil {
			return nil, fmt.Errorf("failed to get user groups: %w", err)
		}

		user, err := r.toDomainEntity(userID, name, email, role, groupIDs, createdAt, updatedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to convert to domain entity: %w", err)
		}

		users = append(users, user)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over user rows: %w", err)
	}

	return users, nil
}

// Update updates an existing user
func (r *UserRepositoryImpl) Update(ctx context.Context, user entity.User) error {
	query := `
		UPDATE users
		SET name = $1, email = $2, role = $3, updated_at = $4
		WHERE id = $5;
	`

	result, err := r.db.Exec(ctx, query,
		user.Name(),
		user.Email().String(),
		user.Role().String(),
		user.UpdatedAt(),
		user.ID().String(),
	)
	if err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("user with ID %s not found", user.ID().String())
	}

	// Update user-group associations
	err = r.saveUserGroups(ctx, user.ID().String(), user.GroupIDs())
	if err != nil {
		return fmt.Errorf("failed to save user groups: %w", err)
	}

	return nil
}

// Delete removes a user by its ID
func (r *UserRepositoryImpl) Delete(ctx context.Context, id value_object.UserID) error {
	// First delete user-group associations
	deleteGroupsQuery := `DELETE FROM user_groups WHERE user_id = $1;`
	_, err := r.db.Exec(ctx, deleteGroupsQuery, id.String())
	if err != nil {
		return fmt.Errorf("failed to delete user groups: %w", err)
	}

	// Then delete the user
	deleteQuery := `DELETE FROM users WHERE id = $1;`
	result, err := r.db.Exec(ctx, deleteQuery, id.String())
	if err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("user with ID %s not found", id.String())
	}

	return nil
}

// getUserGroupIDs retrieves the group IDs for a user
func (r *UserRepositoryImpl) getUserGroupIDs(ctx context.Context, userID string) ([]string, error) {
	query := `
		SELECT group_id
		FROM user_groups
		WHERE user_id = $1
		ORDER BY group_id;
	`

	rows, err := r.db.Query(ctx, query, userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return []string{}, nil
		}
		return nil, fmt.Errorf("failed to query user groups: %w", err)
	}
	defer rows.Close()

	var groupIDs []string
	for rows.Next() {
		var groupID string
		if err := rows.Scan(&groupID); err != nil {
			return nil, fmt.Errorf("failed to scan group ID: %w", err)
		}
		groupIDs = append(groupIDs, groupID)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over group ID rows: %w", err)
	}

	return groupIDs, nil
}

// toDomainEntity converts database data to a domain entity
func (r *UserRepositoryImpl) toDomainEntity(
	userID, name, email, role string,
	groupIDStrs []string,
	createdAt, updatedAt time.Time,
) (entity.User, error) {
	id, err := value_object.NewUserID(userID)
	if err != nil {
		return nil, fmt.Errorf("invalid user ID: %w", err)
	}

	emailVO, err := value_object.NewEmail(email)
	if err != nil {
		return nil, fmt.Errorf("invalid email: %w", err)
	}

	roleVO, err := value_object.NewRole(role)
	if err != nil {
		return nil, fmt.Errorf("invalid role: %w", err)
	}

	groupIDs := make([]value_object.GroupID, 0, len(groupIDStrs))
	for _, gidStr := range groupIDStrs {
		gid, err := value_object.NewGroupID(strings.TrimSpace(gidStr))
		if err != nil {
			return nil, fmt.Errorf("invalid group ID: %w", err)
		}
		groupIDs = append(groupIDs, gid)
	}

	return entity.ReconstructUser(id, name, emailVO, roleVO, groupIDs, createdAt, updatedAt), nil
}
