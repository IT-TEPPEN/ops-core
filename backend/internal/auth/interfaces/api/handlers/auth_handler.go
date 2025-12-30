package handlers

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"opscore/backend/internal/auth/application/service"
	"opscore/backend/internal/auth/domain"
)

// AuthHandler handles authentication-related HTTP requests
type AuthHandler struct {
	providerFactory    *service.ProviderFactory
	jwtService         *service.JWTService
	userRepository     domain.UserRepository
	userIdentityRepo   domain.UserIdentityRepository
	refreshTokenRepo   domain.RefreshTokenRepository
	logger             Logger
}

// Logger interface for structured logging
type Logger interface {
	Info(msg string, args ...any)
	Error(msg string, args ...any)
	Debug(msg string, args ...any)
	Warn(msg string, args ...any)
}

// NewAuthHandler creates a new AuthHandler
func NewAuthHandler(
	providerFactory *service.ProviderFactory,
	jwtService *service.JWTService,
	userRepository domain.UserRepository,
	userIdentityRepo domain.UserIdentityRepository,
	refreshTokenRepo domain.RefreshTokenRepository,
	logger Logger,
) *AuthHandler {
	return &AuthHandler{
		providerFactory:  providerFactory,
		jwtService:       jwtService,
		userRepository:   userRepository,
		userIdentityRepo: userIdentityRepo,
		refreshTokenRepo: refreshTokenRepo,
		logger:           logger,
	}
}

// ProviderLoginResponse represents the response for initiating OIDC login
type ProviderLoginResponse struct {
	AuthURL string `json:"auth_url"`
	State   string `json:"state"`
}

// ProviderCallbackRequest represents the request body for OIDC callback
type ProviderCallbackRequest struct {
	Code       string `json:"code" binding:"required"`
	State      string `json:"state" binding:"required"`
	RememberMe bool   `json:"remember_me"` // Optional, defaults to false
}

// AuthResponse represents the authentication response
type AuthResponse struct {
	Token        string      `json:"token"`
	RefreshToken string      `json:"refresh_token"`
	User         UserProfile `json:"user"`
}

// RefreshTokenRequest represents the refresh token request
type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

// RefreshTokenResponse represents the refresh token response
type RefreshTokenResponse struct {
	Token        string `json:"token"`
	RefreshToken string `json:"refresh_token"`
}

// UserProfile represents user profile information
type UserProfile struct {
	ID      string `json:"id"`
	Email   string `json:"email"`
	Name    string `json:"name"`
	Picture string `json:"picture"`
}

// IdentityResponse represents a user identity
type IdentityResponse struct {
	ID         string `json:"id"`
	Provider   string `json:"provider"`
	Email      string `json:"email"`
	Name       string `json:"name"`
	LinkedAt   string `json:"linked_at"`
	LastUsedAt string `json:"last_used_at"`
}

