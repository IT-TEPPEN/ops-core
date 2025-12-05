package handlers

import (
	"net/http"
	"strconv"

	"opscore/backend/internal/document/application/usecase"
	intererror "opscore/backend/internal/document/interfaces/error"
	"opscore/backend/internal/document/interfaces/api/schema"

	"github.com/gin-gonic/gin"
)

// DocumentHandler holds dependencies for document handlers.
type DocumentHandler struct {
	docUseCase usecase.DocumentUseCase
	logger     Logger
}

// Logger interface defines the logging methods required by the handler.
type Logger interface {
	Info(msg string, args ...any)
	Error(msg string, args ...any)
	Debug(msg string, args ...any)
	Warn(msg string, args ...any)
}

// NewDocumentHandler creates a new DocumentHandler.
func NewDocumentHandler(uc usecase.DocumentUseCase, logger Logger) *DocumentHandler {
	return &DocumentHandler{
		docUseCase: uc,
		logger:     logger,
	}
}

// CreateDocument godoc
// @Summary Create a new document
// @Description Create a new document with an initial version
// @Tags documents
// @Accept json
// @Produce json
// @Param document body schema.CreateDocumentRequest true "Document information"
// @Success 201 {object} schema.DocumentResponse "Document created successfully"
// @Failure 400 {object} schema.ErrorResponse "Invalid request body"
// @Failure 500 {object} schema.ErrorResponse "Internal server error"
// @Router /documents [post]
func (h *DocumentHandler) CreateDocument(c *gin.Context) {
	var req schema.CreateDocumentRequest
	requestID := c.GetString("request_id")

	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Error("Invalid request body", "request_id", requestID, "error", err.Error())
		c.JSON(http.StatusBadRequest, schema.ErrorResponse{Code: "INVALID_REQUEST", Message: "Invalid request body: " + err.Error()})
		return
	}

	// Convert schema to DTO
	dtoReq := schema.ToCreateDocumentDTO(req)

	h.logger.Info("Creating document", "request_id", requestID, "title", dtoReq.Title)
	result, err := h.docUseCase.CreateDocument(c.Request.Context(), &dtoReq)

	if err != nil {
		httpErr := intererror.MapToHTTPError(err, requestID)
		h.logger.Error("Failed to create document", "request_id", requestID, "error", err.Error(), "http_code", httpErr.Code)
		c.JSON(httpErr.StatusCode, schema.ErrorResponse{Code: httpErr.Code, Message: httpErr.Message, Details: httpErr.Details})
		return
	}

	h.logger.Info("Document created successfully", "request_id", requestID, "doc_id", result.ID)
	response := schema.FromDocumentDTO(*result)
	c.JSON(http.StatusCreated, response)
}

// GetDocument godoc
// @Summary Get document details
// @Description Retrieves detailed information about a specific document by ID
// @Tags documents
// @Produce json
// @Param docId path string true "Document ID" example:"a1b2c3d4-e5f6-7890-1234-567890abcdef"
// @Success 200 {object} schema.DocumentResponse "Successfully retrieved document details"
// @Failure 400 {object} schema.ErrorResponse "Invalid document ID format"
// @Failure 404 {object} schema.ErrorResponse "Document not found"
// @Failure 500 {object} schema.ErrorResponse "Internal server error"
// @Router /documents/{docId} [get]
func (h *DocumentHandler) GetDocument(c *gin.Context) {
	docID := c.Param("docId")
	requestID := c.GetString("request_id")

	if docID == "" {
		h.logger.Warn("Missing document ID", "request_id", requestID)
		c.JSON(http.StatusBadRequest, schema.ErrorResponse{Code: "INVALID_ID", Message: "Document ID is required"})
		return
	}

	h.logger.Info("Getting document details", "request_id", requestID, "doc_id", docID)
	result, err := h.docUseCase.GetDocument(c.Request.Context(), docID)

	if err != nil {
		httpErr := intererror.MapToHTTPError(err, requestID)
		h.logger.Error("Failed to get document details", "request_id", requestID, "doc_id", docID, "error", err.Error(), "http_code", httpErr.Code)
		c.JSON(httpErr.StatusCode, schema.ErrorResponse{Code: httpErr.Code, Message: httpErr.Message, Details: httpErr.Details})
		return
	}

	h.logger.Info("Successfully retrieved document details", "request_id", requestID, "doc_id", docID)
	response := schema.FromDocumentDTO(*result)
	c.JSON(http.StatusOK, response)
}

