package dto

import "time"

// CreateDocumentRequest represents the use case request for creating a document
type CreateDocumentRequest struct {
	RepositoryID string
	FilePath     string
	CommitHash   string
	Title        string
	DocType      string   // "procedure" or "knowledge"
	Owner        string
	Tags         []string
	Variables    []VariableDefinitionDTO
	Content      string
	AccessScope  string // "public" or "private"
	IsAutoUpdate bool
}

// UpdateDocumentRequest represents the use case request for updating a document
type UpdateDocumentRequest struct {
	FilePath   string
	CommitHash string
	Title      string
	DocType    string
	Tags       []string
	Variables  []VariableDefinitionDTO
	Content    string
}

// UpdateDocumentMetadataRequest represents the use case request for updating document metadata
type UpdateDocumentMetadataRequest struct {
	Title       *string
	Owner       *string
	Tags        []string
	AccessScope *string
	IsAutoUpdate *bool
}

// VariableDefinitionDTO represents a variable definition for documents
type VariableDefinitionDTO struct {
	Name         string
	Label        string
	Description  string
	Type         string // "string", "number", "boolean", "date"
	Required     bool
	DefaultValue interface{}
}

// DocumentResponse represents the use case response for a document
type DocumentResponse struct {
	ID              string
	RepositoryID    string
	Owner           string
	IsPublished     bool
	IsAutoUpdate    bool
	AccessScope     string
	CurrentVersion  *DocumentVersionResponse
	VersionCount    int
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

// DocumentVersionResponse represents the use case response for a document version
type DocumentVersionResponse struct {
	ID            string
	DocumentID    string
	VersionNumber int
	FilePath      string
	CommitHash    string
	Title         string
	DocType       string
	Tags          []string
	Variables     []VariableDefinitionDTO
	Content       string
	PublishedAt   time.Time
	UnpublishedAt *time.Time
	IsCurrent     bool
}

// DocumentListItemResponse represents a document item in a list response
type DocumentListItemResponse struct {
	ID           string
	RepositoryID string
	Title        string
	Owner        string
	DocType      string
	Tags         []string
	IsPublished  bool
	VersionCount int
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

// VersionHistoryResponse represents the version history for a document
type VersionHistoryResponse struct {
	DocumentID string
	Versions   []DocumentVersionResponse
}

// PublishVersionRequest represents the use case request for publishing a version
type PublishVersionRequest struct {
	VersionNumber int
}

// RollbackVersionRequest represents the use case request for rollback
type RollbackVersionRequest struct {
	VersionNumber int
}
