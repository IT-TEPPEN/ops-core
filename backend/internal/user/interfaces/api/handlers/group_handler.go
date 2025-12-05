package handlers

import (
	"net/http"

	"opscore/backend/internal/user/application/usecase"
	intererror "opscore/backend/internal/user/interfaces/error"
	"opscore/backend/internal/user/interfaces/api/schema"

	"github.com/gin-gonic/gin"
)

// GroupHandler holds dependencies for group handlers
type GroupHandler struct {
	groupUseCase usecase.GroupUseCase
	logger       Logger
}

// NewGroupHandler creates a new GroupHandler
func NewGroupHandler(uc usecase.GroupUseCase, logger Logger) *GroupHandler {
	return &GroupHandler{
		groupUseCase: uc,
		logger:       logger,
	}
}

// CreateGroup godoc
// @Summary Create a new group
// @Description Add a new group to the system
// @Tags groups
// @Accept json
// @Produce json
// @Param group body schema.CreateGroupRequest true "Group information"
// @Success 201 {object} schema.GroupResponse "Group created successfully"
// @Failure 400 {object} schema.ErrorResponse "Invalid request body"
// @Failure 500 {object} schema.ErrorResponse "Internal server error"
// @Router /groups [post]
func (h *GroupHandler) CreateGroup(c *gin.Context) {
	var req schema.CreateGroupRequest
	requestID := c.GetString("request_id")

	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Error("Invalid request body", "request_id", requestID, "error", err.Error())
		c.JSON(http.StatusBadRequest, schema.ErrorResponse{Code: "INVALID_REQUEST", Message: "Invalid request body: " + err.Error()})
		return
	}

	dtoReq := schema.ToCreateGroupDTO(req)

	h.logger.Info("Creating group", "request_id", requestID, "name", dtoReq.Name)
	result, err := h.groupUseCase.Create(c.Request.Context(), dtoReq)

	if err != nil {
		httpErr := intererror.MapToHTTPError(err, requestID)
		h.logger.Error("Failed to create group", "request_id", requestID, "error", err.Error(), "http_code", httpErr.Code)
		c.JSON(httpErr.StatusCode, schema.ErrorResponse{Code: httpErr.Code, Message: httpErr.Message})
		return
	}

	h.logger.Info("Group created successfully", "request_id", requestID, "group_id", result.ID)
	response := schema.FromGroupDTO(*result)
	c.JSON(http.StatusCreated, response)
}

// GetGroup godoc
// @Summary Get group details
// @Description Retrieves detailed information about a specific group by ID
// @Tags groups
// @Produce json
// @Param groupId path string true "Group ID" example:"550e8400-e29b-41d4-a716-446655440000"
// @Success 200 {object} schema.GroupResponse "Successfully retrieved group details"
// @Failure 400 {object} schema.ErrorResponse "Invalid group ID format"
// @Failure 404 {object} schema.ErrorResponse "Group not found"
// @Failure 500 {object} schema.ErrorResponse "Internal server error"
// @Router /groups/{groupId} [get]
func (h *GroupHandler) GetGroup(c *gin.Context) {
	groupID := c.Param("groupId")
	requestID := c.GetString("request_id")

	if groupID == "" {
		h.logger.Warn("Missing group ID", "request_id", requestID)
		c.JSON(http.StatusBadRequest, schema.ErrorResponse{Code: "INVALID_ID", Message: "Group ID is required"})
		return
	}

	h.logger.Info("Getting group details", "request_id", requestID, "group_id", groupID)
	result, err := h.groupUseCase.GetByID(c.Request.Context(), groupID)

	if err != nil {
		httpErr := intererror.MapToHTTPError(err, requestID)
		h.logger.Error("Failed to get group details", "request_id", requestID, "group_id", groupID, "error", err.Error(), "http_code", httpErr.Code)
		c.JSON(httpErr.StatusCode, schema.ErrorResponse{Code: httpErr.Code, Message: httpErr.Message})
		return
	}

	h.logger.Info("Successfully retrieved group details", "request_id", requestID, "group_id", groupID)
	response := schema.FromGroupDTO(*result)
	c.JSON(http.StatusOK, response)
}