// UpdateDocument godoc
// @Summary Update a document
// @Description Update an existing document by creating a new version
// @Tags documents
// @Accept json
// @Produce json
// @Param docId path string true "Document ID" example:"a1b2c3d4-e5f6-7890-1234-567890abcdef"
// @Param document body schema.UpdateDocumentRequest true "Document update information"
// @Success 200 {object} schema.DocumentResponse "Document updated successfully"
// @Failure 400 {object} schema.ErrorResponse "Invalid request body or document ID"
// @Failure 404 {object} schema.ErrorResponse "Document not found"
// @Failure 500 {object} schema.ErrorResponse "Internal server error"
// @Router /documents/{docId} [put]
func (h *DocumentHandler) UpdateDocument(c *gin.Context) {
	docID := c.Param("docId")
	requestID := c.GetString("request_id")
	var req schema.UpdateDocumentRequest

	if docID == "" {
		h.logger.Warn("Missing document ID", "request_id", requestID)
		c.JSON(http.StatusBadRequest, schema.ErrorResponse{Code: "INVALID_ID", Message: "Document ID is required"})
		return
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Error("Invalid request body", "request_id", requestID, "doc_id", docID, "error", err.Error())
		c.JSON(http.StatusBadRequest, schema.ErrorResponse{Code: "INVALID_REQUEST", Message: "Invalid request body: " + err.Error()})
		return
	}

	// Convert schema to DTO
	dtoReq := schema.ToUpdateDocumentDTO(req)

	h.logger.Info("Updating document", "request_id", requestID, "doc_id", docID)
	result, err := h.docUseCase.UpdateDocument(c.Request.Context(), docID, &dtoReq)

	if err != nil {
		httpErr := intererror.MapToHTTPError(err, requestID)
		h.logger.Error("Failed to update document", "request_id", requestID, "doc_id", docID, "error", err.Error(), "http_code", httpErr.Code)
		c.JSON(httpErr.StatusCode, schema.ErrorResponse{Code: httpErr.Code, Message: httpErr.Message, Details: httpErr.Details})
		return
	}

	h.logger.Info("Document updated successfully", "request_id", requestID, "doc_id", docID)
	response := schema.FromDocumentDTO(*result)
	c.JSON(http.StatusOK, response)
}

// ListDocuments godoc
// @Summary List all documents
// @Description Retrieves a list of all published documents
// @Tags documents
// @Produce json
// @Param repository_id query string false "Filter by repository ID"
// @Success 200 {object} schema.ListDocumentsResponse "Successfully retrieved documents"
// @Failure 500 {object} schema.ErrorResponse "Internal server error"
// @Router /documents [get]
func (h *DocumentHandler) ListDocuments(c *gin.Context) {
	requestID := c.GetString("request_id")
	repositoryID := c.Query("repository_id")

	h.logger.Info("Listing documents", "request_id", requestID, "repository_id", repositoryID)

	var docs []schema.DocumentListItemResponse
	var err error

	if repositoryID != "" {
		dtoResp, e := h.docUseCase.ListDocumentsByRepository(c.Request.Context(), repositoryID)
		if e != nil {
			err = e
		} else {
			docs = schema.FromDocumentListDTO(dtoResp)
		}
	} else {
		dtoResp, e := h.docUseCase.ListDocuments(c.Request.Context())
		if e != nil {
			err = e
		} else {
			docs = schema.FromDocumentListDTO(dtoResp)
		}
	}

	if err != nil {
		httpErr := intererror.MapToHTTPError(err, requestID)
		h.logger.Error("Failed to list documents", "request_id", requestID, "error", err.Error(), "http_code", httpErr.Code)
		c.JSON(httpErr.StatusCode, schema.ErrorResponse{Code: httpErr.Code, Message: httpErr.Message, Details: httpErr.Details})
		return
	}

	h.logger.Info("Successfully listed documents", "request_id", requestID, "count", len(docs))
	c.JSON(http.StatusOK, schema.ListDocumentsResponse{Documents: docs})
}

