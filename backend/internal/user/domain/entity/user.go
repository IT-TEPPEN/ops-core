package entity

import (
	"errors"
	"time"

	"opscore/backend/internal/user/domain/value_object"
)

// user represents a user entity with unexported fields
type user struct {
	id        value_object.UserID
	name      string
	email     value_object.Email
	role      value_object.Role
	groupIDs  []value_object.GroupID
	createdAt time.Time
	updatedAt time.Time
}

// User interface defines the methods for a user entity
type User interface {
	ID() value_object.UserID
	Name() string
	Email() value_object.Email
	Role() value_object.Role
	GroupIDs() []value_object.GroupID
	CreatedAt() time.Time
	UpdatedAt() time.Time
	UpdateProfile(name string, email value_object.Email) error
	JoinGroup(groupID value_object.GroupID) error
	LeaveGroup(groupID value_object.GroupID) error
	ChangeRole(role value_object.Role) error
}

// NewUser creates a new User instance with validation
func NewUser(id value_object.UserID, name string, email value_object.Email, role value_object.Role) (User, error) {
	if name == "" {
		return nil, errors.New("user name cannot be empty")
	}

	now := time.Now()
	return &user{
		id:        id,
		name:      name,
		email:     email,
		role:      role,
		groupIDs:  []value_object.GroupID{},
		createdAt: now,
		updatedAt: now,
	}, nil
}

// ReconstructUser reconstructs a User from persistence data
func ReconstructUser(
	id value_object.UserID,
	name string,
	email value_object.Email,
	role value_object.Role,
	groupIDs []value_object.GroupID,
	createdAt, updatedAt time.Time,
) User {
	if groupIDs == nil {
		groupIDs = []value_object.GroupID{}
	}
	return &user{
		id:        id,
		name:      name,
		email:     email,
		role:      role,
		groupIDs:  groupIDs,
		createdAt: createdAt,
		updatedAt: updatedAt,
	}
}

// ID returns the user's unique identifier
func (u *user) ID() value_object.UserID {
	return u.id
}

// Name returns the user's name
func (u *user) Name() string {
	return u.name
}

// Email returns the user's email address
func (u *user) Email() value_object.Email {
	return u.email
}

// Role returns the user's role
func (u *user) Role() value_object.Role {
	return u.role
}

// GroupIDs returns a copy of the user's group IDs
func (u *user) GroupIDs() []value_object.GroupID {
	result := make([]value_object.GroupID, len(u.groupIDs))
	copy(result, u.groupIDs)
	return result
}

// CreatedAt returns the timestamp when the user was created
func (u *user) CreatedAt() time.Time {
	return u.createdAt
}

// UpdatedAt returns the timestamp of the last update
func (u *user) UpdatedAt() time.Time {
	return u.updatedAt
}

// UpdateProfile updates the user's profile information
func (u *user) UpdateProfile(name string, email value_object.Email) error {
	if name == "" {
		return errors.New("user name cannot be empty")
	}

	u.name = name
	u.email = email
	u.updatedAt = time.Now()
	return nil
}

// JoinGroup adds a group to the user's group list
func (u *user) JoinGroup(groupID value_object.GroupID) error {
	// Check for duplicate
	for _, existingID := range u.groupIDs {
		if existingID.Equals(groupID) {
			return errors.New("user is already a member of this group")
		}
	}

	u.groupIDs = append(u.groupIDs, groupID)
	u.updatedAt = time.Now()
	return nil
}

// LeaveGroup removes a group from the user's group list
func (u *user) LeaveGroup(groupID value_object.GroupID) error {
	for i, existingID := range u.groupIDs {
		if existingID.Equals(groupID) {
			u.groupIDs = append(u.groupIDs[:i], u.groupIDs[i+1:]...)
			u.updatedAt = time.Now()
			return nil
		}
	}
	return errors.New("user is not a member of this group")
}

// ChangeRole changes the user's role
func (u *user) ChangeRole(role value_object.Role) error {
	u.role = role
	u.updatedAt = time.Now()
	return nil
}
