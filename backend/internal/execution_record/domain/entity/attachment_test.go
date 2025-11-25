package entity

import (
	"testing"
	"time"

	"opscore/backend/internal/execution_record/domain/value_object"
)

func TestNewAttachment(t *testing.T) {
	validID := value_object.GenerateAttachmentID()
	validRecordID := value_object.GenerateExecutionRecordID()
	validStepID := value_object.GenerateExecutionStepID()

	tests := []struct {
		name        string
		id          value_object.AttachmentID
		recordID    value_object.ExecutionRecordID
		stepID      value_object.ExecutionStepID
		fileName    string
		fileSize    int64
		mimeType    string
		storageType value_object.StorageType
		storagePath string
		uploadedBy  string
		wantErr     bool
	}{
		{
			name:        "valid attachment",
			id:          validID,
			recordID:    validRecordID,
			stepID:      validStepID,
			fileName:    "screenshot.png",
			fileSize:    1024,
			mimeType:    "image/png",
			storageType: value_object.StorageTypeLocal,
			storagePath: "/data/attachments/file.png",
			uploadedBy:  "user-123",
			wantErr:     false,
		},
		{
			name:        "empty attachment ID",
			id:          value_object.AttachmentID(""),
			recordID:    validRecordID,
			stepID:      validStepID,
			fileName:    "screenshot.png",
			fileSize:    1024,
			mimeType:    "image/png",
			storageType: value_object.StorageTypeLocal,
			storagePath: "/data/attachments/file.png",
			uploadedBy:  "user-123",
			wantErr:     true,
		},
		{
			name:        "empty record ID",
			id:          validID,
			recordID:    value_object.ExecutionRecordID(""),
			stepID:      validStepID,
			fileName:    "screenshot.png",
			fileSize:    1024,
			mimeType:    "image/png",
			storageType: value_object.StorageTypeLocal,
			storagePath: "/data/attachments/file.png",
			uploadedBy:  "user-123",
			wantErr:     true,
		},
		{
			name:        "empty step ID",
			id:          validID,
			recordID:    validRecordID,
			stepID:      value_object.ExecutionStepID(""),
			fileName:    "screenshot.png",
			fileSize:    1024,
			mimeType:    "image/png",
			storageType: value_object.StorageTypeLocal,
			storagePath: "/data/attachments/file.png",
			uploadedBy:  "user-123",
			wantErr:     true,
		},
		{
			name:        "empty file name",
			id:          validID,
			recordID:    validRecordID,
			stepID:      validStepID,
			fileName:    "",
			fileSize:    1024,
			mimeType:    "image/png",
			storageType: value_object.StorageTypeLocal,
			storagePath: "/data/attachments/file.png",
			uploadedBy:  "user-123",
			wantErr:     true,
		},
		{
			name:        "zero file size",
			id:          validID,
			recordID:    validRecordID,
			stepID:      validStepID,
			fileName:    "screenshot.png",
			fileSize:    0,
			mimeType:    "image/png",
			storageType: value_object.StorageTypeLocal,
			storagePath: "/data/attachments/file.png",
			uploadedBy:  "user-123",
			wantErr:     true,
		},
		{
			name:        "negative file size",
			id:          validID,
			recordID:    validRecordID,
			stepID:      validStepID,
			fileName:    "screenshot.png",
			fileSize:    -1,
			mimeType:    "image/png",
			storageType: value_object.StorageTypeLocal,
			storagePath: "/data/attachments/file.png",
			uploadedBy:  "user-123",
			wantErr:     true,
		},
		{
			name:        "empty MIME type",
			id:          validID,
			recordID:    validRecordID,
			stepID:      validStepID,
			fileName:    "screenshot.png",
			fileSize:    1024,
			mimeType:    "",
			storageType: value_object.StorageTypeLocal,
			storagePath: "/data/attachments/file.png",
			uploadedBy:  "user-123",
			wantErr:     true,
		},
		{
			name:        "invalid storage type",
			id:          validID,
			recordID:    validRecordID,
			stepID:      validStepID,
			fileName:    "screenshot.png",
			fileSize:    1024,
			mimeType:    "image/png",
			storageType: value_object.StorageType("invalid"),
			storagePath: "/data/attachments/file.png",
			uploadedBy:  "user-123",
			wantErr:     true,
		},
		{
			name:        "empty storage path",
			id:          validID,
			recordID:    validRecordID,
			stepID:      validStepID,
			fileName:    "screenshot.png",
			fileSize:    1024,
			mimeType:    "image/png",
			storageType: value_object.StorageTypeLocal,
			storagePath: "",
			uploadedBy:  "user-123",
			wantErr:     true,
		},
		{
			name:        "empty uploader ID",
			id:          validID,
			recordID:    validRecordID,
			stepID:      validStepID,
			fileName:    "screenshot.png",
			fileSize:    1024,
			mimeType:    "image/png",
			storageType: value_object.StorageTypeLocal,
			storagePath: "/data/attachments/file.png",
			uploadedBy:  "",
			wantErr:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewAttachment(
				tt.id,
				tt.recordID,
				tt.stepID,
				tt.fileName,
				tt.fileSize,
				tt.mimeType,
				tt.storageType,
				tt.storagePath,
				tt.uploadedBy,
			)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewAttachment() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if got == nil {
					t.Error("NewAttachment() returned nil")
					return
				}
				if !got.ID().Equals(tt.id) {
					t.Errorf("ID() = %v, want %v", got.ID(), tt.id)
				}
				if !got.ExecutionRecordID().Equals(tt.recordID) {
					t.Errorf("ExecutionRecordID() = %v, want %v", got.ExecutionRecordID(), tt.recordID)
				}
				if !got.ExecutionStepID().Equals(tt.stepID) {
					t.Errorf("ExecutionStepID() = %v, want %v", got.ExecutionStepID(), tt.stepID)
				}
				if got.FileName() != tt.fileName {
					t.Errorf("FileName() = %v, want %v", got.FileName(), tt.fileName)
				}
				if got.FileSize() != tt.fileSize {
					t.Errorf("FileSize() = %v, want %v", got.FileSize(), tt.fileSize)
				}
				if got.MimeType() != tt.mimeType {
					t.Errorf("MimeType() = %v, want %v", got.MimeType(), tt.mimeType)
				}
				if got.StorageType() != tt.storageType {
					t.Errorf("StorageType() = %v, want %v", got.StorageType(), tt.storageType)
				}
				if got.StoragePath() != tt.storagePath {
					t.Errorf("StoragePath() = %v, want %v", got.StoragePath(), tt.storagePath)
				}
				if got.UploadedBy() != tt.uploadedBy {
					t.Errorf("UploadedBy() = %v, want %v", got.UploadedBy(), tt.uploadedBy)
				}
			}
		})
	}
}