// ListGroups godoc
// @Summary List all groups
// @Description Retrieves a list of all groups in the system
// @Tags groups
// @Produce json
// @Success 200 {object} schema.ListGroupsResponse "Successfully retrieved groups"
// @Failure 500 {object} schema.ErrorResponse "Internal server error"
// @Router /groups [get]
func (h *GroupHandler) ListGroups(c *gin.Context) {
	requestID := c.GetString("request_id")

	h.logger.Info("Listing all groups", "request_id", requestID)

	results, err := h.groupUseCase.GetAll(c.Request.Context())
	if err != nil {
		httpErr := intererror.MapToHTTPError(err, requestID)
		h.logger.Error("Failed to list groups", "request_id", requestID, "error", err.Error(), "http_code", httpErr.Code)
		c.JSON(httpErr.StatusCode, schema.ErrorResponse{Code: httpErr.Code, Message: httpErr.Message})
		return
	}

	h.logger.Info("Successfully listed groups", "request_id", requestID, "group_count", len(results))
	schemaResponses := schema.FromGroupListDTO(results)
	c.JSON(http.StatusOK, schema.ListGroupsResponse{Groups: schemaResponses})
}

// UpdateGroup godoc
// @Summary Update group details
// @Description Updates the information of a specific group
// @Tags groups
// @Accept json
// @Produce json
// @Param groupId path string true "Group ID" example:"550e8400-e29b-41d4-a716-446655440000"
// @Param group body schema.UpdateGroupRequest true "Group information"
// @Success 200 {object} schema.GroupResponse "Group updated successfully"
// @Failure 400 {object} schema.ErrorResponse "Invalid request body or group ID"
// @Failure 404 {object} schema.ErrorResponse "Group not found"
// @Failure 500 {object} schema.ErrorResponse "Internal server error"
// @Router /groups/{groupId} [put]
func (h *GroupHandler) UpdateGroup(c *gin.Context) {
	groupID := c.Param("groupId")
	requestID := c.GetString("request_id")
	var req schema.UpdateGroupRequest

	if groupID == "" {
		h.logger.Warn("Missing group ID", "request_id", requestID)
		c.JSON(http.StatusBadRequest, schema.ErrorResponse{Code: "INVALID_ID", Message: "Group ID is required"})
		return
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Error("Invalid request body", "request_id", requestID, "group_id", groupID, "error", err.Error())
		c.JSON(http.StatusBadRequest, schema.ErrorResponse{Code: "INVALID_REQUEST", Message: "Invalid request body: " + err.Error()})
		return
	}

	dtoReq := schema.ToUpdateGroupDTO(req)

	h.logger.Info("Updating group", "request_id", requestID, "group_id", groupID)
	result, err := h.groupUseCase.Update(c.Request.Context(), groupID, dtoReq)

	if err != nil {
		httpErr := intererror.MapToHTTPError(err, requestID)
		h.logger.Error("Failed to update group", "request_id", requestID, "group_id", groupID, "error", err.Error(), "http_code", httpErr.Code)
		c.JSON(httpErr.StatusCode, schema.ErrorResponse{Code: httpErr.Code, Message: httpErr.Message})
		return
	}

	h.logger.Info("Group updated successfully", "request_id", requestID, "group_id", groupID)
	response := schema.FromGroupDTO(*result)
	c.JSON(http.StatusOK, response)
}

