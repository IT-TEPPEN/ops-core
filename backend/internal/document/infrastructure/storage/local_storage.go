package storage

import (
	context "context"
	"fmt"
	"os"
	"path/filepath"
)

// DocumentStorage persists document contents alongside DB metadata.
type DocumentStorage interface {
	SaveContent(ctx context.Context, documentID string, version int, content string) error
}

// localDocumentStorage writes document content to the local filesystem.
type localDocumentStorage struct {
	basePath string
}

// NewLocalDocumentStorage creates a filesystem-backed storage rooted at basePath.
func NewLocalDocumentStorage(basePath string) (DocumentStorage, error) {
	if basePath == "" {
		return nil, fmt.Errorf("basePath cannot be empty")
	}
	if err := os.MkdirAll(basePath, 0o755); err != nil {
		return nil, fmt.Errorf("failed to create document storage directory: %w", err)
	}
	return &localDocumentStorage{basePath: basePath}, nil
}

// SaveContent writes content to {basePath}/{documentID}/{version}/content.md.
func (s *localDocumentStorage) SaveContent(ctx context.Context, documentID string, version int, content string) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	dir := filepath.Join(s.basePath, documentID, fmt.Sprintf("%d", version))
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return fmt.Errorf("failed to create storage directory: %w", err)
	}

	path := filepath.Join(dir, "content.md")
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		return fmt.Errorf("failed to write document content: %w", err)
	}
	return nil
}
