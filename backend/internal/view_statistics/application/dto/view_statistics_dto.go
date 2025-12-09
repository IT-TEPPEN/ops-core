package dto

import "time"

// DocumentStatisticsResponse represents document statistics.
type DocumentStatisticsResponse struct {
	DocumentID          string    `json:"document_id"`
	TotalViews          int64     `json:"total_views"`
	UniqueViewers       int64     `json:"unique_viewers"`
	LastViewedAt        time.Time `json:"last_viewed_at"`
	AverageViewDuration int       `json:"average_view_duration"`
}

// UserStatisticsResponse represents user statistics.
type UserStatisticsResponse struct {
	UserID              string `json:"user_id"`
	TotalViews          int64  `json:"total_views"`
	UniqueDocuments     int64  `json:"unique_documents"`
}

// PopularDocumentResponse represents a popular document.
type PopularDocumentResponse struct {
	DocumentID    string    `json:"document_id"`
	TotalViews    int64     `json:"total_views"`
	UniqueViewers int64     `json:"unique_viewers"`
	LastViewedAt  time.Time `json:"last_viewed_at"`
}

// PopularDocumentsResponse represents a list of popular documents.
type PopularDocumentsResponse struct {
	Items []PopularDocumentResponse `json:"items"`
	Limit int                       `json:"limit"`
}

// RecentDocumentResponse represents a recently viewed document.
type RecentDocumentResponse struct {
	DocumentID   string    `json:"document_id"`
	LastViewedAt time.Time `json:"last_viewed_at"`
	TotalViews   int64     `json:"total_views"`
}

// RecentDocumentsResponse represents a list of recently viewed documents.
type RecentDocumentsResponse struct {
	Items []RecentDocumentResponse `json:"items"`
	Limit int                      `json:"limit"`
}
