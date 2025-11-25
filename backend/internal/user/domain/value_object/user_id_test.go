package value_object

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewUserID(t *testing.T) {
	t.Run("有効なIDで正常に作成できる", func(t *testing.T) {
		id := "user-123"
		userID, err := NewUserID(id)

		assert.NoError(t, err)
		assert.Equal(t, id, userID.String())
	})

	t.Run("空のIDでエラーになる", func(t *testing.T) {
		_, err := NewUserID("")

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "user ID cannot be empty")
	})
}

func TestUserID_String(t *testing.T) {
	userID, _ := NewUserID("user-456")
	assert.Equal(t, "user-456", userID.String())
}

func TestUserID_Equals(t *testing.T) {
	t.Run("同じIDは等しい", func(t *testing.T) {
		userID1, _ := NewUserID("user-123")
		userID2, _ := NewUserID("user-123")

		assert.True(t, userID1.Equals(userID2))
	})

	t.Run("異なるIDは等しくない", func(t *testing.T) {
		userID1, _ := NewUserID("user-123")
		userID2, _ := NewUserID("user-456")

		assert.False(t, userID1.Equals(userID2))
	})
}
