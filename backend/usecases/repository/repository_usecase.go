package repository

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"path"
	"path/filepath"
	"regexp"
	"strings"

	"opscore/backend/domain/model"
	"opscore/backend/domain/repository"
	"opscore/backend/infrastructure/git"

	"github.com/google/uuid"
)

// ドメイン固有のエラー定義
var (
	ErrRepositoryAlreadyExists = errors.New("repository with this URL already exists")
	ErrRepositoryNotFound      = errors.New("repository not found")
	ErrInvalidRepositoryURL    = errors.New("invalid repository URL format")
	ErrUnsupportedURLScheme    = errors.New("unsupported repository URL scheme: only https is supported")
)

// 有効なGitリポジトリURLのパターン
var (
	validGitURLPattern = regexp.MustCompile(`^https://(?:github\.com|gitlab\.com|bitbucket\.org)/[a-zA-Z0-9_.-]+/[a-zA-Z0-9_.-]+(?:\.git)?$`)
)

// RepositoryUseCase defines the interface for repository related use cases.
type RepositoryUseCase interface {
	Register(ctx context.Context, repoURL string) (*model.Repository, error) // Return created repository
	// ListFiles retrieves the file structure for a given repository ID.
	ListFiles(ctx context.Context, repoID string) ([]*model.FileNode, error) // Use model.FileNode
	// SelectFiles marks specific files within a repository as manageable.
	SelectFiles(ctx context.Context, repoID string, filePaths []string) error
	// GetSelectedMarkdown retrieves the concatenated content of selected Markdown files.
	GetSelectedMarkdown(ctx context.Context, repoID string) (string, error)
}

// repositoryUseCase implements the RepositoryUseCase interface.
type repositoryUseCase struct {
	repo       repository.Repository // Persistence for repository metadata
	gitManager git.GitManager        // For interacting with Git repositories
}

// NewRepositoryUseCase creates a new instance of repositoryUseCase.
func NewRepositoryUseCase(repo repository.Repository, gitManager git.GitManager) RepositoryUseCase {
	return &repositoryUseCase{
		repo:       repo,
		gitManager: gitManager, // Initialize GitManager
	}
}

// validateRepositoryURL validates that the URL is properly formatted and uses supported schemes
func validateRepositoryURL(repoURL string) error {
	// 1. Parse the URL
	parsedURL, err := url.ParseRequestURI(repoURL)
	if err != nil {
		return ErrInvalidRepositoryURL
	}

	// 2. Ensure the scheme is https only (more secure)
	if parsedURL.Scheme != "https" {
		return ErrUnsupportedURLScheme
	}

	// 3. Validate against whitelist pattern
	if !validGitURLPattern.MatchString(repoURL) {
		return ErrInvalidRepositoryURL
	}

	return nil
}

