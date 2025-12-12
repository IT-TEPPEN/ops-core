package handlers

import (
	"context"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"

	"opscore/backend/internal/execution_record/application/dto"
	"opscore/backend/internal/execution_record/interfaces/api/schema"
)

// AttachmentUsecase defines the interface for attachment business logic.
type AttachmentUsecase interface {
	UploadAttachment(ctx context.Context, req *dto.UploadAttachmentRequest) (*dto.AttachmentResponse, error)
	GetAttachment(ctx context.Context, attachmentID string) (*dto.AttachmentResponse, error)
	GetAttachmentFile(ctx context.Context, attachmentID string) (io.ReadCloser, *dto.AttachmentResponse, error)
	ListAttachmentsByRecordID(ctx context.Context, recordID string) ([]*dto.AttachmentResponse, error)
	ListAttachmentsByStepID(ctx context.Context, stepID string) ([]*dto.AttachmentResponse, error)
	DeleteAttachment(ctx context.Context, attachmentID string) error
	GetAttachmentURL(ctx context.Context, attachmentID string, expirationMinutes int) (string, error)
}

// AttachmentHandler handles HTTP requests for attachments.
type AttachmentHandler struct {
	usecase AttachmentUsecase
}

// NewAttachmentHandler creates a new AttachmentHandler.
func NewAttachmentHandler(uc AttachmentUsecase) *AttachmentHandler {
	return &AttachmentHandler{usecase: uc}
}

// UploadAttachment godoc
// @Summary Upload an attachment
// @Description Upload a file attachment to an execution record step
// @Tags attachments
// @Accept multipart/form-data
// @Produce json
// @Param id path string true "Execution Record ID" example:"a1b2c3d4-e5f6-7890-1234-567890abcdef"
// @Param execution_step_id formData string true "Execution Step ID"
// @Param file formData file true "File to upload"
// @Success 201 {object} schema.AttachmentResponse "Attachment uploaded successfully"
// @Failure 400 {object} map[string]string "Invalid request"
// @Failure 401 {object} map[string]string "User not authenticated"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /execution-records/{id}/attachments [post]
func (h *AttachmentHandler) UploadAttachment(c *gin.Context) {
	recordID := c.Param("id")
	if recordID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Record ID is required"})
		return
	}

	var req schema.UploadAttachmentRequest
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request: " + err.Error()})
		return
	}

	// Open the uploaded file
	file, err := req.File.Open()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to open file: " + err.Error()})
		return
	}
	defer file.Close()

	// Get uploader ID from context (should be set by auth middleware)
	uploaderID := c.GetString("user_id")
	if uploaderID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	// Detect MIME type from file header content type
	mimeType := req.File.Header.Get("Content-Type")
	if mimeType == "" {
		mimeType = "application/octet-stream"
	}

	dtoReq := &dto.UploadAttachmentRequest{
		ExecutionRecordID: recordID,
		ExecutionStepID:   req.ExecutionStepID,
		FileName:          req.File.Filename,
		FileSize:          req.File.Size,
		MimeType:          mimeType,
		UploadedBy:        uploaderID,
		File:              file,
	}

	resp, err := h.usecase.UploadAttachment(c.Request.Context(), dtoReq)
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusCreated, schema.FromAttachmentDTO(resp))
}

// GetAttachment godoc
// @Summary Get attachment metadata
// @Description Retrieves metadata information about a specific attachment
// @Tags attachments
// @Produce json
// @Param id path string true "Attachment ID" example:"a1b2c3d4-e5f6-7890-1234-567890abcdef"
// @Success 200 {object} schema.AttachmentResponse "Successfully retrieved attachment metadata"
// @Failure 400 {object} map[string]string "Invalid attachment ID"
// @Failure 404 {object} map[string]string "Attachment not found"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /attachments/{id} [get]
func (h *AttachmentHandler) GetAttachment(c *gin.Context) {
	attachmentID := c.Param("id")
	if attachmentID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Attachment ID is required"})
		return
	}

	resp, err := h.usecase.GetAttachment(c.Request.Context(), attachmentID)
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, schema.FromAttachmentDTO(resp))
}

// DownloadAttachment godoc
// @Summary Download an attachment file
// @Description Download the actual file content of an attachment
// @Tags attachments
// @Produce application/octet-stream
// @Param id path string true "Attachment ID" example:"a1b2c3d4-e5f6-7890-1234-567890abcdef"
// @Success 200 {file} binary "Attachment file content"
// @Failure 400 {object} map[string]string "Invalid attachment ID"
// @Failure 404 {object} map[string]string "Attachment not found"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /attachments/{id}/download [get]
func (h *AttachmentHandler) DownloadAttachment(c *gin.Context) {
	attachmentID := c.Param("id")
	if attachmentID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Attachment ID is required"})
		return
	}

	file, attachment, err := h.usecase.GetAttachmentFile(c.Request.Context(), attachmentID)
	if err != nil {
		handleError(c, err)
		return
	}
	defer file.Close()

	c.Header("Content-Disposition", "attachment; filename="+attachment.FileName)
	c.Header("Content-Type", attachment.MimeType)
	c.DataFromReader(http.StatusOK, attachment.FileSize, attachment.MimeType, file, nil)
}

