package storage

import (
	"context"
	"os"
	"path/filepath"
	"testing"
)

func TestLocalDocumentStorage_SaveContent(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "doc-storage-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	t.Cleanup(func() { _ = os.RemoveAll(tmpDir) })

	store, err := NewLocalDocumentStorage(tmpDir)
	if err != nil {
		t.Fatalf("failed to create storage: %v", err)
	}

	err = store.SaveContent(context.Background(), "doc-1", 2, "hello")
	if err != nil {
		t.Fatalf("failed to save content: %v", err)
	}

	path := filepath.Join(tmpDir, "doc-1", "2", "content.md")
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("failed to read content: %v", err)
	}
	if string(data) != "hello" {
		t.Fatalf("unexpected content: %s", string(data))
	}
}
