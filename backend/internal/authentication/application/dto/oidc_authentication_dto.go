package dto

// GetAuthorizationURLInput represents input for getting OAuth authorization URL.
type GetAuthorizationURLInput struct {
	// Provider is the OAuth provider identifier (e.g., "google", "github")
	Provider string
}

// GetAuthorizationURLOutput represents output from getting OAuth authorization URL.
type GetAuthorizationURLOutput struct {
	// AuthorizationURL is the URL to redirect the user to
	AuthorizationURL string
	// State is the CSRF protection token (managed by infrastructure, returned for session storage)
	State string
}

// AuthenticateInput represents input for OIDC authentication.
// This combines the code exchange and user info retrieval into a single operation.
type AuthenticateInput struct {
	// Provider is the OIDC provider identifier (e.g., "google", "github")
	Provider string
	// Code is the authorization code received from the provider
	Code string
	// State is the state parameter received from the provider (for CSRF validation)
	State string
}
