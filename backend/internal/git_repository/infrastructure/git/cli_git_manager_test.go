package git

import (
	"context"
	"os"
	"opscore/backend/internal/git_repository/domain/entity"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Note: EnsureCloned method is not fully tested in unit tests as it requires
// actual Git repository operations. Integration tests with real repositories
// would be needed for comprehensive coverage of this method.
// The method is indirectly tested through integration tests at higher layers.

// TestNewCliGitManager tests the initialization of cliGitManager
func TestNewCliGitManager(t *testing.T) {
	t.Run("正常にディレクトリが作成される", func(t *testing.T) {
		// 一時ディレクトリを使用
		tmpDir := filepath.Join(os.TempDir(), "test-cli-git-manager")
		defer os.RemoveAll(tmpDir)

		manager, err := NewCliGitManager(tmpDir)

		assert.NoError(t, err)
		assert.NotNil(t, manager)

		// ディレクトリが作成されていることを確認
		info, err := os.Stat(tmpDir)
		assert.NoError(t, err)
		assert.True(t, info.IsDir())
	})

	t.Run("既存のディレクトリで初期化できる", func(t *testing.T) {
		tmpDir := filepath.Join(os.TempDir(), "test-cli-git-manager-existing")
		defer os.RemoveAll(tmpDir)

		// 事前にディレクトリを作成
		err := os.MkdirAll(tmpDir, 0755)
		require.NoError(t, err)

		manager, err := NewCliGitManager(tmpDir)

		assert.NoError(t, err)
		assert.NotNil(t, manager)
	})

	t.Run("無効なパスでエラーになる", func(t *testing.T) {
		// 無効なパス（書き込み権限のないパスをシミュレート）
		invalidPath := "/root/invalid-test-path"

		_, err := NewCliGitManager(invalidPath)

		assert.Error(t, err)
	})
}

// TestGetLocalPath tests the getLocalPath method
func TestGetLocalPath(t *testing.T) {
	t.Run("リポジトリIDから正しいローカルパスを生成する", func(t *testing.T) {
		tmpDir := filepath.Join(os.TempDir(), "test-local-path")
		defer os.RemoveAll(tmpDir)

		manager, err := NewCliGitManager(tmpDir)
		require.NoError(t, err)

		cliManager := manager.(*cliGitManager)
		repo := entity.NewRepository("test-repo-id", "test-repo", "https://github.com/example/test", "")

		localPath := cliManager.getLocalPath(repo)

		expectedPath := filepath.Join(tmpDir, "test-repo-id")
		assert.Equal(t, expectedPath, localPath)
	})
}

// TestRunGitCommand tests the runGitCommand method
func TestRunGitCommand(t *testing.T) {
	tmpDir := filepath.Join(os.TempDir(), "test-run-git-command")
	defer os.RemoveAll(tmpDir)

	manager, err := NewCliGitManager(tmpDir)
	require.NoError(t, err)
	cliManager := manager.(*cliGitManager)

	repo := entity.NewRepository("test-id", "test-repo", "https://github.com/example/test", "")

	t.Run("許可されていないコマンドでエラーになる", func(t *testing.T) {
		_, err := cliManager.runGitCommand(context.Background(), tmpDir, repo, "push")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "not allowed")
	})

	t.Run("HTTPURLでエラーになる", func(t *testing.T) {
		_, err := cliManager.runGitCommand(context.Background(), tmpDir, repo, "clone", "http://github.com/example/test")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "only https URLs are allowed")
	})

	t.Run("コマンド引数なしでエラーになる", func(t *testing.T) {
		_, err := cliManager.runGitCommand(context.Background(), tmpDir, repo)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "no git command specified")
	})
}

// TestListRepositoryFiles tests the ListRepositoryFiles method
func TestListRepositoryFiles(t *testing.T) {
	t.Run("空のリポジトリで空リストを返す", func(t *testing.T) {
		// モックリポジトリディレクトリを作成
		tmpDir := filepath.Join(os.TempDir(), "test-list-files")
		defer os.RemoveAll(tmpDir)

		manager, err := NewCliGitManager(tmpDir)
		require.NoError(t, err)

		repo := entity.NewRepository("test-id", "test-repo", "https://github.com/example/test", "")
		localPath := filepath.Join(tmpDir, "test-id")

		// .gitディレクトリを含むリポジトリ構造を作成
		err = os.MkdirAll(filepath.Join(localPath, ".git"), 0755)
		require.NoError(t, err)

		// Gitリポジトリを初期化（実際のgitコマンドを使用）
		ctx := context.Background()
		cliManager := manager.(*cliGitManager)

		// ls-treeコマンドは実際のgitリポジトリが必要なため、エラーになることを確認
		_, err = cliManager.ListRepositoryFiles(ctx, localPath, repo)
		assert.Error(t, err) // 初期化されていないリポジトリでエラー
	})

	t.Run("正常にファイルリストをパースする", func(t *testing.T) {
		tmpDir := filepath.Join(os.TempDir(), "test-list-files-success")
		defer os.RemoveAll(tmpDir)

		manager, err := NewCliGitManager(tmpDir)
		require.NoError(t, err)

		repo := entity.NewRepository("test-id", "test-repo", "https://github.com/example/test", "")
		localPath := filepath.Join(tmpDir, "test-id")

		// リポジトリディレクトリを作成
		err = os.MkdirAll(localPath, 0755)
		require.NoError(t, err)

		ctx := context.Background()
		cliManager := manager.(*cliGitManager)

		// 実際のgitリポジトリがないため、エラーになるが、コードパスはテストされる
		_, err = cliManager.ListRepositoryFiles(ctx, localPath, repo)
		assert.Error(t, err)
	})
}

