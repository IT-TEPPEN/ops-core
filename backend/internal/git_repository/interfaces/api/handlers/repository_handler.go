package handlers

import (
	"net/http"
	"opscore/backend/internal/git_repository/application/dto"
	"opscore/backend/internal/git_repository/application/usecase"
	intererror "opscore/backend/internal/git_repository/interfaces/error"
	"opscore/backend/internal/git_repository/interfaces/api/schema"

	"github.com/gin-gonic/gin"
)

// RepositoryHandler holds dependencies for repository handlers.
type RepositoryHandler struct {
	repoUseCase repository.RepositoryUseCase
	logger      Logger // ADR 0008に従ってロガーを追加
}

// Logger インターフェースは構造化ロギングの最小インターフェースを定義
type Logger interface {
	Info(msg string, args ...any)
	Error(msg string, args ...any)
	Debug(msg string, args ...any)
	Warn(msg string, args ...any)
}

// NewRepositoryHandler creates a new RepositoryHandler.
func NewRepositoryHandler(uc repository.RepositoryUseCase, logger Logger) *RepositoryHandler {
	return &RepositoryHandler{
		repoUseCase: uc,
		logger:      logger,
	}
}

// --- Handler Methods ---

// RegisterRepository godoc
// @Summary Register a new repository
// @Description Add a new repository to be managed by OpsCore by providing its Git URL and optional access token.
// @Tags repositories
// @Accept  json
// @Produce  json
// @Param   repository body schema.RegisterRepositoryRequest true "Repository information"
// @Success 201 {object} schema.RepositoryResponse "Repository registered successfully"
// @Failure 400 {object} schema.ErrorResponse "Invalid request body or URL format"
// @Failure 409 {object} schema.ErrorResponse "Repository with this URL already exists"
// @Failure 500 {object} schema.ErrorResponse "Internal server error"
// @Router /repositories [post]
func (h *RepositoryHandler) RegisterRepository(c *gin.Context) {
	var req schema.RegisterRepositoryRequest
	requestID := c.GetString("request_id") // ミドルウェアから設定されたリクエストID

	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Error("Invalid request body", "request_id", requestID, "error", err.Error())
		c.JSON(http.StatusBadRequest, schema.ErrorResponse{Code: "INVALID_REQUEST", Message: "Invalid request body: " + err.Error()})
		return
	}

	// Convert schema to DTO
	dtoReq := schema.ToRegisterRepositoryDTO(req)

	h.logger.Info("Registering repository", "request_id", requestID, "url", dtoReq.URL)
	newRepo, err := h.repoUseCase.Register(c.Request.Context(), dtoReq.URL, dtoReq.AccessToken)

	if err != nil {
		// Use error mapper to convert application errors to HTTP errors
		httpErr := intererror.MapToHTTPError(err, requestID)
		h.logger.Error("Failed to register repository", "request_id", requestID, "error", err.Error(), "http_code", httpErr.Code)
		c.JSON(httpErr.StatusCode, schema.ErrorResponse{Code: httpErr.Code, Message: httpErr.Message})
		return
	}

	h.logger.Info("Repository registered successfully", "request_id", requestID, "repo_id", newRepo.ID())
	// Convert domain entity to DTO, then DTO to schema
	dtoResp := dto.ToRepositoryResponse(newRepo)
	response := schema.FromRepositoryDTO(dtoResp)
	c.JSON(http.StatusCreated, response)
}

