package entity

import (
	"time"

	domain_err "opscore/backend/internal/authentication/domain/err"
	"opscore/backend/internal/authentication/domain/value_object"
)

// Identity represents a link between OpScore user and external authentication provider.
// This is a child entity within the User aggregate.
type Identity struct {
	id             value_object.IdentityID
	userID         value_object.UserID
	provider       string
	providerUserID string
	email          string
	name           string
	pictureURL     string
	isPrimary      bool
	createdAt      time.Time
	updatedAt      time.Time
	lastUsedAt     time.Time
}

// NewIdentity creates a new Identity entity with validation.
// This function is used when linking a new external provider account to a user.
//
// Parameters:
//   - id value_object.IdentityID: Unique identifier for this identity
//   - userID value_object.UserID: The user this identity belongs to
//   - provider string: Provider name (e.g., "google", "github")
//   - providerUserID string: User ID from the provider
//   - email string: Email address from the provider
//   - name string: Display name from the provider
//   - pictureURL string: Profile picture URL from the provider
//
// Returns:
//   - *Identity: The newly created identity entity
//   - error: Domain error if validation fails
//
// Errors:
//   - domain_err.ErrInvalidProviderUserID: When providerUserID is empty
func NewIdentity(id value_object.IdentityID, userID value_object.UserID, provider, providerUserID, email, name, pictureURL string) (*Identity, error) {
	if providerUserID == "" {
		return nil, domain_err.NewInvalidProviderUserIDError(providerUserID)
	}

	now := time.Now()
	return &Identity{
		id:             id,
		userID:         userID,
		provider:       provider,
		providerUserID: providerUserID,
		email:          email,
		name:           name,
		pictureURL:     pictureURL,
		isPrimary:      false,
		createdAt:      now,
		updatedAt:      now,
		lastUsedAt:     now,
	}, nil
}

// ReconstructIdentity reconstructs an Identity entity from a trusted data source without validation.
// This function is used when loading data from persistent storage (e.g., database)
// that has already been validated when it was first saved.
//
// Parameters:
//   - id value_object.IdentityID: Unique identifier
//   - userID value_object.UserID: Owner user ID
//   - provider string: Provider name
//   - providerUserID string: Provider user ID
//   - email string: Email address
//   - name string: Display name
//   - pictureURL string: Profile picture URL
//   - isPrimary bool: Whether this is the primary identity
//   - createdAt time.Time: Creation timestamp
//   - updatedAt time.Time: Last update timestamp
//   - lastUsedAt time.Time: Last usage timestamp
//
// Returns:
//   - *Identity: The reconstructed identity entity
func ReconstructIdentity(id value_object.IdentityID, userID value_object.UserID, provider, providerUserID, email, name, pictureURL string,
	isPrimary bool, createdAt, updatedAt, lastUsedAt time.Time) *Identity {
	return &Identity{
		id:             id,
		userID:         userID,
		provider:       provider,
		providerUserID: providerUserID,
		email:          email,
		name:           name,
		pictureURL:     pictureURL,
		isPrimary:      isPrimary,
		createdAt:      createdAt,
		updatedAt:      updatedAt,
		lastUsedAt:     lastUsedAt,
	}
}

// UpdateLastUsed updates the last used timestamp.
// This is called when the identity is used for authentication.
func (i *Identity) UpdateLastUsed() {
	i.lastUsedAt = time.Now()
	i.updatedAt = time.Now()
}

// UpdateInfo updates the identity information (email, name, picture).
// This is called when provider information changes.
//
// Parameters:
//   - email string: New email address
//   - name string: New display name
//   - pictureURL string: New profile picture URL
func (i *Identity) UpdateInfo(email, name, pictureURL string) {
	i.email = email
	i.name = name
	i.pictureURL = pictureURL
	i.updatedAt = time.Now()
}

// IsSameProvider checks if the identity belongs to the same provider and provider user.
//
// Parameters:
//   - provider string: Provider name to check
//   - providerUserID string: Provider user ID to check
//
// Returns:
//   - bool: true if both provider and providerUserID match, false otherwise
func (i *Identity) IsSameProvider(provider, providerUserID string) bool {
	return i.provider == provider && i.providerUserID == providerUserID
}

// ID returns the identity's unique identifier.
//
// Returns:
//   - value_object.IdentityID: The identity ID
func (i *Identity) ID() value_object.IdentityID {
	return i.id
}

// UserID returns the user ID this identity belongs to.
//
// Returns:
//   - value_object.UserID: The user ID
func (i *Identity) UserID() value_object.UserID {
	return i.userID
}

// Provider returns the authentication provider name.
//
// Returns:
//   - string: Provider name (e.g., "google", "github")
func (i *Identity) Provider() string {
	return i.provider
}

// ProviderUserID returns the user ID from the provider.
//
// Returns:
//   - string: Provider user ID
func (i *Identity) ProviderUserID() string {
	return i.providerUserID
}

// Email returns the email address from the provider.
//
// Returns:
//   - string: Email address
func (i *Identity) Email() string {
	return i.email
}

// Name returns the display name from the provider.
//
// Returns:
//   - string: Display name
func (i *Identity) Name() string {
	return i.name
}

// PictureURL returns the profile picture URL from the provider.
//
// Returns:
//   - string: Profile picture URL
func (i *Identity) PictureURL() string {
	return i.pictureURL
}

// CreatedAt returns the creation timestamp.
//
// Returns:
//   - time.Time: Creation timestamp
func (i *Identity) CreatedAt() time.Time {
	return i.createdAt
}

// UpdatedAt returns the last update timestamp.
//
// Returns:
//   - time.Time: Last update timestamp
func (i *Identity) UpdatedAt() time.Time {
	return i.updatedAt
}

// LastUsedAt returns the last usage timestamp.
//
// Returns:
//   - time.Time: Last usage timestamp
func (i *Identity) LastUsedAt() time.Time {
	return i.lastUsedAt
}

// IsPrimary returns whether this is the primary identity.
//
// Returns:
//   - bool: true if this is the primary identity, false otherwise
func (i *Identity) IsPrimary() bool {
	return i.isPrimary
}

// SetAsPrimary marks this identity as primary.
// This should only be called from the User aggregate to maintain invariants.
func (i *Identity) SetAsPrimary() {
	i.isPrimary = true
	i.updatedAt = time.Now()
}

// UnsetPrimary removes the primary flag.
// This should only be called from the User aggregate to maintain invariants.
func (i *Identity) UnsetPrimary() {
	i.isPrimary = false
	i.updatedAt = time.Now()
}