// GetDocumentVersions godoc
// @Summary Get document version history
// @Description Retrieves all versions for a specific document
// @Tags documents
// @Produce json
// @Param docId path string true "Document ID" example:"a1b2c3d4-e5f6-7890-1234-567890abcdef"
// @Success 200 {object} schema.VersionHistoryResponse "Successfully retrieved version history"
// @Failure 400 {object} schema.ErrorResponse "Invalid document ID format"
// @Failure 404 {object} schema.ErrorResponse "Document not found"
// @Failure 500 {object} schema.ErrorResponse "Internal server error"
// @Router /documents/{docId}/versions [get]
func (h *DocumentHandler) GetDocumentVersions(c *gin.Context) {
	docID := c.Param("docId")
	requestID := c.GetString("request_id")

	if docID == "" {
		h.logger.Warn("Missing document ID", "request_id", requestID)
		c.JSON(http.StatusBadRequest, schema.ErrorResponse{Code: "INVALID_ID", Message: "Document ID is required"})
		return
	}

	h.logger.Info("Getting document versions", "request_id", requestID, "doc_id", docID)
	result, err := h.docUseCase.GetDocumentVersions(c.Request.Context(), docID)

	if err != nil {
		httpErr := intererror.MapToHTTPError(err, requestID)
		h.logger.Error("Failed to get document versions", "request_id", requestID, "doc_id", docID, "error", err.Error(), "http_code", httpErr.Code)
		c.JSON(httpErr.StatusCode, schema.ErrorResponse{Code: httpErr.Code, Message: httpErr.Message, Details: httpErr.Details})
		return
	}

	h.logger.Info("Successfully retrieved document versions", "request_id", requestID, "doc_id", docID, "version_count", len(result.Versions))
	response := schema.FromVersionHistoryDTO(*result)
	c.JSON(http.StatusOK, response)
}

// GetDocumentVersion godoc
// @Summary Get a specific document version
// @Description Retrieves details of a specific version of a document
// @Tags documents
// @Produce json
// @Param docId path string true "Document ID" example:"a1b2c3d4-e5f6-7890-1234-567890abcdef"
// @Param version path int true "Version number" example:"1"
// @Success 200 {object} schema.DocumentVersionResponse "Successfully retrieved version details"
// @Failure 400 {object} schema.ErrorResponse "Invalid document ID or version number"
// @Failure 404 {object} schema.ErrorResponse "Document or version not found"
// @Failure 500 {object} schema.ErrorResponse "Internal server error"
// @Router /documents/{docId}/versions/{version} [get]
func (h *DocumentHandler) GetDocumentVersion(c *gin.Context) {
	docID := c.Param("docId")
	versionStr := c.Param("version")
	requestID := c.GetString("request_id")

	if docID == "" {
		h.logger.Warn("Missing document ID", "request_id", requestID)
		c.JSON(http.StatusBadRequest, schema.ErrorResponse{Code: "INVALID_ID", Message: "Document ID is required"})
		return
	}

	versionNumber, err := strconv.Atoi(versionStr)
	if err != nil || versionNumber < 1 {
		h.logger.Warn("Invalid version number", "request_id", requestID, "version", versionStr)
		c.JSON(http.StatusBadRequest, schema.ErrorResponse{Code: "INVALID_VERSION", Message: "Version number must be a positive integer"})
		return
	}

	h.logger.Info("Getting document version", "request_id", requestID, "doc_id", docID, "version", versionNumber)
	result, err := h.docUseCase.GetDocumentVersion(c.Request.Context(), docID, versionNumber)

	if err != nil {
		httpErr := intererror.MapToHTTPError(err, requestID)
		h.logger.Error("Failed to get document version", "request_id", requestID, "doc_id", docID, "version", versionNumber, "error", err.Error(), "http_code", httpErr.Code)
		c.JSON(httpErr.StatusCode, schema.ErrorResponse{Code: httpErr.Code, Message: httpErr.Message, Details: httpErr.Details})
		return
	}

	h.logger.Info("Successfully retrieved document version", "request_id", requestID, "doc_id", docID, "version", versionNumber)
	response := schema.FromDocumentVersionDTO(*result)
	c.JSON(http.StatusOK, response)
}

