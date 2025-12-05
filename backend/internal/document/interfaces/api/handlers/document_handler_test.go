package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"opscore/backend/internal/document/application/dto"
	"opscore/backend/internal/document/application/usecase"
	"opscore/backend/internal/document/interfaces/api/schema"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func setupDocumentTest() (*usecase.MockDocumentUseCase, *MockLogger, *DocumentHandler, *gin.Engine, *httptest.ResponseRecorder) {
	gin.SetMode(gin.TestMode)

	mockUseCase := new(usecase.MockDocumentUseCase)
	mockLogger := new(MockLogger)

	// Set up the logger mock
	mockLogger.On("Info", mock.AnythingOfType("string"), mock.Anything).Maybe()
	mockLogger.On("Error", mock.AnythingOfType("string"), mock.Anything).Maybe()
	mockLogger.On("Warn", mock.AnythingOfType("string"), mock.Anything).Maybe()
	mockLogger.On("Debug", mock.AnythingOfType("string"), mock.Anything).Maybe()

	handler := NewDocumentHandler(mockUseCase, mockLogger)

	router := gin.New()
	router.Use(func(c *gin.Context) {
		c.Set("request_id", "test-request-id")
		c.Next()
	})

	rec := httptest.NewRecorder()

	return mockUseCase, mockLogger, handler, router, rec
}

func TestDocumentHandler_CreateDocument(t *testing.T) {
	t.Run("正常にドキュメントを作成できる", func(t *testing.T) {
		mockUseCase, _, handler, router, rec := setupDocumentTest()

		requestBody := schema.CreateDocumentRequest{
			RepositoryID: "a1b2c3d4-e5f6-7890-1234-567890abcdef",
			FilePath:     "docs/test.md",
			CommitHash:   "abc1234567890",
			Title:        "Test Document",
			DocType:      "procedure",
			Owner:        "test-owner",
			Tags:         []string{"test"},
			Content:      "# Test Content",
			AccessScope:  "public",
			IsAutoUpdate: true,
		}
		jsonBody, _ := json.Marshal(requestBody)

		now := time.Now()
		mockResponse := &dto.DocumentResponse{
			ID:           "doc-uuid-1234",
			RepositoryID: requestBody.RepositoryID,
			Owner:        requestBody.Owner,
			IsPublished:  true,
			IsAutoUpdate: requestBody.IsAutoUpdate,
			AccessScope:  requestBody.AccessScope,
			CurrentVersion: &dto.DocumentVersionResponse{
				ID:            "ver-uuid-1234",
				DocumentID:    "doc-uuid-1234",
				VersionNumber: 1,
				FilePath:      requestBody.FilePath,
				CommitHash:    requestBody.CommitHash,
				Title:         requestBody.Title,
				DocType:       requestBody.DocType,
				Tags:          requestBody.Tags,
				Content:       requestBody.Content,
				PublishedAt:   now,
				IsCurrent:     true,
			},
			VersionCount: 1,
			CreatedAt:    now,
			UpdatedAt:    now,
		}

		mockUseCase.On("CreateDocument", mock.Anything, mock.AnythingOfType("*dto.CreateDocumentRequest")).Return(mockResponse, nil)

		router.POST("/documents", handler.CreateDocument)

		req, _ := http.NewRequest("POST", "/documents", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusCreated, rec.Code)

		var response schema.DocumentResponse
		err := json.Unmarshal(rec.Body.Bytes(), &response)
		require.NoError(t, err)
		assert.Equal(t, mockResponse.ID, response.ID)
		assert.Equal(t, mockResponse.Owner, response.Owner)
		assert.True(t, response.IsPublished)

		mockUseCase.AssertExpectations(t)
	})

	t.Run("無効なリクエストボディでエラーになる", func(t *testing.T) {
		_, _, handler, router, rec := setupDocumentTest()

		invalidJSON := []byte(`{"invalid": json}`)

		router.POST("/documents", handler.CreateDocument)

		req, _ := http.NewRequest("POST", "/documents", bytes.NewBuffer(invalidJSON))
		req.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})
}

