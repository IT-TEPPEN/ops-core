package repository

import (
	"context"

	"opscore/backend/internal/git_repository/domain/entity"
	"opscore/backend/internal/git_repository/domain/repository"
)

// getRepositoryUseCase retrieves a repository by ID.
type getRepositoryUseCase struct {
	repo repository.Repository
}

func newGetRepositoryUseCase(repo repository.Repository) getRepositoryUseCase {
	return getRepositoryUseCase{repo: repo}
}

func (uc getRepositoryUseCase) Execute(ctx context.Context, repoID string) (entity.Repository, error) {
	return loadRepository(ctx, uc.repo, repoID, "failed to retrieve repository")
}
