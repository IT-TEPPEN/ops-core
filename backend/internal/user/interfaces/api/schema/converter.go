package schema

import (
	"opscore/backend/internal/user/application/dto"
)

// ToCreateUserDTO converts API schema to application DTO
func ToCreateUserDTO(req CreateUserRequest) dto.CreateUserRequest {
	return dto.CreateUserRequest{
		Name:  req.Name,
		Email: req.Email,
		Role:  req.Role,
	}
}

// ToUpdateUserDTO converts API schema to application DTO
func ToUpdateUserDTO(req UpdateUserRequest) dto.UpdateUserRequest {
	return dto.UpdateUserRequest{
		Name:  req.Name,
		Email: req.Email,
	}
}

// ToChangeRoleDTO converts API schema to application DTO
func ToChangeRoleDTO(req ChangeRoleRequest) dto.ChangeRoleRequest {
	return dto.ChangeRoleRequest{
		Role: req.Role,
	}
}

// ToJoinGroupDTO converts API schema to application DTO
func ToJoinGroupDTO(req JoinGroupRequest) dto.JoinGroupRequest {
	return dto.JoinGroupRequest{
		GroupID: req.GroupID,
	}
}

// ToLeaveGroupDTO converts API schema to application DTO
func ToLeaveGroupDTO(req LeaveGroupRequest) dto.LeaveGroupRequest {
	return dto.LeaveGroupRequest{
		GroupID: req.GroupID,
	}
}

// FromUserDTO converts application DTO to API schema
func FromUserDTO(dtoResp dto.UserResponse) UserResponse {
	groupIDs := dtoResp.GroupIDs
	if groupIDs == nil {
		groupIDs = []string{}
	}
	return UserResponse{
		ID:        dtoResp.ID,
		Name:      dtoResp.Name,
		Email:     dtoResp.Email,
		Role:      dtoResp.Role,
		GroupIDs:  groupIDs,
		CreatedAt: dtoResp.CreatedAt,
		UpdatedAt: dtoResp.UpdatedAt,
	}
}

// FromUserListDTO converts application DTO list to API schema list
func FromUserListDTO(dtoList []dto.UserResponse) []UserResponse {
	schemas := make([]UserResponse, 0, len(dtoList))
	for _, dtoResp := range dtoList {
		schemas = append(schemas, FromUserDTO(dtoResp))
	}
	return schemas
}

// ToCreateGroupDTO converts API schema to application DTO
func ToCreateGroupDTO(req CreateGroupRequest) dto.CreateGroupRequest {
	return dto.CreateGroupRequest{
		Name:        req.Name,
		Description: req.Description,
	}
}

// ToUpdateGroupDTO converts API schema to application DTO
func ToUpdateGroupDTO(req UpdateGroupRequest) dto.UpdateGroupRequest {
	return dto.UpdateGroupRequest{
		Name:        req.Name,
		Description: req.Description,
	}
}

// ToAddMemberDTO converts API schema to application DTO
func ToAddMemberDTO(req AddMemberRequest) dto.AddMemberRequest {
	return dto.AddMemberRequest{
		UserID: req.UserID,
	}
}

// ToRemoveMemberDTO converts API schema to application DTO
func ToRemoveMemberDTO(req RemoveMemberRequest) dto.RemoveMemberRequest {
	return dto.RemoveMemberRequest{
		UserID: req.UserID,
	}
}

// FromGroupDTO converts application DTO to API schema
func FromGroupDTO(dtoResp dto.GroupResponse) GroupResponse {
	memberIDs := dtoResp.MemberIDs
	if memberIDs == nil {
		memberIDs = []string{}
	}
	return GroupResponse{
		ID:          dtoResp.ID,
		Name:        dtoResp.Name,
		Description: dtoResp.Description,
		MemberIDs:   memberIDs,
		CreatedAt:   dtoResp.CreatedAt,
		UpdatedAt:   dtoResp.UpdatedAt,
	}
}

// FromGroupListDTO converts application DTO list to API schema list
func FromGroupListDTO(dtoList []dto.GroupResponse) []GroupResponse {
	schemas := make([]GroupResponse, 0, len(dtoList))
	for _, dtoResp := range dtoList {
		schemas = append(schemas, FromGroupDTO(dtoResp))
	}
	return schemas
}
