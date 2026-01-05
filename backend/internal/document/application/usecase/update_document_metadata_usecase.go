package usecase

import (
	"context"
	"fmt"

	"opscore/backend/internal/document/application/dto"
	apperror "opscore/backend/internal/document/application/error"
	"opscore/backend/internal/document/domain/repository"
	"opscore/backend/internal/document/domain/value_object"
)

// updateDocumentMetadataUseCase handles metadata updates.
type updateDocumentMetadataUseCase struct {
	repo repository.DocumentRepository
}

func newUpdateDocumentMetadataUseCase(repo repository.DocumentRepository) updateDocumentMetadataUseCase {
	return updateDocumentMetadataUseCase{repo: repo}
}

func (uc updateDocumentMetadataUseCase) Execute(ctx context.Context, documentID string, req *dto.UpdateDocumentMetadataRequest) (*dto.DocumentResponse, error) {
	docID, err := value_object.NewDocumentID(documentID)
	if err != nil {
		return nil, apperror.NewValidationFailedError([]apperror.FieldError{
			{Field: "document_id", Message: err.Error()},
		})
	}

	doc, err := uc.repo.FindByID(ctx, docID)
	if err != nil {
		return nil, fmt.Errorf("failed to find document: %w", err)
	}
	if doc == nil {
		return nil, apperror.NewNotFoundError("Document", documentID, nil)
	}

	if req.AccessScope != nil {
		accessScope, scopeErr := value_object.NewAccessScope(*req.AccessScope)
		if scopeErr != nil {
			return nil, apperror.NewValidationFailedError([]apperror.FieldError{
				{Field: "access_scope", Message: scopeErr.Error()},
			})
		}
		if err := doc.UpdateAccessScope(accessScope); err != nil {
			return nil, fmt.Errorf("failed to update access scope: %w", err)
		}
	}

	if req.IsAutoUpdate != nil {
		if *req.IsAutoUpdate {
			doc.EnableAutoUpdate()
		} else {
			doc.DisableAutoUpdate()
		}
	}

	if err := uc.repo.Update(ctx, doc); err != nil {
		return nil, fmt.Errorf("failed to update document: %w", err)
	}

	response := dto.ToDocumentResponse(doc)
	return &response, nil
}
