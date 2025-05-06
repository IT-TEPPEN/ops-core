package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"opscore/backend/domain/model"
	"opscore/backend/infrastructure/git"
	"opscore/backend/usecases/repository"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// IntegrationTestLogger is a simple logger implementation for integration tests
type IntegrationTestLogger struct{}

func (l *IntegrationTestLogger) Info(msg string, args ...any) {
	log.Printf("[INFO] %s %v", msg, args)
}

func (l *IntegrationTestLogger) Error(msg string, args ...any) {
	log.Printf("[ERROR] %s %v", msg, args)
}

func (l *IntegrationTestLogger) Debug(msg string, args ...any) {
	log.Printf("[DEBUG] %s %v", msg, args)
}

func (l *IntegrationTestLogger) Warn(msg string, args ...any) {
	log.Printf("[WARN] %s %v", msg, args)
}

// setupIntegrationTest creates a test environment with actual implementations
// rather than mocks for integration testing
func setupIntegrationTest() (*RepositoryHandler, *gin.Engine, *httptest.ResponseRecorder) {
	// テスト用にGinをリリースモードに設定（ログを抑制）
	gin.SetMode(gin.TestMode)

	// 実際のロガーを使用
	logger := &IntegrationTestLogger{}

	// インメモリのリポジトリを作成（統合テストなのでモックではなく実際の実装を使用）
	repo := NewInMemoryRepository()

	// テスト後に削除するようにする
	// defer os.RemoveAll(tempDir) // テスト時に削除せず、デバッグのために残しておく

	// 実際のGitマネージャーを作成
	// CLIの実装を避けてモックや簡易的な実装を使うことでテスト環境の依存を減らす
	gitManager := git.NewMockGitManager() // 本来はCLIではなくAPIを使うとよい

	// 実際のUseCaseを作成
	useCase := repository.NewRepositoryUseCase(repo, gitManager)

	// テスト対象ハンドラーを作成（実際の実装を使用）
	handler := NewRepositoryHandler(useCase, logger)

	// Ginのルーターとレスポンスレコーダーを設定
	router := gin.New()
	router.Use(func(c *gin.Context) {
		// リクエストIDをコンテキストに設定するミドルウェアをモック
		c.Set("request_id", "test-integration-request-id")
		c.Next()
	})

	// レスポンスレコーダーを作成
	rec := httptest.NewRecorder()

	return handler, router, rec
}

