package dto

import (
	"testing"
	"time"

	"opscore/backend/internal/user/domain/entity"
	"opscore/backend/internal/user/domain/value_object"

	"github.com/stretchr/testify/assert"
)

func TestToUserResponse(t *testing.T) {
	userID, _ := value_object.NewUserID("user-123")
	email, _ := value_object.NewEmail("test@example.com")
	role, _ := value_object.NewRole("admin")
	groupID, _ := value_object.NewGroupID("group-1")
	createdAt := time.Date(2023, 1, 1, 12, 0, 0, 0, time.UTC)
	updatedAt := time.Date(2023, 1, 2, 12, 0, 0, 0, time.UTC)

	user := entity.ReconstructUser(userID, "Test User", email, role, []value_object.GroupID{groupID}, createdAt, updatedAt)

	response := ToUserResponse(user)

	assert.Equal(t, "user-123", response.ID)
	assert.Equal(t, "Test User", response.Name)
	assert.Equal(t, "test@example.com", response.Email)
	assert.Equal(t, "admin", response.Role)
	assert.Len(t, response.GroupIDs, 1)
	assert.Equal(t, "group-1", response.GroupIDs[0])
	assert.Equal(t, createdAt, response.CreatedAt)
	assert.Equal(t, updatedAt, response.UpdatedAt)
}

func TestToUserResponseList(t *testing.T) {
	userID1, _ := value_object.NewUserID("user-1")
	email1, _ := value_object.NewEmail("user1@example.com")
	role1, _ := value_object.NewRole("user")
	createdAt := time.Now()
	updatedAt := time.Now()

	userID2, _ := value_object.NewUserID("user-2")
	email2, _ := value_object.NewEmail("user2@example.com")
	role2, _ := value_object.NewRole("admin")

	user1 := entity.ReconstructUser(userID1, "User 1", email1, role1, nil, createdAt, updatedAt)
	user2 := entity.ReconstructUser(userID2, "User 2", email2, role2, nil, createdAt, updatedAt)

	responses := ToUserResponseList([]entity.User{user1, user2})

	assert.Len(t, responses, 2)
	assert.Equal(t, "user-1", responses[0].ID)
	assert.Equal(t, "user-2", responses[1].ID)
}

func TestToUserResponseList_EmptyList(t *testing.T) {
	responses := ToUserResponseList([]entity.User{})
	assert.Empty(t, responses)
}

func TestToGroupResponse(t *testing.T) {
	groupID, _ := value_object.NewGroupID("group-123")
	userID, _ := value_object.NewUserID("user-1")
	createdAt := time.Date(2023, 1, 1, 12, 0, 0, 0, time.UTC)
	updatedAt := time.Date(2023, 1, 2, 12, 0, 0, 0, time.UTC)

	group := entity.ReconstructGroup(groupID, "Test Group", "Description", []value_object.UserID{userID}, createdAt, updatedAt)

	response := ToGroupResponse(group)

	assert.Equal(t, "group-123", response.ID)
	assert.Equal(t, "Test Group", response.Name)
	assert.Equal(t, "Description", response.Description)
	assert.Len(t, response.MemberIDs, 1)
	assert.Equal(t, "user-1", response.MemberIDs[0])
	assert.Equal(t, createdAt, response.CreatedAt)
	assert.Equal(t, updatedAt, response.UpdatedAt)
}

func TestToGroupResponseList(t *testing.T) {
	groupID1, _ := value_object.NewGroupID("group-1")
	groupID2, _ := value_object.NewGroupID("group-2")
	createdAt := time.Now()
	updatedAt := time.Now()

	group1 := entity.ReconstructGroup(groupID1, "Group 1", "Desc 1", nil, createdAt, updatedAt)
	group2 := entity.ReconstructGroup(groupID2, "Group 2", "Desc 2", nil, createdAt, updatedAt)

	responses := ToGroupResponseList([]entity.Group{group1, group2})

	assert.Len(t, responses, 2)
	assert.Equal(t, "group-1", responses[0].ID)
	assert.Equal(t, "group-2", responses[1].ID)
}

func TestToGroupResponseList_EmptyList(t *testing.T) {
	responses := ToGroupResponseList([]entity.Group{})
	assert.Empty(t, responses)
}
