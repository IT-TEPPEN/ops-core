package usecase

import "context"

// UnlinkIdentity removes an identity link from a user
type UnlinkIdentity interface {
	// Execute removes the specified identity from the user
	// Returns error if:
	// - The identity is the last one (user must have at least one identity)
	// - The identity is the primary identity (must set another as primary first)
	// - The identity does not belong to the user
	Execute(ctx context.Context, userID string, identityID string) error
}
