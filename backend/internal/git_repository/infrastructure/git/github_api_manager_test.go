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

// Note: EnsureCloned and downloadRepository methods are not fully tested in unit tests
// as they require actual GitHub API interactions. Integration tests with mock GitHub API
// servers or real repositories would be needed for comprehensive coverage of these methods.
// These methods are indirectly tested through integration tests at higher layers.

// TestNewGithubApiManager tests the initialization of githubApiManager
func TestNewGithubApiManager(t *testing.T) {
	t.Run("正常にディレクトリが作成される", func(t *testing.T) {
		tmpDir := filepath.Join(os.TempDir(), "test-github-api-manager")
		defer os.RemoveAll(tmpDir)

		manager, err := NewGithubApiManager(tmpDir)

		assert.NoError(t, err)
		assert.NotNil(t, manager)

		// ディレクトリが作成されていることを確認
		info, err := os.Stat(tmpDir)
		assert.NoError(t, err)
		assert.True(t, info.IsDir())
	})

	t.Run("既存のディレクトリで初期化できる", func(t *testing.T) {
		tmpDir := filepath.Join(os.TempDir(), "test-github-api-manager-existing")
		defer os.RemoveAll(tmpDir)

		// 事前にディレクトリを作成
		err := os.MkdirAll(tmpDir, 0755)
		require.NoError(t, err)

		manager, err := NewGithubApiManager(tmpDir)

		assert.NoError(t, err)
		assert.NotNil(t, manager)
	})

	t.Run("無効なパスでエラーになる", func(t *testing.T) {
		// 無効なパス（書き込み権限のないパスをシミュレート）
		invalidPath := "/root/invalid-test-path"

		_, err := NewGithubApiManager(invalidPath)

		assert.Error(t, err)
	})

	t.Run("クライアントキャッシュが初期化される", func(t *testing.T) {
		tmpDir := filepath.Join(os.TempDir(), "test-github-api-cache")
		defer os.RemoveAll(tmpDir)

		manager, err := NewGithubApiManager(tmpDir)
		require.NoError(t, err)

		githubManager := manager.(*githubApiManager)
		assert.NotNil(t, githubManager.clients)
		assert.Equal(t, 0, len(githubManager.clients))
	})
}

// TestGetLocalPath tests the getLocalPath method for githubApiManager
func TestGetLocalPathGitHub(t *testing.T) {
	t.Run("リポジトリIDから正しいローカルパスを生成する", func(t *testing.T) {
		tmpDir := filepath.Join(os.TempDir(), "test-github-local-path")
		defer os.RemoveAll(tmpDir)

		manager, err := NewGithubApiManager(tmpDir)
		require.NoError(t, err)

		githubManager := manager.(*githubApiManager)
		repo := entity.NewRepository("test-repo-id", "test-repo", "https://github.com/example/test", "")

		localPath := githubManager.getLocalPath(repo)

		expectedPath := filepath.Join(tmpDir, "test-repo-id")
		assert.Equal(t, expectedPath, localPath)
	})
}

// TestParseGitHubURL tests the parseGitHubURL helper function
func TestParseGitHubURL(t *testing.T) {
	t.Run("有効なGitHub URLをパースできる", func(t *testing.T) {
		testCases := []struct {
			name          string
			url           string
			expectedOwner string
			expectedRepo  string
		}{
			{
				name:          "HTTPS URL without .git",
				url:           "https://github.com/example/test-repo",
				expectedOwner: "example",
				expectedRepo:  "test-repo",
			},
			{
				name:          "HTTPS URL with .git",
				url:           "https://github.com/example/test-repo.git",
				expectedOwner: "example",
				expectedRepo:  "test-repo",
			},
			{
				name:          "URL with organization",
				url:           "https://github.com/my-org/my-repo",
				expectedOwner: "my-org",
				expectedRepo:  "my-repo",
			},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				owner, repo, err := parseGitHubURL(tc.url)

				assert.NoError(t, err)
				assert.Equal(t, tc.expectedOwner, owner)
				assert.Equal(t, tc.expectedRepo, repo)
			})
		}
	})

	t.Run("無効なURLでエラーになる", func(t *testing.T) {
		testCases := []struct {
			name string
			url  string
		}{
			{
				name: "短すぎるURL",
				url:  "https://github.com",
			},
			{
				name: "パスが不足",
				url:  "https://github.com/example",
			},
			{
				name: "空のURL",
				url:  "",
			},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				_, _, err := parseGitHubURL(tc.url)
				assert.Error(t, err)
				assert.Contains(t, err.Error(), "invalid GitHub URL format")
			})
		}
	})
}