// ListAttachments godoc
// @Summary List attachments for an execution record
// @Description Retrieves all attachments associated with an execution record
// @Tags attachments
// @Produce json
// @Param id path string true "Execution Record ID" example:"a1b2c3d4-e5f6-7890-1234-567890abcdef"
// @Success 200 {object} schema.ListAttachmentsResponse "List of attachments"
// @Failure 400 {object} map[string]string "Invalid record ID"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /execution-records/{id}/attachments [get]
func (h *AttachmentHandler) ListAttachments(c *gin.Context) {
	recordID := c.Param("id")
	if recordID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Record ID is required"})
		return
	}

	attachments, err := h.usecase.ListAttachmentsByRecordID(c.Request.Context(), recordID)
	if err != nil {
		handleError(c, err)
		return
	}

	responses := make([]schema.AttachmentResponse, len(attachments))
	for i, attachment := range attachments {
		responses[i] = schema.FromAttachmentDTO(attachment)
	}

	c.JSON(http.StatusOK, schema.ListAttachmentsResponse{Attachments: responses})
}

// ListStepAttachments godoc
// @Summary List attachments for an execution step
// @Description Retrieves all attachments associated with a specific execution step
// @Tags attachments
// @Produce json
// @Param id path string true "Execution Record ID" example:"a1b2c3d4-e5f6-7890-1234-567890abcdef"
// @Param stepId path string true "Execution Step ID" example:"s1b2c3d4-e5f6-7890-1234-567890abcdef"
// @Success 200 {object} schema.ListAttachmentsResponse "List of attachments"
// @Failure 400 {object} map[string]string "Invalid step ID"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /execution-records/{id}/steps/{stepId}/attachments [get]
func (h *AttachmentHandler) ListStepAttachments(c *gin.Context) {
	stepID := c.Param("stepId")
	if stepID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Step ID is required"})
		return
	}

	attachments, err := h.usecase.ListAttachmentsByStepID(c.Request.Context(), stepID)
	if err != nil {
		handleError(c, err)
		return
	}

	responses := make([]schema.AttachmentResponse, len(attachments))
	for i, attachment := range attachments {
		responses[i] = schema.FromAttachmentDTO(attachment)
	}

	c.JSON(http.StatusOK, schema.ListAttachmentsResponse{Attachments: responses})
}

// DeleteAttachment godoc
// @Summary Delete an attachment
// @Description Delete a specific attachment by ID
// @Tags attachments
// @Param id path string true "Attachment ID" example:"a1b2c3d4-e5f6-7890-1234-567890abcdef"
// @Success 204 "Attachment deleted successfully"
// @Failure 400 {object} map[string]string "Invalid attachment ID"
// @Failure 404 {object} map[string]string "Attachment not found"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /attachments/{id} [delete]
func (h *AttachmentHandler) DeleteAttachment(c *gin.Context) {
	attachmentID := c.Param("id")
	if attachmentID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Attachment ID is required"})
		return
	}

	err := h.usecase.DeleteAttachment(c.Request.Context(), attachmentID)
	if err != nil {
		handleError(c, err)
		return
	}

	c.Status(http.StatusNoContent)
}

// GetAttachmentURL godoc
// @Summary Get signed URL for attachment
// @Description Get a presigned URL for downloading an attachment (for S3 storage) or download endpoint (for local storage)
// @Tags attachments
// @Produce json
// @Param id path string true "Attachment ID" example:"a1b2c3d4-e5f6-7890-1234-567890abcdef"
// @Success 200 {object} map[string]string "Presigned URL or download endpoint"
// @Failure 400 {object} map[string]string "Invalid attachment ID"
// @Failure 404 {object} map[string]string "Attachment not found"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /attachments/{id}/url [get]
func (h *AttachmentHandler) GetAttachmentURL(c *gin.Context) {
	attachmentID := c.Param("id")
	if attachmentID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Attachment ID is required"})
		return
	}

	// Default expiration is 60 minutes
	expirationMinutes := 60

	url, err := h.usecase.GetAttachmentURL(c.Request.Context(), attachmentID, expirationMinutes)
	if err != nil {
		handleError(c, err)
		return
	}

	// If URL is empty (local storage), return the download endpoint URL
	if url == "" {
		url = "/api/v1/attachments/" + attachmentID + "/download"
	}

	c.JSON(http.StatusOK, gin.H{"url": url})
}
