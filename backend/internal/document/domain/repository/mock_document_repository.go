package repository

import (
	"context"

	"opscore/backend/internal/document/domain/entity"
	"opscore/backend/internal/document/domain/value_object"

	"github.com/stretchr/testify/mock"
)

// MockDocumentRepository is a mock implementation of DocumentRepository for testing.
type MockDocumentRepository struct {
	mock.Mock
}

// Save mocks the Save method.
func (m *MockDocumentRepository) Save(ctx context.Context, document entity.Document) error {
	args := m.Called(ctx, document)
	return args.Error(0)
}

// FindByID mocks the FindByID method.
func (m *MockDocumentRepository) FindByID(ctx context.Context, id value_object.DocumentID) (entity.Document, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(entity.Document), args.Error(1)
}

// FindByRepositoryID mocks the FindByRepositoryID method.
func (m *MockDocumentRepository) FindByRepositoryID(ctx context.Context, repoID value_object.RepositoryID) ([]entity.Document, error) {
	args := m.Called(ctx, repoID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]entity.Document), args.Error(1)
}

// FindPublished mocks the FindPublished method.
func (m *MockDocumentRepository) FindPublished(ctx context.Context, filters ...Filter) ([]entity.Document, error) {
	args := m.Called(ctx, filters)
	if args.Get(0) == nil {
		return []entity.Document{}, args.Error(1)
	}
	return args.Get(0).([]entity.Document), args.Error(1)
}

// Update mocks the Update method.
func (m *MockDocumentRepository) Update(ctx context.Context, document entity.Document) error {
	args := m.Called(ctx, document)
	return args.Error(0)
}

// Delete mocks the Delete method.
func (m *MockDocumentRepository) Delete(ctx context.Context, id value_object.DocumentID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

// SaveVersion mocks the SaveVersion method.
func (m *MockDocumentRepository) SaveVersion(ctx context.Context, version entity.DocumentVersion) error {
	args := m.Called(ctx, version)
	return args.Error(0)
}

// FindVersionsByDocumentID mocks the FindVersionsByDocumentID method.
func (m *MockDocumentRepository) FindVersionsByDocumentID(ctx context.Context, docID value_object.DocumentID) ([]entity.DocumentVersion, error) {
	args := m.Called(ctx, docID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]entity.DocumentVersion), args.Error(1)
}

// FindVersionByNumber mocks the FindVersionByNumber method.
func (m *MockDocumentRepository) FindVersionByNumber(ctx context.Context, docID value_object.DocumentID, versionNumber value_object.VersionNumber) (entity.DocumentVersion, error) {
	args := m.Called(ctx, docID, versionNumber)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(entity.DocumentVersion), args.Error(1)
}