// ProviderLogin godoc
// @Summary Initiate OIDC login
// @Description Returns OIDC provider authorization URL for user to authenticate
// @Tags auth
// @Produce json
// @Param provider path string true "Provider name (google, github, gitlab, microsoft)"
// @Success 200 {object} ProviderLoginResponse
// @Router /auth/{provider}/login [get]
func (h *AuthHandler) ProviderLogin(c *gin.Context) {
	provider := c.Param("provider")

	// Create provider instance
	oidcProvider, err := h.providerFactory.CreateProvider(provider)
	if err != nil {
		h.logger.Error("Unsupported provider", "provider", provider, "error", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Unsupported provider: %s", provider)})
		return
	}

	// Generate random state for CSRF protection
	state, err := generateRandomState()
	if err != nil {
		h.logger.Error("Failed to generate state", "error", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate state"})
		return
	}

	authURL := oidcProvider.GetAuthURL(state)

	c.JSON(http.StatusOK, ProviderLoginResponse{
		AuthURL: authURL,
		State:   state,
	})
}

// ProviderCallback godoc
// @Summary Handle OIDC provider callback
// @Description Exchanges authorization code for JWT token and creates/updates user
// @Tags auth
// @Accept json
// @Produce json
// @Param provider path string true "Provider name (google, github, gitlab, microsoft)"
// @Param request body ProviderCallbackRequest true "Authorization code and state"
// @Success 200 {object} AuthResponse
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /auth/{provider}/callback [post]
func (h *AuthHandler) ProviderCallback(c *gin.Context) {
	provider := c.Param("provider")

	var req ProviderCallbackRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Error("Invalid request body", "error", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	// Create provider instance
	oidcProvider, err := h.providerFactory.CreateProvider(provider)
	if err != nil {
		h.logger.Error("Unsupported provider", "provider", provider, "error", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Unsupported provider: %s", provider)})
		return
	}

	// Exchange code for token
	token, err := oidcProvider.ExchangeToken(c.Request.Context(), req.Code)
	if err != nil {
		h.logger.Error("Failed to exchange code for token", "provider", provider, "error", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to authenticate with provider"})
		return
	}

	// Verify ID token and get user info
	oidcUserInfo, err := oidcProvider.VerifyIDToken(c.Request.Context(), token)
	if err != nil {
		h.logger.Error("Failed to verify ID token", "provider", provider, "error", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to verify user identity"})
		return
	}

	// Find existing identity
	identity, err := h.userIdentityRepo.FindByProviderIdentity(c.Request.Context(), provider, oidcUserInfo.Subject)

	var user *domain.User

	if err != nil {
		// Identity doesn't exist, create new user and identity
		userID := uuid.New().String()
		user = domain.NewUser(userID, *oidcUserInfo)

		if err := h.userRepository.Create(c.Request.Context(), user); err != nil {
			h.logger.Error("Failed to create user", "error", err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
			return
		}

		identityID := uuid.New().String()
		identity = domain.NewUserIdentity(identityID, userID, *oidcUserInfo)

		if err := h.userIdentityRepo.Create(c.Request.Context(), identity); err != nil {
			h.logger.Error("Failed to create user identity", "error", err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user identity"})
			return
		}

		h.logger.Info("New user created", "user_id", user.ID, "email", user.PrimaryEmail, "provider", provider)
	} else {
		// Identity exists, get user and update last used
		user, err = h.userRepository.FindByID(c.Request.Context(), identity.UserID)
		if err != nil {
			h.logger.Error("Failed to find user", "user_id", identity.UserID, "error", err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to find user"})
			return
		}

		// Update last login and identity last used
		user.UpdateLastLogin()
		identity.UpdateLastUsed()

		if err := h.userRepository.Update(c.Request.Context(), user); err != nil {
			h.logger.Error("Failed to update user", "error", err.Error())
			// Non-fatal error, continue
		}

		if err := h.userIdentityRepo.Update(c.Request.Context(), identity); err != nil {
			h.logger.Error("Failed to update identity", "error", err.Error())
			// Non-fatal error, continue
		}

		h.logger.Info("User logged in", "user_id", user.ID, "email", user.PrimaryEmail, "provider", provider)
	}

	// Generate token pair (access + refresh)
	tokenPair, refreshToken, err := h.jwtService.GenerateTokenPair(c.Request.Context(), user, req.RememberMe)
	if err != nil {
		h.logger.Error("Failed to generate token pair", "error", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate session token"})
		return
	}

	// Save refresh token to database
	if err := h.refreshTokenRepo.Create(c.Request.Context(), refreshToken); err != nil {
		h.logger.Error("Failed to save refresh token", "error", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save session"})
		return
	}

	c.JSON(http.StatusOK, AuthResponse{
		Token:        tokenPair.AccessToken,
		RefreshToken: tokenPair.RefreshToken,
		User: UserProfile{
			ID:      user.ID,
			Email:   user.PrimaryEmail,
			Name:    user.DisplayName,
			Picture: user.PictureURL,
		},
	})
}

// GetMe godoc
// @Summary Get current user information
// @Description Returns the currently authenticated user's profile
// @Tags auth
// @Produce json
// @Security BearerAuth
// @Success 200 {object} UserProfile
// @Failure 401 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /auth/me [get]
func (h *AuthHandler) GetMe(c *gin.Context) {
	// Get user ID from context (set by auth middleware)
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	user, err := h.userRepository.FindByID(c.Request.Context(), userID.(string))
	if err != nil {
		h.logger.Error("Failed to find user", "user_id", userID, "error", err.Error())
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, UserProfile{
		ID:      user.ID,
		Email:   user.PrimaryEmail,
		Name:    user.DisplayName,
		Picture: user.PictureURL,
	})
}

// GetIdentities godoc
// @Summary Get user's linked identities
// @Description Returns all provider identities linked to the current user
// @Tags auth
// @Produce json
// @Security BearerAuth
// @Success 200 {object} map[string][]IdentityResponse
// @Failure 401 {object} map[string]string
// @Router /auth/identities [get]
func (h *AuthHandler) GetIdentities(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	identities, err := h.userIdentityRepo.FindByUserID(c.Request.Context(), userID.(string))
	if err != nil {
		h.logger.Error("Failed to get identities", "user_id", userID, "error", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get identities"})
		return
	}

	response := make([]IdentityResponse, len(identities))
	for i, identity := range identities {
		response[i] = IdentityResponse{
			ID:         identity.ID,
			Provider:   identity.Provider,
			Email:      identity.Email,
			Name:       identity.Name,
			LinkedAt:   identity.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
			LastUsedAt: identity.LastUsedAt.Format("2006-01-02T15:04:05Z07:00"),
		}
	}

	c.JSON(http.StatusOK, gin.H{"identities": response})
}

// UnlinkIdentity godoc
// @Summary Unlink a provider identity
// @Description Removes a provider identity from the user's account
// @Tags auth
// @Produce json
// @Security BearerAuth
// @Param id path string true "Identity ID"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /auth/identities/{id} [delete]
func (h *AuthHandler) UnlinkIdentity(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	identityID := c.Param("id")

	// Check that user has more than one identity
	identities, err := h.userIdentityRepo.FindByUserID(c.Request.Context(), userID.(string))
	if err != nil {
		h.logger.Error("Failed to get identities", "user_id", userID, "error", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get identities"})
		return
	}

	if len(identities) <= 1 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Cannot unlink the last identity"})
		return
	}

	// Verify the identity belongs to the user
	found := false
	for _, identity := range identities {
		if identity.ID == identityID {
			found = true
			break
		}
	}

	if !found {
		c.JSON(http.StatusNotFound, gin.H{"error": "Identity not found"})
		return
	}

	// Delete the identity
	if err := h.userIdentityRepo.Delete(c.Request.Context(), identityID); err != nil {
		h.logger.Error("Failed to delete identity", "identity_id", identityID, "error", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to unlink identity"})
		return
	}

	h.logger.Info("Identity unlinked", "user_id", userID, "identity_id", identityID)
	c.JSON(http.StatusOK, gin.H{"message": "Identity unlinked successfully"})
}

// Logout godoc
// @Summary Logout user
// @Description Revokes all refresh tokens for the user
// @Tags auth
// @Produce json
// @Security BearerAuth
// @Success 200 {object} map[string]string
// @Router /auth/logout [post]
func (h *AuthHandler) Logout(c *gin.Context) {
	// Get user ID from context (set by auth middleware)
	userID, exists := c.Get("user_id")
	if exists {
		// Revoke all refresh tokens for the user
		if err := h.refreshTokenRepo.RevokeByUserID(c.Request.Context(), userID.(string)); err != nil {
			h.logger.Error("Failed to revoke refresh tokens", "user_id", userID, "error", err.Error())
			// Don't return error to client, just log it
		}
		h.logger.Info("User logged out", "user_id", userID)
	}

	c.JSON(http.StatusOK, gin.H{"message": "Logged out successfully"})
}

// RefreshToken godoc
// @Summary Refresh access token
// @Description Exchanges a refresh token for a new access token and refresh token pair
// @Tags auth
// @Accept json
// @Produce json
// @Param request body RefreshTokenRequest true "Refresh token"
// @Success 200 {object} RefreshTokenResponse
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Router /auth/refresh [post]
func (h *AuthHandler) RefreshToken(c *gin.Context) {
	var req RefreshTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Error("Invalid refresh token request", "error", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	// Validate and hash the refresh token
	tokenHash, err := h.jwtService.ValidateRefreshToken(req.RefreshToken)
	if err != nil {
		h.logger.Error("Invalid refresh token format", "error", err.Error())
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid refresh token"})
		return
	}

	// Find refresh token in database
	storedToken, err := h.refreshTokenRepo.FindByTokenHash(c.Request.Context(), tokenHash)
	if err != nil {
		h.logger.Error("Refresh token not found", "error", err.Error())
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid refresh token"})
		return
	}

	// Validate refresh token (not expired, not revoked)
	if !storedToken.IsValid() {
		h.logger.Warn("Refresh token is invalid or revoked", "token_id", storedToken.ID, "user_id", storedToken.UserID)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Refresh token expired or revoked"})
		return
	}

	// Get user
	user, err := h.userRepository.FindByID(c.Request.Context(), storedToken.UserID)
	if err != nil {
		h.logger.Error("Failed to find user", "user_id", storedToken.UserID, "error", err.Error())
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
		return
	}

	// Revoke old refresh token (token rotation)
	storedToken.Revoke()
	if err := h.refreshTokenRepo.Update(c.Request.Context(), storedToken); err != nil {
		h.logger.Error("Failed to revoke old refresh token", "error", err.Error())
		// Continue anyway, token rotation is a security enhancement
	}

	// Generate new token pair
	tokenPair, newRefreshToken, err := h.jwtService.GenerateTokenPair(c.Request.Context(), user, storedToken.RememberMe)
	if err != nil {
		h.logger.Error("Failed to generate new token pair", "error", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate tokens"})
		return
	}

	// Save new refresh token
	if err := h.refreshTokenRepo.Create(c.Request.Context(), newRefreshToken); err != nil {
		h.logger.Error("Failed to save new refresh token", "error", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save session"})
		return
	}

	h.logger.Info("Tokens refreshed", "user_id", user.ID)

	c.JSON(http.StatusOK, RefreshTokenResponse{
		Token:        tokenPair.AccessToken,
		RefreshToken: tokenPair.RefreshToken,
	})
}

// generateRandomState generates a cryptographically secure random state for CSRF protection
func generateRandomState() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", fmt.Errorf("failed to generate random bytes: %w", err)
	}
	return base64.URLEncoding.EncodeToString(b), nil
}
