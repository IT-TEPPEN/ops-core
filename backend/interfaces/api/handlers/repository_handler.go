package handlers

import (
	"errors"
	"net/http"
	"opscore/backend/usecases/repository"
	"time"

	"github.com/gin-gonic/gin"
)

// --- Structs (Request/Response/Handler) ---

// RegisterRepositoryRequest represents the request body for registering a repository
type RegisterRepositoryRequest struct {
	URL string `json:"url" binding:"required" example:"https://github.com/user/repo.git"`
}

// RepositoryResponse represents the standard response format for a repository
type RepositoryResponse struct {
	ID        string    `json:"id" example:"a1b2c3d4-e5f6-7890-1234-567890abcdef"`
	Name      string    `json:"name" example:"repo"`
	URL       string    `json:"url" example:"https://github.com/user/repo.git"`
	CreatedAt time.Time `json:"created_at" example:"2025-04-22T10:00:00Z"`
	UpdatedAt time.Time `json:"updated_at" example:"2025-04-22T10:00:00Z"`
}

// ErrorResponse represents a generic error response
type ErrorResponse struct {
	Code    string `json:"code" example:"INVALID_REQUEST"`
	Message string `json:"message" example:"Invalid request body"`
}

// FileNode represents a file or directory within a repository
type FileNode struct {
	Path string `json:"path" example:"src/main.go"`
	Type string `json:"type" example:"file"` // "file" or "dir"
}

// ListFilesResponse represents the response for listing repository files
type ListFilesResponse struct {
	Files []FileNode `json:"files"`
}

// SelectFilesRequest represents the request body for selecting files
type SelectFilesRequest struct {
	FilePaths []string `json:"filePaths" binding:"required,dive,required" example:"[\"README.md\", \"docs/adr/0001.md\"]"` // List of file paths to select
}

// SelectFilesResponse represents the success response for selecting files
type SelectFilesResponse struct {
	Message       string `json:"message" example:"Files selected successfully"`
	RepoID        string `json:"repoId" example:"a1b2c3d4-e5f6-7890-1234-567890abcdef"`
	SelectedFiles int    `json:"selectedFiles" example:"2"`
}

// GetMarkdownResponse represents the response containing selected Markdown content
type GetMarkdownResponse struct {
	RepoID  string `json:"repoId" example:"a1b2c3d4-e5f6-7890-1234-567890abcdef"`
	Content string `json:"content" example:"# Project Title\\n\\n## ADR 1\\n..."` // Concatenated Markdown content
}

// ListRepositoriesResponse represents the response for listing all repositories
type ListRepositoriesResponse struct {
	Repositories []RepositoryResponse `json:"repositories"`
}

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
// @Description Add a new repository to be managed by OpsCore by providing its Git URL.
// @Tags repositories
// @Accept  json
// @Produce  json
// @Param   repository body RegisterRepositoryRequest true "Repository URL"
// @Success 201 {object} RepositoryResponse "Repository registered successfully"
// @Failure 400 {object} ErrorResponse "Invalid request body or URL format"
// @Failure 409 {object} ErrorResponse "Repository with this URL already exists"
// @Failure 500 {object} ErrorResponse "Internal server error"
// @Router /repositories [post]
func (h *RepositoryHandler) RegisterRepository(c *gin.Context) {
	var req RegisterRepositoryRequest
	requestID := c.GetString("request_id") // ミドルウェアから設定されたリクエストID

	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Error("Invalid request body", "request_id", requestID, "error", err.Error())
		c.JSON(http.StatusBadRequest, ErrorResponse{Code: "INVALID_REQUEST", Message: "Invalid request body: " + err.Error()})
		return
	}

	h.logger.Info("Registering repository", "request_id", requestID, "url", req.URL)
	newRepo, err := h.repoUseCase.Register(c.Request.Context(), req.URL)

	if err != nil {
		if errors.Is(err, repository.ErrRepositoryAlreadyExists) {
			h.logger.Warn("Repository already exists", "request_id", requestID, "url", req.URL)
			c.JSON(http.StatusConflict, ErrorResponse{Code: "CONFLICT", Message: err.Error()})
		} else if errors.Is(err, repository.ErrInvalidRepositoryURL) {
			h.logger.Warn("Invalid repository URL", "request_id", requestID, "url", req.URL)
			c.JSON(http.StatusBadRequest, ErrorResponse{Code: "INVALID_URL", Message: err.Error()})
		} else {
			h.logger.Error("Failed to register repository", "request_id", requestID, "error", err.Error())
			c.JSON(http.StatusInternalServerError, ErrorResponse{Code: "INTERNAL_ERROR", Message: "Failed to register repository"})
		}
		return
	}

	h.logger.Info("Repository registered successfully", "request_id", requestID, "repo_id", newRepo.ID())
	response := RepositoryResponse{
		ID:        newRepo.ID(),
		Name:      newRepo.Name(),
		URL:       newRepo.URL(),
		CreatedAt: newRepo.CreatedAt(),
		UpdatedAt: newRepo.UpdatedAt(),
	}
	c.JSON(http.StatusCreated, response)
}

