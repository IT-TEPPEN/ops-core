package service

import (
	"fmt"

	"opscore/backend/internal/auth/domain"
)

// ProviderFactory creates OIDC provider instances
type ProviderFactory struct {
	configs map[string]domain.ProviderConfig
}

// NewProviderFactory creates a new provider factory
func NewProviderFactory(configs map[string]domain.ProviderConfig) *ProviderFactory {
	return &ProviderFactory{
		configs: configs,
	}
}

// CreateProvider creates an OIDC provider instance
func (f *ProviderFactory) CreateProvider(providerName string) (domain.OIDCProvider, error) {
	config, exists := f.configs[providerName]
	if !exists {
		return nil, fmt.Errorf("provider not configured: %s", providerName)
	}

	switch providerName {
	case "google":
		return NewGoogleOAuthService(config.ClientID, config.ClientSecret, config.RedirectURL), nil
	case "github":
		// TODO: Implement GitHub provider
		return nil, fmt.Errorf("GitHub provider not yet implemented")
	case "gitlab":
		// TODO: Implement GitLab provider
		return nil, fmt.Errorf("GitLab provider not yet implemented")
	case "microsoft":
		// TODO: Implement Microsoft provider
		return nil, fmt.Errorf("Microsoft provider not yet implemented")
	default:
		return nil, fmt.Errorf("unsupported provider: %s", providerName)
	}
}

// GetSupportedProviders returns a list of supported provider names
func (f *ProviderFactory) GetSupportedProviders() []string {
	providers := make([]string, 0, len(f.configs))
	for provider := range f.configs {
		providers = append(providers, provider)
	}
	return providers
}
