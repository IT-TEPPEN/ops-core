package entity

import (
	"testing"
	"time"

	"opscore/backend/internal/user/domain/value_object"

	"github.com/stretchr/testify/assert"
)

func createTestGroup(t *testing.T) Group {
	groupID, _ := value_object.NewGroupID("group-123")

	group, err := NewGroup(groupID, "Test Group", "Test Description")
	assert.NoError(t, err)
	return group
}

func TestNewGroup(t *testing.T) {
	t.Run("有効なパラメータでグループが正常に作成される", func(t *testing.T) {
		groupID, _ := value_object.NewGroupID("group-123")

		beforeCreate := time.Now()
		group, err := NewGroup(groupID, "Test Group", "Test Description")
		afterCreate := time.Now()

		assert.NoError(t, err)
		assert.NotNil(t, group)
		assert.Equal(t, "group-123", group.ID().String())
		assert.Equal(t, "Test Group", group.Name())
		assert.Equal(t, "Test Description", group.Description())
		assert.Empty(t, group.MemberIDs())

		// 作成時間が現在時刻に近いことを確認
		assert.True(t, group.CreatedAt().After(beforeCreate) || group.CreatedAt().Equal(beforeCreate))
		assert.True(t, group.CreatedAt().Before(afterCreate) || group.CreatedAt().Equal(afterCreate))
	})

	t.Run("空の名前でエラーになる", func(t *testing.T) {
		groupID, _ := value_object.NewGroupID("group-123")

		_, err := NewGroup(groupID, "", "Description")

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "group name cannot be empty")
	})

	t.Run("空の説明でも正常に作成できる", func(t *testing.T) {
		groupID, _ := value_object.NewGroupID("group-123")

		group, err := NewGroup(groupID, "Test Group", "")

		assert.NoError(t, err)
		assert.Empty(t, group.Description())
	})
}

func TestReconstructGroup(t *testing.T) {
	t.Run("永続化データからグループが正しく再構築される", func(t *testing.T) {
		groupID, _ := value_object.NewGroupID("group-123")
		userID, _ := value_object.NewUserID("user-1")
		createdAt := time.Date(2023, 1, 1, 12, 0, 0, 0, time.UTC)
		updatedAt := time.Date(2023, 1, 2, 12, 0, 0, 0, time.UTC)

		group := ReconstructGroup(groupID, "Test Group", "Description", []value_object.UserID{userID}, createdAt, updatedAt)

		assert.NotNil(t, group)
		assert.Equal(t, "group-123", group.ID().String())
		assert.Equal(t, "Test Group", group.Name())
		assert.Equal(t, "Description", group.Description())
		assert.Len(t, group.MemberIDs(), 1)
		assert.Equal(t, "user-1", group.MemberIDs()[0].String())
		assert.Equal(t, createdAt, group.CreatedAt())
		assert.Equal(t, updatedAt, group.UpdatedAt())
	})

	t.Run("nilのmemberIDsで空のスライスが設定される", func(t *testing.T) {
		groupID, _ := value_object.NewGroupID("group-123")
		createdAt := time.Now()
		updatedAt := time.Now()

		group := ReconstructGroup(groupID, "Test Group", "Description", nil, createdAt, updatedAt)

		assert.NotNil(t, group.MemberIDs())
		assert.Empty(t, group.MemberIDs())
	})
}

func TestGroup_UpdateInfo(t *testing.T) {
	t.Run("グループ情報が正常に更新される", func(t *testing.T) {
		group := createTestGroup(t)
		originalUpdatedAt := group.UpdatedAt()

		time.Sleep(time.Millisecond)
		err := group.UpdateInfo("New Name", "New Description")

		assert.NoError(t, err)
		assert.Equal(t, "New Name", group.Name())
		assert.Equal(t, "New Description", group.Description())
		assert.True(t, group.UpdatedAt().After(originalUpdatedAt))
	})

	t.Run("空の名前でエラーになる", func(t *testing.T) {
		group := createTestGroup(t)

		err := group.UpdateInfo("", "New Description")

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "group name cannot be empty")
	})
}

func TestGroup_AddMember(t *testing.T) {
	t.Run("メンバーが正常に追加される", func(t *testing.T) {
		group := createTestGroup(t)
		userID, _ := value_object.NewUserID("user-1")
		originalUpdatedAt := group.UpdatedAt()

		time.Sleep(time.Millisecond)
		err := group.AddMember(userID)

		assert.NoError(t, err)
		assert.Len(t, group.MemberIDs(), 1)
		assert.Equal(t, "user-1", group.MemberIDs()[0].String())
		assert.True(t, group.UpdatedAt().After(originalUpdatedAt))
	})

	t.Run("同じメンバーを追加しようとするとエラーになる", func(t *testing.T) {
		group := createTestGroup(t)
		userID, _ := value_object.NewUserID("user-1")

		_ = group.AddMember(userID)
		err := group.AddMember(userID)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "already a member")
	})
}

func TestGroup_RemoveMember(t *testing.T) {
	t.Run("メンバーが正常に削除される", func(t *testing.T) {
		group := createTestGroup(t)
		userID, _ := value_object.NewUserID("user-1")
		_ = group.AddMember(userID)
		originalUpdatedAt := group.UpdatedAt()

		time.Sleep(time.Millisecond)
		err := group.RemoveMember(userID)

		assert.NoError(t, err)
		assert.Empty(t, group.MemberIDs())
		assert.True(t, group.UpdatedAt().After(originalUpdatedAt))
	})

	t.Run("存在しないメンバーを削除しようとするとエラーになる", func(t *testing.T) {
		group := createTestGroup(t)
		userID, _ := value_object.NewUserID("user-1")

		err := group.RemoveMember(userID)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "not a member")
	})
}

func TestGroup_MemberIDs_ReturnsCopy(t *testing.T) {
	t.Run("MemberIDsが返すスライスはコピーである", func(t *testing.T) {
		group := createTestGroup(t)
		userID, _ := value_object.NewUserID("user-1")
		_ = group.AddMember(userID)

		memberIDs := group.MemberIDs()
		memberIDs[0] = value_object.UserID{}

		// 元のグループのMemberIDsは変更されていない
		assert.Equal(t, "user-1", group.MemberIDs()[0].String())
	})
}
