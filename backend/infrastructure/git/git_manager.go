package git

import (
	"context"
	"opscore/backend/domain/model"
)

// GitManager defines the interface for interacting with Git repositories.
type GitManager interface {
	// EnsureCloned clones the repository if it's not already present locally
	// or updates it if it exists. Returns the local path to the repository.
	EnsureCloned(ctx context.Context, repo model.Repository) (string, error)
	// ListRepositoryFiles lists all files in the repository at the HEAD commit.
	// Returns a list of file paths relative to the repository root.
	ListRepositoryFiles(ctx context.Context, localPath string, repo model.Repository) ([]string, error)
	// ValidateFilesExist checks if the given file paths exist in the repository.
	ValidateFilesExist(ctx context.Context, localPath string, filePaths []string, repo model.Repository) error
	// ReadManagedFileContent reads the content of a specific file from the local repository.
	ReadManagedFileContent(ctx context.Context, localPath string, filePath string, repo model.Repository) ([]byte, error)
}
