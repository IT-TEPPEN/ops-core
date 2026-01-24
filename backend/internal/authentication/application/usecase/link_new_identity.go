package usecase

import (
	"context"

	"opscore/backend/internal/authentication/application/dto"
)

// LinkNewIdentity links a new external provider identity to an existing user.
// This allows users to sign in with multiple authentication providers.
type LinkNewIdentity interface {
	// Execute creates a new identity link for the user.
	// The provider identity must not already be linked to any user (including this one).
	//
	// Parameters:
	//   - ctx context.Context: The context for cancellation and timeout
	//   - dto dto.LinkIdentityRequest: Request containing userID and provider information
	//
	// Returns:
	//   - *dto.IdentityInfo: The newly created identity information
	//   - error: Error if identity linking fails
	//
	// Errors:
	//   - application_err.ErrInvalidInput: When userID or provider information is invalid
	//   - domain_err.ErrNotFoundUser: When specified user does not exist
	//   - domain_err.ErrDuplicateIdentity: When identity is already linked to this user
	//   - application_err.ErrIdentityAlreadyLinked: When identity is linked to another user
	//   - domain_err.ErrDataPersistFailure: When saving identity data fails
	Execute(ctx context.Context, dto dto.LinkIdentityRequest) (*dto.IdentityInfo, error)
}

type linkNewIdentityImpl struct{}

func NewLinkNewIdentity() LinkNewIdentity {
	return &linkNewIdentityImpl{}
}

func (l *linkNewIdentityImpl) Execute(ctx context.Context, dto dto.LinkIdentityRequest) (*dto.IdentityInfo, error) {
	return nil, nil
}
