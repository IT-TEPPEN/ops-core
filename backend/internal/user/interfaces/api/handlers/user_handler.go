package handlers

import (
	"net/http"

	"opscore/backend/internal/user/application/usecase"
	intererror "opscore/backend/internal/user/interfaces/error"
	"opscore/backend/internal/user/interfaces/api/schema"

	"github.com/gin-gonic/gin"
)

// UserHandler holds dependencies for user handlers
type UserHandler struct {
	userUseCase usecase.UserUseCase
	logger      Logger
}

// Logger interface defines the minimal logging interface
type Logger interface {
	Info(msg string, args ...any)
	Error(msg string, args ...any)
	Debug(msg string, args ...any)
	Warn(msg string, args ...any)
}

// NewUserHandler creates a new UserHandler
func NewUserHandler(uc usecase.UserUseCase, logger Logger) *UserHandler {
	return &UserHandler{
		userUseCase: uc,
		logger:      logger,
	}
}

// CreateUser godoc
// @Summary Create a new user
// @Description Add a new user to the system
// @Tags users
// @Accept json
// @Produce json
// @Param user body schema.CreateUserRequest true "User information"
// @Success 201 {object} schema.UserResponse "User created successfully"
// @Failure 400 {object} schema.ErrorResponse "Invalid request body"
// @Failure 409 {object} schema.ErrorResponse "User with this email already exists"
// @Failure 500 {object} schema.ErrorResponse "Internal server error"
// @Router /users [post]
func (h *UserHandler) CreateUser(c *gin.Context) {
	var req schema.CreateUserRequest
	requestID := c.GetString("request_id")

	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Error("Invalid request body", "request_id", requestID, "error", err.Error())
		c.JSON(http.StatusBadRequest, schema.ErrorResponse{Code: "INVALID_REQUEST", Message: "Invalid request body: " + err.Error()})
		return
	}

	dtoReq := schema.ToCreateUserDTO(req)

	h.logger.Info("Creating user", "request_id", requestID, "email", dtoReq.Email)
	result, err := h.userUseCase.Create(c.Request.Context(), dtoReq)

	if err != nil {
		httpErr := intererror.MapToHTTPError(err, requestID)
		h.logger.Error("Failed to create user", "request_id", requestID, "error", err.Error(), "http_code", httpErr.Code)
		c.JSON(httpErr.StatusCode, schema.ErrorResponse{Code: httpErr.Code, Message: httpErr.Message})
		return
	}

	h.logger.Info("User created successfully", "request_id", requestID, "user_id", result.ID)
	response := schema.FromUserDTO(*result)
	c.JSON(http.StatusCreated, response)
}

// GetUser godoc
// @Summary Get user details
// @Description Retrieves detailed information about a specific user by ID
// @Tags users
// @Produce json
// @Param userId path string true "User ID" example:"550e8400-e29b-41d4-a716-446655440000"
// @Success 200 {object} schema.UserResponse "Successfully retrieved user details"
// @Failure 400 {object} schema.ErrorResponse "Invalid user ID format"
// @Failure 404 {object} schema.ErrorResponse "User not found"
// @Failure 500 {object} schema.ErrorResponse "Internal server error"
// @Router /users/{userId} [get]
func (h *UserHandler) GetUser(c *gin.Context) {
	userID := c.Param("userId")
	requestID := c.GetString("request_id")

	if userID == "" {
		h.logger.Warn("Missing user ID", "request_id", requestID)
		c.JSON(http.StatusBadRequest, schema.ErrorResponse{Code: "INVALID_ID", Message: "User ID is required"})
		return
	}

	h.logger.Info("Getting user details", "request_id", requestID, "user_id", userID)
	result, err := h.userUseCase.GetByID(c.Request.Context(), userID)

	if err != nil {
		httpErr := intererror.MapToHTTPError(err, requestID)
		h.logger.Error("Failed to get user details", "request_id", requestID, "user_id", userID, "error", err.Error(), "http_code", httpErr.Code)
		c.JSON(httpErr.StatusCode, schema.ErrorResponse{Code: httpErr.Code, Message: httpErr.Message})
		return
	}

	h.logger.Info("Successfully retrieved user details", "request_id", requestID, "user_id", userID)
	response := schema.FromUserDTO(*result)
	c.JSON(http.StatusOK, response)
}

