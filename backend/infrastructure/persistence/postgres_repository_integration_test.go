package persistence

import (
	"context"
	"fmt"
	"opscore/backend/domain/model"
	"os"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// PostgreSQLIntegrationTestConfig is the configuration for PostgreSQL integration tests
type PostgreSQLIntegrationTestConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
}

// getPostgreSQLConfig gets the PostgreSQL configuration from environment variables
// or returns default values for local testing
func getPostgreSQLConfig() PostgreSQLIntegrationTestConfig {
	return PostgreSQLIntegrationTestConfig{
		Host:     getEnvOrDefault("POSTGRES_HOST", "localhost"),
		Port:     getEnvOrDefault("POSTGRES_PORT", "5432"),
		User:     getEnvOrDefault("POSTGRES_USER", "postgres"),
		Password: getEnvOrDefault("POSTGRES_PASSWORD", "postgres"),
		DBName:   getEnvOrDefault("POSTGRES_DB", "opscore_test"),
	}
}

// getEnvOrDefault gets an environment variable or returns the default value
func getEnvOrDefault(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

// checkDatabaseConnection checks if the database is available
func checkDatabaseConnection(t *testing.T) bool {
	config := getPostgreSQLConfig()
	dsn := fmt.Sprintf("postgresql://%s:%s@%s:%s/postgres",
		config.User, config.Password, config.Host, config.Port)

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	conn, err := pgxpool.New(ctx, dsn)
	if err != nil {
		t.Logf("Could not connect to PostgreSQL: %v", err)
		return false
	}

	err = conn.Ping(ctx)
	conn.Close()

	if err != nil {
		t.Logf("Could not ping PostgreSQL: %v", err)
		return false
	}

	return true
}

// setupPostgreSQLRepository sets up a PostgreSQL repository for testing
// It returns the repository and a cleanup function to close the database connection
func setupPostgreSQLRepository(t *testing.T) (*PostgresRepository, func()) {
	config := getPostgreSQLConfig()

	// テスト用の一意のデータベース名を生成（並列テスト実行のため）
	testDBName := fmt.Sprintf("%s_%s", config.DBName, uuid.New().String()[:8])

	// PostgreSQLに接続（マスターDBに接続してテスト用DBを作成する）
	masterDSN := fmt.Sprintf("postgresql://%s:%s@%s:%s/postgres",
		config.User, config.Password, config.Host, config.Port)

	masterConn, err := pgxpool.New(context.Background(), masterDSN)
	require.NoError(t, err, "Failed to connect to PostgreSQL master database")
	defer masterConn.Close()

	// テスト用のデータベースを作成
	_, err = masterConn.Exec(context.Background(), fmt.Sprintf("CREATE DATABASE %s", testDBName))
	require.NoError(t, err, "Failed to create test database")

	// 作成したテスト用のDBに接続
	testDSN := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s",
		config.User, config.Password, config.Host, config.Port, testDBName)

	testConn, err := pgxpool.New(context.Background(), testDSN)
	require.NoError(t, err, "Failed to connect to test database")

	// マイグレーションを実行
	_, err = testConn.Exec(context.Background(), `
		CREATE TABLE repositories (
			id VARCHAR(36) PRIMARY KEY,
			name VARCHAR(255) NOT NULL,
			url TEXT NOT NULL UNIQUE,
			access_token TEXT,
			created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
		);
		
		CREATE TABLE managed_files (
			repository_id VARCHAR(36) REFERENCES repositories(id),
			file_path TEXT NOT NULL,
			PRIMARY KEY (repository_id, file_path)
		);
	`)
	require.NoError(t, err, "Failed to apply migrations")

	// リポジトリの作成
	repo := &PostgresRepository{db: testConn}

	// クリーンアップ関数を返す
	cleanup := func() {
		// テスト用のDBへの接続を閉じる
		testConn.Close()

		// マスターDBに再接続してテスト用DBを削除
		cleanupConn, err := pgxpool.New(context.Background(), masterDSN)
		if err != nil {
			t.Logf("Failed to connect to master database for cleanup: %v", err)
			return
		}
		defer cleanupConn.Close()

		// 他の接続を強制的に切断
		_, err = cleanupConn.Exec(context.Background(), fmt.Sprintf("DROP DATABASE IF EXISTS %s WITH (FORCE)", testDBName))
		if err != nil {
			t.Logf("Failed to drop test database: %v", err)
		}
	}

	return repo, cleanup
}

