package repository

import (
	"context"
	"time"

	documentVO "opscore/backend/internal/document/domain/value_object"
	userVO "opscore/backend/internal/user/domain/value_object"
	"opscore/backend/internal/view_history/domain/entity"
	"opscore/backend/internal/view_history/domain/value_object"
)

// ViewHistoryRepository defines the interface for view history persistence.
type ViewHistoryRepository interface {
	// Save saves a view history record.
	Save(ctx context.Context, viewHistory entity.ViewHistory) error

	// FindByID retrieves a view history record by its ID.
	FindByID(ctx context.Context, id value_object.ViewHistoryID) (entity.ViewHistory, error)

	// FindByUserID retrieves view history records for a specific user.
	FindByUserID(ctx context.Context, userID userVO.UserID, limit int, offset int) ([]entity.ViewHistory, int64, error)

	// FindByDocumentID retrieves view history records for a specific document.
	FindByDocumentID(ctx context.Context, documentID documentVO.DocumentID, limit int, offset int) ([]entity.ViewHistory, int64, error)

	// FindByUserIDAndDocumentID retrieves view history records for a specific user and document.
	FindByUserIDAndDocumentID(ctx context.Context, userID userVO.UserID, documentID documentVO.DocumentID) ([]entity.ViewHistory, error)

	// FindRecentByUserID retrieves recent view history for a user within a time range.
	FindRecentByUserID(ctx context.Context, userID userVO.UserID, since time.Time, limit int) ([]entity.ViewHistory, error)

	// CountByDocumentID counts the total number of views for a document.
	CountByDocumentID(ctx context.Context, documentID documentVO.DocumentID) (int64, error)

	// CountUniqueViewersByDocumentID counts unique viewers for a document.
	CountUniqueViewersByDocumentID(ctx context.Context, documentID documentVO.DocumentID) (int64, error)
}
