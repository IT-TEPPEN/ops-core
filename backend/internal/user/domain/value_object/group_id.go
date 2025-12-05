package value_object

import "errors"

// GroupID represents the unique identifier for a Group
type GroupID struct {
	value string
}

// NewGroupID creates a new GroupID with validation
func NewGroupID(id string) (GroupID, error) {
	if id == "" {
		return GroupID{}, errors.New("group ID cannot be empty")
	}
	return GroupID{value: id}, nil
}

// String returns the string representation of GroupID
func (gid GroupID) String() string {
	return gid.value
}

// Equals checks if two GroupIDs are equal
func (gid GroupID) Equals(other GroupID) bool {
	return gid.value == other.value
}
