package usecase

import (
	"context"

	"opscore/backend/internal/authentication/application/dto"
)

// SetPrimaryIdentity sets which identity should be used as the user's primary identity.
// The primary identity's display name and picture are used as defaults for the user.
type SetPrimaryIdentity interface {
	// Execute sets the specified identity as the user's primary identity.
	// The previous primary identity will be automatically demoted.
	//
	// Parameters:
	//   - ctx context.Context: The context for cancellation and timeout
	//   - dto dto.SetPrimaryIdentityRequest: Request containing userID and identityID
	//
	// Returns:
	//   - error: Error if setting primary identity fails
	//
	// Errors:
	//   - application_err.ErrInvalidInput: When userID or identityID format is invalid
	//   - domain_err.ErrNotFoundUser: When specified user does not exist
	//   - application_err.ErrResourceNotOwned: When identity does not belong to the user
	//   - domain_err.ErrDataPersistFailure: When saving changes fails
	Execute(ctx context.Context, dto dto.SetPrimaryIdentityRequest) error
}

type setPrimaryIdentityImpl struct{}

func NewSetPrimaryIdentity() SetPrimaryIdentity {
	return &setPrimaryIdentityImpl{}
}

func (s *setPrimaryIdentityImpl) Execute(ctx context.Context, dto dto.SetPrimaryIdentityRequest) error {
	return nil
}
