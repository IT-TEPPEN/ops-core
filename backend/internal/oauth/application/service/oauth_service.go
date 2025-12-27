package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"

	"opscore/backend/internal/oauth/domain"
)

// OAuthService handles OAuth2.0 token exchange
type OAuthService struct {
	httpClient *http.Client
}

// NewOAuthService creates a new OAuthService
func NewOAuthService() *OAuthService {
	return &OAuthService{
		httpClient: &http.Client{},
	}
}

// ExchangeToken exchanges authorization code for access token
func (s *OAuthService) ExchangeToken(req domain.OAuthCallbackRequest) (*domain.OAuthTokenResponse, error) {
	config, err := s.getOAuthConfig(req.Provider, req.GitLabURL, req.ClientID, req.ClientSecret)
	if err != nil {
		return nil, fmt.Errorf("failed to get OAuth config: %w", err)
	}

	tokenResp, err := s.exchangeCodeForToken(config, req.Code)
	if err != nil {
		return nil, fmt.Errorf("failed to exchange code for token: %w", err)
	}

	return tokenResp, nil
}

// getOAuthConfig returns OAuth configuration for the provider
func (s *OAuthService) getOAuthConfig(provider domain.Provider, gitlabURL, clientID, clientSecret string) (*domain.OAuthConfig, error) {
	redirectURI := os.Getenv("OAUTH_REDIRECT_URI")
	if redirectURI == "" {
		redirectURI = "http://localhost:5173/oauth/callback"
	}

	switch provider {
	case domain.ProviderGitHub:
		clientID := os.Getenv("GITHUB_CLIENT_ID")
		secret := os.Getenv("GITHUB_CLIENT_SECRET")
		if clientID == "" || secret == "" {
			return nil, fmt.Errorf("GitHub OAuth credentials not configured")
		}
		return &domain.OAuthConfig{
			ClientID:     clientID,
			ClientSecret: secret,
			TokenURL:     "https://github.com/login/oauth/access_token",
			RedirectURI:  redirectURI,
		}, nil

	case domain.ProviderGitLab:
		clientID := os.Getenv("GITLAB_CLIENT_ID")
		secret := os.Getenv("GITLAB_CLIENT_SECRET")
		if clientID == "" || secret == "" {
			return nil, fmt.Errorf("GitLab OAuth credentials not configured")
		}
		return &domain.OAuthConfig{
			ClientID:     clientID,
			ClientSecret: secret,
			TokenURL:     "https://gitlab.com/oauth/token",
			RedirectURI:  redirectURI,
		}, nil

	case domain.ProviderGitLabSelfHosted:
		if gitlabURL == "" || clientID == "" || clientSecret == "" {
			return nil, fmt.Errorf("GitLab URL, Client ID, and Client Secret are required for self-hosted GitLab")
		}
		// URLの末尾のスラッシュを削除
		if gitlabURL[len(gitlabURL)-1] == '/' {
			gitlabURL = gitlabURL[:len(gitlabURL)-1]
		}
		return &domain.OAuthConfig{
			ClientID:     clientID,
			ClientSecret: clientSecret,
			TokenURL:     fmt.Sprintf("%s/oauth/token", gitlabURL),
			RedirectURI:  redirectURI,
		}, nil

	default:
		return nil, fmt.Errorf("unsupported provider: %s", provider)
	}
}

// exchangeCodeForToken exchanges authorization code for access token
func (s *OAuthService) exchangeCodeForToken(config *domain.OAuthConfig, code string) (*domain.OAuthTokenResponse, error) {
	data := url.Values{}
	data.Set("client_id", config.ClientID)
	data.Set("client_secret", config.ClientSecret)
	data.Set("code", code)
	data.Set("redirect_uri", config.RedirectURI)
	data.Set("grant_type", "authorization_code")

	req, err := http.NewRequest("POST", config.TokenURL, bytes.NewBufferString(data.Encode()))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Accept", "application/json")

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("token exchange failed with status %d: %s", resp.StatusCode, string(body))
	}

	var tokenResp domain.OAuthTokenResponse
	if err := json.Unmarshal(body, &tokenResp); err != nil {
		return nil, fmt.Errorf("failed to parse token response: %w", err)
	}

	return &tokenResp, nil
}
