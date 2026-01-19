package usecase

import "context"

// RevokeAllUserSessions revokes all active sessions for a specific user
type RevokeAllUserSessions interface {
	// Execute revokes all active sessions belonging to the specified user
	// Returns the number of sessions that were revoked
	Execute(ctx context.Context, userID string) (int, error)
}
