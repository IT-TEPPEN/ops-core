package repository

import (
	"context"
	"fmt"

	apperror "opscore/backend/internal/git_repository/application/error"
	"opscore/backend/internal/git_repository/domain/entity"
	"opscore/backend/internal/git_repository/domain/repository"
)

func loadRepository(ctx context.Context, store repository.Repository, repoID string, wrapMsg string) (entity.Repository, error) {
	repo, err := store.FindByID(ctx, repoID)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", wrapMsg, err)
	}
	if repo == nil {
		return nil, apperror.NewNotFoundError("Repository", repoID, nil)
	}
	return repo, nil
}
