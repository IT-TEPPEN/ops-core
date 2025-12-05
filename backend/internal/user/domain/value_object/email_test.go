package value_object

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewEmail(t *testing.T) {
	t.Run("有効なメールアドレスで正常に作成できる", func(t *testing.T) {
		testCases := []struct {
			input    string
			expected string
		}{
			{"test@example.com", "test@example.com"},
			{"USER@EXAMPLE.COM", "user@example.com"},
			{"test.name@example.co.jp", "test.name@example.co.jp"},
			{"test+tag@example.com", "test+tag@example.com"},
		}

		for _, tc := range testCases {
			email, err := NewEmail(tc.input)
			assert.NoError(t, err)
			assert.Equal(t, tc.expected, email.String())
		}
	})

	t.Run("空のメールアドレスでエラーになる", func(t *testing.T) {
		_, err := NewEmail("")

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "email cannot be empty")
	})

	t.Run("無効なメールアドレス形式でエラーになる", func(t *testing.T) {
		invalidEmails := []string{
			"invalid",
			"invalid@",
			"@example.com",
			"test@.com",
			"test@example.",
		}

		for _, email := range invalidEmails {
			_, err := NewEmail(email)
			assert.Error(t, err)
			assert.Contains(t, err.Error(), "invalid email format")
		}
	})
}

func TestEmail_String(t *testing.T) {
	email, _ := NewEmail("test@example.com")
	assert.Equal(t, "test@example.com", email.String())
}

func TestEmail_Equals(t *testing.T) {
	t.Run("同じメールアドレスは等しい", func(t *testing.T) {
		email1, _ := NewEmail("test@example.com")
		email2, _ := NewEmail("test@example.com")

		assert.True(t, email1.Equals(email2))
	})

	t.Run("大文字小文字は正規化されて等しくなる", func(t *testing.T) {
		email1, _ := NewEmail("test@example.com")
		email2, _ := NewEmail("TEST@EXAMPLE.COM")

		assert.True(t, email1.Equals(email2))
	})

	t.Run("異なるメールアドレスは等しくない", func(t *testing.T) {
		email1, _ := NewEmail("test1@example.com")
		email2, _ := NewEmail("test2@example.com")

		assert.False(t, email1.Equals(email2))
	})
}
