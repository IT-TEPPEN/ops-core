package usecase

import (
	"context"
	"fmt"

	"opscore/backend/internal/document/application/dto"
	apperror "opscore/backend/internal/document/application/error"
	"opscore/backend/internal/document/domain/repository"
	"opscore/backend/internal/document/domain/value_object"
	"opscore/backend/internal/document/infrastructure/parser"
	"opscore/backend/internal/document/infrastructure/storage"
	gitrepo "opscore/backend/internal/git_repository/domain/repository"
	"opscore/backend/internal/git_repository/infrastructure/git"
)

// updateDocumentUseCase handles updating documents by creating a new version.
type updateDocumentUseCase struct {
	repo       repository.DocumentRepository
	gitRepo    gitrepo.Repository
	gitManager git.GitManager
	fmParser   parser.FrontmatterParser
	persister  contentPersister
}

func newUpdateDocumentUseCase(
	repo repository.DocumentRepository,
	gitRepo gitrepo.Repository,
	gitManager git.GitManager,
	fmParser parser.FrontmatterParser,
	storage storage.DocumentStorage,
) updateDocumentUseCase {
	return updateDocumentUseCase{
		repo:       repo,
		gitRepo:    gitRepo,
		gitManager: gitManager,
		fmParser:   fmParser,
		persister:  newContentPersister(storage),
	}
}

func (uc updateDocumentUseCase) Execute(ctx context.Context, documentID string, req *dto.UpdateDocumentRequest) (*dto.DocumentResponse, error) {
	docID, err := value_object.NewDocumentID(documentID)
	if err != nil {
		return nil, apperror.NewValidationFailedError([]apperror.FieldError{
			{Field: "document_id", Message: err.Error()},
		})
	}

	doc, err := uc.repo.FindByID(ctx, docID)
	if err != nil {
		return nil, fmt.Errorf("failed to find document: %w", err)
	}
	if doc == nil {
		return nil, apperror.NewNotFoundError("Document", documentID, nil)
	}

	filePath, err := value_object.NewFilePath(req.FilePath)
	if err != nil {
		return nil, apperror.NewValidationFailedError([]apperror.FieldError{
			{Field: "file_path", Message: err.Error()},
		})
	}

	repoEntity, err := uc.gitRepo.FindByID(ctx, doc.RepositoryID().String())
	if err != nil {
		return nil, fmt.Errorf("failed to find repository: %w", err)
	}
	if repoEntity == nil {
		return nil, apperror.NewNotFoundError("Repository", doc.RepositoryID().String(), nil)
	}

	localPath, err := uc.gitManager.EnsureCloned(ctx, repoEntity)
	if err != nil {
		return nil, fmt.Errorf("failed to clone repository: %w", err)
	}

	fileContent, err := uc.gitManager.ReadManagedFileContent(ctx, localPath, req.FilePath, repoEntity)
	if err != nil {
		return nil, fmt.Errorf("failed to read file content: %w", err)
	}

	frontmatterData, err := uc.fmParser.Parse(string(fileContent))
	if err != nil {
		return nil, apperror.NewValidationFailedError([]apperror.FieldError{
			{Field: "file_content", Message: fmt.Sprintf("failed to parse frontmatter: %s", err.Error())},
		})
	}

	docType, err := value_object.NewDocumentType(frontmatterData.Type)
	if err != nil {
		return nil, apperror.NewValidationFailedError([]apperror.FieldError{
			{Field: "doc_type", Message: err.Error()},
		})
	}

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

	source, err := value_object.NewDocumentSource(filePath, commitHash)
	if err != nil {
		return nil, apperror.NewValidationFailedError([]apperror.FieldError{
			{Field: "source", Message: err.Error()},
		})
	}

	if err := doc.Publish(source, frontmatterData.Title, docType, tags, variables, frontmatterData.Content); err != nil {
		return nil, fmt.Errorf("failed to publish new version: %w", err)
	}

	if err := uc.repo.Update(ctx, doc); err != nil {
		return nil, fmt.Errorf("failed to update document: %w", err)
	}

	uc.persister.persistContent(ctx, doc)

	response := dto.ToDocumentResponse(doc)
	return &response, nil
}
