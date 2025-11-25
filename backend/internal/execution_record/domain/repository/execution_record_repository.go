package repository

import (
	"context"
	"time"

	docvo "opscore/backend/internal/document/domain/value_object"
	"opscore/backend/internal/execution_record/domain/entity"
	"opscore/backend/internal/execution_record/domain/value_object"
)

// SearchCriteria represents the search criteria for execution records.
type SearchCriteria struct {
	ExecutorID      *string
	DocumentID      *docvo.DocumentID
	Status          *value_object.ExecutionStatus
	StartedFrom     *time.Time
	StartedTo       *time.Time
	VariableFilters []VariableFilter
}

// VariableFilter represents a filter for variable values.
type VariableFilter struct {
	Name  string
	Value interface{}
}

// ExecutionRecordRepository defines the interface for execution record persistence.
type ExecutionRecordRepository interface {
	// Save saves a new execution record.
	Save(ctx context.Context, record entity.ExecutionRecord) error

	// FindByID retrieves an execution record by ID.
	FindByID(ctx context.Context, id value_object.ExecutionRecordID) (entity.ExecutionRecord, error)

	// FindByExecutorID retrieves execution records by executor ID.
	FindByExecutorID(ctx context.Context, executorID string) ([]entity.ExecutionRecord, error)

	// FindByDocumentID retrieves execution records by document ID.
	FindByDocumentID(ctx context.Context, documentID docvo.DocumentID) ([]entity.ExecutionRecord, error)

	// Search searches for execution records based on criteria.
	Search(ctx context.Context, criteria SearchCriteria) ([]entity.ExecutionRecord, error)

	// Update updates an existing execution record.
	Update(ctx context.Context, record entity.ExecutionRecord) error

	// Delete deletes an execution record by ID.
	Delete(ctx context.Context, id value_object.ExecutionRecordID) error
}
