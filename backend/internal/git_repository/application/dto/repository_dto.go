package dto

import "time"

// RegisterRepositoryRequest represents the use case request for registering a repository
type RegisterRepositoryRequest struct {
	URL         string
	AccessToken string // Optional access token for private repositories
}

// UpdateAccessTokenRequest represents the use case request for updating a repository's access token
type UpdateAccessTokenRequest struct {
	AccessToken string
}

// RepositoryResponse represents the use case response for a repository
type RepositoryResponse struct {
	ID        string
	Name      string
	URL       string
	CreatedAt time.Time
	UpdatedAt time.Time
}
