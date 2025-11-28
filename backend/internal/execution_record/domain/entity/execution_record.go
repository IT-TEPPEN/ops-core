package entity

import (
	"errors"
	"time"

	docvo "opscore/backend/internal/document/domain/value_object"
	"opscore/backend/internal/execution_record/domain/value_object"
)

// executionRecord represents an execution record (aggregate root).
type executionRecord struct {
	id                value_object.ExecutionRecordID
	documentID        docvo.DocumentID
	documentVersionID docvo.VersionID
	executorID        string // User ID as string (reference to user aggregate)
	title             string
	variableValues    []value_object.VariableValue
	notes             string
	status            value_object.ExecutionStatus
	accessScope       value_object.AccessScope
	steps             []ExecutionStep
	startedAt         time.Time
	completedAt       *time.Time
	createdAt         time.Time
	updatedAt         time.Time
}

// ExecutionRecord is the interface for an execution record (aggregate root).
type ExecutionRecord interface {
	ID() value_object.ExecutionRecordID
	DocumentID() docvo.DocumentID
	DocumentVersionID() docvo.VersionID
	ExecutorID() string
	Title() string
	VariableValues() []value_object.VariableValue
	Notes() string
	Status() value_object.ExecutionStatus
	AccessScope() value_object.AccessScope
	Steps() []ExecutionStep
	StartedAt() time.Time
	CompletedAt() *time.Time
	CreatedAt() time.Time
	UpdatedAt() time.Time

	// Behaviors
	AddStep(stepNumber int, description string) error
	UpdateStepNotes(stepNumber int, notes string) error
	UpdateNotes(notes string)
	UpdateTitle(title string) error
	Complete() error
	MarkAsFailed() error
	UpdateAccessScope(scope value_object.AccessScope)
}

// NewExecutionRecord creates a new ExecutionRecord instance.
func NewExecutionRecord(
	id value_object.ExecutionRecordID,
	documentID docvo.DocumentID,
	versionID docvo.VersionID,
	executorID string,
	title string,
	variableValues []value_object.VariableValue,
) (ExecutionRecord, error) {
	if id.IsEmpty() {
		return nil, errors.New("execution record ID cannot be empty")
	}
	if documentID.IsEmpty() {
		return nil, errors.New("document ID cannot be empty")
	}
	if versionID.IsEmpty() {
		return nil, errors.New("document version ID cannot be empty")
	}
	if executorID == "" {
		return nil, errors.New("executor ID cannot be empty")
	}
	if title == "" {
		return nil, errors.New("title cannot be empty")
	}

	now := time.Now()
	return &executionRecord{
		id:                id,
		documentID:        documentID,
		documentVersionID: versionID,
		executorID:        executorID,
		title:             title,
		variableValues:    variableValues,
		notes:             "",
		status:            value_object.ExecutionStatusInProgress,
		accessScope:       value_object.AccessScopePrivate,
		steps:             []ExecutionStep{},
		startedAt:         now,
		completedAt:       nil,
		createdAt:         now,
		updatedAt:         now,
	}, nil
}

// ReconstructExecutionRecord reconstructs an ExecutionRecord from persistence data.
func ReconstructExecutionRecord(
	id value_object.ExecutionRecordID,
	documentID docvo.DocumentID,
	documentVersionID docvo.VersionID,
	executorID string,
	title string,
	variableValues []value_object.VariableValue,
	notes string,
	status value_object.ExecutionStatus,
	accessScope value_object.AccessScope,
	steps []ExecutionStep,
	startedAt time.Time,
	completedAt *time.Time,
	createdAt time.Time,
	updatedAt time.Time,
) ExecutionRecord {
	return &executionRecord{
		id:                id,
		documentID:        documentID,
		documentVersionID: documentVersionID,
		executorID:        executorID,
		title:             title,
		variableValues:    variableValues,
		notes:             notes,
		status:            status,
		accessScope:       accessScope,
		steps:             steps,
		startedAt:         startedAt,
		completedAt:       completedAt,
		createdAt:         createdAt,
		updatedAt:         updatedAt,
	}
}

