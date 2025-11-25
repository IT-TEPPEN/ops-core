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

func createTestUserForUseCase(t *testing.T, id, name, email, role string) entity.User {
	userID, _ := value_object.NewUserID(id)
	emailVO, _ := value_object.NewEmail(email)
	roleVO, _ := value_object.NewRole(role)
	user, err := entity.NewUser(userID, name, emailVO, roleVO)
	assert.NoError(t, err)
	return user
}

func TestUserUseCase_Create(t *testing.T) {
	t.Run("有効なリクエストでユーザーが正常に作成される", func(t *testing.T) {
		mockUserRepo := new(repository.MockUserRepository)
		mockGroupRepo := new(repository.MockGroupRepository)
		uc := NewUserUseCase(mockUserRepo, mockGroupRepo)

		email, _ := value_object.NewEmail("test@example.com")
		mockUserRepo.On("FindByEmail", mock.Anything, email).Return(nil, nil)
		mockUserRepo.On("Save", mock.Anything, mock.AnythingOfType("*entity.user")).Return(nil)

		req := dto.CreateUserRequest{
			Name:  "Test User",
			Email: "test@example.com",
			Role:  "user",
		}

		result, err := uc.Create(context.Background(), req)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, "Test User", result.Name)
		assert.Equal(t, "test@example.com", result.Email)
		assert.Equal(t, "user", result.Role)
		mockUserRepo.AssertExpectations(t)
	})

	t.Run("既に存在するメールでエラーになる", func(t *testing.T) {
		mockUserRepo := new(repository.MockUserRepository)
		mockGroupRepo := new(repository.MockGroupRepository)
		uc := NewUserUseCase(mockUserRepo, mockGroupRepo)

		existingUser := createTestUserForUseCase(t, "existing-user", "Existing User", "test@example.com", "user")
		email, _ := value_object.NewEmail("test@example.com")
		mockUserRepo.On("FindByEmail", mock.Anything, email).Return(existingUser, nil)

		req := dto.CreateUserRequest{
			Name:  "Test User",
			Email: "test@example.com",
			Role:  "user",
		}

		_, err := uc.Create(context.Background(), req)

		assert.Error(t, err)
		assert.True(t, errors.Is(err, apperror.ErrConflict))
	})

	t.Run("無効なメール形式でエラーになる", func(t *testing.T) {
		mockUserRepo := new(repository.MockUserRepository)
		mockGroupRepo := new(repository.MockGroupRepository)
		uc := NewUserUseCase(mockUserRepo, mockGroupRepo)

		req := dto.CreateUserRequest{
			Name:  "Test User",
			Email: "invalid-email",
			Role:  "user",
		}

		_, err := uc.Create(context.Background(), req)

		assert.Error(t, err)
		assert.True(t, errors.Is(err, apperror.ErrBadRequest))
	})
}

func TestUserUseCase_GetByID(t *testing.T) {
	t.Run("存在するIDでユーザーが取得できる", func(t *testing.T) {
		mockUserRepo := new(repository.MockUserRepository)
		mockGroupRepo := new(repository.MockGroupRepository)
		uc := NewUserUseCase(mockUserRepo, mockGroupRepo)

		user := createTestUserForUseCase(t, "user-123", "Test User", "test@example.com", "user")
		userID, _ := value_object.NewUserID("user-123")
		mockUserRepo.On("FindByID", mock.Anything, userID).Return(user, nil)

		result, err := uc.GetByID(context.Background(), "user-123")

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, "user-123", result.ID)
		assert.Equal(t, "Test User", result.Name)
		mockUserRepo.AssertExpectations(t)
	})

	t.Run("存在しないIDでエラーになる", func(t *testing.T) {
		mockUserRepo := new(repository.MockUserRepository)
		mockGroupRepo := new(repository.MockGroupRepository)
		uc := NewUserUseCase(mockUserRepo, mockGroupRepo)

		userID, _ := value_object.NewUserID("non-existent")
		mockUserRepo.On("FindByID", mock.Anything, userID).Return(nil, nil)

		_, err := uc.GetByID(context.Background(), "non-existent")

		assert.Error(t, err)
		assert.True(t, errors.Is(err, apperror.ErrNotFound))
	})
}

func TestUserUseCase_GetAll(t *testing.T) {
	t.Run("全ユーザーが取得できる", func(t *testing.T) {
		mockUserRepo := new(repository.MockUserRepository)
		mockGroupRepo := new(repository.MockGroupRepository)
		uc := NewUserUseCase(mockUserRepo, mockGroupRepo)

		user1 := createTestUserForUseCase(t, "user-1", "User 1", "user1@example.com", "user")
		user2 := createTestUserForUseCase(t, "user-2", "User 2", "user2@example.com", "admin")
		mockUserRepo.On("FindAll", mock.Anything).Return([]entity.User{user1, user2}, nil)

		result, err := uc.GetAll(context.Background())

		assert.NoError(t, err)
		assert.Len(t, result, 2)
		mockUserRepo.AssertExpectations(t)
	})
}

func TestUserUseCase_Update(t *testing.T) {
	t.Run("有効なリクエストでユーザーが更新される", func(t *testing.T) {
		mockUserRepo := new(repository.MockUserRepository)
		mockGroupRepo := new(repository.MockGroupRepository)
		uc := NewUserUseCase(mockUserRepo, mockGroupRepo)

		user := createTestUserForUseCase(t, "user-123", "Old Name", "old@example.com", "user")
		userID, _ := value_object.NewUserID("user-123")
		newEmail, _ := value_object.NewEmail("new@example.com")

		mockUserRepo.On("FindByID", mock.Anything, userID).Return(user, nil)
		mockUserRepo.On("FindByEmail", mock.Anything, newEmail).Return(nil, nil)
		mockUserRepo.On("Update", mock.Anything, mock.AnythingOfType("*entity.user")).Return(nil)

		req := dto.UpdateUserRequest{
			Name:  "New Name",
			Email: "new@example.com",
		}

		result, err := uc.Update(context.Background(), "user-123", req)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, "New Name", result.Name)
		mockUserRepo.AssertExpectations(t)
	})
}

