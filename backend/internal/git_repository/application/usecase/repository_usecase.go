package repository

import (
	"context"
	"errors"
	"fmt"
	"strings"

	apperror "opscore/backend/internal/git_repository/application/error"
	"opscore/backend/internal/git_repository/domain/entity"
	"opscore/backend/internal/git_repository/domain/repository"
	"opscore/backend/internal/git_repository/infrastructure/git"
)

// Legacy error definitions - kept for backward compatibility during migration
// These will be removed after full migration
var (
	ErrRepositoryNotFound = errors.New("repository not found")
)

// OAuthTokenProvider provides OAuth access tokens for Git providers
type OAuthTokenProvider interface {
	GetAccessTokenForProvider(ctx context.Context, userID string, providerName string) (string, error)
}

// RepositoryUseCase defines the interface for repository related use cases.
type RepositoryUseCase interface {
	// GetRepository retrieves a single repository by its ID
	GetRepository(ctx context.Context, repoID string) (entity.Repository, error)
	// ListRepositories retrieves all registered repositories
	ListRepositories(ctx context.Context) ([]entity.Repository, error)
	// ListFiles retrieves the file structure for a given repository ID.
	ListFiles(ctx context.Context, repoID string, userID string) ([]entity.FileNode, error) // Use entity.FileNode
	// GetDirectoryContents retrieves files and directories at a specific path (non-recursive).
	GetDirectoryContents(ctx context.Context, repoID string, path string, userID string) ([]entity.FileNode, error)
	// GetFileContents retrieves the content of a specific file from a repository.
	// If commitHash is empty, retrieves the latest version. Returns content and actual commit hash.
	GetFileContents(ctx context.Context, repoID string, filePath string, userID string, commitHash string) (string, string, error)
	// GetLatestCommit retrieves the latest commit information for a repository.
	GetLatestCommit(ctx context.Context, repoID string, userID string) (*git.CommitInfo, error)
	// GetFileCommitHistory retrieves the commit history for a specific file.
	GetFileCommitHistory(ctx context.Context, repoID string, filePath string, userID string) ([]git.CommitInfo, error)
}

// repositoryUseCase implements the RepositoryUseCase interface.
type repositoryUseCase struct {
	repo          repository.Repository // Persistence for repository metadata
	gitManager    git.GitManager        // For interacting with Git repositories
	oauthProvider OAuthTokenProvider    // For getting OAuth access tokens
}

// NewRepositoryUseCase creates a new instance of repositoryUseCase.
func NewRepositoryUseCase(repo repository.Repository, gitManager git.GitManager, oauthProvider OAuthTokenProvider) RepositoryUseCase {
	return &repositoryUseCase{
		repo:          repo,
		gitManager:    gitManager,
		oauthProvider: oauthProvider,
	}
}

// getProviderFromURL extracts the Git provider name from a repository URL
func getProviderFromURL(repoURL string) string {
	if strings.Contains(repoURL, "github.com") {
		return "github"
	} else if strings.Contains(repoURL, "gitlab.com") {
		return "gitlab"
	}
	return ""
}

// GetRepository retrieves a single repository by its ID.
func (uc *repositoryUseCase) GetRepository(ctx context.Context, repoID string) (entity.Repository, error) {
	// Find the repository by ID
	repo, err := uc.repo.FindByID(ctx, repoID)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve repository: %w", err)
	}

	if repo == nil {
		return nil, apperror.NewNotFoundError("Repository", repoID, nil)
	}

	return repo, nil
}

// ListRepositories implements the logic for retrieving all registered repositories.
func (uc *repositoryUseCase) ListRepositories(ctx context.Context) ([]entity.Repository, error) {
	// Get all repositories from the repository layer
	repos, err := uc.repo.FindAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve repositories: %w", err)
	}

	return repos, nil
}

// ListFiles implements the logic for listing files in a repository.
func (uc *repositoryUseCase) ListFiles(ctx context.Context, repoID string, userID string) ([]entity.FileNode, error) {
	// 1. Find the repository by ID
	repo, err := uc.repo.FindByID(ctx, repoID)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve repository details: %w", err)
	}
	if repo == nil {
		return nil, apperror.NewNotFoundError("Repository", repoID, nil)
	}

	// 2. Get OAuth access token for the provider
	provider := getProviderFromURL(repo.URL())
	if provider == "" {
		return nil, apperror.NewValidationFailedError([]apperror.FieldError{
			{Field: "url", Message: "unsupported Git provider"},
		})
	}

	accessToken, err := uc.oauthProvider.GetAccessTokenForProvider(ctx, userID, provider)
	if err != nil {
		return nil, apperror.NewValidationFailedError([]apperror.FieldError{
			{Field: "oauth", Message: fmt.Sprintf("OAuth connection required for %s. Please connect your account.", provider)},
		})
	}

	// Set the OAuth token on the repository for GitManager to use
	repo.SetAccessToken(accessToken)

	// 3. List files directly from GitHub API (not from local cache)
	files, err := uc.gitManager.ListRepositoryFiles(ctx, "", repo)
	if err != nil {
		return nil, fmt.Errorf("failed to list repository files: %w", err)
	}

	// 4. Map the output to []entity.FileNode (Basic mapping)
	fileNodes := make([]entity.FileNode, 0, len(files))
	for _, f := range files {
		fileType := "file"
		fileNodes = append(fileNodes, entity.NewFileNode(f, fileType))
	}

	return fileNodes, nil
}

