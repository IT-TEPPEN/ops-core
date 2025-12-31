package service

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"opscore/backend/internal/oauth/domain"

	"github.com/google/uuid"
)

// GitHubUserResponse represents GitHub's user API response
type GitHubUserResponse struct {
	ID    int64  `json:"id"`
	Login string `json:"login"`
}

// OAuthService handles OAuth2.0 token exchange
type OAuthService struct {
	httpClient *http.Client
	repo       domain.OAuthConnectionRepository
}

// NewOAuthService creates a new OAuthService
func NewOAuthService() *OAuthService {
	return &OAuthService{
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
		repo: nil,
	}
}

// NewOAuthServiceWithRepository creates a new OAuthService with a repository
func NewOAuthServiceWithRepository(repo domain.OAuthConnectionRepository) *OAuthService {
	return &OAuthService{
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
		repo: repo,
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

// extractProviderHost extracts the host from the provider
func (s *OAuthService) extractProviderHost(provider domain.Provider, gitlabURL string) (string, error) {
	switch provider {
	case domain.ProviderGitHub:
		return "github.com", nil
	case domain.ProviderGitLab:
		return "gitlab.com", nil
	case domain.ProviderGitLabSelfHosted:
		if gitlabURL == "" {
			return "", fmt.Errorf("gitlabURL is required for self-hosted GitLab")
		}
		// Parse URL and extract host
		parsedURL, err := url.Parse(gitlabURL)
		if err != nil {
			return "", fmt.Errorf("failed to parse GitLab URL: %w", err)
		}
		return parsedURL.Host, nil
	default:
		return "", fmt.Errorf("unsupported provider: %s", provider)
	}
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

// ExchangeAndSaveToken exchanges authorization code for tokens and saves them to the database
func (s *OAuthService) ExchangeAndSaveToken(ctx context.Context, userID string, req domain.OAuthCallbackRequest) (*domain.OAuthConnection, error) {
	if s.repo == nil {
		return nil, fmt.Errorf("repository not configured")
	}

	// Exchange code for token
	tokenResp, err := s.ExchangeToken(req)
	if err != nil {
		return nil, err
	}

	// Extract provider host
	providerHost, err := s.extractProviderHost(req.Provider, req.GitLabURL)
	if err != nil {
		return nil, fmt.Errorf("failed to extract provider host: %w", err)
	}

	// Get user info from provider
	providerUserID, providerUsername, err := s.getProviderUserInfo(ctx, req.Provider, tokenResp.AccessToken, req.GitLabURL)
	if err != nil {
		return nil, fmt.Errorf("failed to get provider user info: %w", err)
	}

	// Calculate token expiration
	var expiresAt *time.Time
	if tokenResp.ExpiresIn > 0 {
		t := time.Now().Add(time.Duration(tokenResp.ExpiresIn) * time.Second)
		expiresAt = &t
	}

	// Parse scopes - GitHub returns space-separated scopes, GitLab uses comma-separated
	var scopes []string
	if tokenResp.Scope != "" {
		// Handle both space and comma separators
		scopeStr := strings.ReplaceAll(tokenResp.Scope, " ", ",")
		for _, s := range strings.Split(scopeStr, ",") {
			s = strings.TrimSpace(s)
			if s != "" {
				scopes = append(scopes, s)
			}
		}
	}
	// If no scopes returned, use default scopes based on provider
	if len(scopes) == 0 {
		switch req.Provider {
		case domain.ProviderGitHub:
			scopes = []string{"repo", "read:user", "user:email"}
		case domain.ProviderGitLab, domain.ProviderGitLabSelfHosted:
			scopes = []string{"api", "read_user", "read_repository"}
		}
	}

	// Build provider metadata
	metadata := domain.ProviderMetadata{
		ProviderUserID:   providerUserID,
		ProviderUsername: providerUsername,
		Scopes:           scopes,
	}

	// Add refresh token if present (GitLab)
	if tokenResp.RefreshToken != "" {
		metadata.RefreshTokenEncrypted = &tokenResp.RefreshToken
	}

	// Add token expiration if present (GitLab)
	if expiresAt != nil {
		metadata.TokenExpiresAt = expiresAt
	}

	// Add client credentials for self-hosted GitLab
	if req.Provider == domain.ProviderGitLabSelfHosted {
		if req.ClientID != "" {
			metadata.ClientIDEncrypted = &req.ClientID
		}
		if req.ClientSecret != "" {
			metadata.ClientSecretEncrypted = &req.ClientSecret
		}
	}

	// Create OAuth connection
	conn, err := domain.NewOAuthConnection(
		uuid.New().String(),
		userID,
		req.Provider,
		providerHost,
		metadata,
		tokenResp.AccessToken,
	)
	if err != nil {
		return nil, err
	}

	// Save to database
	if err := s.repo.Save(ctx, conn); err != nil {
		return nil, fmt.Errorf("failed to save oauth connection: %w", err)
	}

	return conn, nil
}

// getProviderUserInfo fetches user info from the OAuth provider
func (s *OAuthService) getProviderUserInfo(ctx context.Context, provider domain.Provider, accessToken string, gitlabURL string) (string, string, error) {
	switch provider {
	case domain.ProviderGitHub:
		return s.getGitHubUserInfo(ctx, accessToken)
	case domain.ProviderGitLab, domain.ProviderGitLabSelfHosted:
		baseURL := "https://gitlab.com"
		if provider == domain.ProviderGitLabSelfHosted && gitlabURL != "" {
			baseURL = strings.TrimSuffix(gitlabURL, "/")
		}
		return s.getGitLabUserInfo(ctx, accessToken, baseURL)
	default:
		return "", "", fmt.Errorf("unsupported provider: %s", provider)
	}
}

// getGitHubUserInfo fetches user info from GitHub
func (s *OAuthService) getGitHubUserInfo(ctx context.Context, accessToken string) (string, string, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", "https://api.github.com/user", nil)
	if err != nil {
		return "", "", err
	}
	req.Header.Set("Authorization", "Bearer "+accessToken)
	req.Header.Set("Accept", "application/vnd.github+json")

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return "", "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", "", fmt.Errorf("github api error: %s - %s", resp.Status, string(body))
	}

	var user GitHubUserResponse
	if err := json.NewDecoder(resp.Body).Decode(&user); err != nil {
		return "", "", err
	}

	return fmt.Sprintf("%d", user.ID), user.Login, nil
}

// GitLabUserResponse represents GitLab's user API response
type GitLabUserResponse struct {
	ID       int64  `json:"id"`
	Username string `json:"username"`
}

// getGitLabUserInfo fetches user info from GitLab
func (s *OAuthService) getGitLabUserInfo(ctx context.Context, accessToken string, baseURL string) (string, string, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", baseURL+"/api/v4/user", nil)
	if err != nil {
		return "", "", err
	}
	req.Header.Set("Authorization", "Bearer "+accessToken)

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return "", "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", "", fmt.Errorf("gitlab api error: %s - %s", resp.Status, string(body))
	}

	var user GitLabUserResponse
	if err := json.NewDecoder(resp.Body).Decode(&user); err != nil {
		return "", "", err
	}

	return fmt.Sprintf("%d", user.ID), user.Username, nil
}

// GetConnection retrieves a user's OAuth connection for a provider
func (s *OAuthService) GetConnection(ctx context.Context, userID string, provider domain.Provider) (*domain.OAuthConnection, error) {
	if s.repo == nil {
		return nil, fmt.Errorf("repository not configured")
	}
	return s.repo.FindByUserAndProvider(ctx, userID, provider)
}

// GetConnectionByHost retrieves a user's OAuth connection for a provider and host
func (s *OAuthService) GetConnectionByHost(ctx context.Context, userID string, provider domain.Provider, providerHost string) (*domain.OAuthConnection, error) {
	if s.repo == nil {
		return nil, fmt.Errorf("repository not configured")
	}
	return s.repo.FindByUserProviderAndHost(ctx, userID, provider, providerHost)
}

// ExtractHostFromRepoURL extracts the host from a repository URL
func ExtractHostFromRepoURL(repoURL string) (string, error) {
	parsedURL, err := url.Parse(repoURL)
	if err != nil {
		return "", fmt.Errorf("failed to parse repository URL: %w", err)
	}
	if parsedURL.Host == "" {
		return "", fmt.Errorf("no host found in repository URL: %s", repoURL)
	}
	return parsedURL.Host, nil
}

// DetermineProviderFromHost determines the provider type from a host
func DetermineProviderFromHost(host string) domain.Provider {
	switch {
	case strings.Contains(host, "github.com"):
		return domain.ProviderGitHub
	case strings.Contains(host, "gitlab.com"):
		return domain.ProviderGitLab
	default:
		// Assume self-hosted GitLab for other hosts
		return domain.ProviderGitLabSelfHosted
	}
}

// GetAccessToken retrieves a valid access token for a user and provider
// If the token is expired, it will attempt to refresh it
func (s *OAuthService) GetAccessToken(ctx context.Context, userID string, provider domain.Provider) (string, error) {
	conn, err := s.GetConnection(ctx, userID, provider)
	if err != nil {
		return "", err
	}
	if conn == nil {
		return "", fmt.Errorf("no oauth connection found for provider %s", provider)
	}

	// Check if token is expired
	if conn.IsTokenExpired() {
		// Try to refresh the token
		refreshedConn, err := s.RefreshToken(ctx, conn)
		if err != nil {
			return "", fmt.Errorf("token expired and refresh failed: %w", err)
		}
		return refreshedConn.AccessToken(), nil
	}

	return conn.AccessToken(), nil
}

// GetAccessTokenForRepoURL retrieves a valid access token for a repository URL
func (s *OAuthService) GetAccessTokenForRepoURL(ctx context.Context, userID string, repoURL string) (string, error) {
	// Extract host from repository URL
	host, err := ExtractHostFromRepoURL(repoURL)
	if err != nil {
		return "", err
	}

	// Determine provider from host
	provider := DetermineProviderFromHost(host)

	// Get OAuth connection by host
	conn, err := s.GetConnectionByHost(ctx, userID, provider, host)
	if err != nil {
		return "", err
	}
	if conn == nil {
		return "", fmt.Errorf("no oauth connection found for %s (%s)", host, provider)
	}

	// Check if token is expired
	if conn.IsTokenExpired() {
		// Try to refresh the token
		refreshedConn, err := s.RefreshToken(ctx, conn)
		if err != nil {
			return "", fmt.Errorf("token expired and refresh failed: %w", err)
		}
		return refreshedConn.AccessToken(), nil
	}

	return conn.AccessToken(), nil
}

// RefreshToken refreshes an expired access token
func (s *OAuthService) RefreshToken(ctx context.Context, conn *domain.OAuthConnection) (*domain.OAuthConnection, error) {
	if conn.RefreshToken() == "" {
		return nil, fmt.Errorf("no refresh token available")
	}

	switch conn.Provider() {
	case domain.ProviderGitHub:
		return s.refreshGitHubToken(ctx, conn)
	case domain.ProviderGitLab, domain.ProviderGitLabSelfHosted:
		return s.refreshGitLabToken(ctx, conn)
	default:
		return nil, fmt.Errorf("unsupported provider: %s", conn.Provider())
	}
}

func (s *OAuthService) refreshGitHubToken(ctx context.Context, conn *domain.OAuthConnection) (*domain.OAuthConnection, error) {
	clientID := os.Getenv("GITHUB_CLIENT_ID")
	clientSecret := os.Getenv("GITHUB_CLIENT_SECRET")
	if clientID == "" || clientSecret == "" {
		return nil, fmt.Errorf("GitHub OAuth credentials not configured")
	}

	data := url.Values{}
	data.Set("client_id", clientID)
	data.Set("client_secret", clientSecret)
	data.Set("grant_type", "refresh_token")
	data.Set("refresh_token", conn.RefreshToken())

	req, err := http.NewRequestWithContext(ctx, "POST", "https://github.com/login/oauth/access_token", bytes.NewBufferString(data.Encode()))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Accept", "application/json")

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var tokenResp domain.OAuthTokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&tokenResp); err != nil {
		return nil, err
	}

	if tokenResp.AccessToken == "" {
		return nil, fmt.Errorf("failed to refresh github token")
	}

	var expiresAt *time.Time
	if tokenResp.ExpiresIn > 0 {
		t := time.Now().Add(time.Duration(tokenResp.ExpiresIn) * time.Second)
		expiresAt = &t
	}

	if err := conn.UpdateTokens(tokenResp.AccessToken, tokenResp.RefreshToken, expiresAt); err != nil {
		return nil, err
	}

	if s.repo != nil {
		if err := s.repo.Save(ctx, conn); err != nil {
			return nil, err
		}
	}

	return conn, nil
}

