package dto

import "time"

// UserInfo represents the authenticated user's information
type UserInfo struct {
	// UserID is the unique identifier of the user in OpScore
	UserID string
	// Email is the user's primary email address
	Email string
	// DisplayName is the user's display name
	DisplayName string
	// PictureURL is the URL of the user's profile picture
	PictureURL string
	// CreatedAt is when the user was created
	CreatedAt time.Time
	// LastLoginAt is when the user last logged in
	LastLoginAt time.Time
}
