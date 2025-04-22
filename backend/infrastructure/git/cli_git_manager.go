package git

import (
	"bytes"
	"context"
	"fmt"
	"opscore/backend/domain/model"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// cliGitManager implements the GitManager interface using Git CLI commands.
type cliGitManager struct {
	baseClonePath string // Base directory where repositories will be cloned
}

// NewCliGitManager creates a new cliGitManager.
// baseClonePath is the directory where repositories will be stored locally.
func NewCliGitManager(baseClonePath string) (GitManager, error) {
	// Ensure the base clone path exists
	err := os.MkdirAll(baseClonePath, 0755)
	if err != nil {
		return nil, fmt.Errorf("failed to create base clone directory '%s': %w", baseClonePath, err)
	}
	return &cliGitManager{baseClonePath: baseClonePath},
		nil
}

// getLocalPath determines the local directory path for a given repository.
func (g *cliGitManager) getLocalPath(repo *model.Repository) string {
	// Use a sanitized version of the repo ID or name as the directory name
	// Using ID is generally safer to avoid collisions and special characters.
	return filepath.Join(g.baseClonePath, repo.ID())
}

// runGitCommand executes a git command in the specified directory.
func (g *cliGitManager) runGitCommand(ctx context.Context, dir string, args ...string) ([]byte, error) {
	if len(args) == 0 {
		return nil, fmt.Errorf("no git command specified")
	}

	// 許可されたgitコマンド（サブコマンド）のリスト
	allowedCommands := map[string]bool{
		"clone":    true,
		"fetch":    true,
		"reset":    true,
		"ls-tree":  true,
		"ls-files": true,
		// 必要に応じて他の安全なgitコマンドを追加
	}

	// コマンド引数の検証
	if !allowedCommands[args[0]] {
		return nil, fmt.Errorf("git command not allowed: %s", args[0])
	}

	// GitリポジトリのURLを含むコマンドの場合の追加検証（clone時など）
	if args[0] == "clone" && len(args) > 1 {
		// URLはHTTPSのみを許可することを明示的に検証
		if !strings.HasPrefix(args[1], "https://") {
			return nil, fmt.Errorf("only https URLs are allowed for git operations")
		}
	}

	cmd := exec.CommandContext(ctx, "git", args...)
	cmd.Dir = dir
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	// 追加のセキュリティ制限：環境変数を制限する
	cmd.Env = []string{
		"PATH=" + os.Getenv("PATH"), // 最小限のPATH環境変数のみ設定
		"HOME=" + os.Getenv("HOME"), // Gitはホームディレクトリを必要とすることがある
		"GIT_TERMINAL_PROMPT=0",     // インタラクティブプロンプトを無効化
	}

	err := cmd.Run()
	if err != nil {
		return nil, fmt.Errorf("git command failed: %v\nArgs: %v\nStderr: %s\nError: %w", args[0], args, stderr.String(), err)
	}
	return stdout.Bytes(), nil
}

// EnsureCloned clones or updates the repository.
func (g *cliGitManager) EnsureCloned(ctx context.Context, repo *model.Repository) (string, error) {
	localPath := g.getLocalPath(repo)

	// Check if the directory exists
	if _, err := os.Stat(localPath); os.IsNotExist(err) {
		// Directory does not exist, clone the repository
		fmt.Printf("Cloning repository %s to %s\n", repo.URL(), localPath)
		_, err := g.runGitCommand(ctx, g.baseClonePath, "clone", repo.URL(), localPath)
		if err != nil {
			return "", fmt.Errorf("failed to clone repository %s: %w", repo.URL(), err)
		}
	} else if err == nil {
		// Directory exists, update the repository (fetch + reset or pull)
		fmt.Printf("Updating repository %s in %s\n", repo.URL(), localPath)
		// Using fetch + reset --hard to ensure clean state, adjust if needed
		_, err := g.runGitCommand(ctx, localPath, "fetch", "origin")
		if err != nil {
			return "", fmt.Errorf("failed to fetch repository %s: %w", repo.URL(), err)
		}
		// Determine the default branch (e.g., main or master) or use HEAD
		// For simplicity, using origin/HEAD which usually points to the default branch
		_, err = g.runGitCommand(ctx, localPath, "reset", "--hard", "origin/HEAD") // Adjust branch if necessary
		if err != nil {
			// Attempt pull as a fallback? Or just report error.
			return "", fmt.Errorf("failed to reset repository %s: %w", repo.URL(), err)
		}
	} else {
		// Other error checking directory (permissions?)
		return "", fmt.Errorf("failed to check repository directory %s: %w", localPath, err)
	}

	return localPath, nil
}

// ListRepositoryFiles lists all files tracked by git.
func (g *cliGitManager) ListRepositoryFiles(ctx context.Context, localPath string) ([]string, error) {
	output, err := g.runGitCommand(ctx, localPath, "ls-tree", "-r", "--name-only", "HEAD")
	if err != nil {
		return nil, fmt.Errorf("failed to list files in %s: %w", localPath, err)
	}

	files := strings.Split(strings.TrimSpace(string(output)), "\n")
	// Filter out empty strings if any
	result := make([]string, 0, len(files))
	for _, f := range files {
		if f != "" {
			result = append(result, f)
		}
	}
	return result, nil
}

// ValidateFilesExist checks if files exist in the git repository index.
func (g *cliGitManager) ValidateFilesExist(ctx context.Context, localPath string, filePaths []string) error {
	if len(filePaths) == 0 {
		return nil // Nothing to validate
	}
	args := append([]string{"ls-files", "--error-unmatch", "--"}, filePaths...)
	_, err := g.runGitCommand(ctx, localPath, args...)
	if err != nil {
		// Error indicates one or more files were not found
		return fmt.Errorf("one or more specified files do not exist in the repository at %s: %w", localPath, err)
	}
	return nil // All files exist
}

// ReadManagedFileContent reads the content of a specific file.
func (g *cliGitManager) ReadManagedFileContent(ctx context.Context, localPath string, filePath string) ([]byte, error) {
	// セキュリティ強化: filePath内の危険な文字列をチェック
	if strings.Contains(filePath, "..") || strings.Contains(filePath, "~") {
		return nil, fmt.Errorf("invalid file path containing potentially dangerous sequences: %s", filePath)
	}

	// GitリポジトリのファイルのみにアクセスするためにLSしたファイルリストと検証する
	// これにより、Gitトラッキング対象のファイル以外へのアクセスを防ぐ
	files, err := g.ListRepositoryFiles(ctx, localPath)
	if err != nil {
		return nil, fmt.Errorf("failed to list repository files for validation: %w", err)
	}

	// ファイルがGitリポジトリ内の有効なファイルかチェック
	fileExists := false
	for _, repoFile := range files {
		if repoFile == filePath {
			fileExists = true
			break
		}
	}

	if !fileExists {
		return nil, fmt.Errorf("requested file '%s' is not tracked in the repository", filePath)
	}

	// Construct the absolute path to the file within the local clone
	absFilePath := filepath.Join(localPath, filePath)

	// 強化されたパス検証: 絶対パスがローカルリポジトリパス配下にあるか確認
	relPath, err := filepath.Rel(localPath, absFilePath)
	if err != nil || strings.HasPrefix(relPath, "..") || filepath.IsAbs(relPath) {
		return nil, fmt.Errorf("invalid file path attempting to access files outside the repository: %s", filePath)
	}

	content, err := os.ReadFile(absFilePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read file %s: %w", absFilePath, err)
	}
	return content, nil
}
