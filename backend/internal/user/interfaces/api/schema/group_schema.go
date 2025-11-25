package schema

import "time"

// CreateGroupRequest represents the API request body for creating a group
type CreateGroupRequest struct {
	Name        string `json:"name" binding:"required" example:"Engineering"`
	Description string `json:"description" example:"Engineering team group"`
}

// UpdateGroupRequest represents the API request body for updating a group
type UpdateGroupRequest struct {
	Name        string `json:"name" binding:"required" example:"Engineering"`
	Description string `json:"description" example:"Engineering team group"`
}

// AddMemberRequest represents the API request body for adding a member to a group
type AddMemberRequest struct {
	UserID string `json:"userId" binding:"required" example:"550e8400-e29b-41d4-a716-446655440000"`
}

// RemoveMemberRequest represents the API request body for removing a member from a group
type RemoveMemberRequest struct {
	UserID string `json:"userId" binding:"required" example:"550e8400-e29b-41d4-a716-446655440000"`
}

// GroupResponse represents the API response format for a group
type GroupResponse struct {
	ID          string    `json:"id" example:"550e8400-e29b-41d4-a716-446655440000"`
	Name        string    `json:"name" example:"Engineering"`
	Description string    `json:"description" example:"Engineering team group"`
	MemberIDs   []string  `json:"memberIds" example:"550e8400-e29b-41d4-a716-446655440001,550e8400-e29b-41d4-a716-446655440002"`
	CreatedAt   time.Time `json:"createdAt" example:"2025-04-22T10:00:00Z"`
	UpdatedAt   time.Time `json:"updatedAt" example:"2025-04-22T10:00:00Z"`
}

// ListGroupsResponse represents the API response for listing all groups
type ListGroupsResponse struct {
	Groups []GroupResponse `json:"groups"`
}
