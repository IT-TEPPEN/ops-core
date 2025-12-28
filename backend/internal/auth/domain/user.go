package domain

import "time"

// OIDCUserInfo represents user information from any OIDC provider
type OIDCUserInfo struct {
	Provider      string // "google", "github", "gitlab", "microsoft"
	Subject       string // Provider's user ID (sub claim)
	Email         string
	EmailVerified bool
	Name          string
	Picture       string
}

// User represents an authenticated user in OpScore
type User struct {
	ID           string
	PrimaryEmail string
	DisplayName  string
	PictureURL   string
	LastLoginAt  time.Time
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

// UserIdentity represents a link between OpScore user and OIDC provider
type UserIdentity struct {
	ID             string
	UserID         string
	Provider       string
	ProviderUserID string
	Email          string
	Name           string
	PictureURL     string
	CreatedAt      time.Time
	UpdatedAt      time.Time
	LastUsedAt     time.Time
}

// NewUser creates a new User from OIDC user info
func NewUser(id string, oidcInfo OIDCUserInfo) *User {
	now := time.Now()
	return &User{
		ID:           id,
		PrimaryEmail: oidcInfo.Email,
		DisplayName:  oidcInfo.Name,
		PictureURL:   oidcInfo.Picture,
		LastLoginAt:  now,
		CreatedAt:    now,
		UpdatedAt:    now,
	}
}

// NewUserIdentity creates a new UserIdentity
func NewUserIdentity(id, userID string, oidcInfo OIDCUserInfo) *UserIdentity {
	now := time.Now()
	return &UserIdentity{
		ID:             id,
		UserID:         userID,
		Provider:       oidcInfo.Provider,
		ProviderUserID: oidcInfo.Subject,
		Email:          oidcInfo.Email,
		Name:           oidcInfo.Name,
		PictureURL:     oidcInfo.Picture,
		CreatedAt:      now,
		UpdatedAt:      now,
		LastUsedAt:     now,
	}
}

// UpdateLastLogin updates the last login timestamp
func (u *User) UpdateLastLogin() {
	u.LastLoginAt = time.Now()
	u.UpdatedAt = time.Now()
}

// UpdateLastUsed updates the last used timestamp for identity
func (ui *UserIdentity) UpdateLastUsed() {
	ui.LastUsedAt = time.Now()
	ui.UpdatedAt = time.Now()
}