// ListRepositoryFiles godoc
// @Summary List files in a repository
// @Description Retrieves a list of files and directories within a specified repository.
// @Tags repositories
// @Produce  json
// @Param   repoId path string true "Repository ID" example:"a1b2c3d4-e5f6-7890-1234-567890abcdef"
// @Success 200 {object} ListFilesResponse "Successfully retrieved file list"
// @Failure 400 {object} ErrorResponse "Invalid repository ID format"
// @Failure 404 {object} ErrorResponse "Repository not found"
// @Failure 500 {object} ErrorResponse "Internal server error"
// @Router /repositories/{repoId}/files [get]
func (h *RepositoryHandler) ListRepositoryFiles(c *gin.Context) {
	repoId := c.Param("repoId")
	requestID := c.GetString("request_id") // ミドルウェアから設定されたリクエストID

	if repoId == "" {
		h.logger.Warn("Missing repository ID", "request_id", requestID)
		c.JSON(http.StatusBadRequest, ErrorResponse{Code: "INVALID_ID", Message: "Repository ID is required"})
		return
	}

	h.logger.Info("Listing repository files", "request_id", requestID, "repo_id", repoId)
	// Call the use case which now returns []*model.FileNode
	domainFiles, err := h.repoUseCase.ListFiles(c.Request.Context(), repoId)

	if err != nil {
		if errors.Is(err, repository.ErrRepositoryNotFound) {
			h.logger.Warn("Repository not found", "request_id", requestID, "repo_id", repoId)
			c.JSON(http.StatusNotFound, ErrorResponse{Code: "NOT_FOUND", Message: err.Error()})
		} else {
			h.logger.Error("Failed to list repository files", "request_id", requestID, "repo_id", repoId, "error", err.Error())
			c.JSON(http.StatusInternalServerError, ErrorResponse{Code: "INTERNAL_ERROR", Message: "Failed to list repository files"})
		}
		return
	}

	// Map domain model.FileNode to handler FileNode for the response
	responseFiles := make([]FileNode, 0, len(domainFiles))
	for _, df := range domainFiles {
		responseFiles = append(responseFiles, FileNode{
			Path: df.Path(),
			Type: df.Type(),
		})
	}

	h.logger.Info("Successfully listed repository files", "request_id", requestID, "repo_id", repoId, "file_count", len(responseFiles))
	c.JSON(http.StatusOK, ListFilesResponse{Files: responseFiles})
}

// SelectRepositoryFiles godoc
// @Summary Select manageable files in a repository
// @Description Marks specific files within a repository as manageable by OpsCore.
// @Tags repositories
// @Accept  json
// @Produce  json
// @Param   repoId path string true "Repository ID" example:"a1b2c3d4-e5f6-7890-1234-567890abcdef"
// @Param   files body SelectFilesRequest true "List of file paths to select"
// @Success 200 {object} SelectFilesResponse "Files selected successfully"
// @Failure 400 {object} ErrorResponse "Invalid request body or repository ID"
// @Failure 404 {object} ErrorResponse "Repository not found"
// @Failure 500 {object} ErrorResponse "Internal server error"
// @Router /repositories/{repoId}/files/select [post]
func (h *RepositoryHandler) SelectRepositoryFiles(c *gin.Context) {
	repoId := c.Param("repoId")
	requestID := c.GetString("request_id") // ミドルウェアから設定されたリクエストID
	var req SelectFilesRequest

	if repoId == "" {
		h.logger.Warn("Missing repository ID", "request_id", requestID)
		c.JSON(http.StatusBadRequest, ErrorResponse{Code: "INVALID_ID", Message: "Repository ID is required"})
		return
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Error("Invalid request body", "request_id", requestID, "repo_id", repoId, "error", err.Error())
		c.JSON(http.StatusBadRequest, ErrorResponse{Code: "INVALID_REQUEST", Message: "Invalid request body: " + err.Error()})
		return
	}
	if len(req.FilePaths) == 0 {
		h.logger.Warn("Empty file paths", "request_id", requestID, "repo_id", repoId)
		c.JSON(http.StatusBadRequest, ErrorResponse{Code: "INVALID_REQUEST", Message: "filePaths cannot be empty"})
		return
	}

	h.logger.Info("Selecting repository files", "request_id", requestID, "repo_id", repoId, "file_count", len(req.FilePaths))
	err := h.repoUseCase.SelectFiles(c.Request.Context(), repoId, req.FilePaths)

	if err != nil {
		if errors.Is(err, repository.ErrRepositoryNotFound) {
			h.logger.Warn("Repository not found", "request_id", requestID, "repo_id", repoId)
			c.JSON(http.StatusNotFound, ErrorResponse{Code: "NOT_FOUND", Message: err.Error()})
		} else {
			h.logger.Error("Failed to select files", "request_id", requestID, "repo_id", repoId, "error", err.Error())
			c.JSON(http.StatusInternalServerError, ErrorResponse{Code: "INTERNAL_ERROR", Message: "Failed to select files"})
		}
		return
	}

	h.logger.Info("Files selected successfully", "request_id", requestID, "repo_id", repoId, "file_count", len(req.FilePaths))
	c.JSON(http.StatusOK, SelectFilesResponse{
		Message:       "Files selected successfully",
		RepoID:        repoId,
		SelectedFiles: len(req.FilePaths),
	})
}