func (s *OAuthService) refreshGitLabToken(ctx context.Context, conn *domain.OAuthConnection) (*domain.OAuthConnection, error) {
	clientID := os.Getenv("GITLAB_CLIENT_ID")
	clientSecret := os.Getenv("GITLAB_CLIENT_SECRET")
	if clientID == "" || clientSecret == "" {
		return nil, fmt.Errorf("GitLab OAuth credentials not configured")
	}

	tokenURL := "https://gitlab.com/oauth/token"
	// TODO: Handle self-hosted GitLab URLs

	data := url.Values{}
	data.Set("client_id", clientID)
	data.Set("client_secret", clientSecret)
	data.Set("grant_type", "refresh_token")
	data.Set("refresh_token", conn.RefreshToken())

	req, err := http.NewRequestWithContext(ctx, "POST", tokenURL, bytes.NewBufferString(data.Encode()))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Accept", "application/json")

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var tokenResp domain.OAuthTokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&tokenResp); err != nil {
		return nil, err
	}

	if tokenResp.AccessToken == "" {
		return nil, fmt.Errorf("failed to refresh gitlab token")
	}

	var expiresAt *time.Time
	if tokenResp.ExpiresIn > 0 {
		t := time.Now().Add(time.Duration(tokenResp.ExpiresIn) * time.Second)
		expiresAt = &t
	}

	if err := conn.UpdateTokens(tokenResp.AccessToken, tokenResp.RefreshToken, expiresAt); err != nil {
		return nil, err
	}

	if s.repo != nil {
		if err := s.repo.Save(ctx, conn); err != nil {
			return nil, err
		}
	}

	return conn, nil
}

