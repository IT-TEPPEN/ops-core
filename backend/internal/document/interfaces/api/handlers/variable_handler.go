package handlers

import (
	"net/http"

	"opscore/backend/internal/document/application/usecase"
	intererror "opscore/backend/internal/document/interfaces/error"
	"opscore/backend/internal/document/interfaces/api/schema"

	"github.com/gin-gonic/gin"
)

// VariableHandler holds dependencies for variable handlers
type VariableHandler struct {
	varUseCase usecase.VariableUseCase
	logger     Logger
}

// NewVariableHandler creates a new VariableHandler
func NewVariableHandler(uc usecase.VariableUseCase, logger Logger) *VariableHandler {
	return &VariableHandler{
		varUseCase: uc,
		logger:     logger,
	}
}

// GetVariableDefinitions godoc
// @Summary Get variable definitions for a document
// @Description Retrieves all variable definitions from the current version of a document
// @Tags variables
// @Produce json
// @Param docId path string true "Document ID" example:"a1b2c3d4-e5f6-7890-1234-567890abcdef"
// @Success 200 {object} schema.GetVariableDefinitionsResponse "Successfully retrieved variable definitions"
// @Failure 400 {object} schema.ErrorResponse "Invalid document ID format"
// @Failure 404 {object} schema.ErrorResponse "Document not found"
// @Failure 500 {object} schema.ErrorResponse "Internal server error"
// @Router /documents/{docId}/variables [get]
func (h *VariableHandler) GetVariableDefinitions(c *gin.Context) {
	docID := c.Param("docId")
	requestID := c.GetString("request_id")

	if docID == "" {
		h.logger.Warn("Missing document ID", "request_id", requestID)
		c.JSON(http.StatusBadRequest, schema.ErrorResponse{Code: "INVALID_ID", Message: "Document ID is required"})
		return
	}

	h.logger.Info("Getting variable definitions", "request_id", requestID, "doc_id", docID)
	variables, err := h.varUseCase.GetVariableDefinitions(c.Request.Context(), docID)

	if err != nil {
		httpErr := intererror.MapToHTTPError(err, requestID)
		h.logger.Error("Failed to get variable definitions", "request_id", requestID, "doc_id", docID, "error", err.Error(), "http_code", httpErr.Code)
		c.JSON(httpErr.StatusCode, schema.ErrorResponse{Code: httpErr.Code, Message: httpErr.Message, Details: httpErr.Details})
		return
	}

	// Convert to schema DTOs
	schemaDTOs := make([]schema.VariableDefinitionDTO, len(variables))
	for i, v := range variables {
		schemaDTOs[i] = schema.VariableDefinitionDTO{
			Name:         v.Name,
			Label:        v.Label,
			Description:  v.Description,
			Type:         v.Type,
			Required:     v.Required,
			DefaultValue: v.DefaultValue,
		}
	}

	h.logger.Info("Successfully retrieved variable definitions", "request_id", requestID, "doc_id", docID, "count", len(variables))
	c.JSON(http.StatusOK, schema.GetVariableDefinitionsResponse{
		DocumentID: docID,
		Variables:  schemaDTOs,
	})
}

// ValidateVariableValues godoc
// @Summary Validate variable values for a document
// @Description Validates the provided variable values against the document's variable definitions.
// @Description
// @Description **Note on Response Status Codes:**
// @Description This endpoint returns HTTP 200 OK for both successful validations and validation failures.
// @Description The response body's "valid" field indicates whether validation passed or failed.
// @Description This design allows distinguishing between:
// @Description - Validation failures (HTTP 200 with valid:false) - client provided invalid data
// @Description - System errors (HTTP 4xx/5xx) - server-side issues or malformed requests
// @Tags variables
// @Accept json
// @Produce json
// @Param docId path string true "Document ID" example:"a1b2c3d4-e5f6-7890-1234-567890abcdef"
// @Param values body schema.ValidateVariableValuesRequest true "Variable values to validate"
// @Success 200 {object} schema.ValidateVariableValuesResponse "Validation result (valid:true or valid:false with errors)"
// @Failure 400 {object} schema.ErrorResponse "Invalid request body or document ID"
// @Failure 404 {object} schema.ErrorResponse "Document not found"
// @Failure 500 {object} schema.ErrorResponse "Internal server error"
// @Router /documents/{docId}/validate-variables [post]
func (h *VariableHandler) ValidateVariableValues(c *gin.Context) {
	docID := c.Param("docId")
	requestID := c.GetString("request_id")
	var req schema.ValidateVariableValuesRequest

	if docID == "" {
		h.logger.Warn("Missing document ID", "request_id", requestID)
		c.JSON(http.StatusBadRequest, schema.ErrorResponse{Code: "INVALID_ID", Message: "Document ID is required"})
		return
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Error("Invalid request body", "request_id", requestID, "doc_id", docID, "error", err.Error())
		c.JSON(http.StatusBadRequest, schema.ErrorResponse{Code: "INVALID_REQUEST", Message: "Invalid request format"})
		return
	}

	// Convert schema DTOs to usecase values
	values := make([]usecase.VariableValue, len(req.Values))
	for i, v := range req.Values {
		values[i] = usecase.VariableValue{
			Name:  v.Name,
			Value: v.Value,
		}
	}

	h.logger.Info("Validating variable values", "request_id", requestID, "doc_id", docID, "value_count", len(values))
	err := h.varUseCase.ValidateVariableValues(c.Request.Context(), docID, values)

	if err != nil {
		// Check if it's a validation error
		httpErr := intererror.MapToHTTPError(err, requestID)
		
		// For validation errors, return a structured response with valid: false
		if httpErr.StatusCode == http.StatusBadRequest && httpErr.Code == "VALIDATION_FAILED" {
			h.logger.Info("Variable validation failed", "request_id", requestID, "doc_id", docID)
			
			// Extract field errors if available
			validationErrors := []schema.ValidationErrorDTO{}
			if httpErr.Details != nil {
				// Extract validation_errors from Details
				if validationErrsInterface, ok := httpErr.Details["validation_errors"]; ok {
					if fieldErrors, ok := validationErrsInterface.([]map[string]interface{}); ok {
						for _, fieldErr := range fieldErrors {
							// Safely extract field and message with type assertions
							name, nameOk := fieldErr["field"].(string)
							message, messageOk := fieldErr["message"].(string)
							if nameOk && messageOk {
								validationErrors = append(validationErrors, schema.ValidationErrorDTO{
									Name:    name,
									Message: message,
								})
							}
						}
					}
				}
			}
			
			c.JSON(http.StatusOK, schema.ValidateVariableValuesResponse{
				Valid:  false,
				Errors: validationErrors,
			})
			return
		}

		// For other errors, return error response
		h.logger.Error("Failed to validate variable values", "request_id", requestID, "doc_id", docID, "error", err.Error(), "http_code", httpErr.Code)
		c.JSON(httpErr.StatusCode, schema.ErrorResponse{Code: httpErr.Code, Message: httpErr.Message, Details: httpErr.Details})
		return
	}

	h.logger.Info("Variable validation successful", "request_id", requestID, "doc_id", docID)
	c.JSON(http.StatusOK, schema.ValidateVariableValuesResponse{
		Valid:  true,
		Errors: []schema.ValidationErrorDTO{},
	})
}
