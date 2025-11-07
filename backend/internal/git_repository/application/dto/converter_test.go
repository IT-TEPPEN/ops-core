package dto

import (
	"opscore/backend/internal/git_repository/domain/entity"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestToRepositoryResponse(t *testing.T) {
	// Create a test repository entity
	now := time.Now()
	repo := entity.ReconstructRepository(
		"test-id",
		"test-repo",
		"https://github.com/test/repo.git",
		"test-token",
		now,
		now,
	)

	// Convert to DTO
	dto := ToRepositoryResponse(repo)

	// Verify the conversion
	assert.Equal(t, repo.ID(), dto.ID)
	assert.Equal(t, repo.Name(), dto.Name)
	assert.Equal(t, repo.URL(), dto.URL)
	assert.Equal(t, repo.CreatedAt(), dto.CreatedAt)
	assert.Equal(t, repo.UpdatedAt(), dto.UpdatedAt)
}

func TestToRepositoryResponseList(t *testing.T) {
	// Create test repository entities
	now := time.Now()
	repo1 := entity.ReconstructRepository(
		"test-id-1",
		"test-repo-1",
		"https://github.com/test/repo1.git",
		"test-token-1",
		now,
		now,
	)
	repo2 := entity.ReconstructRepository(
		"test-id-2",
		"test-repo-2",
		"https://github.com/test/repo2.git",
		"test-token-2",
		now,
		now,
	)
	repos := []entity.Repository{repo1, repo2}

	// Convert to DTO list
	dtos := ToRepositoryResponseList(repos)

	// Verify the conversion
	assert.Len(t, dtos, 2)
	assert.Equal(t, repo1.ID(), dtos[0].ID)
	assert.Equal(t, repo1.Name(), dtos[0].Name)
	assert.Equal(t, repo2.ID(), dtos[1].ID)
	assert.Equal(t, repo2.Name(), dtos[1].Name)
}

func TestToRepositoryResponseList_EmptyList(t *testing.T) {
	// Create empty list
	var repos []entity.Repository

	// Convert to DTO list
	dtos := ToRepositoryResponseList(repos)

	// Verify the conversion returns empty list
	assert.NotNil(t, dtos)
	assert.Len(t, dtos, 0)
}

func TestToFileNode(t *testing.T) {
	// Create a test file node entity
	fileNode := entity.NewFileNode("README.md", "file")

	// Convert to DTO
	dto := ToFileNode(fileNode)

	// Verify the conversion
	assert.Equal(t, fileNode.Path(), dto.Path)
	assert.Equal(t, fileNode.Type(), dto.Type)
}

func TestToFileNodeList(t *testing.T) {
	// Create test file node entities
	fileNode1 := entity.NewFileNode("README.md", "file")
	fileNode2 := entity.NewFileNode("src/", "dir")
	fileNodes := []entity.FileNode{fileNode1, fileNode2}

	// Convert to DTO list
	dtos := ToFileNodeList(fileNodes)

	// Verify the conversion
	assert.Len(t, dtos, 2)
	assert.Equal(t, fileNode1.Path(), dtos[0].Path)
	assert.Equal(t, fileNode1.Type(), dtos[0].Type)
	assert.Equal(t, fileNode2.Path(), dtos[1].Path)
	assert.Equal(t, fileNode2.Type(), dtos[1].Type)
}

func TestToFileNodeList_EmptyList(t *testing.T) {
	// Create empty list
	var fileNodes []entity.FileNode

	// Convert to DTO list
	dtos := ToFileNodeList(fileNodes)

	// Verify the conversion returns empty list
	assert.NotNil(t, dtos)
	assert.Len(t, dtos, 0)
}
