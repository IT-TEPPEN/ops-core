package repository

import (
	"context"
	"opscore/backend/internal/git_repository/domain/entity"
)

// Repository defines the interface for data persistence operations related to repositories.
type Repository interface {
	// Save persists a new repository or updates an existing one.
	Save(ctx context.Context, repo entity.Repository) error
	// FindByURL retrieves a repository by its URL. Returns nil if not found.
	FindByURL(ctx context.Context, url string) (entity.Repository, error)
	// FindByID retrieves a repository by its ID. Returns nil if not found.
	FindByID(ctx context.Context, id string) (entity.Repository, error)
	// FindAll retrieves all repositories.
	FindAll(ctx context.Context) ([]entity.Repository, error)
	// SaveManagedFiles saves the list of file paths selected for management for a given repository.
	// This should replace any existing selection for the repository.
	SaveManagedFiles(ctx context.Context, repoID string, filePaths []string) error
	// GetManagedFiles retrieves the list of file paths selected for management for a given repository.
	GetManagedFiles(ctx context.Context, repoID string) ([]string, error)
	// UpdateAccessToken updates the access token for a repository.
	UpdateAccessToken(ctx context.Context, repoID string, accessToken string) error
	// TODO: Add other necessary methods (e.g., List, Delete)
}
