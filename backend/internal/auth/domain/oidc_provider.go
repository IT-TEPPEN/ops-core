package domain

import (
	"context"

	"golang.org/x/oauth2"
)

// OIDCProvider defines the interface for OIDC authentication providers
type OIDCProvider interface {
	// GetAuthURL returns the authorization URL for the provider
	GetAuthURL(state string) string

	// ExchangeToken exchanges authorization code for tokens
	ExchangeToken(ctx context.Context, code string) (*oauth2.Token, error)

	// VerifyIDToken verifies the ID token and returns user info
	VerifyIDToken(ctx context.Context, token *oauth2.Token) (*OIDCUserInfo, error)
}

// ProviderConfig holds configuration for an OIDC provider
type ProviderConfig struct {
	Provider     string
	ClientID     string
	ClientSecret string
	RedirectURL  string
}
