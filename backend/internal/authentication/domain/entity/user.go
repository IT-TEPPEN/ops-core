package entity

import (
	"time"

	domain_err "opscore/backend/internal/authentication/domain/err"
	"opscore/backend/internal/authentication/domain/value_object"
)

// User represents an authenticated user in OpScore.
// This is an aggregate root that manages a collection of Identity entities.
// The User aggregate ensures invariants such as:
// - At least one identity must exist
// - Exactly one identity must be marked as primary
type User struct {
	id          value_object.UserID
	email       string
	displayName string
	pictureURL  string
	createdAt   time.Time
	updatedAt   time.Time
	lastLoginAt time.Time

	// Identities collection (child entities)
	identities []*Identity

	// Domain events (for change tracking)
	events []DomainEvent
}

// NewUser creates a new User entity.
// This function is used when creating a new user account.
//
// Parameters:
//   - id value_object.UserID: Unique identifier for the user
//   - email string: User's email address
//   - displayName string: User's display name
//   - pictureURL string: User's profile picture URL
//
// Returns:
//   - *User: The newly created user entity
func NewUser(id value_object.UserID, email, displayName, pictureURL string) *User {
	now := time.Now()
	user := &User{
		id:          id,
		email:       email,
		displayName: displayName,
		pictureURL:  pictureURL,
		createdAt:   now,
		updatedAt:   now,
		lastLoginAt: now,
		identities:  make([]*Identity, 0),
		events:      make([]DomainEvent, 0),
	}

	// Record domain event
	user.recordEvent(NewUserCreatedEvent(id, email, displayName, pictureURL))

	return user
}

// ReconstructUser reconstructs a User entity from a trusted data source without validation.
// This function is used when loading data from persistent storage (e.g., database)
// that has already been validated when it was first saved.
//
// Parameters:
//   - id value_object.UserID: User identifier
//   - email string: Email address
//   - displayName string: Display name
//   - pictureURL string: Profile picture URL
//   - createdAt time.Time: Creation timestamp
//   - updatedAt time.Time: Last update timestamp
//   - lastLoginAt time.Time: Last login timestamp
//
// Returns:
//   - *User: The reconstructed user entity
//
// Note:
//   - Identities are loaded separately using LoadIdentities()
func ReconstructUser(id value_object.UserID, email, displayName, pictureURL string, createdAt, updatedAt, lastLoginAt time.Time) *User {
	return &User{
		id:          id,
		email:       email,
		displayName: displayName,
		pictureURL:  pictureURL,
		createdAt:   createdAt,
		updatedAt:   updatedAt,
		lastLoginAt: lastLoginAt,
		identities:  nil, // Will be loaded separately if needed
		events:      make([]DomainEvent, 0),
	}
}

// UpdateLastLogin updates the last login timestamp.
// This is called when the user successfully authenticates.
func (u *User) UpdateLastLogin() {
	u.lastLoginAt = time.Now()
	u.updatedAt = time.Now()
	u.recordEvent(NewUserLastLoginUpdatedEvent(u.id))
}

// AddIdentity adds a new identity to the user.
// This enforces the invariant that each user must have at least one identity,
// and the first identity is automatically set as primary.
//
// Parameters:
//   - identity *Identity: The identity to add
//
// Returns:
//   - error: Domain error if validation fails
//
// Errors:
//   - domain_err.ErrDuplicateIdentity: When an identity with the same provider already exists
func (u *User) AddIdentity(identity *Identity) error {
	// Check if identity with same provider already exists
	for _, existing := range u.identities {
		if existing.IsSameProvider(identity.Provider(), identity.ProviderUserID()) {
			return domain_err.NewDuplicateIdentityError(identity.Provider(), identity.ProviderUserID())
		}
	}

	// If this is the first identity, set it as primary
	if len(u.identities) == 0 {
		identity.SetAsPrimary()
	}

	u.identities = append(u.identities, identity)
	u.updatedAt = time.Now()
	u.recordEvent(NewIdentityAddedEvent(identity))

	return nil
}

// RemoveIdentity removes an identity from the user.
// This enforces the invariant that at least one identity must remain,
// and the primary identity cannot be removed without first setting another as primary.
//
// Parameters:
//   - identityID value_object.IdentityID: The identity ID to remove
//
// Returns:
//   - error: Domain error if validation fails
//
// Errors:
//   - domain_err.ErrCannotRemoveOnlyIdentity: When trying to remove the last identity
//   - domain_err.ErrCannotUnlinkPrimaryIdentity: When trying to remove the primary identity
func (u *User) RemoveIdentity(identityID value_object.IdentityID) error {
	if len(u.identities) <= 1 {
		return domain_err.NewCannotRemoveOnlyIdentityError(u.id.String())
	}

	var targetIdentity *Identity
	for _, identity := range u.identities {
		if identity.ID() == identityID {
			targetIdentity = identity
			break
		}
	}

	if targetIdentity == nil {
		return domain_err.NewUnexpectedError("identity not found in user's collection")
	}

	if targetIdentity.IsPrimary() {
		return domain_err.NewCannotUnlinkPrimaryIdentityError(identityID.String())
	}

	// Remove from collection
	newIdentities := make([]*Identity, 0, len(u.identities)-1)
	for _, identity := range u.identities {
		if identity.ID() != identityID {
			newIdentities = append(newIdentities, identity)
		}
	}
	u.identities = newIdentities
	u.updatedAt = time.Now()
	u.recordEvent(NewIdentityRemovedEvent(identityID, u.id))

	return nil
}

