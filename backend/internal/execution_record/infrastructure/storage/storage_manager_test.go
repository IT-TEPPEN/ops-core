package storage

import (
	"bytes"
	"context"
	"io"
	"os"
	"path/filepath"
	"testing"

	"opscore/backend/internal/execution_record/domain/value_object"
)

func TestNewLocalStorageManager(t *testing.T) {
	// Create a temporary directory for testing
	tmpDir := t.TempDir()

	manager, err := NewLocalStorageManager(tmpDir)
	if err != nil {
		t.Fatalf("NewLocalStorageManager() error = %v", err)
	}

	if manager == nil {
		t.Fatal("NewLocalStorageManager() returned nil")
	}

	if manager.Type() != value_object.StorageTypeLocal {
		t.Errorf("Type() = %v, want %v", manager.Type(), value_object.StorageTypeLocal)
	}
}

func TestLocalStorageManager_Store(t *testing.T) {
	tmpDir := t.TempDir()
	manager, err := NewLocalStorageManager(tmpDir)
	if err != nil {
		t.Fatalf("NewLocalStorageManager() error = %v", err)
	}

	ctx := context.Background()
	content := []byte("test file content")
	reader := bytes.NewReader(content)

	path, err := manager.Store(ctx, "test/file.txt", reader)
	if err != nil {
		t.Fatalf("Store() error = %v", err)
	}

	if path != "test/file.txt" {
		t.Errorf("Store() path = %v, want %v", path, "test/file.txt")
	}

	// Verify the file was created
	fullPath := filepath.Join(tmpDir, "test/file.txt")
	if _, err := os.Stat(fullPath); os.IsNotExist(err) {
		t.Error("Store() did not create the file")
	}

	// Verify the content
	fileContent, err := os.ReadFile(fullPath)
	if err != nil {
		t.Fatalf("Failed to read file: %v", err)
	}

	if !bytes.Equal(fileContent, content) {
		t.Errorf("File content = %v, want %v", fileContent, content)
	}
}

func TestLocalStorageManager_Retrieve(t *testing.T) {
	tmpDir := t.TempDir()
	manager, err := NewLocalStorageManager(tmpDir)
	if err != nil {
		t.Fatalf("NewLocalStorageManager() error = %v", err)
	}

	ctx := context.Background()
	content := []byte("test file content for retrieve")

	// Store a file first
	_, err = manager.Store(ctx, "retrieve/file.txt", bytes.NewReader(content))
	if err != nil {
		t.Fatalf("Store() error = %v", err)
	}

	// Retrieve the file
	reader, err := manager.Retrieve(ctx, "retrieve/file.txt")
	if err != nil {
		t.Fatalf("Retrieve() error = %v", err)
	}
	defer reader.Close()

	// Read and verify content
	retrievedContent, err := io.ReadAll(reader)
	if err != nil {
		t.Fatalf("Failed to read content: %v", err)
	}

	if !bytes.Equal(retrievedContent, content) {
		t.Errorf("Retrieved content = %v, want %v", retrievedContent, content)
	}
}

func TestLocalStorageManager_Retrieve_NotFound(t *testing.T) {
	tmpDir := t.TempDir()
	manager, err := NewLocalStorageManager(tmpDir)
	if err != nil {
		t.Fatalf("NewLocalStorageManager() error = %v", err)
	}

	ctx := context.Background()

	// Try to retrieve a non-existent file
	_, err = manager.Retrieve(ctx, "nonexistent/file.txt")
	if err == nil {
		t.Error("Retrieve() should return error for non-existent file")
	}
}

func TestLocalStorageManager_Delete(t *testing.T) {
	tmpDir := t.TempDir()
	manager, err := NewLocalStorageManager(tmpDir)
	if err != nil {
		t.Fatalf("NewLocalStorageManager() error = %v", err)
	}

	ctx := context.Background()
	content := []byte("test file content for delete")

	// Store a file first
	_, err = manager.Store(ctx, "delete/file.txt", bytes.NewReader(content))
	if err != nil {
		t.Fatalf("Store() error = %v", err)
	}

	// Delete the file
	err = manager.Delete(ctx, "delete/file.txt")
	if err != nil {
		t.Fatalf("Delete() error = %v", err)
	}

	// Verify the file was deleted
	fullPath := filepath.Join(tmpDir, "delete/file.txt")
	if _, err := os.Stat(fullPath); !os.IsNotExist(err) {
		t.Error("Delete() did not remove the file")
	}
}

func TestLocalStorageManager_Delete_NotFound(t *testing.T) {
	tmpDir := t.TempDir()
	manager, err := NewLocalStorageManager(tmpDir)
	if err != nil {
		t.Fatalf("NewLocalStorageManager() error = %v", err)
	}

	ctx := context.Background()

	// Try to delete a non-existent file (should not return error)
	err = manager.Delete(ctx, "nonexistent/file.txt")
	if err != nil {
		t.Errorf("Delete() should not return error for non-existent file, got: %v", err)
	}
}

func TestNewS3StorageManager(t *testing.T) {
	tests := []struct {
		name        string
		config      S3Config
		storageType value_object.StorageType
		wantErr     bool
	}{
		{
			name: "valid S3 config",
			config: S3Config{
				Bucket:          "test-bucket",
				Endpoint:        "https://s3.amazonaws.com",
				Region:          "us-east-1",
				AccessKeyID:     "access-key",
				SecretAccessKey: "secret-key",
			},
			storageType: value_object.StorageTypeS3,
			wantErr:     false,
		},
		{
			name: "valid MinIO config",
			config: S3Config{
				Bucket:          "test-bucket",
				Endpoint:        "http://localhost:9000",
				AccessKeyID:     "access-key",
				SecretAccessKey: "secret-key",
				UsePathStyle:    true,
			},
			storageType: value_object.StorageTypeMinio,
			wantErr:     false,
		},
		{
			name: "empty bucket",
			config: S3Config{
				Endpoint:        "https://s3.amazonaws.com",
				AccessKeyID:     "access-key",
				SecretAccessKey: "secret-key",
			},
			storageType: value_object.StorageTypeS3,
			wantErr:     true,
		},
		{
			name: "empty credentials",
			config: S3Config{
				Bucket:   "test-bucket",
				Endpoint: "https://s3.amazonaws.com",
			},
			storageType: value_object.StorageTypeS3,
			wantErr:     true,
		},
		{
			name: "invalid storage type",
			config: S3Config{
				Bucket:          "test-bucket",
				AccessKeyID:     "access-key",
				SecretAccessKey: "secret-key",
			},
			storageType: value_object.StorageTypeLocal,
			wantErr:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewS3StorageManager(tt.config, tt.storageType)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewS3StorageManager() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got == nil {
				t.Error("NewS3StorageManager() returned nil")
			}
			if !tt.wantErr && got.Type() != tt.storageType {
				t.Errorf("Type() = %v, want %v", got.Type(), tt.storageType)
			}
		})
	}
}
