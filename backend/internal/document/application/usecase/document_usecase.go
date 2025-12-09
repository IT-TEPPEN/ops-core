package usecase

import (
	"context"
	"fmt"

	apperror "opscore/backend/internal/document/application/error"
	"opscore/backend/internal/document/application/dto"
	"opscore/backend/internal/document/domain/entity"
	"opscore/backend/internal/document/domain/repository"
	"opscore/backend/internal/document/domain/value_object"
)

// DocumentUseCase defines the interface for document related use cases.
type DocumentUseCase interface {
	// CreateDocument creates a new document with an initial version.
	CreateDocument(ctx context.Context, req *dto.CreateDocumentRequest) (*dto.DocumentResponse, error)

	// UpdateDocument updates an existing document by creating a new version.
	UpdateDocument(ctx context.Context, documentID string, req *dto.UpdateDocumentRequest) (*dto.DocumentResponse, error)

	// GetDocument retrieves a document by its ID.
	GetDocument(ctx context.Context, documentID string) (*dto.DocumentResponse, error)

	// GetDocumentVersion retrieves a specific version of a document.
	GetDocumentVersion(ctx context.Context, documentID string, versionNumber int) (*dto.DocumentVersionResponse, error)

	// ListDocuments retrieves all documents.
	ListDocuments(ctx context.Context) ([]dto.DocumentListItemResponse, error)

	// ListDocumentsByRepository retrieves all documents for a given repository.
	ListDocumentsByRepository(ctx context.Context, repositoryID string) ([]dto.DocumentListItemResponse, error)

	// GetDocumentVersions retrieves all versions for a document.
	GetDocumentVersions(ctx context.Context, documentID string) (*dto.VersionHistoryResponse, error)

	// PublishDocumentVersion publishes a specific version.
	PublishDocumentVersion(ctx context.Context, documentID string, versionNumber int) (*dto.DocumentResponse, error)

	// RollbackDocumentVersion rolls back to a previous version.
	RollbackDocumentVersion(ctx context.Context, documentID string, versionNumber int) (*dto.DocumentResponse, error)

	// UpdateDocumentMetadata updates the document metadata (owner, access scope, etc.).
	UpdateDocumentMetadata(ctx context.Context, documentID string, req *dto.UpdateDocumentMetadataRequest) (*dto.DocumentResponse, error)
}

// documentUseCase implements the DocumentUseCase interface.
type documentUseCase struct {
	repo repository.DocumentRepository
}

// NewDocumentUseCase creates a new instance of documentUseCase.
func NewDocumentUseCase(repo repository.DocumentRepository) DocumentUseCase {
	return &documentUseCase{
		repo: repo,
	}
}

// CreateDocument creates a new document with an initial version.
func (uc *documentUseCase) CreateDocument(ctx context.Context, req *dto.CreateDocumentRequest) (*dto.DocumentResponse, error) {
	// Validate repository ID
	repositoryID, err := value_object.NewRepositoryID(req.RepositoryID)
	if err != nil {
		return nil, apperror.NewValidationFailedError([]apperror.FieldError{
			{Field: "repository_id", Message: err.Error()},
		})
	}

	// Validate access scope
	accessScope, err := value_object.NewAccessScope(req.AccessScope)
	if err != nil {
		return nil, apperror.NewValidationFailedError([]apperror.FieldError{
			{Field: "access_scope", Message: err.Error()},
		})
	}

	// Validate document type
	docType, err := value_object.NewDocumentType(req.DocType)
	if err != nil {
		return nil, apperror.NewValidationFailedError([]apperror.FieldError{
			{Field: "doc_type", Message: err.Error()},
		})
	}

	// Create file path and commit hash
	filePath, err := value_object.NewFilePath(req.FilePath)
	if err != nil {
		return nil, apperror.NewValidationFailedError([]apperror.FieldError{
			{Field: "file_path", Message: err.Error()},
		})
	}

	commitHash, err := value_object.NewCommitHash(req.CommitHash)
	if err != nil {
		return nil, apperror.NewValidationFailedError([]apperror.FieldError{
			{Field: "commit_hash", Message: err.Error()},
		})
	}

	// Create document source
	source, err := value_object.NewDocumentSource(filePath, commitHash)
	if err != nil {
		return nil, apperror.NewValidationFailedError([]apperror.FieldError{
			{Field: "source", Message: err.Error()},
		})
	}

	// Convert tags
	tags, err := dto.ToTagSlice(req.Tags)
	if err != nil {
		return nil, apperror.NewValidationFailedError([]apperror.FieldError{
			{Field: "tags", Message: err.Error()},
		})
	}

	// Convert variables
	variables, err := dto.ToVariableDefinitionSlice(req.Variables)
	if err != nil {
		return nil, apperror.NewValidationFailedError([]apperror.FieldError{
			{Field: "variables", Message: err.Error()},
		})
	}

	// Validate title
	if req.Title == "" {
		return nil, apperror.NewValidationFailedError([]apperror.FieldError{
			{Field: "title", Message: "title cannot be empty"},
		})
	}

	// Validate content
	if req.Content == "" {
		return nil, apperror.NewValidationFailedError([]apperror.FieldError{
			{Field: "content", Message: "content cannot be empty"},
		})
	}

	// Validate owner
	if req.Owner == "" {
		return nil, apperror.NewValidationFailedError([]apperror.FieldError{
			{Field: "owner", Message: "owner cannot be empty"},
		})
	}

	// Generate new document ID
	documentID := value_object.GenerateDocumentID()

	// Create the document entity
	doc, err := entity.NewDocument(documentID, repositoryID, req.Owner, accessScope)
	if err != nil {
		return nil, fmt.Errorf("failed to create document: %w", err)
	}

	// Set auto update if specified
	if req.IsAutoUpdate {
		doc.EnableAutoUpdate()
	}

	// Publish the initial version
	err = doc.Publish(source, req.Title, docType, tags, variables, req.Content)
	if err != nil {
		return nil, fmt.Errorf("failed to publish initial version: %w", err)
	}

	// Save the document
	err = uc.repo.Save(ctx, doc)
	if err != nil {
		return nil, fmt.Errorf("failed to save document: %w", err)
	}

	// Return the response
	response := dto.ToDocumentResponse(doc)
	return &response, nil
}

