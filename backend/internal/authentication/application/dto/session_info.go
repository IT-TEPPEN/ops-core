package dto

import "time"

// SessionType represents the type of session
type SessionType string

const (
	// SessionTypeShort represents a short-lived session (7 days)
	SessionTypeShort SessionType = "short"
	// SessionTypeLong represents a long-lived session (30 days)
	SessionTypeLong SessionType = "long"
)

// SessionInfo represents session information
type SessionInfo struct {
	// SessionID is the unique identifier of the session
	SessionID string
	// UserID is the user this session belongs to
	UserID string
	// SessionType indicates if this is a short or long session
	SessionType SessionType
	// CreatedAt is when the session was created
	CreatedAt time.Time
	// ExpiresAt is when the session expires
	ExpiresAt time.Time
	// LastUsedAt is when the session was last used
	LastUsedAt *time.Time
	// IsRevoked indicates if the session has been revoked
	IsRevoked bool
}
