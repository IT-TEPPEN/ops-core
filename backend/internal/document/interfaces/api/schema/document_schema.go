package schema

import "time"

// CreateDocumentRequest represents the API request for creating a document
type CreateDocumentRequest struct {
	RepositoryID string                       `json:"repository_id" binding:"required" example:"a1b2c3d4-e5f6-7890-1234-567890abcdef"`
	FilePath     string                       `json:"file_path" binding:"required" example:"docs/backup-procedure.md"`
	CommitHash   string                       `json:"commit_hash" binding:"required" example:"abc1234567890"`
	Title        string                       `json:"title" binding:"required" example:"Database Backup Procedure"`
	DocType      string                       `json:"doc_type" binding:"required" example:"procedure"`
	Owner        string                       `json:"owner" binding:"required" example:"database-team"`
	Tags         []string                     `json:"tags" example:"[\"database\",\"backup\"]"`
	Variables    []VariableDefinitionRequest  `json:"variables"`
	Content      string                       `json:"content" binding:"required" example:"# Database Backup Procedure\n\nThis document describes..."`
	AccessScope  string                       `json:"access_scope" binding:"required" example:"public"`
	IsAutoUpdate bool                         `json:"is_auto_update" example:"true"`
}

// UpdateDocumentRequest represents the API request for updating a document
type UpdateDocumentRequest struct {
	FilePath   string                      `json:"file_path" binding:"required" example:"docs/backup-procedure.md"`
	CommitHash string                      `json:"commit_hash" binding:"required" example:"def4567890123"`
	Title      string                      `json:"title" binding:"required" example:"Database Backup Procedure v2"`
	DocType    string                      `json:"doc_type" binding:"required" example:"procedure"`
	Tags       []string                    `json:"tags" example:"[\"database\",\"backup\",\"v2\"]"`
	Variables  []VariableDefinitionRequest `json:"variables"`
	Content    string                      `json:"content" binding:"required" example:"# Database Backup Procedure v2\n\nUpdated procedure..."`
}

// UpdateDocumentMetadataRequest represents the API request for updating document metadata
type UpdateDocumentMetadataRequest struct {
	Owner        *string `json:"owner" example:"new-team"`
	AccessScope  *string `json:"access_scope" example:"private"`
	IsAutoUpdate *bool   `json:"is_auto_update" example:"false"`
}

// VariableDefinitionRequest represents a variable definition in API requests
type VariableDefinitionRequest struct {
	Name         string      `json:"name" binding:"required" example:"server_name"`
	Label        string      `json:"label" binding:"required" example:"Server Name"`
	Description  string      `json:"description" example:"The target server name"`
	Type         string      `json:"type" binding:"required" example:"string"`
	Required     bool        `json:"required" example:"true"`
	DefaultValue interface{} `json:"default_value" example:"prod-db-01"`
}

// VariableDefinitionResponse represents a variable definition in API responses
type VariableDefinitionResponse struct {
	Name         string      `json:"name" example:"server_name"`
	Label        string      `json:"label" example:"Server Name"`
	Description  string      `json:"description" example:"The target server name"`
	Type         string      `json:"type" example:"string"`
	Required     bool        `json:"required" example:"true"`
	DefaultValue interface{} `json:"default_value" example:"prod-db-01"`
}

// DocumentResponse represents the API response for a document
type DocumentResponse struct {
	ID              string                   `json:"id" example:"a1b2c3d4-e5f6-7890-1234-567890abcdef"`
	RepositoryID    string                   `json:"repository_id" example:"a1b2c3d4-e5f6-7890-1234-567890abcdef"`
	Owner           string                   `json:"owner" example:"database-team"`
	IsPublished     bool                     `json:"is_published" example:"true"`
	IsAutoUpdate    bool                     `json:"is_auto_update" example:"true"`
	AccessScope     string                   `json:"access_scope" example:"public"`
	CurrentVersion  *DocumentVersionResponse `json:"current_version"`
	VersionCount    int                      `json:"version_count" example:"3"`
	CreatedAt       time.Time                `json:"created_at" example:"2025-04-22T10:00:00Z"`
	UpdatedAt       time.Time                `json:"updated_at" example:"2025-04-22T12:00:00Z"`
}

// DocumentVersionResponse represents the API response for a document version
type DocumentVersionResponse struct {
	ID            string                       `json:"id" example:"v1b2c3d4-e5f6-7890-1234-567890abcdef"`
	DocumentID    string                       `json:"document_id" example:"a1b2c3d4-e5f6-7890-1234-567890abcdef"`
	VersionNumber int                          `json:"version_number" example:"1"`
	FilePath      string                       `json:"file_path" example:"docs/backup-procedure.md"`
	CommitHash    string                       `json:"commit_hash" example:"abc1234567890"`
	Title         string                       `json:"title" example:"Database Backup Procedure"`
	DocType       string                       `json:"doc_type" example:"procedure"`
	Tags          []string                     `json:"tags" example:"[\"database\",\"backup\"]"`
	Variables     []VariableDefinitionResponse `json:"variables"`
	Content       string                       `json:"content" example:"# Database Backup Procedure\n\nThis document describes..."`
	PublishedAt   time.Time                    `json:"published_at" example:"2025-04-22T10:00:00Z"`
	UnpublishedAt *time.Time                   `json:"unpublished_at,omitempty" example:"null"`
	IsCurrent     bool                         `json:"is_current" example:"true"`
}

// DocumentListItemResponse represents a document item in a list response
type DocumentListItemResponse struct {
	ID           string    `json:"id" example:"a1b2c3d4-e5f6-7890-1234-567890abcdef"`
	RepositoryID string    `json:"repository_id" example:"a1b2c3d4-e5f6-7890-1234-567890abcdef"`
	Title        string    `json:"title" example:"Database Backup Procedure"`
	Owner        string    `json:"owner" example:"database-team"`
	DocType      string    `json:"doc_type" example:"procedure"`
	Tags         []string  `json:"tags" example:"[\"database\",\"backup\"]"`
	IsPublished  bool      `json:"is_published" example:"true"`
	VersionCount int       `json:"version_count" example:"3"`
	CreatedAt    time.Time `json:"created_at" example:"2025-04-22T10:00:00Z"`
	UpdatedAt    time.Time `json:"updated_at" example:"2025-04-22T12:00:00Z"`
}

// ListDocumentsResponse represents the API response for listing documents
type ListDocumentsResponse struct {
	Documents []DocumentListItemResponse `json:"documents"`
}

// VersionHistoryResponse represents the API response for version history
type VersionHistoryResponse struct {
	DocumentID string                    `json:"document_id" example:"a1b2c3d4-e5f6-7890-1234-567890abcdef"`
	Versions   []DocumentVersionResponse `json:"versions"`
}

// PublishVersionRequest represents the API request for publishing a version
type PublishVersionRequest struct {
	VersionNumber int `json:"version_number" binding:"required" example:"1"`
}

// RollbackVersionRequest represents the API request for rollback
type RollbackVersionRequest struct {
	VersionNumber int `json:"version_number" binding:"required" example:"1"`
}

// ErrorResponse represents an API error response
type ErrorResponse struct {
	Code    string                 `json:"code" example:"VALIDATION_FAILED"`
	Message string                 `json:"message" example:"Validation failed"`
	Details map[string]interface{} `json:"details,omitempty"`
}
