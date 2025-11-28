package error

import "errors"

var (
	// ErrExecutionRecordNotFound is returned when an execution record is not found.
	ErrExecutionRecordNotFound = errors.New("execution record not found")

	// ErrExecutionStepNotFound is returned when an execution step is not found.
	ErrExecutionStepNotFound = errors.New("execution step not found")

	// ErrAttachmentNotFound is returned when an attachment is not found.
	ErrAttachmentNotFound = errors.New("attachment not found")

	// ErrInvalidExecutionRecordID is returned when an execution record ID is invalid.
	ErrInvalidExecutionRecordID = errors.New("invalid execution record ID")

	// ErrInvalidExecutionStepID is returned when an execution step ID is invalid.
	ErrInvalidExecutionStepID = errors.New("invalid execution step ID")

	// ErrInvalidAttachmentID is returned when an attachment ID is invalid.
	ErrInvalidAttachmentID = errors.New("invalid attachment ID")

	// ErrInvalidStatusTransition is returned when an invalid status transition is attempted.
	ErrInvalidStatusTransition = errors.New("invalid status transition")

	// ErrDuplicateStepNumber is returned when a step with the same number already exists.
	ErrDuplicateStepNumber = errors.New("step with this number already exists")

	// ErrExecutionNotInProgress is returned when an operation requires an in-progress execution.
	ErrExecutionNotInProgress = errors.New("execution is not in progress")

	// ErrEmptyTitle is returned when a title is empty.
	ErrEmptyTitle = errors.New("title cannot be empty")

	// ErrEmptyDescription is returned when a description is empty.
	ErrEmptyDescription = errors.New("description cannot be empty")

	// ErrInvalidStepNumber is returned when a step number is invalid.
	ErrInvalidStepNumber = errors.New("step number must be positive")

	// ErrInvalidFileSize is returned when a file size is invalid.
	ErrInvalidFileSize = errors.New("file size must be positive")

	// ErrEmptyFileName is returned when a file name is empty.
	ErrEmptyFileName = errors.New("file name cannot be empty")

	// ErrEmptyMimeType is returned when a MIME type is empty.
	ErrEmptyMimeType = errors.New("MIME type cannot be empty")

	// ErrEmptyStoragePath is returned when a storage path is empty.
	ErrEmptyStoragePath = errors.New("storage path cannot be empty")

	// ErrInvalidStorageType is returned when a storage type is invalid.
	ErrInvalidStorageType = errors.New("invalid storage type")
)
