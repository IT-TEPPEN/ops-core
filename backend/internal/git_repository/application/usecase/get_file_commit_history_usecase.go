package repository

import (
	"context"
	"fmt"

	"opscore/backend/internal/git_repository/domain/entity"
	"opscore/backend/internal/git_repository/domain/repository"
	"opscore/backend/internal/git_repository/infrastructure/git"
)

// getFileCommitHistoryUseCase retrieves commit history for a file.
type getFileCommitHistoryUseCase struct {
	repo          repository.Repository
	gitManager    git.GitManager
	oauthProvider OAuthTokenProvider
}

func newGetFileCommitHistoryUseCase(repo repository.Repository, gitManager git.GitManager, oauthProvider OAuthTokenProvider) getFileCommitHistoryUseCase {
	return getFileCommitHistoryUseCase{repo: repo, gitManager: gitManager, oauthProvider: oauthProvider}
}

func (uc getFileCommitHistoryUseCase) Execute(ctx context.Context, repoID string, filePath string, userID string) ([]git.CommitInfo, error) {
	repoEntity, err := loadRepository(ctx, uc.repo, repoID, "failed to retrieve repository details")
	if err != nil {
		return nil, err
	}

	repoWithToken := uc.attachTokenIfNeeded(ctx, repoEntity, userID)

	commits, err := uc.gitManager.GetFileCommitHistory(ctx, filePath, repoWithToken)
	if err != nil {
		return nil, fmt.Errorf("failed to get file commit history: %w", err)
	}

	return commits, nil
}

func (uc getFileCommitHistoryUseCase) attachTokenIfNeeded(ctx context.Context, repoEntity entity.Repository, userID string) entity.Repository {
	if repoEntity.AccessToken() != "" {
		return repoEntity
	}
	return attachOptionalAccessToken(ctx, repoEntity, userID, uc.oauthProvider)
}
