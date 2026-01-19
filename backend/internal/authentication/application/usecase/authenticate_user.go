package usecase

import (
	"context"

	"opscore/backend/internal/authentication/application/dto"
)

// AuthenticateUser authenticates a user using information from a trusted external authentication provider
type AuthenticateUser interface {
	// Execute authenticates a user with provider information
	// If the identity is not found, it creates a new user and identity
	// If the identity exists, it updates the user's last login time and identity information
	Execute(ctx context.Context, providerInfo dto.ProviderUserInfo) (*dto.AuthenticationResult, error)
}
