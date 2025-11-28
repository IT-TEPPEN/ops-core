package entity

import (
	"errors"
	"time"

	"opscore/backend/internal/execution_record/domain/value_object"
)

// attachment represents an attachment for an execution step.
type attachment struct {
	id                value_object.AttachmentID
	executionRecordID value_object.ExecutionRecordID
	executionStepID   value_object.ExecutionStepID
	fileName          string
	fileSize          int64
	mimeType          string
	storageType       value_object.StorageType
	storagePath       string
	uploadedBy        string // User ID as string
	uploadedAt        time.Time
}

// Attachment is the interface for an attachment.
type Attachment interface {
	ID() value_object.AttachmentID
	ExecutionRecordID() value_object.ExecutionRecordID
	ExecutionStepID() value_object.ExecutionStepID
	FileName() string
	FileSize() int64
	MimeType() string
	StorageType() value_object.StorageType
	StoragePath() string
	UploadedBy() string
	UploadedAt() time.Time
}

// NewAttachment creates a new Attachment instance.
func NewAttachment(
	id value_object.AttachmentID,
	recordID value_object.ExecutionRecordID,
	stepID value_object.ExecutionStepID,
	fileName string,
	fileSize int64,
	mimeType string,
	storageType value_object.StorageType,
	storagePath string,
	uploadedBy string,
) (Attachment, error) {
	if id.IsEmpty() {
		return nil, errors.New("attachment ID cannot be empty")
	}
	if recordID.IsEmpty() {
		return nil, errors.New("execution record ID cannot be empty")
	}
	if stepID.IsEmpty() {
		return nil, errors.New("execution step ID cannot be empty")
	}
	if fileName == "" {
		return nil, errors.New("file name cannot be empty")
	}
	if fileSize <= 0 {
		return nil, errors.New("file size must be positive")
	}
	if mimeType == "" {
		return nil, errors.New("MIME type cannot be empty")
	}
	if !storageType.IsValid() {
		return nil, errors.New("invalid storage type")
	}
	if storagePath == "" {
		return nil, errors.New("storage path cannot be empty")
	}
	if uploadedBy == "" {
		return nil, errors.New("uploader ID cannot be empty")
	}

	return &attachment{
		id:                id,
		executionRecordID: recordID,
		executionStepID:   stepID,
		fileName:          fileName,
		fileSize:          fileSize,
		mimeType:          mimeType,
		storageType:       storageType,
		storagePath:       storagePath,
		uploadedBy:        uploadedBy,
		uploadedAt:        time.Now(),
	}, nil
}

// ReconstructAttachment reconstructs an Attachment from persistence data.
func ReconstructAttachment(
	id value_object.AttachmentID,
	recordID value_object.ExecutionRecordID,
	stepID value_object.ExecutionStepID,
	fileName string,
	fileSize int64,
	mimeType string,
	storageType value_object.StorageType,
	storagePath string,
	uploadedBy string,
	uploadedAt time.Time,
) Attachment {
	return &attachment{
		id:                id,
		executionRecordID: recordID,
		executionStepID:   stepID,
		fileName:          fileName,
		fileSize:          fileSize,
		mimeType:          mimeType,
		storageType:       storageType,
		storagePath:       storagePath,
		uploadedBy:        uploadedBy,
		uploadedAt:        uploadedAt,
	}
}

// Getter methods

// ID returns the attachment ID.
func (a *attachment) ID() value_object.AttachmentID {
	return a.id
}

// ExecutionRecordID returns the execution record ID.
func (a *attachment) ExecutionRecordID() value_object.ExecutionRecordID {
	return a.executionRecordID
}

// ExecutionStepID returns the execution step ID.
func (a *attachment) ExecutionStepID() value_object.ExecutionStepID {
	return a.executionStepID
}

// FileName returns the file name.
func (a *attachment) FileName() string {
	return a.fileName
}

// FileSize returns the file size in bytes.
func (a *attachment) FileSize() int64 {
	return a.fileSize
}

// MimeType returns the MIME type.
func (a *attachment) MimeType() string {
	return a.mimeType
}

// StorageType returns the storage type.
func (a *attachment) StorageType() value_object.StorageType {
	return a.storageType
}

// StoragePath returns the storage path or key.
func (a *attachment) StoragePath() string {
	return a.storagePath
}

// UploadedBy returns the uploader's user ID.
func (a *attachment) UploadedBy() string {
	return a.uploadedBy
}

// UploadedAt returns the upload timestamp.
func (a *attachment) UploadedAt() time.Time {
	return a.uploadedAt
}
