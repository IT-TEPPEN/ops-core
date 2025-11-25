package value_object

import (
	"errors"
	"strings"
)

// Role represents a user role (admin or user)
type Role struct {
	value string
}

// Predefined roles
const (
	RoleAdmin = "admin"
	RoleUser  = "user"
)

// NewRole creates a new Role with validation
func NewRole(role string) (Role, error) {
	normalized := strings.ToLower(strings.TrimSpace(role))
	
	if normalized != RoleAdmin && normalized != RoleUser {
		return Role{}, errors.New("role must be either 'admin' or 'user'")
	}
	
	return Role{value: normalized}, nil
}

// String returns the string representation of Role
func (r Role) String() string {
	return r.value
}

// IsAdmin returns true if the role is admin
func (r Role) IsAdmin() bool {
	return r.value == RoleAdmin
}

// Equals checks if two Roles are equal
func (r Role) Equals(other Role) bool {
	return r.value == other.value
}