// UpdateDocument updates an existing document by creating a new version.
func (uc *documentUseCase) UpdateDocument(ctx context.Context, documentID string, req *dto.UpdateDocumentRequest) (*dto.DocumentResponse, error) {
	// Validate document ID
	docID, err := value_object.NewDocumentID(documentID)
	if err != nil {
		return nil, apperror.NewValidationFailedError([]apperror.FieldError{
			{Field: "document_id", Message: err.Error()},
		})
	}

	// Find the document
	doc, err := uc.repo.FindByID(ctx, docID)
	if err != nil {
		return nil, fmt.Errorf("failed to find document: %w", err)
	}
	if doc == nil {
		return nil, apperror.NewNotFoundError("Document", documentID, nil)
	}

	// Validate document type
	docType, err := value_object.NewDocumentType(req.DocType)
	if err != nil {
		return nil, apperror.NewValidationFailedError([]apperror.FieldError{
			{Field: "doc_type", Message: err.Error()},
		})
	}

	// Create file path and commit hash
	filePath, err := value_object.NewFilePath(req.FilePath)
	if err != nil {
		return nil, apperror.NewValidationFailedError([]apperror.FieldError{
			{Field: "file_path", Message: err.Error()},
		})
	}

	commitHash, err := value_object.NewCommitHash(req.CommitHash)
	if err != nil {
		return nil, apperror.NewValidationFailedError([]apperror.FieldError{
			{Field: "commit_hash", Message: err.Error()},
		})
	}

	// Create document source
	source, err := value_object.NewDocumentSource(filePath, commitHash)
	if err != nil {
		return nil, apperror.NewValidationFailedError([]apperror.FieldError{
			{Field: "source", Message: err.Error()},
		})
	}

	// Convert tags
	tags, err := dto.ToTagSlice(req.Tags)
	if err != nil {
		return nil, apperror.NewValidationFailedError([]apperror.FieldError{
			{Field: "tags", Message: err.Error()},
		})
	}

	// Convert variables
	variables, err := dto.ToVariableDefinitionSlice(req.Variables)
	if err != nil {
		return nil, apperror.NewValidationFailedError([]apperror.FieldError{
			{Field: "variables", Message: err.Error()},
		})
	}

	// Validate title
	if req.Title == "" {
		return nil, apperror.NewValidationFailedError([]apperror.FieldError{
			{Field: "title", Message: "title cannot be empty"},
		})
	}

	// Validate content
	if req.Content == "" {
		return nil, apperror.NewValidationFailedError([]apperror.FieldError{
			{Field: "content", Message: "content cannot be empty"},
		})
	}

	// Publish the new version
	err = doc.Publish(source, req.Title, docType, tags, variables, req.Content)
	if err != nil {
		return nil, fmt.Errorf("failed to publish new version: %w", err)
	}

	// Update the document
	err = uc.repo.Update(ctx, doc)
	if err != nil {
		return nil, fmt.Errorf("failed to update document: %w", err)
	}

	// Return the response
	response := dto.ToDocumentResponse(doc)
	return &response, nil
}

