package schema

import "time"

// CreateExecutionRecordRequest represents the API request to create an execution record.
type CreateExecutionRecordRequest struct {
	DocumentID        string                       `json:"document_id" binding:"required"`
	DocumentVersionID string                       `json:"document_version_id" binding:"required"`
	Title             string                       `json:"title" binding:"required"`
	VariableValues    []VariableValueRequestSchema `json:"variable_values"`
}

// VariableValueRequestSchema represents a variable value in API requests.
type VariableValueRequestSchema struct {
	Name  string      `json:"name" binding:"required"`
	Value interface{} `json:"value"`
}

// ExecutionRecordResponse represents the API response for an execution record.
type ExecutionRecordResponse struct {
	ID                string                        `json:"id"`
	DocumentID        string                        `json:"document_id"`
	DocumentVersionID string                        `json:"document_version_id"`
	ExecutorID        string                        `json:"executor_id"`
	Title             string                        `json:"title"`
	VariableValues    []VariableValueResponseSchema `json:"variable_values"`
	Notes             string                        `json:"notes"`
	Status            string                        `json:"status"`
	AccessScope       string                        `json:"access_scope"`
	Steps             []ExecutionStepResponseSchema `json:"steps"`
	StartedAt         time.Time                     `json:"started_at"`
	CompletedAt       *time.Time                    `json:"completed_at,omitempty"`
	CreatedAt         time.Time                     `json:"created_at"`
	UpdatedAt         time.Time                     `json:"updated_at"`
}

// VariableValueResponseSchema represents a variable value in API responses.
type VariableValueResponseSchema struct {
	Name  string      `json:"name"`
	Value interface{} `json:"value"`
}

// ExecutionStepResponseSchema represents an execution step in API responses.
type ExecutionStepResponseSchema struct {
	ID                string    `json:"id"`
	ExecutionRecordID string    `json:"execution_record_id"`
	StepNumber        int       `json:"step_number"`
	Description       string    `json:"description"`
	Notes             string    `json:"notes"`
	ExecutedAt        time.Time `json:"executed_at"`
}

// AddStepRequest represents the API request to add a step.
type AddStepRequest struct {
	StepNumber  int    `json:"step_number" binding:"required,min=1"`
	Description string `json:"description" binding:"required"`
}

// UpdateStepNotesRequest represents the API request to update step notes.
type UpdateStepNotesRequest struct {
	Notes string `json:"notes"`
}

// UpdateNotesRequest represents the API request to update overall notes.
type UpdateNotesRequest struct {
	Notes string `json:"notes"`
}

// UpdateTitleRequest represents the API request to update the title.
type UpdateTitleRequest struct {
	Title string `json:"title" binding:"required"`
}

// UpdateAccessScopeRequest represents the API request to update access scope.
type UpdateAccessScopeRequest struct {
	AccessScope string `json:"access_scope" binding:"required,oneof=public private"`
}

// SearchExecutionRecordRequest represents the API request to search execution records.
type SearchExecutionRecordRequest struct {
	ExecutorID  *string    `form:"executor_id"`
	DocumentID  *string    `form:"document_id"`
	Status      *string    `form:"status"`
	StartedFrom *time.Time `form:"started_from" time_format:"2006-01-02T15:04:05Z07:00"`
	StartedTo   *time.Time `form:"started_to" time_format:"2006-01-02T15:04:05Z07:00"`
}
