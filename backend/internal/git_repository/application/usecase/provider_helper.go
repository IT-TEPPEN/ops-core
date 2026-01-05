package repository

import (
	"strings"

	apperror "opscore/backend/internal/git_repository/application/error"
)

func getProviderFromURL(repoURL string) string {
	if strings.Contains(repoURL, "github.com") {
		return "github"
	}
	if strings.Contains(repoURL, "gitlab.com") {
		return "gitlab"
	}
	return ""
}

func requireProvider(repoURL string) (string, error) {
	provider := getProviderFromURL(repoURL)
	if provider == "" {
		return "", apperror.NewValidationFailedError([]apperror.FieldError{{
			Field:   "url",
			Message: "unsupported Git provider",
		}})
	}
	return provider, nil
}
