package domain

import (
	"errors"
	"time"
)

// ProviderMetadata represents provider-specific metadata stored in JSONB
type ProviderMetadata struct {
	ProviderUserID   string   `json:"provider_user_id"`
	ProviderUsername string   `json:"provider_username"`
	Scopes           []string `json:"scopes"`

	// GitLab用（オプショナル）
	RefreshTokenEncrypted *string    `json:"refresh_token_encrypted,omitempty"`
	TokenExpiresAt        *time.Time `json:"token_expires_at,omitempty"`

	// GitLab self-hosted用（オプショナル）
	ClientIDEncrypted     *string `json:"client_id_encrypted,omitempty"`
	ClientSecretEncrypted *string `json:"client_secret_encrypted,omitempty"`
}

// OAuthConnection represents a user's OAuth connection to a Git provider
type OAuthConnection struct {
	id           string
	userID       string
	provider     Provider
	providerHost string
	metadata     ProviderMetadata
	accessToken  string
	createdAt    time.Time
	updatedAt    time.Time
}

// NewOAuthConnection creates a new OAuth connection
func NewOAuthConnection(
	id string,
	userID string,
	provider Provider,
	providerHost string,
	metadata ProviderMetadata,
	accessToken string,
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
	if providerHost == "" {
		return nil, errors.New("provider_host is required")
	}
	if accessToken == "" {
		return nil, errors.New("access_token is required")
	}

	now := time.Now()
	return &OAuthConnection{
		id:           id,
		userID:       userID,
		provider:     provider,
		providerHost: providerHost,
		metadata:     metadata,
		accessToken:  accessToken,
		createdAt:    now,
		updatedAt:    now,
	}, nil
}

// Reconstruct creates an OAuth connection from stored data
func ReconstructOAuthConnection(
	id string,
	userID string,
	provider Provider,
	providerHost string,
	metadata ProviderMetadata,
	accessToken string,
	createdAt time.Time,
	updatedAt time.Time,
) *OAuthConnection {
	return &OAuthConnection{
		id:           id,
		userID:       userID,
		provider:     provider,
		providerHost: providerHost,
		metadata:     metadata,
		accessToken:  accessToken,
		createdAt:    createdAt,
		updatedAt:    updatedAt,
	}
}

// Getters
func (c *OAuthConnection) ID() string                 { return c.id }
func (c *OAuthConnection) UserID() string             { return c.userID }
func (c *OAuthConnection) Provider() Provider         { return c.provider }
func (c *OAuthConnection) ProviderHost() string       { return c.providerHost }
func (c *OAuthConnection) Metadata() ProviderMetadata { return c.metadata }
func (c *OAuthConnection) AccessToken() string        { return c.accessToken }
func (c *OAuthConnection) CreatedAt() time.Time       { return c.createdAt }
func (c *OAuthConnection) UpdatedAt() time.Time       { return c.updatedAt }

// Backward compatibility getters
func (c *OAuthConnection) ProviderUserID() string   { return c.metadata.ProviderUserID }
func (c *OAuthConnection) ProviderUsername() string { return c.metadata.ProviderUsername }
func (c *OAuthConnection) Scopes() []string         { return c.metadata.Scopes }

func (c *OAuthConnection) RefreshToken() string {
	if c.metadata.RefreshTokenEncrypted != nil {
		return *c.metadata.RefreshTokenEncrypted
	}
	return ""
}

func (c *OAuthConnection) TokenExpiresAt() *time.Time {
	return c.metadata.TokenExpiresAt
}

func (c *OAuthConnection) ClientID() string {
	if c.metadata.ClientIDEncrypted != nil {
		return *c.metadata.ClientIDEncrypted
	}
	return ""
}

func (c *OAuthConnection) ClientSecret() string {
	if c.metadata.ClientSecretEncrypted != nil {
		return *c.metadata.ClientSecretEncrypted
	}
	return ""
}

// IsTokenExpired checks if the access token is expired
func (c *OAuthConnection) IsTokenExpired() bool {
	if c.metadata.TokenExpiresAt == nil {
		return false
	}
	// Add 5 minute buffer for token refresh
	return time.Now().Add(5 * time.Minute).After(*c.metadata.TokenExpiresAt)
}

// UpdateTokens updates the access and refresh tokens
func (c *OAuthConnection) UpdateTokens(accessToken string, refreshToken string, expiresAt *time.Time) error {
	if accessToken == "" {
		return errors.New("access_token is required")
	}
	c.accessToken = accessToken
	if refreshToken != "" {
		c.metadata.RefreshTokenEncrypted = &refreshToken
	}
	c.metadata.TokenExpiresAt = expiresAt
	c.updatedAt = time.Now()
	return nil
}