// DisconnectProvider removes a user's OAuth connection
func (s *OAuthService) DisconnectProvider(ctx context.Context, userID string, provider domain.Provider) error {
	if s.repo == nil {
		return fmt.Errorf("repository not configured")
	}
	return s.repo.DeleteByUserAndProvider(ctx, userID, provider)
}

// ListConnections lists all OAuth connections for a user
func (s *OAuthService) ListConnections(ctx context.Context, userID string) ([]*domain.OAuthConnection, error) {
	if s.repo == nil {
		return nil, fmt.Errorf("repository not configured")
	}
	return s.repo.FindByUser(ctx, userID)
}

// GetAccessTokenByConnectionID retrieves a valid access token for a specific connection ID
// If the token is expired, it will attempt to refresh it
func (s *OAuthService) GetAccessTokenByConnectionID(ctx context.Context, userID string, connectionID string) (string, error) {
	if s.repo == nil {
		return "", fmt.Errorf("repository not configured")
	}

	conn, err := s.repo.FindByID(ctx, connectionID)
	if err != nil {
		return "", err
	}
	if conn == nil {
		return "", fmt.Errorf("no oauth connection found for connection ID %s", connectionID)
	}

	// Verify the connection belongs to the user
	if conn.UserID() != userID {
		return "", fmt.Errorf("connection %s does not belong to user %s", connectionID, userID)
	}

	// Check if token is expired
	if conn.IsTokenExpired() {
		// Try to refresh the token
		refreshedConn, err := s.RefreshToken(ctx, conn)
		if err != nil {
			return "", fmt.Errorf("token expired and refresh failed: %w", err)
		}
		return refreshedConn.AccessToken(), nil
	}

	return conn.AccessToken(), nil
}

// GetConnectionByID retrieves a connection by its ID
func (s *OAuthService) GetConnectionByID(ctx context.Context, userID string, connectionID string) (*domain.OAuthConnection, error) {
	if s.repo == nil {
		return nil, fmt.Errorf("repository not configured")
	}

	conn, err := s.repo.FindByID(ctx, connectionID)
	if err != nil {
		return nil, err
	}
	if conn == nil {
		return nil, fmt.Errorf("no oauth connection found for connection ID %s", connectionID)
	}

	// Verify the connection belongs to the user
	if conn.UserID() != userID {
		return nil, fmt.Errorf("connection %s does not belong to user %s", connectionID, userID)
	}

	return conn, nil
}