// TestGetGitHubClient tests the getGitHubClient method
func TestGetGitHubClient(t *testing.T) {
	t.Run("トークンなしでクライアントを作成できる", func(t *testing.T) {
		tmpDir := filepath.Join(os.TempDir(), "test-github-client-1")
		defer os.RemoveAll(tmpDir)

		manager, err := NewGithubApiManager(tmpDir)
		require.NoError(t, err)
		githubManager := manager.(*githubApiManager)

		client := githubManager.getGitHubClient("")

		assert.NotNil(t, client)
	})

	t.Run("トークンありでクライアントを作成できる", func(t *testing.T) {
		tmpDir := filepath.Join(os.TempDir(), "test-github-client-2")
		defer os.RemoveAll(tmpDir)

		manager, err := NewGithubApiManager(tmpDir)
		require.NoError(t, err)
		githubManager := manager.(*githubApiManager)

		token := "test-token-123"
		client := githubManager.getGitHubClient(token)

		assert.NotNil(t, client)

		// クライアントがキャッシュされていることを確認
		assert.Equal(t, 1, len(githubManager.clients))
		cachedClient := githubManager.clients[token]
		assert.Equal(t, client, cachedClient)
	})

	t.Run("同じトークンでキャッシュされたクライアントを返す", func(t *testing.T) {
		tmpDir := filepath.Join(os.TempDir(), "test-github-client-3")
		defer os.RemoveAll(tmpDir)

		manager, err := NewGithubApiManager(tmpDir)
		require.NoError(t, err)
		githubManager := manager.(*githubApiManager)

		token := "cached-token"
		client1 := githubManager.getGitHubClient(token)
		client2 := githubManager.getGitHubClient(token)

		assert.Equal(t, client1, client2)
	})

	t.Run("異なるトークンで別のクライアントを作成する", func(t *testing.T) {
		tmpDir := filepath.Join(os.TempDir(), "test-github-client-4")
		defer os.RemoveAll(tmpDir)

		manager, err := NewGithubApiManager(tmpDir)
		require.NoError(t, err)
		githubManager := manager.(*githubApiManager)

		token1 := "token-1"
		token2 := "token-2"

		client1 := githubManager.getGitHubClient(token1)
		client2 := githubManager.getGitHubClient(token2)

		assert.NotEqual(t, client1, client2)
		assert.Equal(t, 2, len(githubManager.clients))
	})
}

// TestListRepositoryFilesGitHub tests the ListRepositoryFiles method for githubApiManager
func TestListRepositoryFilesGitHub(t *testing.T) {
	tmpDir := filepath.Join(os.TempDir(), "test-github-list-files")
	defer os.RemoveAll(tmpDir)

	manager, err := NewGithubApiManager(tmpDir)
	require.NoError(t, err)

	repo := entity.NewRepository("test-id", "test-repo", "https://github.com/example/test", "")
	localPath := filepath.Join(tmpDir, "test-id")

	t.Run("ローカルディレクトリが存在しない場合はエラー", func(t *testing.T) {
		// ローカルディレクトリが存在しない場合、APIフォールバックが試行される
		// 実際のAPIリクエストは行われず、モックされていないためエラーになる
		_, err := manager.ListRepositoryFiles(context.Background(), localPath, repo)
		assert.Error(t, err)
	})

	t.Run("ローカルディレクトリにファイルがある場合はリストを返す", func(t *testing.T) {
		// テスト用のファイル構造を作成
		err := os.MkdirAll(localPath, 0755)
		require.NoError(t, err)

		// サブディレクトリとファイルを作成
		err = os.MkdirAll(filepath.Join(localPath, "dir1"), 0755)
		require.NoError(t, err)

		err = os.WriteFile(filepath.Join(localPath, "file1.txt"), []byte("content1"), 0644)
		require.NoError(t, err)

		err = os.WriteFile(filepath.Join(localPath, "dir1", "file2.txt"), []byte("content2"), 0644)
		require.NoError(t, err)

		files, err := manager.ListRepositoryFiles(context.Background(), localPath, repo)

		assert.NoError(t, err)
		assert.NotNil(t, files)
		assert.Contains(t, files, "file1.txt")
		assert.Contains(t, files, "dir1/file2.txt")
	})

	t.Run("隠しファイルは除外される", func(t *testing.T) {
		localPath := filepath.Join(tmpDir, "test-hidden-files")
		err := os.MkdirAll(localPath, 0755)
		require.NoError(t, err)
		defer os.RemoveAll(localPath)

		// 通常のファイルと隠しファイルを作成
		err = os.WriteFile(filepath.Join(localPath, "visible.txt"), []byte("content"), 0644)
		require.NoError(t, err)

		err = os.WriteFile(filepath.Join(localPath, ".hidden.txt"), []byte("content"), 0644)
		require.NoError(t, err)

		files, err := manager.ListRepositoryFiles(context.Background(), localPath, repo)

		assert.NoError(t, err)
		assert.Contains(t, files, "visible.txt")
		assert.NotContains(t, files, ".hidden.txt")
	})
}

