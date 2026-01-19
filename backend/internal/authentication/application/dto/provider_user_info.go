package dto

// ProviderUserInfo represents user information provided by an external authentication provider
type ProviderUserInfo struct {
	// Provider is the identifier of the authentication provider (e.g., "google", "github")
	Provider string
	// Subject is the provider's unique user identifier (sub claim in OIDC)
	Subject string
	// Email is the user's email address
	Email string
	// EmailVerified indicates if the email has been verified by the provider
	EmailVerified bool
	// Name is the user's full name
	Name string
	// Picture is the URL of the user's profile picture
	Picture string
}
