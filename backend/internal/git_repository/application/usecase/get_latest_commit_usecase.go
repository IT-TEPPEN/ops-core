package repository

import (
	"context"
	"fmt"

	"opscore/backend/internal/git_repository/domain/entity"
	"opscore/backend/internal/git_repository/domain/repository"
	"opscore/backend/internal/git_repository/infrastructure/git"
)

// getLatestCommitUseCase retrieves the latest commit for a repository.
type getLatestCommitUseCase struct {
	repo          repository.Repository
	gitManager    git.GitManager
	oauthProvider OAuthTokenProvider
}

func newGetLatestCommitUseCase(repo repository.Repository, gitManager git.GitManager, oauthProvider OAuthTokenProvider) getLatestCommitUseCase {
	return getLatestCommitUseCase{repo: repo, gitManager: gitManager, oauthProvider: oauthProvider}
}

func (uc getLatestCommitUseCase) Execute(ctx context.Context, repoID string, userID string) (*git.CommitInfo, error) {
	repoEntity, err := loadRepository(ctx, uc.repo, repoID, "failed to retrieve repository details")
	if err != nil {
		return nil, err
	}

	repoWithToken := uc.attachTokenIfNeeded(ctx, repoEntity, userID)

	commitInfo, err := uc.gitManager.GetLatestCommit(ctx, repoWithToken)
	if err != nil {
		return nil, fmt.Errorf("failed to get latest commit: %w", err)
	}

	return commitInfo, nil
}

func (uc getLatestCommitUseCase) attachTokenIfNeeded(ctx context.Context, repoEntity entity.Repository, userID string) entity.Repository {
	if repoEntity.AccessToken() != "" {
		return repoEntity
	}
	return attachOptionalAccessToken(ctx, repoEntity, userID, uc.oauthProvider)
}
