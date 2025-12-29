package repository

import (
	"context"
	"opscore/backend/internal/git_repository/domain/entity"
	"opscore/backend/internal/git_repository/infrastructure/git"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestRepositoryUseCaseIntegration tests the repository use case with its real dependencies
func TestRepositoryUseCaseIntegration(t *testing.T) {
	// インメモリリポジトリとモックGitマネージャを使用して統合テスト環境を設定
	repo := NewInMemoryRepository()
	gitManager := git.NewMockGitManager()
	oauthProvider := new(MockOAuthTokenProvider)

	// 実際のユースケース実装を使用（モックではなく）
	useCase := NewRepositoryUseCase(repo, gitManager, oauthProvider)
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
}

// テスト用のインメモリリポジトリ実装
type InMemoryRepository struct {
	repositories map[string]entity.Repository
	managedFiles map[string][]string
}

func NewInMemoryRepository() *InMemoryRepository {
	return &InMemoryRepository{
		repositories: make(map[string]entity.Repository),
		managedFiles: make(map[string][]string),
	}
}

func (r *InMemoryRepository) Save(ctx context.Context, repo entity.Repository) error {
	r.repositories[repo.ID()] = repo
	return nil
}

func (r *InMemoryRepository) FindByID(ctx context.Context, id string) (entity.Repository, error) {
	repo, exists := r.repositories[id]
	if !exists {
		return nil, nil // Not found
	}
	return repo, nil
}

func (r *InMemoryRepository) FindByURL(ctx context.Context, url string) (entity.Repository, error) {
	for _, repo := range r.repositories {
		if repo.URL() == url {
			return repo, nil
		}
	}
	return nil, nil // Not found
}

func (r *InMemoryRepository) FindAll(ctx context.Context) ([]entity.Repository, error) {
	repos := make([]entity.Repository, 0, len(r.repositories))
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
	updatedRepo := entity.NewRepository(
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
