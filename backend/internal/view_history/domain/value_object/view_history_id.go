package value_object

import (
	"errors"

	"github.com/google/uuid"
)

// ViewHistoryID represents a unique identifier for a view history record.
type ViewHistoryID string

// NewViewHistoryID creates a new ViewHistoryID from a string.
func NewViewHistoryID(id string) (ViewHistoryID, error) {
	if id == "" {
		return "", errors.New("view history ID cannot be empty")
	}
	// Validate that it's a valid UUID
	if _, err := uuid.Parse(id); err != nil {
		return "", errors.New("view history ID must be a valid UUID")
	}
	return ViewHistoryID(id), nil
}

// GenerateViewHistoryID generates a new ViewHistoryID using UUID v4.
func GenerateViewHistoryID() ViewHistoryID {
	return ViewHistoryID(uuid.New().String())
}

// String returns the string representation of ViewHistoryID.
func (v ViewHistoryID) String() string {
	return string(v)
}

// IsEmpty returns true if the ViewHistoryID is empty.
func (v ViewHistoryID) IsEmpty() bool {
	return string(v) == ""
}

// Equals checks if two ViewHistoryIDs are equal.
func (v ViewHistoryID) Equals(other ViewHistoryID) bool {
	return v == other
}
