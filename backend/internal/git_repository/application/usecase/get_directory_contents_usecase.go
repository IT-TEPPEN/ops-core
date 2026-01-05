package repository

import (
	"context"
	"fmt"

	"opscore/backend/internal/git_repository/domain/entity"
	"opscore/backend/internal/git_repository/domain/repository"
	"opscore/backend/internal/git_repository/infrastructure/git"
)

// getDirectoryContentsUseCase lists items at a path.
type getDirectoryContentsUseCase struct {
	repo          repository.Repository
	gitManager    git.GitManager
	oauthProvider OAuthTokenProvider
}

func newGetDirectoryContentsUseCase(repo repository.Repository, gitManager git.GitManager, oauthProvider OAuthTokenProvider) getDirectoryContentsUseCase {
	return getDirectoryContentsUseCase{repo: repo, gitManager: gitManager, oauthProvider: oauthProvider}
}

func (uc getDirectoryContentsUseCase) Execute(ctx context.Context, repoID string, path string, userID string) ([]entity.FileNode, error) {
	repoEntity, err := loadRepository(ctx, uc.repo, repoID, "failed to retrieve repository details")
	if err != nil {
		return nil, err
	}

	provider, err := requireProvider(repoEntity.URL())
	if err != nil {
		return nil, err
	}

	if err := setRequiredAccessToken(ctx, repoEntity, userID, provider, uc.oauthProvider); err != nil {
		return nil, err
	}

	fileNodes, err := uc.gitManager.ListDirectoryContents(ctx, path, repoEntity)
	if err != nil {
		return nil, fmt.Errorf("failed to list directory contents at path '%s': %w", path, err)
	}

	return fileNodes, nil
}
