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
		// Directory exists, use cached clone
		fmt.Printf("Using cached repository %s in %s\n", repo.URL(), localPath)
		// TODO: Implement update logic with SHA checking or periodic refresh if needed
	} else {
		// Other error checking directory
		return "", fmt.Errorf("failed to check repository directory %s: %w", localPath, err)
	}

	return localPath, nil
}

// downloadRepository recursively downloads repository content from GitHub.
func (g *githubApiManager) downloadRepository(ctx context.Context, client *github.Client, owner string, repo string, path string, localPath string) error {
	// List files and directories at the current path
	fileContent, directoryContents, _, err := client.Repositories.GetContents(ctx, owner, repo, path, &github.RepositoryContentGetOptions{})
	if err != nil {
		return fmt.Errorf("failed to get repository contents: %w", err)
	}

	// If it's a single file, handle it directly
	if fileContent != nil {
		localFilePath := filepath.Join(localPath, path)
		dir := filepath.Dir(localFilePath)
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("failed to create directory %s: %w", dir, err)
		}

		content, err := fileContent.GetContent()
		if err != nil {
			return fmt.Errorf("failed to decode content for %s: %w", path, err)
		}

		if err := os.WriteFile(localFilePath, []byte(content), 0644); err != nil {
			return fmt.Errorf("failed to write file %s: %w", localFilePath, err)
		}
		return nil
	}

	// If it's a directory, process all contents
	for _, content := range directoryContents {
		localFilePath := filepath.Join(localPath, *content.Path)

		switch *content.Type {
		case "file":
			// For files in directory listing, we need to fetch content separately
			// because directory listing doesn't include file contents
			fileContentItem, _, _, err := client.Repositories.GetContents(ctx, owner, repo, *content.Path, &github.RepositoryContentGetOptions{})
			if err != nil {
				return fmt.Errorf("failed to get file content for %s: %w", *content.Path, err)
			}

			// Ensure directory exists
			dir := filepath.Dir(localFilePath)
			if err := os.MkdirAll(dir, 0755); err != nil {
				return fmt.Errorf("failed to create directory %s: %w", dir, err)
			}

			// Get file content
			fileContentStr, err := fileContentItem.GetContent()
			if err != nil {
				return fmt.Errorf("failed to decode content for %s: %w", *content.Path, err)
			}

			fmt.Printf("Writing file %s (size: %d bytes)\n", localFilePath, len(fileContentStr))
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

// ListRepositoryFiles lists all files in the repository directly from GitHub API.
func (g *githubApiManager) ListRepositoryFiles(ctx context.Context, localPath string, repo entity.Repository) ([]string, error) {
	// Always fetch file list from GitHub API to get the latest repository structure
	fmt.Printf("Fetching file list from GitHub API for repository: %s\n", repo.URL())

	// Extract owner and repo name from URL
	owner, repoName, err := parseGitHubURL(repo.URL())
	if err != nil {
		return nil, err
	}

	// Get GitHub client
	client := g.getGitHubClient(repo.AccessToken())

	// Get files from API
	files, err := g.listFilesFromAPI(ctx, client, owner, repoName, "")
	if err != nil {
		return nil, fmt.Errorf("failed to list files in repository: %w", err)
	}

	return files, nil
}

// ListDirectoryContents lists files and directories at a specific path (non-recursive).
func (g *githubApiManager) ListDirectoryContents(ctx context.Context, path string, repo entity.Repository) ([]entity.FileNode, error) {
	fmt.Printf("Fetching directory contents from GitHub API for repository: %s, path: %s\n", repo.URL(), path)

	// Extract owner and repo name from URL
	owner, repoName, err := parseGitHubURL(repo.URL())
	if err != nil {
		return nil, err
	}

	// Get GitHub client
	client := g.getGitHubClient(repo.AccessToken())

	// Get contents at the specified path
	_, contents, _, err := client.Repositories.GetContents(ctx, owner, repoName, path, &github.RepositoryContentGetOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to get repository contents at path '%s': %w", path, err)
	}

	// Map contents to FileNode entities
	fileNodes := make([]entity.FileNode, 0, len(contents))
	for _, content := range contents {
		if content.Name == nil || content.Type == nil {
			continue
		}

		nodeType := "file"
		if *content.Type == "dir" {
			nodeType = "dir"
		}

		fileNodes = append(fileNodes, entity.NewFileNode(*content.Name, nodeType))
	}

	return fileNodes, nil
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

// ReadManagedFileContent reads the content of a repository file with on-demand fetching.
func (g *githubApiManager) ReadManagedFileContent(ctx context.Context, localPath string, filePath string, repo entity.Repository) ([]byte, error) {
	// Security check: ensure the filePath doesn't contain path traversal sequences
	if strings.Contains(filePath, "..") {
		return nil, fmt.Errorf("invalid file path containing path traversal sequences: %s", filePath)
	}

	// Determine the local cache path
	if localPath == "" {
		localPath = g.getLocalPath(repo)
	}

	// Join the local repository path with the requested file path
	fullPath := filepath.Join(localPath, filePath)

	// Ensure the resulting path is still within the repository directory
	absLocalPath, _ := filepath.Abs(localPath)
	absFullPath, _ := filepath.Abs(fullPath)
	if !strings.HasPrefix(absFullPath, absLocalPath) {
		return nil, fmt.Errorf("invalid file path: attempt to access file outside repository directory")
	}

	// Try to read from local cache first
	content, err := os.ReadFile(fullPath)
	if err == nil {
		fmt.Printf("Reading file from cache: %s\n", filePath)
		return content, nil
	}

	// If file not found locally, fetch from GitHub API
	if os.IsNotExist(err) {
		fmt.Printf("Fetching file from GitHub API: %s\n", filePath)

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

		contentStr, err := fileContent.GetContent()
		if err != nil {
			return nil, fmt.Errorf("failed to decode content from API: %w", err)
		}

		content = []byte(contentStr)

		// Cache the file locally for future requests
		dir := filepath.Dir(fullPath)
		if err := os.MkdirAll(dir, 0755); err != nil {
			fmt.Printf("Warning: failed to create cache directory %s: %v\n", dir, err)
		} else {
			if err := os.WriteFile(fullPath, content, 0644); err != nil {
				fmt.Printf("Warning: failed to cache file %s: %v\n", fullPath, err)
			} else {
				fmt.Printf("Cached file: %s\n", fullPath)
			}
		}

		return content, nil
	}

	// Other error reading file
	return nil, fmt.Errorf("failed to read file %s: %w", fullPath, err)
}

// GetLatestCommit retrieves the latest commit information for the repository.
func (g *githubApiManager) GetLatestCommit(ctx context.Context, repo entity.Repository) (*CommitInfo, error) {
	// Extract owner and repo name from URL
	owner, repoName, err := parseGitHubURL(repo.URL())
	if err != nil {
		return nil, err
	}

	// Get GitHub client
	client := g.getGitHubClient(repo.AccessToken())

	// Get the default branch
	repoInfo, _, err := client.Repositories.Get(ctx, owner, repoName)
	if err != nil {
		return nil, fmt.Errorf("failed to get repository info: %w", err)
	}

	defaultBranch := repoInfo.GetDefaultBranch()

	// Get the latest commit on the default branch
	commits, _, err := client.Repositories.ListCommits(ctx, owner, repoName, &github.CommitsListOptions{
		SHA: defaultBranch,
		ListOptions: github.ListOptions{
			Page:    1,
			PerPage: 1,
		},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get latest commit: %w", err)
	}

	if len(commits) == 0 {
		return nil, fmt.Errorf("no commits found in repository")
	}

	commit := commits[0]
	return &CommitInfo{
		Hash:        commit.GetSHA(),
		Message:     commit.GetCommit().GetMessage(),
		Author:      commit.GetCommit().GetAuthor().GetName(),
		AuthorEmail: commit.GetCommit().GetAuthor().GetEmail(),
		Date:        commit.GetCommit().GetAuthor().GetDate().Time,
	}, nil
}

// GetFileCommitHistory retrieves the commit history for a specific file.
func (g *githubApiManager) GetFileCommitHistory(ctx context.Context, filePath string, repo entity.Repository) ([]CommitInfo, error) {
	// Extract owner and repo name from URL
	owner, repoName, err := parseGitHubURL(repo.URL())
	if err != nil {
		return nil, err
	}

	// Get GitHub client
	client := g.getGitHubClient(repo.AccessToken())

	// Get commits for the specific file
	commits, _, err := client.Repositories.ListCommits(ctx, owner, repoName, &github.CommitsListOptions{
		Path: filePath,
		ListOptions: github.ListOptions{
			Page:    1,
			PerPage: 50, // Get up to 50 most recent commits for the file
		},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get file commit history: %w", err)
	}

	// Convert to CommitInfo slice
	result := make([]CommitInfo, 0, len(commits))
	for _, commit := range commits {
		result = append(result, CommitInfo{
			Hash:        commit.GetSHA(),
			Message:     commit.GetCommit().GetMessage(),
			Author:      commit.GetCommit().GetAuthor().GetName(),
			AuthorEmail: commit.GetCommit().GetAuthor().GetEmail(),
			Date:        commit.GetCommit().GetAuthor().GetDate().Time,
		})
	}

	return result, nil
}
