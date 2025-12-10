package storage

import (
	"bytes"
	"context"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestNewLocalStorage(t *testing.T) {
	t.Run("creates storage with valid path", func(t *testing.T) {
		tmpDir := t.TempDir()
		storage, err := NewLocalStorage(tmpDir)

		if err != nil {
			t.Fatalf("NewLocalStorage() error = %v, want nil", err)
		}
		if storage == nil {
			t.Fatal("NewLocalStorage() returned nil storage")
		}
	})

	t.Run("creates base directory if not exists", func(t *testing.T) {
		tmpDir := t.TempDir()
		newDir := filepath.Join(tmpDir, "new", "nested", "dir")

		storage, err := NewLocalStorage(newDir)

		if err != nil {
			t.Fatalf("NewLocalStorage() error = %v, want nil", err)
		}
		if storage == nil {
			t.Fatal("NewLocalStorage() returned nil storage")
		}

		// Verify directory was created
		if _, err := os.Stat(newDir); os.IsNotExist(err) {
			t.Error("NewLocalStorage() did not create base directory")
		}
	})

	t.Run("returns error for empty base path", func(t *testing.T) {
		_, err := NewLocalStorage("")

		if err == nil {
			t.Error("NewLocalStorage() with empty path should return error")
		}
	})
}

func TestLocalStorage_Save(t *testing.T) {
	t.Run("saves file successfully", func(t *testing.T) {
		tmpDir := t.TempDir()
		storage, _ := NewLocalStorage(tmpDir)
		ctx := context.Background()

		content := []byte("test content")
		reader := bytes.NewReader(content)

		err := storage.Save(ctx, "test/file.txt", reader, "text/plain")

		if err != nil {
			t.Fatalf("Save() error = %v, want nil", err)
		}

		// Verify file exists and has correct content
		fullPath := filepath.Join(tmpDir, "test/file.txt")
		savedContent, err := os.ReadFile(fullPath)
		if err != nil {
			t.Fatalf("Failed to read saved file: %v", err)
		}

		if !bytes.Equal(savedContent, content) {
			t.Errorf("Saved content = %v, want %v", savedContent, content)
		}
	})

	t.Run("creates nested directories", func(t *testing.T) {
		tmpDir := t.TempDir()
		storage, _ := NewLocalStorage(tmpDir)
		ctx := context.Background()

		content := []byte("nested content")
		reader := bytes.NewReader(content)

		err := storage.Save(ctx, "a/b/c/file.txt", reader, "text/plain")

		if err != nil {
			t.Fatalf("Save() error = %v, want nil", err)
		}

		fullPath := filepath.Join(tmpDir, "a/b/c/file.txt")
		if _, err := os.Stat(fullPath); os.IsNotExist(err) {
			t.Error("Save() did not create nested directories")
		}
	})

	t.Run("returns error for empty key", func(t *testing.T) {
		tmpDir := t.TempDir()
		storage, _ := NewLocalStorage(tmpDir)
		ctx := context.Background()

		reader := bytes.NewReader([]byte("content"))
		err := storage.Save(ctx, "", reader, "text/plain")

		if err == nil {
			t.Error("Save() with empty key should return error")
		}
	})

	t.Run("returns error for nil reader", func(t *testing.T) {
		tmpDir := t.TempDir()
		storage, _ := NewLocalStorage(tmpDir)
		ctx := context.Background()

		err := storage.Save(ctx, "test.txt", nil, "text/plain")

		if err == nil {
			t.Error("Save() with nil reader should return error")
		}
	})

	t.Run("prevents path traversal attack", func(t *testing.T) {
		tmpDir := t.TempDir()
		storage, _ := NewLocalStorage(tmpDir)
		ctx := context.Background()

		content := bytes.NewReader([]byte("attack"))

		// Try various path traversal attempts
		maliciousKeys := []string{
			"../../../etc/passwd",
			"..\\..\\..\\windows\\system32",
			"test/../../outside.txt",
		}

		for _, key := range maliciousKeys {
			err := storage.Save(ctx, key, content, "text/plain")
			if err == nil {
				t.Errorf("Save() with malicious key %q should return error", key)
			}
			content.Reset([]byte("attack"))
		}
	})
}

func TestLocalStorage_Get(t *testing.T) {
	t.Run("retrieves file successfully", func(t *testing.T) {
		tmpDir := t.TempDir()
		storage, _ := NewLocalStorage(tmpDir)
		ctx := context.Background()

		content := []byte("retrieve test content")
		_ = storage.Save(ctx, "retrieve/file.txt", bytes.NewReader(content), "text/plain")

		reader, err := storage.Get(ctx, "retrieve/file.txt")

		if err != nil {
			t.Fatalf("Get() error = %v, want nil", err)
		}
		defer reader.Close()

		retrievedContent, _ := io.ReadAll(reader)
		if !bytes.Equal(retrievedContent, content) {
			t.Errorf("Get() content = %v, want %v", retrievedContent, content)
		}
	})

	t.Run("returns error for non-existent file", func(t *testing.T) {
		tmpDir := t.TempDir()
		storage, _ := NewLocalStorage(tmpDir)
		ctx := context.Background()

		_, err := storage.Get(ctx, "nonexistent.txt")

		if err == nil {
			t.Error("Get() for non-existent file should return error")
		}
		if !strings.Contains(err.Error(), "not found") {
			t.Errorf("Get() error should mention file not found, got: %v", err)
		}
	})

	t.Run("returns error for empty key", func(t *testing.T) {
		tmpDir := t.TempDir()
		storage, _ := NewLocalStorage(tmpDir)
		ctx := context.Background()

		_, err := storage.Get(ctx, "")

		if err == nil {
			t.Error("Get() with empty key should return error")
		}
	})

	t.Run("prevents path traversal attack", func(t *testing.T) {
		tmpDir := t.TempDir()
		storage, _ := NewLocalStorage(tmpDir)
		ctx := context.Background()

		_, err := storage.Get(ctx, "../../../etc/passwd")

		if err == nil {
			t.Error("Get() with malicious key should return error")
		}
	})
}

func TestLocalStorage_Delete(t *testing.T) {
	t.Run("deletes file successfully", func(t *testing.T) {
		tmpDir := t.TempDir()
		storage, _ := NewLocalStorage(tmpDir)
		ctx := context.Background()

		content := []byte("delete test")
		_ = storage.Save(ctx, "delete/file.txt", bytes.NewReader(content), "text/plain")

		err := storage.Delete(ctx, "delete/file.txt")

		if err != nil {
			t.Fatalf("Delete() error = %v, want nil", err)
		}

		// Verify file is deleted
		fullPath := filepath.Join(tmpDir, "delete/file.txt")
		if _, err := os.Stat(fullPath); !os.IsNotExist(err) {
			t.Error("Delete() did not remove the file")
		}
	})

	t.Run("is idempotent for non-existent file", func(t *testing.T) {
		tmpDir := t.TempDir()
		storage, _ := NewLocalStorage(tmpDir)
		ctx := context.Background()

		err := storage.Delete(ctx, "nonexistent.txt")

		if err != nil {
			t.Errorf("Delete() for non-existent file should not return error, got: %v", err)
		}
	})

	t.Run("returns error for empty key", func(t *testing.T) {
		tmpDir := t.TempDir()
		storage, _ := NewLocalStorage(tmpDir)
		ctx := context.Background()

		err := storage.Delete(ctx, "")

		if err == nil {
			t.Error("Delete() with empty key should return error")
		}
	})

	t.Run("prevents path traversal attack", func(t *testing.T) {
		tmpDir := t.TempDir()
		storage, _ := NewLocalStorage(tmpDir)
		ctx := context.Background()

		err := storage.Delete(ctx, "../../../etc/passwd")

		if err == nil {
			t.Error("Delete() with malicious key should return error")
		}
	})
}

func TestLocalStorage_Exists(t *testing.T) {
	t.Run("returns true for existing file", func(t *testing.T) {
		tmpDir := t.TempDir()
		storage, _ := NewLocalStorage(tmpDir)
		ctx := context.Background()

		content := []byte("exists test")
		_ = storage.Save(ctx, "exists/file.txt", bytes.NewReader(content), "text/plain")

		exists, err := storage.Exists(ctx, "exists/file.txt")

		if err != nil {
			t.Fatalf("Exists() error = %v, want nil", err)
		}
		if !exists {
			t.Error("Exists() = false, want true")
		}
	})

	t.Run("returns false for non-existent file", func(t *testing.T) {
		tmpDir := t.TempDir()
		storage, _ := NewLocalStorage(tmpDir)
		ctx := context.Background()

		exists, err := storage.Exists(ctx, "nonexistent.txt")

		if err != nil {
			t.Fatalf("Exists() error = %v, want nil", err)
		}
		if exists {
			t.Error("Exists() = true, want false")
		}
	})

	t.Run("returns error for empty key", func(t *testing.T) {
		tmpDir := t.TempDir()
		storage, _ := NewLocalStorage(tmpDir)
		ctx := context.Background()

		_, err := storage.Exists(ctx, "")

		if err == nil {
			t.Error("Exists() with empty key should return error")
		}
	})

	t.Run("prevents path traversal attack", func(t *testing.T) {
		tmpDir := t.TempDir()
		storage, _ := NewLocalStorage(tmpDir)
		ctx := context.Background()

		_, err := storage.Exists(ctx, "../../../etc/passwd")

		if err == nil {
			t.Error("Exists() with malicious key should return error")
		}
	})
}

func TestLocalStorage_GetSignedURL(t *testing.T) {
	t.Run("returns empty string", func(t *testing.T) {
		tmpDir := t.TempDir()
		storage, _ := NewLocalStorage(tmpDir)
		ctx := context.Background()

		url, err := storage.GetSignedURL(ctx, "test.txt", 1*time.Hour)

		if err != nil {
			t.Fatalf("GetSignedURL() error = %v, want nil", err)
		}
		if url != "" {
			t.Errorf("GetSignedURL() = %q, want empty string", url)
		}
	})
}

func TestLocalStorage_ResolvePath(t *testing.T) {
	tmpDir := t.TempDir()
	storage, _ := NewLocalStorage(tmpDir)

	t.Run("resolves valid path", func(t *testing.T) {
		path, err := storage.resolvePath("test/file.txt")

		if err != nil {
			t.Fatalf("resolvePath() error = %v, want nil", err)
		}

		expected := filepath.Join(tmpDir, "test/file.txt")
		if path != expected {
			t.Errorf("resolvePath() = %q, want %q", path, expected)
		}
	})

	t.Run("rejects path with double dots", func(t *testing.T) {
		_, err := storage.resolvePath("../outside.txt")

		if err == nil {
			t.Error("resolvePath() with '..' should return error")
		}
	})

	t.Run("rejects path escaping base directory", func(t *testing.T) {
		_, err := storage.resolvePath("test/../../outside.txt")

		if err == nil {
			t.Error("resolvePath() escaping base should return error")
		}
	})
}
