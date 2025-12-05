package dto

import "time"

// CreateUserRequest represents the use case request for creating a user
type CreateUserRequest struct {
	Name  string
	Email string
	Role  string
}

// UpdateUserRequest represents the use case request for updating a user
type UpdateUserRequest struct {
	Name  string
	Email string
}

// ChangeRoleRequest represents the use case request for changing user role
type ChangeRoleRequest struct {
	Role string
}

// UserResponse represents the use case response for a user
type UserResponse struct {
	ID        string
	Name      string
	Email     string
	Role      string
	GroupIDs  []string
	CreatedAt time.Time
	UpdatedAt time.Time
}

// JoinGroupRequest represents the use case request for joining a group
type JoinGroupRequest struct {
	GroupID string
}

// LeaveGroupRequest represents the use case request for leaving a group
type LeaveGroupRequest struct {
	GroupID string
}
