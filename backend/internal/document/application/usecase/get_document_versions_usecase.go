package usecase

import (
	"context"
	"fmt"

	"opscore/backend/internal/document/application/dto"
	apperror "opscore/backend/internal/document/application/error"
	"opscore/backend/internal/document/domain/repository"
	"opscore/backend/internal/document/domain/value_object"
)

// getDocumentVersionsUseCase handles retrieving all versions for a document.
type getDocumentVersionsUseCase struct {
	repo repository.DocumentRepository
}

func newGetDocumentVersionsUseCase(repo repository.DocumentRepository) getDocumentVersionsUseCase {
	return getDocumentVersionsUseCase{repo: repo}
}

func (uc getDocumentVersionsUseCase) Execute(ctx context.Context, documentID string) (*dto.VersionHistoryResponse, error) {
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

	versions, err := uc.repo.FindVersionsByDocumentID(ctx, docID)
	if err != nil {
		return nil, fmt.Errorf("failed to find document versions: %w", err)
	}

	response := dto.ToVersionHistoryResponse(documentID, versions)
	return &response, nil
}