// DeleteGroup godoc
// @Summary Delete a group
// @Description Removes a group from the system
// @Tags groups
// @Produce json
// @Param groupId path string true "Group ID" example:"550e8400-e29b-41d4-a716-446655440000"
// @Success 200 {object} map[string]string "Group deleted successfully"
// @Failure 400 {object} schema.ErrorResponse "Invalid group ID"
// @Failure 404 {object} schema.ErrorResponse "Group not found"
// @Failure 500 {object} schema.ErrorResponse "Internal server error"
// @Router /groups/{groupId} [delete]
func (h *GroupHandler) DeleteGroup(c *gin.Context) {
	groupID := c.Param("groupId")
	requestID := c.GetString("request_id")

	if groupID == "" {
		h.logger.Warn("Missing group ID", "request_id", requestID)
		c.JSON(http.StatusBadRequest, schema.ErrorResponse{Code: "INVALID_ID", Message: "Group ID is required"})
		return
	}

	h.logger.Info("Deleting group", "request_id", requestID, "group_id", groupID)
	err := h.groupUseCase.Delete(c.Request.Context(), groupID)

	if err != nil {
		httpErr := intererror.MapToHTTPError(err, requestID)
		h.logger.Error("Failed to delete group", "request_id", requestID, "group_id", groupID, "error", err.Error(), "http_code", httpErr.Code)
		c.JSON(httpErr.StatusCode, schema.ErrorResponse{Code: httpErr.Code, Message: httpErr.Message})
		return
	}

	h.logger.Info("Group deleted successfully", "request_id", requestID, "group_id", groupID)
	c.JSON(http.StatusOK, map[string]string{
		"message": "Group deleted successfully",
		"groupId": groupID,
	})
}

// AddMember godoc
// @Summary Add member to group
// @Description Adds a user to a specific group
// @Tags groups
// @Accept json
// @Produce json
// @Param groupId path string true "Group ID" example:"550e8400-e29b-41d4-a716-446655440000"
// @Param member body schema.AddMemberRequest true "Member information"
// @Success 200 {object} schema.GroupResponse "Member added successfully"
// @Failure 400 {object} schema.ErrorResponse "Invalid request body or group ID"
// @Failure 404 {object} schema.ErrorResponse "Group or user not found"
// @Failure 409 {object} schema.ErrorResponse "User is already a member of the group"
// @Failure 500 {object} schema.ErrorResponse "Internal server error"
// @Router /groups/{groupId}/members [post]
func (h *GroupHandler) AddMember(c *gin.Context) {
	groupID := c.Param("groupId")
	requestID := c.GetString("request_id")
	var req schema.AddMemberRequest

	if groupID == "" {
		h.logger.Warn("Missing group ID", "request_id", requestID)
		c.JSON(http.StatusBadRequest, schema.ErrorResponse{Code: "INVALID_ID", Message: "Group ID is required"})
		return
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Error("Invalid request body", "request_id", requestID, "group_id", groupID, "error", err.Error())
		c.JSON(http.StatusBadRequest, schema.ErrorResponse{Code: "INVALID_REQUEST", Message: "Invalid request body: " + err.Error()})
		return
	}

	dtoReq := schema.ToAddMemberDTO(req)

	h.logger.Info("Adding member to group", "request_id", requestID, "group_id", groupID, "user_id", dtoReq.UserID)
	result, err := h.groupUseCase.AddMember(c.Request.Context(), groupID, dtoReq)

	if err != nil {
		httpErr := intererror.MapToHTTPError(err, requestID)
		h.logger.Error("Failed to add member to group", "request_id", requestID, "group_id", groupID, "error", err.Error(), "http_code", httpErr.Code)
		c.JSON(httpErr.StatusCode, schema.ErrorResponse{Code: httpErr.Code, Message: httpErr.Message})
		return
	}

	h.logger.Info("Member added to group successfully", "request_id", requestID, "group_id", groupID)
	response := schema.FromGroupDTO(*result)
	c.JSON(http.StatusOK, response)
}

