package schema

import "time"

// CreateUserRequest represents the API request body for creating a user
type CreateUserRequest struct {
	Name  string `json:"name" binding:"required" example:"John Doe"`
	Email string `json:"email" binding:"required,email" example:"john.doe@example.com"`
	Role  string `json:"role" binding:"required" example:"user"`
}

// UpdateUserRequest represents the API request body for updating a user
type UpdateUserRequest struct {
	Name  string `json:"name" binding:"required" example:"John Doe"`
	Email string `json:"email" binding:"required,email" example:"john.doe@example.com"`
}

// ChangeRoleRequest represents the API request body for changing user role
type ChangeRoleRequest struct {
	Role string `json:"role" binding:"required" example:"admin"`
}

// JoinGroupRequest represents the API request body for joining a group
type JoinGroupRequest struct {
	GroupID string `json:"groupId" binding:"required" example:"550e8400-e29b-41d4-a716-446655440000"`
}

// LeaveGroupRequest represents the API request body for leaving a group
type LeaveGroupRequest struct {
	GroupID string `json:"groupId" binding:"required" example:"550e8400-e29b-41d4-a716-446655440000"`
}

// UserResponse represents the API response format for a user
type UserResponse struct {
	ID        string    `json:"id" example:"550e8400-e29b-41d4-a716-446655440000"`
	Name      string    `json:"name" example:"John Doe"`
	Email     string    `json:"email" example:"john.doe@example.com"`
	Role      string    `json:"role" example:"user"`
	GroupIDs  []string  `json:"groupIds" example:"550e8400-e29b-41d4-a716-446655440001,550e8400-e29b-41d4-a716-446655440002"`
	CreatedAt time.Time `json:"createdAt" example:"2025-04-22T10:00:00Z"`
	UpdatedAt time.Time `json:"updatedAt" example:"2025-04-22T10:00:00Z"`
}

// ListUsersResponse represents the API response for listing all users
type ListUsersResponse struct {
	Users []UserResponse `json:"users"`
}