// TestValidateFilesExistGitHub tests the ValidateFilesExist method for githubApiManager
func TestValidateFilesExistGitHub(t *testing.T) {
	tmpDir := filepath.Join(os.TempDir(), "test-github-validate-files")
	defer os.RemoveAll(tmpDir)

	manager, err := NewGithubApiManager(tmpDir)
	require.NoError(t, err)

	repo := entity.NewRepository("test-id", "test-repo", "https://github.com/example/test", "")
	localPath := filepath.Join(tmpDir, "test-id")

	t.Run("ファイルパスが空の場合はエラーなし", func(t *testing.T) {
		err := manager.ValidateFilesExist(context.Background(), localPath, []string{}, repo)
		assert.NoError(t, err)
	})

	t.Run("存在するファイルは検証成功", func(t *testing.T) {
		err := os.MkdirAll(localPath, 0755)
		require.NoError(t, err)

		testFile := filepath.Join(localPath, "test.txt")
		err = os.WriteFile(testFile, []byte("test"), 0644)
		require.NoError(t, err)

		err = manager.ValidateFilesExist(context.Background(), localPath, []string{"test.txt"}, repo)
		assert.NoError(t, err)
	})

	t.Run("存在しないファイルでエラーになる", func(t *testing.T) {
		err := manager.ValidateFilesExist(context.Background(), localPath, []string{"nonexistent.txt"}, repo)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "does not exist")
	})

	t.Run("複数ファイルの一部が存在しない場合はエラー", func(t *testing.T) {
		err := manager.ValidateFilesExist(context.Background(), localPath, []string{"test.txt", "missing.txt"}, repo)
		assert.Error(t, err)
	})
}

// TestReadManagedFileContentGitHub tests the ReadManagedFileContent method for githubApiManager
func TestReadManagedFileContentGitHub(t *testing.T) {
	tmpDir := filepath.Join(os.TempDir(), "test-github-read-file")
	defer os.RemoveAll(tmpDir)

	manager, err := NewGithubApiManager(tmpDir)
	require.NoError(t, err)

	repo := entity.NewRepository("test-id", "test-repo", "https://github.com/example/test", "")
	localPath := filepath.Join(tmpDir, "test-id")

	t.Run("パストラバーサル攻撃を防ぐ - .. を含むパス", func(t *testing.T) {
		_, err := manager.ReadManagedFileContent(context.Background(), localPath, "../etc/passwd", repo)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "path traversal")
	})

	t.Run("正常にファイル内容を読み取れる", func(t *testing.T) {
		err := os.MkdirAll(localPath, 0755)
		require.NoError(t, err)

		testContent := "test file content"
		testFile := filepath.Join(localPath, "test.txt")
		err = os.WriteFile(testFile, []byte(testContent), 0644)
		require.NoError(t, err)

		content, err := manager.ReadManagedFileContent(context.Background(), localPath, "test.txt", repo)

		assert.NoError(t, err)
		assert.Equal(t, testContent, string(content))
	})

	t.Run("サブディレクトリのファイルを読み取れる", func(t *testing.T) {
		subDir := filepath.Join(localPath, "subdir")
		err := os.MkdirAll(subDir, 0755)
		require.NoError(t, err)

		testContent := "nested file content"
		testFile := filepath.Join(subDir, "nested.txt")
		err = os.WriteFile(testFile, []byte(testContent), 0644)
		require.NoError(t, err)

		content, err := manager.ReadManagedFileContent(context.Background(), localPath, "subdir/nested.txt", repo)

		assert.NoError(t, err)
		assert.Equal(t, testContent, string(content))
	})

	t.Run("存在しないファイルでエラーになる", func(t *testing.T) {
		// ローカルに存在しない場合、APIフォールバックが試行される
		// 実際のAPIリクエストは行われず、エラーになる
		_, err := manager.ReadManagedFileContent(context.Background(), localPath, "nonexistent.txt", repo)
		assert.Error(t, err)
	})

	t.Run("リポジトリディレクトリ外へのアクセスを防ぐ", func(t *testing.T) {
		// 絶対パスの検証により、リポジトリディレクトリ外へのアクセスが防がれる
		// この場合はエラーが返されるべき
		_, err := manager.ReadManagedFileContent(context.Background(), localPath, "../../etc/passwd", repo)
		assert.Error(t, err)
	})
}
