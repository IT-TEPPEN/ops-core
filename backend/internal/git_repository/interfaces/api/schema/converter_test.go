package schema

import (
	"opscore/backend/internal/git_repository/application/dto"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestToRegisterRepositoryDTO(t *testing.T) {
	req := RegisterRepositoryRequest{
		URL:         "https://github.com/test/repo.git",
		AccessToken: "test-token",
	}

	dtoReq := ToRegisterRepositoryDTO(req)

	assert.Equal(t, req.URL, dtoReq.URL)
	assert.Equal(t, req.AccessToken, dtoReq.AccessToken)
}

func TestToUpdateAccessTokenDTO(t *testing.T) {
	req := UpdateAccessTokenRequest{
		AccessToken: "new-token",
	}

	dtoReq := ToUpdateAccessTokenDTO(req)

	assert.Equal(t, req.AccessToken, dtoReq.AccessToken)
}

func TestToSelectFilesDTO(t *testing.T) {
	req := SelectFilesRequest{
		FilePaths: []string{"README.md", "docs/guide.md"},
	}

	dtoReq := ToSelectFilesDTO(req)

	assert.Equal(t, req.FilePaths, dtoReq.FilePaths)
}

func TestFromRepositoryDTO(t *testing.T) {
	now := time.Now()
	dtoResp := dto.RepositoryResponse{
		ID:        "test-id",
		Name:      "test-repo",
		URL:       "https://github.com/test/repo.git",
		CreatedAt: now,
		UpdatedAt: now,
	}

	schemaResp := FromRepositoryDTO(dtoResp)

	assert.Equal(t, dtoResp.ID, schemaResp.ID)
	assert.Equal(t, dtoResp.Name, schemaResp.Name)
	assert.Equal(t, dtoResp.URL, schemaResp.URL)
	assert.Equal(t, dtoResp.CreatedAt, schemaResp.CreatedAt)
	assert.Equal(t, dtoResp.UpdatedAt, schemaResp.UpdatedAt)
}

func TestFromRepositoryListDTO(t *testing.T) {
	now := time.Now()
	dtoList := []dto.RepositoryResponse{
		{
			ID:        "test-id-1",
			Name:      "test-repo-1",
			URL:       "https://github.com/test/repo1.git",
			CreatedAt: now,
			UpdatedAt: now,
		},
		{
			ID:        "test-id-2",
			Name:      "test-repo-2",
			URL:       "https://github.com/test/repo2.git",
			CreatedAt: now,
			UpdatedAt: now,
		},
	}

	schemaList := FromRepositoryListDTO(dtoList)

	assert.Len(t, schemaList, 2)
	assert.Equal(t, dtoList[0].ID, schemaList[0].ID)
	assert.Equal(t, dtoList[1].ID, schemaList[1].ID)
}

func TestFromRepositoryListDTO_EmptyList(t *testing.T) {
	var dtoList []dto.RepositoryResponse

	schemaList := FromRepositoryListDTO(dtoList)

	assert.NotNil(t, schemaList)
	assert.Len(t, schemaList, 0)
}

func TestFromFileNodeDTO(t *testing.T) {
	dtoNode := dto.FileNode{
		Path: "README.md",
		Type: "file",
	}

	schemaNode := FromFileNodeDTO(dtoNode)

	assert.Equal(t, dtoNode.Path, schemaNode.Path)
	assert.Equal(t, dtoNode.Type, schemaNode.Type)
}

func TestFromFileNodeListDTO(t *testing.T) {
	dtoList := []dto.FileNode{
		{Path: "README.md", Type: "file"},
		{Path: "src/", Type: "dir"},
	}

	schemaList := FromFileNodeListDTO(dtoList)

	assert.Len(t, schemaList, 2)
	assert.Equal(t, dtoList[0].Path, schemaList[0].Path)
	assert.Equal(t, dtoList[1].Type, schemaList[1].Type)
}

func TestFromFileNodeListDTO_EmptyList(t *testing.T) {
	var dtoList []dto.FileNode

	schemaList := FromFileNodeListDTO(dtoList)

	assert.NotNil(t, schemaList)
	assert.Len(t, schemaList, 0)
}
