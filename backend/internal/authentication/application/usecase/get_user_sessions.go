package usecase

import (
	"context"

	"opscore/backend/internal/authentication/application/dto"
)

// GetUserSessions retrieves all active sessions for a user.
// This allows users to see where they are currently logged in.
type GetUserSessions interface {
	// Execute returns all active (non-revoked, non-expired) sessions for the specified user.
	//
	// Parameters:
	//   - ctx context.Context: The context for cancellation and timeout
	//   - userID string: The user ID whose sessions to retrieve
	//
	// Returns:
	//   - []dto.SessionInfo: List of active sessions (may be empty)
	//   - error: Error if retrieval fails
	//
	// Errors:
	//   - application_err.ErrInvalidInput: When userID format is invalid
	//   - domain_err.ErrNotFoundUser: When specified user does not exist
	//   - domain_err.ErrDataAccessFailure: When loading session data fails
	Execute(ctx context.Context, userID string) ([]dto.SessionInfo, error)
}
