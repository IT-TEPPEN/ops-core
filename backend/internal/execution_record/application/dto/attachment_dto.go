package dto

import (
	"io"
	"time"
)

// UploadAttachmentRequest represents the request to upload an attachment.
type UploadAttachmentRequest struct {
	ExecutionRecordID string
	ExecutionStepID   string
	FileName          string
	FileSize          int64
	MimeType          string
	UploadedBy        string
	File              io.Reader
}

// AttachmentResponse represents an attachment response.
type AttachmentResponse struct {
	ID                string
	ExecutionRecordID string
	ExecutionStepID   string
	FileName          string
	FileSize          int64
	MimeType          string
	StorageType       string
	StoragePath       string
	UploadedBy        string
	UploadedAt        time.Time
}

// GetAttachmentRequest represents the request to get an attachment.
type GetAttachmentRequest struct {
	AttachmentID string
}

// DeleteAttachmentRequest represents the request to delete an attachment.
type DeleteAttachmentRequest struct {
	AttachmentID string
}

// ListAttachmentsRequest represents the request to list attachments.
type ListAttachmentsRequest struct {
	ExecutionRecordID *string
	ExecutionStepID   *string
}
