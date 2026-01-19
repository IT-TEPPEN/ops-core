package usecase

import (
	"context"

	"opscore/backend/internal/authentication/application/dto"
)

// CreateSession creates a new authentication session for a user
type CreateSession interface {
	// Execute creates a new session for the specified user
	// Returns the session identifier that can be used for subsequent requests
	Execute(ctx context.Context, userID string, sessionType dto.SessionType) (*dto.SessionInfo, error)
}
