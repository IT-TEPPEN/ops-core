package handlers

import (
	"net/http"
	"strconv"

	"opscore/backend/internal/view_history/application/dto"
	"opscore/backend/internal/view_history/application/usecase"
	"opscore/backend/internal/view_history/interfaces/api/schema"

	"github.com/gin-gonic/gin"
)

// ViewHistoryHandler holds dependencies for view history handlers.
type ViewHistoryHandler struct {
	useCase usecase.ViewHistoryUseCase
	logger  Logger
}

// Logger interface defines the logging methods required by the handler.
type Logger interface {
	Info(msg string, args ...any)
	Error(msg string, args ...any)
	Debug(msg string, args ...any)
	Warn(msg string, args ...any)
}

// NewViewHistoryHandler creates a new ViewHistoryHandler.
func NewViewHistoryHandler(uc usecase.ViewHistoryUseCase, logger Logger) *ViewHistoryHandler {
	return &ViewHistoryHandler{
		useCase: uc,
		logger:  logger,
	}
}

// RecordView godoc
// @Summary Record a document view
// @Description Records a view of a document by a user
// @Tags view-history
// @Accept json
// @Produce json
// @Param id path string true "Document ID" example:"a1b2c3d4-e5f6-7890-1234-567890abcdef"
// @Param request body schema.RecordViewRequest true "User ID"
// @Success 201 {object} schema.ViewHistoryResponse "View recorded successfully"
// @Failure 400 {object} schema.ErrorResponse "Invalid request"
// @Failure 500 {object} schema.ErrorResponse "Internal server error"
// @Router /api/documents/{id}/views [post]
func (h *ViewHistoryHandler) RecordView(c *gin.Context) {
	documentID := c.Param("id")
	requestID := c.GetString("request_id")

	var req schema.RecordViewRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Error("Invalid request body", "request_id", requestID, "error", err.Error())
		c.JSON(http.StatusBadRequest, schema.ErrorResponse{Code: "INVALID_REQUEST", Message: "Invalid request format"})
		return
	}

	h.logger.Info("Recording view", "request_id", requestID, "document_id", documentID, "user_id", req.UserID)

	result, err := h.useCase.RecordView(c.Request.Context(), &dto.RecordViewRequest{
		DocumentID: documentID,
		UserID:     req.UserID,
	})

	if err != nil {
		h.logger.Error("Failed to record view", "request_id", requestID, "error", err.Error())
		c.JSON(http.StatusInternalServerError, schema.ErrorResponse{Code: "INTERNAL_ERROR", Message: "Failed to record view"})
		return
	}

	h.logger.Info("View recorded successfully", "request_id", requestID, "view_id", result.ID)
	c.JSON(http.StatusCreated, schema.ViewHistoryResponse{
		ID:           result.ID,
		DocumentID:   result.DocumentID,
		UserID:       result.UserID,
		ViewedAt:     result.ViewedAt,
		ViewDuration: result.ViewDuration,
	})
}

// GetUserViewHistory godoc
// @Summary Get user's view history
// @Description Retrieves the view history for a specific user
// @Tags view-history
// @Produce json
// @Param id path string true "User ID" example:"user-123"
// @Param limit query int false "Limit" default(50) minimum(1) maximum(100)
// @Param offset query int false "Offset" default(0) minimum(0)
// @Success 200 {object} schema.ViewHistoryListResponse "View history retrieved successfully"
// @Failure 400 {object} schema.ErrorResponse "Invalid request"
// @Failure 500 {object} schema.ErrorResponse "Internal server error"
// @Router /api/users/{id}/view-history [get]
func (h *ViewHistoryHandler) GetUserViewHistory(c *gin.Context) {
	userID := c.Param("id")
	requestID := c.GetString("request_id")

	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "50"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	h.logger.Info("Getting user view history", "request_id", requestID, "user_id", userID)

	result, err := h.useCase.GetViewHistory(c.Request.Context(), userID, limit, offset)
	if err != nil {
		h.logger.Error("Failed to get view history", "request_id", requestID, "error", err.Error())
		c.JSON(http.StatusInternalServerError, schema.ErrorResponse{Code: "INTERNAL_ERROR", Message: "Failed to retrieve view history"})
		return
	}

	items := make([]schema.ViewHistoryResponse, len(result.Items))
	for i, item := range result.Items {
		items[i] = schema.ViewHistoryResponse{
			ID:           item.ID,
			DocumentID:   item.DocumentID,
			UserID:       item.UserID,
			ViewedAt:     item.ViewedAt,
			ViewDuration: item.ViewDuration,
		}
	}

	h.logger.Info("View history retrieved successfully", "request_id", requestID, "count", len(items))
	c.JSON(http.StatusOK, schema.ViewHistoryListResponse{
		Items:      items,
		TotalCount: result.TotalCount,
		Limit:      result.Limit,
		Offset:     result.Offset,
	})
}

// GetDocumentViewHistory godoc
// @Summary Get document's view history
// @Description Retrieves the view history for a specific document
// @Tags view-history
// @Produce json
// @Param id path string true "Document ID" example:"a1b2c3d4-e5f6-7890-1234-567890abcdef"
// @Param limit query int false "Limit" default(50) minimum(1) maximum(100)
// @Param offset query int false "Offset" default(0) minimum(0)
// @Success 200 {object} schema.ViewHistoryListResponse "View history retrieved successfully"
// @Failure 400 {object} schema.ErrorResponse "Invalid request"
// @Failure 500 {object} schema.ErrorResponse "Internal server error"
// @Router /api/documents/{id}/view-history [get]
func (h *ViewHistoryHandler) GetDocumentViewHistory(c *gin.Context) {
	documentID := c.Param("id")
	requestID := c.GetString("request_id")

	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "50"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	h.logger.Info("Getting document view history", "request_id", requestID, "document_id", documentID)

	result, err := h.useCase.GetDocumentViewHistory(c.Request.Context(), documentID, limit, offset)
	if err != nil {
		h.logger.Error("Failed to get document view history", "request_id", requestID, "error", err.Error())
		c.JSON(http.StatusInternalServerError, schema.ErrorResponse{Code: "INTERNAL_ERROR", Message: "Failed to retrieve view history"})
		return
	}

	items := make([]schema.ViewHistoryResponse, len(result.Items))
	for i, item := range result.Items {
		items[i] = schema.ViewHistoryResponse{
			ID:           item.ID,
			DocumentID:   item.DocumentID,
			UserID:       item.UserID,
			ViewedAt:     item.ViewedAt,
			ViewDuration: item.ViewDuration,
		}
	}

	h.logger.Info("Document view history retrieved successfully", "request_id", requestID, "count", len(items))
	c.JSON(http.StatusOK, schema.ViewHistoryListResponse{
		Items:      items,
		TotalCount: result.TotalCount,
		Limit:      result.Limit,
		Offset:     result.Offset,
	})
}
