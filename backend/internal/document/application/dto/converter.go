package dto

import (
	"opscore/backend/internal/document/domain/entity"
	"opscore/backend/internal/document/domain/value_object"
)

// ToDocumentResponse converts a domain Document to a DTO DocumentResponse
func ToDocumentResponse(doc entity.Document) DocumentResponse {
	var currentVersion *DocumentVersionResponse
	if doc.CurrentVersion() != nil {
		cv := ToDocumentVersionResponse(doc.CurrentVersion())
		currentVersion = &cv
	}

	return DocumentResponse{
		ID:             doc.ID().String(),
		RepositoryID:   doc.RepositoryID().String(),
		Owner:          doc.Owner(),
		IsPublished:    doc.IsPublished(),
		IsAutoUpdate:   doc.IsAutoUpdate(),
		AccessScope:    doc.AccessScope().String(),
		CurrentVersion: currentVersion,
		VersionCount:   len(doc.Versions()),
		CreatedAt:      doc.CreatedAt(),
		UpdatedAt:      doc.UpdatedAt(),
	}
}

// ToDocumentVersionResponse converts a domain DocumentVersion to a DTO DocumentVersionResponse
func ToDocumentVersionResponse(ver entity.DocumentVersion) DocumentVersionResponse {
	tags := make([]string, len(ver.Tags()))
	for i, t := range ver.Tags() {
		tags[i] = t.String()
	}

	variables := make([]VariableDefinitionDTO, len(ver.Variables()))
	for i, v := range ver.Variables() {
		variables[i] = VariableDefinitionDTO{
			Name:         v.Name(),
			Label:        v.Label(),
			Description:  v.Description(),
			Type:         v.Type().String(),
			Required:     v.Required(),
			DefaultValue: v.DefaultValue(),
		}
	}

	return DocumentVersionResponse{
		ID:            ver.ID().String(),
		DocumentID:    ver.DocumentID().String(),
		VersionNumber: ver.VersionNumber().Int(),
		FilePath:      ver.Source().FilePath().String(),
		CommitHash:    ver.Source().CommitHash().String(),
		Title:         ver.Title(),
		DocType:       ver.Type().String(),
		Tags:          tags,
		Variables:     variables,
		Content:       ver.Content(),
		PublishedAt:   ver.PublishedAt(),
		UnpublishedAt: ver.UnpublishedAt(),
		IsCurrent:     ver.IsCurrentVersion(),
	}
}

// ToDocumentListItemResponse converts a domain Document to a DTO DocumentListItemResponse
func ToDocumentListItemResponse(doc entity.Document) DocumentListItemResponse {
	var title, docType string
	var tags []string

	if doc.CurrentVersion() != nil {
		title = doc.CurrentVersion().Title()
		docType = doc.CurrentVersion().Type().String()
		tags = make([]string, len(doc.CurrentVersion().Tags()))
		for i, t := range doc.CurrentVersion().Tags() {
			tags[i] = t.String()
		}
	}

	return DocumentListItemResponse{
		ID:           doc.ID().String(),
		RepositoryID: doc.RepositoryID().String(),
		Title:        title,
		Owner:        doc.Owner(),
		DocType:      docType,
		Tags:         tags,
		IsPublished:  doc.IsPublished(),
		VersionCount: len(doc.Versions()),
		CreatedAt:    doc.CreatedAt(),
		UpdatedAt:    doc.UpdatedAt(),
	}
}

// ToDocumentListResponse converts a slice of domain Documents to DTO list items
func ToDocumentListResponse(docs []entity.Document) []DocumentListItemResponse {
	result := make([]DocumentListItemResponse, len(docs))
	for i, doc := range docs {
		result[i] = ToDocumentListItemResponse(doc)
	}
	return result
}

// ToVersionHistoryResponse converts a slice of domain DocumentVersions to DTO
func ToVersionHistoryResponse(docID string, versions []entity.DocumentVersion) VersionHistoryResponse {
	versionResponses := make([]DocumentVersionResponse, len(versions))
	for i, v := range versions {
		versionResponses[i] = ToDocumentVersionResponse(v)
	}
	return VersionHistoryResponse{
		DocumentID: docID,
		Versions:   versionResponses,
	}
}

// ToTagSlice converts a slice of strings to a slice of Tag value objects
func ToTagSlice(tags []string) ([]value_object.Tag, error) {
	result := make([]value_object.Tag, len(tags))
	for i, t := range tags {
		tag, err := value_object.NewTag(t)
		if err != nil {
			return nil, err
		}
		result[i] = tag
	}
	return result, nil
}

// ToVariableDefinitionSlice converts DTO variable definitions to domain value objects
func ToVariableDefinitionSlice(variables []VariableDefinitionDTO) ([]value_object.VariableDefinition, error) {
	result := make([]value_object.VariableDefinition, len(variables))
	for i, v := range variables {
		varType, err := value_object.NewVariableType(v.Type)
		if err != nil {
			return nil, err
		}
		varDef, err := value_object.NewVariableDefinition(
			v.Name,
			v.Label,
			v.Description,
			varType,
			v.Required,
			v.DefaultValue,
		)
		if err != nil {
			return nil, err
		}
		result[i] = varDef
	}
	return result, nil
}
