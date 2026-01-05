package repository

import (
	"context"

	apperror "opscore/backend/internal/git_repository/application/error"
	"opscore/backend/internal/git_repository/domain/entity"
	"opscore/backend/internal/git_repository/domain/repository"
	"opscore/backend/internal/git_repository/infrastructure/git"
)

// Legacy error definitions - kept for backward compatibility during migration.
// These will be removed after full migration.
var (
	ErrRepositoryNotFound = apperror.ErrNotFound
)

// OAuthTokenProvider provides OAuth access tokens for Git providers.
type OAuthTokenProvider interface {
	GetAccessTokenForProvider(ctx context.Context, userID string, providerName string) (string, error)
}

// RepositoryUseCase defines the interface for repository related use cases.
type RepositoryUseCase interface {
	GetRepository(ctx context.Context, repoID string) (entity.Repository, error)
	ListRepositories(ctx context.Context) ([]entity.Repository, error)
	ListFiles(ctx context.Context, repoID string, userID string) ([]entity.FileNode, error)
	GetDirectoryContents(ctx context.Context, repoID string, path string, userID string) ([]entity.FileNode, error)
	GetFileContents(ctx context.Context, repoID string, filePath string, userID string, commitHash string) (string, string, error)
	GetLatestCommit(ctx context.Context, repoID string, userID string) (*git.CommitInfo, error)
	GetFileCommitHistory(ctx context.Context, repoID string, filePath string, userID string) ([]git.CommitInfo, error)
}

// repositoryUseCase is a facade delegating to per-usecase executors.
type repositoryUseCase struct {
	get             getRepositoryUseCase
	list            listRepositoriesUseCase
	listFiles       listFilesUseCase
	dirContents     getDirectoryContentsUseCase
	fileContents    getFileContentsUseCase
	latestCommit    getLatestCommitUseCase
	fileCommitHists getFileCommitHistoryUseCase
}

// NewRepositoryUseCase creates a new instance wired with executors.
func NewRepositoryUseCase(repo repository.Repository, gitManager git.GitManager, oauthProvider OAuthTokenProvider) RepositoryUseCase {
	return &repositoryUseCase{
		get:             newGetRepositoryUseCase(repo),
		list:            newListRepositoriesUseCase(repo),
		listFiles:       newListFilesUseCase(repo, gitManager, oauthProvider),
		dirContents:     newGetDirectoryContentsUseCase(repo, gitManager, oauthProvider),
		fileContents:    newGetFileContentsUseCase(repo, gitManager, oauthProvider),
		latestCommit:    newGetLatestCommitUseCase(repo, gitManager, oauthProvider),
		fileCommitHists: newGetFileCommitHistoryUseCase(repo, gitManager, oauthProvider),
	}
}

func (uc *repositoryUseCase) GetRepository(ctx context.Context, repoID string) (entity.Repository, error) {
	return uc.get.Execute(ctx, repoID)
}

func (uc *repositoryUseCase) ListRepositories(ctx context.Context) ([]entity.Repository, error) {
	return uc.list.Execute(ctx)
}

func (uc *repositoryUseCase) ListFiles(ctx context.Context, repoID string, userID string) ([]entity.FileNode, error) {
	return uc.listFiles.Execute(ctx, repoID, userID)
}

func (uc *repositoryUseCase) GetDirectoryContents(ctx context.Context, repoID string, path string, userID string) ([]entity.FileNode, error) {
	return uc.dirContents.Execute(ctx, repoID, path, userID)
}

func (uc *repositoryUseCase) GetFileContents(ctx context.Context, repoID string, filePath string, userID string, commitHash string) (string, string, error) {
	return uc.fileContents.Execute(ctx, repoID, filePath, userID, commitHash)
}

func (uc *repositoryUseCase) GetLatestCommit(ctx context.Context, repoID string, userID string) (*git.CommitInfo, error) {
	return uc.latestCommit.Execute(ctx, repoID, userID)
}

func (uc *repositoryUseCase) GetFileCommitHistory(ctx context.Context, repoID string, filePath string, userID string) ([]git.CommitInfo, error) {
	return uc.fileCommitHists.Execute(ctx, repoID, filePath, userID)
}
