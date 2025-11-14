package schema

import (
	"opscore/backend/internal/git_repository/application/dto"
)

// ToRegisterRepositoryDTO converts API schema to application DTO
func ToRegisterRepositoryDTO(req RegisterRepositoryRequest) dto.RegisterRepositoryRequest {
	return dto.RegisterRepositoryRequest{
		URL:         req.URL,
		AccessToken: req.AccessToken,
	}
}

// ToUpdateAccessTokenDTO converts API schema to application DTO
func ToUpdateAccessTokenDTO(req UpdateAccessTokenRequest) dto.UpdateAccessTokenRequest {
	return dto.UpdateAccessTokenRequest{
		AccessToken: req.AccessToken,
	}
}

// ToSelectFilesDTO converts API schema to application DTO
func ToSelectFilesDTO(req SelectFilesRequest) dto.SelectFilesRequest {
	return dto.SelectFilesRequest{
		FilePaths: req.FilePaths,
	}
}

// FromRepositoryDTO converts application DTO to API schema
func FromRepositoryDTO(dtoResp dto.RepositoryResponse) RepositoryResponse {
	return RepositoryResponse{
		ID:        dtoResp.ID,
		Name:      dtoResp.Name,
		URL:       dtoResp.URL,
		CreatedAt: dtoResp.CreatedAt,
		UpdatedAt: dtoResp.UpdatedAt,
	}
}

// FromRepositoryListDTO converts application DTO list to API schema list
func FromRepositoryListDTO(dtoList []dto.RepositoryResponse) []RepositoryResponse {
	schemas := make([]RepositoryResponse, 0, len(dtoList))
	for _, dtoResp := range dtoList {
		schemas = append(schemas, FromRepositoryDTO(dtoResp))
	}
	return schemas
}

// FromFileNodeDTO converts application DTO to API schema
func FromFileNodeDTO(dtoNode dto.FileNode) FileNode {
	return FileNode{
		Path: dtoNode.Path,
		Type: dtoNode.Type,
	}
}

// FromFileNodeListDTO converts application DTO list to API schema list
func FromFileNodeListDTO(dtoList []dto.FileNode) []FileNode {
	schemas := make([]FileNode, 0, len(dtoList))
	for _, dtoNode := range dtoList {
		schemas = append(schemas, FromFileNodeDTO(dtoNode))
	}
	return schemas
}
