package usecase

import (
	"context"

	"opscore/backend/internal/authentication/application/dto"
)

// CreateSession creates a new authentication session for a user.
// This usecase generates a refresh token and stores the session in the database.
type CreateSession interface {
	// Execute creates a new session for the specified user.
	// The session duration is determined by the sessionType (short: 7 days, long: 30 days).
	//
	// Parameters:
	//   - ctx context.Context: The context for cancellation and timeout
	//   - dto dto.CreateSessionRequest: Request containing userID and sessionType
	//
	// Returns:
	//   - *dto.SessionInfo: Created session information with session ID and expiration
	//   - error: Error if session creation fails
	//
	// Errors:
	//   - application_err.ErrInvalidInput: When userID is invalid or sessionType is unknown
	//   - domain_err.ErrNotFoundUser: When specified user does not exist
	//   - domain_err.ErrDataPersistFailure: When saving session data fails
	//   - application_err.ErrUnexpected: When session token generation fails
	Execute(ctx context.Context, dto dto.CreateSessionRequest) (*dto.SessionInfo, error)
}
