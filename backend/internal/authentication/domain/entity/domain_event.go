package entity

import (
	"time"

	"opscore/backend/internal/authentication/domain/value_object"
)

// DomainEvent represents a domain event that occurred in the aggregate
type DomainEvent interface {
	OccurredAt() time.Time
}

// UserCreatedEvent is emitted when a new user is created
type UserCreatedEvent struct {
	UserID      value_object.UserID
	Email       string
	DisplayName string
	PictureURL  string
	occurredAt  time.Time
}

func NewUserCreatedEvent(userID value_object.UserID, email, displayName, pictureURL string) *UserCreatedEvent {
	return &UserCreatedEvent{
		UserID:      userID,
		Email:       email,
		DisplayName: displayName,
		PictureURL:  pictureURL,
		occurredAt:  time.Now(),
	}
}

func (e *UserCreatedEvent) OccurredAt() time.Time {
	return e.occurredAt
}

// IdentityAddedEvent is emitted when an identity is added to a user
type IdentityAddedEvent struct {
	IdentityID     value_object.IdentityID
	UserID         value_object.UserID
	Provider       string
	ProviderUserID string
	Email          string
	Name           string
	PictureURL     string
	IsPrimary      bool
	occurredAt     time.Time
}

func NewIdentityAddedEvent(identity *Identity) *IdentityAddedEvent {
	return &IdentityAddedEvent{
		IdentityID:     identity.ID(),
		UserID:         identity.UserID(),
		Provider:       identity.Provider(),
		ProviderUserID: identity.ProviderUserID(),
		Email:          identity.Email(),
		Name:           identity.Name(),
		PictureURL:     identity.PictureURL(),
		IsPrimary:      identity.IsPrimary(),
		occurredAt:     time.Now(),
	}
}

func (e *IdentityAddedEvent) OccurredAt() time.Time {
	return e.occurredAt
}

// IdentityUpdatedEvent is emitted when an identity is updated
type IdentityUpdatedEvent struct {
	IdentityID value_object.IdentityID
	Email      string
	Name       string
	PictureURL string
	occurredAt time.Time
}

func NewIdentityUpdatedEvent(identityID value_object.IdentityID, email, name, pictureURL string) *IdentityUpdatedEvent {
	return &IdentityUpdatedEvent{
		IdentityID: identityID,
		Email:      email,
		Name:       name,
		PictureURL: pictureURL,
		occurredAt: time.Now(),
	}
}

func (e *IdentityUpdatedEvent) OccurredAt() time.Time {
	return e.occurredAt
}

// IdentityRemovedEvent is emitted when an identity is removed
type IdentityRemovedEvent struct {
	IdentityID value_object.IdentityID
	UserID     value_object.UserID
	occurredAt time.Time
}

func NewIdentityRemovedEvent(identityID value_object.IdentityID, userID value_object.UserID) *IdentityRemovedEvent {
	return &IdentityRemovedEvent{
		IdentityID: identityID,
		UserID:     userID,
		occurredAt: time.Now(),
	}
}

func (e *IdentityRemovedEvent) OccurredAt() time.Time {
	return e.occurredAt
}

// PrimaryIdentityChangedEvent is emitted when the primary identity changes
type PrimaryIdentityChangedEvent struct {
	UserID        value_object.UserID
	OldIdentityID value_object.IdentityID
	NewIdentityID value_object.IdentityID
	occurredAt    time.Time
}

func NewPrimaryIdentityChangedEvent(userID value_object.UserID, oldIdentityID, newIdentityID value_object.IdentityID) *PrimaryIdentityChangedEvent {
	return &PrimaryIdentityChangedEvent{
		UserID:        userID,
		OldIdentityID: oldIdentityID,
		NewIdentityID: newIdentityID,
		occurredAt:    time.Now(),
	}
}

func (e *PrimaryIdentityChangedEvent) OccurredAt() time.Time {
	return e.occurredAt
}

// UserLastLoginUpdatedEvent is emitted when user's last login is updated
type UserLastLoginUpdatedEvent struct {
	UserID     value_object.UserID
	occurredAt time.Time
}

func NewUserLastLoginUpdatedEvent(userID value_object.UserID) *UserLastLoginUpdatedEvent {
	return &UserLastLoginUpdatedEvent{
		UserID:     userID,
		occurredAt: time.Now(),
	}
}

func (e *UserLastLoginUpdatedEvent) OccurredAt() time.Time {
	return e.occurredAt
}
