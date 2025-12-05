package value_object

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewRole(t *testing.T) {
	t.Run("有効なロールで正常に作成できる", func(t *testing.T) {
		testCases := []struct {
			input    string
			expected string
		}{
			{"admin", "admin"},
			{"ADMIN", "admin"},
			{"user", "user"},
			{"USER", "user"},
			{"  admin  ", "admin"},
		}

		for _, tc := range testCases {
			role, err := NewRole(tc.input)
			assert.NoError(t, err)
			assert.Equal(t, tc.expected, role.String())
		}
	})

	t.Run("無効なロールでエラーになる", func(t *testing.T) {
		invalidRoles := []string{
			"",
			"superadmin",
			"moderator",
			"guest",
		}

		for _, role := range invalidRoles {
			_, err := NewRole(role)
			assert.Error(t, err)
			assert.Contains(t, err.Error(), "role must be either 'admin' or 'user'")
		}
	})
}

func TestRole_String(t *testing.T) {
	role, _ := NewRole("admin")
	assert.Equal(t, "admin", role.String())
}

func TestRole_IsAdmin(t *testing.T) {
	t.Run("adminロールはtrueを返す", func(t *testing.T) {
		role, _ := NewRole("admin")
		assert.True(t, role.IsAdmin())
	})

	t.Run("userロールはfalseを返す", func(t *testing.T) {
		role, _ := NewRole("user")
		assert.False(t, role.IsAdmin())
	})
}

func TestRole_Equals(t *testing.T) {
	t.Run("同じロールは等しい", func(t *testing.T) {
		role1, _ := NewRole("admin")
		role2, _ := NewRole("admin")

		assert.True(t, role1.Equals(role2))
	})

	t.Run("異なるロールは等しくない", func(t *testing.T) {
		role1, _ := NewRole("admin")
		role2, _ := NewRole("user")

		assert.False(t, role1.Equals(role2))
	})
}
