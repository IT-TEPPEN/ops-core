package schema

import "time"

// RecordViewRequest represents the request to record a document view.
type RecordViewRequest struct {
	UserID string `json:"user_id" binding:"required" example:"user-123"`
}

// ViewHistoryResponse represents a view history record.
type ViewHistoryResponse struct {
	ID           string    `json:"id" example:"a1b2c3d4-e5f6-7890-1234-567890abcdef"`
	DocumentID   string    `json:"document_id" example:"b2c3d4e5-f6a7-8901-2345-678901bcdefg"`
	UserID       string    `json:"user_id" example:"user-123"`
	ViewedAt     time.Time `json:"viewed_at" example:"2024-01-15T10:30:00Z"`
	ViewDuration int       `json:"view_duration" example:"120"`
}

// ViewHistoryListResponse represents a list of view history records.
type ViewHistoryListResponse struct {
	Items      []ViewHistoryResponse `json:"items"`
	TotalCount int                   `json:"total_count" example:"100"`
	Limit      int                   `json:"limit" example:"50"`
	Offset     int                   `json:"offset" example:"0"`
}

// ErrorResponse represents an error response.
type ErrorResponse struct {
	Code    string `json:"code" example:"INVALID_REQUEST"`
	Message string `json:"message" example:"Invalid request format"`
	Details string `json:"details,omitempty"`
}