// TestRegisterRepositoryIntegration tests the RegisterRepository handler with actual use cases
func TestRegisterRepositoryIntegration(t *testing.T) {
	// テスト用のセットアップ
	handler, router, rec := setupIntegrationTest()

	// ルーターの設定
	router.POST("/repositories", handler.RegisterRepository)

	// テストデータ
	repoURL := "https://github.com/example/test-repo"
	accessToken := "test-token"
	requestBody := RegisterRepositoryRequest{
		URL:         repoURL,
		AccessToken: accessToken,
	}
	jsonBody, _ := json.Marshal(requestBody)

	// リクエスト実行
	req, _ := http.NewRequest("POST", "/repositories", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(rec, req)

	// 検証
	assert.Equal(t, http.StatusCreated, rec.Code)

	var response RepositoryResponse
	err := json.Unmarshal(rec.Body.Bytes(), &response)
	require.NoError(t, err)

	// レスポンスの検証
	assert.NotEmpty(t, response.ID)
	assert.Equal(t, "test-repo", response.Name) // URLから抽出されるはず
	assert.Equal(t, repoURL, response.URL)

	// 登録したリポジトリをGetRepositoryで取得できることを確認
	t.Run("RegisteredRepositoryCanBeRetrieved", func(t *testing.T) {
		getRouter := gin.New()
		getRec := httptest.NewRecorder()

		// ルーターの設定
		getRouter.GET("/repositories/:repoId", handler.GetRepository)

		// リクエスト実行
		req, _ := http.NewRequest("GET", "/repositories/"+response.ID, nil)
		getRouter.ServeHTTP(getRec, req)

		// 検証
		assert.Equal(t, http.StatusOK, getRec.Code)

		var getResponse RepositoryResponse
		err := json.Unmarshal(getRec.Body.Bytes(), &getResponse)
		require.NoError(t, err)

		// 元の登録内容と一致することを確認
		assert.Equal(t, response.ID, getResponse.ID)
		assert.Equal(t, response.Name, getResponse.Name)
		assert.Equal(t, response.URL, getResponse.URL)
	})
}

// TestRepositoryEndToEndFlow tests a complete repository workflow
func TestRepositoryEndToEndFlow(t *testing.T) {
	// テスト用のセットアップ
	handler, router, _ := setupIntegrationTest()

	// Step 1: リポジトリを登録
	// ルーターの設定
	router.POST("/repositories", handler.RegisterRepository)
	router.GET("/repositories/:repoId", handler.GetRepository)
	router.GET("/repositories", handler.ListRepositories)
	router.PUT("/repositories/:repoId/token", handler.UpdateAccessToken)

	// リポジトリ登録
	repoURL := "https://github.com/example/test-repo-flow"
	accessToken := "test-token-flow"
	requestBody := RegisterRepositoryRequest{
		URL:         repoURL,
		AccessToken: accessToken,
	}
	jsonBody, _ := json.Marshal(requestBody)

	rec := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/repositories", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusCreated, rec.Code)

	var repoResponse RepositoryResponse
	err := json.Unmarshal(rec.Body.Bytes(), &repoResponse)
	require.NoError(t, err)
	repoID := repoResponse.ID

	// Step 2: リポジトリの取得
	rec = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/repositories/"+repoID, nil)
	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)

	var getResponse RepositoryResponse
	err = json.Unmarshal(rec.Body.Bytes(), &getResponse)
	require.NoError(t, err)
	assert.Equal(t, repoID, getResponse.ID)

	// Step 3: リポジトリの一覧取得
	rec = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/repositories", nil)
	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)

	var listResponse ListRepositoriesResponse
	err = json.Unmarshal(rec.Body.Bytes(), &listResponse)
	require.NoError(t, err)
	assert.GreaterOrEqual(t, len(listResponse.Repositories), 1)

	found := false
	for _, repo := range listResponse.Repositories {
		if repo.ID == repoID {
			found = true
			break
		}
	}
	assert.True(t, found, "登録したリポジトリが一覧に含まれていない")

	// Step 4: アクセストークンの更新
	newToken := "updated-token"
	tokenRequest := UpdateAccessTokenRequest{
		AccessToken: newToken,
	}
	tokenJSON, _ := json.Marshal(tokenRequest)

	rec = httptest.NewRecorder()
	req, _ = http.NewRequest("PUT", "/repositories/"+repoID+"/token", bytes.NewBuffer(tokenJSON))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)

	// Step 5: 更新後のリポジトリ取得で内部的にアクセストークンが更新されていることを確認
	// Noteː トークンは外部に公開されないので直接確認できないが、取得は成功するはず
	rec = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/repositories/"+repoID, nil)
	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
}

// inMemoryRepository provides an in-memory implementation of the repository.Repository interface for testing
type inMemoryRepository struct {
	repositories map[string]model.Repository // Map of repository ID to Repository
	managedFiles map[string][]string         // Map of repository ID to managed file paths
}

func NewInMemoryRepository() *inMemoryRepository {
	return &inMemoryRepository{
		repositories: make(map[string]model.Repository),
		managedFiles: make(map[string][]string),
	}
}

func (r *inMemoryRepository) Save(ctx context.Context, repo model.Repository) error {
	r.repositories[repo.ID()] = repo
	return nil
}

func (r *inMemoryRepository) FindByID(ctx context.Context, id string) (model.Repository, error) {
	repo, exists := r.repositories[id]
	if !exists {
		return nil, nil // Not found
	}
	return repo, nil
}

func (r *inMemoryRepository) FindByURL(ctx context.Context, url string) (model.Repository, error) {
	for _, repo := range r.repositories {
		if repo.URL() == url {
			return repo, nil
		}
	}
	return nil, nil // Not found
}

func (r *inMemoryRepository) FindAll(ctx context.Context) ([]model.Repository, error) {
	repos := make([]model.Repository, 0, len(r.repositories))
	for _, repo := range r.repositories {
		repos = append(repos, repo)
	}
	return repos, nil
}

func (r *inMemoryRepository) UpdateAccessToken(ctx context.Context, id string, accessToken string) error {
	repo, exists := r.repositories[id]
	if !exists {
		return repository.ErrRepositoryNotFound
	}

	// We need to use SetAccessToken instead of creating a new repository
	repo.SetAccessToken(accessToken)
	return nil
}

func (r *inMemoryRepository) SaveManagedFiles(ctx context.Context, repoID string, filePaths []string) error {
	_, exists := r.repositories[repoID]
	if !exists {
		return repository.ErrRepositoryNotFound
	}
	r.managedFiles[repoID] = filePaths
	return nil
}

func (r *inMemoryRepository) GetManagedFiles(ctx context.Context, repoID string) ([]string, error) {
	filePaths, exists := r.managedFiles[repoID]
	if !exists {
		return []string{}, nil // Empty slice for no files
	}
	return filePaths, nil
}
