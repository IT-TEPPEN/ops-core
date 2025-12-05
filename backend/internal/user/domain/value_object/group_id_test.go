package value_object

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewGroupID(t *testing.T) {
	t.Run("有効なIDで正常に作成できる", func(t *testing.T) {
		id := "group-123"
		groupID, err := NewGroupID(id)

		assert.NoError(t, err)
		assert.Equal(t, id, groupID.String())
	})

	t.Run("空のIDでエラーになる", func(t *testing.T) {
		_, err := NewGroupID("")

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "group ID cannot be empty")
	})
}

func TestGroupID_String(t *testing.T) {
	groupID, _ := NewGroupID("group-456")
	assert.Equal(t, "group-456", groupID.String())
}

func TestGroupID_Equals(t *testing.T) {
	t.Run("同じIDは等しい", func(t *testing.T) {
		groupID1, _ := NewGroupID("group-123")
		groupID2, _ := NewGroupID("group-123")

		assert.True(t, groupID1.Equals(groupID2))
	})

	t.Run("異なるIDは等しくない", func(t *testing.T) {
		groupID1, _ := NewGroupID("group-123")
		groupID2, _ := NewGroupID("group-456")

		assert.False(t, groupID1.Equals(groupID2))
	})
}
