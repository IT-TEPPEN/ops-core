package persistence

import (
	"context"
	"crypto/rand"
	"opscore/backend/internal/git_repository/domain/entity"
	"opscore/backend/internal/git_repository/infrastructure/encryption"
	"os"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

// PostgresRepositoryTestSuite はPostgresRepositoryのテストスイートです。
type PostgresRepositoryTestSuite struct {
	suite.Suite
	db         *pgxpool.Pool
	repository *PostgresRepository
	ctx        context.Context
}

// TestPostgresRepository は統合テストスイートを実行します。
func TestPostgresRepository(t *testing.T) {
	// テスト用のデータベース接続がない場合はスキップ
	// 実際の運用では環境変数やテスト設定から取得します
	dbURL := os.Getenv("TEST_DATABASE_URL")
	if dbURL == "" {
		t.Skip("テスト用のデータベース接続情報がありません。TEST_DATABASE_URL環境変数を設定してください。")
	}

	suite.Run(t, new(PostgresRepositoryTestSuite))
}

// SetupSuite はテストスイート全体の前処理を行います。
func (s *PostgresRepositoryTestSuite) SetupSuite() {
	// テスト用のデータベース接続を設定
	dbURL := os.Getenv("TEST_DATABASE_URL")
	var err error
	s.ctx = context.Background()
	s.db, err = pgxpool.New(s.ctx, dbURL)
	if err != nil {
		s.T().Fatalf("テスト用のデータベース接続に失敗しました: %v", err)
	}

	// テスト用の暗号化キーを生成
	key := make([]byte, 32)
	_, err = rand.Read(key)
	if err != nil {
		s.T().Fatalf("テスト用の暗号化キーの生成に失敗しました: %v", err)
	}

	encryptor, err := encryption.NewEncryptor(key)
	if err != nil {
		s.T().Fatalf("暗号化器の作成に失敗しました: %v", err)
	}

	// PostgresRepositoryのインスタンスを作成
	s.repository = &PostgresRepository{
		db:        s.db,
		encryptor: encryptor,
	}
}

// TearDownSuite はテストスイート全体の後処理を行います。
func (s *PostgresRepositoryTestSuite) TearDownSuite() {
	// データベース接続を閉じる
	if s.db != nil {
		s.db.Close()
	}
}

// SetupTest は各テストの前処理を行います。
func (s *PostgresRepositoryTestSuite) SetupTest() {
	// テスト用のテーブルをクリーンアップ
	_, err := s.db.Exec(s.ctx, "DELETE FROM managed_files")
	if err != nil {
		s.T().Fatalf("managed_filesテーブルのクリーンアップに失敗しました: %v", err)
	}
	_, err = s.db.Exec(s.ctx, "DELETE FROM repositories")
	if err != nil {
		s.T().Fatalf("repositoriesテーブルのクリーンアップに失敗しました: %v", err)
	}
}

// TearDownTest は各テストの後処理を行います。
func (s *PostgresRepositoryTestSuite) TearDownTest() {
	// 必要に応じてテスト後のクリーンアップを行います
}

// createTestRepository はテスト用のリポジトリオブジェクトを作成します。
func (s *PostgresRepositoryTestSuite) createTestRepository() entity.Repository {
	id := uuid.New().String()
	name := "test-repo-" + id[:8]
	url := "https://github.com/example/" + name
	accessToken := "test-token-" + id[:8]
	return entity.NewRepository(id, name, url, accessToken)
}

// TestSave はSaveメソッドのテストを行います。
func (s *PostgresRepositoryTestSuite) TestSave() {
	// テスト：新しいリポジトリが正常に保存されることを確認する
	s.T().Run("新しいリポジトリが正常に保存される", func(t *testing.T) {
		repo := s.createTestRepository()

		err := s.repository.Save(s.ctx, repo)

		assert.NoError(t, err)

		// 保存されたリポジトリを検証
		savedRepo, err := s.repository.FindByID(s.ctx, repo.ID())
		assert.NoError(t, err)
		assert.NotNil(t, savedRepo)
		assert.Equal(t, repo.ID(), savedRepo.ID())
		assert.Equal(t, repo.Name(), savedRepo.Name())
		assert.Equal(t, repo.URL(), savedRepo.URL())
		assert.Equal(t, repo.AccessToken(), savedRepo.AccessToken())
	})

	// テスト：同じIDの既存リポジトリが更新されることを確認する
	s.T().Run("同じIDの既存リポジトリが更新される", func(t *testing.T) {
		// 初回保存
		repo := s.createTestRepository()
		err := s.repository.Save(s.ctx, repo)
		assert.NoError(t, err)

		// 同じIDで別の値を持つリポジトリを作成
		updatedRepo := entity.ReconstructRepository(
			repo.ID(),
			repo.Name()+"-updated",
			repo.URL()+"-updated",
			repo.AccessToken()+"-updated",
			repo.CreatedAt(),
			time.Now().Add(time.Hour), // 1時間後に更新
		)

		// 更新
		err = s.repository.Save(s.ctx, updatedRepo)
		assert.NoError(t, err)

		// 更新されたリポジトリを検証
		savedRepo, err := s.repository.FindByID(s.ctx, repo.ID())
		assert.NoError(t, err)
		assert.NotNil(t, savedRepo)
		assert.Equal(t, updatedRepo.Name(), savedRepo.Name())
		assert.Equal(t, updatedRepo.URL(), savedRepo.URL())
		assert.Equal(t, updatedRepo.AccessToken(), savedRepo.AccessToken())
	})

	// テスト：同じURLの別リポジトリ保存は一意制約違反になることを確認する
	s.T().Run("同じURLの別リポジトリ保存は一意制約違反になる", func(t *testing.T) {
		// 初回保存
		repo1 := s.createTestRepository()
		err := s.repository.Save(s.ctx, repo1)
		assert.NoError(t, err)

		// 同じURLで異なるIDのリポジトリを作成
		repo2 := entity.NewRepository(
			uuid.New().String(),
			repo1.Name()+"-second",
			repo1.URL(), // 同じURL
			"different-token",
		)

		// 保存を試みる
		err = s.repository.Save(s.ctx, repo2)
		assert.Error(t, err) // エラーを期待
		// 注：実装によっては特定のエラータイプをチェックすることも可能
	})
}

// TestFindByURL はFindByURLメソッドのテストを行います。
func (s *PostgresRepositoryTestSuite) TestFindByURL() {
	// テスト：存在するURLでリポジトリが取得できることを確認する
	s.T().Run("存在するURLでリポジトリが取得できる", func(t *testing.T) {
		// テストリポジトリを作成して保存
		repo := s.createTestRepository()
		err := s.repository.Save(s.ctx, repo)
		assert.NoError(t, err)

		// URLで検索
		found, err := s.repository.FindByURL(s.ctx, repo.URL())

		// 検証
		assert.NoError(t, err)
		assert.NotNil(t, found)
		assert.Equal(t, repo.ID(), found.ID())
		assert.Equal(t, repo.Name(), found.Name())
		assert.Equal(t, repo.URL(), found.URL())
		assert.Equal(t, repo.AccessToken(), found.AccessToken())
	})

	// テスト：存在しないURLではnilが返ることを確認する
	s.T().Run("存在しないURLではnilが返る", func(t *testing.T) {
		found, err := s.repository.FindByURL(s.ctx, "https://nonexistent-url.example.com")

		assert.NoError(t, err)
		assert.Nil(t, found)
	})
}

// TestFindByID はFindByIDメソッドのテストを行います。
func (s *PostgresRepositoryTestSuite) TestFindByID() {
	// テスト：存在するIDでリポジトリが取得できることを確認する
	s.T().Run("存在するIDでリポジトリが取得できる", func(t *testing.T) {
		// テストリポジトリを作成して保存
		repo := s.createTestRepository()
		err := s.repository.Save(s.ctx, repo)
		assert.NoError(t, err)

		// IDで検索
		found, err := s.repository.FindByID(s.ctx, repo.ID())

		// 検証
		assert.NoError(t, err)
		assert.NotNil(t, found)
		assert.Equal(t, repo.ID(), found.ID())
		assert.Equal(t, repo.Name(), found.Name())
		assert.Equal(t, repo.URL(), found.URL())
		assert.Equal(t, repo.AccessToken(), found.AccessToken())
	})

	// テスト：存在しないIDではnilが返ることを確認する
	s.T().Run("存在しないIDではnilが返る", func(t *testing.T) {
		found, err := s.repository.FindByID(s.ctx, uuid.New().String())

		assert.NoError(t, err)
		assert.Nil(t, found)
	})
}

// TestFindAll はFindAllメソッドのテストを行います。
func (s *PostgresRepositoryTestSuite) TestFindAll() {
	// テスト：保存したリポジトリがすべて取得できることを確認する
	s.T().Run("保存したリポジトリがすべて取得できる", func(t *testing.T) {
		// テストリポジトリを作成して保存
		repo1 := s.createTestRepository()
		err := s.repository.Save(s.ctx, repo1)
		assert.NoError(t, err)

		repo2 := s.createTestRepository()
		err = s.repository.Save(s.ctx, repo2)
		assert.NoError(t, err)

		// 全リポジトリを検索
		repos, err := s.repository.FindAll(s.ctx)

		// 検証
		assert.NoError(t, err)
		assert.Len(t, repos, 2)

		// IDでリポジトリをマップ化して検索しやすくする
		repoMap := make(map[string]entity.Repository)
		for _, r := range repos {
			repoMap[r.ID()] = r
		}

		// 各リポジトリが存在することを確認
		assert.Contains(t, repoMap, repo1.ID())
		assert.Contains(t, repoMap, repo2.ID())
		assert.Equal(t, repo1.URL(), repoMap[repo1.ID()].URL())
		assert.Equal(t, repo2.URL(), repoMap[repo2.ID()].URL())
	})

	// テスト：リポジトリが存在しない場合は空のスライスが返ることを確認する
	s.T().Run("リポジトリが存在しない場合は空のスライスが返る", func(t *testing.T) {
		// 事前にすべてのリポジトリを削除
		_, err := s.db.Exec(s.ctx, "DELETE FROM managed_files")
		assert.NoError(t, err)
		_, err = s.db.Exec(s.ctx, "DELETE FROM repositories")
		assert.NoError(t, err)

		// 全リポジトリを検索
		repos, err := s.repository.FindAll(s.ctx)

		// 検証
		assert.NoError(t, err)
		assert.NotNil(t, repos)
		assert.Empty(t, repos)
	})
}

// TestSaveAndGetManagedFiles はSaveManagedFilesとGetManagedFilesメソッドのテストを行います。
func (s *PostgresRepositoryTestSuite) TestSaveAndGetManagedFiles() {
	// テスト：管理対象ファイルが正しく保存・取得できることを確認する
	s.T().Run("管理対象ファイルが正しく保存・取得できる", func(t *testing.T) {
		// テストリポジトリを作成して保存
		repo := s.createTestRepository()
		err := s.repository.Save(s.ctx, repo)
		assert.NoError(t, err)

		// ファイルパスのリスト
		filePaths := []string{
			"docs/README.md",
			"docs/CONTRIBUTING.md",
			"src/main.go",
		}

		// 管理対象ファイルを保存
		err = s.repository.SaveManagedFiles(s.ctx, repo.ID(), filePaths)
		assert.NoError(t, err)

		// 管理対象ファイルを取得
		savedPaths, err := s.repository.GetManagedFiles(s.ctx, repo.ID())

		// 検証
		assert.NoError(t, err)
		assert.ElementsMatch(t, filePaths, savedPaths)
	})

	// テスト：管理対象ファイルが更新されることを確認する
	s.T().Run("管理対象ファイルが更新される", func(t *testing.T) {
		// テストリポジトリを作成して保存
		repo := s.createTestRepository()
		err := s.repository.Save(s.ctx, repo)
		assert.NoError(t, err)

		// 最初のファイルパスのリスト
		filePaths1 := []string{
			"docs/README.md",
			"docs/CONTRIBUTING.md",
		}

		// 管理対象ファイルを保存
		err = s.repository.SaveManagedFiles(s.ctx, repo.ID(), filePaths1)
		assert.NoError(t, err)

		// 更新後のファイルパスのリスト
		filePaths2 := []string{
			"docs/README.md", // 変更なし
			"src/main.go",    // 追加
			// "docs/CONTRIBUTING.md" は削除
		}

		// 管理対象ファイルを更新
		err = s.repository.SaveManagedFiles(s.ctx, repo.ID(), filePaths2)
		assert.NoError(t, err)

		// 管理対象ファイルを取得
		savedPaths, err := s.repository.GetManagedFiles(s.ctx, repo.ID())

		// 検証
		assert.NoError(t, err)
		assert.ElementsMatch(t, filePaths2, savedPaths)
	})

	// テスト：存在しないリポジトリに対する管理対象ファイルの保存はエラーになることを確認する
	s.T().Run("存在しないリポジトリに対する管理対象ファイルの保存はエラーになる", func(t *testing.T) {
		// 存在しないリポジトリID
		nonExistentID := uuid.New().String()
		filePaths := []string{"docs/README.md"}

		// 管理対象ファイルを保存
		err := s.repository.SaveManagedFiles(s.ctx, nonExistentID, filePaths)

		// 検証
		assert.Error(t, err)
	})

	// テスト：リポジトリが存在しても管理対象ファイルが設定されていない場合は空のスライスが返ることを確認する
	s.T().Run("管理対象ファイルが設定されていない場合は空のスライスが返る", func(t *testing.T) {
		// テストリポジトリを作成して保存
		repo := s.createTestRepository()
		err := s.repository.Save(s.ctx, repo)
		assert.NoError(t, err)

		// 管理対象ファイルを取得
		savedPaths, err := s.repository.GetManagedFiles(s.ctx, repo.ID())

		// 検証
		assert.NoError(t, err)
		assert.Empty(t, savedPaths)
	})
}

// TestUpdateAccessToken はUpdateAccessTokenメソッドのテストを行います。
func (s *PostgresRepositoryTestSuite) TestUpdateAccessToken() {
	// テスト：アクセストークンが正しく更新されることを確認する
	s.T().Run("アクセストークンが正しく更新される", func(t *testing.T) {
		// テストリポジトリを作成して保存
		repo := s.createTestRepository()
		err := s.repository.Save(s.ctx, repo)
		assert.NoError(t, err)

		// アクセストークンを更新
		newToken := "updated-token-" + uuid.New().String()
		err = s.repository.UpdateAccessToken(s.ctx, repo.ID(), newToken)
		assert.NoError(t, err)

		// 更新されたリポジトリを取得
		updatedRepo, err := s.repository.FindByID(s.ctx, repo.ID())

		// 検証
		assert.NoError(t, err)
		assert.NotNil(t, updatedRepo)
		assert.Equal(t, newToken, updatedRepo.AccessToken())
	})

	// テスト：存在しないリポジトリに対するアクセストークンの更新はエラーになることを確認する
	s.T().Run("存在しないリポジトリに対するアクセストークンの更新はエラーになる", func(t *testing.T) {
		// 存在しないリポジトリID
		nonExistentID := uuid.New().String()

		// アクセストークンを更新
		err := s.repository.UpdateAccessToken(s.ctx, nonExistentID, "new-token")

		// 検証
		assert.Error(t, err)
	})
}
