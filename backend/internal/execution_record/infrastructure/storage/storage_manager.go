package storage

import (
	"context"
	"io"

	"opscore/backend/internal/execution_record/domain/value_object"
)

// StorageManager defines the interface for file storage operations.
type StorageManager interface {
	// Store saves a file and returns the storage path.
	Store(ctx context.Context, path string, file io.Reader) (string, error)

	// Retrieve retrieves a file by its path.
	Retrieve(ctx context.Context, path string) (io.ReadCloser, error)

	// Delete deletes a file by its path.
	Delete(ctx context.Context, path string) error

	// GeneratePresignedURL generates a presigned URL for accessing a file.
	// Returns empty string if not supported (e.g., local storage).
	GeneratePresignedURL(ctx context.Context, path string, expirationMinutes int) (string, error)

	// Type returns the storage type.
	Type() value_object.StorageType
}
