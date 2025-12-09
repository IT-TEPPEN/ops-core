package storage

import (
"context"
"errors"
"fmt"
"io"
"os"
"path/filepath"
"strings"
"time"
)

// LocalStorage implements Storage interface for local file system.
type LocalStorage struct {
basePath string
}

// NewLocalStorage creates a new LocalStorage instance.
// basePath is the root directory where files will be stored.
// The directory will be created if it doesn't exist.
func NewLocalStorage(basePath string) (*LocalStorage, error) {
if basePath == "" {
return nil, errors.New("base path cannot be empty")
}

// Create the base directory if it doesn't exist
if err := os.MkdirAll(basePath, 0755); err != nil {
return nil, fmt.Errorf("failed to create base directory: %w", err)
}

// Get absolute path to prevent issues
absPath, err := filepath.Abs(basePath)
if err != nil {
return nil, fmt.Errorf("failed to resolve absolute path: %w", err)
}

return &LocalStorage{
basePath: absPath,
}, nil
}

// Save stores a file in the local file system.
func (l *LocalStorage) Save(ctx context.Context, key string, reader io.Reader, contentType string) error {
if key == "" {
return errors.New("key cannot be empty")
}
if reader == nil {
return errors.New("reader cannot be nil")
}

// Validate and sanitize the key to prevent path traversal attacks
fullPath, err := l.resolvePath(key)
if err != nil {
return err
}

// Create directory structure if it doesn't exist
dir := filepath.Dir(fullPath)
if err := os.MkdirAll(dir, 0755); err != nil {
return fmt.Errorf("failed to create directory: %w", err)
}

// Create the file
file, err := os.Create(fullPath)
if err != nil {
return fmt.Errorf("failed to create file: %w", err)
}
defer file.Close()

// Copy content from reader to file
if _, err := io.Copy(file, reader); err != nil {
return fmt.Errorf("failed to write file: %w", err)
}

return nil
}

// Get retrieves a file from the local file system.
func (l *LocalStorage) Get(ctx context.Context, key string) (io.ReadCloser, error) {
if key == "" {
return nil, errors.New("key cannot be empty")
}

fullPath, err := l.resolvePath(key)
if err != nil {
return nil, err
}

file, err := os.Open(fullPath)
if err != nil {
if os.IsNotExist(err) {
return nil, fmt.Errorf("file not found: %s", key)
}
return nil, fmt.Errorf("failed to open file: %w", err)
}

return file, nil
}

// Delete removes a file from the local file system.
func (l *LocalStorage) Delete(ctx context.Context, key string) error {
if key == "" {
return errors.New("key cannot be empty")
}

fullPath, err := l.resolvePath(key)
if err != nil {
return err
}

if err := os.Remove(fullPath); err != nil {
if os.IsNotExist(err) {
// File doesn't exist, which is fine (idempotent)
return nil
}
return fmt.Errorf("failed to delete file: %w", err)
}

return nil
}

// Exists checks if a file exists in the local file system.
func (l *LocalStorage) Exists(ctx context.Context, key string) (bool, error) {
if key == "" {
return false, errors.New("key cannot be empty")
}

fullPath, err := l.resolvePath(key)
if err != nil {
return false, err
}

_, err = os.Stat(fullPath)
if err != nil {
if os.IsNotExist(err) {
return false, nil
}
return false, fmt.Errorf("failed to check file existence: %w", err)
}

return true, nil
}

// GetSignedURL returns an empty string for local storage.
// Local storage doesn't support pre-signed URLs.
func (l *LocalStorage) GetSignedURL(ctx context.Context, key string, expiration time.Duration) (string, error) {
// Local storage doesn't support signed URLs
// We could return a file:// URL, but that's typically not useful
return "", nil
}

// resolvePath validates the key and returns the full file system path.
// This prevents path traversal attacks by ensuring the resolved path
// is within the base directory.
func (l *LocalStorage) resolvePath(key string) (string, error) {
// Clean the key to remove any ".." or other path manipulation attempts
cleanKey := filepath.Clean(key)

// Join with base path
fullPath := filepath.Join(l.basePath, cleanKey)

// Ensure the resolved path is still within basePath
absPath, err := filepath.Abs(fullPath)
if err != nil {
return "", fmt.Errorf("failed to resolve path: %w", err)
}

// Use filepath.Rel to verify the path is within base directory
// This is more robust than string prefix checking
relPath, err := filepath.Rel(l.basePath, absPath)
if err != nil {
return "", fmt.Errorf("failed to compute relative path: %w", err)
}

// If the relative path starts with "..", it's outside the base directory
if strings.HasPrefix(relPath, "..") {
return "", errors.New("invalid key: path outside base directory")
}

return absPath, nil
}