// ListUsers godoc
// @Summary List all users
// @Description Retrieves a list of all users in the system
// @Tags users
// @Produce json
// @Success 200 {object} schema.ListUsersResponse "Successfully retrieved users"
// @Failure 500 {object} schema.ErrorResponse "Internal server error"
// @Router /users [get]
func (h *UserHandler) ListUsers(c *gin.Context) {
	requestID := c.GetString("request_id")

	h.logger.Info("Listing all users", "request_id", requestID)

	results, err := h.userUseCase.GetAll(c.Request.Context())
	if err != nil {
		httpErr := intererror.MapToHTTPError(err, requestID)
		h.logger.Error("Failed to list users", "request_id", requestID, "error", err.Error(), "http_code", httpErr.Code)
		c.JSON(httpErr.StatusCode, schema.ErrorResponse{Code: httpErr.Code, Message: httpErr.Message})
		return
	}

	h.logger.Info("Successfully listed users", "request_id", requestID, "user_count", len(results))
	schemaResponses := schema.FromUserListDTO(results)
	c.JSON(http.StatusOK, schema.ListUsersResponse{Users: schemaResponses})
}

// UpdateUser godoc
// @Summary Update user details
// @Description Updates the profile information of a specific user
// @Tags users
// @Accept json
// @Produce json
// @Param userId path string true "User ID" example:"550e8400-e29b-41d4-a716-446655440000"
// @Param user body schema.UpdateUserRequest true "User information"
// @Success 200 {object} schema.UserResponse "User updated successfully"
// @Failure 400 {object} schema.ErrorResponse "Invalid request body or user ID"
// @Failure 404 {object} schema.ErrorResponse "User not found"
// @Failure 409 {object} schema.ErrorResponse "Email already in use"
// @Failure 500 {object} schema.ErrorResponse "Internal server error"
// @Router /users/{userId} [put]
func (h *UserHandler) UpdateUser(c *gin.Context) {
	userID := c.Param("userId")
	requestID := c.GetString("request_id")
	var req schema.UpdateUserRequest

	if userID == "" {
		h.logger.Warn("Missing user ID", "request_id", requestID)
		c.JSON(http.StatusBadRequest, schema.ErrorResponse{Code: "INVALID_ID", Message: "User ID is required"})
		return
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Error("Invalid request body", "request_id", requestID, "user_id", userID, "error", err.Error())
		c.JSON(http.StatusBadRequest, schema.ErrorResponse{Code: "INVALID_REQUEST", Message: "Invalid request body: " + err.Error()})
		return
	}

	dtoReq := schema.ToUpdateUserDTO(req)

	h.logger.Info("Updating user", "request_id", requestID, "user_id", userID)
	result, err := h.userUseCase.Update(c.Request.Context(), userID, dtoReq)

	if err != nil {
		httpErr := intererror.MapToHTTPError(err, requestID)
		h.logger.Error("Failed to update user", "request_id", requestID, "user_id", userID, "error", err.Error(), "http_code", httpErr.Code)
		c.JSON(httpErr.StatusCode, schema.ErrorResponse{Code: httpErr.Code, Message: httpErr.Message})
		return
	}

	h.logger.Info("User updated successfully", "request_id", requestID, "user_id", userID)
	response := schema.FromUserDTO(*result)
	c.JSON(http.StatusOK, response)
}

