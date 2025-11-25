package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"opscore/backend/internal/execution_record/application/dto"
	"opscore/backend/internal/execution_record/application/usecase"
	"opscore/backend/internal/execution_record/interfaces/api/schema"
)

// AttachmentHandler handles HTTP requests for attachments.
type AttachmentHandler struct {
	usecase *usecase.AttachmentUsecase
}

// NewAttachmentHandler creates a new AttachmentHandler.
func NewAttachmentHandler(uc *usecase.AttachmentUsecase) *AttachmentHandler {
	return &AttachmentHandler{usecase: uc}
}

// UploadAttachment handles POST /execution-records/:id/attachments
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

	// Get uploader ID from context (would normally come from auth middleware)
	uploaderID := c.GetString("user_id")
	if uploaderID == "" {
		uploaderID = "anonymous" // Fallback for testing
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

// GetAttachment handles GET /attachments/:id
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

// DownloadAttachment handles GET /attachments/:id/download
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

// ListAttachments handles GET /execution-records/:id/attachments
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

// ListStepAttachments handles GET /execution-records/:id/steps/:stepId/attachments
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

// DeleteAttachment handles DELETE /attachments/:id
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
