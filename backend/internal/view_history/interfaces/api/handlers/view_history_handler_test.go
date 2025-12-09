package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"opscore/backend/internal/view_history/application/dto"
	"opscore/backend/internal/view_history/interfaces/api/schema"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockViewHistoryUseCase is a mock implementation of ViewHistoryUseCase.
type MockViewHistoryUseCase struct {
	mock.Mock
}

func (m *MockViewHistoryUseCase) RecordView(ctx context.Context, req *dto.RecordViewRequest) (*dto.ViewHistoryResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.ViewHistoryResponse), args.Error(1)
}

func (m *MockViewHistoryUseCase) GetViewHistory(ctx context.Context, userID string, limit int, offset int) (*dto.ViewHistoryListResponse, error) {
	args := m.Called(ctx, userID, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.ViewHistoryListResponse), args.Error(1)
}

func (m *MockViewHistoryUseCase) GetDocumentViewHistory(ctx context.Context, documentID string, limit int, offset int) (*dto.ViewHistoryListResponse, error) {
	args := m.Called(ctx, documentID, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.ViewHistoryListResponse), args.Error(1)
}

// MockLogger is a simple mock logger for testing.
type MockLogger struct{}

func (m *MockLogger) Info(msg string, args ...any)  {}
func (m *MockLogger) Error(msg string, args ...any) {}
func (m *MockLogger) Debug(msg string, args ...any) {}
func (m *MockLogger) Warn(msg string, args ...any)  {}

func TestViewHistoryHandler_RecordView(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("正常に閲覧を記録できる", func(t *testing.T) {
		mockUseCase := new(MockViewHistoryUseCase)
		handler := NewViewHistoryHandler(mockUseCase, &MockLogger{})

		expectedResp := &dto.ViewHistoryResponse{
			ID:           "view-123",
			DocumentID:   "doc-456",
			UserID:       "user-789",
			ViewedAt:     time.Now(),
			ViewDuration: 0,
		}

		mockUseCase.On("RecordView", mock.Anything, mock.AnythingOfType("*dto.RecordViewRequest")).Return(expectedResp, nil)

		reqBody := schema.RecordViewRequest{
			UserID: "user-789",
		}
		body, _ := json.Marshal(reqBody)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest("POST", "/api/documents/doc-456/views", bytes.NewBuffer(body))
		c.Request.Header.Set("Content-Type", "application/json")
		c.Params = gin.Params{{Key: "id", Value: "doc-456"}}

		handler.RecordView(c)

		assert.Equal(t, http.StatusCreated, w.Code)
		mockUseCase.AssertExpectations(t)
	})

	t.Run("無効なリクエストでエラーになる", func(t *testing.T) {
		mockUseCase := new(MockViewHistoryUseCase)
		handler := NewViewHistoryHandler(mockUseCase, &MockLogger{})

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest("POST", "/api/documents/doc-456/views", bytes.NewBuffer([]byte("invalid json")))
		c.Request.Header.Set("Content-Type", "application/json")
		c.Params = gin.Params{{Key: "id", Value: "doc-456"}}

		handler.RecordView(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("UseCaseエラーで500を返す", func(t *testing.T) {
		mockUseCase := new(MockViewHistoryUseCase)
		handler := NewViewHistoryHandler(mockUseCase, &MockLogger{})

		mockUseCase.On("RecordView", mock.Anything, mock.AnythingOfType("*dto.RecordViewRequest")).Return(nil, errors.New("database error"))

		reqBody := schema.RecordViewRequest{
			UserID: "user-789",
		}
		body, _ := json.Marshal(reqBody)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest("POST", "/api/documents/doc-456/views", bytes.NewBuffer(body))
		c.Request.Header.Set("Content-Type", "application/json")
		c.Params = gin.Params{{Key: "id", Value: "doc-456"}}

		handler.RecordView(c)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		mockUseCase.AssertExpectations(t)
	})
}

func TestViewHistoryHandler_GetUserViewHistory(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("正常にユーザー閲覧履歴を取得できる", func(t *testing.T) {
		mockUseCase := new(MockViewHistoryUseCase)
		handler := NewViewHistoryHandler(mockUseCase, &MockLogger{})

		expectedResp := &dto.ViewHistoryListResponse{
			Items: []dto.ViewHistoryResponse{
				{
					ID:         "view-1",
					DocumentID: "doc-1",
					UserID:     "user-123",
					ViewedAt:   time.Now(),
				},
			},
			TotalCount: 1,
			Limit:      50,
			Offset:     0,
		}

		mockUseCase.On("GetViewHistory", mock.Anything, "user-123", 50, 0).Return(expectedResp, nil)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest("GET", "/api/users/user-123/view-history", nil)
		c.Params = gin.Params{{Key: "id", Value: "user-123"}}

		handler.GetUserViewHistory(c)

		assert.Equal(t, http.StatusOK, w.Code)
		mockUseCase.AssertExpectations(t)
	})
}

func TestViewHistoryHandler_GetDocumentViewHistory(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("正常にドキュメント閲覧履歴を取得できる", func(t *testing.T) {
		mockUseCase := new(MockViewHistoryUseCase)
		handler := NewViewHistoryHandler(mockUseCase, &MockLogger{})

		expectedResp := &dto.ViewHistoryListResponse{
			Items: []dto.ViewHistoryResponse{
				{
					ID:         "view-1",
					DocumentID: "doc-123",
					UserID:     "user-1",
					ViewedAt:   time.Now(),
				},
			},
			TotalCount: 1,
			Limit:      50,
			Offset:     0,
		}

		mockUseCase.On("GetDocumentViewHistory", mock.Anything, "doc-123", 50, 0).Return(expectedResp, nil)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest("GET", "/api/documents/doc-123/view-history", nil)
		c.Params = gin.Params{{Key: "id", Value: "doc-123"}}

		handler.GetDocumentViewHistory(c)

		assert.Equal(t, http.StatusOK, w.Code)
		mockUseCase.AssertExpectations(t)
	})
}
