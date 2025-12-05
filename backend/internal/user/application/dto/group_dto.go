package dto

import "time"

// CreateGroupRequest represents the use case request for creating a group
type CreateGroupRequest struct {
	Name        string
	Description string
}

// UpdateGroupRequest represents the use case request for updating a group
type UpdateGroupRequest struct {
	Name        string
	Description string
}

// GroupResponse represents the use case response for a group
type GroupResponse struct {
	ID          string
	Name        string
	Description string
	MemberIDs   []string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// AddMemberRequest represents the use case request for adding a member to a group
type AddMemberRequest struct {
	UserID string
}

// RemoveMemberRequest represents the use case request for removing a member from a group
type RemoveMemberRequest struct {
	UserID string
}
