package handlers

import (
	"net/http"

	"opscore/backend/internal/oauth/application/service"
	"opscore/backend/internal/oauth/domain"

	"github.com/gin-gonic/gin"
)

// OAuthHandler handles OAuth-related HTTP requests
type OAuthHandler struct {
	oauthService       *service.OAuthService
	gitProviderService *service.GitProviderService
	logger             domain.Logger
}

// NewOAuthHandler creates a new OAuthHandler
func NewOAuthHandler(oauthService *service.OAuthService, gitProviderService *service.GitProviderService, logger domain.Logger) *OAuthHandler {
	return &OAuthHandler{
		oauthService:       oauthService,
		gitProviderService: gitProviderService,
		logger:             logger,
	}
}

// OAuthConnectionResponse represents an OAuth connection in API responses
type OAuthConnectionResponse struct {
	ID               string   `json:"id"`
	Provider         string   `json:"provider"`
	ProviderHost     string   `json:"provider_host"`
	ProviderUsername string   `json:"provider_username"`
	Scopes           []string `json:"scopes"`
	ConnectedAt      string   `json:"connected_at"`
}

// HandleCallback handles OAuth callback request
// @Summary OAuth callback endpoint
// @Description Exchanges OAuth authorization code for access token
// @Tags OAuth
// @Accept json
// @Produce json
// @Param request body domain.OAuthCallbackRequest true "OAuth callback request"
// @Success 200 {object} domain.OAuthTokenResponse
// @Failure 400 {object} domain.OAuthTokenResponse "Invalid request"
// @Failure 500 {object} domain.OAuthTokenResponse "Internal server error"
// @Router /auth/oauth/callback [post]
func (h *OAuthHandler) HandleCallback(c *gin.Context) {
	var req domain.OAuthCallbackRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Error("Failed to bind request", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "invalid_request",
			"message": "Invalid request format",
		})
		return
	}

	// Validate provider
	if !req.Provider.IsValid() {
		h.logger.Error("Invalid provider", "provider", req.Provider)
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "invalid_provider",
			"message": "Unsupported Git provider",
		})
		return
	}

	// Validate required fields for self-hosted GitLab
	if req.Provider == domain.ProviderGitLabSelfHosted {
		if req.GitLabURL == "" || req.ClientID == "" || req.ClientSecret == "" {
			h.logger.Error("Missing required fields for self-hosted GitLab")
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   "invalid_request",
				"message": "Fields 'gitlabUrl', 'clientId', and 'clientSecret' are required for self-hosted GitLab",
			})
			return
		}
	}

	// Exchange code for token
	tokenResp, err := h.oauthService.ExchangeToken(req)
	if err != nil {
		h.logger.Error("Failed to exchange token", "error", err, "provider", req.Provider)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "token_exchange_failed",
			"message": err.Error(),
		})
		return
	}

	h.logger.Info("OAuth token exchanged successfully", "provider", req.Provider)
	c.JSON(http.StatusOK, tokenResp)
}

// HandleCallbackWithSave handles OAuth callback and saves tokens to database
// @Summary OAuth callback endpoint with token storage
// @Description Exchanges OAuth authorization code for access token and saves to database
// @Tags OAuth
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body domain.OAuthCallbackRequest true "OAuth callback request"
// @Success 200 {object} OAuthConnectionResponse
// @Failure 400 {object} map[string]interface{} "Invalid request"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /auth/oauth/connect [post]
func (h *OAuthHandler) HandleCallbackWithSave(c *gin.Context) {
	var req domain.OAuthCallbackRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Error("Failed to bind request", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "invalid_request",
			"message": "Invalid request format",
		})
		return
	}

	// Validate provider
	if !req.Provider.IsValid() {
		h.logger.Error("Invalid provider", "provider", req.Provider)
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "invalid_provider",
			"message": "Unsupported Git provider",
		})
		return
	}

	// Get user ID from context (requires authentication middleware)
	userID, exists := c.Get("user_id")
	if !exists {
		h.logger.Error("User ID not found in context")
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":   "unauthorized",
			"message": "Authentication required",
		})
		return
	}

	// Exchange code for token and save
	conn, err := h.oauthService.ExchangeAndSaveToken(c.Request.Context(), userID.(string), req)
	if err != nil {
		h.logger.Error("Failed to exchange and save token", "error", err, "provider", req.Provider)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "token_exchange_failed",
			"message": err.Error(),
		})
		return
	}

	h.logger.Info("OAuth connection saved successfully", "provider", req.Provider, "user_id", userID)
	c.JSON(http.StatusOK, OAuthConnectionResponse{
		ID:               conn.ID(),
		Provider:         string(conn.Provider()),
		ProviderHost:     conn.ProviderHost(),
		ProviderUsername: conn.ProviderUsername(),
		Scopes:           conn.Scopes(),
		ConnectedAt:      conn.CreatedAt().Format("2006-01-02T15:04:05Z07:00"),
	})
}

