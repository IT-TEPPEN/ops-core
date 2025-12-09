package value_object

import (
	"errors"

	"github.com/google/uuid"
)

// ViewStatisticsID represents a unique identifier for view statistics.
type ViewStatisticsID string

// NewViewStatisticsID creates a new ViewStatisticsID from a string.
func NewViewStatisticsID(id string) (ViewStatisticsID, error) {
	if id == "" {
		return "", errors.New("view statistics ID cannot be empty")
	}
	// Validate that it's a valid UUID
	if _, err := uuid.Parse(id); err != nil {
		return "", errors.New("view statistics ID must be a valid UUID")
	}
	return ViewStatisticsID(id), nil
}

// GenerateViewStatisticsID generates a new ViewStatisticsID using UUID v4.
func GenerateViewStatisticsID() ViewStatisticsID {
	return ViewStatisticsID(uuid.New().String())
}

// String returns the string representation of ViewStatisticsID.
func (v ViewStatisticsID) String() string {
	return string(v)
}

// IsEmpty returns true if the ViewStatisticsID is empty.
func (v ViewStatisticsID) IsEmpty() bool {
	return string(v) == ""
}

// Equals checks if two ViewStatisticsIDs are equal.
func (v ViewStatisticsID) Equals(other ViewStatisticsID) bool {
	return v == other
}
