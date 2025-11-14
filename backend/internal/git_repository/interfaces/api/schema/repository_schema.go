package schema

import "time"

// RegisterRepositoryRequest represents the API request body for registering a repository
type RegisterRepositoryRequest struct {
	URL         string `json:"url" binding:"required" example:"https://github.com/user/repo.git"`
	AccessToken string `json:"accessToken" example:"ghp_1234567890abcdefghijklmnopqrstuvwxyz"` // Optional access token for private repositories
}

// UpdateAccessTokenRequest represents the API request body for updating a repository's access token
type UpdateAccessTokenRequest struct {
	AccessToken string `json:"accessToken" binding:"required" example:"ghp_1234567890abcdefghijklmnopqrstuvwxyz"`
}

// RepositoryResponse represents the API response format for a repository
type RepositoryResponse struct {
	ID        string    `json:"id" example:"a1b2c3d4-e5f6-7890-1234-567890abcdef"`
	Name      string    `json:"name" example:"repo"`
	URL       string    `json:"url" example:"https://github.com/user/repo.git"`
	CreatedAt time.Time `json:"created_at" example:"2025-04-22T10:00:00Z"`
	UpdatedAt time.Time `json:"updated_at" example:"2025-04-22T10:00:00Z"`
}

// ListRepositoriesResponse represents the API response for listing all repositories
type ListRepositoriesResponse struct {
	Repositories []RepositoryResponse `json:"repositories"`
}
