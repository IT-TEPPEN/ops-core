package usecase

import (
	"context"
	"errors"
	"testing"
	"time"

	apperror "opscore/backend/internal/user/application/error"
	"opscore/backend/internal/user/application/dto"
	"opscore/backend/internal/user/domain/entity"
	"opscore/backend/internal/user/domain/repository"
	"opscore/backend/internal/user/domain/value_object"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func createTestGroupForUseCase(t *testing.T, id, name, description string) entity.Group {
	groupID, _ := value_object.NewGroupID(id)
	group, err := entity.NewGroup(groupID, name, description)
	assert.NoError(t, err)
	return group
}

func TestGroupUseCase_Create(t *testing.T) {
	t.Run("有効なリクエストでグループが正常に作成される", func(t *testing.T) {
		mockGroupRepo := new(repository.MockGroupRepository)
		mockUserRepo := new(repository.MockUserRepository)
		uc := NewGroupUseCase(mockGroupRepo, mockUserRepo)

		mockGroupRepo.On("Save", mock.Anything, mock.AnythingOfType("*entity.group")).Return(nil)

		req := dto.CreateGroupRequest{
			Name:        "Test Group",
			Description: "Test Description",
		}

		result, err := uc.Create(context.Background(), req)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, "Test Group", result.Name)
		assert.Equal(t, "Test Description", result.Description)
		mockGroupRepo.AssertExpectations(t)
	})

	t.Run("空の名前でエラーになる", func(t *testing.T) {
		mockGroupRepo := new(repository.MockGroupRepository)
		mockUserRepo := new(repository.MockUserRepository)
		uc := NewGroupUseCase(mockGroupRepo, mockUserRepo)

		req := dto.CreateGroupRequest{
			Name:        "",
			Description: "Test Description",
		}

		_, err := uc.Create(context.Background(), req)

		assert.Error(t, err)
		assert.True(t, errors.Is(err, apperror.ErrBadRequest))
	})
}

func TestGroupUseCase_GetByID(t *testing.T) {
	t.Run("存在するIDでグループが取得できる", func(t *testing.T) {
		mockGroupRepo := new(repository.MockGroupRepository)
		mockUserRepo := new(repository.MockUserRepository)
		uc := NewGroupUseCase(mockGroupRepo, mockUserRepo)

		group := createTestGroupForUseCase(t, "group-123", "Test Group", "Description")
		groupID, _ := value_object.NewGroupID("group-123")
		mockGroupRepo.On("FindByID", mock.Anything, groupID).Return(group, nil)

		result, err := uc.GetByID(context.Background(), "group-123")

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, "group-123", result.ID)
		assert.Equal(t, "Test Group", result.Name)
		mockGroupRepo.AssertExpectations(t)
	})

	t.Run("存在しないIDでエラーになる", func(t *testing.T) {
		mockGroupRepo := new(repository.MockGroupRepository)
		mockUserRepo := new(repository.MockUserRepository)
		uc := NewGroupUseCase(mockGroupRepo, mockUserRepo)

		groupID, _ := value_object.NewGroupID("non-existent")
		mockGroupRepo.On("FindByID", mock.Anything, groupID).Return(nil, nil)

		_, err := uc.GetByID(context.Background(), "non-existent")

		assert.Error(t, err)
		assert.True(t, errors.Is(err, apperror.ErrNotFound))
	})
}

func TestGroupUseCase_GetAll(t *testing.T) {
	t.Run("全グループが取得できる", func(t *testing.T) {
		mockGroupRepo := new(repository.MockGroupRepository)
		mockUserRepo := new(repository.MockUserRepository)
		uc := NewGroupUseCase(mockGroupRepo, mockUserRepo)

		group1 := createTestGroupForUseCase(t, "group-1", "Group 1", "Description 1")
		group2 := createTestGroupForUseCase(t, "group-2", "Group 2", "Description 2")
		mockGroupRepo.On("FindAll", mock.Anything).Return([]entity.Group{group1, group2}, nil)

		result, err := uc.GetAll(context.Background())

		assert.NoError(t, err)
		assert.Len(t, result, 2)
		mockGroupRepo.AssertExpectations(t)
	})
}

func TestGroupUseCase_Update(t *testing.T) {
	t.Run("有効なリクエストでグループが更新される", func(t *testing.T) {
		mockGroupRepo := new(repository.MockGroupRepository)
		mockUserRepo := new(repository.MockUserRepository)
		uc := NewGroupUseCase(mockGroupRepo, mockUserRepo)

		group := createTestGroupForUseCase(t, "group-123", "Old Name", "Old Description")
		groupID, _ := value_object.NewGroupID("group-123")

		mockGroupRepo.On("FindByID", mock.Anything, groupID).Return(group, nil)
		mockGroupRepo.On("Update", mock.Anything, mock.AnythingOfType("*entity.group")).Return(nil)

		req := dto.UpdateGroupRequest{
			Name:        "New Name",
			Description: "New Description",
		}

		result, err := uc.Update(context.Background(), "group-123", req)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, "New Name", result.Name)
		assert.Equal(t, "New Description", result.Description)
		mockGroupRepo.AssertExpectations(t)
	})
}