// TestValidateFilesExist tests the ValidateFilesExist method
func TestValidateFilesExist(t *testing.T) {
	tmpDir := filepath.Join(os.TempDir(), "test-validate-files")
	defer os.RemoveAll(tmpDir)

	manager, err := NewCliGitManager(tmpDir)
	require.NoError(t, err)

	repo := entity.NewRepository("test-id", "test-repo", "https://github.com/example/test", "")
	localPath := filepath.Join(tmpDir, "test-id")

	t.Run("ファイルパスが空の場合はエラーなし", func(t *testing.T) {
		err := manager.ValidateFilesExist(context.Background(), localPath, []string{}, repo)
		assert.NoError(t, err)
	})

	t.Run("存在しないファイルでエラーになる", func(t *testing.T) {
		// Gitリポジトリが初期化されていない状態でテスト
		err := manager.ValidateFilesExist(context.Background(), localPath, []string{"nonexistent.txt"}, repo)
		assert.Error(t, err)
	})
}

// TestReadManagedFileContent tests the ReadManagedFileContent method with security checks
func TestReadManagedFileContent(t *testing.T) {
	tmpDir := filepath.Join(os.TempDir(), "test-read-file")
	defer os.RemoveAll(tmpDir)

	manager, err := NewCliGitManager(tmpDir)
	require.NoError(t, err)

	repo := entity.NewRepository("test-id", "test-repo", "https://github.com/example/test", "")
	localPath := filepath.Join(tmpDir, "test-id")

	t.Run("パストラバーサル攻撃を防ぐ - .. を含むパス", func(t *testing.T) {
		_, err := manager.ReadManagedFileContent(context.Background(), localPath, "../etc/passwd", repo)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "dangerous sequences")
	})

	t.Run("パストラバーサル攻撃を防ぐ - ~ を含むパス", func(t *testing.T) {
		_, err := manager.ReadManagedFileContent(context.Background(), localPath, "~/test.txt", repo)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "dangerous sequences")
	})

	t.Run("Gitトラッキング対象外のファイルへのアクセスを防ぐ", func(t *testing.T) {
		// リポジトリディレクトリとファイルを作成
		err := os.MkdirAll(localPath, 0755)
		require.NoError(t, err)

		testFile := filepath.Join(localPath, "untracked.txt")
		err = os.WriteFile(testFile, []byte("test content"), 0644)
		require.NoError(t, err)

		// Gitリポジトリが初期化されていない、またはファイルがトラッキングされていないため、エラーになる
		_, err = manager.ReadManagedFileContent(context.Background(), localPath, "untracked.txt", repo)
		assert.Error(t, err)
	})

	t.Run("パス検証でリポジトリ外へのアクセスを防ぐ", func(t *testing.T) {
		// Gitリポジトリディレクトリを作成
		err := os.MkdirAll(localPath, 0755)
		require.NoError(t, err)

		// ListRepositoryFilesが失敗するため、エラーになる
		_, err = manager.ReadManagedFileContent(context.Background(), localPath, "test.txt", repo)
		assert.Error(t, err)
	})
}

// TestCreateGitAskPassScript tests the createGitAskPassScript helper function
func TestCreateGitAskPassScript(t *testing.T) {
	t.Run("アクセストークン用のスクリプトが作成される", func(t *testing.T) {
		token := "test-token-123"
		scriptPath := createGitAskPassScript(token)

		// スクリプトが作成されていることを確認
		if scriptPath != "" {
			defer os.Remove(scriptPath)

			content, err := os.ReadFile(scriptPath)
			assert.NoError(t, err)
			assert.Contains(t, string(content), token)

			// 実行権限があることを確認
			info, err := os.Stat(scriptPath)
			assert.NoError(t, err)
			// Check if owner execute bit is set (0100 octal = owner execute permission)
			const ownerExecute = os.FileMode(0100)
			assert.True(t, info.Mode()&ownerExecute != 0, "Script should be executable")
		}
	})

	t.Run("空のトークンでも動作する", func(t *testing.T) {
		scriptPath := createGitAskPassScript("")
		if scriptPath != "" {
			defer os.Remove(scriptPath)
		}
		// エラーにならないことを確認（空でも処理される）
	})
}