// Register implements the logic for registering a new repository.
func (uc *repositoryUseCase) Register(ctx context.Context, repoURL string) (*model.Repository, error) {
	// 1. Validate URL with enhanced security
	if err := validateRepositoryURL(repoURL); err != nil {
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
		return nil, ErrRepositoryAlreadyExists
	}

	// 3. Extract repository name from URL (simple approach)
	repoName := path.Base(parsedURL.Path)
	repoName = strings.TrimSuffix(repoName, ".git") // Remove .git suffix if present
	if repoName == "" || repoName == "." {
		// Fallback or error if name extraction fails
		repoName = "unknown" // Or return an error
	}

	// 4. Create new repository model
	newRepo := model.NewRepository(
		uuid.NewString(), // Generate new UUID
		repoName,
		repoURL,
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

// ListFiles implements the logic for listing files in a repository.
func (uc *repositoryUseCase) ListFiles(ctx context.Context, repoID string) ([]*model.FileNode, error) {
	// 1. Find the repository by ID
	repo, err := uc.repo.FindByID(ctx, repoID)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve repository details: %w", err)
	}
	if repo == nil {
		return nil, ErrRepositoryNotFound
	}

	// 3. Ensure repository is cloned/updated locally and get its path
	localPath, err := uc.gitManager.EnsureCloned(ctx, repo)
	if err != nil {
		return nil, fmt.Errorf("failed to ensure repository is cloned: %w", err)
	}

	// 4. List files using GitManager
	files, err := uc.gitManager.ListRepositoryFiles(ctx, localPath)
	if err != nil {
		return nil, fmt.Errorf("failed to list repository files: %w", err)
	}

	// 5. Map the output to []model.FileNode (Basic mapping)
	fileNodes := make([]*model.FileNode, 0, len(files))
	for _, f := range files {
		fileType := "file"
		fileNodes = append(fileNodes, model.NewFileNode(f, fileType))
	}

	return fileNodes, nil
}

// SelectFiles implements the logic for selecting manageable files.
func (uc *repositoryUseCase) SelectFiles(ctx context.Context, repoID string, filePaths []string) error {
	// 1. Find the repository by ID to ensure it exists
	repo, err := uc.repo.FindByID(ctx, repoID)
	if err != nil {
		return fmt.Errorf("failed to retrieve repository details: %w", err)
	}
	if repo == nil {
		return ErrRepositoryNotFound
	}

	// 3. Ensure repository is cloned/updated locally and get its path
	localPath, err := uc.gitManager.EnsureCloned(ctx, repo)
	if err != nil {
		return fmt.Errorf("failed to ensure repository is cloned: %w", err)
	}

	// 4. Validate that each file path in filePaths exists in the repository
	err = uc.gitManager.ValidateFilesExist(ctx, localPath, filePaths)
	if err != nil {
		return err
	}

	// 6. Persist the selection using the repository method
	err = uc.repo.SaveManagedFiles(ctx, repoID, filePaths)
	if err != nil {
		return fmt.Errorf("failed to save managed files selection: %w", err)
	}

	return nil
}

// GetSelectedMarkdown implements the logic for retrieving selected Markdown content.
func (uc *repositoryUseCase) GetSelectedMarkdown(ctx context.Context, repoID string) (string, error) {
	// 1. Find the repository by ID to ensure it exists
	repo, err := uc.repo.FindByID(ctx, repoID)
	if err != nil {
		return "", fmt.Errorf("failed to retrieve repository details: %w", err)
	}
	if repo == nil {
		return "", ErrRepositoryNotFound
	}

	// 3. Retrieve the list of selected file paths for this repoID
	selectedPaths, err := uc.repo.GetManagedFiles(ctx, repoID)
	if err != nil {
		return "", fmt.Errorf("failed to retrieve managed files: %w", err)
	}

	if len(selectedPaths) == 0 {
		return "", nil // No files selected, return empty string
	}

	// 5. Ensure repository is cloned/updated locally and get its path
	localPath, err := uc.gitManager.EnsureCloned(ctx, repo)
	if err != nil {
		return "", fmt.Errorf("failed to ensure repository is cloned: %w", err)
	}

	// 6 & 7. Read and concatenate content of selected Markdown files
	var concatenatedContent strings.Builder
	filesRead := 0
	for _, p := range selectedPaths {
		// a. Verify it's a Markdown file (case-insensitive check)
		ext := strings.ToLower(filepath.Ext(p))
		if ext == ".md" || ext == ".markdown" {
			// b. Read the file content from the local repository path
			contentBytes, readErr := uc.gitManager.ReadManagedFileContent(ctx, localPath, p)
			if readErr != nil {
				return "", fmt.Errorf("failed to read content of file '%s': %w", p, readErr)
			}
			// Add a separator (like ---) between files? Optional.
			if concatenatedContent.Len() > 0 {
				concatenatedContent.WriteString("\n\n---\n\n") // Markdown horizontal rule
			}
			concatenatedContent.Write(contentBytes)
			filesRead++
		}
	}

	// 8. Return the concatenated content
	return concatenatedContent.String(), nil
}
