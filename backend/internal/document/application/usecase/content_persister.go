package usecase

import (
	"context"

	"opscore/backend/internal/document/domain/entity"
	"opscore/backend/internal/document/infrastructure/storage"
)

// contentPersister writes document content to storage when configured.
type contentPersister struct {
	storage storage.DocumentStorage
}

func newContentPersister(storage storage.DocumentStorage) contentPersister {
	return contentPersister{storage: storage}
}

// persistContent writes the current version content to storage if available.
func (p contentPersister) persistContent(ctx context.Context, doc entity.Document) {
	if p.storage == nil {
		return
	}
	current := doc.CurrentVersion()
	if current == nil {
		return
	}
	_ = p.storage.SaveContent(ctx, doc.ID().String(), current.VersionNumber().Int(), current.Content())
}
