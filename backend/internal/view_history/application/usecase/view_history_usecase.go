package usecase

import (
	"context"
	"fmt"

	documentVO "opscore/backend/internal/document/domain/value_object"
	userVO "opscore/backend/internal/user/domain/value_object"
	"opscore/backend/internal/view_history/application/dto"
	"opscore/backend/internal/view_history/domain/entity"
	"opscore/backend/internal/view_history/domain/repository"
)

// ViewHistoryUseCase defines the interface for view history use cases.
type ViewHistoryUseCase interface {
	// RecordView records a document view.
	RecordView(ctx context.Context, req *dto.RecordViewRequest) (*dto.ViewHistoryResponse, error)

	// GetViewHistory retrieves view history for a user.
	GetViewHistory(ctx context.Context, userID string, limit int, offset int) (*dto.ViewHistoryListResponse, error)

	// GetDocumentViewHistory retrieves view history for a document.
	GetDocumentViewHistory(ctx context.Context, documentID string, limit int, offset int) (*dto.ViewHistoryListResponse, error)
}

// viewHistoryUseCase implements the ViewHistoryUseCase interface.
type viewHistoryUseCase struct {
	repo repository.ViewHistoryRepository
}

// NewViewHistoryUseCase creates a new instance of viewHistoryUseCase.
func NewViewHistoryUseCase(repo repository.ViewHistoryRepository) ViewHistoryUseCase {
	return &viewHistoryUseCase{
		repo: repo,
	}
}

// RecordView records a document view.
func (uc *viewHistoryUseCase) RecordView(ctx context.Context, req *dto.RecordViewRequest) (*dto.ViewHistoryResponse, error) {
	// Validate document ID
	documentID, err := documentVO.NewDocumentID(req.DocumentID)
	if err != nil {
		return nil, fmt.Errorf("invalid document ID: %w", err)
	}

	// Validate user ID
	userID, err := userVO.NewUserID(req.UserID)
	if err != nil {
		return nil, fmt.Errorf("invalid user ID: %w", err)
	}

	// Create view history record
	viewHistory, err := entity.RecordViewHistory(documentID, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to create view history: %w", err)
	}

	// Save to repository
	if err := uc.repo.Save(ctx, viewHistory); err != nil {
		return nil, fmt.Errorf("failed to save view history: %w", err)
	}

	// Return response
	return &dto.ViewHistoryResponse{
		ID:           viewHistory.ID().String(),
		DocumentID:   viewHistory.DocumentID().String(),
		UserID:       viewHistory.UserID().String(),
		ViewedAt:     viewHistory.ViewedAt(),
		ViewDuration: viewHistory.ViewDuration(),
	}, nil
}

// GetViewHistory retrieves view history for a user.
func (uc *viewHistoryUseCase) GetViewHistory(ctx context.Context, userID string, limit int, offset int) (*dto.ViewHistoryListResponse, error) {
	// Validate user ID
	uid, err := userVO.NewUserID(userID)
	if err != nil {
		return nil, fmt.Errorf("invalid user ID: %w", err)
	}

	// Set default limit if not specified
	if limit <= 0 {
		limit = 50
	}
	if offset < 0 {
		offset = 0
	}

	// Retrieve view history
	histories, totalCount, err := uc.repo.FindByUserID(ctx, uid, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve view history: %w", err)
	}

	// Convert to response
	items := make([]dto.ViewHistoryResponse, len(histories))
	for i, h := range histories {
		items[i] = dto.ViewHistoryResponse{
			ID:           h.ID().String(),
			DocumentID:   h.DocumentID().String(),
			UserID:       h.UserID().String(),
			ViewedAt:     h.ViewedAt(),
			ViewDuration: h.ViewDuration(),
		}
	}

	return &dto.ViewHistoryListResponse{
		Items:      items,
		TotalCount: int(totalCount),
		Limit:      limit,
		Offset:     offset,
	}, nil
}

// GetDocumentViewHistory retrieves view history for a document.
func (uc *viewHistoryUseCase) GetDocumentViewHistory(ctx context.Context, documentID string, limit int, offset int) (*dto.ViewHistoryListResponse, error) {
	// Validate document ID
	docID, err := documentVO.NewDocumentID(documentID)
	if err != nil {
		return nil, fmt.Errorf("invalid document ID: %w", err)
	}

	// Set default limit if not specified
	if limit <= 0 {
		limit = 50
	}
	if offset < 0 {
		offset = 0
	}

	// Retrieve view history
	histories, totalCount, err := uc.repo.FindByDocumentID(ctx, docID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve document view history: %w", err)
	}

	// Convert to response
	items := make([]dto.ViewHistoryResponse, len(histories))
	for i, h := range histories {
		items[i] = dto.ViewHistoryResponse{
			ID:           h.ID().String(),
			DocumentID:   h.DocumentID().String(),
			UserID:       h.UserID().String(),
			ViewedAt:     h.ViewedAt(),
			ViewDuration: h.ViewDuration(),
		}
	}

	return &dto.ViewHistoryListResponse{
		Items:      items,
		TotalCount: int(totalCount),
		Limit:      limit,
		Offset:     offset,
	}, nil
}
