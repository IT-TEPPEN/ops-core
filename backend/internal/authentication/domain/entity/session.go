package entity

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"time"

	domain_err "opscore/backend/internal/authentication/domain/err"
	"opscore/backend/internal/authentication/domain/value_object"
)

// Session represents an authentication session.
// This entity corresponds to the refresh_tokens table in the database.
type Session struct {
	id         value_object.SessionID
	userID     value_object.UserID
	tokenHash  string
	jti        string
	expiresAt  time.Time
	createdAt  time.Time
	lastUsedAt *time.Time
	revokedAt  *time.Time
	isRevoked  bool
	rememberMe bool
}

// NewSession creates a new Session entity with a generated token.
// This function generates a secure random token and returns both the Session entity
// and the raw token string. The raw token should be sent to the client, while the
// hashed version is stored in the entity.
//
// Parameters:
//   - userID value_object.UserID: The user this session belongs to
//   - expirationDuration time.Duration: How long the session should be valid
//   - rememberMe bool: Whether this is a "remember me" session
//
// Returns:
//   - *Session: The newly created session entity (with token hash)
//   - string: The raw token to be sent to the client
//   - error: Domain error if token generation fails
//
// Errors:
//   - domain_err.ErrUnexpected: When random token generation fails
func NewSession(userID value_object.UserID, expirationDuration time.Duration, rememberMe bool) (*Session, string, error) {
	id := value_object.NewSessionID()
	jti := value_object.NewSessionID().String()
	now := time.Now()
	expiresAt := now.Add(expirationDuration)

	// Generate random token
	tokenBytes := make([]byte, 32)
	if _, err := rand.Read(tokenBytes); err != nil {
		return nil, "", domain_err.NewUnexpectedError("failed to generate session token").WithParent(err)
	}
	rawToken := hex.EncodeToString(tokenBytes)

	// Hash the token for storage
	hash := sha256.Sum256([]byte(rawToken))
	tokenHash := hex.EncodeToString(hash[:])

	session := &Session{
		id:         id,
		userID:     userID,
		tokenHash:  tokenHash,
		jti:        jti,
		expiresAt:  expiresAt,
		createdAt:  now,
		lastUsedAt: nil,
		revokedAt:  nil,
		isRevoked:  false,
		rememberMe: rememberMe,
	}

	return session, rawToken, nil
}

// ReconstructSession reconstructs a Session entity from a trusted data source without validation.
// This function is used when loading data from persistent storage (e.g., database)
// that has already been validated when it was first saved.
//
// Parameters:
//   - id value_object.SessionID: Session identifier
//   - userID value_object.UserID: User identifier
//   - tokenHash string: Hashed token
//   - jti string: JWT ID
//   - expiresAt time.Time: Expiration timestamp
//   - createdAt time.Time: Creation timestamp
//   - lastUsedAt *time.Time: Last usage timestamp (nullable)
//   - revokedAt *time.Time: Revocation timestamp (nullable)
//   - isRevoked bool: Whether the session is revoked
//   - rememberMe bool: Whether this is a "remember me" session
//
// Returns:
//   - *Session: The reconstructed session entity
func ReconstructSession(id value_object.SessionID, userID value_object.UserID, tokenHash, jti string, expiresAt, createdAt time.Time,
	lastUsedAt, revokedAt *time.Time, isRevoked, rememberMe bool) *Session {
	return &Session{
		id:         id,
		userID:     userID,
		tokenHash:  tokenHash,
		jti:        jti,
		expiresAt:  expiresAt,
		createdAt:  createdAt,
		lastUsedAt: lastUsedAt,
		revokedAt:  revokedAt,
		isRevoked:  isRevoked,
		rememberMe: rememberMe,
	}
}

// IsValid checks if the session is valid (not expired and not revoked).
//
// Returns:
//   - bool: true if the session is valid, false otherwise
func (s *Session) IsValid() bool {
	if s.isRevoked {
		return false
	}
	if time.Now().After(s.expiresAt) {
		return false
	}
	return true
}

// UpdateLastUsed updates the last used timestamp.
// This is called when the session is used for token refresh.
func (s *Session) UpdateLastUsed() {
	now := time.Now()
	s.lastUsedAt = &now
}

// Revoke marks the session as revoked.
// A revoked session cannot be used anymore.
func (s *Session) Revoke() {
	now := time.Now()
	s.revokedAt = &now
	s.isRevoked = true
}

// ID returns the session's unique identifier.
//
// Returns:
//   - value_object.SessionID: The session ID
func (s *Session) ID() value_object.SessionID {
	return s.id
}

// UserID returns the user ID this session belongs to.
//
// Returns:
//   - value_object.UserID: The user ID
func (s *Session) UserID() value_object.UserID {
	return s.userID
}

// TokenHash returns the hashed token.
//
// Returns:
//   - string: SHA256 hash of the session token
func (s *Session) TokenHash() string {
	return s.tokenHash
}

// JTI returns the JWT ID for this session.
//
// Returns:
//   - string: JWT ID
func (s *Session) JTI() string {
	return s.jti
}

// ExpiresAt returns the expiration timestamp.
//
// Returns:
//   - time.Time: Expiration timestamp
func (s *Session) ExpiresAt() time.Time {
	return s.expiresAt
}

// CreatedAt returns the creation timestamp.
//
// Returns:
//   - time.Time: Creation timestamp
func (s *Session) CreatedAt() time.Time {
	return s.createdAt
}

// LastUsedAt returns the last usage timestamp.
//
// Returns:
//   - *time.Time: Last usage timestamp (nil if never used)
func (s *Session) LastUsedAt() *time.Time {
	return s.lastUsedAt
}

// RevokedAt returns the revocation timestamp.
//
// Returns:
//   - *time.Time: Revocation timestamp (nil if not revoked)
func (s *Session) RevokedAt() *time.Time {
	return s.revokedAt
}

// IsRevoked returns whether the session is revoked.
//
// Returns:
//   - bool: true if revoked, false otherwise
func (s *Session) IsRevoked() bool {
	return s.isRevoked
}

// RememberMe returns whether this is a "remember me" session.
//
// Returns:
//   - bool: true for long-lived sessions, false for short-lived sessions
func (s *Session) RememberMe() bool {
	return s.rememberMe
}
