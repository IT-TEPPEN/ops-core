package repository

import (
	"context"
	"fmt"

	apperror "opscore/backend/internal/git_repository/application/error"
	"opscore/backend/internal/git_repository/domain/entity"
)

func setRequiredAccessToken(ctx context.Context, repo entity.Repository, userID string, provider string, oauth OAuthTokenProvider) error {
	accessToken, err := oauth.GetAccessTokenForProvider(ctx, userID, provider)
	if err != nil {
		return apperror.NewValidationFailedError([]apperror.FieldError{{
			Field:   "oauth",
			Message: fmt.Sprintf("OAuth connection required for %s. Please connect your account.", provider),
		}})
	}

	repo.SetAccessToken(accessToken)
	return nil
}

func attachOptionalAccessToken(ctx context.Context, repo entity.Repository, userID string, oauth OAuthTokenProvider) entity.Repository {
	if repo.AccessToken() != "" || userID == "" {
		return repo
	}

	provider := getProviderFromURL(repo.URL())
	if provider == "" {
		return repo
	}

	token, err := oauth.GetAccessTokenForProvider(ctx, userID, provider)
	if err == nil && token != "" {
		return entity.ReconstructRepository(repo.ID(), repo.Name(), repo.URL(), token, repo.CreatedAt(), repo.UpdatedAt())
	}

	return repo
}
