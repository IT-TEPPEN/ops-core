package persistence

import (
	"context"
	"crypto/rand"
	"testing"

	"opscore/backend/internal/git_repository/domain/entity"
	"opscore/backend/internal/git_repository/infrastructure/encryption"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestAccessTokenEncryption verifies that access tokens are encrypted in the database
func TestAccessTokenEncryption(t *testing.T) {
	// Skip test if PostgreSQL is not available
	if !checkDatabaseConnection(t) {
		t.Skip("Skipping PostgreSQL integration test - database is not available")
		return
	}

	// PostgreSQLリポジトリのセットアップ
	repo, cleanup := setupPostgreSQLRepository(t)
	defer cleanup()

	ctx := context.Background()

	t.Run("Access tokens are encrypted in database", func(t *testing.T) {
		// テストデータの準備
		repoID := uuid.New().String()
		repoName := "encryption-test-repo"
		repoURL := "https://github.com/example/encryption-test"
		plainTextToken := "ghp_1234567890abcdefghijklmnopqrstuvwxyz"

		testRepo := entity.NewRepository(repoID, repoName, repoURL, plainTextToken)

		// リポジトリの保存
		err := repo.Save(ctx, testRepo)
		require.NoError(t, err)

		// データベースから直接取得して、暗号化されていることを確認
		var storedToken string
		err = repo.db.QueryRow(ctx, "SELECT access_token FROM repositories WHERE id = $1", repoID).Scan(&storedToken)
		require.NoError(t, err)

		// 暗号化されているため、平文トークンと異なるはず
		assert.NotEqual(t, plainTextToken, storedToken, "Token should be encrypted in database")
		assert.NotEmpty(t, storedToken, "Encrypted token should not be empty")

		// Base64エンコードされた暗号文であることを確認（簡易チェック）
		assert.Greater(t, len(storedToken), len(plainTextToken), "Encrypted token should be longer than plaintext")

		// リポジトリ経由で取得すると、復号化された値が返ることを確認
		retrieved, err := repo.FindByID(ctx, repoID)
		require.NoError(t, err)
		require.NotNil(t, retrieved)

		// 復号化されて元の平文トークンが返ることを確認
		assert.Equal(t, plainTextToken, retrieved.AccessToken(), "Retrieved token should be decrypted")
	})

	t.Run("Empty access token handling", func(t *testing.T) {
		// テストデータの準備（空のトークン）
		repoID := uuid.New().String()
		repoName := "empty-token-repo"
		repoURL := "https://github.com/example/empty-token"
		emptyToken := ""

		testRepo := entity.NewRepository(repoID, repoName, repoURL, emptyToken)

		// リポジトリの保存
		err := repo.Save(ctx, testRepo)
		require.NoError(t, err)

		// リポジトリ経由で取得
		retrieved, err := repo.FindByID(ctx, repoID)
		require.NoError(t, err)
		require.NotNil(t, retrieved)

		// 空のトークンが正しく処理されることを確認
		assert.Equal(t, emptyToken, retrieved.AccessToken())
	})

	t.Run("Different tokens produce different ciphertexts", func(t *testing.T) {
		// 同じ平文トークンでも、毎回異なる暗号文が生成されることを確認
		plainTextToken := "ghp_sametoken1234567890abcdefghijk"

		// 1つ目のリポジトリ
		repoID1 := uuid.New().String()
		testRepo1 := entity.NewRepository(repoID1, "repo1", "https://github.com/example/repo1", plainTextToken)
		err := repo.Save(ctx, testRepo1)
		require.NoError(t, err)

		// 2つ目のリポジトリ
		repoID2 := uuid.New().String()
		testRepo2 := entity.NewRepository(repoID2, "repo2", "https://github.com/example/repo2", plainTextToken)
		err = repo.Save(ctx, testRepo2)
		require.NoError(t, err)

		// データベースから直接取得
		var storedToken1, storedToken2 string
		err = repo.db.QueryRow(ctx, "SELECT access_token FROM repositories WHERE id = $1", repoID1).Scan(&storedToken1)
		require.NoError(t, err)
		err = repo.db.QueryRow(ctx, "SELECT access_token FROM repositories WHERE id = $2", repoID2).Scan(&storedToken2)
		require.NoError(t, err)

		// 同じ平文でも暗号文は異なる（nonceが異なるため）
		assert.NotEqual(t, storedToken1, storedToken2, "Different encryptions of same plaintext should produce different ciphertexts")

		// 両方とも復号化すると同じ平文が得られる
		retrieved1, _ := repo.FindByID(ctx, repoID1)
		retrieved2, _ := repo.FindByID(ctx, repoID2)
		assert.Equal(t, plainTextToken, retrieved1.AccessToken())
		assert.Equal(t, plainTextToken, retrieved2.AccessToken())
	})

	t.Run("UpdateAccessToken encrypts new token", func(t *testing.T) {
		// テストデータの準備
		repoID := uuid.New().String()
		originalToken := "ghp_original_token_1234567890"
		testRepo := entity.NewRepository(repoID, "update-test", "https://github.com/example/update-test", originalToken)

		err := repo.Save(ctx, testRepo)
		require.NoError(t, err)

		// トークンを更新
		newToken := "ghp_updated_token_0987654321"
		err = repo.UpdateAccessToken(ctx, repoID, newToken)
		require.NoError(t, err)

		// データベースから直接取得して暗号化されていることを確認
		var storedToken string
		err = repo.db.QueryRow(ctx, "SELECT access_token FROM repositories WHERE id = $1", repoID).Scan(&storedToken)
		require.NoError(t, err)

		// 暗号化されているため、平文と異なるはず
		assert.NotEqual(t, newToken, storedToken)

		// リポジトリ経由で取得すると、復号化された新しいトークンが返る
		retrieved, err := repo.FindByID(ctx, repoID)
		require.NoError(t, err)
		assert.Equal(t, newToken, retrieved.AccessToken())
	})

	t.Run("FindAll returns decrypted tokens", func(t *testing.T) {
		// 複数のリポジトリを保存
		tokens := map[string]string{
			uuid.New().String(): "ghp_token1_abcdefgh",
			uuid.New().String(): "ghp_token2_ijklmnop",
			uuid.New().String(): "ghp_token3_qrstuvwx",
		}

		for id, token := range tokens {
			testRepo := entity.NewRepository(id, "repo-"+id[:8], "https://github.com/example/"+id[:8], token)
			err := repo.Save(ctx, testRepo)
			require.NoError(t, err)
		}

		// FindAllで取得
		allRepos, err := repo.FindAll(ctx)
		require.NoError(t, err)

		// すべてのトークンが復号化されていることを確認
		foundTokens := make(map[string]string)
		for _, r := range allRepos {
			if token, ok := tokens[r.ID()]; ok {
				foundTokens[r.ID()] = r.AccessToken()
				assert.Equal(t, token, r.AccessToken(), "Token should be decrypted in FindAll")
			}
		}

		assert.Equal(t, len(tokens), len(foundTokens), "All saved tokens should be retrieved")
	})

	t.Run("FindByURL returns decrypted token", func(t *testing.T) {
		repoID := uuid.New().String()
		repoURL := "https://github.com/example/url-encryption-test"
		plainToken := "ghp_url_test_token_123456"

		testRepo := entity.NewRepository(repoID, "url-test", repoURL, plainToken)
		err := repo.Save(ctx, testRepo)
		require.NoError(t, err)

		// URLで検索
		retrieved, err := repo.FindByURL(ctx, repoURL)
		require.NoError(t, err)
		require.NotNil(t, retrieved)

		// トークンが復号化されていることを確認
		assert.Equal(t, plainToken, retrieved.AccessToken())
	})
}

// TestEncryptionKeyRotation tests the scenario of encryption key rotation
func TestEncryptionKeyRotation(t *testing.T) {
	// Skip test if PostgreSQL is not available
	if !checkDatabaseConnection(t) {
		t.Skip("Skipping PostgreSQL integration test - database is not available")
		return
	}

	ctx := context.Background()

	t.Run("Cannot decrypt with different key", func(t *testing.T) {
		// Setup repository with first key
		repo1, cleanup1 := setupPostgreSQLRepository(t)
		defer cleanup1()

		// Save a repository with a token using first key
		repoID := uuid.New().String()
		plainToken := "ghp_rotation_test_token"
		testRepo := entity.NewRepository(repoID, "rotation-test", "https://github.com/example/rotation", plainToken)
		err := repo1.Save(ctx, testRepo)
		require.NoError(t, err)

		// Get the encrypted token from database
		var encryptedToken string
		err = repo1.db.QueryRow(ctx, "SELECT access_token FROM repositories WHERE id = $1", repoID).Scan(&encryptedToken)
		require.NoError(t, err)

		// Create a new encryptor with a different key
		differentKey := make([]byte, 32)
		_, err = rand.Read(differentKey)
		require.NoError(t, err)

		differentEncryptor, err := encryption.NewEncryptor(differentKey)
		require.NoError(t, err)

		// Try to decrypt with the different key - should fail
		_, err = differentEncryptor.Decrypt(encryptedToken)
		assert.Error(t, err, "Decryption with different key should fail")
	})
}
