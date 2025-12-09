package usecase

import (
	"context"
	"fmt"
	"time"

	documentVO "opscore/backend/internal/document/domain/value_object"
	userVO "opscore/backend/internal/user/domain/value_object"
	"opscore/backend/internal/view_statistics/application/dto"
	"opscore/backend/internal/view_statistics/domain/repository"
)

// ViewStatisticsUseCase defines the interface for view statistics use cases.
type ViewStatisticsUseCase interface {
	// GetDocumentStatistics retrieves statistics for a document.
	GetDocumentStatistics(ctx context.Context, documentID string) (*dto.DocumentStatisticsResponse, error)

	// GetUserStatistics retrieves statistics for a user.
	GetUserStatistics(ctx context.Context, userID string) (*dto.UserStatisticsResponse, error)

	// GetPopularDocuments retrieves popular documents ranking.
	GetPopularDocuments(ctx context.Context, limit int, days int) (*dto.PopularDocumentsResponse, error)

	// GetRecentlyViewedDocuments retrieves recently viewed documents.
	GetRecentlyViewedDocuments(ctx context.Context, limit int) (*dto.RecentDocumentsResponse, error)
}

// viewStatisticsUseCase implements the ViewStatisticsUseCase interface.
type viewStatisticsUseCase struct {
	repo repository.ViewStatisticsRepository
}

// NewViewStatisticsUseCase creates a new instance of viewStatisticsUseCase.
func NewViewStatisticsUseCase(repo repository.ViewStatisticsRepository) ViewStatisticsUseCase {
	return &viewStatisticsUseCase{
		repo: repo,
	}
}

// GetDocumentStatistics retrieves statistics for a document.
func (uc *viewStatisticsUseCase) GetDocumentStatistics(ctx context.Context, documentID string) (*dto.DocumentStatisticsResponse, error) {
	// Validate document ID
	docID, err := documentVO.NewDocumentID(documentID)
	if err != nil {
		return nil, fmt.Errorf("invalid document ID: %w", err)
	}

	// Retrieve statistics
	stats, err := uc.repo.FindByDocumentID(ctx, docID)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve document statistics: %w", err)
	}

	// Convert to response
	return &dto.DocumentStatisticsResponse{
		DocumentID:          stats.DocumentID().String(),
		TotalViews:          stats.TotalViews(),
		UniqueViewers:       stats.UniqueViewers(),
		LastViewedAt:        stats.LastViewedAt(),
		AverageViewDuration: stats.AverageViewDuration(),
	}, nil
}

// GetUserStatistics retrieves statistics for a user.
func (uc *viewStatisticsUseCase) GetUserStatistics(ctx context.Context, userID string) (*dto.UserStatisticsResponse, error) {
	// Validate user ID
	uid, err := userVO.NewUserID(userID)
	if err != nil {
		return nil, fmt.Errorf("invalid user ID: %w", err)
	}

	// Retrieve statistics
	totalViews, err := uc.repo.GetUserViewCount(ctx, uid)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve user view count: %w", err)
	}

	uniqueDocs, err := uc.repo.GetUserUniqueDocumentCount(ctx, uid)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve user unique document count: %w", err)
	}

	// Convert to response
	return &dto.UserStatisticsResponse{
		UserID:          userID,
		TotalViews:      totalViews,
		UniqueDocuments: uniqueDocs,
	}, nil
}

// GetPopularDocuments retrieves popular documents ranking.
func (uc *viewStatisticsUseCase) GetPopularDocuments(ctx context.Context, limit int, days int) (*dto.PopularDocumentsResponse, error) {
	// Set default values
	if limit <= 0 {
		limit = 10
	}
	if days <= 0 {
		days = 30
	}

	// Calculate since time
	since := time.Now().AddDate(0, 0, -days)

	// Retrieve popular documents
	popularDocs, err := uc.repo.FindPopularDocuments(ctx, limit, since)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve popular documents: %w", err)
	}

	// Convert to response
	items := make([]dto.PopularDocumentResponse, len(popularDocs))
	for i, doc := range popularDocs {
		items[i] = dto.PopularDocumentResponse{
			DocumentID:    doc.DocumentID.String(),
			TotalViews:    doc.TotalViews,
			UniqueViewers: doc.UniqueViewers,
			LastViewedAt:  doc.LastViewedAt,
		}
	}

	return &dto.PopularDocumentsResponse{
		Items: items,
		Limit: limit,
	}, nil
}

// GetRecentlyViewedDocuments retrieves recently viewed documents.
func (uc *viewStatisticsUseCase) GetRecentlyViewedDocuments(ctx context.Context, limit int) (*dto.RecentDocumentsResponse, error) {
	// Set default limit
	if limit <= 0 {
		limit = 10
	}

	// Retrieve recently viewed documents
	recentDocs, err := uc.repo.FindRecentlyViewedDocuments(ctx, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve recently viewed documents: %w", err)
	}

	// Convert to response
	items := make([]dto.RecentDocumentResponse, len(recentDocs))
	for i, doc := range recentDocs {
		items[i] = dto.RecentDocumentResponse{
			DocumentID:   doc.DocumentID.String(),
			LastViewedAt: doc.LastViewedAt,
			TotalViews:   doc.TotalViews,
		}
	}

	return &dto.RecentDocumentsResponse{
		Items: items,
		Limit: limit,
	}, nil
}
