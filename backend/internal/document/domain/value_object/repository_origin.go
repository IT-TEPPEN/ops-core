package value_object

import (
	"errors"
	"strings"
)

// RepositoryOrigin represents the origin of a document in a provider repository.
type RepositoryOrigin struct {
	providerRepositoryID ProviderRepositoryID
	owner                string
	repository           string
}

// NewRepositoryOrigin constructs a RepositoryOrigin.
func NewRepositoryOrigin(providerRepoID string, owner string, repository string) (RepositoryOrigin, error) {
	pid, err := NewProviderRepositoryID(providerRepoID)
	if err != nil {
		return RepositoryOrigin{}, err
	}
	if strings.TrimSpace(owner) == "" {
		return RepositoryOrigin{}, errors.New("repository owner cannot be empty")
	}
	if strings.TrimSpace(repository) == "" {
		return RepositoryOrigin{}, errors.New("repository name cannot be empty")
	}
	return RepositoryOrigin{
		providerRepositoryID: pid,
		owner:                strings.TrimSpace(owner),
		repository:           strings.TrimSpace(repository),
	}, nil
}

// ProviderRepositoryID returns the provider repository ID.
func (o RepositoryOrigin) ProviderRepositoryID() ProviderRepositoryID {
	return o.providerRepositoryID
}

// Owner returns the repository owner.
func (o RepositoryOrigin) Owner() string {
	return o.owner
}

// Repository returns the repository name.
func (o RepositoryOrigin) Repository() string {
	return o.repository
}

// Equals checks equality.
func (o RepositoryOrigin) Equals(other RepositoryOrigin) bool {
	return o.providerRepositoryID == other.providerRepositoryID && o.owner == other.owner && o.repository == other.repository
}