// DeleteUser godoc
// @Summary Delete a user
// @Description Removes a user from the system
// @Tags users
// @Produce json
// @Param userId path string true "User ID" example:"550e8400-e29b-41d4-a716-446655440000"
// @Success 200 {object} map[string]string "User deleted successfully"
// @Failure 400 {object} schema.ErrorResponse "Invalid user ID"
// @Failure 404 {object} schema.ErrorResponse "User not found"
// @Failure 500 {object} schema.ErrorResponse "Internal server error"
// @Router /users/{userId} [delete]
func (h *UserHandler) DeleteUser(c *gin.Context) {
	userID := c.Param("userId")
	requestID := c.GetString("request_id")

	if userID == "" {
		h.logger.Warn("Missing user ID", "request_id", requestID)
		c.JSON(http.StatusBadRequest, schema.ErrorResponse{Code: "INVALID_ID", Message: "User ID is required"})
		return
	}

	h.logger.Info("Deleting user", "request_id", requestID, "user_id", userID)
	err := h.userUseCase.Delete(c.Request.Context(), userID)

	if err != nil {
		httpErr := intererror.MapToHTTPError(err, requestID)
		h.logger.Error("Failed to delete user", "request_id", requestID, "user_id", userID, "error", err.Error(), "http_code", httpErr.Code)
		c.JSON(httpErr.StatusCode, schema.ErrorResponse{Code: httpErr.Code, Message: httpErr.Message})
		return
	}

	h.logger.Info("User deleted successfully", "request_id", requestID, "user_id", userID)
	c.JSON(http.StatusOK, map[string]string{
		"message": "User deleted successfully",
		"userId":  userID,
	})
}

// ChangeUserRole godoc
// @Summary Change user role
// @Description Changes the role of a specific user
// @Tags users
// @Accept json
// @Produce json
// @Param userId path string true "User ID" example:"550e8400-e29b-41d4-a716-446655440000"
// @Param role body schema.ChangeRoleRequest true "Role information"
// @Success 200 {object} schema.UserResponse "Role changed successfully"
// @Failure 400 {object} schema.ErrorResponse "Invalid request body or user ID"
// @Failure 404 {object} schema.ErrorResponse "User not found"
// @Failure 500 {object} schema.ErrorResponse "Internal server error"
// @Router /users/{userId}/role [put]
func (h *UserHandler) ChangeUserRole(c *gin.Context) {
	userID := c.Param("userId")
	requestID := c.GetString("request_id")
	var req schema.ChangeRoleRequest

	if userID == "" {
		h.logger.Warn("Missing user ID", "request_id", requestID)
		c.JSON(http.StatusBadRequest, schema.ErrorResponse{Code: "INVALID_ID", Message: "User ID is required"})
		return
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Error("Invalid request body", "request_id", requestID, "user_id", userID, "error", err.Error())
		c.JSON(http.StatusBadRequest, schema.ErrorResponse{Code: "INVALID_REQUEST", Message: "Invalid request body: " + err.Error()})
		return
	}

	dtoReq := schema.ToChangeRoleDTO(req)

	h.logger.Info("Changing user role", "request_id", requestID, "user_id", userID, "role", dtoReq.Role)
	result, err := h.userUseCase.ChangeRole(c.Request.Context(), userID, dtoReq)

	if err != nil {
		httpErr := intererror.MapToHTTPError(err, requestID)
		h.logger.Error("Failed to change user role", "request_id", requestID, "user_id", userID, "error", err.Error(), "http_code", httpErr.Code)
		c.JSON(httpErr.StatusCode, schema.ErrorResponse{Code: httpErr.Code, Message: httpErr.Message})
		return
	}

	h.logger.Info("User role changed successfully", "request_id", requestID, "user_id", userID)
	response := schema.FromUserDTO(*result)
	c.JSON(http.StatusOK, response)
}

