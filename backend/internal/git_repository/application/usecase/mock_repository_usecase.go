package repository

import (
	"context"
	"opscore/backend/internal/git_repository/domain/entity"

	"github.com/stretchr/testify/mock"
)

// MockRepositoryUseCase is a mock implementation of the RepositoryUseCase interface for testing
type MockRepositoryUseCase struct {
	mock.Mock
}

// Register is a mock implementation of the RepositoryUseCase.Register method
func (m *MockRepositoryUseCase) Register(ctx context.Context, repoURL string, accessToken string) (entity.Repository, error) {
	args := m.Called(ctx, repoURL, accessToken)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(entity.Repository), args.Error(1)
}

// GetRepository is a mock implementation of the RepositoryUseCase.GetRepository method
func (m *MockRepositoryUseCase) GetRepository(ctx context.Context, repoID string) (entity.Repository, error) {
	args := m.Called(ctx, repoID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(entity.Repository), args.Error(1)
}

// ListRepositories is a mock implementation of the RepositoryUseCase.ListRepositories method
func (m *MockRepositoryUseCase) ListRepositories(ctx context.Context) ([]entity.Repository, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]entity.Repository), args.Error(1)
}

// ListFiles is a mock implementation of the RepositoryUseCase.ListFiles method
func (m *MockRepositoryUseCase) ListFiles(ctx context.Context, repoID string) ([]entity.FileNode, error) {
	args := m.Called(ctx, repoID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]entity.FileNode), args.Error(1)
}

// SelectFiles is a mock implementation of the RepositoryUseCase.SelectFiles method
func (m *MockRepositoryUseCase) SelectFiles(ctx context.Context, repoID string, filePaths []string) error {
	args := m.Called(ctx, repoID, filePaths)
	return args.Error(0)
}

// GetSelectedMarkdown is a mock implementation of the RepositoryUseCase.GetSelectedMarkdown method
func (m *MockRepositoryUseCase) GetSelectedMarkdown(ctx context.Context, repoID string) (string, error) {
	args := m.Called(ctx, repoID)
	return args.String(0), args.Error(1)
}

// UpdateAccessToken is a mock implementation of the RepositoryUseCase.UpdateAccessToken method
func (m *MockRepositoryUseCase) UpdateAccessToken(ctx context.Context, repoID string, accessToken string) error {
	args := m.Called(ctx, repoID, accessToken)
	return args.Error(0)
}
