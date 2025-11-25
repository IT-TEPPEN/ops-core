package repository

import (
	"context"
	"io"

	"opscore/backend/internal/execution_record/domain/entity"
	"opscore/backend/internal/execution_record/domain/value_object"
)

// AttachmentRepository defines the interface for attachment persistence.
type AttachmentRepository interface {
	// Save saves a new attachment with its file content.
	Save(ctx context.Context, attachment entity.Attachment, file io.Reader) error

	// FindByID retrieves an attachment by ID.
	FindByID(ctx context.Context, id value_object.AttachmentID) (entity.Attachment, error)

	// FindByExecutionRecordID retrieves attachments by execution record ID.
	FindByExecutionRecordID(ctx context.Context, recordID value_object.ExecutionRecordID) ([]entity.Attachment, error)

	// FindByExecutionStepID retrieves attachments by execution step ID.
	FindByExecutionStepID(ctx context.Context, stepID value_object.ExecutionStepID) ([]entity.Attachment, error)

	// GetFile retrieves the file content for an attachment.
	GetFile(ctx context.Context, id value_object.AttachmentID) (io.ReadCloser, error)

	// Delete deletes an attachment by ID.
	Delete(ctx context.Context, id value_object.AttachmentID) error
}
