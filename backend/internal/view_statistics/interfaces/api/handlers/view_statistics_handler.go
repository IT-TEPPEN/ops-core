package handlers

import (
	"net/http"
	"strconv"

	"opscore/backend/internal/view_statistics/application/usecase"
	"opscore/backend/internal/view_statistics/interfaces/api/schema"

	"github.com/gin-gonic/gin"
)

// ViewStatisticsHandler holds dependencies for view statistics handlers.
type ViewStatisticsHandler struct {
	useCase usecase.ViewStatisticsUseCase
	logger  Logger
}

// Logger interface defines the logging methods required by the handler.
type Logger interface {
	Info(msg string, args ...any)
	Error(msg string, args ...any)
	Debug(msg string, args ...any)
	Warn(msg string, args ...any)
}

// NewViewStatisticsHandler creates a new ViewStatisticsHandler.
func NewViewStatisticsHandler(uc usecase.ViewStatisticsUseCase, logger Logger) *ViewStatisticsHandler {
	return &ViewStatisticsHandler{
		useCase: uc,
		logger:  logger,
	}
}

// GetDocumentStatistics godoc
// @Summary Get document statistics
// @Description Retrieves view statistics for a specific document
// @Tags statistics
// @Produce json
// @Param id path string true "Document ID" example:"a1b2c3d4-e5f6-7890-1234-567890abcdef"
// @Success 200 {object} schema.DocumentStatisticsResponse "Statistics retrieved successfully"
// @Failure 400 {object} schema.ErrorResponse "Invalid request"
// @Failure 500 {object} schema.ErrorResponse "Internal server error"
// @Router /api/documents/{id}/statistics [get]
func (h *ViewStatisticsHandler) GetDocumentStatistics(c *gin.Context) {
	documentID := c.Param("id")
	requestID := c.GetString("request_id")

	h.logger.Info("Getting document statistics", "request_id", requestID, "document_id", documentID)

	result, err := h.useCase.GetDocumentStatistics(c.Request.Context(), documentID)
	if err != nil {
		h.logger.Error("Failed to get document statistics", "request_id", requestID, "error", err.Error())
		c.JSON(http.StatusInternalServerError, schema.ErrorResponse{Code: "INTERNAL_ERROR", Message: "Failed to retrieve statistics"})
		return
	}

	h.logger.Info("Document statistics retrieved successfully", "request_id", requestID)
	c.JSON(http.StatusOK, schema.DocumentStatisticsResponse{
		DocumentID:          result.DocumentID,
		TotalViews:          result.TotalViews,
		UniqueViewers:       result.UniqueViewers,
		LastViewedAt:        result.LastViewedAt,
		AverageViewDuration: result.AverageViewDuration,
	})
}

// GetUserStatistics godoc
// @Summary Get user statistics
// @Description Retrieves view statistics for a specific user
// @Tags statistics
// @Produce json
// @Param id path string true "User ID" example:"user-123"
// @Success 200 {object} schema.UserStatisticsResponse "Statistics retrieved successfully"
// @Failure 400 {object} schema.ErrorResponse "Invalid request"
// @Failure 500 {object} schema.ErrorResponse "Internal server error"
// @Router /api/users/{id}/statistics [get]
func (h *ViewStatisticsHandler) GetUserStatistics(c *gin.Context) {
	userID := c.Param("id")
	requestID := c.GetString("request_id")

	h.logger.Info("Getting user statistics", "request_id", requestID, "user_id", userID)

	result, err := h.useCase.GetUserStatistics(c.Request.Context(), userID)
	if err != nil {
		h.logger.Error("Failed to get user statistics", "request_id", requestID, "error", err.Error())
		c.JSON(http.StatusInternalServerError, schema.ErrorResponse{Code: "INTERNAL_ERROR", Message: "Failed to retrieve statistics"})
		return
	}

	h.logger.Info("User statistics retrieved successfully", "request_id", requestID)
	c.JSON(http.StatusOK, schema.UserStatisticsResponse{
		UserID:          result.UserID,
		TotalViews:      result.TotalViews,
		UniqueDocuments: result.UniqueDocuments,
	})
}

// GetPopularDocuments godoc
// @Summary Get popular documents
// @Description Retrieves the most popular documents based on view count
// @Tags statistics
// @Produce json
// @Param limit query int false "Limit" default(10) minimum(1) maximum(100)
// @Param days query int false "Days to look back" default(30) minimum(1)
// @Success 200 {object} schema.PopularDocumentsResponse "Popular documents retrieved successfully"
// @Failure 500 {object} schema.ErrorResponse "Internal server error"
// @Router /api/statistics/popular-documents [get]
func (h *ViewStatisticsHandler) GetPopularDocuments(c *gin.Context) {
	requestID := c.GetString("request_id")

	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	days, _ := strconv.Atoi(c.DefaultQuery("days", "30"))

	h.logger.Info("Getting popular documents", "request_id", requestID, "limit", limit, "days", days)

	result, err := h.useCase.GetPopularDocuments(c.Request.Context(), limit, days)
	if err != nil {
		h.logger.Error("Failed to get popular documents", "request_id", requestID, "error", err.Error())
		c.JSON(http.StatusInternalServerError, schema.ErrorResponse{Code: "INTERNAL_ERROR", Message: "Failed to retrieve popular documents"})
		return
	}

	items := make([]schema.PopularDocumentResponse, len(result.Items))
	for i, item := range result.Items {
		items[i] = schema.PopularDocumentResponse{
			DocumentID:    item.DocumentID,
			TotalViews:    item.TotalViews,
			UniqueViewers: item.UniqueViewers,
			LastViewedAt:  item.LastViewedAt,
		}
	}

	h.logger.Info("Popular documents retrieved successfully", "request_id", requestID, "count", len(items))
	c.JSON(http.StatusOK, schema.PopularDocumentsResponse{
		Items: items,
		Limit: result.Limit,
	})
}

// GetRecentlyViewedDocuments godoc
// @Summary Get recently viewed documents
// @Description Retrieves the most recently viewed documents
// @Tags statistics
// @Produce json
// @Param limit query int false "Limit" default(10) minimum(1) maximum(100)
// @Success 200 {object} schema.RecentDocumentsResponse "Recent documents retrieved successfully"
// @Failure 500 {object} schema.ErrorResponse "Internal server error"
// @Router /api/statistics/recent-documents [get]
func (h *ViewStatisticsHandler) GetRecentlyViewedDocuments(c *gin.Context) {
	requestID := c.GetString("request_id")

	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	h.logger.Info("Getting recently viewed documents", "request_id", requestID, "limit", limit)

	result, err := h.useCase.GetRecentlyViewedDocuments(c.Request.Context(), limit)
	if err != nil {
		h.logger.Error("Failed to get recently viewed documents", "request_id", requestID, "error", err.Error())
		c.JSON(http.StatusInternalServerError, schema.ErrorResponse{Code: "INTERNAL_ERROR", Message: "Failed to retrieve recently viewed documents"})
		return
	}

	items := make([]schema.RecentDocumentResponse, len(result.Items))
	for i, item := range result.Items {
		items[i] = schema.RecentDocumentResponse{
			DocumentID:   item.DocumentID,
			LastViewedAt: item.LastViewedAt,
			TotalViews:   item.TotalViews,
		}
	}

	h.logger.Info("Recently viewed documents retrieved successfully", "request_id", requestID, "count", len(items))
	c.JSON(http.StatusOK, schema.RecentDocumentsResponse{
		Items: items,
		Limit: result.Limit,
	})
}
