package service

import (
	"context"
	"fmt"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	oauth2api "google.golang.org/api/oauth2/v2"
	"google.golang.org/api/option"

	"opscore/backend/internal/auth/domain"
)

// GoogleOAuthService handles Google OIDC authentication
type GoogleOAuthService struct {
	config *oauth2.Config
}

// NewGoogleOAuthService creates a new GoogleOAuthService
func NewGoogleOAuthService(clientID, clientSecret, redirectURL string) *GoogleOAuthService {
	return &GoogleOAuthService{
		config: &oauth2.Config{
			ClientID:     clientID,
			ClientSecret: clientSecret,
			RedirectURL:  redirectURL,
			Scopes: []string{
				"openid",
				"https://www.googleapis.com/auth/userinfo.email",
				"https://www.googleapis.com/auth/userinfo.profile",
			},
			Endpoint: google.Endpoint,
		},
	}
}

// GetAuthURL returns the Google OAuth authorization URL
func (s *GoogleOAuthService) GetAuthURL(state string) string {
	return s.config.AuthCodeURL(state, oauth2.AccessTypeOffline)
}

// ExchangeToken exchanges authorization code for tokens (implements domain.OIDCProvider)
func (s *GoogleOAuthService) ExchangeToken(ctx context.Context, code string) (*oauth2.Token, error) {
	token, err := s.config.Exchange(ctx, code)
	if err != nil {
		return nil, fmt.Errorf("failed to exchange code for token: %w", err)
	}
	return token, nil
}

// VerifyIDToken verifies the ID token and retrieves user information (implements domain.OIDCProvider)
func (s *GoogleOAuthService) VerifyIDToken(ctx context.Context, token *oauth2.Token) (*domain.OIDCUserInfo, error) {
	// Create OAuth2 service client
	oauth2Service, err := oauth2api.NewService(ctx, option.WithTokenSource(s.config.TokenSource(ctx, token)))
	if err != nil {
		return nil, fmt.Errorf("failed to create oauth2 service: %w", err)
	}

	// Get user info
	userInfo, err := oauth2Service.Userinfo.Get().Do()
	if err != nil {
		return nil, fmt.Errorf("failed to get user info: %w", err)
	}

	emailVerified := false
	if userInfo.VerifiedEmail != nil {
		emailVerified = *userInfo.VerifiedEmail
	}

	return &domain.OIDCUserInfo{
		Provider:      "google",
		Subject:       userInfo.Id,
		Email:         userInfo.Email,
		EmailVerified: emailVerified,
		Name:          userInfo.Name,
		Picture:       userInfo.Picture,
	}, nil
}