// RemoveMember godoc
// @Summary Remove member from group
// @Description Removes a user from a specific group
// @Tags groups
// @Accept json
// @Produce json
// @Param groupId path string true "Group ID" example:"550e8400-e29b-41d4-a716-446655440000"
// @Param member body schema.RemoveMemberRequest true "Member information"
// @Success 200 {object} schema.GroupResponse "Member removed successfully"
// @Failure 400 {object} schema.ErrorResponse "Invalid request body or group ID"
// @Failure 404 {object} schema.ErrorResponse "Group or user not found"
// @Failure 409 {object} schema.ErrorResponse "User is not a member of the group"
// @Failure 500 {object} schema.ErrorResponse "Internal server error"
// @Router /groups/{groupId}/members [delete]
func (h *GroupHandler) RemoveMember(c *gin.Context) {
	groupID := c.Param("groupId")
	requestID := c.GetString("request_id")
	var req schema.RemoveMemberRequest

	if groupID == "" {
		h.logger.Warn("Missing group ID", "request_id", requestID)
		c.JSON(http.StatusBadRequest, schema.ErrorResponse{Code: "INVALID_ID", Message: "Group ID is required"})
		return
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Error("Invalid request body", "request_id", requestID, "group_id", groupID, "error", err.Error())
		c.JSON(http.StatusBadRequest, schema.ErrorResponse{Code: "INVALID_REQUEST", Message: "Invalid request body: " + err.Error()})
		return
	}

	dtoReq := schema.ToRemoveMemberDTO(req)

	h.logger.Info("Removing member from group", "request_id", requestID, "group_id", groupID, "user_id", dtoReq.UserID)
	result, err := h.groupUseCase.RemoveMember(c.Request.Context(), groupID, dtoReq)

	if err != nil {
		httpErr := intererror.MapToHTTPError(err, requestID)
		h.logger.Error("Failed to remove member from group", "request_id", requestID, "group_id", groupID, "error", err.Error(), "http_code", httpErr.Code)
		c.JSON(httpErr.StatusCode, schema.ErrorResponse{Code: httpErr.Code, Message: httpErr.Message})
		return
	}

	h.logger.Info("Member removed from group successfully", "request_id", requestID, "group_id", groupID)
	response := schema.FromGroupDTO(*result)
	c.JSON(http.StatusOK, response)
}

// GetUserGroups godoc
// @Summary Get groups for a user
// @Description Retrieves all groups that a specific user is a member of
// @Tags groups
// @Produce json
// @Param userId path string true "User ID" example:"550e8400-e29b-41d4-a716-446655440000"
// @Success 200 {object} schema.ListGroupsResponse "Successfully retrieved user groups"
// @Failure 400 {object} schema.ErrorResponse "Invalid user ID format"
// @Failure 404 {object} schema.ErrorResponse "User not found"
// @Failure 500 {object} schema.ErrorResponse "Internal server error"
// @Router /users/{userId}/groups [get]
func (h *GroupHandler) GetUserGroups(c *gin.Context) {
	userID := c.Param("userId")
	requestID := c.GetString("request_id")

	if userID == "" {
		h.logger.Warn("Missing user ID", "request_id", requestID)
		c.JSON(http.StatusBadRequest, schema.ErrorResponse{Code: "INVALID_ID", Message: "User ID is required"})
		return
	}

	h.logger.Info("Getting groups for user", "request_id", requestID, "user_id", userID)
	results, err := h.groupUseCase.GetByMemberID(c.Request.Context(), userID)

	if err != nil {
		httpErr := intererror.MapToHTTPError(err, requestID)
		h.logger.Error("Failed to get user groups", "request_id", requestID, "user_id", userID, "error", err.Error(), "http_code", httpErr.Code)
		c.JSON(httpErr.StatusCode, schema.ErrorResponse{Code: httpErr.Code, Message: httpErr.Message})
		return
	}

	h.logger.Info("Successfully retrieved user groups", "request_id", requestID, "user_id", userID, "group_count", len(results))
	schemaResponses := schema.FromGroupListDTO(results)
	c.JSON(http.StatusOK, schema.ListGroupsResponse{Groups: schemaResponses})
}
