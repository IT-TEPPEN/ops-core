package service

import (
	"context"

	"opscore/backend/internal/authentication/application/dto"
	"opscore/backend/internal/authentication/domain/value_object"
)

// SessionService handles session creation and management.
// This application service orchestrates JWT and refresh token generation
// for authenticated users.
type SessionService interface {
	// CreateSession creates a new authentication session for the user.
	//
	// This method:
	//   1. Generates a JWT access token with appropriate claims
	//   2. Creates a refresh token (Session entity) and persists it
	//   3. Returns both tokens along with user information
	//
	// Parameters:
	//   - ctx: Context for cancellation and timeout
	//   - userID: The user ID to create session for
	//   - rememberMe: Whether to create a long-lived session (30 days) or short-lived (7 days)
	//
	// Returns:
	//   - SessionResult containing access token, refresh token, and user information
	//   - Error if JWT generation, session creation, or persistence fails
	//
	// Errors:
	//   - application_err.ErrAuthenticationFailed: When JWT or refresh token generation fails
	//   - application_err.ErrDataPersistFailure: When session cannot be saved to database
	//   - application_err.ErrDataAccessFailure: When user information cannot be loaded
	//   - application_err.ErrUnexpected: When an unexpected error occurs
	CreateSession(ctx context.Context, userID value_object.UserID, rememberMe bool) (*dto.SessionResult, error)
}