// TestPostgresRepositoryIntegration is an integration test for the PostgresRepository
func TestPostgresRepositoryIntegration(t *testing.T) {
	// Skip test if PostgreSQL is not available
	if !checkDatabaseConnection(t) {
		t.Skip("Skipping PostgreSQL integration test - database is not available")
		return
	}

	// PostgreSQLリポジトリのセットアップ
	repo, cleanup := setupPostgreSQLRepository(t)
	defer cleanup()

	ctx := context.Background()

	// テスト: リポジトリの保存と取得
	t.Run("Save and FindByID", func(t *testing.T) {
		// テストデータの準備
		repoID := uuid.New().String()
		repoName := "test-repo"
		repoURL := "https://github.com/example/test-repo"
		accessToken := "test-token"

		testRepo := model.NewRepository(repoID, repoName, repoURL, accessToken)

		// リポジトリの保存
		err := repo.Save(ctx, testRepo)
		require.NoError(t, err)

		// IDで取得
		retrieved, err := repo.FindByID(ctx, repoID)
		require.NoError(t, err)
		require.NotNil(t, retrieved)

		// 取得したリポジトリの検証
		assert.Equal(t, repoID, retrieved.ID())
		assert.Equal(t, repoName, retrieved.Name())
		assert.Equal(t, repoURL, retrieved.URL())
		// Note: アクセストークンはエンティティ外部からはアクセスできないため検証できない
	})

	// テスト: URLでリポジトリを検索
	t.Run("FindByURL", func(t *testing.T) {
		// テストデータの準備
		repoID := uuid.New().String()
		repoName := "url-repo"
		repoURL := "https://github.com/example/url-repo"
		accessToken := "url-token"

		testRepo := model.NewRepository(repoID, repoName, repoURL, accessToken)

		// リポジトリの保存
		err := repo.Save(ctx, testRepo)
		require.NoError(t, err)

		// URLで取得
		retrieved, err := repo.FindByURL(ctx, repoURL)
		require.NoError(t, err)
		require.NotNil(t, retrieved)

		// 取得したリポジトリの検証
		assert.Equal(t, repoID, retrieved.ID())
		assert.Equal(t, repoName, retrieved.Name())
		assert.Equal(t, repoURL, retrieved.URL())
	})

	// テスト: 全リポジトリの取得
	t.Run("FindAll", func(t *testing.T) {
		// 追加のテストデータ
		newRepo := model.NewRepository(
			uuid.New().String(),
			"another-repo",
			"https://github.com/example/another-repo",
			"another-token",
		)

		// 保存
		err := repo.Save(ctx, newRepo)
		require.NoError(t, err)

		// 全リポジトリ取得
		repos, err := repo.FindAll(ctx)
		require.NoError(t, err)

		// 少なくとも2つ以上のリポジトリが存在すること
		assert.GreaterOrEqual(t, len(repos), 2)

		// 保存したリポジトリが含まれていることを確認
		found := false
		for _, r := range repos {
			if r.ID() == newRepo.ID() {
				found = true
				break
			}
		}
		assert.True(t, found)
	})

	// テスト: アクセストークンの更新
	t.Run("UpdateAccessToken", func(t *testing.T) {
		// テストデータの準備
		repoID := uuid.New().String()
		repoName := "token-update-repo"
		repoURL := "https://github.com/example/token-update-repo"
		initialToken := "initial-token"

		testRepo := model.NewRepository(repoID, repoName, repoURL, initialToken)

		// リポジトリの保存
		err := repo.Save(ctx, testRepo)
		require.NoError(t, err)

		// アクセストークンの更新
		updatedToken := "updated-token"
		err = repo.UpdateAccessToken(ctx, repoID, updatedToken)
		require.NoError(t, err)

		// リポジトリの再取得
		retrieved, err := repo.FindByID(ctx, repoID)
		require.NoError(t, err)
		require.NotNil(t, retrieved)

		// 基本情報が変わっていないことを確認
		assert.Equal(t, repoID, retrieved.ID())
		assert.Equal(t, repoName, retrieved.Name())
		assert.Equal(t, repoURL, retrieved.URL())
		// トークンが更新されていることは直接検証できないが、エラーなく処理が完了したことで更新は成功したと判断
	})

	// テスト: マネージドファイルの保存と取得
	t.Run("SaveManagedFiles and GetManagedFiles", func(t *testing.T) {
		// テストデータの準備
		repoID := uuid.New().String()
		repoName := "files-repo"
		repoURL := "https://github.com/example/files-repo"
		accessToken := "files-token"

		testRepo := model.NewRepository(repoID, repoName, repoURL, accessToken)

		// リポジトリの保存
		err := repo.Save(ctx, testRepo)
		require.NoError(t, err)

		// マネージドファイルの保存
		filePaths := []string{
			"README.md",
			"docs/adr/0001-test-adr.md",
			"src/main.go",
		}
		err = repo.SaveManagedFiles(ctx, repoID, filePaths)
		require.NoError(t, err)

		// マネージドファイルの取得
		retrievedFiles, err := repo.GetManagedFiles(ctx, repoID)
		require.NoError(t, err)

		// 件数の確認
		assert.Equal(t, len(filePaths), len(retrievedFiles))

		// 全ファイルパスが含まれていることを確認
		for _, path := range filePaths {
			assert.Contains(t, retrievedFiles, path)
		}

		// マネージドファイルの上書き
		newFilePaths := []string{
			"README.md",   // 既存の一つを残す
			"new-file.md", // 新しいファイル
		}
		err = repo.SaveManagedFiles(ctx, repoID, newFilePaths)
		require.NoError(t, err)

		// 更新後のマネージドファイルを取得
		updatedFiles, err := repo.GetManagedFiles(ctx, repoID)
		require.NoError(t, err)

		// 更新後の件数と内容を確認
		assert.Equal(t, len(newFilePaths), len(updatedFiles))
		for _, path := range newFilePaths {
			assert.Contains(t, updatedFiles, path)
		}
		// 削除されたファイルが含まれていないことを確認
		assert.NotContains(t, updatedFiles, "docs/adr/0001-test-adr.md")
	})

	// テスト: 存在しないリポジトリの処理
	t.Run("Non-existent repository", func(t *testing.T) {
		nonExistentID := "non-existent-id"

		// FindByID
		retrieved, err := repo.FindByID(ctx, nonExistentID)
		require.NoError(t, err)  // エラーは返さない
		assert.Nil(t, retrieved) // 結果はnilになる

		// UpdateAccessToken
		err = repo.UpdateAccessToken(ctx, nonExistentID, "token")
		require.Error(t, err) // エラーが発生するべき

		// SaveManagedFiles
		err = repo.SaveManagedFiles(ctx, nonExistentID, []string{"file.md"})
		require.Error(t, err) // エラーが発生するべき

		// GetManagedFiles
		files, err := repo.GetManagedFiles(ctx, nonExistentID)
		require.NoError(t, err) // エラーは返さない
		assert.Empty(t, files)  // 空のスライスが返るべき
	})
}
