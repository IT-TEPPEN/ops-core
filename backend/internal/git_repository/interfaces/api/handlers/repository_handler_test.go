package handlers

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	apperror "opscore/backend/internal/git_repository/application/error"
	"opscore/backend/internal/git_repository/application/usecase"
	"opscore/backend/internal/git_repository/domain/entity"
	"opscore/backend/internal/git_repository/interfaces/api/schema"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func setupTest() (*repository.MockRepositoryUseCase, *MockLogger, *RepositoryHandler, *gin.Engine, *httptest.ResponseRecorder) {
	// テスト用にGinをリリースモードに設定（ログを抑制）
	gin.SetMode(gin.TestMode)

	// モックの準備
	mockUseCase := new(repository.MockRepositoryUseCase)
	mockLogger := new(MockLogger)

	// Set up the logger mock to handle any number of arguments
	mockLogger.On("Info", mock.AnythingOfType("string"), mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Maybe()
	mockLogger.On("Error", mock.AnythingOfType("string"), mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Maybe()
	mockLogger.On("Warn", mock.AnythingOfType("string"), mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Maybe()
	mockLogger.On("Debug", mock.AnythingOfType("string"), mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Maybe()

	// テスト対象ハンドラーを作成
	handler := NewRepositoryHandler(mockUseCase, mockLogger)

	// Ginのルーターとレスポンスレコーダーを設定
	router := gin.New() // ミドルウェアなしで新しいエンジンを作成
	router.Use(func(c *gin.Context) {
		// リクエストIDをコンテキストに設定するミドルウェアをモック
		c.Set("request_id", "test-request-id")
		c.Next()
	})

	// レスポンスレコーダーを作成
	rec := httptest.NewRecorder()

	return mockUseCase, mockLogger, handler, router, rec
}

func TestRegisterRepository(t *testing.T) {
	t.Run("正常に登録できる場合", func(t *testing.T) {
		// テスト用のセットアップ
		mockUseCase, _, handler, router, rec := setupTest()

		// テストデータ
		repoURL := "https://github.com/example/test-repo"
		accessToken := "test-token"
		requestBody := schema.RegisterRepositoryRequest{
			URL:         repoURL,
			AccessToken: accessToken,
		}
		jsonBody, _ := json.Marshal(requestBody)

		// Using matchers for better flexibility
		contextMatcher := mock.MatchedBy(func(ctx interface{}) bool { return true })
		urlMatcher := mock.MatchedBy(func(url string) bool { return url == repoURL })
		tokenMatcher := mock.MatchedBy(func(token string) bool { return token == accessToken })

		// モックの振る舞いを定義
		createdRepo := entity.ReconstructRepository(
			uuid.NewString(),
			"test-repo",
			repoURL,
			accessToken,
			time.Now(),
			time.Now(),
		)
		mockUseCase.On("Register", contextMatcher, urlMatcher, tokenMatcher).Return(createdRepo, nil)

		// ルーターの設定
		router.POST("/repositories", handler.RegisterRepository)

		// リクエスト実行
		req, _ := http.NewRequest("POST", "/repositories", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(rec, req)

		// 検証
		assert.Equal(t, http.StatusCreated, rec.Code)

		var response schema.RepositoryResponse
		err := json.Unmarshal(rec.Body.Bytes(), &response)
		require.NoError(t, err)
		assert.Equal(t, createdRepo.ID(), response.ID)
		assert.Equal(t, createdRepo.Name(), response.Name)
		assert.Equal(t, createdRepo.URL(), response.URL)

		// モックの呼び出しを検証
		mockUseCase.AssertExpectations(t)
	})

	t.Run("不正なリクエスト本文の場合", func(t *testing.T) {
		// テスト用のセットアップ
		_, _, handler, router, rec := setupTest()

		// 不正なJSONデータ
		invalidJSON := `{"url": }`

		// ルーターの設定
		router.POST("/repositories", handler.RegisterRepository)

		// リクエスト実行
		req, _ := http.NewRequest("POST", "/repositories", bytes.NewBufferString(invalidJSON))
		req.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(rec, req)

		// 検証
		assert.Equal(t, http.StatusBadRequest, rec.Code)

		var response schema.ErrorResponse
		err := json.Unmarshal(rec.Body.Bytes(), &response)
		require.NoError(t, err)
		assert.Equal(t, "INVALID_REQUEST", response.Code)
		assert.Contains(t, response.Message, "Invalid request body")
	})

	t.Run("URL必須パラメータが欠けている場合", func(t *testing.T) {
		// テスト用のセットアップ
		_, _, handler, router, rec := setupTest()

		// URLが欠けている
		requestBody := schema.RegisterRepositoryRequest{
			AccessToken: "test-token",
		}
		jsonBody, _ := json.Marshal(requestBody)

		// ルーターの設定
		router.POST("/repositories", handler.RegisterRepository)

		// リクエスト実行
		req, _ := http.NewRequest("POST", "/repositories", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(rec, req)

		// 検証
		assert.Equal(t, http.StatusBadRequest, rec.Code)

		var response schema.ErrorResponse
		err := json.Unmarshal(rec.Body.Bytes(), &response)
		require.NoError(t, err)
		assert.Equal(t, "INVALID_REQUEST", response.Code)
		assert.Contains(t, response.Message, "Invalid request body")
	})

	t.Run("リポジトリがすでに存在する場合", func(t *testing.T) {
		// テスト用のセットアップ
		mockUseCase, _, handler, router, rec := setupTest()

		// テストデータ
		repoURL := "https://github.com/example/existing-repo"
		accessToken := "test-token"
		requestBody := schema.RegisterRepositoryRequest{
			URL:         repoURL,
			AccessToken: accessToken,
		}
		jsonBody, _ := json.Marshal(requestBody)

		// モックの振る舞いを定義
		mockUseCase.On("Register", mock.Anything, repoURL, accessToken).Return(nil, &apperror.ConflictError{
			Code:         apperror.CodeResourceConflict,
			ResourceType: "Repository",
			Identifier:   repoURL,
			Reason:       "repository with this URL already exists",
		})

		// ルーターの設定
		router.POST("/repositories", handler.RegisterRepository)

		// リクエスト実行
		req, _ := http.NewRequest("POST", "/repositories", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(rec, req)

		// 検証
		assert.Equal(t, http.StatusConflict, rec.Code)

		var response schema.ErrorResponse
		err := json.Unmarshal(rec.Body.Bytes(), &response)
		require.NoError(t, err)
		assert.Equal(t, "CONFLICT", response.Code)

		// モックの呼び出しを検証
		mockUseCase.AssertExpectations(t)
	})

	t.Run("無効なリポジトリURLの場合", func(t *testing.T) {
		// テスト用のセットアップ
		mockUseCase, _, handler, router, rec := setupTest()

		// テストデータ
		repoURL := "invalid-url"
		accessToken := "test-token"
		requestBody := schema.RegisterRepositoryRequest{
			URL:         repoURL,
			AccessToken: accessToken,
		}
		jsonBody, _ := json.Marshal(requestBody)

		// モックの振る舞いを定義
		mockUseCase.On("Register", mock.Anything, repoURL, accessToken).Return(nil, &apperror.ValidationFailedError{
			Code: apperror.CodeValidationFailed,
			Errors: []apperror.FieldError{
				{Field: "url", Message: "invalid repository URL format"},
			},
		})

		// ルーターの設定
		router.POST("/repositories", handler.RegisterRepository)

		// リクエスト実行
		req, _ := http.NewRequest("POST", "/repositories", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(rec, req)

		// 検証
		assert.Equal(t, http.StatusBadRequest, rec.Code)

		var response schema.ErrorResponse
		err := json.Unmarshal(rec.Body.Bytes(), &response)
		require.NoError(t, err)
		assert.Equal(t, "BAD_REQUEST", response.Code)

		// モックの呼び出しを検証
		mockUseCase.AssertExpectations(t)
	})

	t.Run("内部エラーが発生した場合", func(t *testing.T) {
		// テスト用のセットアップ
		mockUseCase, _, handler, router, rec := setupTest()

		// テストデータ
		repoURL := "https://github.com/example/test-repo"
		accessToken := "test-token"
		requestBody := schema.RegisterRepositoryRequest{
			URL:         repoURL,
			AccessToken: accessToken,
		}
		jsonBody, _ := json.Marshal(requestBody)

		// モックの振る舞いを定義
		mockUseCase.On("Register", mock.Anything, repoURL, accessToken).Return(nil, errors.New("internal error"))

		// ルーターの設定
		router.POST("/repositories", handler.RegisterRepository)

		// リクエスト実行
		req, _ := http.NewRequest("POST", "/repositories", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(rec, req)

		// 検証
		assert.Equal(t, http.StatusInternalServerError, rec.Code)

		var response schema.ErrorResponse
		err := json.Unmarshal(rec.Body.Bytes(), &response)
		require.NoError(t, err)
		assert.Equal(t, "INTERNAL_SERVER_ERROR", response.Code)

		// モックの呼び出しを検証
		mockUseCase.AssertExpectations(t)
	})
}

func TestGetRepository(t *testing.T) {
	t.Run("正常に取得できる場合", func(t *testing.T) {
		// テスト用のセットアップ
		mockUseCase, _, handler, router, rec := setupTest()

		// テストデータ
		repoID := uuid.NewString()
		repo := entity.NewRepository(
			repoID,
			"test-repo",
			"https://github.com/example/test-repo",
			"test-token",
		)

		// モックの振る舞いを定義
		mockUseCase.On("GetRepository", mock.Anything, repoID).Return(repo, nil)

		// ルーターの設定
		router.GET("/repositories/:repoId", handler.GetRepository)

		// リクエスト実行
		req, _ := http.NewRequest("GET", "/repositories/"+repoID, nil)
		router.ServeHTTP(rec, req)

		// 検証
		assert.Equal(t, http.StatusOK, rec.Code)

		var response schema.RepositoryResponse
		err := json.Unmarshal(rec.Body.Bytes(), &response)
		require.NoError(t, err)
		assert.Equal(t, repo.ID(), response.ID)
		assert.Equal(t, repo.Name(), response.Name)
		assert.Equal(t, repo.URL(), response.URL)

		// モックの呼び出しを検証
		mockUseCase.AssertExpectations(t)
	})

	t.Run("リポジトリIDが欠けている場合", func(t *testing.T) {
		// テスト用のセットアップ
		_, _, handler, router, rec := setupTest()

		// ルーターの設定
		router.GET("/repositories/:repoId", handler.GetRepository)

		// リクエスト実行（空のIDパラメータを使用）
		req, _ := http.NewRequest("GET", "/repositories/", nil)
		router.ServeHTTP(rec, req)

		// 検証
		assert.Equal(t, http.StatusNotFound, rec.Code) // Ginルーターのデフォルト挙動でNotFoundになる
	})

	t.Run("リポジトリが存在しない場合", func(t *testing.T) {
		// テスト用のセットアップ
		mockUseCase, _, handler, router, rec := setupTest()

		// テストデータ
		repoID := uuid.NewString()

		// モックの振る舞いを定義
		mockUseCase.On("GetRepository", mock.Anything, repoID).Return(nil, &apperror.NotFoundError{
			Code:         apperror.CodeResourceNotFound,
			ResourceType: "Repository",
			ResourceID:   repoID,
		})

		// ルーターの設定
		router.GET("/repositories/:repoId", handler.GetRepository)

		// リクエスト実行
		req, _ := http.NewRequest("GET", "/repositories/"+repoID, nil)
		router.ServeHTTP(rec, req)

		// 検証
		assert.Equal(t, http.StatusNotFound, rec.Code)

		var response schema.ErrorResponse
		err := json.Unmarshal(rec.Body.Bytes(), &response)
		require.NoError(t, err)
		assert.Equal(t, "NOT_FOUND", response.Code)

		// モックの呼び出しを検証
		mockUseCase.AssertExpectations(t)
	})

	t.Run("内部エラーが発生した場合", func(t *testing.T) {
		// テスト用のセットアップ
		mockUseCase, _, handler, router, rec := setupTest()

		// テストデータ
		repoID := uuid.NewString()

		// モックの振る舞いを定義
		mockUseCase.On("GetRepository", mock.Anything, repoID).Return(nil, errors.New("internal error"))

		// ルーターの設定
		router.GET("/repositories/:repoId", handler.GetRepository)

		// リクエスト実行
		req, _ := http.NewRequest("GET", "/repositories/"+repoID, nil)
		router.ServeHTTP(rec, req)

		// 検証
		assert.Equal(t, http.StatusInternalServerError, rec.Code)

		var response schema.ErrorResponse
		err := json.Unmarshal(rec.Body.Bytes(), &response)
		require.NoError(t, err)
		assert.Equal(t, "INTERNAL_SERVER_ERROR", response.Code)

		// モックの呼び出しを検証
		mockUseCase.AssertExpectations(t)
	})
}

func TestListRepositories(t *testing.T) {
	t.Run("正常にリポジトリ一覧を取得できる場合", func(t *testing.T) {
		// テスト用のセットアップ
		mockUseCase, _, handler, router, rec := setupTest()

		// テストデータ
		repo1 := entity.NewRepository(
			uuid.NewString(),
			"repo1",
			"https://github.com/example/repo1",
			"token1",
		)
		repo2 := entity.NewRepository(
			uuid.NewString(),
			"repo2",
			"https://github.com/example/repo2",
			"token2",
		)
		repos := []entity.Repository{repo1, repo2}

		// モックの振る舞いを定義
		mockUseCase.On("ListRepositories", mock.Anything).Return(repos, nil)

		// ルーターの設定
		router.GET("/repositories", handler.ListRepositories)

		// リクエスト実行
		req, _ := http.NewRequest("GET", "/repositories", nil)
		router.ServeHTTP(rec, req)

		// 検証
		assert.Equal(t, http.StatusOK, rec.Code)

		var response schema.ListRepositoriesResponse
		err := json.Unmarshal(rec.Body.Bytes(), &response)
		require.NoError(t, err)
		assert.Len(t, response.Repositories, len(repos))
		assert.Equal(t, repo1.ID(), response.Repositories[0].ID)
		assert.Equal(t, repo2.ID(), response.Repositories[1].ID)

		// モックの呼び出しを検証
		mockUseCase.AssertExpectations(t)
	})

	t.Run("空のリポジトリ一覧を取得できる場合", func(t *testing.T) {
		// テスト用のセットアップ
		mockUseCase, _, handler, router, rec := setupTest()

		// 空のリポジトリリスト
		var emptyRepos []entity.Repository

		// モックの振る舞いを定義
		mockUseCase.On("ListRepositories", mock.Anything).Return(emptyRepos, nil)

		// ルーターの設定
		router.GET("/repositories", handler.ListRepositories)

		// リクエスト実行
		req, _ := http.NewRequest("GET", "/repositories", nil)
		router.ServeHTTP(rec, req)

		// 検証
		assert.Equal(t, http.StatusOK, rec.Code)

		var response schema.ListRepositoriesResponse
		err := json.Unmarshal(rec.Body.Bytes(), &response)
		require.NoError(t, err)
		assert.Empty(t, response.Repositories)

		// モックの呼び出しを検証
		mockUseCase.AssertExpectations(t)
	})

	t.Run("内部エラーが発生した場合", func(t *testing.T) {
		// テスト用のセットアップ
		mockUseCase, _, handler, router, rec := setupTest()

		// モックの振る舞いを定義
		mockUseCase.On("ListRepositories", mock.Anything).Return(nil, errors.New("internal error"))

		// ルーターの設定
		router.GET("/repositories", handler.ListRepositories)

		// リクエスト実行
		req, _ := http.NewRequest("GET", "/repositories", nil)
		router.ServeHTTP(rec, req)

		// 検証
		assert.Equal(t, http.StatusInternalServerError, rec.Code)

		var response schema.ErrorResponse
		err := json.Unmarshal(rec.Body.Bytes(), &response)
		require.NoError(t, err)
		assert.Equal(t, "INTERNAL_SERVER_ERROR", response.Code)

		// モックの呼び出しを検証
		mockUseCase.AssertExpectations(t)
	})
}

func TestListRepositoryFiles(t *testing.T) {
	t.Run("正常にファイル一覧を取得できる場合", func(t *testing.T) {
		// テスト用のセットアップ
		mockUseCase, _, handler, router, rec := setupTest()

		// テストデータ
		repoID := uuid.NewString()
		fileNodes := []entity.FileNode{
			entity.NewFileNode("README.md", "file"),
			entity.NewFileNode("src/main.go", "file"),
			entity.NewFileNode("docs/index.md", "file"),
		}

		// モックの振る舞いを定義
		mockUseCase.On("ListFiles", mock.Anything, repoID).Return(fileNodes, nil)

		// ルーターの設定
		router.GET("/repositories/:repoId/files", handler.ListRepositoryFiles)

		// リクエスト実行
		req, _ := http.NewRequest("GET", "/repositories/"+repoID+"/files", nil)
		router.ServeHTTP(rec, req)

		// 検証
		assert.Equal(t, http.StatusOK, rec.Code)

		var response schema.ListFilesResponse
		err := json.Unmarshal(rec.Body.Bytes(), &response)
		require.NoError(t, err)
		assert.Len(t, response.Files, len(fileNodes))
		assert.Equal(t, fileNodes[0].Path(), response.Files[0].Path)
		assert.Equal(t, fileNodes[0].Type(), response.Files[0].Type)

		// モックの呼び出しを検証
		mockUseCase.AssertExpectations(t)
	})

	t.Run("リポジトリIDが欠けている場合", func(t *testing.T) {
		// テスト用のセットアップ
		_, _, handler, router, rec := setupTest()

		// ルーターの設定
		router.GET("/repositories/:repoId/files", handler.ListRepositoryFiles)

		// リクエスト実行（空のIDパラメータを使用）
		req, _ := http.NewRequest("GET", "/repositories//files", nil)
		router.ServeHTTP(rec, req)

		// 検証
		assert.Equal(t, http.StatusBadRequest, rec.Code)

		var response schema.ErrorResponse
		err := json.Unmarshal(rec.Body.Bytes(), &response)
		require.NoError(t, err)
		assert.Equal(t, "INVALID_ID", response.Code)
	})

	t.Run("リポジトリが存在しない場合", func(t *testing.T) {
		// テスト用のセットアップ
		mockUseCase, _, handler, router, rec := setupTest()

		// テストデータ
		repoID := uuid.NewString()

		// モックの振る舞いを定義
		mockUseCase.On("ListFiles", mock.Anything, repoID).Return(nil, &apperror.NotFoundError{
			Code:         apperror.CodeResourceNotFound,
			ResourceType: "Repository",
			ResourceID:   repoID,
		})

		// ルーターの設定
		router.GET("/repositories/:repoId/files", handler.ListRepositoryFiles)

		// リクエスト実行
		req, _ := http.NewRequest("GET", "/repositories/"+repoID+"/files", nil)
		router.ServeHTTP(rec, req)

		// 検証
		assert.Equal(t, http.StatusNotFound, rec.Code)

		var response schema.ErrorResponse
		err := json.Unmarshal(rec.Body.Bytes(), &response)
		require.NoError(t, err)
		assert.Equal(t, "NOT_FOUND", response.Code)

		// モックの呼び出しを検証
		mockUseCase.AssertExpectations(t)
	})

	t.Run("アクセストークンが必要な場合", func(t *testing.T) {
		// テスト用のセットアップ
		mockUseCase, _, handler, router, rec := setupTest()

		// テストデータ
		repoID := uuid.NewString()

		// モックの振る舞いを定義
		mockUseCase.On("ListFiles", mock.Anything, repoID).Return(nil, &apperror.ValidationFailedError{
			Code: apperror.CodeValidationFailed,
			Errors: []apperror.FieldError{
				{Field: "access_token", Message: "access token is required for this operation"},
			},
		})

		// ルーターの設定
		router.GET("/repositories/:repoId/files", handler.ListRepositoryFiles)

		// リクエスト実行
		req, _ := http.NewRequest("GET", "/repositories/"+repoID+"/files", nil)
		router.ServeHTTP(rec, req)

		// 検証
		assert.Equal(t, http.StatusBadRequest, rec.Code)

		var response schema.ErrorResponse
		err := json.Unmarshal(rec.Body.Bytes(), &response)
		require.NoError(t, err)
		assert.Equal(t, "BAD_REQUEST", response.Code)

		// モックの呼び出しを検証
		mockUseCase.AssertExpectations(t)
	})

	t.Run("内部エラーが発生した場合", func(t *testing.T) {
		// テスト用のセットアップ
		mockUseCase, _, handler, router, rec := setupTest()

		// テストデータ
		repoID := uuid.NewString()

		// モックの振る舞いを定義
		mockUseCase.On("ListFiles", mock.Anything, repoID).Return(nil, errors.New("internal error"))

		// ルーターの設定
		router.GET("/repositories/:repoId/files", handler.ListRepositoryFiles)

		// リクエスト実行
		req, _ := http.NewRequest("GET", "/repositories/"+repoID+"/files", nil)
		router.ServeHTTP(rec, req)

		// 検証
		assert.Equal(t, http.StatusInternalServerError, rec.Code)

		var response schema.ErrorResponse
		err := json.Unmarshal(rec.Body.Bytes(), &response)
		require.NoError(t, err)
		assert.Equal(t, "INTERNAL_SERVER_ERROR", response.Code)

		// モックの呼び出しを検証
		mockUseCase.AssertExpectations(t)
	})
}

func TestSelectRepositoryFiles(t *testing.T) {
	t.Run("正常にファイルを選択できる場合", func(t *testing.T) {
		// テスト用のセットアップ
		mockUseCase, _, handler, router, rec := setupTest()

		// テストデータ
		repoID := uuid.NewString()
		filePaths := []string{"README.md", "docs/index.md"}
		requestBody := schema.SelectFilesRequest{
			FilePaths: filePaths,
		}
		jsonBody, _ := json.Marshal(requestBody)

		// モックの振る舞いを定義
		mockUseCase.On("SelectFiles", mock.Anything, repoID, filePaths).Return(nil)

		// ルーターの設定
		router.POST("/repositories/:repoId/files/select", handler.SelectRepositoryFiles)

		// リクエスト実行
		req, _ := http.NewRequest("POST", "/repositories/"+repoID+"/files/select", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(rec, req)

		// 検証
		assert.Equal(t, http.StatusOK, rec.Code)

		var response schema.SelectFilesResponse
		err := json.Unmarshal(rec.Body.Bytes(), &response)
		require.NoError(t, err)
		assert.Equal(t, repoID, response.RepoID)
		assert.Equal(t, len(filePaths), response.SelectedFiles)

		// モックの呼び出しを検証
		mockUseCase.AssertExpectations(t)
	})

	t.Run("リポジトリIDが欠けている場合", func(t *testing.T) {
		// テスト用のセットアップ
		_, _, handler, router, rec := setupTest()

		// テストデータ
		filePaths := []string{"README.md", "docs/index.md"}
		requestBody := schema.SelectFilesRequest{
			FilePaths: filePaths,
		}
		jsonBody, _ := json.Marshal(requestBody)

		// ルーターの設定
		router.POST("/repositories/:repoId/files/select", handler.SelectRepositoryFiles)

		// リクエスト実行（空のIDパラメータを使用）
		req, _ := http.NewRequest("POST", "/repositories//files/select", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(rec, req)

		// 検証
		assert.Equal(t, http.StatusBadRequest, rec.Code)

		var response schema.ErrorResponse
		err := json.Unmarshal(rec.Body.Bytes(), &response)
		require.NoError(t, err)
		assert.Equal(t, "INVALID_ID", response.Code)
	})

	t.Run("不正なリクエスト本文の場合", func(t *testing.T) {
		// テスト用のセットアップ
		_, _, handler, router, rec := setupTest()

		// テストデータ
		repoID := uuid.NewString()

		// 不正なJSONデータ
		invalidJSON := `{"filePaths": [}`

		// ルーターの設定
		router.POST("/repositories/:repoId/files/select", handler.SelectRepositoryFiles)

		// リクエスト実行
		req, _ := http.NewRequest("POST", "/repositories/"+repoID+"/files/select", bytes.NewBufferString(invalidJSON))
		req.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(rec, req)

		// 検証
		assert.Equal(t, http.StatusBadRequest, rec.Code)

		var response schema.ErrorResponse
		err := json.Unmarshal(rec.Body.Bytes(), &response)
		require.NoError(t, err)
		assert.Equal(t, "INVALID_REQUEST", response.Code)
	})

	t.Run("空のファイルパスリストの場合", func(t *testing.T) {
		// テスト用のセットアップ
		_, _, handler, router, rec := setupTest()

		// テストデータ
		repoID := uuid.NewString()
		requestBody := schema.SelectFilesRequest{
			FilePaths: []string{},
		}
		jsonBody, _ := json.Marshal(requestBody)

		// ルーターの設定
		router.POST("/repositories/:repoId/files/select", handler.SelectRepositoryFiles)

		// リクエスト実行
		req, _ := http.NewRequest("POST", "/repositories/"+repoID+"/files/select", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(rec, req)

		// 検証
		assert.Equal(t, http.StatusBadRequest, rec.Code)

		var response schema.ErrorResponse
		err := json.Unmarshal(rec.Body.Bytes(), &response)
		require.NoError(t, err)
		assert.Equal(t, "INVALID_REQUEST", response.Code)
		assert.Contains(t, response.Message, "filePaths cannot be empty")
	})

	t.Run("リポジトリが存在しない場合", func(t *testing.T) {
		// テスト用のセットアップ
		mockUseCase, _, handler, router, rec := setupTest()

		// テストデータ
		repoID := uuid.NewString()
		filePaths := []string{"README.md"}
		requestBody := schema.SelectFilesRequest{
			FilePaths: filePaths,
		}
		jsonBody, _ := json.Marshal(requestBody)

		// モックの振る舞いを定義
		mockUseCase.On("SelectFiles", mock.Anything, repoID, filePaths).Return(&apperror.NotFoundError{
			Code:         apperror.CodeResourceNotFound,
			ResourceType: "Repository",
			ResourceID:   repoID,
		})

		// ルーターの設定
		router.POST("/repositories/:repoId/files/select", handler.SelectRepositoryFiles)

		// リクエスト実行
		req, _ := http.NewRequest("POST", "/repositories/"+repoID+"/files/select", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(rec, req)

		// 検証
		assert.Equal(t, http.StatusNotFound, rec.Code)

		var response schema.ErrorResponse
		err := json.Unmarshal(rec.Body.Bytes(), &response)
		require.NoError(t, err)
		assert.Equal(t, "NOT_FOUND", response.Code)

		// モックの呼び出しを検証
		mockUseCase.AssertExpectations(t)
	})

	t.Run("内部エラーが発生した場合", func(t *testing.T) {
		// テスト用のセットアップ
		mockUseCase, _, handler, router, rec := setupTest()

		// テストデータ
		repoID := uuid.NewString()
		filePaths := []string{"README.md"}
		requestBody := schema.SelectFilesRequest{
			FilePaths: filePaths,
		}
		jsonBody, _ := json.Marshal(requestBody)

		// モックの振る舞いを定義
		mockUseCase.On("SelectFiles", mock.Anything, repoID, filePaths).Return(errors.New("internal error"))

		// ルーターの設定
		router.POST("/repositories/:repoId/files/select", handler.SelectRepositoryFiles)

		// リクエスト実行
		req, _ := http.NewRequest("POST", "/repositories/"+repoID+"/files/select", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(rec, req)

		// 検証
		assert.Equal(t, http.StatusInternalServerError, rec.Code)

		var response schema.ErrorResponse
		err := json.Unmarshal(rec.Body.Bytes(), &response)
		require.NoError(t, err)
		assert.Equal(t, "INTERNAL_SERVER_ERROR", response.Code)

		// モックの呼び出しを検証
		mockUseCase.AssertExpectations(t)
	})
}

func TestGetSelectedMarkdown(t *testing.T) {
	t.Run("正常にMarkdown内容を取得できる場合", func(t *testing.T) {
		// テスト用のセットアップ
		mockUseCase, _, handler, router, rec := setupTest()

		// テストデータ
		repoID := uuid.NewString()
		markdown := "# Title\n\nContent\n\n## Section\n\nMore content"

		// モックの振る舞いを定義
		mockUseCase.On("GetSelectedMarkdown", mock.Anything, repoID).Return(markdown, nil)

		// ルーターの設定
		router.GET("/repositories/:repoId/markdown", handler.GetSelectedMarkdown)

		// リクエスト実行
		req, _ := http.NewRequest("GET", "/repositories/"+repoID+"/markdown", nil)
		router.ServeHTTP(rec, req)

		// 検証
		assert.Equal(t, http.StatusOK, rec.Code)

		var response schema.GetMarkdownResponse
		err := json.Unmarshal(rec.Body.Bytes(), &response)
		require.NoError(t, err)
		assert.Equal(t, repoID, response.RepoID)
		assert.Equal(t, markdown, response.Content)

		// モックの呼び出しを検証
		mockUseCase.AssertExpectations(t)
	})

	t.Run("リポジトリIDが欠けている場合", func(t *testing.T) {
		// テスト用のセットアップ
		_, _, handler, router, rec := setupTest()

		// ルーターの設定
		router.GET("/repositories/:repoId/markdown", handler.GetSelectedMarkdown)

		// リクエスト実行（空のIDパラメータを使用）
		req, _ := http.NewRequest("GET", "/repositories//markdown", nil)
		router.ServeHTTP(rec, req)

		// 検証
		assert.Equal(t, http.StatusBadRequest, rec.Code)

		var response schema.ErrorResponse
		err := json.Unmarshal(rec.Body.Bytes(), &response)
		require.NoError(t, err)
		assert.Equal(t, "INVALID_ID", response.Code)
	})

	t.Run("リポジトリが存在しない場合", func(t *testing.T) {
		// テスト用のセットアップ
		mockUseCase, _, handler, router, rec := setupTest()

		// テストデータ
		repoID := uuid.NewString()

		// モックの振る舞いを定義
		mockUseCase.On("GetSelectedMarkdown", mock.Anything, repoID).Return("", &apperror.NotFoundError{
			Code:         apperror.CodeResourceNotFound,
			ResourceType: "Repository",
			ResourceID:   repoID,
		})

		// ルーターの設定
		router.GET("/repositories/:repoId/markdown", handler.GetSelectedMarkdown)

		// リクエスト実行
		req, _ := http.NewRequest("GET", "/repositories/"+repoID+"/markdown", nil)
		router.ServeHTTP(rec, req)

		// 検証
		assert.Equal(t, http.StatusNotFound, rec.Code)

		var response schema.ErrorResponse
		err := json.Unmarshal(rec.Body.Bytes(), &response)
		require.NoError(t, err)
		assert.Equal(t, "NOT_FOUND", response.Code)

		// モックの呼び出しを検証
		mockUseCase.AssertExpectations(t)
	})

	t.Run("内部エラーが発生した場合", func(t *testing.T) {
		// テスト用のセットアップ
		mockUseCase, _, handler, router, rec := setupTest()

		// テストデータ
		repoID := uuid.NewString()

		// モックの振る舞いを定義
		mockUseCase.On("GetSelectedMarkdown", mock.Anything, repoID).Return("", errors.New("internal error"))

		// ルーターの設定
		router.GET("/repositories/:repoId/markdown", handler.GetSelectedMarkdown)

		// リクエスト実行
		req, _ := http.NewRequest("GET", "/repositories/"+repoID+"/markdown", nil)
		router.ServeHTTP(rec, req)

		// 検証
		assert.Equal(t, http.StatusInternalServerError, rec.Code)

		var response schema.ErrorResponse
		err := json.Unmarshal(rec.Body.Bytes(), &response)
		require.NoError(t, err)
		assert.Equal(t, "INTERNAL_SERVER_ERROR", response.Code)

		// モックの呼び出しを検証
		mockUseCase.AssertExpectations(t)
	})
}

func TestUpdateAccessToken(t *testing.T) {
	t.Run("正常にアクセストークンを更新できる場合", func(t *testing.T) {
		// テスト用のセットアップ
		mockUseCase, _, handler, router, rec := setupTest()

		// テストデータ
		repoID := uuid.NewString()
		accessToken := "new-token"
		requestBody := schema.UpdateAccessTokenRequest{
			AccessToken: accessToken,
		}
		jsonBody, _ := json.Marshal(requestBody)

		// モックの振る舞いを定義
		mockUseCase.On("UpdateAccessToken", mock.Anything, repoID, accessToken).Return(nil)

		// ルーターの設定
		router.PUT("/repositories/:repoId/token", handler.UpdateAccessToken)

		// リクエスト実行
		req, _ := http.NewRequest("PUT", "/repositories/"+repoID+"/token", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(rec, req)

		// 検証
		assert.Equal(t, http.StatusOK, rec.Code)

		var response map[string]string
		err := json.Unmarshal(rec.Body.Bytes(), &response)
		require.NoError(t, err)
		assert.Equal(t, repoID, response["repoId"])
		assert.Contains(t, response["message"], "Access token updated successfully")

		// モックの呼び出しを検証
		mockUseCase.AssertExpectations(t)
	})

	t.Run("リポジトリIDが欠けている場合", func(t *testing.T) {
		// テスト用のセットアップ
		_, _, handler, router, rec := setupTest()

		// テストデータ
		accessToken := "new-token"
		requestBody := schema.UpdateAccessTokenRequest{
			AccessToken: accessToken,
		}
		jsonBody, _ := json.Marshal(requestBody)

		// ルーターの設定
		router.PUT("/repositories/:repoId/token", handler.UpdateAccessToken)

		// リクエスト実行（空のIDパラメータを使用）
		req, _ := http.NewRequest("PUT", "/repositories//token", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(rec, req)

		// 検証
		assert.Equal(t, http.StatusBadRequest, rec.Code)

		var response schema.ErrorResponse
		err := json.Unmarshal(rec.Body.Bytes(), &response)
		require.NoError(t, err)
		assert.Equal(t, "INVALID_ID", response.Code)
	})

	t.Run("不正なリクエスト本文の場合", func(t *testing.T) {
		// テスト用のセットアップ
		_, _, handler, router, rec := setupTest()

		// テストデータ
		repoID := uuid.NewString()

		// 不正なJSONデータ
		invalidJSON := `{"accessToken": }`

		// ルーターの設定
		router.PUT("/repositories/:repoId/token", handler.UpdateAccessToken)

		// リクエスト実行
		req, _ := http.NewRequest("PUT", "/repositories/"+repoID+"/token", bytes.NewBufferString(invalidJSON))
		req.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(rec, req)

		// 検証
		assert.Equal(t, http.StatusBadRequest, rec.Code)

		var response schema.ErrorResponse
		err := json.Unmarshal(rec.Body.Bytes(), &response)
		require.NoError(t, err)
		assert.Equal(t, "INVALID_REQUEST", response.Code)
	})

	t.Run("アクセストークンが欠けている場合", func(t *testing.T) {
		// テスト用のセットアップ
		_, _, handler, router, rec := setupTest()

		// テストデータ
		repoID := uuid.NewString()
		// アクセストークンが空のリクエスト
		invalidRequest := `{"accessToken": ""}`

		// ルーターの設定
		router.PUT("/repositories/:repoId/token", handler.UpdateAccessToken)

		// リクエスト実行
		req, _ := http.NewRequest("PUT", "/repositories/"+repoID+"/token", bytes.NewBufferString(invalidRequest))
		req.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(rec, req)

		// 検証
		assert.Equal(t, http.StatusBadRequest, rec.Code)

		var response schema.ErrorResponse
		err := json.Unmarshal(rec.Body.Bytes(), &response)
		require.NoError(t, err)
		assert.Equal(t, "INVALID_REQUEST", response.Code)
	})

	t.Run("リポジトリが存在しない場合", func(t *testing.T) {
		// テスト用のセットアップ
		mockUseCase, _, handler, router, rec := setupTest()

		// テストデータ
		repoID := uuid.NewString()
		accessToken := "new-token"
		requestBody := schema.UpdateAccessTokenRequest{
			AccessToken: accessToken,
		}
		jsonBody, _ := json.Marshal(requestBody)

		// モックの振る舞いを定義
		mockUseCase.On("UpdateAccessToken", mock.Anything, repoID, accessToken).Return(&apperror.NotFoundError{
			Code:         apperror.CodeResourceNotFound,
			ResourceType: "Repository",
			ResourceID:   repoID,
		})

		// ルーターの設定
		router.PUT("/repositories/:repoId/token", handler.UpdateAccessToken)

		// リクエスト実行
		req, _ := http.NewRequest("PUT", "/repositories/"+repoID+"/token", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(rec, req)

		// 検証
		assert.Equal(t, http.StatusNotFound, rec.Code)

		var response schema.ErrorResponse
		err := json.Unmarshal(rec.Body.Bytes(), &response)
		require.NoError(t, err)
		assert.Equal(t, "NOT_FOUND", response.Code)

		// モックの呼び出しを検証
		mockUseCase.AssertExpectations(t)
	})

	t.Run("内部エラーが発生した場合", func(t *testing.T) {
		// テスト用のセットアップ
		mockUseCase, _, handler, router, rec := setupTest()

		// テストデータ
		repoID := uuid.NewString()
		accessToken := "new-token"
		requestBody := schema.UpdateAccessTokenRequest{
			AccessToken: accessToken,
		}
		jsonBody, _ := json.Marshal(requestBody)

		// モックの振る舞いを定義
		mockUseCase.On("UpdateAccessToken", mock.Anything, repoID, accessToken).Return(errors.New("internal error"))

		// ルーターの設定
		router.PUT("/repositories/:repoId/token", handler.UpdateAccessToken)

		// リクエスト実行
		req, _ := http.NewRequest("PUT", "/repositories/"+repoID+"/token", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(rec, req)

		// 検証
		assert.Equal(t, http.StatusInternalServerError, rec.Code)

		var response schema.ErrorResponse
		err := json.Unmarshal(rec.Body.Bytes(), &response)
		require.NoError(t, err)
		assert.Equal(t, "INTERNAL_SERVER_ERROR", response.Code)

		// モックの呼び出しを検証
		mockUseCase.AssertExpectations(t)
	})
}
