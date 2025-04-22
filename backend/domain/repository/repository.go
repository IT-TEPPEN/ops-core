package repository

import (
	"context"
	"opscore/backend/domain/model" // Import the domain model
)

// Repository defines the interface for data persistence operations related to repositories.
type Repository interface {
	// Save persists a new repository or updates an existing one.
	Save(ctx context.Context, repo *model.Repository) error
	// FindByURL retrieves a repository by its URL. Returns nil if not found.
	FindByURL(ctx context.Context, url string) (*model.Repository, error)
	// FindByID retrieves a repository by its ID. Returns nil if not found.
	FindByID(ctx context.Context, id string) (*model.Repository, error)
	// FindAll retrieves all repositories.
	FindAll(ctx context.Context) ([]*model.Repository, error)
	// SaveManagedFiles saves the list of file paths selected for management for a given repository.
	// This should replace any existing selection for the repository.
	SaveManagedFiles(ctx context.Context, repoID string, filePaths []string) error
	// GetManagedFiles retrieves the list of file paths selected for management for a given repository.
	GetManagedFiles(ctx context.Context, repoID string) ([]string, error)
	// TODO: Add other necessary methods (e.g., List, Delete)
}
