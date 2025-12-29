package oauth

import (
	"context"
	"fmt"

	oauthservice "opscore/backend/internal/oauth/application/service"
	oauthdomain "opscore/backend/internal/oauth/domain"
)

// OAuthTokenProviderAdapter adapts OAuthService to OAuthTokenProvider interface
type OAuthTokenProviderAdapter struct {
	oauthService *oauthservice.OAuthService
}

// NewOAuthTokenProviderAdapter creates a new adapter
func NewOAuthTokenProviderAdapter(oauthService *oauthservice.OAuthService) *OAuthTokenProviderAdapter {
	return &OAuthTokenProviderAdapter{
		oauthService: oauthService,
	}
}

// GetAccessTokenForProvider retrieves an OAuth access token for the specified provider
// Deprecated: Use GetAccessTokenForRepoURL instead for better multi-instance support
func (a *OAuthTokenProviderAdapter) GetAccessTokenForProvider(ctx context.Context, userID string, providerName string) (string, error) {
	// Convert provider name to domain.Provider
	var provider oauthdomain.Provider
	switch providerName {
	case "github":
		provider = oauthdomain.ProviderGitHub
	case "gitlab":
		provider = oauthdomain.ProviderGitLab
	default:
		return "", fmt.Errorf("unsupported provider: %s", providerName)
	}

	return a.oauthService.GetAccessToken(ctx, userID, provider)
}

// GetAccessTokenForRepoURL retrieves an OAuth access token for a repository URL
// This method supports multiple self-hosted GitLab instances
func (a *OAuthTokenProviderAdapter) GetAccessTokenForRepoURL(ctx context.Context, userID string, repoURL string) (string, error) {
	return a.oauthService.GetAccessTokenForRepoURL(ctx, userID, repoURL)
}
