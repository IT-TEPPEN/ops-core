package handlers

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"opscore/backend/internal/view_statistics/application/dto"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockViewStatisticsUseCase is a mock implementation of ViewStatisticsUseCase.
type MockViewStatisticsUseCase struct {
	mock.Mock
}

func (m *MockViewStatisticsUseCase) GetDocumentStatistics(ctx context.Context, documentID string) (*dto.DocumentStatisticsResponse, error) {
	args := m.Called(ctx, documentID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.DocumentStatisticsResponse), args.Error(1)
}

func (m *MockViewStatisticsUseCase) GetUserStatistics(ctx context.Context, userID string) (*dto.UserStatisticsResponse, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.UserStatisticsResponse), args.Error(1)
}

func (m *MockViewStatisticsUseCase) GetPopularDocuments(ctx context.Context, limit int, days int) (*dto.PopularDocumentsResponse, error) {
	args := m.Called(ctx, limit, days)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.PopularDocumentsResponse), args.Error(1)
}

func (m *MockViewStatisticsUseCase) GetRecentlyViewedDocuments(ctx context.Context, limit int) (*dto.RecentDocumentsResponse, error) {
	args := m.Called(ctx, limit)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.RecentDocumentsResponse), args.Error(1)
}

// MockLogger is a simple mock logger for testing.
type MockLogger struct{}

func (m *MockLogger) Info(msg string, args ...any)  {}
func (m *MockLogger) Error(msg string, args ...any) {}
func (m *MockLogger) Debug(msg string, args ...any) {}
func (m *MockLogger) Warn(msg string, args ...any)  {}

func TestViewStatisticsHandler_GetDocumentStatistics(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("正常にドキュメント統計を取得できる", func(t *testing.T) {
		mockUseCase := new(MockViewStatisticsUseCase)
		handler := NewViewStatisticsHandler(mockUseCase, &MockLogger{})

		expectedResp := &dto.DocumentStatisticsResponse{
			DocumentID:          "doc-123",
			TotalViews:          1000,
			UniqueViewers:       500,
			LastViewedAt:        time.Now(),
			AverageViewDuration: 120,
		}

		mockUseCase.On("GetDocumentStatistics", mock.Anything, "doc-123").Return(expectedResp, nil)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest("GET", "/api/documents/doc-123/statistics", nil)
		c.Params = gin.Params{{Key: "id", Value: "doc-123"}}

		handler.GetDocumentStatistics(c)

		assert.Equal(t, http.StatusOK, w.Code)
		mockUseCase.AssertExpectations(t)
	})

	t.Run("UseCaseエラーで500を返す", func(t *testing.T) {
		mockUseCase := new(MockViewStatisticsUseCase)
		handler := NewViewStatisticsHandler(mockUseCase, &MockLogger{})

		mockUseCase.On("GetDocumentStatistics", mock.Anything, "doc-123").Return(nil, errors.New("database error"))

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest("GET", "/api/documents/doc-123/statistics", nil)
		c.Params = gin.Params{{Key: "id", Value: "doc-123"}}

		handler.GetDocumentStatistics(c)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		mockUseCase.AssertExpectations(t)
	})
}

func TestViewStatisticsHandler_GetUserStatistics(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("正常にユーザー統計を取得できる", func(t *testing.T) {
		mockUseCase := new(MockViewStatisticsUseCase)
		handler := NewViewStatisticsHandler(mockUseCase, &MockLogger{})

		expectedResp := &dto.UserStatisticsResponse{
			UserID:          "user-123",
			TotalViews:      150,
			UniqueDocuments: 30,
		}

		mockUseCase.On("GetUserStatistics", mock.Anything, "user-123").Return(expectedResp, nil)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest("GET", "/api/users/user-123/statistics", nil)
		c.Params = gin.Params{{Key: "id", Value: "user-123"}}

		handler.GetUserStatistics(c)

		assert.Equal(t, http.StatusOK, w.Code)
		mockUseCase.AssertExpectations(t)
	})
}

func TestViewStatisticsHandler_GetPopularDocuments(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("正常に人気ドキュメントを取得できる", func(t *testing.T) {
		mockUseCase := new(MockViewStatisticsUseCase)
		handler := NewViewStatisticsHandler(mockUseCase, &MockLogger{})

		expectedResp := &dto.PopularDocumentsResponse{
			Items: []dto.PopularDocumentResponse{
				{
					DocumentID:    "doc-1",
					TotalViews:    1000,
					UniqueViewers: 500,
					LastViewedAt:  time.Now(),
				},
			},
			Limit: 10,
		}

		mockUseCase.On("GetPopularDocuments", mock.Anything, 10, 30).Return(expectedResp, nil)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest("GET", "/api/statistics/popular-documents", nil)

		handler.GetPopularDocuments(c)

		assert.Equal(t, http.StatusOK, w.Code)
		mockUseCase.AssertExpectations(t)
	})
}

func TestViewStatisticsHandler_GetRecentlyViewedDocuments(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("正常に最近閲覧されたドキュメントを取得できる", func(t *testing.T) {
		mockUseCase := new(MockViewStatisticsUseCase)
		handler := NewViewStatisticsHandler(mockUseCase, &MockLogger{})

		expectedResp := &dto.RecentDocumentsResponse{
			Items: []dto.RecentDocumentResponse{
				{
					DocumentID:   "doc-1",
					LastViewedAt: time.Now(),
					TotalViews:   50,
				},
			},
			Limit: 10,
		}

		mockUseCase.On("GetRecentlyViewedDocuments", mock.Anything, 10).Return(expectedResp, nil)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest("GET", "/api/statistics/recent-documents", nil)

		handler.GetRecentlyViewedDocuments(c)

		assert.Equal(t, http.StatusOK, w.Code)
		mockUseCase.AssertExpectations(t)
	})
}