// UpdateAccessToken godoc
// @Summary Update repository access token
// @Description Updates the access token used for accessing a private repository
// @Tags repositories
// @Accept  json
// @Produce  json
// @Param   repoId path string true "Repository ID" example:"a1b2c3d4-e5f6-7890-1234-567890abcdef"
// @Param   tokenInfo body schema.UpdateAccessTokenRequest true "Access token information"
// @Success 200 {object} map[string]string "Access token updated successfully"
// @Failure 400 {object} schema.ErrorResponse "Invalid request body or repository ID"
// @Failure 404 {object} schema.ErrorResponse "Repository not found"
// @Failure 500 {object} schema.ErrorResponse "Internal server error"
// @Router /repositories/{repoId}/token [put]
func (h *RepositoryHandler) UpdateAccessToken(c *gin.Context) {
	repoId := c.Param("repoId")
	requestID := c.GetString("request_id")
	var req schema.UpdateAccessTokenRequest

	if repoId == "" {
		h.logger.Warn("Missing repository ID", "request_id", requestID)
		c.JSON(http.StatusBadRequest, schema.ErrorResponse{Code: "INVALID_ID", Message: "Repository ID is required"})
		return
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Error("Invalid request body", "request_id", requestID, "repo_id", repoId, "error", err.Error())
		c.JSON(http.StatusBadRequest, schema.ErrorResponse{Code: "INVALID_REQUEST", Message: "Invalid request body: " + err.Error()})
		return
	}

	// Convert schema to DTO
	dtoReq := schema.ToUpdateAccessTokenDTO(req)

	h.logger.Info("Updating repository access token", "request_id", requestID, "repo_id", repoId)
	err := h.repoUseCase.UpdateAccessToken(c.Request.Context(), repoId, dtoReq.AccessToken)

	if err != nil {
		// Use error mapper to convert application errors to HTTP errors
		httpErr := intererror.MapToHTTPError(err, requestID)
		h.logger.Error("Failed to update access token", "request_id", requestID, "repo_id", repoId, "error", err.Error(), "http_code", httpErr.Code)
		c.JSON(httpErr.StatusCode, schema.ErrorResponse{Code: httpErr.Code, Message: httpErr.Message})
		return
	}

	h.logger.Info("Repository access token updated successfully", "request_id", requestID, "repo_id", repoId)
	c.JSON(http.StatusOK, map[string]string{
		"message": "Access token updated successfully",
		"repoId":  repoId,
	})
}

// ListRepositoryFiles godoc
// @Summary List files in a repository
// @Description Retrieves a list of files and directories within a specified repository.
// @Tags repositories
// @Produce  json
// @Param   repoId path string true "Repository ID" example:"a1b2c3d4-e5f6-7890-1234-567890abcdef"
// @Success 200 {object} schema.ListFilesResponse "Successfully retrieved file list"
// @Failure 400 {object} schema.ErrorResponse "Invalid repository ID format or access token missing"
// @Failure 404 {object} schema.ErrorResponse "Repository not found"
// @Failure 500 {object} schema.ErrorResponse "Internal server error"
// @Router /repositories/{repoId}/files [get]
func (h *RepositoryHandler) ListRepositoryFiles(c *gin.Context) {
	repoId := c.Param("repoId")
	requestID := c.GetString("request_id")

	if repoId == "" {
		h.logger.Warn("Missing repository ID", "request_id", requestID)
		c.JSON(http.StatusBadRequest, schema.ErrorResponse{Code: "INVALID_ID", Message: "Repository ID is required"})
		return
	}

	h.logger.Info("Listing repository files", "request_id", requestID, "repo_id", repoId)
	// Call the use case which now returns []entity.FileNode
	domainFiles, err := h.repoUseCase.ListFiles(c.Request.Context(), repoId)

	if err != nil {
		// Use error mapper to convert application errors to HTTP errors
		httpErr := intererror.MapToHTTPError(err, requestID)
		h.logger.Error("Failed to list repository files", "request_id", requestID, "repo_id", repoId, "error", err.Error(), "http_code", httpErr.Code)
		c.JSON(httpErr.StatusCode, schema.ErrorResponse{Code: httpErr.Code, Message: httpErr.Message})
		return
	}

	// Map domain entity.FileNode to DTO FileNode, then to schema
	dtoFiles := dto.ToFileNodeList(domainFiles)
	responseFiles := schema.FromFileNodeListDTO(dtoFiles)

	h.logger.Info("Successfully listed repository files", "request_id", requestID, "repo_id", repoId, "file_count", len(responseFiles))
	c.JSON(http.StatusOK, schema.ListFilesResponse{Files: responseFiles})
}