// Getter methods

// ID returns the execution record ID.
func (e *executionRecord) ID() value_object.ExecutionRecordID {
	return e.id
}

// DocumentID returns the document ID.
func (e *executionRecord) DocumentID() docvo.DocumentID {
	return e.documentID
}

// DocumentVersionID returns the document version ID.
func (e *executionRecord) DocumentVersionID() docvo.VersionID {
	return e.documentVersionID
}

// ExecutorID returns the executor user ID.
func (e *executionRecord) ExecutorID() string {
	return e.executorID
}

// Title returns the execution title.
func (e *executionRecord) Title() string {
	return e.title
}

// VariableValues returns the variable values used during execution.
func (e *executionRecord) VariableValues() []value_object.VariableValue {
	return e.variableValues
}

// Notes returns the overall notes.
func (e *executionRecord) Notes() string {
	return e.notes
}

// Status returns the execution status.
func (e *executionRecord) Status() value_object.ExecutionStatus {
	return e.status
}

// AccessScope returns the access scope.
func (e *executionRecord) AccessScope() value_object.AccessScope {
	return e.accessScope
}

// Steps returns the execution steps.
func (e *executionRecord) Steps() []ExecutionStep {
	return e.steps
}

// StartedAt returns the start timestamp.
func (e *executionRecord) StartedAt() time.Time {
	return e.startedAt
}

// CompletedAt returns the completion timestamp.
func (e *executionRecord) CompletedAt() *time.Time {
	return e.completedAt
}

// CreatedAt returns the creation timestamp.
func (e *executionRecord) CreatedAt() time.Time {
	return e.createdAt
}

// UpdatedAt returns the last update timestamp.
func (e *executionRecord) UpdatedAt() time.Time {
	return e.updatedAt
}

// Behavior methods

// AddStep adds a new step to the execution record.
func (e *executionRecord) AddStep(stepNumber int, description string) error {
	if !e.status.IsInProgress() {
		return errors.New("cannot add step to a completed or failed execution")
	}

	// Check for duplicate step number
	for _, step := range e.steps {
		if step.StepNumber() == stepNumber {
			return errors.New("step with this number already exists")
		}
	}

	stepID := value_object.GenerateExecutionStepID()
	step, err := NewExecutionStep(stepID, e.id, stepNumber, description)
	if err != nil {
		return err
	}

	e.steps = append(e.steps, step)
	e.updatedAt = time.Now()
	return nil
}

// UpdateStepNotes updates the notes for a specific step.
func (e *executionRecord) UpdateStepNotes(stepNumber int, notes string) error {
	for _, step := range e.steps {
		if step.StepNumber() == stepNumber {
			step.UpdateNotes(notes)
			e.updatedAt = time.Now()
			return nil
		}
	}
	return errors.New("step not found")
}

// UpdateNotes updates the overall notes.
func (e *executionRecord) UpdateNotes(notes string) {
	e.notes = notes
	e.updatedAt = time.Now()
}

// UpdateTitle updates the execution title.
func (e *executionRecord) UpdateTitle(title string) error {
	if title == "" {
		return errors.New("title cannot be empty")
	}
	e.title = title
	e.updatedAt = time.Now()
	return nil
}

// Complete marks the execution as completed.
func (e *executionRecord) Complete() error {
	if !e.status.IsInProgress() {
		return errors.New("only in-progress executions can be completed")
	}
	e.status = value_object.ExecutionStatusCompleted
	now := time.Now()
	e.completedAt = &now
	e.updatedAt = now
	return nil
}

// MarkAsFailed marks the execution as failed.
func (e *executionRecord) MarkAsFailed() error {
	if !e.status.IsInProgress() {
		return errors.New("only in-progress executions can be marked as failed")
	}
	e.status = value_object.ExecutionStatusFailed
	now := time.Now()
	e.completedAt = &now
	e.updatedAt = now
	return nil
}

// UpdateAccessScope updates the access scope.
func (e *executionRecord) UpdateAccessScope(scope value_object.AccessScope) {
	e.accessScope = scope
	e.updatedAt = time.Now()
}