// GetDirectoryContents implements the logic for listing directory contents at a specific path.
func (uc *repositoryUseCase) GetDirectoryContents(ctx context.Context, repoID string, path string, userID string) ([]entity.FileNode, error) {
	// 1. Find the repository by ID
	repo, err := uc.repo.FindByID(ctx, repoID)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve repository details: %w", err)
	}
	if repo == nil {
		return nil, apperror.NewNotFoundError("Repository", repoID, nil)
	}

	// 2. Get OAuth access token for the provider
	provider := getProviderFromURL(repo.URL())
	if provider == "" {
		return nil, apperror.NewValidationFailedError([]apperror.FieldError{
			{Field: "url", Message: "unsupported Git provider"},
		})
	}

	accessToken, err := uc.oauthProvider.GetAccessTokenForProvider(ctx, userID, provider)
	if err != nil {
		return nil, apperror.NewValidationFailedError([]apperror.FieldError{
			{Field: "oauth", Message: fmt.Sprintf("OAuth connection required for %s. Please connect your account.", provider)},
		})
	}

	// Set the OAuth token on the repository for GitManager to use
	repo.SetAccessToken(accessToken)

	// 3. Get directory contents from GitHub API (non-recursive)
	fileNodes, err := uc.gitManager.ListDirectoryContents(ctx, path, repo)
	if err != nil {
		return nil, fmt.Errorf("failed to list directory contents at path '%s': %w", path, err)
	}

	return fileNodes, nil
}

// GetFileContents implements the logic for retrieving a specific file's content.
// If commitHash is empty, retrieves the latest version. Returns content and actual commit hash.
func (uc *repositoryUseCase) GetFileContents(ctx context.Context, repoID string, filePath string, userID string, commitHash string) (string, string, error) {
	// 1. Find the repository by ID to ensure it exists
	repo, err := uc.repo.FindByID(ctx, repoID)
	if err != nil {
		return "", "", fmt.Errorf("failed to retrieve repository details: %w", err)
	}
	if repo == nil {
		return "", "", apperror.NewNotFoundError("Repository", repoID, nil)
	}

	// 2. Get OAuth access token for the provider
	provider := getProviderFromURL(repo.URL())
	if provider == "" {
		return "", "", apperror.NewValidationFailedError([]apperror.FieldError{
			{Field: "url", Message: "unsupported Git provider"},
		})
	}

	accessToken, err := uc.oauthProvider.GetAccessTokenForProvider(ctx, userID, provider)
	if err != nil {
		return "", "", apperror.NewValidationFailedError([]apperror.FieldError{
			{Field: "oauth", Message: fmt.Sprintf("OAuth connection required for %s. Please connect your account.", provider)},
		})
	}

	// Set the OAuth token on the repository for GitManager to use
	repo.SetAccessToken(accessToken)

	// 3. Read the file content at specific commit (or latest if commitHash is empty)
	contentBytes, actualCommit, err := uc.gitManager.ReadFileAtCommit(ctx, filePath, commitHash, repo)
	if err != nil {
		return "", "", fmt.Errorf("failed to read content of file '%s': %w", filePath, err)
	}

	return string(contentBytes), actualCommit, nil
}

// GetLatestCommit retrieves the latest commit information for a repository.
func (uc *repositoryUseCase) GetLatestCommit(ctx context.Context, repoID string, userID string) (*git.CommitInfo, error) {
	// 1. Find the repository by ID
	repo, err := uc.repo.FindByID(ctx, repoID)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve repository details: %w", err)
	}
	if repo == nil {
		return nil, apperror.NewNotFoundError("Repository", repoID, nil)
	}

	// 2. Try to get OAuth token if no stored access token
	if repo.AccessToken() == "" && userID != "" {
		provider := getProviderFromURL(repo.URL())
		if provider != "" {
			token, err := uc.oauthProvider.GetAccessTokenForProvider(ctx, userID, provider)
			if err == nil && token != "" {
				repo = entity.ReconstructRepository(repo.ID(), repo.Name(), repo.URL(), token, repo.CreatedAt(), repo.UpdatedAt())
			}
		}
	}

	// 3. Get the latest commit information
	commitInfo, err := uc.gitManager.GetLatestCommit(ctx, repo)
	if err != nil {
		return nil, fmt.Errorf("failed to get latest commit: %w", err)
	}

	return commitInfo, nil
}

// GetFileCommitHistory retrieves the commit history for a specific file.
func (uc *repositoryUseCase) GetFileCommitHistory(ctx context.Context, repoID string, filePath string, userID string) ([]git.CommitInfo, error) {
	// 1. Find the repository by ID
	repo, err := uc.repo.FindByID(ctx, repoID)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve repository details: %w", err)
	}
	if repo == nil {
		return nil, apperror.NewNotFoundError("Repository", repoID, nil)
	}

	// 2. Try to get OAuth token if no stored access token
	if repo.AccessToken() == "" && userID != "" {
		provider := getProviderFromURL(repo.URL())
		if provider != "" {
			token, err := uc.oauthProvider.GetAccessTokenForProvider(ctx, userID, provider)
			if err == nil && token != "" {
				repo = entity.ReconstructRepository(repo.ID(), repo.Name(), repo.URL(), token, repo.CreatedAt(), repo.UpdatedAt())
			}
		}
	}

	// 3. Get the file commit history
	commits, err := uc.gitManager.GetFileCommitHistory(ctx, filePath, repo)
	if err != nil {
		return nil, fmt.Errorf("failed to get file commit history: %w", err)
	}

	return commits, nil
}
