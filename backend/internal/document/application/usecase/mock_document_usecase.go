package usecase

import (
	"context"

	"opscore/backend/internal/document/application/dto"

	"github.com/stretchr/testify/mock"
)

// MockDocumentUseCase is a mock implementation of DocumentUseCase for testing.
type MockDocumentUseCase struct {
	mock.Mock
}

// CreateDocument mocks the CreateDocument method.
func (m *MockDocumentUseCase) CreateDocument(ctx context.Context, req *dto.CreateDocumentRequest) (*dto.DocumentResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.DocumentResponse), args.Error(1)
}

// UpdateDocument mocks the UpdateDocument method.
func (m *MockDocumentUseCase) UpdateDocument(ctx context.Context, documentID string, req *dto.UpdateDocumentRequest) (*dto.DocumentResponse, error) {
	args := m.Called(ctx, documentID, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.DocumentResponse), args.Error(1)
}

// GetDocument mocks the GetDocument method.
func (m *MockDocumentUseCase) GetDocument(ctx context.Context, documentID string) (*dto.DocumentResponse, error) {
	args := m.Called(ctx, documentID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.DocumentResponse), args.Error(1)
}

// GetDocumentVersion mocks the GetDocumentVersion method.
func (m *MockDocumentUseCase) GetDocumentVersion(ctx context.Context, documentID string, versionNumber int) (*dto.DocumentVersionResponse, error) {
	args := m.Called(ctx, documentID, versionNumber)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.DocumentVersionResponse), args.Error(1)
}

// ListDocuments mocks the ListDocuments method.
func (m *MockDocumentUseCase) ListDocuments(ctx context.Context) ([]dto.DocumentListItemResponse, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]dto.DocumentListItemResponse), args.Error(1)
}

// ListDocumentsByRepository mocks the ListDocumentsByRepository method.
func (m *MockDocumentUseCase) ListDocumentsByRepository(ctx context.Context, repositoryID string) ([]dto.DocumentListItemResponse, error) {
	args := m.Called(ctx, repositoryID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]dto.DocumentListItemResponse), args.Error(1)
}

// GetDocumentVersions mocks the GetDocumentVersions method.
func (m *MockDocumentUseCase) GetDocumentVersions(ctx context.Context, documentID string) (*dto.VersionHistoryResponse, error) {
	args := m.Called(ctx, documentID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.VersionHistoryResponse), args.Error(1)
}

// PublishDocumentVersion mocks the PublishDocumentVersion method.
func (m *MockDocumentUseCase) PublishDocumentVersion(ctx context.Context, documentID string, versionNumber int) (*dto.DocumentResponse, error) {
	args := m.Called(ctx, documentID, versionNumber)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.DocumentResponse), args.Error(1)
}

// RollbackDocumentVersion mocks the RollbackDocumentVersion method.
func (m *MockDocumentUseCase) RollbackDocumentVersion(ctx context.Context, documentID string, versionNumber int) (*dto.DocumentResponse, error) {
	args := m.Called(ctx, documentID, versionNumber)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.DocumentResponse), args.Error(1)
}

// UpdateDocumentMetadata mocks the UpdateDocumentMetadata method.
func (m *MockDocumentUseCase) UpdateDocumentMetadata(ctx context.Context, documentID string, req *dto.UpdateDocumentMetadataRequest) (*dto.DocumentResponse, error) {
	args := m.Called(ctx, documentID, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.DocumentResponse), args.Error(1)
}
