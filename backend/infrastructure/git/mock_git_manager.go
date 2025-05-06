package git

import (
	"context"
	"opscore/backend/domain/model"

	"github.com/stretchr/testify/mock"
)

// MockGitManager is a mock implementation of GitManager interface for testing
type MockGitManager struct {
	mock.Mock
	mockFiles map[string]map[string]string // repoID -> filePath -> content
}

// NewMockGitManager creates a new instance of MockGitManager for testing
func NewMockGitManager() *MockGitManager {
	return &MockGitManager{
		mockFiles: make(map[string]map[string]string),
	}
}

// AddMockFile adds a mock file to the MockGitManager for testing purposes
func (m *MockGitManager) AddMockFile(repoID, filePath, content string) {
	if _, exists := m.mockFiles[repoID]; !exists {
		m.mockFiles[repoID] = make(map[string]string)
	}
	m.mockFiles[repoID][filePath] = content
}

// EnsureCloned is a mock implementation of the GitManager.EnsureCloned method
func (m *MockGitManager) EnsureCloned(ctx context.Context, repo model.Repository) (string, error) {
	args := m.Called(ctx, repo)
	return args.String(0), args.Error(1)
}

// ListRepositoryFiles is a mock implementation of the GitManager.ListRepositoryFiles method
func (m *MockGitManager) ListRepositoryFiles(ctx context.Context, localPath string, repo model.Repository) ([]string, error) {
	args := m.Called(ctx, localPath, repo)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]string), args.Error(1)
}

// ValidateFilesExist is a mock implementation of the GitManager.ValidateFilesExist method
func (m *MockGitManager) ValidateFilesExist(ctx context.Context, localPath string, filePaths []string, repo model.Repository) error {
	args := m.Called(ctx, localPath, filePaths, repo)
	return args.Error(0)
}

// ReadManagedFileContent is a mock implementation of the GitManager.ReadManagedFileContent method
func (m *MockGitManager) ReadManagedFileContent(ctx context.Context, localPath string, filePath string, repo model.Repository) ([]byte, error) {
	args := m.Called(ctx, localPath, filePath, repo)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]byte), args.Error(1)
}
