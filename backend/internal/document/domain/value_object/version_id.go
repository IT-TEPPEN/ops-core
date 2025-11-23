package value_object

import (
	"errors"

	"github.com/google/uuid"
)

// VersionID represents a unique identifier for a document version.
type VersionID string

// NewVersionID creates a new VersionID from a string.
func NewVersionID(id string) (VersionID, error) {
	if id == "" {
		return "", errors.New("version ID cannot be empty")
	}
	// Validate that it's a valid UUID
	if _, err := uuid.Parse(id); err != nil {
		return "", errors.New("version ID must be a valid UUID")
	}
	return VersionID(id), nil
}

// GenerateVersionID generates a new VersionID using UUID v4.
func GenerateVersionID() VersionID {
	return VersionID(uuid.New().String())
}

// String returns the string representation of VersionID.
func (v VersionID) String() string {
	return string(v)
}

// IsEmpty returns true if the VersionID is empty.
func (v VersionID) IsEmpty() bool {
	return string(v) == ""
}

// Equals checks if two VersionIDs are equal.
func (v VersionID) Equals(other VersionID) bool {
	return v == other
}
