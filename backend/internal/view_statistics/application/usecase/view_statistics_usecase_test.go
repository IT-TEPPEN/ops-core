package usecase

import (
	"context"
	"errors"
	"testing"
	"time"

	documentVO "opscore/backend/internal/document/domain/value_object"
	userVO "opscore/backend/internal/user/domain/value_object"
	"opscore/backend/internal/view_statistics/domain/entity"
	"opscore/backend/internal/view_statistics/domain/repository"
	"opscore/backend/internal/view_statistics/domain/value_object"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockViewStatisticsRepository is a mock implementation of ViewStatisticsRepository.
type MockViewStatisticsRepository struct {
	mock.Mock
}

func (m *MockViewStatisticsRepository) Save(ctx context.Context, stats entity.ViewStatistics) error {
	args := m.Called(ctx, stats)
	return args.Error(0)
}

func (m *MockViewStatisticsRepository) FindByDocumentID(ctx context.Context, documentID documentVO.DocumentID) (entity.ViewStatistics, error) {
	args := m.Called(ctx, documentID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(entity.ViewStatistics), args.Error(1)
}

func (m *MockViewStatisticsRepository) FindPopularDocuments(ctx context.Context, limit int, since time.Time) ([]repository.PopularDocument, error) {
	args := m.Called(ctx, limit, since)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]repository.PopularDocument), args.Error(1)
}

func (m *MockViewStatisticsRepository) FindRecentlyViewedDocuments(ctx context.Context, limit int) ([]repository.PopularDocument, error) {
	args := m.Called(ctx, limit)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]repository.PopularDocument), args.Error(1)
}

func (m *MockViewStatisticsRepository) GetUserViewCount(ctx context.Context, userID userVO.UserID) (int64, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockViewStatisticsRepository) GetUserUniqueDocumentCount(ctx context.Context, userID userVO.UserID) (int64, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).(int64), args.Error(1)
}

func TestViewStatisticsUseCase_GetDocumentStatistics(t *testing.T) {
	t.Run("正常にドキュメント統計を取得できる", func(t *testing.T) {
		mockRepo := new(MockViewStatisticsRepository)
		docID := documentVO.GenerateDocumentID()

		stats, _ := entity.NewViewStatistics(value_object.GenerateViewStatisticsID(), docID)
		_ = stats.IncrementView(true)
		_ = stats.IncrementView(false)
		_ = stats.UpdateLastViewedAt(time.Now())
		_ = stats.UpdateAverageViewDuration(120)

		mockRepo.On("FindByDocumentID", mock.Anything, docID).Return(stats, nil)

		uc := NewViewStatisticsUseCase(mockRepo)
		result, err := uc.GetDocumentStatistics(context.Background(), docID.String())

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, docID.String(), result.DocumentID)
		assert.Equal(t, int64(2), result.TotalViews)
		assert.Equal(t, int64(1), result.UniqueViewers)
		assert.Equal(t, 120, result.AverageViewDuration)

		mockRepo.AssertExpectations(t)
	})

	t.Run("無効なDocumentIDでエラーになる", func(t *testing.T) {
		mockRepo := new(MockViewStatisticsRepository)

		uc := NewViewStatisticsUseCase(mockRepo)
		result, err := uc.GetDocumentStatistics(context.Background(), "invalid-id")

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "invalid document ID")
	})

	t.Run("取得失敗でエラーになる", func(t *testing.T) {
		mockRepo := new(MockViewStatisticsRepository)
		docID := documentVO.GenerateDocumentID()

		mockRepo.On("FindByDocumentID", mock.Anything, docID).Return(nil, errors.New("database error"))

		uc := NewViewStatisticsUseCase(mockRepo)
		result, err := uc.GetDocumentStatistics(context.Background(), docID.String())

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "failed to retrieve document statistics")

		mockRepo.AssertExpectations(t)
	})
}

func TestViewStatisticsUseCase_GetUserStatistics(t *testing.T) {
	t.Run("正常にユーザー統計を取得できる", func(t *testing.T) {
		mockRepo := new(MockViewStatisticsRepository)
		userID, _ := userVO.NewUserID("user-123")

		mockRepo.On("GetUserViewCount", mock.Anything, userID).Return(int64(100), nil)
		mockRepo.On("GetUserUniqueDocumentCount", mock.Anything, userID).Return(int64(25), nil)

		uc := NewViewStatisticsUseCase(mockRepo)
		result, err := uc.GetUserStatistics(context.Background(), userID.String())

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, userID.String(), result.UserID)
		assert.Equal(t, int64(100), result.TotalViews)
		assert.Equal(t, int64(25), result.UniqueDocuments)

		mockRepo.AssertExpectations(t)
	})
}

func TestViewStatisticsUseCase_GetPopularDocuments(t *testing.T) {
	t.Run("正常に人気ドキュメントを取得できる", func(t *testing.T) {
		mockRepo := new(MockViewStatisticsRepository)

		popularDocs := []repository.PopularDocument{
			{
				DocumentID:    documentVO.GenerateDocumentID(),
				TotalViews:    1000,
				UniqueViewers: 500,
				LastViewedAt:  time.Now(),
			},
			{
				DocumentID:    documentVO.GenerateDocumentID(),
				TotalViews:    800,
				UniqueViewers: 400,
				LastViewedAt:  time.Now(),
			},
		}

		mockRepo.On("FindPopularDocuments", mock.Anything, 10, mock.AnythingOfType("time.Time")).Return(popularDocs, nil)

		uc := NewViewStatisticsUseCase(mockRepo)
		result, err := uc.GetPopularDocuments(context.Background(), 10, 30)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, 2, len(result.Items))
		assert.Equal(t, 10, result.Limit)
		assert.Equal(t, int64(1000), result.Items[0].TotalViews)

		mockRepo.AssertExpectations(t)
	})

	t.Run("デフォルト値で人気ドキュメントを取得できる", func(t *testing.T) {
		mockRepo := new(MockViewStatisticsRepository)

		mockRepo.On("FindPopularDocuments", mock.Anything, 10, mock.AnythingOfType("time.Time")).Return([]repository.PopularDocument{}, nil)

		uc := NewViewStatisticsUseCase(mockRepo)
		result, err := uc.GetPopularDocuments(context.Background(), 0, 0)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, 10, result.Limit)

		mockRepo.AssertExpectations(t)
	})
}

func TestViewStatisticsUseCase_GetRecentlyViewedDocuments(t *testing.T) {
	t.Run("正常に最近閲覧されたドキュメントを取得できる", func(t *testing.T) {
		mockRepo := new(MockViewStatisticsRepository)

		recentDocs := []repository.PopularDocument{
			{
				DocumentID:   documentVO.GenerateDocumentID(),
				TotalViews:   50,
				LastViewedAt: time.Now(),
			},
			{
				DocumentID:   documentVO.GenerateDocumentID(),
				TotalViews:   30,
				LastViewedAt: time.Now().Add(-1 * time.Hour),
			},
		}

		mockRepo.On("FindRecentlyViewedDocuments", mock.Anything, 10).Return(recentDocs, nil)

		uc := NewViewStatisticsUseCase(mockRepo)
		result, err := uc.GetRecentlyViewedDocuments(context.Background(), 10)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, 2, len(result.Items))
		assert.Equal(t, 10, result.Limit)
		assert.Equal(t, int64(50), result.Items[0].TotalViews)

		mockRepo.AssertExpectations(t)
	})
}
