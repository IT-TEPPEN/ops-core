package value_object

import "errors"

// ProviderRepositoryID represents the repository/project ID issued by the provider (e.g., GitHub repo ID, GitLab project ID).
type ProviderRepositoryID string

// NewProviderRepositoryID validates and creates a ProviderRepositoryID.
func NewProviderRepositoryID(id string) (ProviderRepositoryID, error) {
	if id == "" {
		return "", errors.New("provider repository ID cannot be empty")
	}
	return ProviderRepositoryID(id), nil
}

// String returns the string representation.
func (p ProviderRepositoryID) String() string {
	return string(p)
}

// IsEmpty returns true if empty.
func (p ProviderRepositoryID) IsEmpty() bool {
	return string(p) == ""
}

// Equals compares two ProviderRepositoryIDs.
func (p ProviderRepositoryID) Equals(other ProviderRepositoryID) bool {
	return p == other
}
