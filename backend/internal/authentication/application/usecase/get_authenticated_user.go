package usecase

import (
	"context"

	"opscore/backend/internal/authentication/application/dto"
)

// GetAuthenticatedUser retrieves the currently authenticated user from a session.
// This usecase validates the session and returns the user information.
type GetAuthenticatedUser interface {
	// Execute validates the session and returns the authenticated user's information.
	// The session must be valid, not expired, and not revoked.
	//
	// Parameters:
	//   - ctx context.Context: The context for cancellation and timeout
	//   - dto dto.SessionIDRequest: Request containing sessionID to validate
	//
	// Returns:
	//   - *dto.UserInfo: Authenticated user's information
	//   - error: Error if session is invalid or user retrieval fails
	//
	// Errors:
	//   - application_err.ErrInvalidInput: When sessionID format is invalid
	//   - domain_err.ErrNotFoundSession: When session with given ID does not exist
	//   - application_err.ErrSessionExpired: When session has expired
	//   - application_err.ErrSessionRevoked: When session has been revoked
	//   - domain_err.ErrNotFoundUser: When user associated with session does not exist
	//   - domain_err.ErrDataAccessFailure: When loading session or user data fails
	Execute(ctx context.Context, dto dto.SessionIDRequest) (*dto.UserInfo, error)
}