// SelectRepositoryFiles godoc
// @Summary Select manageable files in a repository
// @Description Marks specific files within a repository as manageable by OpsCore.
// @Tags repositories
// @Accept  json
// @Produce  json
// @Param   repoId path string true "Repository ID" example:"a1b2c3d4-e5f6-7890-1234-567890abcdef"
// @Param   files body schema.SelectFilesRequest true "List of file paths to select"
// @Success 200 {object} schema.SelectFilesResponse "Files selected successfully"
// @Failure 400 {object} schema.ErrorResponse "Invalid request body or repository ID"
// @Failure 404 {object} schema.ErrorResponse "Repository not found"
// @Failure 500 {object} schema.ErrorResponse "Internal server error"
// @Router /repositories/{repoId}/files/select [post]
func (h *RepositoryHandler) SelectRepositoryFiles(c *gin.Context) {
	repoId := c.Param("repoId")
	requestID := c.GetString("request_id") // ミドルウェアから設定されたリクエストID
	var req schema.SelectFilesRequest

	if repoId == "" {
		h.logger.Warn("Missing repository ID", "request_id", requestID)
		c.JSON(http.StatusBadRequest, schema.ErrorResponse{Code: "INVALID_ID", Message: "Repository ID is required"})
		return
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Error("Invalid request body", "request_id", requestID, "repo_id", repoId, "error", err.Error())
		c.JSON(http.StatusBadRequest, schema.ErrorResponse{Code: "INVALID_REQUEST", Message: "Invalid request body: " + err.Error()})
		return
	}
	if len(req.FilePaths) == 0 {
		h.logger.Warn("Empty file paths", "request_id", requestID, "repo_id", repoId)
		c.JSON(http.StatusBadRequest, schema.ErrorResponse{Code: "INVALID_REQUEST", Message: "filePaths cannot be empty"})
		return
	}

	// Convert schema to DTO
	dtoReq := schema.ToSelectFilesDTO(req)

	h.logger.Info("Selecting repository files", "request_id", requestID, "repo_id", repoId, "file_count", len(dtoReq.FilePaths))
	err := h.repoUseCase.SelectFiles(c.Request.Context(), repoId, dtoReq.FilePaths)

	if err != nil {
		// Use error mapper to convert application errors to HTTP errors
		httpErr := intererror.MapToHTTPError(err, requestID)
		h.logger.Error("Failed to select files", "request_id", requestID, "repo_id", repoId, "error", err.Error(), "http_code", httpErr.Code)
		c.JSON(httpErr.StatusCode, schema.ErrorResponse{Code: httpErr.Code, Message: httpErr.Message})
		return
	}

	h.logger.Info("Files selected successfully", "request_id", requestID, "repo_id", repoId, "file_count", len(dtoReq.FilePaths))
	c.JSON(http.StatusOK, schema.SelectFilesResponse{
		Message:       "Files selected successfully",
		RepoID:        repoId,
		SelectedFiles: len(dtoReq.FilePaths),
	})
}

// GetSelectedMarkdown godoc
// @Summary Get selected Markdown content from a repository
// @Description Retrieves the concatenated content of all selected Markdown files for a given repository.
// @Tags repositories
// @Produce  json
// @Param   repoId path string true "Repository ID" example:"a1b2c3d4-e5f6-7890-1234-567890abcdef"
// @Success 200 {object} schema.GetMarkdownResponse "Successfully retrieved Markdown content"
// @Failure 400 {object} schema.ErrorResponse "Invalid repository ID format"
// @Failure 404 {object} schema.ErrorResponse "Repository not found or no files selected"
// @Failure 500 {object} schema.ErrorResponse "Internal server error"
// @Router /repositories/{repoId}/markdown [get]
func (h *RepositoryHandler) GetSelectedMarkdown(c *gin.Context) {
	repoId := c.Param("repoId")
	requestID := c.GetString("request_id") // ミドルウェアから設定されたリクエストID

	if repoId == "" {
		h.logger.Warn("Missing repository ID", "request_id", requestID)
		c.JSON(http.StatusBadRequest, schema.ErrorResponse{Code: "INVALID_ID", Message: "Repository ID is required"})
		return
	}

	h.logger.Info("Getting selected Markdown content", "request_id", requestID, "repo_id", repoId)
	markdownContent, err := h.repoUseCase.GetSelectedMarkdown(c.Request.Context(), repoId)

	if err != nil {
		// Use error mapper to convert application errors to HTTP errors
		httpErr := intererror.MapToHTTPError(err, requestID)
		h.logger.Error("Failed to retrieve Markdown content", "request_id", requestID, "repo_id", repoId, "error", err.Error(), "http_code", httpErr.Code)
		c.JSON(httpErr.StatusCode, schema.ErrorResponse{Code: httpErr.Code, Message: httpErr.Message})
		return
	}

	h.logger.Info("Successfully retrieved Markdown content", "request_id", requestID, "repo_id", repoId, "content_length", len(markdownContent))
	c.JSON(http.StatusOK, schema.GetMarkdownResponse{
		RepoID:  repoId,
		Content: markdownContent,
	})
}

