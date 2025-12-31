package service

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"opscore/backend/internal/oauth/domain"
)

// GitProviderService provides Git provider specific operations using OAuth tokens
type GitProviderService struct {
	oauthService *OAuthService
	httpClient   *http.Client
}

// NewGitProviderService creates a new GitProviderService
func NewGitProviderService(oauthService *OAuthService) *GitProviderService {
	return &GitProviderService{
		oauthService: oauthService,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// GitRepositoryOwner represents the owner of a repository
type GitRepositoryOwner struct {
	Login     string `json:"login"`
	AvatarURL string `json:"avatarUrl"`
}

// GitRepository represents a repository from a Git provider
type GitRepository struct {
	ID            int64              `json:"id"`
	Name          string             `json:"name"`
	FullName      string             `json:"fullName"`
	Description   string             `json:"description"`
	Private       bool               `json:"private"`
	HTMLURL       string             `json:"htmlUrl"`
	CloneURL      string             `json:"cloneUrl"`
	DefaultBranch string             `json:"defaultBranch"`
	Provider      string             `json:"provider"`
	Owner         GitRepositoryOwner `json:"owner"`
}

// ListUserRepositories lists repositories accessible to the user for a given provider
func (s *GitProviderService) ListUserRepositories(ctx context.Context, userID string, provider domain.Provider) ([]GitRepository, error) {
	// Get access token for the user
	token, err := s.oauthService.GetAccessToken(ctx, userID, provider)
	if err != nil {
		return nil, fmt.Errorf("failed to get access token: %w", err)
	}

	switch provider {
	case domain.ProviderGitHub:
		return s.listGitHubRepositories(ctx, token)
	case domain.ProviderGitLab, domain.ProviderGitLabSelfHosted:
		return s.listGitLabRepositories(ctx, token, "https://gitlab.com")
	default:
		return nil, fmt.Errorf("unsupported provider: %s", provider)
	}
}

// ListUserRepositoriesByConnectionID lists repositories accessible via a specific OAuth connection
func (s *GitProviderService) ListUserRepositoriesByConnectionID(ctx context.Context, userID string, connectionID string) ([]GitRepository, error) {
	// Get the OAuth connection
	conn, err := s.oauthService.GetConnectionByID(ctx, userID, connectionID)
	if err != nil {
		return nil, fmt.Errorf("failed to get connection: %w", err)
	}

	// Get access token for the connection
	token, err := s.oauthService.GetAccessTokenByConnectionID(ctx, userID, connectionID)
	if err != nil {
		return nil, fmt.Errorf("failed to get access token: %w", err)
	}

	switch conn.Provider() {
	case domain.ProviderGitHub:
		return s.listGitHubRepositories(ctx, token)
	case domain.ProviderGitLab:
		return s.listGitLabRepositories(ctx, token, "https://gitlab.com")
	case domain.ProviderGitLabSelfHosted:
		// For self-hosted GitLab, use the provider host from the connection
		baseURL := fmt.Sprintf("https://%s", conn.ProviderHost())
		return s.listGitLabRepositories(ctx, token, baseURL)
	default:
		return nil, fmt.Errorf("unsupported provider: %s", conn.Provider())
	}
}

// listGitHubRepositories lists repositories from GitHub
func (s *GitProviderService) listGitHubRepositories(ctx context.Context, accessToken string) ([]GitRepository, error) {
	var allRepos []GitRepository
	page := 1
	perPage := 100

	for {
		// visibility=all でプライベートリポジトリも含める
		// affiliation でowner, collaborator, organization_member全て含める
		url := fmt.Sprintf("https://api.github.com/user/repos?per_page=%d&page=%d&sort=updated&visibility=all&affiliation=owner,collaborator,organization_member", perPage, page)
		req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
		if err != nil {
			return nil, err
		}
		req.Header.Set("Authorization", "Bearer "+accessToken)
		req.Header.Set("Accept", "application/vnd.github+json")
		req.Header.Set("X-GitHub-Api-Version", "2022-11-28")

		resp, err := s.httpClient.Do(req)
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()

		// Debug: Log response headers to check OAuth scopes
		xOAuthScopes := resp.Header.Get("X-OAuth-Scopes")
		fmt.Printf("[DEBUG] GitHub API Response - Status: %s, X-OAuth-Scopes: %s\n", resp.Status, xOAuthScopes)

		if resp.StatusCode != http.StatusOK {
			body, _ := io.ReadAll(resp.Body)
			return nil, fmt.Errorf("github api error: %s - %s", resp.Status, string(body))
		}

		var repos []struct {
			ID            int64  `json:"id"`
			Name          string `json:"name"`
			FullName      string `json:"full_name"`
			Description   string `json:"description"`
			Private       bool   `json:"private"`
			HTMLURL       string `json:"html_url"`
			CloneURL      string `json:"clone_url"`
			DefaultBranch string `json:"default_branch"`
			Owner         struct {
				Login     string `json:"login"`
				AvatarURL string `json:"avatar_url"`
			} `json:"owner"`
		}

		if err := json.NewDecoder(resp.Body).Decode(&repos); err != nil {
			return nil, err
		}

		// Debug: Log repo count and private repos
		privateCount := 0
		for _, r := range repos {
			if r.Private {
				privateCount++
			}
		}
		fmt.Printf("[DEBUG] GitHub API - Page %d: Total repos: %d, Private repos: %d\n", page, len(repos), privateCount)

		for _, r := range repos {
			allRepos = append(allRepos, GitRepository{
				ID:            r.ID,
				Name:          r.Name,
				FullName:      r.FullName,
				Description:   r.Description,
				Private:       r.Private,
				HTMLURL:       r.HTMLURL,
				CloneURL:      r.CloneURL,
				DefaultBranch: r.DefaultBranch,
				Provider:      string(domain.ProviderGitHub),
				Owner: GitRepositoryOwner{
					Login:     r.Owner.Login,
					AvatarURL: r.Owner.AvatarURL,
				},
			})
		}

		// Check if there are more pages
		if len(repos) < perPage {
			break
		}
		page++
	}

	return allRepos, nil
}

// listGitLabRepositories lists repositories from GitLab
func (s *GitProviderService) listGitLabRepositories(ctx context.Context, accessToken string, baseURL string) ([]GitRepository, error) {
	var allRepos []GitRepository
	page := 1
	perPage := 100

	for {
		url := fmt.Sprintf("%s/api/v4/projects?membership=true&per_page=%d&page=%d&order_by=updated_at", baseURL, perPage, page)
		req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
		if err != nil {
			return nil, err
		}
		req.Header.Set("Authorization", "Bearer "+accessToken)

		resp, err := s.httpClient.Do(req)
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			body, _ := io.ReadAll(resp.Body)
			return nil, fmt.Errorf("gitlab api error: %s - %s", resp.Status, string(body))
		}

		var repos []struct {
			ID                int64  `json:"id"`
			Name              string `json:"name"`
			PathWithNamespace string `json:"path_with_namespace"`
			Description       string `json:"description"`
			Visibility        string `json:"visibility"`
			WebURL            string `json:"web_url"`
			HTTPURLToRepo     string `json:"http_url_to_repo"`
			DefaultBranch     string `json:"default_branch"`
			Namespace         struct {
				Name      string `json:"name"`
				AvatarURL string `json:"avatar_url"`
			} `json:"namespace"`
			AvatarURL string `json:"avatar_url"`
		}

		if err := json.NewDecoder(resp.Body).Decode(&repos); err != nil {
			return nil, err
		}

		for _, r := range repos {
			// GitLabではnamespaceからオーナー情報を取得
			ownerName := r.Namespace.Name
			avatarURL := r.Namespace.AvatarURL
			if avatarURL == "" {
				avatarURL = r.AvatarURL
			}

			allRepos = append(allRepos, GitRepository{
				ID:            r.ID,
				Name:          r.Name,
				FullName:      r.PathWithNamespace,
				Description:   r.Description,
				Private:       r.Visibility == "private",
				HTMLURL:       r.WebURL,
				CloneURL:      r.HTTPURLToRepo,
				DefaultBranch: r.DefaultBranch,
				Provider:      string(domain.ProviderGitLab),
				Owner: GitRepositoryOwner{
					Login:     ownerName,
					AvatarURL: avatarURL,
				},
			})
		}

		// Check if there are more pages
		if len(repos) < perPage {
			break
		}
		page++
	}

	return allRepos, nil
}
