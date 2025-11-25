package usecase

import (
	"context"
	"fmt"
	"io"
	"path/filepath"

	apperror "opscore/backend/internal/execution_record/application/error"
	"opscore/backend/internal/execution_record/application/dto"
	"opscore/backend/internal/execution_record/domain/entity"
	"opscore/backend/internal/execution_record/domain/repository"
	"opscore/backend/internal/execution_record/domain/value_object"
	"opscore/backend/internal/execution_record/infrastructure/storage"
)

// AttachmentUsecase handles attachment business logic.
type AttachmentUsecase struct {
	attachmentRepo repository.AttachmentRepository
	recordRepo     repository.ExecutionRecordRepository
	storageManager storage.StorageManager
}

// NewAttachmentUsecase creates a new AttachmentUsecase.
func NewAttachmentUsecase(
	attachmentRepo repository.AttachmentRepository,
	recordRepo repository.ExecutionRecordRepository,
	storageManager storage.StorageManager,
) *AttachmentUsecase {
	return &AttachmentUsecase{
		attachmentRepo: attachmentRepo,
		recordRepo:     recordRepo,
		storageManager: storageManager,
	}
}

// UploadAttachment uploads a new attachment.
func (uc *AttachmentUsecase) UploadAttachment(
	ctx context.Context,
	req *dto.UploadAttachmentRequest,
) (*dto.AttachmentResponse, error) {
	// Validate execution record ID
	recordID, err := value_object.NewExecutionRecordID(req.ExecutionRecordID)
	if err != nil {
		return nil, &apperror.ValidationError{
			Field:   "executionRecordID",
			Message: "invalid execution record ID format",
		}
	}

	// Validate execution step ID
	stepID, err := value_object.NewExecutionStepID(req.ExecutionStepID)
	if err != nil {
		return nil, &apperror.ValidationError{
			Field:   "executionStepID",
			Message: "invalid execution step ID format",
		}
	}

	// Verify execution record exists
	record, err := uc.recordRepo.FindByID(ctx, recordID)
	if err != nil {
		return nil, err
	}
	if record == nil {
		return nil, &apperror.NotFoundError{
			ResourceType: "ExecutionRecord",
			ResourceID:   req.ExecutionRecordID,
		}
	}

	// Generate attachment ID
	attachmentID := value_object.GenerateAttachmentID()

	// Generate storage path
	storagePath := filepath.Join(
		recordID.String(),
		stepID.String(),
		attachmentID.String()+filepath.Ext(req.FileName),
	)

	// Store the file
	_, err = uc.storageManager.Store(ctx, storagePath, req.File)
	if err != nil {
		return nil, fmt.Errorf("failed to store file: %w", err)
	}

	// Create the attachment entity
	attachment, err := entity.NewAttachment(
		attachmentID,
		recordID,
		stepID,
		req.FileName,
		req.FileSize,
		req.MimeType,
		uc.storageManager.Type(),
		storagePath,
		req.UploadedBy,
	)
	if err != nil {
		// Clean up the stored file
		_ = uc.storageManager.Delete(ctx, storagePath)
		return nil, &apperror.ValidationError{
			Field:   "attachment",
			Message: err.Error(),
		}
	}

	// Save to repository (without the file, as it's already stored)
	if err := uc.attachmentRepo.Save(ctx, attachment, nil); err != nil {
		// Clean up the stored file
		_ = uc.storageManager.Delete(ctx, storagePath)
		return nil, err
	}

	return toAttachmentResponse(attachment), nil
}

// GetAttachment retrieves an attachment by ID.
func (uc *AttachmentUsecase) GetAttachment(
	ctx context.Context,
	attachmentID string,
) (*dto.AttachmentResponse, error) {
	id, err := value_object.NewAttachmentID(attachmentID)
	if err != nil {
		return nil, &apperror.ValidationError{
			Field:   "attachmentID",
			Message: "invalid attachment ID format",
		}
	}

	attachment, err := uc.attachmentRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if attachment == nil {
		return nil, &apperror.NotFoundError{
			ResourceType: "Attachment",
			ResourceID:   attachmentID,
		}
	}

	return toAttachmentResponse(attachment), nil
}

