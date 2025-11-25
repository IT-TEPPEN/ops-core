package value_object

import "errors"

// ExecutionStatus represents the status of an execution record.
type ExecutionStatus string

const (
	// ExecutionStatusInProgress represents an execution in progress.
	ExecutionStatusInProgress ExecutionStatus = "in_progress"
	// ExecutionStatusCompleted represents a completed execution.
	ExecutionStatusCompleted ExecutionStatus = "completed"
	// ExecutionStatusFailed represents a failed execution.
	ExecutionStatusFailed ExecutionStatus = "failed"
)

// NewExecutionStatus creates a new ExecutionStatus from a string.
func NewExecutionStatus(status string) (ExecutionStatus, error) {
	execStatus := ExecutionStatus(status)
	if !execStatus.IsValid() {
		return "", errors.New("invalid execution status: must be 'in_progress', 'completed', or 'failed'")
	}
	return execStatus, nil
}

// IsValid checks if the ExecutionStatus is valid.
func (e ExecutionStatus) IsValid() bool {
	return e == ExecutionStatusInProgress || e == ExecutionStatusCompleted || e == ExecutionStatusFailed
}

// String returns the string representation of ExecutionStatus.
func (e ExecutionStatus) String() string {
	return string(e)
}

// IsInProgress returns true if the status is in progress.
func (e ExecutionStatus) IsInProgress() bool {
	return e == ExecutionStatusInProgress
}

// IsCompleted returns true if the status is completed.
func (e ExecutionStatus) IsCompleted() bool {
	return e == ExecutionStatusCompleted
}

// IsFailed returns true if the status is failed.
func (e ExecutionStatus) IsFailed() bool {
	return e == ExecutionStatusFailed
}

// Equals checks if two ExecutionStatuses are equal.
func (e ExecutionStatus) Equals(other ExecutionStatus) bool {
	return e == other
}
