package repository

import (
	"context"

	"opscore/backend/internal/user/domain/entity"
	"opscore/backend/internal/user/domain/value_object"
)

// UserRepository defines the interface for user data persistence operations
type UserRepository interface {
	// Save persists a new user or updates an existing one
	Save(ctx context.Context, user entity.User) error
	// FindByID retrieves a user by its ID. Returns nil if not found.
	FindByID(ctx context.Context, id value_object.UserID) (entity.User, error)
	// FindByEmail retrieves a user by its email. Returns nil if not found.
	FindByEmail(ctx context.Context, email value_object.Email) (entity.User, error)
	// FindAll retrieves all users
	FindAll(ctx context.Context) ([]entity.User, error)
	// Update updates an existing user
	Update(ctx context.Context, user entity.User) error
	// Delete removes a user by its ID
	Delete(ctx context.Context, id value_object.UserID) error
}
