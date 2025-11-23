package repository

import (
	"context"

	"opscore/backend/internal/document/domain/entity"
	"opscore/backend/internal/document/domain/value_object"
)

// Filter represents a filter for querying documents.
type Filter interface {
	Apply(ctx context.Context) context.Context
}

// DocumentRepository defines the interface for document persistence.
type DocumentRepository interface {
	// Save creates a new document.
	Save(ctx context.Context, document entity.Document) error

	// FindByID retrieves a document by its ID.
	FindByID(ctx context.Context, id value_object.DocumentID) (entity.Document, error)

	// FindByRepositoryID retrieves all documents for a given repository.
	FindByRepositoryID(ctx context.Context, repoID value_object.RepositoryID) ([]entity.Document, error)

	// FindPublished retrieves published documents with optional filters.
	FindPublished(ctx context.Context, filters ...Filter) ([]entity.Document, error)

	// Update updates an existing document.
	Update(ctx context.Context, document entity.Document) error

	// Delete deletes a document by its ID.
	Delete(ctx context.Context, id value_object.DocumentID) error

	// SaveVersion saves a document version.
	SaveVersion(ctx context.Context, version entity.DocumentVersion) error

	// FindVersionsByDocumentID retrieves all versions for a given document.
	FindVersionsByDocumentID(ctx context.Context, docID value_object.DocumentID) ([]entity.DocumentVersion, error)

	// FindVersionByNumber retrieves a specific version by document ID and version number.
	FindVersionByNumber(ctx context.Context, docID value_object.DocumentID, versionNumber value_object.VersionNumber) (entity.DocumentVersion, error)
}
