package usecase

import (
	"context"

	"opscore/backend/internal/authentication/application/dto"
)

// ListUserIdentities retrieves all identities linked to a user.
// This allows users to see which authentication providers they have connected.
type ListUserIdentities interface {
	// Execute returns all identities (external provider links) for the specified user.
	//
	// Parameters:
	//   - ctx context.Context: The context for cancellation and timeout
	//   - userID string: The user ID whose identities to retrieve
	//
	// Returns:
	//   - []dto.IdentityInfo: List of identities (may be empty)
	//   - error: Error if retrieval fails
	//
	// Errors:
	//   - application_err.ErrInvalidInput: When userID format is invalid
	//   - domain_err.ErrNotFoundUser: When specified user does not exist
	//   - domain_err.ErrDataAccessFailure: When loading identity data fails
	Execute(ctx context.Context, userID string) ([]dto.IdentityInfo, error)
}
