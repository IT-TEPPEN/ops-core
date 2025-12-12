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

// CreateExecutionRecord godoc
// @Summary Create a new execution record
// @Description Create a new execution record (work record) for tracking procedure execution
// @Tags execution-records
// @Accept json
// @Produce json
// @Param execution-record body schema.CreateExecutionRecordRequest true "Execution record information"
// @Success 201 {object} schema.ExecutionRecordResponse "Execution record created successfully"
// @Failure 400 {object} map[string]string "Invalid request body"
// @Failure 401 {object} map[string]string "User not authenticated"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /execution-records [post]
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

// GetExecutionRecord godoc
// @Summary Get execution record details
// @Description Retrieves detailed information about a specific execution record by ID
// @Tags execution-records
// @Produce json
// @Param id path string true "Execution Record ID" example:"a1b2c3d4-e5f6-7890-1234-567890abcdef"
// @Success 200 {object} schema.ExecutionRecordResponse "Successfully retrieved execution record details"
// @Failure 400 {object} map[string]string "Invalid execution record ID"
// @Failure 404 {object} map[string]string "Execution record not found"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /execution-records/{id} [get]
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

// AddStep godoc
// @Summary Add a step to execution record
// @Description Add a new step to an existing execution record
// @Tags execution-records
// @Accept json
// @Produce json
// @Param id path string true "Execution Record ID" example:"a1b2c3d4-e5f6-7890-1234-567890abcdef"
// @Param step body schema.AddStepRequest true "Step information"
// @Success 200 {object} schema.ExecutionRecordResponse "Step added successfully"
// @Failure 400 {object} map[string]string "Invalid request body or record ID"
// @Failure 404 {object} map[string]string "Execution record not found"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /execution-records/{id}/steps [post]
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

// UpdateStepNotes godoc
// @Summary Update step notes
// @Description Update the notes of a specific step in an execution record
// @Tags execution-records
// @Accept json
// @Produce json
// @Param id path string true "Execution Record ID" example:"a1b2c3d4-e5f6-7890-1234-567890abcdef"
// @Param stepNumber path int true "Step Number" example:1
// @Param notes body schema.UpdateStepNotesRequest true "Step notes"
// @Success 200 {object} schema.ExecutionRecordResponse "Step notes updated successfully"
// @Failure 400 {object} map[string]string "Invalid request body or step number"
// @Failure 404 {object} map[string]string "Execution record or step not found"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /execution-records/{id}/steps/{stepNumber}/notes [put]
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

// UpdateNotes godoc
// @Summary Update execution record notes
// @Description Update the overall notes of an execution record
// @Tags execution-records
// @Accept json
// @Produce json
// @Param id path string true "Execution Record ID" example:"a1b2c3d4-e5f6-7890-1234-567890abcdef"
// @Param notes body schema.UpdateNotesRequest true "Execution record notes"
// @Success 200 {object} schema.ExecutionRecordResponse "Notes updated successfully"
// @Failure 400 {object} map[string]string "Invalid request body or record ID"
// @Failure 404 {object} map[string]string "Execution record not found"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /execution-records/{id}/notes [put]
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

// UpdateTitle godoc
// @Summary Update execution record title
// @Description Update the title of an execution record
// @Tags execution-records
// @Accept json
// @Produce json
// @Param id path string true "Execution Record ID" example:"a1b2c3d4-e5f6-7890-1234-567890abcdef"
// @Param title body schema.UpdateTitleRequest true "Execution record title"
// @Success 200 {object} schema.ExecutionRecordResponse "Title updated successfully"
// @Failure 400 {object} map[string]string "Invalid request body or record ID"
// @Failure 404 {object} map[string]string "Execution record not found"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /execution-records/{id}/title [put]
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

// Complete godoc
// @Summary Complete an execution record
// @Description Mark an execution record as completed
// @Tags execution-records
// @Produce json
// @Param id path string true "Execution Record ID" example:"a1b2c3d4-e5f6-7890-1234-567890abcdef"
// @Success 200 {object} schema.ExecutionRecordResponse "Execution record completed successfully"
// @Failure 400 {object} map[string]string "Invalid record ID"
// @Failure 404 {object} map[string]string "Execution record not found"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /execution-records/{id}/complete [post]
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

// MarkAsFailed godoc
// @Summary Mark execution record as failed
// @Description Mark an execution record as failed
// @Tags execution-records
// @Produce json
// @Param id path string true "Execution Record ID" example:"a1b2c3d4-e5f6-7890-1234-567890abcdef"
// @Success 200 {object} schema.ExecutionRecordResponse "Execution record marked as failed successfully"
// @Failure 400 {object} map[string]string "Invalid record ID"
// @Failure 404 {object} map[string]string "Execution record not found"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /execution-records/{id}/fail [post]
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

// UpdateAccessScope godoc
// @Summary Update execution record access scope
// @Description Update the access scope of an execution record (public or private)
// @Tags execution-records
// @Accept json
// @Produce json
// @Param id path string true "Execution Record ID" example:"a1b2c3d4-e5f6-7890-1234-567890abcdef"
// @Param access_scope body schema.UpdateAccessScopeRequest true "Access scope"
// @Success 200 {object} schema.ExecutionRecordResponse "Access scope updated successfully"
// @Failure 400 {object} map[string]string "Invalid request body or record ID"
// @Failure 404 {object} map[string]string "Execution record not found"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /execution-records/{id}/access-scope [put]
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

// SearchExecutionRecords godoc
// @Summary Search execution records
// @Description Search and filter execution records by various criteria
// @Tags execution-records
// @Produce json
// @Param executor_id query string false "Executor User ID" example:"user-123"
// @Param document_id query string false "Document ID" example:"a1b2c3d4-e5f6-7890-1234-567890abcdef"
// @Param status query string false "Status" example:"in_progress"
// @Param started_from query string false "Started From" example:"2025-01-01T00:00:00Z"
// @Param started_to query string false "Started To" example:"2025-12-31T23:59:59Z"
// @Success 200 {object} map[string][]schema.ExecutionRecordResponse "List of execution records"
// @Failure 400 {object} map[string]string "Invalid query parameters"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /execution-records [get]
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

// DeleteExecutionRecord godoc
// @Summary Delete an execution record
// @Description Delete a specific execution record by ID
// @Tags execution-records
// @Param id path string true "Execution Record ID" example:"a1b2c3d4-e5f6-7890-1234-567890abcdef"
// @Success 204 "Execution record deleted successfully"
// @Failure 400 {object} map[string]string "Invalid record ID"
// @Failure 404 {object} map[string]string "Execution record not found"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /execution-records/{id} [delete]
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