func TestDocumentHandler_GetDocument(t *testing.T) {
	t.Run("正常にドキュメントを取得できる", func(t *testing.T) {
		mockUseCase, _, handler, router, rec := setupDocumentTest()

		docID := "a1b2c3d4-e5f6-7890-1234-567890abcdef"
		now := time.Now()
		mockResponse := &dto.DocumentResponse{
			ID:           docID,
			RepositoryID: "repo-uuid-1234",
			Owner:        "test-owner",
			IsPublished:  true,
			IsAutoUpdate: false,
			AccessScope:  "public",
			VersionCount: 1,
			CreatedAt:    now,
			UpdatedAt:    now,
		}

		mockUseCase.On("GetDocument", mock.Anything, docID).Return(mockResponse, nil)

		router.GET("/documents/:docId", handler.GetDocument)

		req, _ := http.NewRequest("GET", "/documents/"+docID, nil)
		router.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)

		var response schema.DocumentResponse
		err := json.Unmarshal(rec.Body.Bytes(), &response)
		require.NoError(t, err)
		assert.Equal(t, docID, response.ID)

		mockUseCase.AssertExpectations(t)
	})

	t.Run("存在しないドキュメントでNotFoundエラーになる", func(t *testing.T) {
		mockUseCase, _, handler, router, rec := setupDocumentTest()

		docID := "nonexistent-uuid-1234-5678-abcdefabcdef"
		mockUseCase.On("GetDocument", mock.Anything, docID).Return(nil, &mockNotFoundError{})

		router.GET("/documents/:docId", handler.GetDocument)

		req, _ := http.NewRequest("GET", "/documents/"+docID, nil)
		router.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusNotFound, rec.Code)

		mockUseCase.AssertExpectations(t)
	})
}

func TestDocumentHandler_ListDocuments(t *testing.T) {
	t.Run("正常にドキュメント一覧を取得できる", func(t *testing.T) {
		mockUseCase, _, handler, router, rec := setupDocumentTest()

		now := time.Now()
		mockResponse := []dto.DocumentListItemResponse{
			{
				ID:           "doc-uuid-1",
				RepositoryID: "repo-uuid-1",
				Title:        "Document 1",
				Owner:        "owner-1",
				DocType:      "procedure",
				Tags:         []string{"tag1"},
				IsPublished:  true,
				VersionCount: 1,
				CreatedAt:    now,
				UpdatedAt:    now,
			},
			{
				ID:           "doc-uuid-2",
				RepositoryID: "repo-uuid-2",
				Title:        "Document 2",
				Owner:        "owner-2",
				DocType:      "knowledge",
				Tags:         []string{"tag2"},
				IsPublished:  true,
				VersionCount: 2,
				CreatedAt:    now,
				UpdatedAt:    now,
			},
		}

		mockUseCase.On("ListDocuments", mock.Anything).Return(mockResponse, nil)

		router.GET("/documents", handler.ListDocuments)

		req, _ := http.NewRequest("GET", "/documents", nil)
		router.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)

		var response schema.ListDocumentsResponse
		err := json.Unmarshal(rec.Body.Bytes(), &response)
		require.NoError(t, err)
		assert.Len(t, response.Documents, 2)

		mockUseCase.AssertExpectations(t)
	})

	t.Run("リポジトリIDでフィルタリングできる", func(t *testing.T) {
		mockUseCase, _, handler, router, rec := setupDocumentTest()

		repoID := "repo-uuid-1234"
		mockResponse := []dto.DocumentListItemResponse{}

		mockUseCase.On("ListDocumentsByRepository", mock.Anything, repoID).Return(mockResponse, nil)

		router.GET("/documents", handler.ListDocuments)

		req, _ := http.NewRequest("GET", "/documents?repository_id="+repoID, nil)
		router.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)

		mockUseCase.AssertExpectations(t)
	})
}

func TestDocumentHandler_GetDocumentVersions(t *testing.T) {
	t.Run("正常にバージョン履歴を取得できる", func(t *testing.T) {
		mockUseCase, _, handler, router, rec := setupDocumentTest()

		docID := "a1b2c3d4-e5f6-7890-1234-567890abcdef"
		now := time.Now()
		mockResponse := &dto.VersionHistoryResponse{
			DocumentID: docID,
			Versions: []dto.DocumentVersionResponse{
				{
					ID:            "ver-uuid-1",
					DocumentID:    docID,
					VersionNumber: 1,
					FilePath:      "docs/test.md",
					CommitHash:    "abc1234",
					Title:         "Version 1",
					DocType:       "procedure",
					Content:       "# V1",
					PublishedAt:   now,
					IsCurrent:     false,
				},
				{
					ID:            "ver-uuid-2",
					DocumentID:    docID,
					VersionNumber: 2,
					FilePath:      "docs/test.md",
					CommitHash:    "def5678",
					Title:         "Version 2",
					DocType:       "procedure",
					Content:       "# V2",
					PublishedAt:   now,
					IsCurrent:     true,
				},
			},
		}

		mockUseCase.On("GetDocumentVersions", mock.Anything, docID).Return(mockResponse, nil)

		router.GET("/documents/:docId/versions", handler.GetDocumentVersions)

		req, _ := http.NewRequest("GET", "/documents/"+docID+"/versions", nil)
		router.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)

		var response schema.VersionHistoryResponse
		err := json.Unmarshal(rec.Body.Bytes(), &response)
		require.NoError(t, err)
		assert.Len(t, response.Versions, 2)

		mockUseCase.AssertExpectations(t)
	})
}