// GetAttachmentFile retrieves the file content for an attachment.
func (uc *AttachmentUsecase) GetAttachmentFile(
	ctx context.Context,
	attachmentID string,
) (io.ReadCloser, *dto.AttachmentResponse, error) {
	id, err := value_object.NewAttachmentID(attachmentID)
	if err != nil {
		return nil, nil, &apperror.ValidationError{
			Field:   "attachmentID",
			Message: "invalid attachment ID format",
		}
	}

	attachment, err := uc.attachmentRepo.FindByID(ctx, id)
	if err != nil {
		return nil, nil, err
	}
	if attachment == nil {
		return nil, nil, &apperror.NotFoundError{
			ResourceType: "Attachment",
			ResourceID:   attachmentID,
		}
	}

	// Retrieve the file from storage
	file, err := uc.storageManager.Retrieve(ctx, attachment.StoragePath())
	if err != nil {
		return nil, nil, fmt.Errorf("failed to retrieve file: %w", err)
	}

	return file, toAttachmentResponse(attachment), nil
}

// ListAttachmentsByRecordID lists attachments by execution record ID.
func (uc *AttachmentUsecase) ListAttachmentsByRecordID(
	ctx context.Context,
	recordID string,
) ([]*dto.AttachmentResponse, error) {
	id, err := value_object.NewExecutionRecordID(recordID)
	if err != nil {
		return nil, &apperror.ValidationError{
			Field:   "recordID",
			Message: "invalid execution record ID format",
		}
	}

	attachments, err := uc.attachmentRepo.FindByExecutionRecordID(ctx, id)
	if err != nil {
		return nil, err
	}

	responses := make([]*dto.AttachmentResponse, len(attachments))
	for i, attachment := range attachments {
		responses[i] = toAttachmentResponse(attachment)
	}

	return responses, nil
}

// ListAttachmentsByStepID lists attachments by execution step ID.
func (uc *AttachmentUsecase) ListAttachmentsByStepID(
	ctx context.Context,
	stepID string,
) ([]*dto.AttachmentResponse, error) {
	id, err := value_object.NewExecutionStepID(stepID)
	if err != nil {
		return nil, &apperror.ValidationError{
			Field:   "stepID",
			Message: "invalid execution step ID format",
		}
	}

	attachments, err := uc.attachmentRepo.FindByExecutionStepID(ctx, id)
	if err != nil {
		return nil, err
	}

	responses := make([]*dto.AttachmentResponse, len(attachments))
	for i, attachment := range attachments {
		responses[i] = toAttachmentResponse(attachment)
	}

	return responses, nil
}

// DeleteAttachment deletes an attachment.
func (uc *AttachmentUsecase) DeleteAttachment(
	ctx context.Context,
	attachmentID string,
) error {
	id, err := value_object.NewAttachmentID(attachmentID)
	if err != nil {
		return &apperror.ValidationError{
			Field:   "attachmentID",
			Message: "invalid attachment ID format",
		}
	}

	attachment, err := uc.attachmentRepo.FindByID(ctx, id)
	if err != nil {
		return err
	}
	if attachment == nil {
		return &apperror.NotFoundError{
			ResourceType: "Attachment",
			ResourceID:   attachmentID,
		}
	}

	// Delete from storage
	if err := uc.storageManager.Delete(ctx, attachment.StoragePath()); err != nil {
		return fmt.Errorf("failed to delete file from storage: %w", err)
	}

	// Delete from repository
	return uc.attachmentRepo.Delete(ctx, id)
}

// Helper function to convert entity to DTO response
func toAttachmentResponse(attachment entity.Attachment) *dto.AttachmentResponse {
	return &dto.AttachmentResponse{
		ID:                attachment.ID().String(),
		ExecutionRecordID: attachment.ExecutionRecordID().String(),
		ExecutionStepID:   attachment.ExecutionStepID().String(),
		FileName:          attachment.FileName(),
		FileSize:          attachment.FileSize(),
		MimeType:          attachment.MimeType(),
		StorageType:       attachment.StorageType().String(),
		StoragePath:       attachment.StoragePath(),
		UploadedBy:        attachment.UploadedBy(),
		UploadedAt:        attachment.UploadedAt(),
	}
}
