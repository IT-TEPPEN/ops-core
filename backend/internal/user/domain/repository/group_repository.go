package repository

import (
	"context"

	"opscore/backend/internal/user/domain/entity"
	"opscore/backend/internal/user/domain/value_object"
)

// GroupRepository defines the interface for group data persistence operations
type GroupRepository interface {
	// Save persists a new group or updates an existing one
	Save(ctx context.Context, group entity.Group) error
	// FindByID retrieves a group by its ID. Returns nil if not found.
	FindByID(ctx context.Context, id value_object.GroupID) (entity.Group, error)
	// FindByMemberID retrieves all groups that contain the specified user as a member
	FindByMemberID(ctx context.Context, userID value_object.UserID) ([]entity.Group, error)
	// FindAll retrieves all groups
	FindAll(ctx context.Context) ([]entity.Group, error)
	// Update updates an existing group
	Update(ctx context.Context, group entity.Group) error
	// Delete removes a group by its ID
	Delete(ctx context.Context, id value_object.GroupID) error
}
