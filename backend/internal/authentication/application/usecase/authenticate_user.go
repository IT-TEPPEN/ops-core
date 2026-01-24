package usecase

import (
	"context"

	"opscore/backend/internal/authentication/application/dto"
)

// AuthenticateUser authenticates a user using information from a trusted external authentication provider.
// This usecase handles both new user registration and existing user authentication.
type AuthenticateUser interface {
	// Execute authenticates a user with provider information.
	// If the identity is not found, it creates a new user and identity.
	// If the identity exists, it updates the user's last login time and identity information.
	//
	// Parameters:
	//   - ctx context.Context: The context for cancellation and timeout
	//   - dto dto.ProviderUserInfo: User information from the external provider
	//
	// Returns:
	//   - *dto.AuthenticationResult: Authentication result with user and identity information
	//   - error: Error if authentication fails
	//
	// Errors:
	//   - application_err.ErrInvalidInput: When provider information is invalid or incomplete
	//   - application_err.ErrAuthenticationFailed: When authentication process fails
	//   - domain_err.ErrDataPersistFailure: When saving user/identity data fails
	//   - domain_err.ErrDataConflict: When unique constraint is violated (e.g., email already exists)
	//   - application_err.ErrUnexpected: When unexpected error occurs
	Execute(ctx context.Context, dto dto.ProviderUserInfo) (*dto.AuthenticationResult, error)
}

type authenticateUserImpl struct{}

func NewAuthenticateUser() AuthenticateUser {
	return &authenticateUserImpl{}
}

func (a *authenticateUserImpl) Execute(ctx context.Context, dto dto.ProviderUserInfo) (*dto.AuthenticationResult, error) {
	return nil, nil
}
