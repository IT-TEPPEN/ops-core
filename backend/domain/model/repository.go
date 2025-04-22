package model

import "time"

// Repository represents a registered code repository.
// Fields are based on ADR-0005.
type Repository struct {
	id        string    // Unique identifier (e.g., UUID)
	name      string    // Repository name (e.g., derived from URL)
	url       string    // Git repository URL
	createdAt time.Time // Timestamp of registration
	updatedAt time.Time // Timestamp of last update
}

// NewRepository creates a new Repository instance.
func NewRepository(id, name, url string) *Repository {
	now := time.Now()
	return &Repository{
		id:        id,
		name:      name,
		url:       url,
		createdAt: now,
		updatedAt: now,
	}
}

// ReconstructRepository reconstructs a Repository from persistence data.
func ReconstructRepository(id, name, url string, createdAt, updatedAt time.Time) *Repository {
	return &Repository{
		id:        id,
		name:      name,
		url:       url,
		createdAt: createdAt,
		updatedAt: updatedAt,
	}
}

// ID returns the repository's unique identifier.
func (r *Repository) ID() string {
	return r.id
}

// Name returns the repository name.
func (r *Repository) Name() string {
	return r.name
}

// URL returns the Git repository URL.
func (r *Repository) URL() string {
	return r.url
}

// CreatedAt returns the timestamp when the repository was registered.
func (r *Repository) CreatedAt() time.Time {
	return r.createdAt
}

// UpdatedAt returns the timestamp of the last update.
func (r *Repository) UpdatedAt() time.Time {
	return r.updatedAt
}

// SetUpdatedAt updates the updatedAt timestamp to the current time.
func (r *Repository) SetUpdatedAt() {
	r.updatedAt = time.Now()
}