// PublishDocumentVersion godoc
// @Summary Publish a specific document version
// @Description Publishes a specific version of a document, making it the current version
// @Tags documents
// @Accept json
// @Produce json
// @Param docId path string true "Document ID" example:"a1b2c3d4-e5f6-7890-1234-567890abcdef"
// @Param version path int true "Version number" example:"1"
// @Success 200 {object} schema.DocumentResponse "Version published successfully"
// @Failure 400 {object} schema.ErrorResponse "Invalid document ID or version number"
// @Failure 404 {object} schema.ErrorResponse "Document or version not found"
// @Failure 500 {object} schema.ErrorResponse "Internal server error"
// @Router /documents/{docId}/versions/{version}/publish [post]
func (h *DocumentHandler) PublishDocumentVersion(c *gin.Context) {
	docID := c.Param("docId")
	versionStr := c.Param("version")
	requestID := c.GetString("request_id")

	if docID == "" {
		h.logger.Warn("Missing document ID", "request_id", requestID)
		c.JSON(http.StatusBadRequest, schema.ErrorResponse{Code: "INVALID_ID", Message: "Document ID is required"})
		return
	}

	versionNumber, err := strconv.Atoi(versionStr)
	if err != nil || versionNumber < 1 {
		h.logger.Warn("Invalid version number", "request_id", requestID, "version", versionStr)
		c.JSON(http.StatusBadRequest, schema.ErrorResponse{Code: "INVALID_VERSION", Message: "Version number must be a positive integer"})
		return
	}

	h.logger.Info("Publishing document version", "request_id", requestID, "doc_id", docID, "version", versionNumber)
	result, err := h.docUseCase.PublishDocumentVersion(c.Request.Context(), docID, versionNumber)

	if err != nil {
		httpErr := intererror.MapToHTTPError(err, requestID)
		h.logger.Error("Failed to publish document version", "request_id", requestID, "doc_id", docID, "version", versionNumber, "error", err.Error(), "http_code", httpErr.Code)
		c.JSON(httpErr.StatusCode, schema.ErrorResponse{Code: httpErr.Code, Message: httpErr.Message, Details: httpErr.Details})
		return
	}

	h.logger.Info("Document version published successfully", "request_id", requestID, "doc_id", docID, "version", versionNumber)
	response := schema.FromDocumentDTO(*result)
	c.JSON(http.StatusOK, response)
}

