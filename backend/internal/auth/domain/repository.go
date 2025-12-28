package domain

import "context"

// UserRepository defines the interface for user persistence
type UserRepository interface {
	// FindByID finds a user by their OpScore ID
	FindByID(ctx context.Context, id string) (*User, error)

	// FindByEmail finds a user by their email address
	FindByEmail(ctx context.Context, email string) (*User, error)

	// Create creates a new user
	Create(ctx context.Context, user *User) error

	// Update updates an existing user
	Update(ctx context.Context, user *User) error
}

// UserIdentityRepository defines the interface for user identity persistence
type UserIdentityRepository interface {
	// FindByProviderIdentity finds a user identity by provider and provider user ID
	FindByProviderIdentity(ctx context.Context, provider, providerUserID string) (*UserIdentity, error)

	// FindByUserID finds all identities for a user
	FindByUserID(ctx context.Context, userID string) ([]*UserIdentity, error)

	// Create creates a new user identity
	Create(ctx context.Context, identity *UserIdentity) error

	// Update updates an existing user identity
	Update(ctx context.Context, identity *UserIdentity) error

	// Delete deletes a user identity
	Delete(ctx context.Context, id string) error
}