func TestGroupUseCase_Delete(t *testing.T) {
	t.Run("存在するIDでグループが削除できる", func(t *testing.T) {
		mockGroupRepo := new(repository.MockGroupRepository)
		mockUserRepo := new(repository.MockUserRepository)
		uc := NewGroupUseCase(mockGroupRepo, mockUserRepo)

		group := createTestGroupForUseCase(t, "group-123", "Test Group", "Description")
		groupID, _ := value_object.NewGroupID("group-123")

		mockGroupRepo.On("FindByID", mock.Anything, groupID).Return(group, nil)
		mockGroupRepo.On("Delete", mock.Anything, groupID).Return(nil)

		err := uc.Delete(context.Background(), "group-123")

		assert.NoError(t, err)
		mockGroupRepo.AssertExpectations(t)
	})

	t.Run("存在しないIDでエラーになる", func(t *testing.T) {
		mockGroupRepo := new(repository.MockGroupRepository)
		mockUserRepo := new(repository.MockUserRepository)
		uc := NewGroupUseCase(mockGroupRepo, mockUserRepo)

		groupID, _ := value_object.NewGroupID("non-existent")
		mockGroupRepo.On("FindByID", mock.Anything, groupID).Return(nil, nil)

		err := uc.Delete(context.Background(), "non-existent")

		assert.Error(t, err)
		assert.True(t, errors.Is(err, apperror.ErrNotFound))
	})
}

func TestGroupUseCase_AddMember(t *testing.T) {
	t.Run("グループにメンバーを追加できる", func(t *testing.T) {
		mockGroupRepo := new(repository.MockGroupRepository)
		mockUserRepo := new(repository.MockUserRepository)
		uc := NewGroupUseCase(mockGroupRepo, mockUserRepo)

		group := createTestGroupForUseCase(t, "group-123", "Test Group", "Description")
		groupID, _ := value_object.NewGroupID("group-123")

		user := createTestUserForUseCase(t, "user-123", "Test User", "test@example.com", "user")
		userID, _ := value_object.NewUserID("user-123")

		mockGroupRepo.On("FindByID", mock.Anything, groupID).Return(group, nil)
		mockUserRepo.On("FindByID", mock.Anything, userID).Return(user, nil)
		mockGroupRepo.On("Update", mock.Anything, mock.AnythingOfType("*entity.group")).Return(nil)
		mockUserRepo.On("Update", mock.Anything, mock.AnythingOfType("*entity.user")).Return(nil)

		req := dto.AddMemberRequest{
			UserID: "user-123",
		}

		result, err := uc.AddMember(context.Background(), "group-123", req)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Contains(t, result.MemberIDs, "user-123")
		mockGroupRepo.AssertExpectations(t)
		mockUserRepo.AssertExpectations(t)
	})
}

func TestGroupUseCase_RemoveMember(t *testing.T) {
	t.Run("グループからメンバーを削除できる", func(t *testing.T) {
		mockGroupRepo := new(repository.MockGroupRepository)
		mockUserRepo := new(repository.MockUserRepository)
		uc := NewGroupUseCase(mockGroupRepo, mockUserRepo)

		// Create group with member
		groupID, _ := value_object.NewGroupID("group-123")
		userID, _ := value_object.NewUserID("user-123")
		group := entity.ReconstructGroup(groupID, "Test Group", "Description", []value_object.UserID{userID}, time.Now(), time.Now())

		email, _ := value_object.NewEmail("test@example.com")
		role, _ := value_object.NewRole("user")
		user := entity.ReconstructUser(userID, "Test User", email, role, []value_object.GroupID{groupID}, time.Now(), time.Now())

		mockGroupRepo.On("FindByID", mock.Anything, groupID).Return(group, nil)
		mockUserRepo.On("FindByID", mock.Anything, userID).Return(user, nil)
		mockGroupRepo.On("Update", mock.Anything, mock.AnythingOfType("*entity.group")).Return(nil)
		mockUserRepo.On("Update", mock.Anything, mock.AnythingOfType("*entity.user")).Return(nil)

		req := dto.RemoveMemberRequest{
			UserID: "user-123",
		}

		result, err := uc.RemoveMember(context.Background(), "group-123", req)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Empty(t, result.MemberIDs)
		mockGroupRepo.AssertExpectations(t)
		mockUserRepo.AssertExpectations(t)
	})
}

func TestGroupUseCase_GetByMemberID(t *testing.T) {
	t.Run("ユーザーが所属するグループが取得できる", func(t *testing.T) {
		mockGroupRepo := new(repository.MockGroupRepository)
		mockUserRepo := new(repository.MockUserRepository)
		uc := NewGroupUseCase(mockGroupRepo, mockUserRepo)

		user := createTestUserForUseCase(t, "user-123", "Test User", "test@example.com", "user")
		userID, _ := value_object.NewUserID("user-123")

		group1 := createTestGroupForUseCase(t, "group-1", "Group 1", "Description 1")
		group2 := createTestGroupForUseCase(t, "group-2", "Group 2", "Description 2")

		mockUserRepo.On("FindByID", mock.Anything, userID).Return(user, nil)
		mockGroupRepo.On("FindByMemberID", mock.Anything, userID).Return([]entity.Group{group1, group2}, nil)

		result, err := uc.GetByMemberID(context.Background(), "user-123")

		assert.NoError(t, err)
		assert.Len(t, result, 2)
		mockUserRepo.AssertExpectations(t)
		mockGroupRepo.AssertExpectations(t)
	})
}
