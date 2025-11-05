package repository

import (
	"context"
	"errors"
	"opscore/backend/internal/git_repository/domain/entity"
	"opscore/backend/internal/git_repository/domain/repository"
	"opscore/backend/internal/git_repository/infrastructure/git"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// TestRegister はRegisterメソッドのテストです
func TestRegister(t *testing.T) {
	// テスト：有効なリポジトリURLで正常に登録できることを確認する
	t.Run("有効なリポジトリURLで正常に登録できる", func(t *testing.T) {
		// モックの準備
		mockRepo := new(repository.MockRepository)
		mockGitManager := new(git.MockGitManager)

		url := "https://github.com/example/valid-repo"
		accessToken := "test-token"

		// Using matchers for better flexibility
		contextMatcher := mock.MatchedBy(func(ctx context.Context) bool { return true })
		urlMatcher := mock.MatchedBy(func(u string) bool { return u == url })
		repoMatcher := mock.MatchedBy(func(r entity.Repository) bool { return true })

		// モックの振る舞いを定義
		mockRepo.On("FindByURL", contextMatcher, urlMatcher).Return(nil, nil)
		mockRepo.On("Save", contextMatcher, repoMatcher).Return(nil)

		// テスト対象の UseCase を作成
		uc := NewRepositoryUseCase(mockRepo, mockGitManager)

		// テスト実行
		repo, err := uc.Register(context.Background(), url, accessToken)

		// 検証
		assert.NoError(t, err)
		assert.NotNil(t, repo)
		assert.Equal(t, url, repo.URL())
		assert.Equal(t, "valid-repo", repo.Name()) // URLからリポジトリ名が抽出されているか

		// モックの呼び出しを検証
		mockRepo.AssertExpectations(t)
		mockGitManager.AssertExpectations(t)
	})

	// テスト：既に存在するリポジトリURLでエラーになることを確認する
	t.Run("既に存在するリポジトリURLでエラーになる", func(t *testing.T) {
		// モックの準備
		mockRepo := new(repository.MockRepository)
		mockGitManager := new(git.MockGitManager)

		url := "https://github.com/example/existing-repo"
		accessToken := "test-token"
		existingRepo := entity.NewRepository(uuid.NewString(), "existing-repo", url, "")

		// モックの振る舞いを定義
		mockRepo.On("FindByURL", mock.Anything, url).Return(existingRepo, nil)

		// テスト対象の UseCase を作成
		uc := NewRepositoryUseCase(mockRepo, mockGitManager)

		// テスト実行
		repo, err := uc.Register(context.Background(), url, accessToken)

		// 検証
		assert.Error(t, err)
		assert.Equal(t, ErrRepositoryAlreadyExists, err)
		assert.Nil(t, repo)

		// モックの呼び出しを検証
		mockRepo.AssertExpectations(t)
		mockGitManager.AssertExpectations(t)
	})

	// テスト：無効なURLフォーマットでエラーになることを確認する
	t.Run("無効なURLフォーマットでエラーになる", func(t *testing.T) {
		// モックの準備
		mockRepo := new(repository.MockRepository)
		mockGitManager := new(git.MockGitManager)

		invalidURLs := []string{
			"http://github.com/example/repo", // HTTPスキーム (HTTPSのみ許可)
			"https://example.com/repo",       // ホワイトリストにないドメイン
			"github.com/example/repo",        // スキームなし
			"ftp://github.com/example/repo",  // サポートされていないスキーム
			"",                               // 空文字列
		}

		// テスト対象の UseCase を作成
		uc := NewRepositoryUseCase(mockRepo, mockGitManager)

		// 各無効URLでテスト
		for _, url := range invalidURLs {
			repo, err := uc.Register(context.Background(), url, "test-token")
			assert.Error(t, err)
			assert.Nil(t, repo)
		}

		// モックの呼び出しを検証 (この場合、FindByURLは呼ばれない)
		mockRepo.AssertNotCalled(t, "FindByURL")
		mockRepo.AssertNotCalled(t, "Save")
		mockGitManager.AssertExpectations(t)
	})

	// テスト：リポジトリの保存中にエラーが発生した場合を確認する
	t.Run("リポジトリの保存中にエラーが発生した場合", func(t *testing.T) {
		// モックの準備
		mockRepo := new(repository.MockRepository)
		mockGitManager := new(git.MockGitManager)

		url := "https://github.com/example/save-error-repo"
		accessToken := "test-token"
		saveError := errors.New("database error")

		// モックの振る舞いを定義
		mockRepo.On("FindByURL", mock.Anything, url).Return(nil, nil)
		mockRepo.On("Save", mock.Anything, mock.Anything).Return(saveError)

		// テスト対象の UseCase を作成
		uc := NewRepositoryUseCase(mockRepo, mockGitManager)

		// テスト実行
		repo, err := uc.Register(context.Background(), url, accessToken)

		// 検証
		assert.Error(t, err)
		assert.Nil(t, repo)
		assert.Contains(t, err.Error(), "failed to save repository")

		// モックの呼び出しを検証
		mockRepo.AssertExpectations(t)
		mockGitManager.AssertExpectations(t)
	})
}

// TestGetRepository はGetRepositoryメソッドのテストです
func TestGetRepository(t *testing.T) {
	// テスト：存在するリポジトリIDで正常に取得できることを確認する
	t.Run("存在するリポジトリIDで正常に取得できる", func(t *testing.T) {
		// モックの準備
		mockRepo := new(repository.MockRepository)
		mockGitManager := new(git.MockGitManager)

		repoID := uuid.NewString()
		expectedRepo := entity.NewRepository(repoID, "test-repo", "https://github.com/example/test-repo", "test-token")

		// モックの振る舞いを定義
		mockRepo.On("FindByID", mock.Anything, repoID).Return(expectedRepo, nil)

		// テスト対象の UseCase を作成
		uc := NewRepositoryUseCase(mockRepo, mockGitManager)

		// テスト実行
		repo, err := uc.GetRepository(context.Background(), repoID)

		// 検証
		assert.NoError(t, err)
		assert.Equal(t, expectedRepo, repo)

		// モックの呼び出しを検証
		mockRepo.AssertExpectations(t)
		mockGitManager.AssertExpectations(t)
	})

	// テスト：存在しないリポジトリIDでエラーになることを確認する
	t.Run("存在しないリポジトリIDでエラーになる", func(t *testing.T) {
		// モックの準備
		mockRepo := new(repository.MockRepository)
		mockGitManager := new(git.MockGitManager)

		repoID := uuid.NewString()

		// モックの振る舞いを定義
		mockRepo.On("FindByID", mock.Anything, repoID).Return(nil, nil)

		// テスト対象の UseCase を作成
		uc := NewRepositoryUseCase(mockRepo, mockGitManager)

		// テスト実行
		repo, err := uc.GetRepository(context.Background(), repoID)

		// 検証
		assert.Error(t, err)
		assert.Equal(t, ErrRepositoryNotFound, err)
		assert.Nil(t, repo)

		// モックの呼び出しを検証
		mockRepo.AssertExpectations(t)
		mockGitManager.AssertExpectations(t)
	})

	// テスト：リポジトリ検索中にエラーが発生した場合を確認する
	t.Run("リポジトリ検索中にエラーが発生した場合", func(t *testing.T) {
		// モックの準備
		mockRepo := new(repository.MockRepository)
		mockGitManager := new(git.MockGitManager)

		repoID := uuid.NewString()
		findError := errors.New("database error")

		// モックの振る舞いを定義
		mockRepo.On("FindByID", mock.Anything, repoID).Return(nil, findError)

		// テスト対象の UseCase を作成
		uc := NewRepositoryUseCase(mockRepo, mockGitManager)

		// テスト実行
		repo, err := uc.GetRepository(context.Background(), repoID)

		// 検証
		assert.Error(t, err)
		assert.Nil(t, repo)
		assert.Contains(t, err.Error(), "failed to retrieve repository")

		// モックの呼び出しを検証
		mockRepo.AssertExpectations(t)
		mockGitManager.AssertExpectations(t)
	})
}

// TestListRepositories はListRepositoriesメソッドのテストです
func TestListRepositories(t *testing.T) {
	// テスト：リポジトリ一覧が正常に取得できることを確認する
	t.Run("リポジトリ一覧が正常に取得できる", func(t *testing.T) {
		// モックの準備
		mockRepo := new(repository.MockRepository)
		mockGitManager := new(git.MockGitManager)

		expectedRepos := []entity.Repository{
			entity.NewRepository(uuid.NewString(), "repo1", "https://github.com/example/repo1", "token1"),
			entity.NewRepository(uuid.NewString(), "repo2", "https://github.com/example/repo2", "token2"),
		}

		// モックの振る舞いを定義
		mockRepo.On("FindAll", mock.Anything).Return(expectedRepos, nil)

		// テスト対象の UseCase を作成
		uc := NewRepositoryUseCase(mockRepo, mockGitManager)

		// テスト実行
		repos, err := uc.ListRepositories(context.Background())

		// 検証
		assert.NoError(t, err)
		assert.Equal(t, expectedRepos, repos)

		// モックの呼び出しを検証
		mockRepo.AssertExpectations(t)
		mockGitManager.AssertExpectations(t)
	})

	// テスト：リポジトリが1つも存在しない場合は空のスライスが返ることを確認する
	t.Run("リポジトリが1つも存在しない場合は空のスライスが返る", func(t *testing.T) {
		// モックの準備
		mockRepo := new(repository.MockRepository)
		mockGitManager := new(git.MockGitManager)

		// 空のリポジトリリスト
		var emptyRepos []entity.Repository

		// モックの振る舞いを定義
		mockRepo.On("FindAll", mock.Anything).Return(emptyRepos, nil)

		// テスト対象の UseCase を作成
		uc := NewRepositoryUseCase(mockRepo, mockGitManager)

		// テスト実行
		repos, err := uc.ListRepositories(context.Background())

		// 検証
		assert.NoError(t, err)
		assert.Empty(t, repos)

		// モックの呼び出しを検証
		mockRepo.AssertExpectations(t)
		mockGitManager.AssertExpectations(t)
	})

	// テスト：リポジトリ一覧取得中にエラーが発生した場合を確認する
	t.Run("リポジトリ一覧取得中にエラーが発生した場合", func(t *testing.T) {
		// モックの準備
		mockRepo := new(repository.MockRepository)
		mockGitManager := new(git.MockGitManager)

		findAllError := errors.New("database error")

		// モックの振る舞いを定義
		mockRepo.On("FindAll", mock.Anything).Return(nil, findAllError)

		// テスト対象の UseCase を作成
		uc := NewRepositoryUseCase(mockRepo, mockGitManager)

		// テスト実行
		repos, err := uc.ListRepositories(context.Background())

		// 検証
		assert.Error(t, err)
		assert.Nil(t, repos)
		assert.Contains(t, err.Error(), "failed to retrieve repositories")

		// モックの呼び出しを検証
		mockRepo.AssertExpectations(t)
		mockGitManager.AssertExpectations(t)
	})
}

// TestListFiles はListFilesメソッドのテストです
func TestListFiles(t *testing.T) {
	// テスト：存在するリポジトリのファイル一覧が正常に取得できることを確認する
	t.Run("存在するリポジトリのファイル一覧が正常に取得できる", func(t *testing.T) {
		// モックの準備
		mockRepo := new(repository.MockRepository)
		mockGitManager := new(git.MockGitManager)

		repoID := uuid.NewString()
		testRepo := entity.NewRepository(repoID, "test-repo", "https://github.com/example/test-repo", "test-token")
		localPath := "/tmp/repos/" + repoID

		fileList := []string{
			"README.md",
			"src/main.go",
			"docs/index.md",
		}

		// モックの振る舞いを定義
		mockRepo.On("FindByID", mock.Anything, repoID).Return(testRepo, nil)
		mockGitManager.On("EnsureCloned", mock.Anything, testRepo).Return(localPath, nil)
		mockGitManager.On("ListRepositoryFiles", mock.Anything, localPath, testRepo).Return(fileList, nil)

		// テスト対象の UseCase を作成
		uc := NewRepositoryUseCase(mockRepo, mockGitManager)

		// テスト実行
		fileNodes, err := uc.ListFiles(context.Background(), repoID)

		// 検証
		assert.NoError(t, err)
		assert.Len(t, fileNodes, len(fileList))

		// 返されたFileNodeが期待通りかチェック
		for i, f := range fileNodes {
			assert.Equal(t, fileList[i], f.Path())
			assert.Equal(t, "file", f.Type())
		}

		// モックの呼び出しを検証
		mockRepo.AssertExpectations(t)
		mockGitManager.AssertExpectations(t)
	})

	// テスト：存在しないリポジトリIDでエラーになることを確認する
	t.Run("存在しないリポジトリIDでエラーになる", func(t *testing.T) {
		// モックの準備
		mockRepo := new(repository.MockRepository)
		mockGitManager := new(git.MockGitManager)

		repoID := uuid.NewString()

		// モックの振る舞いを定義
		mockRepo.On("FindByID", mock.Anything, repoID).Return(nil, nil)

		// テスト対象の UseCase を作成
		uc := NewRepositoryUseCase(mockRepo, mockGitManager)

		// テスト実行
		fileNodes, err := uc.ListFiles(context.Background(), repoID)

		// 検証
		assert.Error(t, err)
		assert.Equal(t, ErrRepositoryNotFound, err)
		assert.Nil(t, fileNodes)

		// モックの呼び出しを検証
		mockRepo.AssertExpectations(t)
		mockGitManager.AssertExpectations(t)
	})

	// テスト：アクセストークンが設定されていない場合はエラーになることを確認する
	t.Run("アクセストークンが設定されていない場合はエラーになる", func(t *testing.T) {
		// モックの準備
		mockRepo := new(repository.MockRepository)
		mockGitManager := new(git.MockGitManager)

		repoID := uuid.NewString()
		testRepo := entity.NewRepository(repoID, "test-repo", "https://github.com/example/test-repo", "") // 空のアクセストークン

		// モックの振る舞いを定義
		mockRepo.On("FindByID", mock.Anything, repoID).Return(testRepo, nil)

		// テスト対象の UseCase を作成
		uc := NewRepositoryUseCase(mockRepo, mockGitManager)

		// テスト実行
		fileNodes, err := uc.ListFiles(context.Background(), repoID)

		// 検証
		assert.Error(t, err)
		assert.Equal(t, ErrAccessTokenRequired, err)
		assert.Nil(t, fileNodes)

		// モックの呼び出しを検証
		mockRepo.AssertExpectations(t)
		mockGitManager.AssertExpectations(t)
	})

	// テスト：リポジトリのクローン中にエラーが発生した場合を確認する
	t.Run("リポジトリのクローン中にエラーが発生した場合", func(t *testing.T) {
		// モックの準備
		mockRepo := new(repository.MockRepository)
		mockGitManager := new(git.MockGitManager)

		repoID := uuid.NewString()
		testRepo := entity.NewRepository(repoID, "test-repo", "https://github.com/example/test-repo", "test-token")
		cloneError := errors.New("git clone error")

		// モックの振る舞いを定義
		mockRepo.On("FindByID", mock.Anything, repoID).Return(testRepo, nil)
		mockGitManager.On("EnsureCloned", mock.Anything, testRepo).Return("", cloneError)

		// テスト対象の UseCase を作成
		uc := NewRepositoryUseCase(mockRepo, mockGitManager)

		// テスト実行
		fileNodes, err := uc.ListFiles(context.Background(), repoID)

		// 検証
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to ensure repository is cloned")
		assert.Nil(t, fileNodes)

		// モックの呼び出しを検証
		mockRepo.AssertExpectations(t)
		mockGitManager.AssertExpectations(t)
	})

	// テスト：ファイル一覧取得中にエラーが発生した場合を確認する
	t.Run("ファイル一覧取得中にエラーが発生した場合", func(t *testing.T) {
		// モックの準備
		mockRepo := new(repository.MockRepository)
		mockGitManager := new(git.MockGitManager)

		repoID := uuid.NewString()
		testRepo := entity.NewRepository(repoID, "test-repo", "https://github.com/example/test-repo", "test-token")
		localPath := "/tmp/repos/" + repoID
		listError := errors.New("file listing error")

		// モックの振る舞いを定義
		mockRepo.On("FindByID", mock.Anything, repoID).Return(testRepo, nil)
		mockGitManager.On("EnsureCloned", mock.Anything, testRepo).Return(localPath, nil)
		mockGitManager.On("ListRepositoryFiles", mock.Anything, localPath, testRepo).Return(nil, listError)

		// テスト対象の UseCase を作成
		uc := NewRepositoryUseCase(mockRepo, mockGitManager)

		// テスト実行
		fileNodes, err := uc.ListFiles(context.Background(), repoID)

		// 検証
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to list repository files")
		assert.Nil(t, fileNodes)

		// モックの呼び出しを検証
		mockRepo.AssertExpectations(t)
		mockGitManager.AssertExpectations(t)
	})
}

// TestSelectFiles はSelectFilesメソッドのテストです
func TestSelectFiles(t *testing.T) {
	// テスト：ファイル選択が正常に行われることを確認する
	t.Run("ファイル選択が正常に行われる", func(t *testing.T) {
		// モックの準備
		mockRepo := new(repository.MockRepository)
		mockGitManager := new(git.MockGitManager)

		repoID := uuid.NewString()
		testRepo := entity.NewRepository(repoID, "test-repo", "https://github.com/example/test-repo", "test-token")
		localPath := "/tmp/repos/" + repoID
		filePaths := []string{"README.md", "docs/index.md"}

		// モックの振る舞いを定義
		mockRepo.On("FindByID", mock.Anything, repoID).Return(testRepo, nil)
		mockGitManager.On("EnsureCloned", mock.Anything, testRepo).Return(localPath, nil)
		mockGitManager.On("ValidateFilesExist", mock.Anything, localPath, filePaths, testRepo).Return(nil)
		mockRepo.On("SaveManagedFiles", mock.Anything, repoID, filePaths).Return(nil)

		// テスト対象の UseCase を作成
		uc := NewRepositoryUseCase(mockRepo, mockGitManager)

		// テスト実行
		err := uc.SelectFiles(context.Background(), repoID, filePaths)

		// 検証
		assert.NoError(t, err)

		// モックの呼び出しを検証
		mockRepo.AssertExpectations(t)
		mockGitManager.AssertExpectations(t)
	})

	// テスト：存在しないリポジトリIDでエラーになることを確認する
	t.Run("存在しないリポジトリIDでエラーになる", func(t *testing.T) {
		// モックの準備
		mockRepo := new(repository.MockRepository)
		mockGitManager := new(git.MockGitManager)

		repoID := uuid.NewString()
		filePaths := []string{"README.md"}

		// モックの振る舞いを定義
		mockRepo.On("FindByID", mock.Anything, repoID).Return(nil, nil)

		// テスト対象の UseCase を作成
		uc := NewRepositoryUseCase(mockRepo, mockGitManager)

		// テスト実行
		err := uc.SelectFiles(context.Background(), repoID, filePaths)

		// 検証
		assert.Error(t, err)
		assert.Equal(t, ErrRepositoryNotFound, err)

		// モックの呼び出しを検証
		mockRepo.AssertExpectations(t)
		mockGitManager.AssertExpectations(t)
	})

	// テスト：リポジトリのクローン中にエラーが発生した場合を確認する
	t.Run("リポジトリのクローン中にエラーが発生した場合", func(t *testing.T) {
		// モックの準備
		mockRepo := new(repository.MockRepository)
		mockGitManager := new(git.MockGitManager)

		repoID := uuid.NewString()
		testRepo := entity.NewRepository(repoID, "test-repo", "https://github.com/example/test-repo", "test-token")
		filePaths := []string{"README.md"}
		cloneError := errors.New("git clone error")

		// モックの振る舞いを定義
		mockRepo.On("FindByID", mock.Anything, repoID).Return(testRepo, nil)
		mockGitManager.On("EnsureCloned", mock.Anything, testRepo).Return("", cloneError)

		// テスト対象の UseCase を作成
		uc := NewRepositoryUseCase(mockRepo, mockGitManager)

		// テスト実行
		err := uc.SelectFiles(context.Background(), repoID, filePaths)

		// 検証
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to ensure repository is cloned")

		// モックの呼び出しを検証
		mockRepo.AssertExpectations(t)
		mockGitManager.AssertExpectations(t)
	})

	// テスト：ファイル検証中にエラーが発生した場合を確認する
	t.Run("ファイル検証中にエラーが発生した場合", func(t *testing.T) {
		// モックの準備
		mockRepo := new(repository.MockRepository)
		mockGitManager := new(git.MockGitManager)

		repoID := uuid.NewString()
		testRepo := entity.NewRepository(repoID, "test-repo", "https://github.com/example/test-repo", "test-token")
		localPath := "/tmp/repos/" + repoID
		filePaths := []string{"README.md", "nonexistent-file.md"}
		validateError := errors.New("file validation error")

		// モックの振る舞いを定義
		mockRepo.On("FindByID", mock.Anything, repoID).Return(testRepo, nil)
		mockGitManager.On("EnsureCloned", mock.Anything, testRepo).Return(localPath, nil)
		mockGitManager.On("ValidateFilesExist", mock.Anything, localPath, filePaths, testRepo).Return(validateError)

		// テスト対象の UseCase を作成
		uc := NewRepositoryUseCase(mockRepo, mockGitManager)

		// テスト実行
		err := uc.SelectFiles(context.Background(), repoID, filePaths)

		// 検証
		assert.Error(t, err)
		assert.Equal(t, validateError, err)

		// モックの呼び出しを検証
		mockRepo.AssertExpectations(t)
		mockGitManager.AssertExpectations(t)
	})

	// テスト：ファイル保存中にエラーが発生した場合を確認する
	t.Run("ファイル保存中にエラーが発生した場合", func(t *testing.T) {
		// モックの準備
		mockRepo := new(repository.MockRepository)
		mockGitManager := new(git.MockGitManager)

		repoID := uuid.NewString()
		testRepo := entity.NewRepository(repoID, "test-repo", "https://github.com/example/test-repo", "test-token")
		localPath := "/tmp/repos/" + repoID
		filePaths := []string{"README.md"}
		saveError := errors.New("file save error")

		// モックの振る舞いを定義
		mockRepo.On("FindByID", mock.Anything, repoID).Return(testRepo, nil)
		mockGitManager.On("EnsureCloned", mock.Anything, testRepo).Return(localPath, nil)
		mockGitManager.On("ValidateFilesExist", mock.Anything, localPath, filePaths, testRepo).Return(nil)
		mockRepo.On("SaveManagedFiles", mock.Anything, repoID, filePaths).Return(saveError)

		// テスト対象の UseCase を作成
		uc := NewRepositoryUseCase(mockRepo, mockGitManager)

		// テスト実行
		err := uc.SelectFiles(context.Background(), repoID, filePaths)

		// 検証
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to save managed files selection")

		// モックの呼び出しを検証
		mockRepo.AssertExpectations(t)
		mockGitManager.AssertExpectations(t)
	})
}

// TestGetSelectedMarkdown はGetSelectedMarkdownメソッドのテストです
func TestGetSelectedMarkdown(t *testing.T) {
	// テスト：選択されたMarkdownファイルの内容が正常に取得できることを確認する
	t.Run("選択されたMarkdownファイルの内容が正常に取得できる", func(t *testing.T) {
		// モックの準備
		mockRepo := new(repository.MockRepository)
		mockGitManager := new(git.MockGitManager)

		repoID := uuid.NewString()
		testRepo := entity.NewRepository(repoID, "test-repo", "https://github.com/example/test-repo", "test-token")
		localPath := "/tmp/repos/" + repoID
		filePaths := []string{"README.md", "docs/index.md"}

		// ファイル内容
		readme := []byte("# README\n\nThis is a test repository.")
		index := []byte("# Documentation\n\nThis is the index page.")

		// モックの振る舞いを定義
		mockRepo.On("FindByID", mock.Anything, repoID).Return(testRepo, nil)
		mockRepo.On("GetManagedFiles", mock.Anything, repoID).Return(filePaths, nil)
		mockGitManager.On("EnsureCloned", mock.Anything, testRepo).Return(localPath, nil)
		mockGitManager.On("ReadManagedFileContent", mock.Anything, localPath, "README.md", testRepo).Return(readme, nil)
		mockGitManager.On("ReadManagedFileContent", mock.Anything, localPath, "docs/index.md", testRepo).Return(index, nil)

		// テスト対象の UseCase を作成
		uc := NewRepositoryUseCase(mockRepo, mockGitManager)

		// テスト実行
		content, err := uc.GetSelectedMarkdown(context.Background(), repoID)

		// 検証
		assert.NoError(t, err)
		assert.Contains(t, content, string(readme))
		assert.Contains(t, content, string(index))
		assert.Contains(t, content, "---") // セパレータ

		// モックの呼び出しを検証
		mockRepo.AssertExpectations(t)
		mockGitManager.AssertExpectations(t)
	})

	// テスト：存在しないリポジトリIDでエラーになることを確認する
	t.Run("存在しないリポジトリIDでエラーになる", func(t *testing.T) {
		// モックの準備
		mockRepo := new(repository.MockRepository)
		mockGitManager := new(git.MockGitManager)

		repoID := uuid.NewString()

		// モックの振る舞いを定義
		mockRepo.On("FindByID", mock.Anything, repoID).Return(nil, nil)

		// テスト対象の UseCase を作成
		uc := NewRepositoryUseCase(mockRepo, mockGitManager)

		// テスト実行
		content, err := uc.GetSelectedMarkdown(context.Background(), repoID)

		// 検証
		assert.Error(t, err)
		assert.Equal(t, ErrRepositoryNotFound, err)
		assert.Empty(t, content)

		// モックの呼び出しを検証
		mockRepo.AssertExpectations(t)
		mockGitManager.AssertExpectations(t)
	})

	// テスト：管理対象ファイルが選択されていない場合は空文字が返ることを確認する
	t.Run("管理対象ファイルが選択されていない場合は空文字が返る", func(t *testing.T) {
		// モックの準備
		mockRepo := new(repository.MockRepository)
		mockGitManager := new(git.MockGitManager)

		repoID := uuid.NewString()
		testRepo := entity.NewRepository(repoID, "test-repo", "https://github.com/example/test-repo", "test-token")
		var emptyPaths []string

		// モックの振る舞いを定義
		mockRepo.On("FindByID", mock.Anything, repoID).Return(testRepo, nil)
		mockRepo.On("GetManagedFiles", mock.Anything, repoID).Return(emptyPaths, nil)

		// テスト対象の UseCase を作成
		uc := NewRepositoryUseCase(mockRepo, mockGitManager)

		// テスト実行
		content, err := uc.GetSelectedMarkdown(context.Background(), repoID)

		// 検証
		assert.NoError(t, err)
		assert.Empty(t, content)

		// モックの呼び出しを検証
		mockRepo.AssertExpectations(t)
		mockGitManager.AssertExpectations(t)
	})

	// テスト：管理対象ファイル取得中にエラーが発生した場合を確認する
	t.Run("管理対象ファイル取得中にエラーが発生した場合", func(t *testing.T) {
		// モックの準備
		mockRepo := new(repository.MockRepository)
		mockGitManager := new(git.MockGitManager)

		repoID := uuid.NewString()
		testRepo := entity.NewRepository(repoID, "test-repo", "https://github.com/example/test-repo", "test-token")
		getFilesError := errors.New("db error")

		// モックの振る舞いを定義
		mockRepo.On("FindByID", mock.Anything, repoID).Return(testRepo, nil)
		mockRepo.On("GetManagedFiles", mock.Anything, repoID).Return(nil, getFilesError)

		// テスト対象の UseCase を作成
		uc := NewRepositoryUseCase(mockRepo, mockGitManager)

		// テスト実行
		content, err := uc.GetSelectedMarkdown(context.Background(), repoID)

		// 検証
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to retrieve managed files")
		assert.Empty(t, content)

		// モックの呼び出しを検証
		mockRepo.AssertExpectations(t)
		mockGitManager.AssertExpectations(t)
	})

	// テスト：リポジトリのクローン中にエラーが発生した場合を確認する
	t.Run("リポジトリのクローン中にエラーが発生した場合", func(t *testing.T) {
		// モックの準備
		mockRepo := new(repository.MockRepository)
		mockGitManager := new(git.MockGitManager)

		repoID := uuid.NewString()
		testRepo := entity.NewRepository(repoID, "test-repo", "https://github.com/example/test-repo", "test-token")
		filePaths := []string{"README.md"}
		cloneError := errors.New("git clone error")

		// モックの振る舞いを定義
		mockRepo.On("FindByID", mock.Anything, repoID).Return(testRepo, nil)
		mockRepo.On("GetManagedFiles", mock.Anything, repoID).Return(filePaths, nil)
		mockGitManager.On("EnsureCloned", mock.Anything, testRepo).Return("", cloneError)

		// テスト対象の UseCase を作成
		uc := NewRepositoryUseCase(mockRepo, mockGitManager)

		// テスト実行
		content, err := uc.GetSelectedMarkdown(context.Background(), repoID)

		// 検証
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to ensure repository is cloned")
		assert.Empty(t, content)

		// モックの呼び出しを検証
		mockRepo.AssertExpectations(t)
		mockGitManager.AssertExpectations(t)
	})

	// テスト：ファイル読み込み中にエラーが発生した場合を確認する
	t.Run("ファイル読み込み中にエラーが発生した場合", func(t *testing.T) {
		// モックの準備
		mockRepo := new(repository.MockRepository)
		mockGitManager := new(git.MockGitManager)

		repoID := uuid.NewString()
		testRepo := entity.NewRepository(repoID, "test-repo", "https://github.com/example/test-repo", "test-token")
		localPath := "/tmp/repos/" + repoID
		filePaths := []string{"README.md", "error.md"}
		readError := errors.New("file read error")

		// モックの振る舞いを定義
		mockRepo.On("FindByID", mock.Anything, repoID).Return(testRepo, nil)
		mockRepo.On("GetManagedFiles", mock.Anything, repoID).Return(filePaths, nil)
		mockGitManager.On("EnsureCloned", mock.Anything, testRepo).Return(localPath, nil)
		mockGitManager.On("ReadManagedFileContent", mock.Anything, localPath, "README.md", testRepo).Return([]byte("# README"), nil)
		mockGitManager.On("ReadManagedFileContent", mock.Anything, localPath, "error.md", testRepo).Return(nil, readError)

		// テスト対象の UseCase を作成
		uc := NewRepositoryUseCase(mockRepo, mockGitManager)

		// テスト実行
		content, err := uc.GetSelectedMarkdown(context.Background(), repoID)

		// 検証
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to read content of file 'error.md'")
		assert.Empty(t, content)

		// モックの呼び出しを検証
		mockRepo.AssertExpectations(t)
		mockGitManager.AssertExpectations(t)
	})
}

// TestUpdateAccessToken はUpdateAccessTokenメソッドのテストです
func TestUpdateAccessToken(t *testing.T) {
	// テスト：アクセストークンが正常に更新されることを確認する
	t.Run("アクセストークンが正常に更新される", func(t *testing.T) {
		// モックの準備
		mockRepo := new(repository.MockRepository)
		mockGitManager := new(git.MockGitManager)

		repoID := uuid.NewString()
		testRepo := entity.NewRepository(repoID, "test-repo", "https://github.com/example/test-repo", "old-token")
		newToken := "new-token"

		// モックの振る舞いを定義
		mockRepo.On("FindByID", mock.Anything, repoID).Return(testRepo, nil)
		mockRepo.On("UpdateAccessToken", mock.Anything, repoID, newToken).Return(nil)

		// テスト対象の UseCase を作成
		uc := NewRepositoryUseCase(mockRepo, mockGitManager)

		// テスト実行
		err := uc.UpdateAccessToken(context.Background(), repoID, newToken)

		// 検証
		assert.NoError(t, err)

		// モックの呼び出しを検証
		mockRepo.AssertExpectations(t)
		mockGitManager.AssertExpectations(t)
	})

	// テスト：空のアクセストークンでも更新できることを確認する
	t.Run("空のアクセストークンでも更新できる", func(t *testing.T) {
		// モックの準備
		mockRepo := new(repository.MockRepository)
		mockGitManager := new(git.MockGitManager)

		repoID := uuid.NewString()
		testRepo := entity.NewRepository(repoID, "test-repo", "https://github.com/example/test-repo", "old-token")
		emptyToken := ""

		// モックの振る舞いを定義
		mockRepo.On("FindByID", mock.Anything, repoID).Return(testRepo, nil)
		mockRepo.On("UpdateAccessToken", mock.Anything, repoID, emptyToken).Return(nil)

		// テスト対象の UseCase を作成
		uc := NewRepositoryUseCase(mockRepo, mockGitManager)

		// テスト実行
		err := uc.UpdateAccessToken(context.Background(), repoID, emptyToken)

		// 検証
		assert.NoError(t, err)

		// モックの呼び出しを検証
		mockRepo.AssertExpectations(t)
		mockGitManager.AssertExpectations(t)
	})

	// テスト：存在しないリポジトリIDでエラーになることを確認する
	t.Run("存在しないリポジトリIDでエラーになる", func(t *testing.T) {
		// モックの準備
		mockRepo := new(repository.MockRepository)
		mockGitManager := new(git.MockGitManager)

		repoID := uuid.NewString()

		// モックの振る舞いを定義
		mockRepo.On("FindByID", mock.Anything, repoID).Return(nil, nil)

		// テスト対象の UseCase を作成
		uc := NewRepositoryUseCase(mockRepo, mockGitManager)

		// テスト実行
		err := uc.UpdateAccessToken(context.Background(), repoID, "new-token")

		// 検証
		assert.Error(t, err)
		assert.Equal(t, ErrRepositoryNotFound, err)

		// モックの呼び出しを検証
		mockRepo.AssertExpectations(t)
		mockGitManager.AssertExpectations(t)
	})

	// テスト：アクセストークン更新中にエラーが発生した場合を確認する
	t.Run("アクセストークン更新中にエラーが発生した場合", func(t *testing.T) {
		// モックの準備
		mockRepo := new(repository.MockRepository)
		mockGitManager := new(git.MockGitManager)

		repoID := uuid.NewString()
		testRepo := entity.NewRepository(repoID, "test-repo", "https://github.com/example/test-repo", "old-token")
		newToken := "new-token"
		updateError := errors.New("db error")

		// モックの振る舞いを定義
		mockRepo.On("FindByID", mock.Anything, repoID).Return(testRepo, nil)
		mockRepo.On("UpdateAccessToken", mock.Anything, repoID, newToken).Return(updateError)

		// テスト対象の UseCase を作成
		uc := NewRepositoryUseCase(mockRepo, mockGitManager)

		// テスト実行
		err := uc.UpdateAccessToken(context.Background(), repoID, newToken)

		// 検証
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to update repository access token")

		// モックの呼び出しを検証
		mockRepo.AssertExpectations(t)
		mockGitManager.AssertExpectations(t)
	})
}
