package usecase

import (
	"context"
	"fmt"

	"opscore/backend/internal/document/application/dto"
	apperror "opscore/backend/internal/document/application/error"
	"opscore/backend/internal/document/domain/repository"
	"opscore/backend/internal/document/domain/value_object"
)

// getDocumentVersionUseCase handles retrieving a specific document version.
type getDocumentVersionUseCase struct {
	repo repository.DocumentRepository
}

func newGetDocumentVersionUseCase(repo repository.DocumentRepository) getDocumentVersionUseCase {
	return getDocumentVersionUseCase{repo: repo}
}

func (uc getDocumentVersionUseCase) Execute(ctx context.Context, documentID string, versionNumber int) (*dto.DocumentVersionResponse, error) {
	docID, err := value_object.NewDocumentID(documentID)
	if err != nil {
		return nil, apperror.NewValidationFailedError([]apperror.FieldError{{Field: "document_id", Message: err.Error()}})
	}

	verNum, err := value_object.NewVersionNumber(versionNumber)
	if err != nil {
		return nil, apperror.NewValidationFailedError([]apperror.FieldError{{Field: "version_number", Message: err.Error()}})
	}

	version, err := uc.repo.FindVersionByNumber(ctx, docID, verNum)
	if err != nil {
		return nil, fmt.Errorf("failed to find document version: %w", err)
	}
	if version == nil {
		return nil, apperror.NewNotFoundError("DocumentVersion", fmt.Sprintf("%s@v%d", documentID, versionNumber), nil)
	}

	response := dto.ToDocumentVersionResponse(version)
	return &response, nil
}
