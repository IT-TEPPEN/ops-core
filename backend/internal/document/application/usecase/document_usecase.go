package usecase

import (
	"context"

	"opscore/backend/internal/document/application/dto"
	"opscore/backend/internal/document/domain/repository"
	"opscore/backend/internal/document/infrastructure/parser"
	"opscore/backend/internal/document/infrastructure/storage"
	gitrepo "opscore/backend/internal/git_repository/domain/repository"
	"opscore/backend/internal/git_repository/infrastructure/git"
	oauthservice "opscore/backend/internal/oauth/application/service"
	oauthdomain "opscore/backend/internal/oauth/domain"
)

// DocumentUseCase defines the interface for document related use cases.
type DocumentUseCase interface {
	CreateDocument(ctx context.Context, req *dto.CreateDocumentRequest) (*dto.DocumentResponse, error)
	UpdateDocument(ctx context.Context, documentID string, req *dto.UpdateDocumentRequest) (*dto.DocumentResponse, error)
	GetDocument(ctx context.Context, documentID string) (*dto.DocumentResponse, error)
	GetDocumentVersion(ctx context.Context, documentID string, versionNumber int) (*dto.DocumentVersionResponse, error)
	ListDocuments(ctx context.Context, filter dto.DocumentListFilter) ([]dto.DocumentListItemResponse, error)
	GetDocumentVersions(ctx context.Context, documentID string) (*dto.VersionHistoryResponse, error)
	PublishDocumentVersion(ctx context.Context, documentID string, versionNumber int) (*dto.DocumentResponse, error)
	RollbackDocumentVersion(ctx context.Context, documentID string, versionNumber int) (*dto.DocumentResponse, error)
	UpdateDocumentMetadata(ctx context.Context, documentID string, req *dto.UpdateDocumentMetadataRequest) (*dto.DocumentResponse, error)
	PublishDocument(ctx context.Context, userID string, req *dto.PublishDocumentRequest) (*dto.DocumentResponse, error)
}

// OAuthService interface for getting OAuth connections.
type OAuthService interface {
	GetConnectionByID(ctx context.Context, userID string, connectionID string) (*oauthdomain.OAuthConnection, error)
	GetAccessTokenByConnectionID(ctx context.Context, userID string, connectionID string) (string, error)
}

// GitProviderService interface for getting file content from Git providers.
type GitProviderService interface {
	GetFileContent(ctx context.Context, userID string, connectionID string, repositoryID string, owner string, repo string, filePath string, ref string) (*oauthservice.FileContent, error)
}

// documentUseCase is a thin facade delegating to per-usecase executors.
type documentUseCase struct {
	create         createDocumentUseCase
	update         updateDocumentUseCase
	get            getDocumentUseCase
	getVersion     getDocumentVersionUseCase
	list           listDocumentsUseCase
	getVersions    getDocumentVersionsUseCase
	publishVersion publishDocumentVersionUseCase
	rollback       rollbackDocumentVersionUseCase
	updateMetadata updateDocumentMetadataUseCase
	publish        publishDocumentUseCase
}

// NewDocumentUseCase wires a facade that delegates to per-usecase executors.
func NewDocumentUseCase(
	repo repository.DocumentRepository,
	gitRepo gitrepo.Repository,
	gitManager git.GitManager,
	fmParser parser.FrontmatterParser,
	oauthService OAuthService,
	gitProviderService GitProviderService,
	storage storage.DocumentStorage,
) DocumentUseCase {
	return &documentUseCase{
		create:         newCreateDocumentUseCase(repo, gitRepo, gitManager, fmParser, storage),
		update:         newUpdateDocumentUseCase(repo, gitRepo, gitManager, fmParser, storage),
		get:            newGetDocumentUseCase(repo),
		getVersion:     newGetDocumentVersionUseCase(repo),
		list:           newListDocumentsUseCase(repo),
		getVersions:    newGetDocumentVersionsUseCase(repo),
		publishVersion: newPublishDocumentVersionUseCase(repo, storage),
		rollback:       newRollbackDocumentVersionUseCase(repo),
		updateMetadata: newUpdateDocumentMetadataUseCase(repo),
		publish:        newPublishDocumentUseCase(repo, gitRepo, gitManager, fmParser, oauthService, gitProviderService, storage),
	}
}

func (uc *documentUseCase) CreateDocument(ctx context.Context, req *dto.CreateDocumentRequest) (*dto.DocumentResponse, error) {
	return uc.create.Execute(ctx, req)
}

func (uc *documentUseCase) UpdateDocument(ctx context.Context, documentID string, req *dto.UpdateDocumentRequest) (*dto.DocumentResponse, error) {
	return uc.update.Execute(ctx, documentID, req)
}

func (uc *documentUseCase) GetDocument(ctx context.Context, documentID string) (*dto.DocumentResponse, error) {
	return uc.get.Execute(ctx, documentID)
}

func (uc *documentUseCase) GetDocumentVersion(ctx context.Context, documentID string, versionNumber int) (*dto.DocumentVersionResponse, error) {
	return uc.getVersion.Execute(ctx, documentID, versionNumber)
}

func (uc *documentUseCase) ListDocuments(ctx context.Context, filter dto.DocumentListFilter) ([]dto.DocumentListItemResponse, error) {
	return uc.list.Execute(ctx, filter)
}

func (uc *documentUseCase) GetDocumentVersions(ctx context.Context, documentID string) (*dto.VersionHistoryResponse, error) {
	return uc.getVersions.Execute(ctx, documentID)
}

func (uc *documentUseCase) PublishDocumentVersion(ctx context.Context, documentID string, versionNumber int) (*dto.DocumentResponse, error) {
	return uc.publishVersion.Execute(ctx, documentID, versionNumber)
}

func (uc *documentUseCase) RollbackDocumentVersion(ctx context.Context, documentID string, versionNumber int) (*dto.DocumentResponse, error) {
	return uc.rollback.Execute(ctx, documentID, versionNumber)
}

func (uc *documentUseCase) UpdateDocumentMetadata(ctx context.Context, documentID string, req *dto.UpdateDocumentMetadataRequest) (*dto.DocumentResponse, error) {
	return uc.updateMetadata.Execute(ctx, documentID, req)
}

func (uc *documentUseCase) PublishDocument(ctx context.Context, userID string, req *dto.PublishDocumentRequest) (*dto.DocumentResponse, error) {
	return uc.publish.Execute(ctx, userID, req)
}