func TestReconstructAttachment(t *testing.T) {
	id := value_object.GenerateAttachmentID()
	recordID := value_object.GenerateExecutionRecordID()
	stepID := value_object.GenerateExecutionStepID()
	fileName := "document.pdf"
	fileSize := int64(2048)
	mimeType := "application/pdf"
	storageType := value_object.StorageTypeS3
	storagePath := "attachments/doc.pdf"
	uploadedBy := "user-789"
	uploadedAt := time.Date(2024, 1, 15, 10, 30, 0, 0, time.UTC)

	attachment := ReconstructAttachment(
		id,
		recordID,
		stepID,
		fileName,
		fileSize,
		mimeType,
		storageType,
		storagePath,
		uploadedBy,
		uploadedAt,
	)

	if !attachment.ID().Equals(id) {
		t.Errorf("ID() = %v, want %v", attachment.ID(), id)
	}
	if !attachment.ExecutionRecordID().Equals(recordID) {
		t.Errorf("ExecutionRecordID() = %v, want %v", attachment.ExecutionRecordID(), recordID)
	}
	if !attachment.ExecutionStepID().Equals(stepID) {
		t.Errorf("ExecutionStepID() = %v, want %v", attachment.ExecutionStepID(), stepID)
	}
	if attachment.FileName() != fileName {
		t.Errorf("FileName() = %v, want %v", attachment.FileName(), fileName)
	}
	if attachment.FileSize() != fileSize {
		t.Errorf("FileSize() = %v, want %v", attachment.FileSize(), fileSize)
	}
	if attachment.MimeType() != mimeType {
		t.Errorf("MimeType() = %v, want %v", attachment.MimeType(), mimeType)
	}
	if attachment.StorageType() != storageType {
		t.Errorf("StorageType() = %v, want %v", attachment.StorageType(), storageType)
	}
	if attachment.StoragePath() != storagePath {
		t.Errorf("StoragePath() = %v, want %v", attachment.StoragePath(), storagePath)
	}
	if attachment.UploadedBy() != uploadedBy {
		t.Errorf("UploadedBy() = %v, want %v", attachment.UploadedBy(), uploadedBy)
	}
	if !attachment.UploadedAt().Equal(uploadedAt) {
		t.Errorf("UploadedAt() = %v, want %v", attachment.UploadedAt(), uploadedAt)
	}
}
