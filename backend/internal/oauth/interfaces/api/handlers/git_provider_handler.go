package handlers

import (
	"net/http"

	"opscore/backend/internal/oauth/application/service"
	"opscore/backend/internal/oauth/domain"

	"github.com/gin-gonic/gin"
)

// GitProviderHandler handles Git provider related HTTP requests
type GitProviderHandler struct {
	gitProviderService *service.GitProviderService
	logger             domain.Logger
}

// NewGitProviderHandler creates a new GitProviderHandler
func NewGitProviderHandler(gitProviderService *service.GitProviderService, logger domain.Logger) *GitProviderHandler {
	return &GitProviderHandler{
		gitProviderService: gitProviderService,
		logger:             logger,
	}
}

// ListUserRepositories lists repositories accessible to the user
// @Summary List user repositories from Git provider
// @Description List all repositories accessible to the authenticated user from a Git provider
// @Tags GitProvider
// @Produce json
// @Security BearerAuth
// @Param provider path string true "Provider name (github, gitlab)"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{} "Invalid provider"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /git-providers/{provider}/repositories [get]
func (h *GitProviderHandler) ListUserRepositories(c *gin.Context) {
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

	repos, err := h.gitProviderService.ListUserRepositories(c.Request.Context(), userID.(string), provider)
	if err != nil {
		h.logger.Error("Failed to list repositories", "error", err, "provider", provider)
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
