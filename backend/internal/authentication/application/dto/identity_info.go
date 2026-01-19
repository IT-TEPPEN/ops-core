package dto

import "time"

// IdentityInfo represents a user's identity linked to an external provider
type IdentityInfo struct {
	// IdentityID is the unique identifier of the identity
	IdentityID string
	// UserID is the user this identity belongs to
	UserID string
	// Provider is the authentication provider name
	Provider string
	// ProviderUserID is the user's ID in the provider's system
	ProviderUserID string
	// Email is the email associated with this identity
	Email string
	// Name is the name associated with this identity
	Name string
	// PictureURL is the profile picture URL from this identity
	PictureURL string
	// IsPrimary indicates if this is the user's primary identity
	IsPrimary bool
	// CreatedAt is when the identity was linked
	CreatedAt time.Time
	// LastUsedAt is when the identity was last used for authentication
	LastUsedAt time.Time
}