// ListConnections lists all OAuth connections for the current user
// @Summary List OAuth connections
// @Description List all OAuth connections for the authenticated user
// @Tags OAuth
// @Produce json
// @Security BearerAuth
// @Success 200 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /auth/oauth/connections [get]
func (h *OAuthHandler) ListConnections(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":   "unauthorized",
			"message": "Authentication required",
		})
		return
	}

	connections, err := h.oauthService.ListConnections(c.Request.Context(), userID.(string))
	if err != nil {
		h.logger.Error("Failed to list connections", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "internal_error",
			"message": err.Error(),
		})
		return
	}

	response := make([]OAuthConnectionResponse, len(connections))
	for i, conn := range connections {
		response[i] = OAuthConnectionResponse{
			ID:               conn.ID(),
			Provider:         string(conn.Provider()),
			ProviderHost:     conn.ProviderHost(),
			ProviderUsername: conn.ProviderUsername(),
			Scopes:           conn.Scopes(),
			ConnectedAt:      conn.CreatedAt().Format("2006-01-02T15:04:05Z07:00"),
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"connections": response,
	})
}

// DisconnectProvider removes an OAuth connection
// @Summary Disconnect OAuth provider
// @Description Remove an OAuth connection for the authenticated user
// @Tags OAuth
// @Security BearerAuth
// @Param provider path string true "Provider name (github, gitlab)"
// @Success 204 "No Content"
// @Failure 400 {object} map[string]interface{} "Invalid provider"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /auth/oauth/connections/{provider} [delete]
func (h *OAuthHandler) DisconnectProvider(c *gin.Context) {
	providerStr := c.Param("provider")

	provider := domain.Provider(providerStr)
	if !provider.IsValid() {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "invalid_provider",
			"message": "Invalid provider",
		})
		return
	}

	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":   "unauthorized",
			"message": "Authentication required",
		})
		return
	}

	if err := h.oauthService.DisconnectProvider(c.Request.Context(), userID.(string), provider); err != nil {
		h.logger.Error("Failed to disconnect provider", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "internal_error",
			"message": err.Error(),
		})
		return
	}

	c.Status(http.StatusNoContent)
}

// GetAccessToken returns the access token for a provider (internal use)
func (h *OAuthHandler) GetAccessToken(c *gin.Context) {
	providerStr := c.Param("provider")

	provider := domain.Provider(providerStr)
	if !provider.IsValid() {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "invalid_provider",
			"message": "Invalid provider",
		})
		return
	}

	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":   "unauthorized",
			"message": "Authentication required",
		})
		return
	}

	token, err := h.oauthService.GetAccessToken(c.Request.Context(), userID.(string), provider)
	if err != nil {
		h.logger.Error("Failed to get access token", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "internal_error",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"access_token": token,
	})
}

// ListRepositoriesByConnection lists repositories for a specific OAuth connection
// @Summary List repositories by OAuth connection
// @Description List all repositories accessible via a specific OAuth connection
// @Tags OAuth
// @Produce json
// @Security BearerAuth
// @Param connectionId path string true "OAuth connection ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{} "Invalid connection ID"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /auth/oauth/connections/{connectionId}/repositories [get]
func (h *OAuthHandler) ListRepositoriesByConnection(c *gin.Context) {
	connectionID := c.Param("connectionId")
	if connectionID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "invalid_request",
			"message": "Connection ID is required",
		})
		return
	}

	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":   "unauthorized",
			"message": "Authentication required",
		})
		return
	}

	repos, err := h.gitProviderService.ListUserRepositoriesByConnectionID(c.Request.Context(), userID.(string), connectionID)
	if err != nil {
		h.logger.Error("Failed to list repositories", "error", err, "connectionId", connectionID)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "internal_error",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"repositories": repos,
		"count":        len(repos),
	})
}

