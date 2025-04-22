package persistence

import (
	"context"
	"sync"

	"opscore/backend/domain/model"
	"opscore/backend/domain/repository"
)

// InMemoryRepository is an in-memory implementation of the repository interface.
// It's useful for testing and development.
type InMemoryRepository struct {
	repositories     map[string]*model.Repository
	urls             map[string]string
	managedFilesById map[string][]string
	mu               sync.RWMutex
}

// NewInMemoryRepository creates a new in-memory repository.
func NewInMemoryRepository() repository.Repository {
	return &InMemoryRepository{
		repositories:     make(map[string]*model.Repository),
		urls:             make(map[string]string), // map url -> id
		managedFilesById: make(map[string][]string),
	}
}

// Save stores a repository in memory.
func (r *InMemoryRepository) Save(ctx context.Context, repo *model.Repository) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	repoURL := repo.URL()
	repoID := repo.ID()

	if existingID, found := r.urls[repoURL]; found && existingID != repoID {
		// Another repository with same URL but different ID exists
		return nil // Return error in the future
	}

	r.repositories[repoID] = repo
	r.urls[repoURL] = repoID
	return nil
}

// FindByURL finds a repository by URL.
func (r *InMemoryRepository) FindByURL(ctx context.Context, url string) (*model.Repository, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	if id, found := r.urls[url]; found {
		return r.repositories[id], nil
	}
	return nil, nil
}

// FindByID finds a repository by ID.
func (r *InMemoryRepository) FindByID(ctx context.Context, id string) (*model.Repository, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	return r.repositories[id], nil
}

// SaveManagedFiles stores the managed files for a repository.
func (r *InMemoryRepository) SaveManagedFiles(ctx context.Context, repoID string, filePaths []string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	// Check if repository exists
	if _, exists := r.repositories[repoID]; !exists {
		return nil // Return error in the future
	}

	// Replace any existing selection
	r.managedFilesById[repoID] = make([]string, len(filePaths))
	copy(r.managedFilesById[repoID], filePaths)
	return nil
}

// GetManagedFiles retrieves the managed files for a repository.
func (r *InMemoryRepository) GetManagedFiles(ctx context.Context, repoID string) ([]string, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	// Check if repository exists
	if _, exists := r.repositories[repoID]; !exists {
		return nil, nil // Return error in the future
	}

	files := r.managedFilesById[repoID]
	if files == nil {
		return []string{}, nil // Return empty slice, not nil
	}

	// Return a copy to prevent external modification
	result := make([]string, len(files))
	copy(result, files)
	return result, nil
}
