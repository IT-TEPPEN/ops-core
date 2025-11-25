package dto

import (
	"opscore/backend/internal/user/domain/entity"
)

// ToUserResponse converts a User entity to UserResponse DTO
func ToUserResponse(user entity.User) UserResponse {
	groupIDs := make([]string, len(user.GroupIDs()))
	for i, id := range user.GroupIDs() {
		groupIDs[i] = id.String()
	}

	return UserResponse{
		ID:        user.ID().String(),
		Name:      user.Name(),
		Email:     user.Email().String(),
		Role:      user.Role().String(),
		GroupIDs:  groupIDs,
		CreatedAt: user.CreatedAt(),
		UpdatedAt: user.UpdatedAt(),
	}
}

// ToUserResponseList converts a slice of User entities to a slice of UserResponse DTOs
func ToUserResponseList(users []entity.User) []UserResponse {
	result := make([]UserResponse, 0, len(users))
	for _, user := range users {
		result = append(result, ToUserResponse(user))
	}
	return result
}

// ToGroupResponse converts a Group entity to GroupResponse DTO
func ToGroupResponse(group entity.Group) GroupResponse {
	memberIDs := make([]string, len(group.MemberIDs()))
	for i, id := range group.MemberIDs() {
		memberIDs[i] = id.String()
	}

	return GroupResponse{
		ID:          group.ID().String(),
		Name:        group.Name(),
		Description: group.Description(),
		MemberIDs:   memberIDs,
		CreatedAt:   group.CreatedAt(),
		UpdatedAt:   group.UpdatedAt(),
	}
}

// ToGroupResponseList converts a slice of Group entities to a slice of GroupResponse DTOs
func ToGroupResponseList(groups []entity.Group) []GroupResponse {
	result := make([]GroupResponse, 0, len(groups))
	for _, group := range groups {
		result = append(result, ToGroupResponse(group))
	}
	return result
}
