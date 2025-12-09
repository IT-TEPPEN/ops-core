package repository

import (
	"context"
	"time"

	documentVO "opscore/backend/internal/document/domain/value_object"
	userVO "opscore/backend/internal/user/domain/value_object"
	"opscore/backend/internal/view_statistics/domain/entity"
)

// PopularDocument represents a document with its view statistics for ranking.
type PopularDocument struct {
	DocumentID    documentVO.DocumentID
	TotalViews    int64
	UniqueViewers int64
	LastViewedAt  time.Time
}

// ViewStatisticsRepository defines the interface for view statistics persistence.
type ViewStatisticsRepository interface {
	// Save saves view statistics.
	Save(ctx context.Context, stats entity.ViewStatistics) error

	// FindByDocumentID retrieves view statistics for a specific document.
	FindByDocumentID(ctx context.Context, documentID documentVO.DocumentID) (entity.ViewStatistics, error)

	// FindPopularDocuments retrieves the most popular documents.
	FindPopularDocuments(ctx context.Context, limit int, since time.Time) ([]PopularDocument, error)

	// FindRecentlyViewedDocuments retrieves recently viewed documents.
	FindRecentlyViewedDocuments(ctx context.Context, limit int) ([]PopularDocument, error)

	// GetUserViewCount returns the total number of views by a user.
	GetUserViewCount(ctx context.Context, userID userVO.UserID) (int64, error)

	// GetUserUniqueDocumentCount returns the number of unique documents viewed by a user.
	GetUserUniqueDocumentCount(ctx context.Context, userID userVO.UserID) (int64, error)
}
