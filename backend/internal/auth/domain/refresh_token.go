package domain

import "time"

// RefreshToken represents a refresh token for JWT authentication
type RefreshToken struct {
	ID         string
	UserID     string
	TokenHash  string    // SHA-256 hash of the token
	JTI        string    // JWT ID (unique identifier)
	ExpiresAt  time.Time
	CreatedAt  time.Time
	LastUsedAt *time.Time
	RevokedAt  *time.Time
	IsRevoked  bool
	RememberMe bool // true = 30 days, false = 7 days
}

// NewRefreshToken creates a new RefreshToken
func NewRefreshToken(id, userID, tokenHash, jti string, expiresAt time.Time, rememberMe bool) *RefreshToken {
	now := time.Now()
	return &RefreshToken{
		ID:         id,
		UserID:     userID,
		TokenHash:  tokenHash,
		JTI:        jti,
		ExpiresAt:  expiresAt,
		CreatedAt:  now,
		LastUsedAt: nil,
		RevokedAt:  nil,
		IsRevoked:  false,
		RememberMe: rememberMe,
	}
}

// IsValid checks if the token is valid (not expired and not revoked)
func (rt *RefreshToken) IsValid() bool {
	if rt.IsRevoked {
		return false
	}
	if time.Now().After(rt.ExpiresAt) {
		return false
	}
	return true
}

// UpdateLastUsed updates the last used timestamp
func (rt *RefreshToken) UpdateLastUsed() {
	now := time.Now()
	rt.LastUsedAt = &now
}

// Revoke marks the token as revoked
func (rt *RefreshToken) Revoke() {
	now := time.Now()
	rt.RevokedAt = &now
	rt.IsRevoked = true
}
