package usecase

import (
	"context"
	"encoding/base64"
	"fmt"

	"opscore/backend/internal/document/application/dto"
	apperror "opscore/backend/internal/document/application/error"
	"opscore/backend/internal/document/domain/entity"
	"opscore/backend/internal/document/domain/repository"
	"opscore/backend/internal/document/domain/value_object"
	"opscore/backend/internal/document/infrastructure/parser"
	"opscore/backend/internal/document/infrastructure/storage"
	gitentity "opscore/backend/internal/git_repository/domain/entity"
	gitrepo "opscore/backend/internal/git_repository/domain/repository"
	"opscore/backend/internal/git_repository/infrastructure/git"
	oauthdomain "opscore/backend/internal/oauth/domain"
)

// publishDocumentUseCase handles publishing documents from OAuth connections.
type publishDocumentUseCase struct {
	repo               repository.DocumentRepository
	gitRepo            gitrepo.Repository
	gitManager         git.GitManager
	fmParser           parser.FrontmatterParser
	oauthService       OAuthService
	gitProviderService GitProviderService
	persister          contentPersister
}

func newPublishDocumentUseCase(
	repo repository.DocumentRepository,
	gitRepo gitrepo.Repository,
	gitManager git.GitManager,
	fmParser parser.FrontmatterParser,
	oauthService OAuthService,
	gitProviderService GitProviderService,
	storage storage.DocumentStorage,
) publishDocumentUseCase {
	return publishDocumentUseCase{
		repo:               repo,
		gitRepo:            gitRepo,
		gitManager:         gitManager,
		fmParser:           fmParser,
		oauthService:       oauthService,
		gitProviderService: gitProviderService,
		persister:          newContentPersister(storage),
	}
}

