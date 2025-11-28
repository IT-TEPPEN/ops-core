package dto

import (
	"time"
)

// CreateExecutionRecordRequest represents the request to create an execution record.
type CreateExecutionRecordRequest struct {
	DocumentID        string
	DocumentVersionID string
	ExecutorID        string
	Title             string
	VariableValues    []VariableValueDTO
}

// ExecutionRecordResponse represents an execution record response.
type ExecutionRecordResponse struct {
	ID                string
	DocumentID        string
	DocumentVersionID string
	ExecutorID        string
	Title             string
	VariableValues    []VariableValueDTO
	Notes             string
	Status            string
	AccessScope       string
	Steps             []ExecutionStepResponse
	StartedAt         time.Time
	CompletedAt       *time.Time
	CreatedAt         time.Time
	UpdatedAt         time.Time
}

// ExecutionStepResponse represents an execution step response.
type ExecutionStepResponse struct {
	ID                string
	ExecutionRecordID string
	StepNumber        int
	Description       string
	Notes             string
	ExecutedAt        time.Time
}

// VariableValueDTO represents a variable value.
type VariableValueDTO struct {
	Name  string
	Value interface{}
}

// AddStepRequest represents the request to add a step.
type AddStepRequest struct {
	ExecutionRecordID string
	StepNumber        int
	Description       string
}

// UpdateStepNotesRequest represents the request to update step notes.
type UpdateStepNotesRequest struct {
	ExecutionRecordID string
	StepNumber        int
	Notes             string
}

// UpdateNotesRequest represents the request to update overall notes.
type UpdateNotesRequest struct {
	ExecutionRecordID string
	Notes             string
}

// UpdateTitleRequest represents the request to update the title.
type UpdateTitleRequest struct {
	ExecutionRecordID string
	Title             string
}

// UpdateAccessScopeRequest represents the request to update access scope.
type UpdateAccessScopeRequest struct {
	ExecutionRecordID string
	AccessScope       string
}

// CompleteExecutionRequest represents the request to complete an execution.
type CompleteExecutionRequest struct {
	ExecutionRecordID string
}

// MarkAsFailedRequest represents the request to mark an execution as failed.
type MarkAsFailedRequest struct {
	ExecutionRecordID string
}

// SearchExecutionRecordRequest represents the search criteria for execution records.
type SearchExecutionRecordRequest struct {
	ExecutorID      *string
	DocumentID      *string
	Status          *string
	StartedFrom     *time.Time
	StartedTo       *time.Time
	VariableFilters []VariableFilterDTO
}

// VariableFilterDTO represents a variable filter for searching.
type VariableFilterDTO struct {
	Name  string
	Value interface{}
}
