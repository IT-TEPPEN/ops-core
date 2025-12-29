package repository

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"path"
	"regexp"
	"strings"

	apperror "opscore/backend/internal/git_repository/application/error"
	"opscore/backend/internal/git_repository/domain/entity"
	domainerror "opscore/backend/internal/git_repository/domain/error"
	"opscore/backend/internal/git_repository/domain/repository"
	"opscore/backend/internal/git_repository/infrastructure/git"

	"github.com/google/uuid"
)

// Legacy error definitions - kept for backward compatibility during migration
// These will be removed after full migration
var (
	ErrRepositoryAlreadyExists = errors.New("repository with this URL already exists")
	ErrRepositoryNotFound      = errors.New("repository not found")
	ErrInvalidRepositoryURL    = errors.New("invalid repository URL format")
	ErrUnsupportedURLScheme    = errors.New("unsupported repository URL scheme: only https is supported")
	ErrAccessTokenRequired     = errors.New("access token is required for this operation")
)

// 有効なGitリポジトリURLのパターン
var (
	validGitURLPattern = regexp.MustCompile(`^https://(?:github\.com|gitlab\.com|bitbucket\.org)/[a-zA-Z0-9_.-]+/[a-zA-Z0-9_.-]+(?:\.git)?$`)
)

// OAuthTokenProvider provides OAuth access tokens for Git providers
type OAuthTokenProvider interface {
	GetAccessTokenForProvider(ctx context.Context, userID string, providerName string) (string, error)
}

// RepositoryUseCase defines the interface for repository related use cases.
type RepositoryUseCase interface {
	Register(ctx context.Context, repoURL string, accessToken string) (entity.Repository, error) // Return created repository
	// GetRepository retrieves a single repository by its ID
	GetRepository(ctx context.Context, repoID string) (entity.Repository, error)
	// ListRepositories retrieves all registered repositories
	ListRepositories(ctx context.Context) ([]entity.Repository, error)
	// ListFiles retrieves the file structure for a given repository ID.
	ListFiles(ctx context.Context, repoID string, userID string) ([]entity.FileNode, error) // Use entity.FileNode
	// GetFileContents retrieves the content of a specific file from a repository.
	GetFileContents(ctx context.Context, repoID string, filePath string, userID string) (string, error)
	// UpdateAccessToken updates the access token for a repository.
	UpdateAccessToken(ctx context.Context, repoID string, accessToken string) error
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

// validateRepositoryURL validates that the URL is properly formatted and uses supported schemes
func validateRepositoryURL(repoURL string) error {
	// 1. Parse the URL
	parsedURL, err := url.ParseRequestURI(repoURL)
	if err != nil {
		return domainerror.NewURLValidationError("url", repoURL, "invalid repository URL format", false)
	}

	// 2. Ensure the scheme is https only (more secure)
	if parsedURL.Scheme != "https" {
		return domainerror.NewURLValidationError("url", repoURL, "unsupported repository URL scheme: only https is supported", true)
	}

	// 3. Validate against whitelist pattern
	if !validGitURLPattern.MatchString(repoURL) {
		return domainerror.NewURLValidationError("url", repoURL, "repository URL must match pattern: https://github.com|gitlab.com|bitbucket.org/owner/repo", false)
	}

	return nil
}

// Register implements the logic for registering a new repository.
func (uc *repositoryUseCase) Register(ctx context.Context, repoURL string, accessToken string) (entity.Repository, error) {
	// 1. Validate URL with enhanced security
	if err := validateRepositoryURL(repoURL); err != nil {
		// Wrap domain validation error with application context
		if errors.Is(err, domainerror.ErrInvalidEntity) {
			return nil, apperror.NewValidationFailedError([]apperror.FieldError{
				{Field: "url", Message: err.Error()},
			})
		}
		return nil, err
	}

	parsedURL, _ := url.ParseRequestURI(repoURL)

	// 2. Check if repository already exists
	existingRepo, err := uc.repo.FindByURL(ctx, repoURL)
	if err != nil {
		// Handle potential database errors (log them)
		// For now, return a generic error
		return nil, fmt.Errorf("failed to check for existing repository: %w", err)
	}
	if existingRepo != nil {
		return nil, apperror.NewConflictError("Repository", repoURL, "repository with this URL already exists", nil)
	}

	// 3. Extract repository name from URL (simple approach)
	repoName := path.Base(parsedURL.Path)
	repoName = strings.TrimSuffix(repoName, ".git") // Remove .git suffix if present
	if repoName == "" || repoName == "." {
		// Fallback or error if name extraction fails
		repoName = "unknown" // Or return an error
	}

	// 4. Create new repository model
	newRepo := entity.NewRepository(
		uuid.NewString(), // Generate new UUID
		repoName,
		repoURL,
		accessToken,
	)

	// 5. Persist the new repository
	err = uc.repo.Save(ctx, newRepo)
	if err != nil {
		// Handle potential database errors (log them)
		return nil, fmt.Errorf("failed to save repository: %w", err)
	}

	// 6. Return the newly created repository
	return newRepo, nil
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

// GetFileContents implements the logic for retrieving a specific file's content.
func (uc *repositoryUseCase) GetFileContents(ctx context.Context, repoID string, filePath string, userID string) (string, error) {
	// 1. Find the repository by ID to ensure it exists
	repo, err := uc.repo.FindByID(ctx, repoID)
	if err != nil {
		return "", fmt.Errorf("failed to retrieve repository details: %w", err)
	}
	if repo == nil {
		return "", apperror.NewNotFoundError("Repository", repoID, nil)
	}

	// 2. Get OAuth access token for the provider
	provider := getProviderFromURL(repo.URL())
	if provider == "" {
		return "", apperror.NewValidationFailedError([]apperror.FieldError{
			{Field: "url", Message: "unsupported Git provider"},
		})
	}

	accessToken, err := uc.oauthProvider.GetAccessTokenForProvider(ctx, userID, provider)
	if err != nil {
		return "", apperror.NewValidationFailedError([]apperror.FieldError{
			{Field: "oauth", Message: fmt.Sprintf("OAuth connection required for %s. Please connect your account.", provider)},
		})
	}

	// Set the OAuth token on the repository for GitManager to use
	repo.SetAccessToken(accessToken)

	// 3. Read the file content (on-demand fetching with caching)
	contentBytes, err := uc.gitManager.ReadManagedFileContent(ctx, "", filePath, repo)
	if err != nil {
		return "", fmt.Errorf("failed to read content of file '%s': %w", filePath, err)
	}

	return string(contentBytes), nil
}

// UpdateAccessToken updates the access token for a repository.
func (uc *repositoryUseCase) UpdateAccessToken(ctx context.Context, repoID string, accessToken string) error {
	// 1. Find the repository by ID to ensure it exists
	repo, err := uc.repo.FindByID(ctx, repoID)
	if err != nil {
		return fmt.Errorf("failed to retrieve repository details: %w", err)
	}
	if repo == nil {
		return apperror.NewNotFoundError("Repository", repoID, nil)
	}

	// 2. Update the access token in the repository
	err = uc.repo.UpdateAccessToken(ctx, repoID, accessToken)
	if err != nil {
		return fmt.Errorf("failed to update repository access token: %w", err)
	}

	return nil
}
