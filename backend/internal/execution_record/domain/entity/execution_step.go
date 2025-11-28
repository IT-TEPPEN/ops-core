package entity

import (
	"errors"
	"time"

	"opscore/backend/internal/execution_record/domain/value_object"
)

// executionStep represents an execution step within an execution record.
type executionStep struct {
	id                value_object.ExecutionStepID
	executionRecordID value_object.ExecutionRecordID
	stepNumber        int
	description       string
	notes             string
	executedAt        time.Time
}

// ExecutionStep is the interface for an execution step.
type ExecutionStep interface {
	ID() value_object.ExecutionStepID
	ExecutionRecordID() value_object.ExecutionRecordID
	StepNumber() int
	Description() string
	Notes() string
	ExecutedAt() time.Time

	// Behaviors
	UpdateNotes(notes string)
}

// NewExecutionStep creates a new ExecutionStep instance.
func NewExecutionStep(
	id value_object.ExecutionStepID,
	recordID value_object.ExecutionRecordID,
	stepNumber int,
	description string,
) (ExecutionStep, error) {
	if id.IsEmpty() {
		return nil, errors.New("execution step ID cannot be empty")
	}
	if recordID.IsEmpty() {
		return nil, errors.New("execution record ID cannot be empty")
	}
	if stepNumber < 1 {
		return nil, errors.New("step number must be positive")
	}
	if description == "" {
		return nil, errors.New("description cannot be empty")
	}

	return &executionStep{
		id:                id,
		executionRecordID: recordID,
		stepNumber:        stepNumber,
		description:       description,
		notes:             "",
		executedAt:        time.Now(),
	}, nil
}

// ReconstructExecutionStep reconstructs an ExecutionStep from persistence data.
func ReconstructExecutionStep(
	id value_object.ExecutionStepID,
	recordID value_object.ExecutionRecordID,
	stepNumber int,
	description string,
	notes string,
	executedAt time.Time,
) ExecutionStep {
	return &executionStep{
		id:                id,
		executionRecordID: recordID,
		stepNumber:        stepNumber,
		description:       description,
		notes:             notes,
		executedAt:        executedAt,
	}
}

// ID returns the execution step ID.
func (e *executionStep) ID() value_object.ExecutionStepID {
	return e.id
}

// ExecutionRecordID returns the parent execution record ID.
func (e *executionStep) ExecutionRecordID() value_object.ExecutionRecordID {
	return e.executionRecordID
}

// StepNumber returns the step number.
func (e *executionStep) StepNumber() int {
	return e.stepNumber
}

// Description returns the step description.
func (e *executionStep) Description() string {
	return e.description
}

// Notes returns the step notes.
func (e *executionStep) Notes() string {
	return e.notes
}

// ExecutedAt returns the execution timestamp.
func (e *executionStep) ExecutedAt() time.Time {
	return e.executedAt
}

// UpdateNotes updates the step notes.
func (e *executionStep) UpdateNotes(notes string) {
	e.notes = notes
}
