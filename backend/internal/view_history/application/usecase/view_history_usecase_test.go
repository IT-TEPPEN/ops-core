package usecase

import (
	"context"
	"errors"
	"testing"
	"time"

	documentVO "opscore/backend/internal/document/domain/value_object"
	userVO "opscore/backend/internal/user/domain/value_object"
	"opscore/backend/internal/view_history/application/dto"
	"opscore/backend/internal/view_history/domain/entity"
	"opscore/backend/internal/view_history/domain/value_object"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockViewHistoryRepository is a mock implementation of ViewHistoryRepository.
type MockViewHistoryRepository struct {
	mock.Mock
}

func (m *MockViewHistoryRepository) Save(ctx context.Context, viewHistory entity.ViewHistory) error {
	args := m.Called(ctx, viewHistory)
	return args.Error(0)
}

func (m *MockViewHistoryRepository) FindByID(ctx context.Context, id value_object.ViewHistoryID) (entity.ViewHistory, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(entity.ViewHistory), args.Error(1)
}

func (m *MockViewHistoryRepository) FindByUserID(ctx context.Context, userID userVO.UserID, limit int, offset int) ([]entity.ViewHistory, int64, error) {
	args := m.Called(ctx, userID, limit, offset)
	if args.Get(0) == nil {
		return nil, 0, args.Error(2)
	}
	return args.Get(0).([]entity.ViewHistory), args.Get(1).(int64), args.Error(2)
}

func (m *MockViewHistoryRepository) FindByDocumentID(ctx context.Context, documentID documentVO.DocumentID, limit int, offset int) ([]entity.ViewHistory, int64, error) {
	args := m.Called(ctx, documentID, limit, offset)
	if args.Get(0) == nil {
		return nil, 0, args.Error(2)
	}
	return args.Get(0).([]entity.ViewHistory), args.Get(1).(int64), args.Error(2)
}

func (m *MockViewHistoryRepository) FindByUserIDAndDocumentID(ctx context.Context, userID userVO.UserID, documentID documentVO.DocumentID) ([]entity.ViewHistory, error) {
	args := m.Called(ctx, userID, documentID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]entity.ViewHistory), args.Error(1)
}

func (m *MockViewHistoryRepository) FindRecentByUserID(ctx context.Context, userID userVO.UserID, since time.Time, limit int) ([]entity.ViewHistory, error) {
	args := m.Called(ctx, userID, since, limit)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]entity.ViewHistory), args.Error(1)
}

func (m *MockViewHistoryRepository) CountByDocumentID(ctx context.Context, documentID documentVO.DocumentID) (int64, error) {
	args := m.Called(ctx, documentID)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockViewHistoryRepository) CountUniqueViewersByDocumentID(ctx context.Context, documentID documentVO.DocumentID) (int64, error) {
	args := m.Called(ctx, documentID)
	return args.Get(0).(int64), args.Error(1)
}

func TestViewHistoryUseCase_RecordView(t *testing.T) {
	t.Run("正常に閲覧を記録できる", func(t *testing.T) {
		mockRepo := new(MockViewHistoryRepository)
		docID := documentVO.GenerateDocumentID()
		userID, _ := userVO.NewUserID("user-123")

		req := &dto.RecordViewRequest{
			DocumentID: docID.String(),
			UserID:     userID.String(),
		}

		mockRepo.On("Save", mock.Anything, mock.AnythingOfType("*entity.viewHistory")).Return(nil)

		uc := NewViewHistoryUseCase(mockRepo)
		result, err := uc.RecordView(context.Background(), req)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.NotEmpty(t, result.ID)
		assert.Equal(t, docID.String(), result.DocumentID)
		assert.Equal(t, userID.String(), result.UserID)
		assert.False(t, result.ViewedAt.IsZero())

		mockRepo.AssertExpectations(t)
	})

	t.Run("無効なDocumentIDでエラーになる", func(t *testing.T) {
		mockRepo := new(MockViewHistoryRepository)

		req := &dto.RecordViewRequest{
			DocumentID: "invalid-id",
			UserID:     "user-123",
		}

		uc := NewViewHistoryUseCase(mockRepo)
		result, err := uc.RecordView(context.Background(), req)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "invalid document ID")
	})

	t.Run("保存失敗でエラーになる", func(t *testing.T) {
		mockRepo := new(MockViewHistoryRepository)
		docID := documentVO.GenerateDocumentID()
		userID, _ := userVO.NewUserID("user-123")

		req := &dto.RecordViewRequest{
			DocumentID: docID.String(),
			UserID:     userID.String(),
		}

		mockRepo.On("Save", mock.Anything, mock.AnythingOfType("*entity.viewHistory")).Return(errors.New("database error"))

		uc := NewViewHistoryUseCase(mockRepo)
		result, err := uc.RecordView(context.Background(), req)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "failed to save view history")

		mockRepo.AssertExpectations(t)
	})
}

func TestViewHistoryUseCase_GetViewHistory(t *testing.T) {
	t.Run("正常にユーザーの閲覧履歴を取得できる", func(t *testing.T) {
		mockRepo := new(MockViewHistoryRepository)
		userID, _ := userVO.NewUserID("user-123")

		// Create mock view histories
		vh1, _ := entity.RecordViewHistory(documentVO.GenerateDocumentID(), userID)
		vh2, _ := entity.RecordViewHistory(documentVO.GenerateDocumentID(), userID)
		histories := []entity.ViewHistory{vh1, vh2}

		mockRepo.On("FindByUserID", mock.Anything, userID, 50, 0).Return(histories, int64(100), nil)

		uc := NewViewHistoryUseCase(mockRepo)
		result, err := uc.GetViewHistory(context.Background(), userID.String(), 50, 0)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, 2, len(result.Items))
		assert.Equal(t, 100, result.TotalCount)
		assert.Equal(t, 50, result.Limit)
		assert.Equal(t, 0, result.Offset)

		mockRepo.AssertExpectations(t)
	})
}

func TestViewHistoryUseCase_GetDocumentViewHistory(t *testing.T) {
	t.Run("正常にドキュメントの閲覧履歴を取得できる", func(t *testing.T) {
		mockRepo := new(MockViewHistoryRepository)
		docID := documentVO.GenerateDocumentID()
		userID1, _ := userVO.NewUserID("user-1")
		userID2, _ := userVO.NewUserID("user-2")

		// Create mock view histories
		vh1, _ := entity.RecordViewHistory(docID, userID1)
		vh2, _ := entity.RecordViewHistory(docID, userID2)
		histories := []entity.ViewHistory{vh1, vh2}

		mockRepo.On("FindByDocumentID", mock.Anything, docID, 50, 0).Return(histories, int64(50), nil)

		uc := NewViewHistoryUseCase(mockRepo)
		result, err := uc.GetDocumentViewHistory(context.Background(), docID.String(), 50, 0)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, 2, len(result.Items))
		assert.Equal(t, 50, result.TotalCount)
		assert.Equal(t, 50, result.Limit)
		assert.Equal(t, 0, result.Offset)

		mockRepo.AssertExpectations(t)
	})
}