// RollbackDocumentVersion godoc
// @Summary Rollback to a previous document version
// @Description Rolls back the document to a previous version
// @Tags documents
// @Accept json
// @Produce json
// @Param docId path string true "Document ID" example:"a1b2c3d4-e5f6-7890-1234-567890abcdef"
// @Param version path int true "Version number" example:"1"
// @Success 200 {object} schema.DocumentResponse "Rollback successful"
// @Failure 400 {object} schema.ErrorResponse "Invalid document ID or version number"
// @Failure 404 {object} schema.ErrorResponse "Document or version not found"
// @Failure 500 {object} schema.ErrorResponse "Internal server error"
// @Router /documents/{docId}/versions/{version}/rollback [post]
func (h *DocumentHandler) RollbackDocumentVersion(c *gin.Context) {
	docID := c.Param("docId")
	versionStr := c.Param("version")
	requestID := c.GetString("request_id")

	if docID == "" {
		h.logger.Warn("Missing document ID", "request_id", requestID)
		c.JSON(http.StatusBadRequest, schema.ErrorResponse{Code: "INVALID_ID", Message: "Document ID is required"})
		return
	}

	versionNumber, err := strconv.Atoi(versionStr)
	if err != nil || versionNumber < 1 {
		h.logger.Warn("Invalid version number", "request_id", requestID, "version", versionStr)
		c.JSON(http.StatusBadRequest, schema.ErrorResponse{Code: "INVALID_VERSION", Message: "Version number must be a positive integer"})
		return
	}

	h.logger.Info("Rolling back document version", "request_id", requestID, "doc_id", docID, "version", versionNumber)
	result, err := h.docUseCase.RollbackDocumentVersion(c.Request.Context(), docID, versionNumber)

	if err != nil {
		httpErr := intererror.MapToHTTPError(err, requestID)
		h.logger.Error("Failed to rollback document version", "request_id", requestID, "doc_id", docID, "version", versionNumber, "error", err.Error(), "http_code", httpErr.Code)
		c.JSON(httpErr.StatusCode, schema.ErrorResponse{Code: httpErr.Code, Message: httpErr.Message, Details: httpErr.Details})
		return
	}

	h.logger.Info("Document version rolled back successfully", "request_id", requestID, "doc_id", docID, "version", versionNumber)
	response := schema.FromDocumentDTO(*result)
	c.JSON(http.StatusOK, response)
}

// UpdateDocumentMetadata godoc
// @Summary Update document metadata
// @Description Updates document metadata such as owner, access scope, and auto-update setting
// @Tags documents
// @Accept json
// @Produce json
// @Param docId path string true "Document ID" example:"a1b2c3d4-e5f6-7890-1234-567890abcdef"
// @Param metadata body schema.UpdateDocumentMetadataRequest true "Metadata update information"
// @Success 200 {object} schema.DocumentResponse "Metadata updated successfully"
// @Failure 400 {object} schema.ErrorResponse "Invalid request body or document ID"
// @Failure 404 {object} schema.ErrorResponse "Document not found"
// @Failure 500 {object} schema.ErrorResponse "Internal server error"
// @Router /documents/{docId}/metadata [patch]
func (h *DocumentHandler) UpdateDocumentMetadata(c *gin.Context) {
	docID := c.Param("docId")
	requestID := c.GetString("request_id")
	var req schema.UpdateDocumentMetadataRequest

	if docID == "" {
		h.logger.Warn("Missing document ID", "request_id", requestID)
		c.JSON(http.StatusBadRequest, schema.ErrorResponse{Code: "INVALID_ID", Message: "Document ID is required"})
		return
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Error("Invalid request body", "request_id", requestID, "doc_id", docID, "error", err.Error())
		c.JSON(http.StatusBadRequest, schema.ErrorResponse{Code: "INVALID_REQUEST", Message: "Invalid request body: " + err.Error()})
		return
	}

	// Convert schema to DTO
	dtoReq := schema.ToUpdateDocumentMetadataDTO(req)

	h.logger.Info("Updating document metadata", "request_id", requestID, "doc_id", docID)
	result, err := h.docUseCase.UpdateDocumentMetadata(c.Request.Context(), docID, &dtoReq)

	if err != nil {
		httpErr := intererror.MapToHTTPError(err, requestID)
		h.logger.Error("Failed to update document metadata", "request_id", requestID, "doc_id", docID, "error", err.Error(), "http_code", httpErr.Code)
		c.JSON(httpErr.StatusCode, schema.ErrorResponse{Code: httpErr.Code, Message: httpErr.Message, Details: httpErr.Details})
		return
	}

	h.logger.Info("Document metadata updated successfully", "request_id", requestID, "doc_id", docID)
	response := schema.FromDocumentDTO(*result)
	c.JSON(http.StatusOK, response)
}
