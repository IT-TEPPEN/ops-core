package usecase

import (
	"context"
	"fmt"

	"opscore/backend/internal/document/application/dto"
	apperror "opscore/backend/internal/document/application/error"
	"opscore/backend/internal/document/domain/repository"
	"opscore/backend/internal/document/domain/value_object"
	"opscore/backend/internal/document/infrastructure/storage"
)

// publishDocumentVersionUseCase handles publishing a specific version.
type publishDocumentVersionUseCase struct {
	repo      repository.DocumentRepository
	persister contentPersister
}

func newPublishDocumentVersionUseCase(repo repository.DocumentRepository, storage storage.DocumentStorage) publishDocumentVersionUseCase {
	return publishDocumentVersionUseCase{repo: repo, persister: newContentPersister(storage)}
}

func (uc publishDocumentVersionUseCase) Execute(ctx context.Context, documentID string, versionNumber int) (*dto.DocumentResponse, error) {
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

	if doc.IsPublished() && doc.CurrentVersion() != nil && doc.CurrentVersion().VersionNumber().Equals(verNum) {
		response := dto.ToDocumentResponse(doc)
		return &response, nil
	}

	// TODO: consider adding PublishWithVersion to domain entity instead of RollbackToVersion.
	if err := doc.RollbackToVersion(verNum); err != nil {
		return nil, fmt.Errorf("failed to publish version: %w", err)
	}

	if err := uc.repo.Update(ctx, doc); err != nil {
		return nil, fmt.Errorf("failed to update document: %w", err)
	}

	uc.persister.persistContent(ctx, doc)

	response := dto.ToDocumentResponse(doc)
	return &response, nil
}
