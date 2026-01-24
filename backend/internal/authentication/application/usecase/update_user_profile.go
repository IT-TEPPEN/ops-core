package usecase

import (
	"context"

	"opscore/backend/internal/authentication/application/dto"
)

// UpdateUserProfile updates the user's customizable profile information.
// This allows users to override their display name and picture from provider defaults.
type UpdateUserProfile interface {
	// Execute updates the user's profile.
	// Nil or empty string values mean "use the primary identity's value".
	//
	// Parameters:
	//   - ctx context.Context: The context for cancellation and timeout
	//   - dto dto.UserProfile: Profile data to update (userID must be set)
	//
	// Returns:
	//   - *dto.UserProfile: Updated user profile
	//   - error: Error if profile update fails
	//
	// Errors:
	//   - application_err.ErrInvalidInput: When userID is invalid or profile data is malformed
	//   - domain_err.ErrNotFoundUser: When specified user does not exist
	//   - domain_err.ErrDataPersistFailure: When saving profile changes fails
	Execute(ctx context.Context, dto dto.UserProfile) (*dto.UserProfile, error)
}
