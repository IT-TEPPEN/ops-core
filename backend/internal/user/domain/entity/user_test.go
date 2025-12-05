package entity

import (
	"testing"
	"time"

	"opscore/backend/internal/user/domain/value_object"

	"github.com/stretchr/testify/assert"
)

func createTestUser(t *testing.T) User {
	userID, _ := value_object.NewUserID("user-123")
	email, _ := value_object.NewEmail("test@example.com")
	role, _ := value_object.NewRole("user")

	user, err := NewUser(userID, "Test User", email, role)
	assert.NoError(t, err)
	return user
}

func TestNewUser(t *testing.T) {
	t.Run("有効なパラメータでユーザーが正常に作成される", func(t *testing.T) {
		userID, _ := value_object.NewUserID("user-123")
		email, _ := value_object.NewEmail("test@example.com")
		role, _ := value_object.NewRole("admin")

		beforeCreate := time.Now()
		user, err := NewUser(userID, "Test User", email, role)
		afterCreate := time.Now()

		assert.NoError(t, err)
		assert.NotNil(t, user)
		assert.Equal(t, "user-123", user.ID().String())
		assert.Equal(t, "Test User", user.Name())
		assert.Equal(t, "test@example.com", user.Email().String())
		assert.Equal(t, "admin", user.Role().String())
		assert.Empty(t, user.GroupIDs())

		// 作成時間が現在時刻に近いことを確認
		assert.True(t, user.CreatedAt().After(beforeCreate) || user.CreatedAt().Equal(beforeCreate))
		assert.True(t, user.CreatedAt().Before(afterCreate) || user.CreatedAt().Equal(afterCreate))
	})

	t.Run("空の名前でエラーになる", func(t *testing.T) {
		userID, _ := value_object.NewUserID("user-123")
		email, _ := value_object.NewEmail("test@example.com")
		role, _ := value_object.NewRole("user")

		_, err := NewUser(userID, "", email, role)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "user name cannot be empty")
	})
}

func TestReconstructUser(t *testing.T) {
	t.Run("永続化データからユーザーが正しく再構築される", func(t *testing.T) {
		userID, _ := value_object.NewUserID("user-123")
		email, _ := value_object.NewEmail("test@example.com")
		role, _ := value_object.NewRole("admin")
		groupID, _ := value_object.NewGroupID("group-1")
		createdAt := time.Date(2023, 1, 1, 12, 0, 0, 0, time.UTC)
		updatedAt := time.Date(2023, 1, 2, 12, 0, 0, 0, time.UTC)

		user := ReconstructUser(userID, "Test User", email, role, []value_object.GroupID{groupID}, createdAt, updatedAt)

		assert.NotNil(t, user)
		assert.Equal(t, "user-123", user.ID().String())
		assert.Equal(t, "Test User", user.Name())
		assert.Equal(t, "test@example.com", user.Email().String())
		assert.Equal(t, "admin", user.Role().String())
		assert.Len(t, user.GroupIDs(), 1)
		assert.Equal(t, "group-1", user.GroupIDs()[0].String())
		assert.Equal(t, createdAt, user.CreatedAt())
		assert.Equal(t, updatedAt, user.UpdatedAt())
	})

	t.Run("nilのgroupIDsで空のスライスが設定される", func(t *testing.T) {
		userID, _ := value_object.NewUserID("user-123")
		email, _ := value_object.NewEmail("test@example.com")
		role, _ := value_object.NewRole("user")
		createdAt := time.Now()
		updatedAt := time.Now()

		user := ReconstructUser(userID, "Test User", email, role, nil, createdAt, updatedAt)

		assert.NotNil(t, user.GroupIDs())
		assert.Empty(t, user.GroupIDs())
	})
}

func TestUser_UpdateProfile(t *testing.T) {
	t.Run("プロフィールが正常に更新される", func(t *testing.T) {
		user := createTestUser(t)
		originalUpdatedAt := user.UpdatedAt()

		time.Sleep(time.Millisecond)
		newEmail, _ := value_object.NewEmail("new@example.com")
		err := user.UpdateProfile("New Name", newEmail)

		assert.NoError(t, err)
		assert.Equal(t, "New Name", user.Name())
		assert.Equal(t, "new@example.com", user.Email().String())
		assert.True(t, user.UpdatedAt().After(originalUpdatedAt))
	})

	t.Run("空の名前でエラーになる", func(t *testing.T) {
		user := createTestUser(t)
		newEmail, _ := value_object.NewEmail("new@example.com")

		err := user.UpdateProfile("", newEmail)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "user name cannot be empty")
	})
}

func TestUser_JoinGroup(t *testing.T) {
	t.Run("グループに正常に参加できる", func(t *testing.T) {
		user := createTestUser(t)
		groupID, _ := value_object.NewGroupID("group-1")
		originalUpdatedAt := user.UpdatedAt()

		time.Sleep(time.Millisecond)
		err := user.JoinGroup(groupID)

		assert.NoError(t, err)
		assert.Len(t, user.GroupIDs(), 1)
		assert.Equal(t, "group-1", user.GroupIDs()[0].String())
		assert.True(t, user.UpdatedAt().After(originalUpdatedAt))
	})

	t.Run("同じグループに参加しようとするとエラーになる", func(t *testing.T) {
		user := createTestUser(t)
		groupID, _ := value_object.NewGroupID("group-1")

		_ = user.JoinGroup(groupID)
		err := user.JoinGroup(groupID)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "already a member")
	})
}

func TestUser_LeaveGroup(t *testing.T) {
	t.Run("グループから正常に脱退できる", func(t *testing.T) {
		user := createTestUser(t)
		groupID, _ := value_object.NewGroupID("group-1")
		_ = user.JoinGroup(groupID)
		originalUpdatedAt := user.UpdatedAt()

		time.Sleep(time.Millisecond)
		err := user.LeaveGroup(groupID)

		assert.NoError(t, err)
		assert.Empty(t, user.GroupIDs())
		assert.True(t, user.UpdatedAt().After(originalUpdatedAt))
	})

	t.Run("参加していないグループから脱退しようとするとエラーになる", func(t *testing.T) {
		user := createTestUser(t)
		groupID, _ := value_object.NewGroupID("group-1")

		err := user.LeaveGroup(groupID)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "not a member")
	})
}

func TestUser_ChangeRole(t *testing.T) {
	t.Run("ロールが正常に変更される", func(t *testing.T) {
		user := createTestUser(t)
		adminRole, _ := value_object.NewRole("admin")
		originalUpdatedAt := user.UpdatedAt()

		time.Sleep(time.Millisecond)
		err := user.ChangeRole(adminRole)

		assert.NoError(t, err)
		assert.Equal(t, "admin", user.Role().String())
		assert.True(t, user.UpdatedAt().After(originalUpdatedAt))
	})
}

func TestUser_GroupIDs_ReturnsCopy(t *testing.T) {
	t.Run("GroupIDsが返すスライスはコピーである", func(t *testing.T) {
		user := createTestUser(t)
		groupID, _ := value_object.NewGroupID("group-1")
		_ = user.JoinGroup(groupID)

		groupIDs := user.GroupIDs()
		groupIDs[0] = value_object.GroupID{}

		// 元のユーザーのGroupIDsは変更されていない
		assert.Equal(t, "group-1", user.GroupIDs()[0].String())
	})
}
