package usecase

import (
	"context"

	"opscore/backend/internal/authentication/application/dto"
)

// RefreshSession refreshes an existing session to extend its validity
type RefreshSession interface {
	// Execute refreshes the session and returns updated session information
	// Returns error if the session is invalid, expired, or revoked
	Execute(ctx context.Context, sessionID string) (*dto.SessionInfo, error)
}