func TestUserUseCase_Delete(t *testing.T) {
	t.Run("存在するIDでユーザーが削除できる", func(t *testing.T) {
		mockUserRepo := new(repository.MockUserRepository)
		mockGroupRepo := new(repository.MockGroupRepository)
		uc := NewUserUseCase(mockUserRepo, mockGroupRepo)

		user := createTestUserForUseCase(t, "user-123", "Test User", "test@example.com", "user")
		userID, _ := value_object.NewUserID("user-123")

		mockUserRepo.On("FindByID", mock.Anything, userID).Return(user, nil)
		mockUserRepo.On("Delete", mock.Anything, userID).Return(nil)

		err := uc.Delete(context.Background(), "user-123")

		assert.NoError(t, err)
		mockUserRepo.AssertExpectations(t)
	})

	t.Run("存在しないIDでエラーになる", func(t *testing.T) {
		mockUserRepo := new(repository.MockUserRepository)
		mockGroupRepo := new(repository.MockGroupRepository)
		uc := NewUserUseCase(mockUserRepo, mockGroupRepo)

		userID, _ := value_object.NewUserID("non-existent")
		mockUserRepo.On("FindByID", mock.Anything, userID).Return(nil, nil)

		err := uc.Delete(context.Background(), "non-existent")

		assert.Error(t, err)
		assert.True(t, errors.Is(err, apperror.ErrNotFound))
	})
}

func TestUserUseCase_ChangeRole(t *testing.T) {
	t.Run("ロールが正常に変更される", func(t *testing.T) {
		mockUserRepo := new(repository.MockUserRepository)
		mockGroupRepo := new(repository.MockGroupRepository)
		uc := NewUserUseCase(mockUserRepo, mockGroupRepo)

		user := createTestUserForUseCase(t, "user-123", "Test User", "test@example.com", "user")
		userID, _ := value_object.NewUserID("user-123")

		mockUserRepo.On("FindByID", mock.Anything, userID).Return(user, nil)
		mockUserRepo.On("Update", mock.Anything, mock.AnythingOfType("*entity.user")).Return(nil)

		req := dto.ChangeRoleRequest{
			Role: "admin",
		}

		result, err := uc.ChangeRole(context.Background(), "user-123", req)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, "admin", result.Role)
		mockUserRepo.AssertExpectations(t)
	})
}

func TestUserUseCase_JoinGroup(t *testing.T) {
	t.Run("ユーザーがグループに参加できる", func(t *testing.T) {
		mockUserRepo := new(repository.MockUserRepository)
		mockGroupRepo := new(repository.MockGroupRepository)
		uc := NewUserUseCase(mockUserRepo, mockGroupRepo)

		user := createTestUserForUseCase(t, "user-123", "Test User", "test@example.com", "user")
		userID, _ := value_object.NewUserID("user-123")

		groupID, _ := value_object.NewGroupID("group-123")
		group, _ := entity.NewGroup(groupID, "Test Group", "Description")

		mockUserRepo.On("FindByID", mock.Anything, userID).Return(user, nil)
		mockGroupRepo.On("FindByID", mock.Anything, groupID).Return(group, nil)
		mockUserRepo.On("Update", mock.Anything, mock.AnythingOfType("*entity.user")).Return(nil)
		mockGroupRepo.On("Update", mock.Anything, mock.AnythingOfType("*entity.group")).Return(nil)

		req := dto.JoinGroupRequest{
			GroupID: "group-123",
		}

		result, err := uc.JoinGroup(context.Background(), "user-123", req)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Contains(t, result.GroupIDs, "group-123")
		mockUserRepo.AssertExpectations(t)
		mockGroupRepo.AssertExpectations(t)
	})
}

func TestUserUseCase_LeaveGroup(t *testing.T) {
	t.Run("ユーザーがグループから脱退できる", func(t *testing.T) {
		mockUserRepo := new(repository.MockUserRepository)
		mockGroupRepo := new(repository.MockGroupRepository)
		uc := NewUserUseCase(mockUserRepo, mockGroupRepo)

		// Create user with group membership
		userID, _ := value_object.NewUserID("user-123")
		email, _ := value_object.NewEmail("test@example.com")
		role, _ := value_object.NewRole("user")
		groupID, _ := value_object.NewGroupID("group-123")
		user := entity.ReconstructUser(userID, "Test User", email, role, []value_object.GroupID{groupID}, time.Now(), time.Now())

		group, _ := entity.NewGroup(groupID, "Test Group", "Description")
		_ = group.AddMember(userID)

		mockUserRepo.On("FindByID", mock.Anything, userID).Return(user, nil)
		mockGroupRepo.On("FindByID", mock.Anything, groupID).Return(group, nil)
		mockUserRepo.On("Update", mock.Anything, mock.AnythingOfType("*entity.user")).Return(nil)
		mockGroupRepo.On("Update", mock.Anything, mock.AnythingOfType("*entity.group")).Return(nil)

		req := dto.LeaveGroupRequest{
			GroupID: "group-123",
		}

		result, err := uc.LeaveGroup(context.Background(), "user-123", req)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Empty(t, result.GroupIDs)
		mockUserRepo.AssertExpectations(t)
		mockGroupRepo.AssertExpectations(t)
	})
}
