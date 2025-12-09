package dto

import "time"

// RecordViewRequest represents a request to record a view.
type RecordViewRequest struct {
	DocumentID string `json:"document_id"`
	UserID     string `json:"user_id"`
}

// ViewHistoryResponse represents a view history record response.
type ViewHistoryResponse struct {
	ID           string    `json:"id"`
	DocumentID   string    `json:"document_id"`
	UserID       string    `json:"user_id"`
	ViewedAt     time.Time `json:"viewed_at"`
	ViewDuration int       `json:"view_duration"`
}

// ViewHistoryListResponse represents a list of view history records.
type ViewHistoryListResponse struct {
	Items      []ViewHistoryResponse `json:"items"`
	TotalCount int                   `json:"total_count"`
	Limit      int                   `json:"limit"`
	Offset     int                   `json:"offset"`
}
