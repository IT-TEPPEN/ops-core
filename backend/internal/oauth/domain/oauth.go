package domain

// Provider represents a Git provider type
type Provider string

const (
	ProviderGitHub           Provider = "github"
	ProviderGitLab           Provider = "gitlab"
	ProviderGitLabSelfHosted Provider = "gitlab-self-hosted"
)

// IsValid checks if the provider is valid
func (p Provider) IsValid() bool {
	switch p {
	case ProviderGitHub, ProviderGitLab, ProviderGitLabSelfHosted:
		return true
	default:
		return false
	}
}

// OAuthCallbackRequest represents the request body for OAuth callback
type OAuthCallbackRequest struct {
	Provider     Provider `json:"provider" binding:"required"`
	Code         string   `json:"code" binding:"required"`
	State        string   `json:"state" binding:"required"`
	GitLabURL    string   `json:"gitlabUrl,omitempty"`
	ClientID     string   `json:"clientId,omitempty"`
	ClientSecret string   `json:"clientSecret,omitempty"`
}

// OAuthTokenResponse represents the response from OAuth token exchange
type OAuthTokenResponse struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	Scope        string `json:"scope,omitempty"`
	RefreshToken string `json:"refresh_token,omitempty"`
	ExpiresIn    int    `json:"expires_in,omitempty"`
}

// OAuthConfig holds OAuth configuration for a provider
type OAuthConfig struct {
	ClientID     string
	ClientSecret string
	TokenURL     string
	RedirectURI  string
}
