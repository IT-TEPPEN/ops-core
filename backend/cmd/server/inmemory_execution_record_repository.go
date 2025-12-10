package main

import (
	"context"
	"errors"
	"sync"

	docvo "opscore/backend/internal/document/domain/value_object"
	"opscore/backend/internal/execution_record/domain/entity"
	"opscore/backend/internal/execution_record/domain/repository"
	"opscore/backend/internal/execution_record/domain/value_object"
)

// InMemoryExecutionRecordRepository is an in-memory implementation of ExecutionRecordRepository.
type InMemoryExecutionRecordRepository struct {
	mu      sync.RWMutex
	records map[string]entity.ExecutionRecord
}

// NewInMemoryExecutionRecordRepository creates a new InMemoryExecutionRecordRepository.
func NewInMemoryExecutionRecordRepository() repository.ExecutionRecordRepository {
	return &InMemoryExecutionRecordRepository{
		records: make(map[string]entity.ExecutionRecord),
	}
}

// Save saves a new execution record.
func (r *InMemoryExecutionRecordRepository) Save(ctx context.Context, record entity.ExecutionRecord) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.records[record.ID().String()] = record
	return nil
}

// FindByID retrieves an execution record by ID.
func (r *InMemoryExecutionRecordRepository) FindByID(ctx context.Context, id value_object.ExecutionRecordID) (entity.ExecutionRecord, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	record, exists := r.records[id.String()]
	if !exists {
		return nil, nil
	}
	return record, nil
}

// FindByExecutorID retrieves execution records by executor ID.
func (r *InMemoryExecutionRecordRepository) FindByExecutorID(ctx context.Context, executorID string) ([]entity.ExecutionRecord, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var results []entity.ExecutionRecord
	for _, record := range r.records {
		if record.ExecutorID() == executorID {
			results = append(results, record)
		}
	}
	return results, nil
}

// FindByDocumentID retrieves execution records by document ID.
func (r *InMemoryExecutionRecordRepository) FindByDocumentID(ctx context.Context, documentID docvo.DocumentID) ([]entity.ExecutionRecord, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var results []entity.ExecutionRecord
	for _, record := range r.records {
		if record.DocumentID().String() == documentID.String() {
			results = append(results, record)
		}
	}
	return results, nil
}

// Search searches for execution records based on criteria.
func (r *InMemoryExecutionRecordRepository) Search(ctx context.Context, criteria repository.SearchCriteria) ([]entity.ExecutionRecord, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var results []entity.ExecutionRecord
	for _, record := range r.records {
		// Filter by executor ID
		if criteria.ExecutorID != nil && record.ExecutorID() != *criteria.ExecutorID {
			continue
		}

		// Filter by document ID
		if criteria.DocumentID != nil && record.DocumentID().String() != criteria.DocumentID.String() {
			continue
		}

		// Filter by status
		if criteria.Status != nil && record.Status().String() != criteria.Status.String() {
			continue
		}

		// Filter by start date range
		if criteria.StartedFrom != nil && record.StartedAt().Before(*criteria.StartedFrom) {
			continue
		}
		if criteria.StartedTo != nil && record.StartedAt().After(*criteria.StartedTo) {
			continue
		}

		// Filter by variable values
		if len(criteria.VariableFilters) > 0 {
			match := true
			for _, filter := range criteria.VariableFilters {
				found := false
				for _, vv := range record.VariableValues() {
					if vv.Name() == filter.Name && vv.Value() == filter.Value {
						found = true
						break
					}
				}
				if !found {
					match = false
					break
				}
			}
			if !match {
				continue
			}
		}

		results = append(results, record)
	}
	return results, nil
}

// Update updates an existing execution record.
func (r *InMemoryExecutionRecordRepository) Update(ctx context.Context, record entity.ExecutionRecord) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.records[record.ID().String()]; !exists {
		return errors.New("execution record not found")
	}

	r.records[record.ID().String()] = record
	return nil
}

// Delete deletes an execution record by ID.
func (r *InMemoryExecutionRecordRepository) Delete(ctx context.Context, id value_object.ExecutionRecordID) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.records[id.String()]; !exists {
		return errors.New("execution record not found")
	}

	delete(r.records, id.String())
	return nil
}