func TestDocumentHandler_UpdateDocument(t *testing.T) {
	t.Run("正常にドキュメントを更新できる", func(t *testing.T) {
		mockUseCase, _, handler, router, rec := setupDocumentTest()

		docID := "a1b2c3d4-e5f6-7890-1234-567890abcdef"
		requestBody := schema.UpdateDocumentRequest{
			FilePath:   "docs/updated.md",
			CommitHash: "def4567890123",
			Title:      "Updated Document",
			DocType:    "procedure",
			Tags:       []string{"updated"},
			Content:    "# Updated Content",
		}
		jsonBody, _ := json.Marshal(requestBody)

		now := time.Now()
		mockResponse := &dto.DocumentResponse{
			ID:           docID,
			RepositoryID: "repo-uuid-1234",
			Owner:        "test-owner",
			IsPublished:  true,
			IsAutoUpdate: false,
			AccessScope:  "public",
			CurrentVersion: &dto.DocumentVersionResponse{
				ID:            "ver-uuid-2",
				DocumentID:    docID,
				VersionNumber: 2,
				FilePath:      requestBody.FilePath,
				CommitHash:    requestBody.CommitHash,
				Title:         requestBody.Title,
				DocType:       requestBody.DocType,
				Tags:          requestBody.Tags,
				Content:       requestBody.Content,
				PublishedAt:   now,
				IsCurrent:     true,
			},
			VersionCount: 2,
			CreatedAt:    now,
			UpdatedAt:    now,
		}

		mockUseCase.On("UpdateDocument", mock.Anything, docID, mock.AnythingOfType("*dto.UpdateDocumentRequest")).Return(mockResponse, nil)

		router.PUT("/documents/:docId", handler.UpdateDocument)

		req, _ := http.NewRequest("PUT", "/documents/"+docID, bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)

		var response schema.DocumentResponse
		err := json.Unmarshal(rec.Body.Bytes(), &response)
		require.NoError(t, err)
		assert.Equal(t, 2, response.VersionCount)

		mockUseCase.AssertExpectations(t)
	})
}

func TestDocumentHandler_RollbackDocumentVersion(t *testing.T) {
	t.Run("正常にロールバックできる", func(t *testing.T) {
		mockUseCase, _, handler, router, rec := setupDocumentTest()

		docID := "a1b2c3d4-e5f6-7890-1234-567890abcdef"
		versionNumber := 1

		now := time.Now()
		mockResponse := &dto.DocumentResponse{
			ID:           docID,
			RepositoryID: "repo-uuid-1234",
			Owner:        "test-owner",
			IsPublished:  true,
			AccessScope:  "public",
			CurrentVersion: &dto.DocumentVersionResponse{
				ID:            "ver-uuid-1",
				DocumentID:    docID,
				VersionNumber: versionNumber,
				Title:         "Version 1",
				DocType:       "procedure",
				Content:       "# V1",
				PublishedAt:   now,
				IsCurrent:     true,
			},
			VersionCount: 2,
			CreatedAt:    now,
			UpdatedAt:    now,
		}

		mockUseCase.On("RollbackDocumentVersion", mock.Anything, docID, versionNumber).Return(mockResponse, nil)

		router.POST("/documents/:docId/versions/:version/rollback", handler.RollbackDocumentVersion)

		req, _ := http.NewRequest("POST", "/documents/"+docID+"/versions/1/rollback", nil)
		router.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)

		var response schema.DocumentResponse
		err := json.Unmarshal(rec.Body.Bytes(), &response)
		require.NoError(t, err)
		assert.Equal(t, versionNumber, response.CurrentVersion.VersionNumber)

		mockUseCase.AssertExpectations(t)
	})

	t.Run("無効なバージョン番号でエラーになる", func(t *testing.T) {
		_, _, handler, router, rec := setupDocumentTest()

		docID := "a1b2c3d4-e5f6-7890-1234-567890abcdef"

		router.POST("/documents/:docId/versions/:version/rollback", handler.RollbackDocumentVersion)

		req, _ := http.NewRequest("POST", "/documents/"+docID+"/versions/invalid/rollback", nil)
		router.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})
}

// Mock error types for testing
type mockNotFoundError struct{}

func (e *mockNotFoundError) Error() string {
	return "not found"
}

func (e *mockNotFoundError) Is(target error) bool {
	return target.Error() == "resource not found"
}
