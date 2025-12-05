package main

import (
	"context"
	"sync"

	"opscore/backend/internal/document/domain/entity"
	"opscore/backend/internal/document/domain/repository"
	"opscore/backend/internal/document/domain/value_object"
)

// InMemoryDocumentRepository is an in-memory implementation of DocumentRepository for development.
// This is a temporary implementation until the database persistence layer is complete.
type InMemoryDocumentRepository struct {
	documents map[string]entity.Document
	versions  map[string][]entity.DocumentVersion
	mu        sync.RWMutex
}

// NewInMemoryDocumentRepository creates a new InMemoryDocumentRepository.
func NewInMemoryDocumentRepository() repository.DocumentRepository {
	return &InMemoryDocumentRepository{
		documents: make(map[string]entity.Document),
		versions:  make(map[string][]entity.DocumentVersion),
	}
}

// Save creates a new document.
func (r *InMemoryDocumentRepository) Save(ctx context.Context, document entity.Document) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.documents[document.ID().String()] = document

	// Store versions
	if document.CurrentVersion() != nil {
		r.versions[document.ID().String()] = document.Versions()
	}

	return nil
}

// FindByID retrieves a document by its ID.
func (r *InMemoryDocumentRepository) FindByID(ctx context.Context, id value_object.DocumentID) (entity.Document, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	doc, exists := r.documents[id.String()]
	if !exists {
		return nil, nil
	}

	return doc, nil
}

// FindByRepositoryID retrieves all documents for a given repository.
func (r *InMemoryDocumentRepository) FindByRepositoryID(ctx context.Context, repoID value_object.RepositoryID) ([]entity.Document, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var result []entity.Document
	for _, doc := range r.documents {
		if doc.RepositoryID().Equals(repoID) {
			result = append(result, doc)
		}
	}

	return result, nil
}

// FindPublished retrieves published documents with optional filters.
func (r *InMemoryDocumentRepository) FindPublished(ctx context.Context, filters ...repository.Filter) ([]entity.Document, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var result []entity.Document
	for _, doc := range r.documents {
		if doc.IsPublished() {
			result = append(result, doc)
		}
	}

	return result, nil
}

// Update updates an existing document.
func (r *InMemoryDocumentRepository) Update(ctx context.Context, document entity.Document) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.documents[document.ID().String()] = document

	// Update versions
	if document.CurrentVersion() != nil {
		r.versions[document.ID().String()] = document.Versions()
	}

	return nil
}

// Delete deletes a document by its ID.
func (r *InMemoryDocumentRepository) Delete(ctx context.Context, id value_object.DocumentID) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	delete(r.documents, id.String())
	delete(r.versions, id.String())

	return nil
}

// SaveVersion saves a document version.
func (r *InMemoryDocumentRepository) SaveVersion(ctx context.Context, version entity.DocumentVersion) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	docID := version.DocumentID().String()
	r.versions[docID] = append(r.versions[docID], version)

	return nil
}

// FindVersionsByDocumentID retrieves all versions for a given document.
func (r *InMemoryDocumentRepository) FindVersionsByDocumentID(ctx context.Context, docID value_object.DocumentID) ([]entity.DocumentVersion, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	versions, exists := r.versions[docID.String()]
	if !exists {
		return nil, nil
	}

	return versions, nil
}

// FindVersionByNumber retrieves a specific version by document ID and version number.
func (r *InMemoryDocumentRepository) FindVersionByNumber(ctx context.Context, docID value_object.DocumentID, versionNumber value_object.VersionNumber) (entity.DocumentVersion, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	versions, exists := r.versions[docID.String()]
	if !exists {
		return nil, nil
	}

	for _, v := range versions {
		if v.VersionNumber().Equals(versionNumber) {
			return v, nil
		}
	}

	return nil, nil
}
