package usecase

import (
	"context"

	"opscore/backend/internal/authentication/application/dto"
)

// GetAuthenticatedUser retrieves the currently authenticated user from a session
type GetAuthenticatedUser interface {
	// Execute validates the session and returns the authenticated user's information
	// Returns error if the session is invalid, expired, or revoked
	Execute(ctx context.Context, sessionID string) (*dto.UserInfo, error)
}
