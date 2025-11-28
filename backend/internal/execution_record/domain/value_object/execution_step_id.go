package value_object

import (
	"errors"

	"github.com/google/uuid"
)

// ExecutionStepID represents a unique identifier for an execution step.
type ExecutionStepID string

// NewExecutionStepID creates a new ExecutionStepID from a string.
func NewExecutionStepID(id string) (ExecutionStepID, error) {
	if id == "" {
		return "", errors.New("execution step ID cannot be empty")
	}
	// Validate that it's a valid UUID
	if _, err := uuid.Parse(id); err != nil {
		return "", errors.New("execution step ID must be a valid UUID")
	}
	return ExecutionStepID(id), nil
}

// GenerateExecutionStepID generates a new ExecutionStepID using UUID v4.
func GenerateExecutionStepID() ExecutionStepID {
	return ExecutionStepID(uuid.New().String())
}

// String returns the string representation of ExecutionStepID.
func (e ExecutionStepID) String() string {
	return string(e)
}

// IsEmpty returns true if the ExecutionStepID is empty.
func (e ExecutionStepID) IsEmpty() bool {
	return string(e) == ""
}

// Equals checks if two ExecutionStepIDs are equal.
func (e ExecutionStepID) Equals(other ExecutionStepID) bool {
	return e == other
}
