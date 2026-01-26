package usecase

import (
	"context"

	"opscore/backend/internal/authentication/application/dto"
)

// IssueSessionFromProviderUser issues an authentication session (JWT tokens) from provider user information.
// This usecase handles user registration/update and session token generation after external provider authentication.
type IssueSessionFromProviderUser interface {
	// Execute issues authentication session from provider user information.
	// If the user identity is not found, it creates a new user and identity.
	// If the user identity exists, it updates the user's last login time and identity information.
	// Finally, it generates and returns JWT access token and refresh token.
	//
	// Parameters:
	//   - ctx context.Context: The context for cancellation and timeout
	//   - dto dto.ProviderUserInfo: User information from the external provider
	//
	// Returns:
	//   - *dto.SessionResult: Session result with JWT tokens and user information
	//   - error: Error if session issuance fails
	//
	// Errors:
	//   - application_err.ErrInvalidInput: When provider information is invalid or incomplete
	//   - application_err.ErrAuthenticationFailed: When session issuance fails
	//   - domain_err.ErrDataPersistFailure: When saving user/identity data fails
	//   - domain_err.ErrDataConflict: When unique constraint is violated (e.g., email already exists)
	//   - application_err.ErrUnexpected: When unexpected error occurs
	Execute(ctx context.Context, dto dto.ProviderUserInfo) (*dto.SessionResult, error)
}

type issueSessionFromProviderUserImpl struct{}

func NewIssueSessionFromProviderUser() IssueSessionFromProviderUser {
	return &issueSessionFromProviderUserImpl{}
}

func (i *issueSessionFromProviderUserImpl) Execute(ctx context.Context, dto dto.ProviderUserInfo) (*dto.SessionResult, error) {
	return nil, nil
}
