package value_object

import (
	"errors"

	"github.com/google/uuid"
)

// DocumentID represents a unique identifier for a document.
type DocumentID string

// NewDocumentID creates a new DocumentID from a string.
func NewDocumentID(id string) (DocumentID, error) {
	if id == "" {
		return "", errors.New("document ID cannot be empty")
	}
	// Validate that it's a valid UUID
	if _, err := uuid.Parse(id); err != nil {
		return "", errors.New("document ID must be a valid UUID")
	}
	return DocumentID(id), nil
}

// GenerateDocumentID generates a new DocumentID using UUID v4.
func GenerateDocumentID() DocumentID {
	return DocumentID(uuid.New().String())
}

// String returns the string representation of DocumentID.
func (d DocumentID) String() string {
	return string(d)
}

// IsEmpty returns true if the DocumentID is empty.
func (d DocumentID) IsEmpty() bool {
	return string(d) == ""
}

// Equals checks if two DocumentIDs are equal.
func (d DocumentID) Equals(other DocumentID) bool {
	return d == other
}