// GetSelectedMarkdown godoc
// @Summary Get selected Markdown content from a repository
// @Description Retrieves the concatenated content of all selected Markdown files for a given repository.
// @Tags repositories
// @Produce  json
// @Param   repoId path string true "Repository ID" example:"a1b2c3d4-e5f6-7890-1234-567890abcdef"
// @Success 200 {object} GetMarkdownResponse "Successfully retrieved Markdown content"
// @Failure 400 {object} ErrorResponse "Invalid repository ID format"
// @Failure 404 {object} ErrorResponse "Repository not found or no files selected"
// @Failure 500 {object} ErrorResponse "Internal server error"
// @Router /repositories/{repoId}/markdown [get]
func (h *RepositoryHandler) GetSelectedMarkdown(c *gin.Context) {
	repoId := c.Param("repoId")
	requestID := c.GetString("request_id") // ミドルウェアから設定されたリクエストID

	if repoId == "" {
		h.logger.Warn("Missing repository ID", "request_id", requestID)
		c.JSON(http.StatusBadRequest, ErrorResponse{Code: "INVALID_ID", Message: "Repository ID is required"})
		return
	}

	h.logger.Info("Getting selected Markdown content", "request_id", requestID, "repo_id", repoId)
	markdownContent, err := h.repoUseCase.GetSelectedMarkdown(c.Request.Context(), repoId)

	if err != nil {
		if errors.Is(err, repository.ErrRepositoryNotFound) {
			h.logger.Warn("Repository not found", "request_id", requestID, "repo_id", repoId)
			c.JSON(http.StatusNotFound, ErrorResponse{Code: "NOT_FOUND", Message: err.Error()})
		} else {
			h.logger.Error("Failed to retrieve Markdown content", "request_id", requestID, "repo_id", repoId, "error", err.Error())
			c.JSON(http.StatusInternalServerError, ErrorResponse{Code: "INTERNAL_ERROR", Message: "Failed to retrieve Markdown content"})
		}
		return
	}

	h.logger.Info("Successfully retrieved Markdown content", "request_id", requestID, "repo_id", repoId, "content_length", len(markdownContent))
	c.JSON(http.StatusOK, GetMarkdownResponse{
		RepoID:  repoId,
		Content: markdownContent,
	})
}

// ListRepositories godoc
// @Summary List all repositories
// @Description Retrieves a list of all repositories registered in OpsCore
// @Tags repositories
// @Produce json
// @Success 200 {object} ListRepositoriesResponse "Successfully retrieved repositories"
// @Failure 500 {object} ErrorResponse "Internal server error"
// @Router /repositories [get]
func (h *RepositoryHandler) ListRepositories(c *gin.Context) {
	requestID := c.GetString("request_id") // ミドルウェアから設定されたリクエストID

	h.logger.Info("Listing all repositories", "request_id", requestID)

	repos, err := h.repoUseCase.ListRepositories(c.Request.Context())
	if err != nil {
		h.logger.Error("Failed to list repositories", "request_id", requestID, "error", err.Error())
		c.JSON(http.StatusInternalServerError, ErrorResponse{Code: "INTERNAL_ERROR", Message: "Failed to retrieve repositories"})
		return
	}

	// Map domain models to response DTOs
	repoResponses := make([]RepositoryResponse, 0, len(repos))
	for _, repo := range repos {
		repoResponses = append(repoResponses, RepositoryResponse{
			ID:        repo.ID(),
			Name:      repo.Name(),
			URL:       repo.URL(),
			CreatedAt: repo.CreatedAt(),
			UpdatedAt: repo.UpdatedAt(),
		})
	}

	h.logger.Info("Successfully listed repositories", "request_id", requestID, "repo_count", len(repos))
	c.JSON(http.StatusOK, ListRepositoriesResponse{
		Repositories: repoResponses,
	})
}
