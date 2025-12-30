package usecase

import (
	"context"
	"fmt"

	apperror "opscore/backend/internal/document/application/error"
	"opscore/backend/internal/document/application/dto"
	"opscore/backend/internal/document/domain/entity"
	"opscore/backend/internal/document/domain/repository"
	"opscore/backend/internal/document/domain/value_object"
	"opscore/backend/internal/document/infrastructure/parser"
	gitrepo "opscore/backend/internal/git_repository/domain/repository"
	"opscore/backend/internal/git_repository/infrastructure/git"
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
	repo        repository.DocumentRepository
	gitRepo     gitrepo.Repository
	gitManager  git.GitManager
	fmParser    parser.FrontmatterParser
}

// NewDocumentUseCase creates a new instance of documentUseCase.
func NewDocumentUseCase(
	repo repository.DocumentRepository,
	gitRepo gitrepo.Repository,
	gitManager git.GitManager,
	fmParser parser.FrontmatterParser,
) DocumentUseCase {
	return &documentUseCase{
		repo:       repo,
		gitRepo:    gitRepo,
		gitManager: gitManager,
		fmParser:   fmParser,
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

	// Validate file path
	filePath, err := value_object.NewFilePath(req.FilePath)
	if err != nil {
		return nil, apperror.NewValidationFailedError([]apperror.FieldError{
			{Field: "file_path", Message: err.Error()},
		})
	}

	// Fetch repository entity
	repoEntity, err := uc.gitRepo.FindByID(ctx, req.RepositoryID)
	if err != nil {
		return nil, fmt.Errorf("failed to find repository: %w", err)
	}
	if repoEntity == nil {
		return nil, apperror.NewNotFoundError("Repository", req.RepositoryID, nil)
	}

	// Ensure repository is cloned
	localPath, err := uc.gitManager.EnsureCloned(ctx, repoEntity)
	if err != nil {
		return nil, fmt.Errorf("failed to clone repository: %w", err)
	}

	// Read file content from repository
	fileContent, err := uc.gitManager.ReadManagedFileContent(ctx, localPath, req.FilePath, repoEntity)
	if err != nil {
		return nil, fmt.Errorf("failed to read file content: %w", err)
	}

	// Parse frontmatter from markdown file
	frontmatterData, err := uc.fmParser.Parse(string(fileContent))
	if err != nil {
		return nil, apperror.NewValidationFailedError([]apperror.FieldError{
			{Field: "file_content", Message: fmt.Sprintf("failed to parse frontmatter: %s", err.Error())},
		})
	}

	// Validate and create document type
	docType, err := value_object.NewDocumentType(frontmatterData.Type)
	if err != nil {
		return nil, apperror.NewValidationFailedError([]apperror.FieldError{
			{Field: "doc_type", Message: err.Error()},
		})
	}

	// Convert tags from frontmatter
	tags := make([]value_object.Tag, len(frontmatterData.Tags))
	for i, tagStr := range frontmatterData.Tags {
		tag, err := value_object.NewTag(tagStr)
		if err != nil {
			return nil, apperror.NewValidationFailedError([]apperror.FieldError{
				{Field: fmt.Sprintf("tags[%d]", i), Message: err.Error()},
			})
		}
		tags[i] = tag
	}

	// Convert variables from frontmatter
	variables := make([]value_object.VariableDefinition, len(frontmatterData.Variables))
	for i, v := range frontmatterData.Variables {
		varType, err := value_object.NewVariableType(v.Type)
		if err != nil {
			return nil, apperror.NewValidationFailedError([]apperror.FieldError{
				{Field: fmt.Sprintf("variables[%d].type", i), Message: err.Error()},
			})
		}
		varDef, err := value_object.NewVariableDefinition(
			v.Name,
			v.Label,
			v.Description,
			varType,
			v.Required,
			v.DefaultValue,
		)
		if err != nil {
			return nil, apperror.NewValidationFailedError([]apperror.FieldError{
				{Field: fmt.Sprintf("variables[%d]", i), Message: err.Error()},
			})
		}
		variables[i] = varDef
	}

	// Determine commit hash (use provided or "HEAD" as default)
	commitHashStr := req.CommitHash
	if commitHashStr == "" {
		commitHashStr = "HEAD"
	}
	commitHash, err := value_object.NewCommitHash(commitHashStr)
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

	// Generate new document ID
	documentID := value_object.GenerateDocumentID()

	// Create the document entity
	doc, err := entity.NewDocument(documentID, repositoryID, frontmatterData.Owner, accessScope)
	if err != nil {
		return nil, fmt.Errorf("failed to create document: %w", err)
	}

	// Set auto update if specified
	if req.IsAutoUpdate {
		doc.EnableAutoUpdate()
	}

	// Publish the initial version
	err = doc.Publish(source, frontmatterData.Title, docType, tags, variables, frontmatterData.Content)
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

	// Validate file path
	filePath, err := value_object.NewFilePath(req.FilePath)
	if err != nil {
		return nil, apperror.NewValidationFailedError([]apperror.FieldError{
			{Field: "file_path", Message: err.Error()},
		})
	}

	// Fetch repository entity
	repoEntity, err := uc.gitRepo.FindByID(ctx, doc.RepositoryID().String())
	if err != nil {
		return nil, fmt.Errorf("failed to find repository: %w", err)
	}
	if repoEntity == nil {
		return nil, apperror.NewNotFoundError("Repository", doc.RepositoryID().String(), nil)
	}

	// Ensure repository is cloned
	localPath, err := uc.gitManager.EnsureCloned(ctx, repoEntity)
	if err != nil {
		return nil, fmt.Errorf("failed to clone repository: %w", err)
	}

	// Read file content from repository
	fileContent, err := uc.gitManager.ReadManagedFileContent(ctx, localPath, req.FilePath, repoEntity)
	if err != nil {
		return nil, fmt.Errorf("failed to read file content: %w", err)
	}

	// Parse frontmatter from markdown file
	frontmatterData, err := uc.fmParser.Parse(string(fileContent))
	if err != nil {
		return nil, apperror.NewValidationFailedError([]apperror.FieldError{
			{Field: "file_content", Message: fmt.Sprintf("failed to parse frontmatter: %s", err.Error())},
		})
	}

	// Validate and create document type
	docType, err := value_object.NewDocumentType(frontmatterData.Type)
	if err != nil {
		return nil, apperror.NewValidationFailedError([]apperror.FieldError{
			{Field: "doc_type", Message: err.Error()},
		})
	}

	// Convert tags from frontmatter
	tags := make([]value_object.Tag, len(frontmatterData.Tags))
	for i, tagStr := range frontmatterData.Tags {
		tag, err := value_object.NewTag(tagStr)
		if err != nil {
			return nil, apperror.NewValidationFailedError([]apperror.FieldError{
				{Field: fmt.Sprintf("tags[%d]", i), Message: err.Error()},
			})
		}
		tags[i] = tag
	}

	// Convert variables from frontmatter
	variables := make([]value_object.VariableDefinition, len(frontmatterData.Variables))
	for i, v := range frontmatterData.Variables {
		varType, err := value_object.NewVariableType(v.Type)
		if err != nil {
			return nil, apperror.NewValidationFailedError([]apperror.FieldError{
				{Field: fmt.Sprintf("variables[%d].type", i), Message: err.Error()},
			})
		}
		varDef, err := value_object.NewVariableDefinition(
			v.Name,
			v.Label,
			v.Description,
			varType,
			v.Required,
			v.DefaultValue,
		)
		if err != nil {
			return nil, apperror.NewValidationFailedError([]apperror.FieldError{
				{Field: fmt.Sprintf("variables[%d]", i), Message: err.Error()},
			})
		}
		variables[i] = varDef
	}

	// Determine commit hash (use provided or "HEAD" as default)
	commitHashStr := req.CommitHash
	if commitHashStr == "" {
		commitHashStr = "HEAD"
	}
	commitHash, err := value_object.NewCommitHash(commitHashStr)
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

	// Publish the new version
	err = doc.Publish(source, frontmatterData.Title, docType, tags, variables, frontmatterData.Content)
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
		// TODO: FIXME - This is a workaround. The PublishWithVersion method doesn't exist in the domain entity.
		// RollbackToVersion is being used here but it may have different semantics.
		// The proper fix would be to:
		// 1. Add PublishWithVersion method to the Document entity that accepts a version number
		// 2. Implement the logic to mark a specific version as published
		// 3. Update this code to use the correct method
		// For now, this allows the code to compile and may work, but should be reviewed.
		err = doc.RollbackToVersion(verNum)
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
