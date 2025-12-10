package usecase

import (
	"bytes"
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"

	docvo "opscore/backend/internal/document/domain/value_object"
	apperror "opscore/backend/internal/execution_record/application/error"
	"opscore/backend/internal/execution_record/application/dto"
	"opscore/backend/internal/execution_record/domain/entity"
	"opscore/backend/internal/execution_record/domain/value_object"
)

func TestAttachmentUsecase_UploadAttachment(t *testing.T) {
	ctx := context.Background()
	docID := docvo.GenerateDocumentID()
	versionID := docvo.GenerateVersionID()
	recordID := value_object.GenerateExecutionRecordID()
	stepID := value_object.GenerateExecutionStepID()

	record, _ := entity.NewExecutionRecord(
		recordID,
		docID,
		versionID,
		"user-123",
		"Test Execution",
		[]value_object.VariableValue{},
	)

	mockRecordRepo := &MockExecutionRecordRepository{
		FindByIDFunc: func(ctx context.Context, id value_object.ExecutionRecordID) (entity.ExecutionRecord, error) {
			if id.Equals(recordID) {
				return record, nil
			}
			return nil, nil
		},
	}

	mockAttachmentRepo := &MockAttachmentRepository{}
	mockStorageManager := &MockStorageManager{}

	uc := NewAttachmentUsecase(mockAttachmentRepo, mockRecordRepo, mockStorageManager)

	fileContent := bytes.NewReader([]byte("test file content"))
	req := &dto.UploadAttachmentRequest{
		ExecutionRecordID: recordID.String(),
		ExecutionStepID:   stepID.String(),
		FileName:          "test.png",
		FileSize:          1024,
		MimeType:          "image/png",
		UploadedBy:        "user-123",
		File:              fileContent,
	}

	resp, err := uc.UploadAttachment(ctx, req)
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, req.FileName, resp.FileName)
	assert.Equal(t, req.FileSize, resp.FileSize)
	assert.Equal(t, req.MimeType, resp.MimeType)
	assert.Equal(t, req.UploadedBy, resp.UploadedBy)
}

func TestAttachmentUsecase_UploadAttachment_InvalidRecordID(t *testing.T) {
	ctx := context.Background()
	mockRecordRepo := &MockExecutionRecordRepository{}
	mockAttachmentRepo := &MockAttachmentRepository{}
	mockStorageManager := &MockStorageManager{}

	uc := NewAttachmentUsecase(mockAttachmentRepo, mockRecordRepo, mockStorageManager)

	req := &dto.UploadAttachmentRequest{
		ExecutionRecordID: "invalid-id",
		ExecutionStepID:   value_object.GenerateExecutionStepID().String(),
		FileName:          "test.png",
		FileSize:          1024,
		MimeType:          "image/png",
		UploadedBy:        "user-123",
		File:              bytes.NewReader([]byte("test")),
	}

	_, err := uc.UploadAttachment(ctx, req)
	assert.Error(t, err)
	var validationErr *apperror.ValidationError
	assert.True(t, errors.As(err, &validationErr))
}

func TestAttachmentUsecase_GetAttachment(t *testing.T) {
	ctx := context.Background()
	attachmentID := value_object.GenerateAttachmentID()
	recordID := value_object.GenerateExecutionRecordID()
	stepID := value_object.GenerateExecutionStepID()

	attachment, _ := entity.NewAttachment(
		attachmentID,
		recordID,
		stepID,
		"test.png",
		1024,
		"image/png",
		value_object.StorageTypeLocal,
		"/path/to/file",
		"user-123",
	)

	mockAttachmentRepo := &MockAttachmentRepository{
		FindByIDFunc: func(ctx context.Context, id value_object.AttachmentID) (entity.Attachment, error) {
			if id.Equals(attachmentID) {
				return attachment, nil
			}
			return nil, nil
		},
	}
	mockRecordRepo := &MockExecutionRecordRepository{}
	mockStorageManager := &MockStorageManager{}

	uc := NewAttachmentUsecase(mockAttachmentRepo, mockRecordRepo, mockStorageManager)

	resp, err := uc.GetAttachment(ctx, attachmentID.String())
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, attachmentID.String(), resp.ID)
	assert.Equal(t, "test.png", resp.FileName)
}

