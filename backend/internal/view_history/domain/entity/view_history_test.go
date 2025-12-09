package entity

import (
	"testing"
	"time"

	documentVO "opscore/backend/internal/document/domain/value_object"
	userVO "opscore/backend/internal/user/domain/value_object"
	"opscore/backend/internal/view_history/domain/value_object"

	"github.com/stretchr/testify/assert"
)

func TestNewViewHistory(t *testing.T) {
	t.Run("正常にViewHistoryを作成できる", func(t *testing.T) {
		id := value_object.GenerateViewHistoryID()
		docID := documentVO.GenerateDocumentID()
		userID, _ := userVO.NewUserID("user-123")
		viewedAt := time.Now()

		vh, err := NewViewHistory(id, docID, userID, viewedAt)

		assert.NoError(t, err)
		assert.NotNil(t, vh)
		assert.Equal(t, id, vh.ID())
		assert.Equal(t, docID, vh.DocumentID())
		assert.Equal(t, userID, vh.UserID())
		assert.Equal(t, viewedAt, vh.ViewedAt())
		assert.Equal(t, 0, vh.ViewDuration())
	})

	t.Run("空のViewHistoryIDでエラーになる", func(t *testing.T) {
		docID := documentVO.GenerateDocumentID()
		userID, _ := userVO.NewUserID("user-123")
		viewedAt := time.Now()

		vh, err := NewViewHistory(value_object.ViewHistoryID(""), docID, userID, viewedAt)

		assert.Error(t, err)
		assert.Nil(t, vh)
		assert.Contains(t, err.Error(), "view history ID cannot be empty")
	})

	t.Run("空のDocumentIDでエラーになる", func(t *testing.T) {
		id := value_object.GenerateViewHistoryID()
		userID, _ := userVO.NewUserID("user-123")
		viewedAt := time.Now()

		vh, err := NewViewHistory(id, documentVO.DocumentID(""), userID, viewedAt)

		assert.Error(t, err)
		assert.Nil(t, vh)
		assert.Contains(t, err.Error(), "document ID cannot be empty")
	})

	t.Run("空のUserIDでエラーになる", func(t *testing.T) {
		id := value_object.GenerateViewHistoryID()
		docID := documentVO.GenerateDocumentID()
		viewedAt := time.Now()
		emptyUserID, _ := userVO.NewUserID("")

		vh, err := NewViewHistory(id, docID, emptyUserID, viewedAt)

		assert.Error(t, err)
		assert.Nil(t, vh)
		assert.Contains(t, err.Error(), "user ID cannot be empty")
	})

	t.Run("ゼロ時刻でエラーになる", func(t *testing.T) {
		id := value_object.GenerateViewHistoryID()
		docID := documentVO.GenerateDocumentID()
		userID, _ := userVO.NewUserID("user-123")

		vh, err := NewViewHistory(id, docID, userID, time.Time{})

		assert.Error(t, err)
		assert.Nil(t, vh)
		assert.Contains(t, err.Error(), "viewed at time cannot be zero")
	})
}

func TestRecordViewHistory(t *testing.T) {
	t.Run("正常に閲覧履歴を記録できる", func(t *testing.T) {
		docID := documentVO.GenerateDocumentID()
		userID, _ := userVO.NewUserID("user-123")

		vh, err := RecordViewHistory(docID, userID)

		assert.NoError(t, err)
		assert.NotNil(t, vh)
		assert.False(t, vh.ID().IsEmpty())
		assert.Equal(t, docID, vh.DocumentID())
		assert.Equal(t, userID, vh.UserID())
		assert.False(t, vh.ViewedAt().IsZero())
	})
}
