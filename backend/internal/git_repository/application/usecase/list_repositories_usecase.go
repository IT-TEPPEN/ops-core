package repository

import (
	"context"
	"fmt"

	"opscore/backend/internal/git_repository/domain/entity"
	"opscore/backend/internal/git_repository/domain/repository"
)

// listRepositoriesUseCase retrieves all repositories.
type listRepositoriesUseCase struct {
	repo repository.Repository
}

func newListRepositoriesUseCase(repo repository.Repository) listRepositoriesUseCase {
	return listRepositoriesUseCase{repo: repo}
}

func (uc listRepositoriesUseCase) Execute(ctx context.Context) ([]entity.Repository, error) {
	repos, err := uc.repo.FindAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve repositories: %w", err)
	}
	return repos, nil
}
