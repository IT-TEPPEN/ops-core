package service

import (
	"context"

	"opscore/backend/internal/authentication/application/dto"
)

// OIDCAuthenticationService defines the interface for user authentication via external OIDC providers.
// This service handles ONLY authentication (not repository access), exchanging authorization codes
// for user information. After obtaining user info, provider tokens are discarded as OpScore
// manages its own session via JWT.
// Implementation is provider-specific and should be in the Infrastructure layer.
type OIDCAuthenticationService interface {
	// GetAuthorizationURL generates the URL to redirect the user to the provider's authorization page.
	// This is the first step of the OIDC authentication flow.
	//
	// Parameters:
	//   - ctx context.Context: The context for cancellation and timeout
	//   - input dto.GetAuthorizationURLInput: Input parameters including provider
	//
	// Returns:
	//   - *dto.GetAuthorizationURLOutput: The authorization URL and state for CSRF protection
	//   - error: Application error if URL generation fails
	//
	// Errors:
	//   - application_err.ErrInvalidInput: When provider is empty or invalid
	//   - application_err.ErrUnsupportedProvider: When provider is not supported (e.g., not google, github, gitlab, microsoft)
	//   - application_err.ErrUnexpected: When unexpected error occurs
	GetAuthorizationURL(ctx context.Context, input dto.GetAuthorizationURLInput) (*dto.GetAuthorizationURLOutput, error)

	// Authenticate performs complete OIDC authentication flow.
	// Exchanges authorization code for tokens, verifies them, retrieves user info,
	// then discards provider tokens (as they are not needed for our session management).
	//
	// Parameters:
	//   - ctx context.Context: The context for cancellation and timeout
	//   - input dto.AuthenticateInput: Input parameters including provider, code, and state
	//
	// Returns:
	//   - *dto.ProviderUserInfo: Verified user information from the provider
	//   - error: Application error if authentication fails
	//
	// Errors:
	//   - application_err.ErrInvalidInput: When provider, code, or state is empty or invalid
	//   - application_err.ErrUnsupportedProvider: When provider is not supported
	//   - application_err.ErrAuthenticationFailed: When authentication fails for any of the following reasons:
	//     * Authorization code exchange fails (invalid code, expired code, etc.)
	//     * State validation fails (possible CSRF attack)
	//     * ID token verification fails (invalid signature, expired token, wrong audience, etc.)
	//     * UserInfo endpoint fails to return user data
	//   - application_err.ErrUnexpected: When unexpected error occurs
	Authenticate(ctx context.Context, input dto.AuthenticateInput) (*dto.ProviderUserInfo, error)
}