// SetPrimaryIdentity sets the specified identity as primary.
// This enforces the invariant that exactly one identity must be marked as primary.
//
// Parameters:
//   - identityID value_object.IdentityID: The identity ID to set as primary
//
// Returns:
//   - error: Domain error if validation fails
//
// Errors:
//   - domain_err.ErrUnexpected: When the specified identity is not found
func (u *User) SetPrimaryIdentity(identityID value_object.IdentityID) error {
	var targetIdentity *Identity
	var oldPrimaryIdentity *Identity

	for _, identity := range u.identities {
		if identity.IsPrimary() {
			oldPrimaryIdentity = identity
		}
		if identity.ID() == identityID {
			targetIdentity = identity
		}
	}

	if targetIdentity == nil {
		return domain_err.NewUnexpectedError("identity not found in user's collection")
	}

	if targetIdentity.IsPrimary() {
		return nil // Already primary
	}

	// Switch primary flag
	if oldPrimaryIdentity != nil {
		oldPrimaryIdentity.UnsetPrimary()
	}
	targetIdentity.SetAsPrimary()

	u.updatedAt = time.Now()
	oldIdentityID := value_object.IdentityID("")
	if oldPrimaryIdentity != nil {
		oldIdentityID = oldPrimaryIdentity.ID()
	}
	u.recordEvent(NewPrimaryIdentityChangedEvent(u.id, oldIdentityID, identityID))

	return nil
}

// GetPrimaryIdentity returns the primary identity.
//
// Returns:
//   - *Identity: The primary identity, or nil if not found
func (u *User) GetPrimaryIdentity() *Identity {
	for _, identity := range u.identities {
		if identity.IsPrimary() {
			return identity
		}
	}
	return nil
}

// LoadIdentities loads identities into the aggregate.
// This is called by the repository after reconstructing a User.
//
// Parameters:
//   - identities []*Identity: The identities to load
func (u *User) LoadIdentities(identities []*Identity) {
	u.identities = identities
}

// Identities returns all identities.
//
// Returns:
//   - []*Identity: The list of identities
func (u *User) Identities() []*Identity {
	return u.identities
}

// HasIdentitiesLoaded checks if identities are loaded.
//
// Returns:
//   - bool: true if identities have been loaded, false otherwise
func (u *User) HasIdentitiesLoaded() bool {
	return u.identities != nil
}

// recordEvent records a domain event for persistence tracking.
func (u *User) recordEvent(event DomainEvent) {
	u.events = append(u.events, event)
}

// GetEvents returns all recorded domain events.
// These events are used by the repository to determine what has changed.
//
// Returns:
//   - []DomainEvent: The list of domain events
func (u *User) GetEvents() []DomainEvent {
	return u.events
}

// ClearEvents clears all recorded domain events.
// This is called by the repository after successfully persisting changes.
func (u *User) ClearEvents() {
	u.events = make([]DomainEvent, 0)
}

// ID returns the user's unique identifier.
//
// Returns:
//   - value_object.UserID: The user ID
func (u *User) ID() value_object.UserID {
	return u.id
}

// Email returns the user's email address.
//
// Returns:
//   - string: Email address
func (u *User) Email() string {
	return u.email
}

// DisplayName returns the user's display name.
//
// Returns:
//   - string: Display name
func (u *User) DisplayName() string {
	return u.displayName
}

// PictureURL returns the user's profile picture URL.
//
// Returns:
//   - string: Profile picture URL
func (u *User) PictureURL() string {
	return u.pictureURL
}

// CreatedAt returns the creation timestamp.
//
// Returns:
//   - time.Time: Creation timestamp
func (u *User) CreatedAt() time.Time {
	return u.createdAt
}

// UpdatedAt returns the last update timestamp.
//
// Returns:
//   - time.Time: Last update timestamp
func (u *User) UpdatedAt() time.Time {
	return u.updatedAt
}

// LastLoginAt returns the last login timestamp.
//
// Returns:
//   - time.Time: Last login timestamp
func (u *User) LastLoginAt() time.Time {
	return u.lastLoginAt
}
