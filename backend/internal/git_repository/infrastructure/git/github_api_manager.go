package git

import (
	"context"
	"fmt"
	"net/http"
	"opscore/backend/internal/git_repository/domain/entity"
	"os"
	"path/filepath"
	"strings"

	"github.com/google/go-github/v60/github"
	"golang.org/x/oauth2"
)

// GitHubAPIManager implements the GitManager interface using the GitHub API.
type githubApiManager struct {
	baseClonePath string                    // Base directory where repositories will be stored locally
	clients       map[string]*github.Client // Cache of GitHub clients by token
}

// NewGithubApiManager creates a new githubApiManager.
// baseClonePath is the directory where repositories will be stored locally.
func NewGithubApiManager(baseClonePath string) (GitManager, error) {
	// Ensure the base clone path exists
	err := os.MkdirAll(baseClonePath, 0755)
	if err != nil {
		return nil, fmt.Errorf("failed to create base clone directory '%s': %w", baseClonePath, err)
	}
	return &githubApiManager{
		baseClonePath: baseClonePath,
		clients:       make(map[string]*github.Client),
	}, nil
}

// getLocalPath determines the local directory path for a given repository.
func (g *githubApiManager) getLocalPath(repo entity.Repository) string {
	// Use a sanitized version of the repo ID as the directory name
	return filepath.Join(g.baseClonePath, repo.ID())
}

// getGitHubClient returns a GitHub API client, authenticated if a token is provided.
func (g *githubApiManager) getGitHubClient(accessToken string) *github.Client {
	// Check if we have a cached client for this token
	if client, ok := g.clients[accessToken]; ok {
		return client
	}

	// Create a new client
	var httpClient *http.Client
	if accessToken != "" {
		// Create an authenticated client if we have a token
		ts := oauth2.StaticTokenSource(
			&oauth2.Token{AccessToken: accessToken},
		)
		httpClient = oauth2.NewClient(context.Background(), ts)
	}

	client := github.NewClient(httpClient)

	// Cache the client for future use
	if accessToken != "" {
		g.clients[accessToken] = client
	}

	return client
}

// parseGitHubURL extracts owner and repo name from a GitHub URL.
func parseGitHubURL(repoURL string) (string, string, error) {
	// Handle URLs like https://github.com/owner/repo.git or https://github.com/owner/repo
	urlParts := strings.Split(repoURL, "/")
	if len(urlParts) < 5 {
		return "", "", fmt.Errorf("invalid GitHub URL format: %s", repoURL)
	}

	owner := urlParts[len(urlParts)-2]
	repo := strings.TrimSuffix(urlParts[len(urlParts)-1], ".git")

	return owner, repo, nil
}

// EnsureCloned ensures the repository is available locally, either by cloning it or updating an existing clone.
func (g *githubApiManager) EnsureCloned(ctx context.Context, repo entity.Repository) (string, error) {
	localPath := g.getLocalPath(repo)

	// Extract owner and repo name from URL
	owner, repoName, err := parseGitHubURL(repo.URL())
	if err != nil {
		return "", err
	}

	// Get GitHub client
	client := g.getGitHubClient(repo.AccessToken())

	// Check if the directory exists
	if _, err := os.Stat(localPath); os.IsNotExist(err) {
		// Directory does not exist, clone the repository
		fmt.Printf("Cloning repository %s to %s\n", repo.URL(), localPath)

		// Create the directory
		if err := os.MkdirAll(localPath, 0755); err != nil {
			return "", fmt.Errorf("failed to create directory for repository: %w", err)
		}

		// Download the default branch content
		err = g.downloadRepository(ctx, client, owner, repoName, "", localPath)
		if err != nil {
			return "", fmt.Errorf("failed to clone repository %s: %w", repo.URL(), err)
		}
	} else if err == nil {
		// Directory exists, update the repository
		fmt.Printf("Updating repository %s in %s\n", repo.URL(), localPath)

		// Clear the directory and re-download
		files, err := os.ReadDir(localPath)
		if err != nil {
			return "", fmt.Errorf("failed to read repository directory: %w", err)
		}

		// Only remove files and directories inside the local path, not the directory itself
		for _, file := range files {
			path := filepath.Join(localPath, file.Name())
			if err := os.RemoveAll(path); err != nil {
				return "", fmt.Errorf("failed to remove old files: %w", err)
			}
		}

		// Download the repository content again
		err = g.downloadRepository(ctx, client, owner, repoName, "", localPath)
		if err != nil {
			return "", fmt.Errorf("failed to update repository %s: %w", repo.URL(), err)
		}
	} else {
		// Other error checking directory
		return "", fmt.Errorf("failed to check repository directory %s: %w", localPath, err)
	}

	return localPath, nil
}

