package usecase

import (
	"context"

	"opscore/backend/internal/authentication/application/dto"
)

// UnlinkIdentity removes an identity link from a user.
// This allows users to remove authentication providers they no longer want to use.
type UnlinkIdentity interface {
	// Execute removes the specified identity from the user.
	// The user must have at least one identity remaining after the operation.
	// Cannot remove the primary identity without setting another as primary first.
	//
	// Parameters:
	//   - ctx context.Context: The context for cancellation and timeout
	//   - dto dto.UnlinkIdentityRequest: Request containing userID and identityID
	//
	// Returns:
	//   - error: Error if identity unlinking fails
	//
	// Errors:
	//   - application_err.ErrInvalidInput: When userID or identityID format is invalid
	//   - domain_err.ErrNotFoundUser: When specified user does not exist
	//   - application_err.ErrResourceNotOwned: When identity does not belong to the user
	//   - domain_err.ErrCannotRemoveOnlyIdentity: When trying to remove the last identity
	//   - domain_err.ErrCannotUnlinkPrimaryIdentity: When trying to remove the primary identity
	//   - domain_err.ErrDataPersistFailure: When saving changes fails
	Execute(ctx context.Context, dto dto.UnlinkIdentityRequest) error
}
