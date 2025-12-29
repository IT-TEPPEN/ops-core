package repository

import (
	"context"

	"github.com/stretchr/testify/mock"
)

// MockOAuthTokenProvider is a mock implementation of the OAuthTokenProvider interface for testing
type MockOAuthTokenProvider struct {
	mock.Mock
}

// GetAccessTokenForProvider is a mock implementation
func (m *MockOAuthTokenProvider) GetAccessTokenForProvider(ctx context.Context, userID string, providerName string) (string, error) {
	args := m.Called(ctx, userID, providerName)
	return args.String(0), args.Error(1)
}
