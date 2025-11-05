package git

import (
	"context"
	"opscore/backend/internal/git_repository/domain/entity"
)

// GitManager defines the interface for interacting with Git repositories.
type GitManager interface {
	// EnsureCloned clones the repository if it's not already present locally
	// or updates it if it exists. Returns the local path to the repository.
	EnsureCloned(ctx context.Context, repo entity.Repository) (string, error)
	// ListRepositoryFiles lists all files in the repository at the HEAD commit.
	// Returns a list of file paths relative to the repository root.
	ListRepositoryFiles(ctx context.Context, localPath string, repo entity.Repository) ([]string, error)
	// ValidateFilesExist checks if the given file paths exist in the repository.
	ValidateFilesExist(ctx context.Context, localPath string, filePaths []string, repo entity.Repository) error
	// ReadManagedFileContent reads the content of a specific file from the local repository.
	ReadManagedFileContent(ctx context.Context, localPath string, filePath string, repo entity.Repository) ([]byte, error)
}
