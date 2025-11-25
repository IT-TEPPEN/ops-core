package entity

import (
	"errors"
	"time"

	"opscore/backend/internal/user/domain/value_object"
)

// group represents a group entity with unexported fields
type group struct {
	id          value_object.GroupID
	name        string
	description string
	memberIDs   []value_object.UserID
	createdAt   time.Time
	updatedAt   time.Time
}

// Group interface defines the methods for a group entity
type Group interface {
	ID() value_object.GroupID
	Name() string
	Description() string
	MemberIDs() []value_object.UserID
	CreatedAt() time.Time
	UpdatedAt() time.Time
	UpdateInfo(name string, description string) error
	AddMember(userID value_object.UserID) error
	RemoveMember(userID value_object.UserID) error
}

// NewGroup creates a new Group instance with validation
func NewGroup(id value_object.GroupID, name string, description string) (Group, error) {
	if name == "" {
		return nil, errors.New("group name cannot be empty")
	}

	now := time.Now()
	return &group{
		id:          id,
		name:        name,
		description: description,
		memberIDs:   []value_object.UserID{},
		createdAt:   now,
		updatedAt:   now,
	}, nil
}

// ReconstructGroup reconstructs a Group from persistence data
func ReconstructGroup(
	id value_object.GroupID,
	name string,
	description string,
	memberIDs []value_object.UserID,
	createdAt, updatedAt time.Time,
) Group {
	if memberIDs == nil {
		memberIDs = []value_object.UserID{}
	}
	return &group{
		id:          id,
		name:        name,
		description: description,
		memberIDs:   memberIDs,
		createdAt:   createdAt,
		updatedAt:   updatedAt,
	}
}

// ID returns the group's unique identifier
func (g *group) ID() value_object.GroupID {
	return g.id
}

// Name returns the group's name
func (g *group) Name() string {
	return g.name
}

// Description returns the group's description
func (g *group) Description() string {
	return g.description
}

// MemberIDs returns a copy of the group's member IDs
func (g *group) MemberIDs() []value_object.UserID {
	result := make([]value_object.UserID, len(g.memberIDs))
	copy(result, g.memberIDs)
	return result
}

// CreatedAt returns the timestamp when the group was created
func (g *group) CreatedAt() time.Time {
	return g.createdAt
}

// UpdatedAt returns the timestamp of the last update
func (g *group) UpdatedAt() time.Time {
	return g.updatedAt
}

// UpdateInfo updates the group's information
func (g *group) UpdateInfo(name string, description string) error {
	if name == "" {
		return errors.New("group name cannot be empty")
	}

	g.name = name
	g.description = description
	g.updatedAt = time.Now()
	return nil
}

// AddMember adds a user to the group's member list
func (g *group) AddMember(userID value_object.UserID) error {
	// Check for duplicate
	for _, existingID := range g.memberIDs {
		if existingID.Equals(userID) {
			return errors.New("user is already a member of this group")
		}
	}

	g.memberIDs = append(g.memberIDs, userID)
	g.updatedAt = time.Now()
	return nil
}

// RemoveMember removes a user from the group's member list
func (g *group) RemoveMember(userID value_object.UserID) error {
	for i, existingID := range g.memberIDs {
		if existingID.Equals(userID) {
			g.memberIDs = append(g.memberIDs[:i], g.memberIDs[i+1:]...)
			g.updatedAt = time.Now()
			return nil
		}
	}
	return errors.New("user is not a member of this group")
}
