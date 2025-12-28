package domain

import (
	"errors"
	"time"
)

// OAuthConnection represents a user's OAuth connection to a Git provider
type OAuthConnection struct {
	id               string
	userID           string
	provider         Provider
	providerUserID   string
	providerUsername string
	accessToken      string
	refreshToken     string
	tokenExpiresAt   *time.Time
	scopes           []string
	createdAt        time.Time
	updatedAt        time.Time
}

// NewOAuthConnection creates a new OAuth connection
func NewOAuthConnection(
	id string,
	userID string,
	provider Provider,
	providerUserID string,
	providerUsername string,
	accessToken string,
	refreshToken string,
	tokenExpiresAt *time.Time,
	scopes []string,
) (*OAuthConnection, error) {
	if id == "" {
		return nil, errors.New("id is required")
	}
	if userID == "" {
		return nil, errors.New("user_id is required")
	}
	if !provider.IsValid() {
		return nil, errors.New("invalid provider")
	}
	if accessToken == "" {
		return nil, errors.New("access_token is required")
	}

	now := time.Now()
	return &OAuthConnection{
		id:               id,
		userID:           userID,
		provider:         provider,
		providerUserID:   providerUserID,
		providerUsername: providerUsername,
		accessToken:      accessToken,
		refreshToken:     refreshToken,
		tokenExpiresAt:   tokenExpiresAt,
		scopes:           scopes,
		createdAt:        now,
		updatedAt:        now,
	}, nil
}

// Reconstruct creates an OAuth connection from stored data
func ReconstructOAuthConnection(
	id string,
	userID string,
	provider Provider,
	providerUserID string,
	providerUsername string,
	accessToken string,
	refreshToken string,
	tokenExpiresAt *time.Time,
	scopes []string,
	createdAt time.Time,
	updatedAt time.Time,
) *OAuthConnection {
	return &OAuthConnection{
		id:               id,
		userID:           userID,
		provider:         provider,
		providerUserID:   providerUserID,
		providerUsername: providerUsername,
		accessToken:      accessToken,
		refreshToken:     refreshToken,
		tokenExpiresAt:   tokenExpiresAt,
		scopes:           scopes,
		createdAt:        createdAt,
		updatedAt:        updatedAt,
	}
}

// Getters
func (c *OAuthConnection) ID() string                 { return c.id }
func (c *OAuthConnection) UserID() string             { return c.userID }
func (c *OAuthConnection) Provider() Provider         { return c.provider }
func (c *OAuthConnection) ProviderUserID() string     { return c.providerUserID }
func (c *OAuthConnection) ProviderUsername() string   { return c.providerUsername }
func (c *OAuthConnection) AccessToken() string        { return c.accessToken }
func (c *OAuthConnection) RefreshToken() string       { return c.refreshToken }
func (c *OAuthConnection) TokenExpiresAt() *time.Time { return c.tokenExpiresAt }
func (c *OAuthConnection) Scopes() []string           { return c.scopes }
func (c *OAuthConnection) CreatedAt() time.Time       { return c.createdAt }
func (c *OAuthConnection) UpdatedAt() time.Time       { return c.updatedAt }

// IsTokenExpired checks if the access token is expired
func (c *OAuthConnection) IsTokenExpired() bool {
	if c.tokenExpiresAt == nil {
		return false
	}
	// Add 5 minute buffer for token refresh
	return time.Now().Add(5 * time.Minute).After(*c.tokenExpiresAt)
}

// UpdateTokens updates the access and refresh tokens
func (c *OAuthConnection) UpdateTokens(accessToken string, refreshToken string, expiresAt *time.Time) error {
	if accessToken == "" {
		return errors.New("access_token is required")
	}
	c.accessToken = accessToken
	if refreshToken != "" {
		c.refreshToken = refreshToken
	}
	c.tokenExpiresAt = expiresAt
	c.updatedAt = time.Now()
	return nil
}