// JoinGroup godoc
// @Summary Add user to group
// @Description Adds a user to a specific group
// @Tags users
// @Accept json
// @Produce json
// @Param userId path string true "User ID" example:"550e8400-e29b-41d4-a716-446655440000"
// @Param group body schema.JoinGroupRequest true "Group information"
// @Success 200 {object} schema.UserResponse "User added to group successfully"
// @Failure 400 {object} schema.ErrorResponse "Invalid request body or user ID"
// @Failure 404 {object} schema.ErrorResponse "User or group not found"
// @Failure 409 {object} schema.ErrorResponse "User is already a member of the group"
// @Failure 500 {object} schema.ErrorResponse "Internal server error"
// @Router /users/{userId}/groups [post]
func (h *UserHandler) JoinGroup(c *gin.Context) {
	userID := c.Param("userId")
	requestID := c.GetString("request_id")
	var req schema.JoinGroupRequest

	if userID == "" {
		h.logger.Warn("Missing user ID", "request_id", requestID)
		c.JSON(http.StatusBadRequest, schema.ErrorResponse{Code: "INVALID_ID", Message: "User ID is required"})
		return
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Error("Invalid request body", "request_id", requestID, "user_id", userID, "error", err.Error())
		c.JSON(http.StatusBadRequest, schema.ErrorResponse{Code: "INVALID_REQUEST", Message: "Invalid request body: " + err.Error()})
		return
	}

	dtoReq := schema.ToJoinGroupDTO(req)

	h.logger.Info("Adding user to group", "request_id", requestID, "user_id", userID, "group_id", dtoReq.GroupID)
	result, err := h.userUseCase.JoinGroup(c.Request.Context(), userID, dtoReq)

	if err != nil {
		httpErr := intererror.MapToHTTPError(err, requestID)
		h.logger.Error("Failed to add user to group", "request_id", requestID, "user_id", userID, "error", err.Error(), "http_code", httpErr.Code)
		c.JSON(httpErr.StatusCode, schema.ErrorResponse{Code: httpErr.Code, Message: httpErr.Message})
		return
	}

	h.logger.Info("User added to group successfully", "request_id", requestID, "user_id", userID)
	response := schema.FromUserDTO(*result)
	c.JSON(http.StatusOK, response)
}

// LeaveGroup godoc
// @Summary Remove user from group
// @Description Removes a user from a specific group
// @Tags users
// @Accept json
// @Produce json
// @Param userId path string true "User ID" example:"550e8400-e29b-41d4-a716-446655440000"
// @Param group body schema.LeaveGroupRequest true "Group information"
// @Success 200 {object} schema.UserResponse "User removed from group successfully"
// @Failure 400 {object} schema.ErrorResponse "Invalid request body or user ID"
// @Failure 404 {object} schema.ErrorResponse "User or group not found"
// @Failure 409 {object} schema.ErrorResponse "User is not a member of the group"
// @Failure 500 {object} schema.ErrorResponse "Internal server error"
// @Router /users/{userId}/groups [delete]
func (h *UserHandler) LeaveGroup(c *gin.Context) {
	userID := c.Param("userId")
	requestID := c.GetString("request_id")
	var req schema.LeaveGroupRequest

	if userID == "" {
		h.logger.Warn("Missing user ID", "request_id", requestID)
		c.JSON(http.StatusBadRequest, schema.ErrorResponse{Code: "INVALID_ID", Message: "User ID is required"})
		return
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Error("Invalid request body", "request_id", requestID, "user_id", userID, "error", err.Error())
		c.JSON(http.StatusBadRequest, schema.ErrorResponse{Code: "INVALID_REQUEST", Message: "Invalid request body: " + err.Error()})
		return
	}

	dtoReq := schema.ToLeaveGroupDTO(req)

	h.logger.Info("Removing user from group", "request_id", requestID, "user_id", userID, "group_id", dtoReq.GroupID)
	result, err := h.userUseCase.LeaveGroup(c.Request.Context(), userID, dtoReq)

	if err != nil {
		httpErr := intererror.MapToHTTPError(err, requestID)
		h.logger.Error("Failed to remove user from group", "request_id", requestID, "user_id", userID, "error", err.Error(), "http_code", httpErr.Code)
		c.JSON(httpErr.StatusCode, schema.ErrorResponse{Code: httpErr.Code, Message: httpErr.Message})
		return
	}

	h.logger.Info("User removed from group successfully", "request_id", requestID, "user_id", userID)
	response := schema.FromUserDTO(*result)
	c.JSON(http.StatusOK, response)
}
