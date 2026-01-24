package usecase

import (
	"context"

	"opscore/backend/internal/authentication/application/dto"
)

// RefreshSession refreshes an existing session to extend its validity.
// This usecase updates the session's last used timestamp without changing its expiration.
type RefreshSession interface {
	// Execute refreshes the session and returns updated session information.
	// The session must be valid, not expired, and not revoked.
	//
	// Parameters:
	//   - ctx context.Context: The context for cancellation and timeout
	//   - sessionID string: The session ID to refresh
	//
	// Returns:
	//   - *dto.SessionInfo: Updated session information
	//   - error: Error if session refresh fails
	//
	// Errors:
	//   - application_err.ErrInvalidInput: When sessionID format is invalid
	//   - domain_err.ErrNotFoundSession: When session with given ID does not exist
	//   - application_err.ErrSessionExpired: When session has expired
	//   - application_err.ErrSessionRevoked: When session has been revoked
	//   - domain_err.ErrDataPersistFailure: When updating session data fails
	Execute(ctx context.Context, sessionID string) (*dto.SessionInfo, error)
}
