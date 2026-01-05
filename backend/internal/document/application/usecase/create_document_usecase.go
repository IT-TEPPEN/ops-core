package usecase

import (
	"context"
	"fmt"
	"strings"

	"opscore/backend/internal/document/application/dto"
	apperror "opscore/backend/internal/document/application/error"
	"opscore/backend/internal/document/domain/entity"
	"opscore/backend/internal/document/domain/repository"
	"opscore/backend/internal/document/domain/value_object"
	"opscore/backend/internal/document/infrastructure/parser"
	"opscore/backend/internal/document/infrastructure/storage"
	gitrepo "opscore/backend/internal/git_repository/domain/repository"
	"opscore/backend/internal/git_repository/infrastructure/git"
)

// createDocumentUseCase handles creation of documents.
type createDocumentUseCase struct {
	repo       repository.DocumentRepository
	gitRepo    gitrepo.Repository
	gitManager git.GitManager
	fmParser   parser.FrontmatterParser
	persister  contentPersister
}

func newCreateDocumentUseCase(
	repo repository.DocumentRepository,
	gitRepo gitrepo.Repository,
	gitManager git.GitManager,
	fmParser parser.FrontmatterParser,
	storage storage.DocumentStorage,
) createDocumentUseCase {
	return createDocumentUseCase{
		repo:       repo,
		gitRepo:    gitRepo,
		gitManager: gitManager,
		fmParser:   fmParser,
		persister:  newContentPersister(storage),
	}
}

func (uc createDocumentUseCase) Execute(ctx context.Context, req *dto.CreateDocumentRequest) (*dto.DocumentResponse, error) {
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

	// Validate repository origin
	origin, err := value_object.NewRepositoryOrigin(req.ProviderRepositoryID, req.Owner, req.Repository)
	if err != nil {
		field := "repository_origin"
		switch {
		case strings.Contains(err.Error(), "provider repository"):
			field = "provider_repository_id"
		case strings.Contains(err.Error(), "owner"):
			field = "owner"
		case strings.Contains(err.Error(), "name"):
			field = "repository"
		}
		return nil, apperror.NewValidationFailedError([]apperror.FieldError{
			{Field: field, Message: err.Error()},
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
		tag, tagErr := value_object.NewTag(tagStr)
		if tagErr != nil {
			return nil, apperror.NewValidationFailedError([]apperror.FieldError{
				{Field: fmt.Sprintf("tags[%d]", i), Message: tagErr.Error()},
			})
		}
		tags[i] = tag
	}

	// Convert variables from frontmatter
	variables := make([]value_object.VariableDefinition, len(frontmatterData.Variables))
	for i, v := range frontmatterData.Variables {
		varType, typeErr := value_object.NewVariableType(v.Type)
		if typeErr != nil {
			return nil, apperror.NewValidationFailedError([]apperror.FieldError{
				{Field: fmt.Sprintf("variables[%d].type", i), Message: typeErr.Error()},
			})
		}
		varDef, defErr := value_object.NewVariableDefinition(
			v.Name,
			v.Label,
			v.Description,
			varType,
			v.Required,
			v.DefaultValue,
		)
		if defErr != nil {
			return nil, apperror.NewValidationFailedError([]apperror.FieldError{
				{Field: fmt.Sprintf("variables[%d]", i), Message: defErr.Error()},
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

	// Create the document entity with origin metadata
	doc, err := entity.NewDocument(documentID, repositoryID, &origin, accessScope)
	if err != nil {
		return nil, fmt.Errorf("failed to create document: %w", err)
	}

	// Set auto update if specified
	if req.IsAutoUpdate {
		doc.EnableAutoUpdate()
	}

	// Publish the initial version
	if err := doc.Publish(source, frontmatterData.Title, docType, tags, variables, frontmatterData.Content); err != nil {
		return nil, fmt.Errorf("failed to publish initial version: %w", err)
	}

	// Save the document
	if err := uc.repo.Save(ctx, doc); err != nil {
		return nil, fmt.Errorf("failed to save document: %w", err)
	}

	uc.persister.persistContent(ctx, doc)

	response := dto.ToDocumentResponse(doc)
	return &response, nil
}
