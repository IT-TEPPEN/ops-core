package domain

import (
	"context"
)

// OAuthConnectionRepository defines the interface for OAuth connection persistence
type OAuthConnectionRepository interface {
	// Save creates or updates an OAuth connection
	Save(ctx context.Context, connection *OAuthConnection) error

	// FindByID finds an OAuth connection by ID
	FindByID(ctx context.Context, id string) (*OAuthConnection, error)

	// FindByUserAndProvider finds an OAuth connection by user ID and provider
	FindByUserAndProvider(ctx context.Context, userID string, provider Provider) (*OAuthConnection, error)

	// FindByUserProviderAndHost finds an OAuth connection by user ID, provider, and host
	FindByUserProviderAndHost(ctx context.Context, userID string, provider Provider, providerHost string) (*OAuthConnection, error)

	// FindByUser finds all OAuth connections for a user
	FindByUser(ctx context.Context, userID string) ([]*OAuthConnection, error)

	// Delete deletes an OAuth connection
	Delete(ctx context.Context, id string) error

	// DeleteByUserAndProvider deletes an OAuth connection by user ID and provider
	DeleteByUserAndProvider(ctx context.Context, userID string, provider Provider) error
}
