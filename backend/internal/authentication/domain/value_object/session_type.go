package value_object

import (
	"time"

	domain_err "opscore/backend/internal/authentication/domain/err"
)

// SessionType represents the type of session (short or long-lived).
// This value object defines the lifetime of authentication sessions.
type SessionType string

const (
	// SessionTypeShort represents a short-lived session (7 days)
	SessionTypeShort SessionType = "short"
	// SessionTypeLong represents a long-lived session (30 days)
	SessionTypeLong SessionType = "long"
)

// SessionTypeFromString creates a SessionType from a string with validation.
//
// Parameters:
//   - value string: The string representation of the session type ("short" or "long")
//
// Returns:
//   - SessionType: The validated session type
//   - error: Domain error if validation fails
//
// Errors:
//   - domain_err.ErrInvalidSessionType: When value is not "short" or "long"
func SessionTypeFromString(value string) (SessionType, error) {
	st := SessionType(value)
	if st != SessionTypeShort && st != SessionTypeLong {
		return "", domain_err.NewInvalidSessionTypeError(value)
	}
	return st, nil
}

// ReconstructSessionType reconstructs a SessionType from a trusted data source without validation.
// This function is used when loading data from persistent storage (e.g., database)
// that has already been validated when it was first saved.
//
// Parameters:
//   - value string: The string representation of the session type from a trusted source
//
// Returns:
//   - SessionType: The reconstructed session type
//
// Note:
//   - This function does not perform validation and assumes the input is valid
//   - Use SessionTypeFromString() for untrusted input that requires validation
func ReconstructSessionType(value string) SessionType {
	return SessionType(value)
}

// ExpirationDuration returns the duration for the session type.
//
// Returns:
//   - time.Duration: 7 days for short sessions, 30 days for long sessions
func (st SessionType) ExpirationDuration() time.Duration {
	switch st {
	case SessionTypeShort:
		return 7 * 24 * time.Hour // 7 days
	case SessionTypeLong:
		return 30 * 24 * time.Hour // 30 days
	default:
		return 7 * 24 * time.Hour // Default to short
	}
}

// RememberMe returns true if this is a long session type.
// This is typically used to determine if the "Remember Me" option was selected.
//
// Returns:
//   - bool: true for long sessions, false for short sessions
func (st SessionType) RememberMe() bool {
	return st == SessionTypeLong
}

// String returns the string representation of the SessionType.
//
// Returns:
//   - string: "short" or "long"
func (st SessionType) String() string {
	return string(st)
}
