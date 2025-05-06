package repository

import (
	"context"
	"opscore/backend/domain/model"
	"opscore/backend/infrastructure/git"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// TestRepositoryUseCaseIntegration tests the repository use case with its real dependencies
func TestRepositoryUseCaseIntegration(t *testing.T) {
	// インメモリリポジトリとモックGitマネージャを使用して統合テスト環境を設定
	repo := NewInMemoryRepository()
	gitManager := git.NewMockGitManager()

	// 実際のユースケース実装を使用（モックではなく）
	useCase := NewRepositoryUseCase(repo, gitManager)
	ctx := context.Background()

	// テスト: Register と GetRepository メソッド
	t.Run("Register and GetRepository", func(t *testing.T) {
		// リポジトリの登録
		repoURL := "https://github.com/example/test-usecase-integration"
		accessToken := "test-token"

		newRepo, err := useCase.Register(ctx, repoURL, accessToken)
		require.NoError(t, err)
		require.NotNil(t, newRepo)

		// 登録したリポジトリがユースケースを通じて取得可能であることを確認
		retrieved, err := useCase.GetRepository(ctx, newRepo.ID())
		require.NoError(t, err)
		require.NotNil(t, retrieved)

		assert.Equal(t, newRepo.ID(), retrieved.ID())
		assert.Equal(t, newRepo.URL(), retrieved.URL())
		assert.Equal(t, "test-usecase-integration", retrieved.Name()) // URLから抽出されるはず
	})

	// テスト: 同じURLで2回登録すると conflict エラーになる
	t.Run("Register with duplicate URL", func(t *testing.T) {
		repoURL := "https://github.com/example/duplicate-repo"
		accessToken := "test-token"

		// 1回目の登録
		_, err := useCase.Register(ctx, repoURL, accessToken)
		require.NoError(t, err)

		// 2回目の登録（同じURL）
		_, err = useCase.Register(ctx, repoURL, accessToken)
		assert.ErrorIs(t, err, ErrRepositoryAlreadyExists)
	})

	// テスト: ListRepositories でリポジトリ一覧が取得できる
	t.Run("ListRepositories", func(t *testing.T) {
		// すでに少なくとも2つのリポジトリが登録されているはず
		repos, err := useCase.ListRepositories(ctx)
		require.NoError(t, err)
		assert.GreaterOrEqual(t, len(repos), 2)
	})

	// テスト: UpdateAccessToken でアクセストークンが更新できる
	t.Run("UpdateAccessToken", func(t *testing.T) {
		// 新しいリポジトリを登録
		repoURL := "https://github.com/example/update-token-repo"
		initialToken := "initial-token"

		newRepo, err := useCase.Register(ctx, repoURL, initialToken)
		require.NoError(t, err)

		// アクセストークンを更新
		updatedToken := "updated-token"
		err = useCase.UpdateAccessToken(ctx, newRepo.ID(), updatedToken)
		require.NoError(t, err)

		// Note: リポジトリエンティティはアクセストークンを外部に公開しないので
		// 直接検証はできないが、内部的に更新されていることを前提とする
	})

	// テスト: SelectFiles と GetSelectedMarkdown
	t.Run("SelectFiles and GetSelectedMarkdown", func(t *testing.T) {
		// マネージドファイルの選択と Markdown 取得のテストのため
		// GitManager をモックしてファイルの内容を返すように設定

		// 新しいリポジトリを登録
		repoURL := "https://github.com/example/markdown-repo"
		accessToken := "token"

		newRepo, err := useCase.Register(ctx, repoURL, accessToken)
		require.NoError(t, err)
		repoID := newRepo.ID()

		// MockGitManager にテスト用の振る舞いを設定
		mockGitManager := git.NewMockGitManager()

		// Using contextMatcher to match any context
		contextMatcher := mock.MatchedBy(func(ctx context.Context) bool { return true })
		repoMatcher := mock.MatchedBy(func(repo model.Repository) bool {
			return repo.ID() == repoID
		})

		// Set up mock expectations correctly with context and repo parameters
		mockPath := "/mock/path/to/repo"
		mockGitManager.On("EnsureCloned", contextMatcher, repoMatcher).Return(mockPath, nil)

		// Add mock files
		mockGitManager.AddMockFile(repoID, "README.md", "# Test Repository\n\nThis is a test.")
		mockGitManager.AddMockFile(repoID, "docs/adr/0001-test-adr.md", "# ADR 0001\n\nThis is a test ADR.")

		// Set up expectations for other methods
		filePathsMatcher := mock.MatchedBy(func(paths []string) bool { return true })
		mockGitManager.On("ValidateFilesExist", contextMatcher, mockPath, filePathsMatcher, repoMatcher).Return(nil)

		readmeMatcher := mock.MatchedBy(func(path string) bool { return path == "README.md" })
		adrMatcher := mock.MatchedBy(func(path string) bool { return path == "docs/adr/0001-test-adr.md" })

		mockGitManager.On("ReadManagedFileContent", contextMatcher, mockPath, readmeMatcher, repoMatcher).
			Return([]byte("# Test Repository\n\nThis is a test."), nil)
		mockGitManager.On("ReadManagedFileContent", contextMatcher, mockPath, adrMatcher, repoMatcher).
			Return([]byte("# ADR 0001\n\nThis is a test ADR."), nil)

		// 実際のユースケースのGitManagerをモックに置き換え
		useCaseWithMock := NewRepositoryUseCase(repo, mockGitManager)

		// ファイルを選択
		filePaths := []string{"README.md", "docs/adr/0001-test-adr.md"}
		err = useCaseWithMock.SelectFiles(ctx, repoID, filePaths)
		require.NoError(t, err)

		// マークダウン内容を取得
		markdown, err := useCaseWithMock.GetSelectedMarkdown(ctx, repoID)
		require.NoError(t, err)

		// 両方のファイルの内容が連結されていることを確認
		assert.Contains(t, markdown, "# Test Repository")
		assert.Contains(t, markdown, "# ADR 0001")
	})

	// テスト: 存在しないリポジトリIDでのエラー処理
	t.Run("Error handling with non-existent repository", func(t *testing.T) {
		nonExistentID := "non-existent-id"

		// GetRepository
		_, err := useCase.GetRepository(ctx, nonExistentID)
		assert.ErrorIs(t, err, ErrRepositoryNotFound)

		// UpdateAccessToken
		err = useCase.UpdateAccessToken(ctx, nonExistentID, "token")
		assert.ErrorIs(t, err, ErrRepositoryNotFound)

		// SelectFiles
		err = useCase.SelectFiles(ctx, nonExistentID, []string{"file.md"})
		assert.ErrorIs(t, err, ErrRepositoryNotFound)

		// GetSelectedMarkdown
		_, err = useCase.GetSelectedMarkdown(ctx, nonExistentID)
		assert.ErrorIs(t, err, ErrRepositoryNotFound)
	})
}

// テスト用のインメモリリポジトリ実装
type InMemoryRepository struct {
	repositories map[string]model.Repository
	managedFiles map[string][]string
}

func NewInMemoryRepository() *InMemoryRepository {
	return &InMemoryRepository{
		repositories: make(map[string]model.Repository),
		managedFiles: make(map[string][]string),
	}
}

func (r *InMemoryRepository) Save(ctx context.Context, repo model.Repository) error {
	r.repositories[repo.ID()] = repo
	return nil
}

func (r *InMemoryRepository) FindByID(ctx context.Context, id string) (model.Repository, error) {
	repo, exists := r.repositories[id]
	if !exists {
		return nil, nil // Not found
	}
	return repo, nil
}

func (r *InMemoryRepository) FindByURL(ctx context.Context, url string) (model.Repository, error) {
	for _, repo := range r.repositories {
		if repo.URL() == url {
			return repo, nil
		}
	}
	return nil, nil // Not found
}

func (r *InMemoryRepository) FindAll(ctx context.Context) ([]model.Repository, error) {
	repos := make([]model.Repository, 0, len(r.repositories))
	for _, repo := range r.repositories {
		repos = append(repos, repo)
	}
	return repos, nil
}

func (r *InMemoryRepository) UpdateAccessToken(ctx context.Context, id string, accessToken string) error {
	repo, exists := r.repositories[id]
	if !exists {
		return ErrRepositoryNotFound
	}

	// Create a new repository with updated token
	updatedRepo := model.NewRepository(
		repo.ID(),
		repo.Name(),
		repo.URL(),
		accessToken,
	)
	r.repositories[id] = updatedRepo
	return nil
}

func (r *InMemoryRepository) SaveManagedFiles(ctx context.Context, repoID string, filePaths []string) error {
	_, exists := r.repositories[repoID]
	if !exists {
		return ErrRepositoryNotFound
	}
	r.managedFiles[repoID] = filePaths
	return nil
}

func (r *InMemoryRepository) GetManagedFiles(ctx context.Context, repoID string) ([]string, error) {
	filePaths, exists := r.managedFiles[repoID]
	if !exists {
		return []string{}, nil // Empty slice for no files
	}
	return filePaths, nil
}
