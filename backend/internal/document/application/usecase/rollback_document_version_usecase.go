package usecase

import (
	"context"
	"fmt"

	"opscore/backend/internal/document/application/dto"
	apperror "opscore/backend/internal/document/application/error"
	"opscore/backend/internal/document/domain/repository"
	"opscore/backend/internal/document/domain/value_object"
)

// rollbackDocumentVersionUseCase handles rolling back to a previous version.
type rollbackDocumentVersionUseCase struct {
	repo repository.DocumentRepository
}

func newRollbackDocumentVersionUseCase(repo repository.DocumentRepository) rollbackDocumentVersionUseCase {
	return rollbackDocumentVersionUseCase{repo: repo}
}

func (uc rollbackDocumentVersionUseCase) Execute(ctx context.Context, documentID string, versionNumber int) (*dto.DocumentResponse, error) {
	docID, err := value_object.NewDocumentID(documentID)
	if err != nil {
		return nil, apperror.NewValidationFailedError([]apperror.FieldError{{Field: "document_id", Message: err.Error()}})
	}

	verNum, err := value_object.NewVersionNumber(versionNumber)
	if err != nil {
		return nil, apperror.NewValidationFailedError([]apperror.FieldError{{Field: "version_number", Message: err.Error()}})
	}

	doc, err := uc.repo.FindByID(ctx, docID)
	if err != nil {
		return nil, fmt.Errorf("failed to find document: %w", err)
	}
	if doc == nil {
		return nil, apperror.NewNotFoundError("Document", documentID, nil)
	}

	if err := doc.RollbackToVersion(verNum); err != nil {
		return nil, fmt.Errorf("failed to rollback to version: %w", err)
	}

	if err := uc.repo.Update(ctx, doc); err != nil {
		return nil, fmt.Errorf("failed to update document: %w", err)
	}

	response := dto.ToDocumentResponse(doc)
	return &response, nil
}
