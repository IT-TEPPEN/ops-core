package usecase

import "context"

// SetPrimaryIdentity sets which identity should be used as the user's primary identity.
// The primary identity's display name and picture are used as defaults for the user.
type SetPrimaryIdentity interface {
	// Execute sets the specified identity as the user's primary identity.
	// The previous primary identity will be automatically demoted.
	//
	// Parameters:
	//   - ctx context.Context: The context for cancellation and timeout
	//   - userID string: The user ID who owns the identity
	//   - identityID string: The identity ID to set as primary
	//
	// Returns:
	//   - error: Error if setting primary identity fails
	//
	// Errors:
	//   - application_err.ErrInvalidInput: When userID or identityID format is invalid
	//   - domain_err.ErrNotFoundUser: When specified user does not exist
	//   - application_err.ErrResourceNotOwned: When identity does not belong to the user
	//   - domain_err.ErrDataPersistFailure: When saving changes fails
	Execute(ctx context.Context, userID string, identityID string) error
}
