package usecase

import (
	"context"

	"opscore/backend/internal/authentication/application/dto"
)

// IssueSessionFromProviderUser issues an authentication session for a user authenticated by an external provider.
// This usecase orchestrates user registration/update and session creation based on provider authentication.
type IssueSessionFromProviderUser interface {
	// Execute issues an authentication session from external provider user information.
	//
	// This usecase:
	//   1. Validates provider user information
	//   2. Registers a new user or updates an existing user based on provider identity
	//   3. Issues an authentication session (access token and refresh token)
	//
	// Parameters:
	//   - ctx: Context for cancellation and timeout
	//   - providerUserInfo: User information from the external authentication provider
	//
	// Returns:
	//   - SessionResult containing access token, refresh token, and user information
	//   - Error if validation, user management, or session issuance fails
	//
	// Errors:
	//   - ErrInvalidInput: When provider information is invalid or domain rules are violated (wraps domain errors with WithParent)
	//   - ErrDataAccessFailure: When user data cannot be read from the database
	//   - ErrDataPersistFailure: When user data cannot be saved to the database
	//   - ErrDataConflict: When unique constraints are violated (e.g., email already exists)
	//   - ErrAuthenticationFailed: When session cannot be issued (JWT/RefreshToken generation failure)
	//   - ErrUnexpected: When an unexpected error occurs
	Execute(ctx context.Context, providerUserInfo dto.ProviderUserInfo) (*dto.SessionResult, error)
}

type issueSessionFromProviderUserImpl struct{}

func NewIssueSessionFromProviderUser() IssueSessionFromProviderUser {
	return &issueSessionFromProviderUserImpl{}
}

func (i *issueSessionFromProviderUserImpl) Execute(ctx context.Context, dto dto.ProviderUserInfo) (*dto.SessionResult, error) {
	return nil, nil
}