func (uc publishDocumentUseCase) Execute(ctx context.Context, userID string, req *dto.PublishDocumentRequest) (*dto.DocumentResponse, error) {
	accessScope, err := value_object.NewAccessScope(req.AccessScope)
	if err != nil {
		return nil, apperror.NewValidationFailedError([]apperror.FieldError{{Field: "access_scope", Message: err.Error()}})
	}

	filePath, err := value_object.NewFilePath(req.FilePath)
	if err != nil {
		return nil, apperror.NewValidationFailedError([]apperror.FieldError{{Field: "file_path", Message: err.Error()}})
	}

	conn, err := uc.oauthService.GetConnectionByID(ctx, userID, req.ConnectionID)
	if err != nil {
		return nil, fmt.Errorf("failed to get OAuth connection: %w", err)
	}

	accessToken, err := uc.oauthService.GetAccessTokenByConnectionID(ctx, userID, req.ConnectionID)
	if err != nil {
		return nil, fmt.Errorf("failed to get access token: %w", err)
	}

	var repoURL string
	provider := conn.Provider()
	switch provider {
	case oauthdomain.ProviderGitHub:
		repoURL = fmt.Sprintf("https://github.com/%s/%s", req.Owner, req.Repository)
	case oauthdomain.ProviderGitLab:
		repoURL = fmt.Sprintf("https://gitlab.com/%s/%s", req.Owner, req.Repository)
	case oauthdomain.ProviderGitLabSelfHosted:
		repoURL = fmt.Sprintf("https://%s/%s/%s", conn.ProviderHost(), req.Owner, req.Repository)
	default:
		return nil, apperror.NewValidationFailedError([]apperror.FieldError{{Field: "connection_id", Message: fmt.Sprintf("unsupported provider: %s", provider)}})
	}

	repoEntity, err := uc.gitRepo.FindByURL(ctx, repoURL)
	if err != nil {
		return nil, fmt.Errorf("failed to find repository: %w", err)
	}

	if repoEntity == nil {
		repoID := value_object.GenerateDocumentID() // reuse generator for repository ID
		repoName := fmt.Sprintf("%s/%s", req.Owner, req.Repository)

		newRepo := gitentity.NewRepository(repoID.String(), repoName, repoURL, accessToken)

		if err := uc.gitRepo.Save(ctx, newRepo); err != nil {
			return nil, fmt.Errorf("failed to save repository: %w", err)
		}
		repoEntity = newRepo
	} else {
		repoEntity.SetAccessToken(accessToken)
		if err := uc.gitRepo.Save(ctx, repoEntity); err != nil {
			return nil, fmt.Errorf("failed to update repository: %w", err)
		}
	}

	fileContent, err := uc.gitProviderService.GetFileContent(ctx, userID, req.ConnectionID, req.ProviderRepositoryID, req.Owner, req.Repository, req.FilePath, req.Ref)
	if err != nil {
		return nil, fmt.Errorf("failed to get file content: %w", err)
	}

	var markdownContent string
	if fileContent.Encoding == "base64" {
		decoded, decErr := base64.StdEncoding.DecodeString(fileContent.Content)
		if decErr != nil {
			return nil, fmt.Errorf("failed to decode base64 content: %w", decErr)
		}
		markdownContent = string(decoded)
	} else {
		markdownContent = fileContent.Content
	}

	frontmatterData, err := uc.fmParser.Parse(markdownContent)
	if err != nil {
		return nil, apperror.NewValidationFailedError([]apperror.FieldError{{Field: "file_content", Message: fmt.Sprintf("failed to parse frontmatter: %s", err.Error())}})
	}

	docType, err := value_object.NewDocumentType(frontmatterData.Type)
	if err != nil {
		return nil, apperror.NewValidationFailedError([]apperror.FieldError{{Field: "doc_type", Message: err.Error()}})
	}

	tags := make([]value_object.Tag, len(frontmatterData.Tags))
	for i, tagStr := range frontmatterData.Tags {
		tag, tagErr := value_object.NewTag(tagStr)
		if tagErr != nil {
			return nil, apperror.NewValidationFailedError([]apperror.FieldError{{Field: fmt.Sprintf("tags[%d]", i), Message: tagErr.Error()}})
		}
		tags[i] = tag
	}

	variables := make([]value_object.VariableDefinition, len(frontmatterData.Variables))
	for i, v := range frontmatterData.Variables {
		varType, typeErr := value_object.NewVariableType(v.Type)
		if typeErr != nil {
			return nil, apperror.NewValidationFailedError([]apperror.FieldError{{Field: fmt.Sprintf("variables[%d].type", i), Message: typeErr.Error()}})
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
			return nil, apperror.NewValidationFailedError([]apperror.FieldError{{Field: fmt.Sprintf("variables[%d]", i), Message: defErr.Error()}})
		}
		variables[i] = varDef
	}

	commitHash, err := value_object.NewCommitHash(fileContent.SHA)
	if err != nil {
		commitHash, _ = value_object.NewCommitHash("HEAD")
	}

	source, err := value_object.NewDocumentSource(filePath, commitHash)
	if err != nil {
		return nil, apperror.NewValidationFailedError([]apperror.FieldError{{Field: "source", Message: err.Error()}})
	}

	documentID := value_object.GenerateDocumentID()

	repositoryID, err := value_object.NewRepositoryID(repoEntity.ID())
	if err != nil {
		return nil, fmt.Errorf("failed to create repository ID: %w", err)
	}

	origin, err := value_object.NewRepositoryOrigin(req.ProviderRepositoryID, req.Owner, req.Repository)
	if err != nil {
		return nil, fmt.Errorf("failed to create repository origin: %w", err)
	}

	doc, err := entity.NewDocument(documentID, repositoryID, &origin, accessScope)
	if err != nil {
		return nil, fmt.Errorf("failed to create document: %w", err)
	}

	if req.IsAutoUpdate {
		doc.EnableAutoUpdate()
	}

	if err := doc.Publish(source, frontmatterData.Title, docType, tags, variables, frontmatterData.Content); err != nil {
		return nil, fmt.Errorf("failed to publish initial version: %w", err)
	}

	if err := uc.repo.Save(ctx, doc); err != nil {
		return nil, fmt.Errorf("failed to save document: %w", err)
	}

	uc.persister.persistContent(ctx, doc)

	response := dto.ToDocumentResponse(doc)
	return &response, nil
}
