package main

import (
	"context"
	"sync"
	"time"

	documentVO "opscore/backend/internal/document/domain/value_object"
	userVO "opscore/backend/internal/user/domain/value_object"
	"opscore/backend/internal/view_history/domain/entity"
	"opscore/backend/internal/view_history/domain/repository"
	"opscore/backend/internal/view_history/domain/value_object"
)

// InMemoryViewHistoryRepository is an in-memory implementation of ViewHistoryRepository for development.
type InMemoryViewHistoryRepository struct {
	records map[string]entity.ViewHistory
	mu      sync.RWMutex
}

// NewInMemoryViewHistoryRepository creates a new InMemoryViewHistoryRepository.
func NewInMemoryViewHistoryRepository() repository.ViewHistoryRepository {
	return &InMemoryViewHistoryRepository{
		records: make(map[string]entity.ViewHistory),
	}
}

// Save saves a view history record.
func (r *InMemoryViewHistoryRepository) Save(ctx context.Context, viewHistory entity.ViewHistory) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.records[viewHistory.ID().String()] = viewHistory
	return nil
}

// FindByID retrieves a view history record by its ID.
func (r *InMemoryViewHistoryRepository) FindByID(ctx context.Context, id value_object.ViewHistoryID) (entity.ViewHistory, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	record, exists := r.records[id.String()]
	if !exists {
		return nil, nil
	}

	return record, nil
}

// FindByUserID retrieves view history records for a specific user.
func (r *InMemoryViewHistoryRepository) FindByUserID(ctx context.Context, userID userVO.UserID, limit int, offset int) ([]entity.ViewHistory, int64, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var results []entity.ViewHistory
	for _, record := range r.records {
		if record.UserID().Equals(userID) {
			results = append(results, record)
		}
	}

	total := int64(len(results))

	// Apply offset and limit
	if offset > len(results) {
		return []entity.ViewHistory{}, total, nil
	}
	results = results[offset:]

	if limit > 0 && len(results) > limit {
		results = results[:limit]
	}

	return results, total, nil
}

// FindByDocumentID retrieves view history records for a specific document.
func (r *InMemoryViewHistoryRepository) FindByDocumentID(ctx context.Context, documentID documentVO.DocumentID, limit int, offset int) ([]entity.ViewHistory, int64, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var results []entity.ViewHistory
	for _, record := range r.records {
		if record.DocumentID().Equals(documentID) {
			results = append(results, record)
		}
	}

	total := int64(len(results))

	// Apply offset and limit
	if offset > len(results) {
		return []entity.ViewHistory{}, total, nil
	}
	results = results[offset:]

	if limit > 0 && len(results) > limit {
		results = results[:limit]
	}

	return results, total, nil
}

// FindByUserIDAndDocumentID retrieves view history records for a specific user and document.
func (r *InMemoryViewHistoryRepository) FindByUserIDAndDocumentID(ctx context.Context, userID userVO.UserID, documentID documentVO.DocumentID) ([]entity.ViewHistory, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var results []entity.ViewHistory
	for _, record := range r.records {
		if record.UserID().Equals(userID) && record.DocumentID().Equals(documentID) {
			results = append(results, record)
		}
	}

	return results, nil
}

// FindRecentByUserID retrieves recent view history for a user within a time range.
func (r *InMemoryViewHistoryRepository) FindRecentByUserID(ctx context.Context, userID userVO.UserID, since time.Time, limit int) ([]entity.ViewHistory, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var results []entity.ViewHistory
	for _, record := range r.records {
		if record.UserID().Equals(userID) && record.ViewedAt().After(since) {
			results = append(results, record)
		}
	}

	if limit > 0 && len(results) > limit {
		results = results[:limit]
	}

	return results, nil
}

// CountByDocumentID counts the total number of views for a document.
func (r *InMemoryViewHistoryRepository) CountByDocumentID(ctx context.Context, documentID documentVO.DocumentID) (int64, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var count int64
	for _, record := range r.records {
		if record.DocumentID().Equals(documentID) {
			count++
		}
	}

	return count, nil
}

// CountUniqueViewersByDocumentID counts unique viewers for a document.
func (r *InMemoryViewHistoryRepository) CountUniqueViewersByDocumentID(ctx context.Context, documentID documentVO.DocumentID) (int64, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	uniqueViewers := make(map[string]bool)
	for _, record := range r.records {
		if record.DocumentID().Equals(documentID) {
			uniqueViewers[record.UserID().String()] = true
		}
	}

	return int64(len(uniqueViewers)), nil
}
