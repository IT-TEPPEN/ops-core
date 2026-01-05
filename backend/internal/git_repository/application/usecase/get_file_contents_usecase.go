package repository

import (
	"context"
	"fmt"

	"opscore/backend/internal/git_repository/domain/repository"
	"opscore/backend/internal/git_repository/infrastructure/git"
)

// getFileContentsUseCase reads a file at a commit.
type getFileContentsUseCase struct {
	repo          repository.Repository
	gitManager    git.GitManager
	oauthProvider OAuthTokenProvider
}

func newGetFileContentsUseCase(repo repository.Repository, gitManager git.GitManager, oauthProvider OAuthTokenProvider) getFileContentsUseCase {
	return getFileContentsUseCase{repo: repo, gitManager: gitManager, oauthProvider: oauthProvider}
}

func (uc getFileContentsUseCase) Execute(ctx context.Context, repoID string, filePath string, userID string, commitHash string) (string, string, error) {
	repoEntity, err := loadRepository(ctx, uc.repo, repoID, "failed to retrieve repository details")
	if err != nil {
		return "", "", err
	}

	provider, err := requireProvider(repoEntity.URL())
	if err != nil {
		return "", "", err
	}

	if err := setRequiredAccessToken(ctx, repoEntity, userID, provider, uc.oauthProvider); err != nil {
		return "", "", err
	}

	contentBytes, actualCommit, err := uc.gitManager.ReadFileAtCommit(ctx, filePath, commitHash, repoEntity)
	if err != nil {
		return "", "", fmt.Errorf("failed to read content of file '%s': %w", filePath, err)
	}

	return string(contentBytes), actualCommit, nil
}
