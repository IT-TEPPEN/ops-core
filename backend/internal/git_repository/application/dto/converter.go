package dto

import (
	"opscore/backend/internal/git_repository/domain/entity"
)

// ToRepositoryResponse converts a domain Repository entity to RepositoryResponse DTO
func ToRepositoryResponse(repo entity.Repository) RepositoryResponse {
	return RepositoryResponse{
		ID:        repo.ID(),
		Name:      repo.Name(),
		URL:       repo.URL(),
		CreatedAt: repo.CreatedAt(),
		UpdatedAt: repo.UpdatedAt(),
	}
}

// ToRepositoryResponseList converts a slice of domain Repository entities to a slice of RepositoryResponse DTOs
func ToRepositoryResponseList(repos []entity.Repository) []RepositoryResponse {
	responses := make([]RepositoryResponse, 0, len(repos))
	for _, repo := range repos {
		responses = append(responses, ToRepositoryResponse(repo))
	}
	return responses
}

// ToFileNode converts a domain FileNode entity to FileNode DTO
func ToFileNode(domainFile entity.FileNode) FileNode {
	return FileNode{
		Path: domainFile.Path(),
		Type: domainFile.Type(),
	}
}

// ToFileNodeList converts a slice of domain FileNode entities to a slice of FileNode DTOs
func ToFileNodeList(domainFiles []entity.FileNode) []FileNode {
	files := make([]FileNode, 0, len(domainFiles))
	for _, df := range domainFiles {
		files = append(files, ToFileNode(df))
	}
	return files
}