// downloadRepository recursively downloads repository content from GitHub.
func (g *githubApiManager) downloadRepository(ctx context.Context, client *github.Client, owner string, repo string, path string, localPath string) error {
	// List files and directories at the current path
	_, contents, _, err := client.Repositories.GetContents(ctx, owner, repo, path, &github.RepositoryContentGetOptions{})
	if err != nil {
		return fmt.Errorf("failed to get repository contents: %w", err)
	}

	for _, content := range contents {
		localFilePath := filepath.Join(localPath, *content.Path)

		switch *content.Type {
		case "file":
			// Ensure directory exists
			dir := filepath.Dir(localFilePath)
			if err := os.MkdirAll(dir, 0755); err != nil {
				return fmt.Errorf("failed to create directory %s: %w", dir, err)
			}

			fileContentStr, err := content.GetContent()
			if err != nil {
				fileContent, _, _, fetchErr := client.Repositories.GetContents(ctx, owner, repo, *content.Path, &github.RepositoryContentGetOptions{})
				if fetchErr != nil {
					return fmt.Errorf("failed to get file content for %s: %w", *content.Path, fetchErr)
				}

				fileContentStr, err = fileContent.GetContent()
				if err != nil {
					return fmt.Errorf("failed to decode content for %s: %w", *content.Path, err)
				}
			}

			if err := os.WriteFile(localFilePath, []byte(fileContentStr), 0644); err != nil {
				return fmt.Errorf("failed to write file %s: %w", localFilePath, err)
			}

		case "dir":
			// Create directory and recursively download its contents
			if err := os.MkdirAll(localFilePath, 0755); err != nil {
				return fmt.Errorf("failed to create directory %s: %w", localFilePath, err)
			}

			if err := g.downloadRepository(ctx, client, owner, repo, *content.Path, localPath); err != nil {
				return err
			}
		}
	}

	return nil
}

// ListRepositoryFiles lists all files in the repository.
func (g *githubApiManager) ListRepositoryFiles(ctx context.Context, localPath string, repo entity.Repository) ([]string, error) {
	// Get the file list from the local clone
	var files []string
	err := filepath.Walk(localPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		// Skip the root directory and hidden files/dirs
		if path == localPath || filepath.Base(path)[0] == '.' {
			return nil
		}
		// Only include files, not directories
		if !info.IsDir() {
			// Get path relative to repository root
			relPath, err := filepath.Rel(localPath, path)
			if err != nil {
				return err
			}
			files = append(files, filepath.ToSlash(relPath))
		}
		return nil
	})

	if err != nil {
		// If walking the local directory fails, try to get files directly from the API
		fmt.Printf("Failed to walk local directory, falling back to API: %v\n", err)

		// Extract owner and repo name from URL
		owner, repoName, err := parseGitHubURL(repo.URL())
		if err != nil {
			return nil, err
		}

		// Get GitHub client
		client := g.getGitHubClient(repo.AccessToken())

		files, err = g.listFilesFromAPI(ctx, client, owner, repoName, "")
		if err != nil {
			return nil, fmt.Errorf("failed to list files in repository: %w", err)
		}
	}

	return files, nil
}

// listFilesFromAPI recursively lists files from the GitHub API.
func (g *githubApiManager) listFilesFromAPI(ctx context.Context, client *github.Client, owner string, repo string, path string) ([]string, error) {
	var files []string

	// List files and directories at the current path
	_, contents, _, err := client.Repositories.GetContents(ctx, owner, repo, path, &github.RepositoryContentGetOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to get repository contents: %w", err)
	}

	for _, content := range contents {
		if *content.Type == "file" {
			files = append(files, *content.Path)
		} else if *content.Type == "dir" {
			// Recursively get files from subdirectory
			subFiles, err := g.listFilesFromAPI(ctx, client, owner, repo, *content.Path)
			if err != nil {
				return nil, err
			}
			files = append(files, subFiles...)
		}
	}

	return files, nil
}

// ValidateFilesExist checks if files exist in the repository.
func (g *githubApiManager) ValidateFilesExist(ctx context.Context, localPath string, filePaths []string, repo entity.Repository) error {
	if len(filePaths) == 0 {
		return nil // Nothing to validate
	}

	// Check if files exist in the local clone
	for _, filePath := range filePaths {
		fullPath := filepath.Join(localPath, filePath)
		if _, err := os.Stat(fullPath); os.IsNotExist(err) {
			return fmt.Errorf("file does not exist in repository: %s", filePath)
		} else if err != nil {
			return fmt.Errorf("error checking if file exists: %s: %w", filePath, err)
		}
	}

	return nil
}

// ReadManagedFileContent reads the content of a repository file.
func (g *githubApiManager) ReadManagedFileContent(ctx context.Context, localPath string, filePath string, repo entity.Repository) ([]byte, error) {
	// Security check: ensure the filePath doesn't contain path traversal sequences
	if strings.Contains(filePath, "..") {
		return nil, fmt.Errorf("invalid file path containing path traversal sequences: %s", filePath)
	}

	// Join the local repository path with the requested file path
	fullPath := filepath.Join(localPath, filePath)

	// Ensure the resulting path is still within the repository directory
	absLocalPath, _ := filepath.Abs(localPath)
	absFullPath, _ := filepath.Abs(fullPath)
	if !strings.HasPrefix(absFullPath, absLocalPath) {
		return nil, fmt.Errorf("invalid file path: attempt to access file outside repository directory")
	}

	// Read the file content from the local filesystem
	content, err := os.ReadFile(fullPath)
	if err != nil {
		// If file not found or other error reading locally, try to get it from the API
		if os.IsNotExist(err) {
			// Extract owner and repo name from URL
			owner, repoName, err := parseGitHubURL(repo.URL())
			if err != nil {
				return nil, err
			}

			// Get file content via GitHub API
			client := g.getGitHubClient(repo.AccessToken())
			fileContent, _, _, err := client.Repositories.GetContents(ctx, owner, repoName, filePath, &github.RepositoryContentGetOptions{})
			if err != nil {
				return nil, fmt.Errorf("failed to get file content from API: %w", err)
			}

			content, err := fileContent.GetContent()
			if err != nil {
				return nil, fmt.Errorf("failed to decode content from API: %w", err)
			}

			return []byte(content), nil
		}
		return nil, fmt.Errorf("failed to read file %s: %w", fullPath, err)
	}

	return content, nil
}