func TestAttachmentUsecase_GetAttachment_NotFound(t *testing.T) {
	ctx := context.Background()
	attachmentID := value_object.GenerateAttachmentID()

	mockAttachmentRepo := &MockAttachmentRepository{
		FindByIDFunc: func(ctx context.Context, id value_object.AttachmentID) (entity.Attachment, error) {
			return nil, nil
		},
	}
	mockRecordRepo := &MockExecutionRecordRepository{}
	mockStorageManager := &MockStorageManager{}

	uc := NewAttachmentUsecase(mockAttachmentRepo, mockRecordRepo, mockStorageManager)

	_, err := uc.GetAttachment(ctx, attachmentID.String())
	assert.Error(t, err)
	var notFoundErr *apperror.NotFoundError
	assert.True(t, errors.As(err, &notFoundErr))
}

func TestAttachmentUsecase_DeleteAttachment(t *testing.T) {
	ctx := context.Background()
	attachmentID := value_object.GenerateAttachmentID()
	recordID := value_object.GenerateExecutionRecordID()
	stepID := value_object.GenerateExecutionStepID()

	attachment, _ := entity.NewAttachment(
		attachmentID,
		recordID,
		stepID,
		"test.png",
		1024,
		"image/png",
		value_object.StorageTypeLocal,
		"/path/to/file",
		"user-123",
	)

	mockAttachmentRepo := &MockAttachmentRepository{
		FindByIDFunc: func(ctx context.Context, id value_object.AttachmentID) (entity.Attachment, error) {
			if id.Equals(attachmentID) {
				return attachment, nil
			}
			return nil, nil
		},
	}
	mockRecordRepo := &MockExecutionRecordRepository{}
	mockStorageManager := &MockStorageManager{}

	uc := NewAttachmentUsecase(mockAttachmentRepo, mockRecordRepo, mockStorageManager)

	err := uc.DeleteAttachment(ctx, attachmentID.String())
	assert.NoError(t, err)
}

func TestAttachmentUsecase_ListAttachmentsByRecordID(t *testing.T) {
	ctx := context.Background()
	recordID := value_object.GenerateExecutionRecordID()
	stepID := value_object.GenerateExecutionStepID()

	attachment1, _ := entity.NewAttachment(
		value_object.GenerateAttachmentID(),
		recordID,
		stepID,
		"test1.png",
		1024,
		"image/png",
		value_object.StorageTypeLocal,
		"/path/to/file1",
		"user-123",
	)

	attachment2, _ := entity.NewAttachment(
		value_object.GenerateAttachmentID(),
		recordID,
		stepID,
		"test2.png",
		2048,
		"image/png",
		value_object.StorageTypeLocal,
		"/path/to/file2",
		"user-123",
	)

	mockAttachmentRepo := &MockAttachmentRepository{
		FindByExecutionRecordIDFunc: func(ctx context.Context, id value_object.ExecutionRecordID) ([]entity.Attachment, error) {
			if id.Equals(recordID) {
				return []entity.Attachment{attachment1, attachment2}, nil
			}
			return nil, nil
		},
	}
	mockRecordRepo := &MockExecutionRecordRepository{}
	mockStorageManager := &MockStorageManager{}

	uc := NewAttachmentUsecase(mockAttachmentRepo, mockRecordRepo, mockStorageManager)

	resp, err := uc.ListAttachmentsByRecordID(ctx, recordID.String())
	assert.NoError(t, err)
	assert.Len(t, resp, 2)
	assert.Equal(t, "test1.png", resp[0].FileName)
	assert.Equal(t, "test2.png", resp[1].FileName)
}

func TestAttachmentUsecase_GetAttachmentURL(t *testing.T) {
	ctx := context.Background()
	attachmentID := value_object.GenerateAttachmentID()
	recordID := value_object.GenerateExecutionRecordID()
	stepID := value_object.GenerateExecutionStepID()

	attachment, _ := entity.NewAttachment(
		attachmentID,
		recordID,
		stepID,
		"test.png",
		1024,
		"image/png",
		value_object.StorageTypeLocal,
		"/path/to/file",
		"user-123",
	)

	mockAttachmentRepo := &MockAttachmentRepository{
		FindByIDFunc: func(ctx context.Context, id value_object.AttachmentID) (entity.Attachment, error) {
			if id.Equals(attachmentID) {
				return attachment, nil
			}
			return nil, nil
		},
	}
	mockRecordRepo := &MockExecutionRecordRepository{}
	mockStorageManager := &MockStorageManager{
		GeneratePresignedURLFunc: func(ctx context.Context, path string, expirationMinutes int) (string, error) {
			return "", nil // Local storage returns empty string
		},
	}

	uc := NewAttachmentUsecase(mockAttachmentRepo, mockRecordRepo, mockStorageManager)

	url, err := uc.GetAttachmentURL(ctx, attachmentID.String(), 60)
	assert.NoError(t, err)
	assert.Equal(t, "", url) // Local storage returns empty string
}