// GetDocument retrieves a document by its ID.
func (uc *documentUseCase) GetDocument(ctx context.Context, documentID string) (*dto.DocumentResponse, error) {
	// Validate document ID
	docID, err := value_object.NewDocumentID(documentID)
	if err != nil {
		return nil, apperror.NewValidationFailedError([]apperror.FieldError{
			{Field: "document_id", Message: err.Error()},
		})
	}

	// Find the document
	doc, err := uc.repo.FindByID(ctx, docID)
	if err != nil {
		return nil, fmt.Errorf("failed to find document: %w", err)
	}
	if doc == nil {
		return nil, apperror.NewNotFoundError("Document", documentID, nil)
	}

	// Return the response
	response := dto.ToDocumentResponse(doc)
	return &response, nil
}

// GetDocumentVersion retrieves a specific version of a document.
func (uc *documentUseCase) GetDocumentVersion(ctx context.Context, documentID string, versionNumber int) (*dto.DocumentVersionResponse, error) {
	// Validate document ID
	docID, err := value_object.NewDocumentID(documentID)
	if err != nil {
		return nil, apperror.NewValidationFailedError([]apperror.FieldError{
			{Field: "document_id", Message: err.Error()},
		})
	}

	// Validate version number
	verNum, err := value_object.NewVersionNumber(versionNumber)
	if err != nil {
		return nil, apperror.NewValidationFailedError([]apperror.FieldError{
			{Field: "version_number", Message: err.Error()},
		})
	}

	// Find the version
	version, err := uc.repo.FindVersionByNumber(ctx, docID, verNum)
	if err != nil {
		return nil, fmt.Errorf("failed to find document version: %w", err)
	}
	if version == nil {
		return nil, apperror.NewNotFoundError("DocumentVersion", fmt.Sprintf("%s@v%d", documentID, versionNumber), nil)
	}

	// Return the response
	response := dto.ToDocumentVersionResponse(version)
	return &response, nil
}

// ListDocuments retrieves all documents.
func (uc *documentUseCase) ListDocuments(ctx context.Context) ([]dto.DocumentListItemResponse, error) {
	// Find all published documents
	docs, err := uc.repo.FindPublished(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to find documents: %w", err)
	}

	// Return the response
	return dto.ToDocumentListResponse(docs), nil
}

// ListDocumentsByRepository retrieves all documents for a given repository.
func (uc *documentUseCase) ListDocumentsByRepository(ctx context.Context, repositoryID string) ([]dto.DocumentListItemResponse, error) {
	// Validate repository ID
	repoID, err := value_object.NewRepositoryID(repositoryID)
	if err != nil {
		return nil, apperror.NewValidationFailedError([]apperror.FieldError{
			{Field: "repository_id", Message: err.Error()},
		})
	}

	// Find all documents for the repository
	docs, err := uc.repo.FindByRepositoryID(ctx, repoID)
	if err != nil {
		return nil, fmt.Errorf("failed to find documents: %w", err)
	}

	// Return the response
	return dto.ToDocumentListResponse(docs), nil
}

// GetDocumentVersions retrieves all versions for a document.
func (uc *documentUseCase) GetDocumentVersions(ctx context.Context, documentID string) (*dto.VersionHistoryResponse, error) {
	// Validate document ID
	docID, err := value_object.NewDocumentID(documentID)
	if err != nil {
		return nil, apperror.NewValidationFailedError([]apperror.FieldError{
			{Field: "document_id", Message: err.Error()},
		})
	}

	// Check if document exists
	doc, err := uc.repo.FindByID(ctx, docID)
	if err != nil {
		return nil, fmt.Errorf("failed to find document: %w", err)
	}
	if doc == nil {
		return nil, apperror.NewNotFoundError("Document", documentID, nil)
	}

	// Find all versions
	versions, err := uc.repo.FindVersionsByDocumentID(ctx, docID)
	if err != nil {
		return nil, fmt.Errorf("failed to find document versions: %w", err)
	}

	// Return the response
	response := dto.ToVersionHistoryResponse(documentID, versions)
	return &response, nil
}

