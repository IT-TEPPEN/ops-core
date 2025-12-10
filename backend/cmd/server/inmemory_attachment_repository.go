package main

import (
	"context"
	"io"
	"sync"

	"opscore/backend/internal/execution_record/domain/entity"
	"opscore/backend/internal/execution_record/domain/repository"
	"opscore/backend/internal/execution_record/domain/value_object"
)

// InMemoryAttachmentRepository is an in-memory implementation of AttachmentRepository.
type InMemoryAttachmentRepository struct {
	mu          sync.RWMutex
	attachments map[string]entity.Attachment
}

// NewInMemoryAttachmentRepository creates a new InMemoryAttachmentRepository.
func NewInMemoryAttachmentRepository() repository.AttachmentRepository {
	return &InMemoryAttachmentRepository{
		attachments: make(map[string]entity.Attachment),
	}
}

// Save saves a new attachment with its file content.
// Note: File storage is handled by the StorageManager, not the repository.
func (r *InMemoryAttachmentRepository) Save(ctx context.Context, attachment entity.Attachment, file io.Reader) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.attachments[attachment.ID().String()] = attachment
	return nil
}

// FindByID retrieves an attachment by ID.
func (r *InMemoryAttachmentRepository) FindByID(ctx context.Context, id value_object.AttachmentID) (entity.Attachment, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	attachment, exists := r.attachments[id.String()]
	if !exists {
		return nil, nil
	}
	return attachment, nil
}

// FindByExecutionRecordID retrieves attachments by execution record ID.
func (r *InMemoryAttachmentRepository) FindByExecutionRecordID(ctx context.Context, recordID value_object.ExecutionRecordID) ([]entity.Attachment, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var results []entity.Attachment
	for _, attachment := range r.attachments {
		if attachment.ExecutionRecordID().Equals(recordID) {
			results = append(results, attachment)
		}
	}
	return results, nil
}

// FindByExecutionStepID retrieves attachments by execution step ID.
func (r *InMemoryAttachmentRepository) FindByExecutionStepID(ctx context.Context, stepID value_object.ExecutionStepID) ([]entity.Attachment, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var results []entity.Attachment
	for _, attachment := range r.attachments {
		if attachment.ExecutionStepID().Equals(stepID) {
			results = append(results, attachment)
		}
	}
	return results, nil
}

// GetFile retrieves the file content for an attachment.
// Note: File retrieval is handled by the StorageManager, not the repository.
func (r *InMemoryAttachmentRepository) GetFile(ctx context.Context, id value_object.AttachmentID) (io.ReadCloser, error) {
	// This method is not used in the current implementation
	// File retrieval is done through the StorageManager in the usecase layer
	return nil, nil
}

// Delete deletes an attachment by ID.
func (r *InMemoryAttachmentRepository) Delete(ctx context.Context, id value_object.AttachmentID) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	delete(r.attachments, id.String())
	return nil
}
