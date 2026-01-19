package usecase

import (
	"context"

	"opscore/backend/internal/authentication/application/dto"
)

// LinkNewIdentity links a new external provider identity to an existing user
type LinkNewIdentity interface {
	// Execute creates a new identity link for the user
	// Returns error if the identity already exists or is linked to another user
	Execute(ctx context.Context, userID string, providerInfo dto.ProviderUserInfo) (*dto.IdentityInfo, error)
}
