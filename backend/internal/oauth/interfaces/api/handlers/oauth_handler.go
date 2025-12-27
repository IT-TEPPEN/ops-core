package handlers

import (
	"net/http"

	"opscore/backend/internal/oauth/application/service"
	"opscore/backend/internal/oauth/domain"

	"github.com/gin-gonic/gin"
)

// OAuthHandler handles OAuth-related HTTP requests
type OAuthHandler struct {
	oauthService *service.OAuthService
	logger       domain.Logger
}

// NewOAuthHandler creates a new OAuthHandler
func NewOAuthHandler(oauthService *service.OAuthService, logger domain.Logger) *OAuthHandler {
	return &OAuthHandler{
		oauthService: oauthService,
		logger:       logger,
	}
}

// HandleCallback handles OAuth callback request
// @Summary OAuth callback endpoint
// @Description Exchanges OAuth authorization code for access token
// @Tags OAuth
// @Accept json
// @Produce json
// @Param request body domain.OAuthCallbackRequest true "OAuth callback request"
// @Success 200 {object} domain.OAuthTokenResponse
// @Failure 400 {object} map[string]interface{} "Invalid request"
// @Failure 500 {object} map[string]interface{} "Internal server error"
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
