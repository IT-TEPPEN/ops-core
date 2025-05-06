package repository

import (
	"context"

	"opscore/backend/domain/model"

	"github.com/stretchr/testify/mock"
)

// MockRepository is a mock implementation of Repository interface for testing
type MockRepository struct {
	mock.Mock
}

// Save is a mock implementation of the repository.Save method
func (m *MockRepository) Save(ctx context.Context, repo model.Repository) error {
	args := m.Called(ctx, repo)
	return args.Error(0)
}

// FindByURL is a mock implementation of the repository.FindByURL method
func (m *MockRepository) FindByURL(ctx context.Context, url string) (model.Repository, error) {
	args := m.Called(ctx, url)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(model.Repository), args.Error(1)
}

// FindByID is a mock implementation of the repository.FindByID method
func (m *MockRepository) FindByID(ctx context.Context, id string) (model.Repository, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(model.Repository), args.Error(1)
}

// FindAll is a mock implementation of the repository.FindAll method
func (m *MockRepository) FindAll(ctx context.Context) ([]model.Repository, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]model.Repository), args.Error(1)
}

// SaveManagedFiles is a mock implementation of the repository.SaveManagedFiles method
func (m *MockRepository) SaveManagedFiles(ctx context.Context, repoID string, filePaths []string) error {
	args := m.Called(ctx, repoID, filePaths)
	return args.Error(0)
}

// GetManagedFiles is a mock implementation of the repository.GetManagedFiles method
func (m *MockRepository) GetManagedFiles(ctx context.Context, repoID string) ([]string, error) {
	args := m.Called(ctx, repoID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]string), args.Error(1)
}

// UpdateAccessToken is a mock implementation of the repository.UpdateAccessToken method
func (m *MockRepository) UpdateAccessToken(ctx context.Context, repoID string, accessToken string) error {
	args := m.Called(ctx, repoID, accessToken)
	return args.Error(0)
}
