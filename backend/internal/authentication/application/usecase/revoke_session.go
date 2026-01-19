package usecase

import "context"

// RevokeSession revokes a specific authentication session
type RevokeSession interface {
	// Execute marks the specified session as revoked
	// Returns error if the session does not exist
	Execute(ctx context.Context, sessionID string) error
}
