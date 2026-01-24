package usecase

import (
	"context"

	"opscore/backend/internal/authentication/application/dto"
)

// RevokeAllUserSessions revokes all active sessions for a specific user.
// This is used for security events (e.g., password change, account compromise).
type RevokeAllUserSessions interface {
	// Execute revokes all active sessions belonging to the specified user.
	// This effectively logs the user out from all devices.
	//
	// Parameters:
	//   - ctx context.Context: The context for cancellation and timeout
	//   - dto dto.UserIDRequest: Request containing userID whose sessions to revoke
	//
	// Returns:
	//   - int: Number of sessions that were revoked
	//   - error: Error if revocation fails
	//
	// Errors:
	//   - application_err.ErrInvalidInput: When userID format is invalid
	//   - domain_err.ErrNotFoundUser: When specified user does not exist
	//   - domain_err.ErrDataPersistFailure: When revoking sessions fails
	Execute(ctx context.Context, dto dto.UserIDRequest) (int, error)
}