// ListRepositories godoc
// @Summary List all repositories
// @Description Retrieves a list of all repositories registered in OpsCore
// @Tags repositories
// @Produce json
// @Success 200 {object} schema.ListRepositoriesResponse "Successfully retrieved repositories"
// @Failure 500 {object} schema.ErrorResponse "Internal server error"
// @Router /repositories [get]
func (h *RepositoryHandler) ListRepositories(c *gin.Context) {
	requestID := c.GetString("request_id") // ミドルウェアから設定されたリクエストID

	h.logger.Info("Listing all repositories", "request_id", requestID)

	repos, err := h.repoUseCase.ListRepositories(c.Request.Context())
	if err != nil {
		// Use error mapper to convert application errors to HTTP errors
		httpErr := intererror.MapToHTTPError(err, requestID)
		h.logger.Error("Failed to list repositories", "request_id", requestID, "error", err.Error(), "http_code", httpErr.Code)
		c.JSON(httpErr.StatusCode, schema.ErrorResponse{Code: httpErr.Code, Message: httpErr.Message})
		return
	}

	// Map domain models to DTOs, then to schema
	dtoResponses := dto.ToRepositoryResponseList(repos)
	schemaResponses := schema.FromRepositoryListDTO(dtoResponses)

	h.logger.Info("Successfully listed repositories", "request_id", requestID, "repo_count", len(repos))
	c.JSON(http.StatusOK, schema.ListRepositoriesResponse{
		Repositories: schemaResponses,
	})
}

// GetRepository godoc
// @Summary Get repository details
// @Description Retrieves detailed information about a specific repository by ID
// @Tags repositories
// @Produce json
// @Param   repoId path string true "Repository ID" example:"a1b2c3d4-e5f6-7890-1234-567890abcdef"
// @Success 200 {object} schema.RepositoryResponse "Successfully retrieved repository details"
// @Failure 400 {object} schema.ErrorResponse "Invalid repository ID format"
// @Failure 404 {object} schema.ErrorResponse "Repository not found"
// @Failure 500 {object} schema.ErrorResponse "Internal server error"
// @Router /repositories/{repoId} [get]
func (h *RepositoryHandler) GetRepository(c *gin.Context) {
	repoId := c.Param("repoId")
	requestID := c.GetString("request_id")

	if repoId == "" {
		h.logger.Warn("Missing repository ID", "request_id", requestID)
		c.JSON(http.StatusBadRequest, schema.ErrorResponse{Code: "INVALID_ID", Message: "Repository ID is required"})
		return
	}

	h.logger.Info("Getting repository details", "request_id", requestID, "repo_id", repoId)
	repo, err := h.repoUseCase.GetRepository(c.Request.Context(), repoId)

	if err != nil {
		// Use error mapper to convert application errors to HTTP errors
		httpErr := intererror.MapToHTTPError(err, requestID)
		h.logger.Error("Failed to get repository details", "request_id", requestID, "repo_id", repoId, "error", err.Error(), "http_code", httpErr.Code)
		c.JSON(httpErr.StatusCode, schema.ErrorResponse{Code: httpErr.Code, Message: httpErr.Message})
		return
	}

	h.logger.Info("Successfully retrieved repository details", "request_id", requestID, "repo_id", repoId)
	dtoResp := dto.ToRepositoryResponse(repo)
	response := schema.FromRepositoryDTO(dtoResp)
	c.JSON(http.StatusOK, response)
}
