package usecase

import "context"

// RevokeSession revokes a specific authentication session.
// This usecase marks a session as revoked, preventing further use (logout).
type RevokeSession interface {
	// Execute marks the specified session as revoked.
	// A revoked session cannot be used for authentication or refresh.
	//
	// Parameters:
	//   - ctx context.Context: The context for cancellation and timeout
	//   - sessionID string: The session ID to revoke
	//
	// Returns:
	//   - error: Error if session revocation fails
	//
	// Errors:
	//   - application_err.ErrInvalidInput: When sessionID format is invalid
	//   - domain_err.ErrNotFoundSession: When session with given ID does not exist
	//   - domain_err.ErrDataPersistFailure: When updating session revocation status fails
	Execute(ctx context.Context, sessionID string) error
}
