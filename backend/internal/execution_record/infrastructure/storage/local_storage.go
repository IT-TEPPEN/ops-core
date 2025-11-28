package storage

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"opscore/backend/internal/execution_record/domain/value_object"
)

// LocalStorageManager implements StorageManager for local file system storage.
type LocalStorageManager struct {
	basePath string
}

// NewLocalStorageManager creates a new LocalStorageManager.
func NewLocalStorageManager(basePath string) (*LocalStorageManager, error) {
	// Ensure the base path exists
	if err := os.MkdirAll(basePath, 0755); err != nil {
		return nil, fmt.Errorf("failed to create storage directory: %w", err)
	}

	return &LocalStorageManager{
		basePath: basePath,
	}, nil
}

// Store saves a file to the local file system.
func (l *LocalStorageManager) Store(ctx context.Context, path string, file io.Reader) (string, error) {
	fullPath := filepath.Join(l.basePath, path)

	// Create the directory structure if it doesn't exist
	dir := filepath.Dir(fullPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return "", fmt.Errorf("failed to create directory: %w", err)
	}

	// Create the file
	f, err := os.Create(fullPath)
	if err != nil {
		return "", fmt.Errorf("failed to create file: %w", err)
	}
	defer f.Close()

	// Copy the content
	if _, err := io.Copy(f, file); err != nil {
		return "", fmt.Errorf("failed to write file: %w", err)
	}

	return path, nil
}

// Retrieve retrieves a file from the local file system.
func (l *LocalStorageManager) Retrieve(ctx context.Context, path string) (io.ReadCloser, error) {
	fullPath := filepath.Join(l.basePath, path)

	file, err := os.Open(fullPath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("file not found: %s", path)
		}
		return nil, fmt.Errorf("failed to open file: %w", err)
	}

	return file, nil
}

// Delete deletes a file from the local file system.
func (l *LocalStorageManager) Delete(ctx context.Context, path string) error {
	fullPath := filepath.Join(l.basePath, path)

	if err := os.Remove(fullPath); err != nil {
		if os.IsNotExist(err) {
			return nil // File already deleted
		}
		return fmt.Errorf("failed to delete file: %w", err)
	}

	return nil
}

// Type returns the storage type.
func (l *LocalStorageManager) Type() value_object.StorageType {
	return value_object.StorageTypeLocal
}
