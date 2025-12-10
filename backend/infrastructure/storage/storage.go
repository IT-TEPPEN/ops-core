package storage

import (
	"context"
	"io"
	"time"
)

// Storage defines the interface for file storage operations.
// This abstraction supports local file system, S3, and MinIO storage.
type Storage interface {
	// Save stores a file with the given key and content type.
	// Returns an error if the operation fails.
	Save(ctx context.Context, key string, reader io.Reader, contentType string) error

	// Get retrieves a file by its key.
	// Returns io.ReadCloser that must be closed by the caller.
	// Returns an error if the file is not found or operation fails.
	Get(ctx context.Context, key string) (io.ReadCloser, error)

	// Delete removes a file by its key.
	// Returns an error if the operation fails.
	// Returns nil if the file doesn't exist (idempotent).
	Delete(ctx context.Context, key string) error

	// Exists checks if a file exists for the given key.
	// Returns true if the file exists, false otherwise.
	Exists(ctx context.Context, key string) (bool, error)

	// GetSignedURL generates a pre-signed URL for temporary file access.
	// For local storage, this may return a local path or empty string.
	// For S3/MinIO, this returns a time-limited URL.
	// expiration specifies how long the URL should be valid.
	GetSignedURL(ctx context.Context, key string, expiration time.Duration) (string, error)
}
