package main

import (
	"context"
	"sort"
	"sync"
	"time"

	documentVO "opscore/backend/internal/document/domain/value_object"
	userVO "opscore/backend/internal/user/domain/value_object"
	"opscore/backend/internal/view_statistics/domain/entity"
	"opscore/backend/internal/view_statistics/domain/repository"
)

// InMemoryViewStatisticsRepository is an in-memory implementation of ViewStatisticsRepository for development.
type InMemoryViewStatisticsRepository struct {
	stats map[string]entity.ViewStatistics
	mu    sync.RWMutex
}

// NewInMemoryViewStatisticsRepository creates a new InMemoryViewStatisticsRepository.
func NewInMemoryViewStatisticsRepository() repository.ViewStatisticsRepository {
	return &InMemoryViewStatisticsRepository{
		stats: make(map[string]entity.ViewStatistics),
	}
}

// Save saves view statistics.
func (r *InMemoryViewStatisticsRepository) Save(ctx context.Context, stats entity.ViewStatistics) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.stats[stats.DocumentID().String()] = stats
	return nil
}

// FindByDocumentID retrieves view statistics for a specific document.
func (r *InMemoryViewStatisticsRepository) FindByDocumentID(ctx context.Context, documentID documentVO.DocumentID) (entity.ViewStatistics, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	stat, exists := r.stats[documentID.String()]
	if !exists {
		return nil, nil
	}

	return stat, nil
}

// FindPopularDocuments retrieves the most popular documents.
func (r *InMemoryViewStatisticsRepository) FindPopularDocuments(ctx context.Context, limit int, since time.Time) ([]repository.PopularDocument, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var results []repository.PopularDocument
	for _, stat := range r.stats {
		if stat.LastViewedAt().After(since) {
			results = append(results, repository.PopularDocument{
				DocumentID:    stat.DocumentID(),
				TotalViews:    stat.TotalViews(),
				UniqueViewers: stat.UniqueViewers(),
				LastViewedAt:  stat.LastViewedAt(),
			})
		}
	}

	// Sort by total views descending
	sort.Slice(results, func(i, j int) bool {
		return results[i].TotalViews > results[j].TotalViews
	})

	if limit > 0 && len(results) > limit {
		results = results[:limit]
	}

	return results, nil
}

// FindRecentlyViewedDocuments retrieves recently viewed documents.
func (r *InMemoryViewStatisticsRepository) FindRecentlyViewedDocuments(ctx context.Context, limit int) ([]repository.PopularDocument, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var results []repository.PopularDocument
	for _, stat := range r.stats {
		results = append(results, repository.PopularDocument{
			DocumentID:    stat.DocumentID(),
			TotalViews:    stat.TotalViews(),
			UniqueViewers: stat.UniqueViewers(),
			LastViewedAt:  stat.LastViewedAt(),
		})
	}

	// Sort by last viewed at descending
	sort.Slice(results, func(i, j int) bool {
		return results[i].LastViewedAt.After(results[j].LastViewedAt)
	})

	if limit > 0 && len(results) > limit {
		results = results[:limit]
	}

	return results, nil
}

// GetUserViewCount returns the total number of views by a user.
func (r *InMemoryViewStatisticsRepository) GetUserViewCount(ctx context.Context, userID userVO.UserID) (int64, error) {
	// This in-memory implementation doesn't track user view counts
	// Would need a separate data structure or actual view history repository integration
	return 0, nil
}

// GetUserUniqueDocumentCount returns the number of unique documents viewed by a user.
func (r *InMemoryViewStatisticsRepository) GetUserUniqueDocumentCount(ctx context.Context, userID userVO.UserID) (int64, error) {
	// This in-memory implementation doesn't track user document counts
	// Would need a separate data structure or actual view history repository integration
	return 0, nil
}
