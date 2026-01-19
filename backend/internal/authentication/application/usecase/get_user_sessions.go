package usecase

import (
	"context"

	"opscore/backend/internal/authentication/application/dto"
)

// GetUserSessions retrieves all active sessions for a user
type GetUserSessions interface {
	// Execute returns all active (non-revoked) sessions for the specified user
	Execute(ctx context.Context, userID string) ([]dto.SessionInfo, error)
}
