package usecase

import (
	"context"
	"fmt"

	"opscore/backend/internal/document/application/dto"
	apperror "opscore/backend/internal/document/application/error"
	"opscore/backend/internal/document/domain/entity"
	"opscore/backend/internal/document/domain/repository"
	"opscore/backend/internal/document/domain/value_object"
)

// listDocumentsUseCase handles listing published documents with filters.
type listDocumentsUseCase struct {
	repo repository.DocumentRepository
}

func newListDocumentsUseCase(repo repository.DocumentRepository) listDocumentsUseCase {
	return listDocumentsUseCase{repo: repo}
}

func (uc listDocumentsUseCase) Execute(ctx context.Context, filter dto.DocumentListFilter) ([]dto.DocumentListItemResponse, error) {
	filters := make([]repository.Filter, 0, 4)

	if filter.RepositoryID != "" {
		repoID, vErr := value_object.NewRepositoryID(filter.RepositoryID)
		if vErr != nil {
			return nil, apperror.NewValidationFailedError([]apperror.FieldError{{Field: "repository_id", Message: vErr.Error()}})
		}
		filters = append(filters, repository.NewRepositoryIDFilter(repoID.String()))
	}
	if filter.ProviderRepositoryID != "" {
		filters = append(filters, repository.NewProviderRepositoryIDFilter(filter.ProviderRepositoryID))
	}
	if filter.Owner != "" {
		filters = append(filters, repository.NewOwnerFilter(filter.Owner))
	}
	if filter.Repository != "" {
		filters = append(filters, repository.NewRepositoryNameFilter(filter.Repository))
	}

	docs, err := uc.repo.FindPublished(ctx, filters...)
	if err != nil {
		return nil, fmt.Errorf("failed to find documents: %w", err)
	}

	// In-memory filter as safety net in case repository ignores optional filters
	filtered := make([]entity.Document, 0, len(docs))
	for _, d := range docs {
		if filter.RepositoryID != "" && d.RepositoryID().String() != filter.RepositoryID {
			continue
		}
		if filter.ProviderRepositoryID != "" {
			if d.Origin() == nil || d.Origin().ProviderRepositoryID().String() != filter.ProviderRepositoryID {
				continue
			}
		}
		if filter.Owner != "" {
			if d.Origin() == nil || d.Origin().Owner() != filter.Owner {
				continue
			}
		}
		if filter.Repository != "" {
			if d.Origin() == nil || d.Origin().Repository() != filter.Repository {
				continue
			}
		}
		filtered = append(filtered, d)
	}

	return dto.ToDocumentListResponse(filtered), nil
}
