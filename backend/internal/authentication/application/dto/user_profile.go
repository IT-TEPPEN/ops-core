package dto

// UserProfile represents user's customizable profile information
type UserProfile struct {
	// UserID is the user this profile belongs to
	UserID string
	// DisplayName is the custom display name (nil means use primary identity name)
	DisplayName *string
	// PictureURL is the custom profile picture URL (nil means use primary identity picture)
	PictureURL *string
	// Bio is the user's biography
	Bio *string
}
