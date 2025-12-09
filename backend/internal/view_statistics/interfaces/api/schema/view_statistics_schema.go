package schema

import "time"

// DocumentStatisticsResponse represents document statistics.
type DocumentStatisticsResponse struct {
	DocumentID          string    `json:"document_id" example:"a1b2c3d4-e5f6-7890-1234-567890abcdef"`
	TotalViews          int64     `json:"total_views" example:"1234"`
	UniqueViewers       int64     `json:"unique_viewers" example:"567"`
	LastViewedAt        time.Time `json:"last_viewed_at" example:"2024-01-15T10:30:00Z"`
	AverageViewDuration int       `json:"average_view_duration" example:"120"`
}

// UserStatisticsResponse represents user statistics.
type UserStatisticsResponse struct {
	UserID          string `json:"user_id" example:"user-123"`
	TotalViews      int64  `json:"total_views" example:"345"`
	UniqueDocuments int64  `json:"unique_documents" example:"78"`
}

// PopularDocumentResponse represents a popular document.
type PopularDocumentResponse struct {
	DocumentID    string    `json:"document_id" example:"a1b2c3d4-e5f6-7890-1234-567890abcdef"`
	TotalViews    int64     `json:"total_views" example:"1234"`
	UniqueViewers int64     `json:"unique_viewers" example:"567"`
	LastViewedAt  time.Time `json:"last_viewed_at" example:"2024-01-15T10:30:00Z"`
}

// PopularDocumentsResponse represents a list of popular documents.
type PopularDocumentsResponse struct {
	Items []PopularDocumentResponse `json:"items"`
	Limit int                       `json:"limit" example:"10"`
}

// RecentDocumentResponse represents a recently viewed document.
type RecentDocumentResponse struct {
	DocumentID   string    `json:"document_id" example:"a1b2c3d4-e5f6-7890-1234-567890abcdef"`
	LastViewedAt time.Time `json:"last_viewed_at" example:"2024-01-15T10:30:00Z"`
	TotalViews   int64     `json:"total_views" example:"234"`
}

// RecentDocumentsResponse represents a list of recently viewed documents.
type RecentDocumentsResponse struct {
	Items []RecentDocumentResponse `json:"items"`
	Limit int                      `json:"limit" example:"10"`
}

// ErrorResponse represents an error response.
type ErrorResponse struct {
	Code    string `json:"code" example:"INVALID_REQUEST"`
	Message string `json:"message" example:"Invalid request format"`
	Details string `json:"details,omitempty"`
}
