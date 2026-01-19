package repository

import (
	"context"

	"opscore/backend/internal/authentication/domain/entity"
	"opscore/backend/internal/authentication/domain/value_object"
)

// UserRepository defines the interface for User aggregate persistence.
// This repository is responsible for managing the User aggregate root and its child Identity entities.
// All operations maintain aggregate consistency by processing domain events.
type UserRepository interface {
	// Save saves the user aggregate (handles all changes via domain events).
	// This includes creating/updating the user and all identities.
	// The implementation processes domain events from user.GetEvents() to determine what has changed.
	//
	// Parameters:
	//   - ctx context.Context: The context for cancellation and timeout
	//   - user *entity.User: The user aggregate to save
	//
	// Returns:
	//   - error: Domain error if save fails
	//
	// Errors:
	//   - domain_err.ErrDataPersistFailure: When saving data fails
	//   - domain_err.ErrDataConflict: When unique constraint is violated
	Save(ctx context.Context, user *entity.User) error

	// FindByID finds a user by ID (without identities loaded).
	// Use FindByIDWithIdentities if you need the full aggregate with identities.
	//
	// Parameters:
	//   - ctx context.Context: The context for cancellation and timeout
	//   - userID value_object.UserID: The user ID to find
	//
	// Returns:
	//   - *entity.User: The user entity (identities not loaded)
	//   - error: Domain error if operation fails
	//
	// Errors:
	//   - domain_err.ErrNotFoundUser: When user with given ID does not exist
	//   - domain_err.ErrDataAccessFailure: When loading data fails
	FindByID(ctx context.Context, userID value_object.UserID) (*entity.User, error)

	// FindByIDWithIdentities finds a user by ID with all identities loaded.
	// This returns the full User aggregate including all Identity entities.
	//
	// Parameters:
	//   - ctx context.Context: The context for cancellation and timeout
	//   - userID value_object.UserID: The user ID to find
	//
	// Returns:
	//   - *entity.User: The user entity with identities loaded
	//   - error: Domain error if operation fails
	//
	// Errors:
	//   - domain_err.ErrNotFoundUser: When user with given ID does not exist
	//   - domain_err.ErrDataAccessFailure: When loading data fails
	FindByIDWithIdentities(ctx context.Context, userID value_object.UserID) (*entity.User, error)

	// FindByProviderIdentity finds a user by external provider identity.
	// This is used to authenticate users via external providers (e.g., Google, GitHub).
	//
	// Parameters:
	//   - ctx context.Context: The context for cancellation and timeout
	//   - provider string: Provider name (e.g., "google", "github")
	//   - providerUserID string: User ID from the provider
	//
	// Returns:
	//   - *entity.User: The user entity with identities loaded
	//   - error: Domain error if operation fails
	//
	// Errors:
	//   - domain_err.ErrNotFoundUser: When no user with given provider identity exists
	//   - domain_err.ErrDataAccessFailure: When loading data fails
	FindByProviderIdentity(ctx context.Context, provider, providerUserID string) (*entity.User, error)

	// Delete deletes a user aggregate (cascades to all identities).
	// This removes the user and all associated identities from the system.
	//
	// Parameters:
	//   - ctx context.Context: The context for cancellation and timeout
	//   - userID value_object.UserID: The user ID to delete
	//
	// Returns:
	//   - error: Domain error if delete fails
	//
	// Errors:
	//   - domain_err.ErrNotFoundUser: When user with given ID does not exist
	//   - domain_err.ErrDataPersistFailure: When deleting data fails
	Delete(ctx context.Context, userID value_object.UserID) error
}
