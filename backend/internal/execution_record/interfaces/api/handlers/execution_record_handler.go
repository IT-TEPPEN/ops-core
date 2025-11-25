package handlers

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	apperror "opscore/backend/internal/execution_record/application/error"
	"opscore/backend/internal/execution_record/application/dto"
	"opscore/backend/internal/execution_record/application/usecase"
	"opscore/backend/internal/execution_record/interfaces/api/schema"
)

// ExecutionRecordHandler handles HTTP requests for execution records.
type ExecutionRecordHandler struct {
	usecase *usecase.ExecutionRecordUsecase
}

// NewExecutionRecordHandler creates a new ExecutionRecordHandler.
func NewExecutionRecordHandler(uc *usecase.ExecutionRecordUsecase) *ExecutionRecordHandler {
	return &ExecutionRecordHandler{usecase: uc}
}

// CreateExecutionRecord handles POST /execution-records
func (h *ExecutionRecordHandler) CreateExecutionRecord(c *gin.Context) {
	var req schema.CreateExecutionRecordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body: " + err.Error()})
		return
	}

	// Get executor ID from context (should be set by auth middleware)
	executorID := c.GetString("user_id")
	if executorID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	dtoReq := schema.ToCreateExecutionRecordDTO(req, executorID)
	resp, err := h.usecase.CreateExecutionRecord(c.Request.Context(), dtoReq)
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusCreated, schema.FromExecutionRecordDTO(resp))
}

// GetExecutionRecord handles GET /execution-records/:id
func (h *ExecutionRecordHandler) GetExecutionRecord(c *gin.Context) {
	recordID := c.Param("id")
	if recordID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Record ID is required"})
		return
	}

	resp, err := h.usecase.GetExecutionRecord(c.Request.Context(), recordID)
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, schema.FromExecutionRecordDTO(resp))
}

// AddStep handles POST /execution-records/:id/steps
func (h *ExecutionRecordHandler) AddStep(c *gin.Context) {
	recordID := c.Param("id")
	if recordID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Record ID is required"})
		return
	}

	var req schema.AddStepRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body: " + err.Error()})
		return
	}

	dtoReq := schema.ToAddStepDTO(req, recordID)
	resp, err := h.usecase.AddStep(c.Request.Context(), dtoReq)
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, schema.FromExecutionRecordDTO(resp))
}

// UpdateStepNotes handles PUT /execution-records/:id/steps/:stepNumber/notes
func (h *ExecutionRecordHandler) UpdateStepNotes(c *gin.Context) {
	recordID := c.Param("id")
	stepNumberStr := c.Param("stepNumber")

	stepNumber, err := strconv.Atoi(stepNumberStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid step number"})
		return
	}

	var req schema.UpdateStepNotesRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body: " + err.Error()})
		return
	}

	dtoReq := schema.ToUpdateStepNotesDTO(req, recordID, stepNumber)
	resp, err := h.usecase.UpdateStepNotes(c.Request.Context(), dtoReq)
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, schema.FromExecutionRecordDTO(resp))
}

// UpdateNotes handles PUT /execution-records/:id/notes
func (h *ExecutionRecordHandler) UpdateNotes(c *gin.Context) {
	recordID := c.Param("id")
	if recordID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Record ID is required"})
		return
	}

	var req schema.UpdateNotesRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body: " + err.Error()})
		return
	}

	dtoReq := schema.ToUpdateNotesDTO(req, recordID)
	resp, err := h.usecase.UpdateNotes(c.Request.Context(), dtoReq)
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, schema.FromExecutionRecordDTO(resp))
}

// UpdateTitle handles PUT /execution-records/:id/title
func (h *ExecutionRecordHandler) UpdateTitle(c *gin.Context) {
	recordID := c.Param("id")
	if recordID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Record ID is required"})
		return
	}

	var req schema.UpdateTitleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body: " + err.Error()})
		return
	}

	dtoReq := schema.ToUpdateTitleDTO(req, recordID)
	resp, err := h.usecase.UpdateTitle(c.Request.Context(), dtoReq)
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, schema.FromExecutionRecordDTO(resp))
}

// Complete handles POST /execution-records/:id/complete
func (h *ExecutionRecordHandler) Complete(c *gin.Context) {
	recordID := c.Param("id")
	if recordID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Record ID is required"})
		return
	}

	req := &dto.CompleteExecutionRequest{
		ExecutionRecordID: recordID,
	}
	resp, err := h.usecase.Complete(c.Request.Context(), req)
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, schema.FromExecutionRecordDTO(resp))
}

// MarkAsFailed handles POST /execution-records/:id/fail
func (h *ExecutionRecordHandler) MarkAsFailed(c *gin.Context) {
	recordID := c.Param("id")
	if recordID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Record ID is required"})
		return
	}

	req := &dto.MarkAsFailedRequest{
		ExecutionRecordID: recordID,
	}
	resp, err := h.usecase.MarkAsFailed(c.Request.Context(), req)
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, schema.FromExecutionRecordDTO(resp))
}

// UpdateAccessScope handles PUT /execution-records/:id/access-scope
func (h *ExecutionRecordHandler) UpdateAccessScope(c *gin.Context) {
	recordID := c.Param("id")
	if recordID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Record ID is required"})
		return
	}

	var req schema.UpdateAccessScopeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body: " + err.Error()})
		return
	}

	dtoReq := schema.ToUpdateAccessScopeDTO(req, recordID)
	resp, err := h.usecase.UpdateAccessScope(c.Request.Context(), dtoReq)
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, schema.FromExecutionRecordDTO(resp))
}

// SearchExecutionRecords handles GET /execution-records
func (h *ExecutionRecordHandler) SearchExecutionRecords(c *gin.Context) {
	var req schema.SearchExecutionRecordRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid query parameters: " + err.Error()})
		return
	}

	dtoReq := schema.ToSearchExecutionRecordDTO(req)
	records, err := h.usecase.SearchExecutionRecords(c.Request.Context(), dtoReq)
	if err != nil {
		handleError(c, err)
		return
	}

	responses := make([]schema.ExecutionRecordResponse, len(records))
	for i, record := range records {
		responses[i] = schema.FromExecutionRecordDTO(record)
	}

	c.JSON(http.StatusOK, gin.H{"execution_records": responses})
}

// DeleteExecutionRecord handles DELETE /execution-records/:id
func (h *ExecutionRecordHandler) DeleteExecutionRecord(c *gin.Context) {
	recordID := c.Param("id")
	if recordID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Record ID is required"})
		return
	}

	err := h.usecase.DeleteExecutionRecord(c.Request.Context(), recordID)
	if err != nil {
		handleError(c, err)
		return
	}

	c.Status(http.StatusNoContent)
}

// handleError maps application errors to HTTP responses.
func handleError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, apperror.ErrNotFound):
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
	case errors.Is(err, apperror.ErrValidationFailed):
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	case errors.Is(err, apperror.ErrConflict):
		c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
	case errors.Is(err, apperror.ErrForbidden):
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
	case errors.Is(err, apperror.ErrUnauthorized):
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
	default:
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
	}
}
