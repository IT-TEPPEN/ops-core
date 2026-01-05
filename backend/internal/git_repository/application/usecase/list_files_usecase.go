package repository

import (
	"context"
	"fmt"

	"opscore/backend/internal/git_repository/domain/entity"
	"opscore/backend/internal/git_repository/domain/repository"
	"opscore/backend/internal/git_repository/infrastructure/git"
)

// listFilesUseCase lists files in a repository.
type listFilesUseCase struct {
	repo          repository.Repository
	gitManager    git.GitManager
	oauthProvider OAuthTokenProvider
}

func newListFilesUseCase(repo repository.Repository, gitManager git.GitManager, oauthProvider OAuthTokenProvider) listFilesUseCase {
	return listFilesUseCase{repo: repo, gitManager: gitManager, oauthProvider: oauthProvider}
}

func (uc listFilesUseCase) Execute(ctx context.Context, repoID string, userID string) ([]entity.FileNode, error) {
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

	files, err := uc.gitManager.ListRepositoryFiles(ctx, "", repoEntity)
	if err != nil {
		return nil, fmt.Errorf("failed to list repository files: %w", err)
	}

	fileNodes := make([]entity.FileNode, 0, len(files))
	for _, f := range files {
		fileNodes = append(fileNodes, entity.NewFileNode(f, "file"))
	}

	return fileNodes, nil
}
