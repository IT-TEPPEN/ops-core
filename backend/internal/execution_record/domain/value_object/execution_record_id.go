package value_object

import (
	"errors"

	"github.com/google/uuid"
)

// ExecutionRecordID represents a unique identifier for an execution record.
type ExecutionRecordID string

// NewExecutionRecordID creates a new ExecutionRecordID from a string.
func NewExecutionRecordID(id string) (ExecutionRecordID, error) {
	if id == "" {
		return "", errors.New("execution record ID cannot be empty")
	}
	// Validate that it's a valid UUID
	if _, err := uuid.Parse(id); err != nil {
		return "", errors.New("execution record ID must be a valid UUID")
	}
	return ExecutionRecordID(id), nil
}

// GenerateExecutionRecordID generates a new ExecutionRecordID using UUID v4.
func GenerateExecutionRecordID() ExecutionRecordID {
	return ExecutionRecordID(uuid.New().String())
}

// String returns the string representation of ExecutionRecordID.
func (e ExecutionRecordID) String() string {
	return string(e)
}

// IsEmpty returns true if the ExecutionRecordID is empty.
func (e ExecutionRecordID) IsEmpty() bool {
	return string(e) == ""
}

// Equals checks if two ExecutionRecordIDs are equal.
func (e ExecutionRecordID) Equals(other ExecutionRecordID) bool {
	return e == other
}