// PublishDocumentVersion publishes a specific version.
func (uc *documentUseCase) PublishDocumentVersion(ctx context.Context, documentID string, versionNumber int) (*dto.DocumentResponse, error) {
	// Validate document ID
	docID, err := value_object.NewDocumentID(documentID)
	if err != nil {
		return nil, apperror.NewValidationFailedError([]apperror.FieldError{
			{Field: "document_id", Message: err.Error()},
		})
	}

	// Find the document
	doc, err := uc.repo.FindByID(ctx, docID)
	if err != nil {
		return nil, fmt.Errorf("failed to find document: %w", err)
	}
	if doc == nil {
		return nil, apperror.NewNotFoundError("Document", documentID, nil)
	}

	// Validate version number
	verNum, err := value_object.NewVersionNumber(versionNumber)
	if err != nil {
		return nil, apperror.NewValidationFailedError([]apperror.FieldError{
			{Field: "version_number", Message: err.Error()},
		})
	}

	// Find the version to check if it exists
	version, err := uc.repo.FindVersionByNumber(ctx, docID, verNum)
	if err != nil {
		return nil, fmt.Errorf("failed to find document version: %w", err)
	}
	if version == nil {
		return nil, apperror.NewNotFoundError("DocumentVersion", fmt.Sprintf("%s@v%d", documentID, versionNumber), nil)
	}

	// If the document is already published with this version as current, just return it
	if doc.IsPublished() && doc.CurrentVersion() != nil && doc.CurrentVersion().VersionNumber().Equals(verNum) {
		response := dto.ToDocumentResponse(doc)
		return &response, nil
	}

	// Rollback to the specified version (which effectively publishes it)
	if doc.IsPublished() {
		err = doc.RollbackToVersion(verNum)
		if err != nil {
			return nil, fmt.Errorf("failed to publish version: %w", err)
		}
	} else {
		// If the document is not published, publish it with the specified version
		err = doc.PublishWithVersion(verNum)
		if err != nil {
			return nil, fmt.Errorf("failed to publish version: %w", err)
		}
	}

	// Update the document
	err = uc.repo.Update(ctx, doc)
	if err != nil {
		return nil, fmt.Errorf("failed to update document: %w", err)
	}

	// Return the response
	response := dto.ToDocumentResponse(doc)
	return &response, nil
}

// RollbackDocumentVersion rolls back to a previous version.
func (uc *documentUseCase) RollbackDocumentVersion(ctx context.Context, documentID string, versionNumber int) (*dto.DocumentResponse, error) {
	// Validate document ID
	docID, err := value_object.NewDocumentID(documentID)
	if err != nil {
		return nil, apperror.NewValidationFailedError([]apperror.FieldError{
			{Field: "document_id", Message: err.Error()},
		})
	}

	// Validate version number
	verNum, err := value_object.NewVersionNumber(versionNumber)
	if err != nil {
		return nil, apperror.NewValidationFailedError([]apperror.FieldError{
			{Field: "version_number", Message: err.Error()},
		})
	}

	// Find the document
	doc, err := uc.repo.FindByID(ctx, docID)
	if err != nil {
		return nil, fmt.Errorf("failed to find document: %w", err)
	}
	if doc == nil {
		return nil, apperror.NewNotFoundError("Document", documentID, nil)
	}

	// Rollback to the specified version
	err = doc.RollbackToVersion(verNum)
	if err != nil {
		return nil, fmt.Errorf("failed to rollback to version: %w", err)
	}

	// Update the document
	err = uc.repo.Update(ctx, doc)
	if err != nil {
		return nil, fmt.Errorf("failed to update document: %w", err)
	}

	// Return the response
	response := dto.ToDocumentResponse(doc)
	return &response, nil
}

// UpdateDocumentMetadata updates the document metadata.
func (uc *documentUseCase) UpdateDocumentMetadata(ctx context.Context, documentID string, req *dto.UpdateDocumentMetadataRequest) (*dto.DocumentResponse, error) {
	// Validate document ID
	docID, err := value_object.NewDocumentID(documentID)
	if err != nil {
		return nil, apperror.NewValidationFailedError([]apperror.FieldError{
			{Field: "document_id", Message: err.Error()},
		})
	}

	// Find the document
	doc, err := uc.repo.FindByID(ctx, docID)
	if err != nil {
		return nil, fmt.Errorf("failed to find document: %w", err)
	}
	if doc == nil {
		return nil, apperror.NewNotFoundError("Document", documentID, nil)
	}

	// Update access scope if provided
	if req.AccessScope != nil {
		accessScope, err := value_object.NewAccessScope(*req.AccessScope)
		if err != nil {
			return nil, apperror.NewValidationFailedError([]apperror.FieldError{
				{Field: "access_scope", Message: err.Error()},
			})
		}
		err = doc.UpdateAccessScope(accessScope)
		if err != nil {
			return nil, fmt.Errorf("failed to update access scope: %w", err)
		}
	}

	// Update auto update setting if provided
	if req.IsAutoUpdate != nil {
		if *req.IsAutoUpdate {
			doc.EnableAutoUpdate()
		} else {
			doc.DisableAutoUpdate()
		}
	}

	// Update the document
	err = uc.repo.Update(ctx, doc)
	if err != nil {
		return nil, fmt.Errorf("failed to update document: %w", err)
	}

	// Return the response
	response := dto.ToDocumentResponse(doc)
	return &response, nil
}
