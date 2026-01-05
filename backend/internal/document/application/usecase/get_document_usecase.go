package usecase

import (
	"context"
	"fmt"

	"opscore/backend/internal/document/application/dto"
	apperror "opscore/backend/internal/document/application/error"
	"opscore/backend/internal/document/domain/repository"
	"opscore/backend/internal/document/domain/value_object"
)

// getDocumentUseCase handles retrieving a single document.
type getDocumentUseCase struct {
	repo repository.DocumentRepository
}

func newGetDocumentUseCase(repo repository.DocumentRepository) getDocumentUseCase {
	return getDocumentUseCase{repo: repo}
}

func (uc getDocumentUseCase) Execute(ctx context.Context, documentID string) (*dto.DocumentResponse, error) {
	docID, err := value_object.NewDocumentID(documentID)
	if err != nil {
		return nil, apperror.NewValidationFailedError([]apperror.FieldError{{Field: "document_id", Message: err.Error()}})
	}

	doc, err := uc.repo.FindByID(ctx, docID)
	if err != nil {
		return nil, fmt.Errorf("failed to find document: %w", err)
	}
	if doc == nil {
		return nil, apperror.NewNotFoundError("Document", documentID, nil)
	}

	response := dto.ToDocumentResponse(doc)
	return &response, nil
}
