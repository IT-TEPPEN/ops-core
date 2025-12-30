package handlers

import (
	"net/http"
	"opscore/backend/internal/git_repository/application/dto"
	repository "opscore/backend/internal/git_repository/application/usecase"
	"opscore/backend/internal/git_repository/interfaces/api/schema"
	intererror "opscore/backend/internal/git_repository/interfaces/error"

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
// @Security BearerAuth
// @Param   repoId path string true "Repository ID" example:"a1b2c3d4-e5f6-7890-1234-567890abcdef"
// @Success 200 {object} schema.ListFilesResponse "Successfully retrieved file list"
// @Failure 400 {object} schema.ErrorResponse "Invalid repository ID format or access token missing"
// @Failure 401 {object} schema.ErrorResponse "Authentication required"
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

	// Get user ID from context (requires authentication middleware)
	userID, exists := c.Get("user_id")
	if !exists {
		h.logger.Error("User ID not found in context", "request_id", requestID)
		c.JSON(http.StatusUnauthorized, schema.ErrorResponse{Code: "UNAUTHORIZED", Message: "Authentication required"})
		return
	}

	h.logger.Info("Listing repository files", "request_id", requestID, "repo_id", repoId, "user_id", userID)
	// Call the use case which now returns []entity.FileNode
	domainFiles, err := h.repoUseCase.ListFiles(c.Request.Context(), repoId, userID.(string))

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

// GetFileContents godoc
// @Summary Get file contents from a repository
// @Description Retrieves the content of a specific file from a repository by its file path provided as a query parameter.
// @Tags repositories
// @Produce  json
// @Security BearerAuth
// @Param   repoId path string true "Repository ID" example:"a1b2c3d4-e5f6-7890-1234-567890abcdef"
// @Param   path query string true "File path in the repository" example:"README.md" example:"docs/adr/0001.md"
// @Success 200 {object} schema.GetFileContentsResponse "Successfully retrieved file contents"
// @Failure 400 {object} schema.ErrorResponse "Invalid repository ID or file path"
// @Failure 401 {object} schema.ErrorResponse "Authentication required"
// @Failure 404 {object} schema.ErrorResponse "Repository or file not found"
// @Failure 500 {object} schema.ErrorResponse "Internal server error"
// @Router /repositories/{repoId}/files/content [get]
func (h *RepositoryHandler) GetFileContents(c *gin.Context) {
	repoId := c.Param("repoId")
	filePath := c.Query("path")
	requestID := c.GetString("request_id")

	if repoId == "" {
		h.logger.Warn("Missing repository ID", "request_id", requestID)
		c.JSON(http.StatusBadRequest, schema.ErrorResponse{Code: "INVALID_ID", Message: "Repository ID is required"})
		return
	}

	if filePath == "" {
		h.logger.Warn("Missing file path", "request_id", requestID, "repo_id", repoId)
		c.JSON(http.StatusBadRequest, schema.ErrorResponse{Code: "INVALID_REQUEST", Message: "File path is required"})
		return
	}

	// Get user ID from context (requires authentication middleware)
	userID, exists := c.Get("user_id")
	if !exists {
		h.logger.Error("User ID not found in context", "request_id", requestID)
		c.JSON(http.StatusUnauthorized, schema.ErrorResponse{Code: "UNAUTHORIZED", Message: "Authentication required"})
		return
	}

	h.logger.Info("Getting file contents", "request_id", requestID, "repo_id", repoId, "file_path", filePath, "user_id", userID)
	content, err := h.repoUseCase.GetFileContents(c.Request.Context(), repoId, filePath, userID.(string))

	if err != nil {
		// Use error mapper to convert application errors to HTTP errors
		httpErr := intererror.MapToHTTPError(err, requestID)
		h.logger.Error("Failed to retrieve file contents", "request_id", requestID, "repo_id", repoId, "file_path", filePath, "error", err.Error(), "http_code", httpErr.Code)
		c.JSON(httpErr.StatusCode, schema.ErrorResponse{Code: httpErr.Code, Message: httpErr.Message})
		return
	}

	h.logger.Info("Successfully retrieved file contents", "request_id", requestID, "repo_id", repoId, "file_path", filePath, "content_length", len(content))
	c.JSON(http.StatusOK, schema.GetFileContentsResponse{
		RepoID:   repoId,
		FilePath: filePath,
		Content:  content,
	})
}

// GetRepositoryContents godoc
// @Summary Get repository contents at a specific path
// @Description Retrieves files and directories at the specified path (non-recursive). Returns only direct children of the specified directory.
// @Tags repositories
// @Produce  json
// @Security BearerAuth
// @Param   repoId path string true "Repository ID" example:"a1b2c3d4-e5f6-7890-1234-567890abcdef"
// @Param   path query string false "Directory path relative to repository root (empty = root)" example:"docs/adr"
// @Success 200 {object} schema.ListFilesResponse "Successfully retrieved directory contents"
// @Failure 400 {object} schema.ErrorResponse "Invalid repository ID"
// @Failure 401 {object} schema.ErrorResponse "Authentication required"
// @Failure 404 {object} schema.ErrorResponse "Repository or path not found"
// @Failure 500 {object} schema.ErrorResponse "Internal server error"
// @Router /repositories/{repoId}/contents [get]
func (h *RepositoryHandler) GetRepositoryContents(c *gin.Context) {
	repoId := c.Param("repoId")
	path := c.DefaultQuery("path", "") // Default to root if not provided
	requestID := c.GetString("request_id")

	if repoId == "" {
		h.logger.Warn("Missing repository ID", "request_id", requestID)
		c.JSON(http.StatusBadRequest, schema.ErrorResponse{Code: "INVALID_ID", Message: "Repository ID is required"})
		return
	}

	// Get user ID from context (requires authentication middleware)
	userID, exists := c.Get("user_id")
	if !exists {
		h.logger.Error("User ID not found in context", "request_id", requestID)
		c.JSON(http.StatusUnauthorized, schema.ErrorResponse{Code: "UNAUTHORIZED", Message: "Authentication required"})
		return
	}

	h.logger.Info("Getting repository contents", "request_id", requestID, "repo_id", repoId, "path", path, "user_id", userID)
	// Call the use case to get directory contents
	domainFiles, err := h.repoUseCase.GetDirectoryContents(c.Request.Context(), repoId, path, userID.(string))

	if err != nil {
		// Use error mapper to convert application errors to HTTP errors
		httpErr := intererror.MapToHTTPError(err, requestID)
		h.logger.Error("Failed to get repository contents", "request_id", requestID, "repo_id", repoId, "path", path, "error", err.Error(), "http_code", httpErr.Code)
		c.JSON(httpErr.StatusCode, schema.ErrorResponse{Code: httpErr.Code, Message: httpErr.Message})
		return
	}

	// Map domain entity.FileNode to DTO FileNode, then to schema
	dtoFiles := dto.ToFileNodeList(domainFiles)
	responseFiles := schema.FromFileNodeListDTO(dtoFiles)

	h.logger.Info("Successfully retrieved repository contents", "request_id", requestID, "repo_id", repoId, "path", path, "item_count", len(responseFiles))
	c.JSON(http.StatusOK, schema.ListFilesResponse{Files: responseFiles})
}

// ListRepositories godoc
// @Summary List all repositories
// @Description Retrieves a list of all repositories registered in OpsCore
// @Tags repositories
// @Produce json
// @Security BearerAuth
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
// @Security BearerAuth
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

// GetFileHistory godoc
// @Summary Get file commit history
// @Description Retrieves the commit history for a specific file in a repository
// @Tags repositories
// @Produce  json
// @Security BearerAuth
// @Param   repoId path string true "Repository ID" example:"a1b2c3d4-e5f6-7890-1234-567890abcdef"
// @Param   path query string true "File path relative to repository root" example:"docs/procedure.md"
// @Success 200 {object} schema.GetFileHistoryResponse "File history retrieved successfully"
// @Failure 400 {object} schema.ErrorResponse "Invalid repository ID or file path"
// @Failure 404 {object} schema.ErrorResponse "Repository or file not found"
// @Failure 500 {object} schema.ErrorResponse "Internal server error"
// @Router /repositories/{repoId}/files/history [get]
func (h *RepositoryHandler) GetFileHistory(c *gin.Context) {
	repoId := c.Param("repoId")
	filePath := c.Query("path")
	requestID := c.GetString("request_id")
	userID := c.GetString("user_id") // Get user ID from context (set by auth middleware)

	if repoId == "" {
		h.logger.Warn("Missing repository ID", "request_id", requestID)
		c.JSON(http.StatusBadRequest, schema.ErrorResponse{Code: "INVALID_ID", Message: "Repository ID is required"})
		return
	}

	if filePath == "" {
		h.logger.Warn("Missing file path", "request_id", requestID)
		c.JSON(http.StatusBadRequest, schema.ErrorResponse{Code: "INVALID_PATH", Message: "File path is required"})
		return
	}

	h.logger.Info("Getting file commit history", "request_id", requestID, "repo_id", repoId, "file_path", filePath)
	commits, err := h.repoUseCase.GetFileCommitHistory(c.Request.Context(), repoId, filePath, userID)

	if err != nil {
		httpErr := intererror.MapToHTTPError(err, requestID)
		h.logger.Error("Failed to get file commit history", "request_id", requestID, "repo_id", repoId, "file_path", filePath, "error", err.Error(), "http_code", httpErr.Code)
		c.JSON(httpErr.StatusCode, schema.ErrorResponse{Code: httpErr.Code, Message: httpErr.Message})
		return
	}

	h.logger.Info("Successfully retrieved file commit history", "request_id", requestID, "repo_id", repoId, "file_path", filePath, "commit_count", len(commits))

	// Convert to schema
	commitSchemas := make([]schema.FileCommitInfo, len(commits))
	for i, commit := range commits {
		commitSchemas[i] = schema.FileCommitInfo{
			CommitHash:  commit.Hash,
			Message:     commit.Message,
			Author:      commit.Author,
			AuthorEmail: commit.AuthorEmail,
			Date:        commit.Date,
		}
	}

	response := schema.GetFileHistoryResponse{
		FilePath: filePath,
		Commits:  commitSchemas,
	}
	c.JSON(http.StatusOK, response)
}
