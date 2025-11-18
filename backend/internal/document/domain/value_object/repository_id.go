package value_object

import (
	"errors"

	"github.com/google/uuid"
)

// RepositoryID represents a unique identifier for a repository.
type RepositoryID string

// NewRepositoryID creates a new RepositoryID from a string.
func NewRepositoryID(id string) (RepositoryID, error) {
	if id == "" {
		return "", errors.New("repository ID cannot be empty")
	}
	// Validate that it's a valid UUID
	if _, err := uuid.Parse(id); err != nil {
		return "", errors.New("repository ID must be a valid UUID")
	}
	return RepositoryID(id), nil
}

// String returns the string representation of RepositoryID.
func (r RepositoryID) String() string {
	return string(r)
}

// IsEmpty returns true if the RepositoryID is empty.
func (r RepositoryID) IsEmpty() bool {
	return string(r) == ""
}

// Equals checks if two RepositoryIDs are equal.
func (r RepositoryID) Equals(other RepositoryID) bool {
	return r == other
}
