package entity

import "time"

// Repository represents a registered code repository.
// Fields are based on ADR-0005.
type repository struct {
	id          string    // Unique identifier (e.g., UUID)
	name        string    // Repository name (e.g., derived from URL)
	url         string    // Git repository URL
	accessToken string    // Access token for private repositories (e.g., PAT, Deploy Key)
	createdAt   time.Time // Timestamp of registration
	updatedAt   time.Time // Timestamp of last update
}

// Repository interface defines the methods for a repository.
type Repository interface {
	ID() string
	Name() string
	URL() string
	AccessToken() string
	CreatedAt() time.Time
	UpdatedAt() time.Time
	SetUpdatedAt()
	SetAccessToken(token string)
}

// NewRepository creates a new Repository instance.
func NewRepository(id, name, url, accessToken string) Repository {
	now := time.Now()
	return &repository{
		id:          id,
		name:        name,
		url:         url,
		accessToken: accessToken,
		createdAt:   now,
		updatedAt:   now,
	}
}

// ReconstructRepository reconstructs a Repository from persistence data.
func ReconstructRepository(id, name, url, accessToken string, createdAt, updatedAt time.Time) Repository {
	return &repository{
		id:          id,
		name:        name,
		url:         url,
		accessToken: accessToken,
		createdAt:   createdAt,
		updatedAt:   updatedAt,
	}
}

// ID returns the repository's unique identifier.
func (r *repository) ID() string {
	return r.id
}

// Name returns the repository name.
func (r *repository) Name() string {
	return r.name
}

// URL returns the Git repository URL.
func (r *repository) URL() string {
	return r.url
}

// AccessToken returns the repository access token.
func (r *repository) AccessToken() string {
	return r.accessToken
}

// CreatedAt returns the timestamp when the repository was registered.
func (r *repository) CreatedAt() time.Time {
	return r.createdAt
}

// UpdatedAt returns the timestamp of the last update.
func (r *repository) UpdatedAt() time.Time {
	return r.updatedAt
}

// SetUpdatedAt updates the updatedAt timestamp to the current time.
func (r *repository) SetUpdatedAt() {
	r.updatedAt = time.Now()
}

// SetAccessToken updates the access token for the repository.
func (r *repository) SetAccessToken(token string) {
	r.accessToken = token
	r.SetUpdatedAt()
}
