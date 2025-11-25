package value_object

import "errors"

// UserID represents the unique identifier for a User
type UserID struct {
	value string
}

// NewUserID creates a new UserID with validation
func NewUserID(id string) (UserID, error) {
	if id == "" {
		return UserID{}, errors.New("user ID cannot be empty")
	}
	return UserID{value: id}, nil
}

// String returns the string representation of UserID
func (uid UserID) String() string {
	return uid.value
}

// Equals checks if two UserIDs are equal
func (uid UserID) Equals(other UserID) bool {
	return uid.value == other.value
}
