package usecase

import "context"

// SetPrimaryIdentity sets which identity should be used as the user's primary identity
type SetPrimaryIdentity interface {
	// Execute sets the specified identity as the user's primary identity
	// The primary identity is used for default display name and profile picture
	// Returns error if the identity does not belong to the user
	Execute(ctx context.Context, userID string, identityID string) error
}