// GetRepositoryContents lists files and directories at a specific path in a repository
// @Summary Get repository contents
// @Description Get files and directories at a specific path via OAuth connection
// @Tags OAuth
// @Produce json
// @Security BearerAuth
// @Param connectionId path string true "OAuth connection ID"
// @Param owner path string true "Repository owner"
// @Param repo path string true "Repository name"
// @Param repositoryId query string true "Repository ID"
// @Param path query string false "Directory path (empty for root)"
// @Param commit_sha query string false "Git commit SHA to fetch contents from"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{} "Invalid request"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /auth/oauth/connections/{connectionId}/repositories/{owner}/{repo}/contents [get]
func (h *OAuthHandler) GetRepositoryContents(c *gin.Context) {
	connectionID := c.Param("connectionId")
	owner := c.Param("owner")
	repo := c.Param("repo")
	path := c.DefaultQuery("path", "")
	commitSha := c.DefaultQuery("commit_sha", "")
	repositoryID := c.Query("repositoryId")

	if connectionID == "" || owner == "" || repo == "" || repositoryID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "invalid_request",
			"message": "Connection ID, repository ID, owner, and repo are required",
		})
		return
	}

	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":   "unauthorized",
			"message": "Authentication required",
		})
		return
	}

	contents, err := h.gitProviderService.GetRepositoryContents(c.Request.Context(), userID.(string), connectionID, repositoryID, owner, repo, path, commitSha)
	if err != nil {
		h.logger.Error("Failed to get repository contents", "error", err, "connectionId", connectionID, "owner", owner, "repo", repo, "path", path)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "internal_error",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"files": contents,
		"path":  path,
	})
}

// GetRepositoryFileContent gets the content of a specific file
// @Summary Get file content
// @Description Get the content of a specific file via OAuth connection
// @Tags OAuth
// @Produce json
// @Security BearerAuth
// @Param connectionId path string true "OAuth connection ID"
// @Param owner path string true "Repository owner"
// @Param repo path string true "Repository name"
// @Param repositoryId query string true "Repository ID"
// @Param filePath path string true "File path"
// @Param commit_sha query string false "Git reference (branch, tag, commit SHA)"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{} "Invalid request"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /auth/oauth/connections/{connectionId}/repositories/{owner}/{repo}/files/{filePath} [get]
func (h *OAuthHandler) GetRepositoryFileContent(c *gin.Context) {
	connectionID := c.Param("connectionId")
	owner := c.Param("owner")
	repo := c.Param("repo")
	filePath := c.Param("filePath")
	repositoryID := c.Query("repositoryId")
	commitSha := c.DefaultQuery("commit_sha", "")

	if connectionID == "" || owner == "" || repo == "" || filePath == "" || repositoryID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "invalid_request",
			"message": "Connection ID, repository ID, owner, repo, and file path are required",
		})
		return
	}

	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":   "unauthorized",
			"message": "Authentication required",
		})
		return
	}

	content, err := h.gitProviderService.GetFileContent(c.Request.Context(), userID.(string), connectionID, repositoryID, owner, repo, filePath, commitSha)
	if err != nil {
		h.logger.Error("Failed to get file content", "error", err, "connectionId", connectionID, "owner", owner, "repo", repo, "filePath", filePath)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "internal_error",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"content":  content.Content,
		"path":     filePath,
		"sha":      content.SHA,
		"encoding": content.Encoding,
	})
}

// GetFileCommitHistory gets commit history for a specific file
// @Summary Get file commit history
// @Description Get the commit history for a specific file via OAuth connection
// @Tags OAuth
// @Produce json
// @Security BearerAuth
// @Param connectionId path string true "OAuth connection ID"
// @Param owner path string true "Repository owner"
// @Param repo path string true "Repository name"
// @Param repositoryId query string true "Repository ID"
// @Param path query string true "File path relative to repository root"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{} "Invalid request"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /auth/oauth/connections/{connectionId}/repositories/{owner}/{repo}/files/history [get]
func (h *OAuthHandler) GetFileCommitHistory(c *gin.Context) {
	connectionID := c.Param("connectionId")
	owner := c.Param("owner")
	repo := c.Param("repo")
	filePath := c.Param("filePath")
	repositoryID := c.Query("repositoryId")

	if connectionID == "" || owner == "" || repo == "" || repositoryID == "" || filePath == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "invalid_request",
			"message": "Connection ID, repository ID, owner, repo, and file path are required",
		})
		return
	}

	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":   "unauthorized",
			"message": "Authentication required",
		})
		return
	}

	commits, err := h.gitProviderService.GetFileCommitHistory(c.Request.Context(), userID.(string), connectionID, repositoryID, owner, repo, filePath)
	if err != nil {
		h.logger.Error("Failed to get file commit history", "error", err, "connectionId", connectionID, "owner", owner, "repo", repo, "filePath", filePath)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "internal_error",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"file_path": filePath,
		"commits":   commits,
		"count":     len(commits),
	})
}
