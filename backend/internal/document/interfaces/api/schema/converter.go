package schema

import (
	"opscore/backend/internal/document/application/dto"
)

// ToCreateDocumentDTO converts API schema to application DTO
func ToCreateDocumentDTO(req CreateDocumentRequest) dto.CreateDocumentRequest {
	return dto.CreateDocumentRequest{
		RepositoryID: req.RepositoryID,
		FilePath:     req.FilePath,
		CommitHash:   req.CommitHash,
		AccessScope:  req.AccessScope,
		IsAutoUpdate: req.IsAutoUpdate,
	}
}

// ToPublishDocumentDTO converts API schema to application DTO
func ToPublishDocumentDTO(req PublishDocumentRequest) dto.PublishDocumentRequest {
	return dto.PublishDocumentRequest{
		ConnectionID: req.ConnectionID,
		Owner:        req.Owner,
		Repository:   req.Repository,
		FilePath:     req.FilePath,
		AccessScope:  req.AccessScope,
		IsAutoUpdate: req.IsAutoUpdate,
	}
}

// ToUpdateDocumentDTO converts API schema to application DTO
func ToUpdateDocumentDTO(req UpdateDocumentRequest) dto.UpdateDocumentRequest {
	return dto.UpdateDocumentRequest{
		FilePath:   req.FilePath,
		CommitHash: req.CommitHash,
	}
}

// ToUpdateDocumentMetadataDTO converts API schema to application DTO
func ToUpdateDocumentMetadataDTO(req UpdateDocumentMetadataRequest) dto.UpdateDocumentMetadataRequest {
	return dto.UpdateDocumentMetadataRequest{
		Owner:        req.Owner,
		AccessScope:  req.AccessScope,
		IsAutoUpdate: req.IsAutoUpdate,
	}
}

// FromDocumentDTO converts application DTO to API schema
func FromDocumentDTO(dtoResp dto.DocumentResponse) DocumentResponse {
	var currentVersion *DocumentVersionResponse
	if dtoResp.CurrentVersion != nil {
		cv := FromDocumentVersionDTO(*dtoResp.CurrentVersion)
		currentVersion = &cv
	}

	return DocumentResponse{
		ID:             dtoResp.ID,
		RepositoryID:   dtoResp.RepositoryID,
		Owner:          dtoResp.Owner,
		IsPublished:    dtoResp.IsPublished,
		IsAutoUpdate:   dtoResp.IsAutoUpdate,
		AccessScope:    dtoResp.AccessScope,
		CurrentVersion: currentVersion,
		VersionCount:   dtoResp.VersionCount,
		CreatedAt:      dtoResp.CreatedAt,
		UpdatedAt:      dtoResp.UpdatedAt,
	}
}

// FromDocumentVersionDTO converts application DTO to API schema
func FromDocumentVersionDTO(dtoResp dto.DocumentVersionResponse) DocumentVersionResponse {
	variables := make([]VariableDefinitionResponse, len(dtoResp.Variables))
	for i, v := range dtoResp.Variables {
		variables[i] = VariableDefinitionResponse{
			Name:         v.Name,
			Label:        v.Label,
			Description:  v.Description,
			Type:         v.Type,
			Required:     v.Required,
			DefaultValue: v.DefaultValue,
		}
	}

	return DocumentVersionResponse{
		ID:            dtoResp.ID,
		DocumentID:    dtoResp.DocumentID,
		VersionNumber: dtoResp.VersionNumber,
		FilePath:      dtoResp.FilePath,
		CommitHash:    dtoResp.CommitHash,
		Title:         dtoResp.Title,
		DocType:       dtoResp.DocType,
		Tags:          dtoResp.Tags,
		Variables:     variables,
		Content:       dtoResp.Content,
		PublishedAt:   dtoResp.PublishedAt,
		UnpublishedAt: dtoResp.UnpublishedAt,
		IsCurrent:     dtoResp.IsCurrent,
	}
}

// FromDocumentListItemDTO converts application DTO to API schema
func FromDocumentListItemDTO(dtoResp dto.DocumentListItemResponse) DocumentListItemResponse {
	return DocumentListItemResponse{
		ID:           dtoResp.ID,
		RepositoryID: dtoResp.RepositoryID,
		Title:        dtoResp.Title,
		Owner:        dtoResp.Owner,
		DocType:      dtoResp.DocType,
		Tags:         dtoResp.Tags,
		IsPublished:  dtoResp.IsPublished,
		VersionCount: dtoResp.VersionCount,
		CreatedAt:    dtoResp.CreatedAt,
		UpdatedAt:    dtoResp.UpdatedAt,
	}
}

// FromDocumentListDTO converts a slice of application DTOs to API schemas
func FromDocumentListDTO(dtoResp []dto.DocumentListItemResponse) []DocumentListItemResponse {
	result := make([]DocumentListItemResponse, len(dtoResp))
	for i, d := range dtoResp {
		result[i] = FromDocumentListItemDTO(d)
	}
	return result
}

// FromVersionHistoryDTO converts application DTO to API schema
func FromVersionHistoryDTO(dtoResp dto.VersionHistoryResponse) VersionHistoryResponse {
	versions := make([]DocumentVersionResponse, len(dtoResp.Versions))
	for i, v := range dtoResp.Versions {
		versions[i] = FromDocumentVersionDTO(v)
	}
	return VersionHistoryResponse{
		DocumentID: dtoResp.DocumentID,
		Versions:   versions,
	}
}
