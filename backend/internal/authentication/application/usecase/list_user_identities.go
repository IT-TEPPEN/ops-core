package usecase

import (
	"context"

	"opscore/backend/internal/authentication/application/dto"
)

// ListUserIdentities retrieves all identities linked to a user
type ListUserIdentities interface {
	// Execute returns all identities (external provider links) for the specified user
	Execute(ctx context.Context, userID string) ([]dto.IdentityInfo, error)
}
