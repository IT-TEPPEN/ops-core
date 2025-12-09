package entity

import (
	"testing"
	"time"

	documentVO "opscore/backend/internal/document/domain/value_object"
	"opscore/backend/internal/view_statistics/domain/value_object"

	"github.com/stretchr/testify/assert"
)

func TestNewViewStatistics(t *testing.T) {
	t.Run("正常にViewStatisticsを作成できる", func(t *testing.T) {
		id := value_object.GenerateViewStatisticsID()
		docID := documentVO.GenerateDocumentID()

		stats, err := NewViewStatistics(id, docID)

		assert.NoError(t, err)
		assert.NotNil(t, stats)
		assert.Equal(t, id, stats.ID())
		assert.Equal(t, docID, stats.DocumentID())
		assert.Equal(t, int64(0), stats.TotalViews())
		assert.Equal(t, int64(0), stats.UniqueViewers())
		assert.True(t, stats.LastViewedAt().IsZero())
		assert.Equal(t, 0, stats.AverageViewDuration())
		assert.False(t, stats.CreatedAt().IsZero())
		assert.False(t, stats.UpdatedAt().IsZero())
	})

	t.Run("空のIDでエラーになる", func(t *testing.T) {
		docID := documentVO.GenerateDocumentID()

		stats, err := NewViewStatistics(value_object.ViewStatisticsID(""), docID)

		assert.Error(t, err)
		assert.Nil(t, stats)
		assert.Contains(t, err.Error(), "view statistics ID cannot be empty")
	})

	t.Run("空のDocumentIDでエラーになる", func(t *testing.T) {
		id := value_object.GenerateViewStatisticsID()

		stats, err := NewViewStatistics(id, documentVO.DocumentID(""))

		assert.Error(t, err)
		assert.Nil(t, stats)
		assert.Contains(t, err.Error(), "document ID cannot be empty")
	})
}

func TestViewStatistics_IncrementView(t *testing.T) {
	t.Run("通常の閲覧でカウントを増やせる", func(t *testing.T) {
		id := value_object.GenerateViewStatisticsID()
		docID := documentVO.GenerateDocumentID()
		stats, _ := NewViewStatistics(id, docID)

		err := stats.IncrementView(false)

		assert.NoError(t, err)
		assert.Equal(t, int64(1), stats.TotalViews())
		assert.Equal(t, int64(0), stats.UniqueViewers())
	})

	t.Run("ユニーク閲覧者でカウントを増やせる", func(t *testing.T) {
		id := value_object.GenerateViewStatisticsID()
		docID := documentVO.GenerateDocumentID()
		stats, _ := NewViewStatistics(id, docID)

		err := stats.IncrementView(true)

		assert.NoError(t, err)
		assert.Equal(t, int64(1), stats.TotalViews())
		assert.Equal(t, int64(1), stats.UniqueViewers())
	})

	t.Run("複数回の閲覧でカウントを増やせる", func(t *testing.T) {
		id := value_object.GenerateViewStatisticsID()
		docID := documentVO.GenerateDocumentID()
		stats, _ := NewViewStatistics(id, docID)

		_ = stats.IncrementView(true)
		_ = stats.IncrementView(false)
		_ = stats.IncrementView(false)

		assert.Equal(t, int64(3), stats.TotalViews())
		assert.Equal(t, int64(1), stats.UniqueViewers())
	})
}

func TestViewStatistics_UpdateLastViewedAt(t *testing.T) {
	t.Run("最終閲覧日時を更新できる", func(t *testing.T) {
		id := value_object.GenerateViewStatisticsID()
		docID := documentVO.GenerateDocumentID()
		stats, _ := NewViewStatistics(id, docID)
		viewedAt := time.Now()

		err := stats.UpdateLastViewedAt(viewedAt)

		assert.NoError(t, err)
		assert.Equal(t, viewedAt, stats.LastViewedAt())
	})

	t.Run("ゼロ時刻でエラーになる", func(t *testing.T) {
		id := value_object.GenerateViewStatisticsID()
		docID := documentVO.GenerateDocumentID()
		stats, _ := NewViewStatistics(id, docID)

		err := stats.UpdateLastViewedAt(time.Time{})

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "viewed at time cannot be zero")
	})
}

func TestViewStatistics_UpdateAverageViewDuration(t *testing.T) {
	t.Run("平均閲覧時間を更新できる", func(t *testing.T) {
		id := value_object.GenerateViewStatisticsID()
		docID := documentVO.GenerateDocumentID()
		stats, _ := NewViewStatistics(id, docID)

		err := stats.UpdateAverageViewDuration(120)

		assert.NoError(t, err)
		assert.Equal(t, 120, stats.AverageViewDuration())
	})

	t.Run("負の値でエラーになる", func(t *testing.T) {
		id := value_object.GenerateViewStatisticsID()
		docID := documentVO.GenerateDocumentID()
		stats, _ := NewViewStatistics(id, docID)

		err := stats.UpdateAverageViewDuration(-10)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "duration cannot be negative")
	})
}
