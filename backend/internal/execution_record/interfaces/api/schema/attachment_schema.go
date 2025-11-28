package schema

import (
	"mime/multipart"
	"time"
)

// UploadAttachmentRequest represents the API request to upload an attachment.
type UploadAttachmentRequest struct {
	ExecutionStepID string                `form:"execution_step_id" binding:"required"`
	File            *multipart.FileHeader `form:"file" binding:"required"`
}

// AttachmentResponse represents the API response for an attachment.
type AttachmentResponse struct {
	ID                string    `json:"id"`
	ExecutionRecordID string    `json:"execution_record_id"`
	ExecutionStepID   string    `json:"execution_step_id"`
	FileName          string    `json:"file_name"`
	FileSize          int64     `json:"file_size"`
	MimeType          string    `json:"mime_type"`
	StorageType       string    `json:"storage_type"`
	UploadedBy        string    `json:"uploaded_by"`
	UploadedAt        time.Time `json:"uploaded_at"`
}

// ListAttachmentsResponse represents the API response for listing attachments.
type ListAttachmentsResponse struct {
	Attachments []AttachmentResponse `json:"attachments"`
}
