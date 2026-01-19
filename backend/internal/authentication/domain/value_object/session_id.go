package value_object

import (
	domain_err "opscore/backend/internal/authentication/domain/err"

	"github.com/google/uuid"
)

// SessionID represents a unique identifier for a session.
// This is a value object that wraps a UUID string to provide type safety
// and prevent mixing different ID types.
type SessionID string

// NewSessionID creates a new SessionID with a generated UUID.
// This function never returns an error as UUID generation is guaranteed to succeed.
//
// Returns:
//   - SessionID: A newly generated session identifier
func NewSessionID() SessionID {
	return SessionID(uuid.New().String())
}

// SessionIDFromString creates a SessionID from a string with validation.
// This function is used when receiving session IDs from external sources
// (e.g., cookies, authorization headers).
//
// Parameters:
//   - id string: The string representation of the session ID
//
// Returns:
//   - SessionID: The validated session identifier
//   - error: Domain error if validation fails
//
// Errors:
//   - domain_err.ErrInvalidSessionID: When id is empty or not a valid UUID format
func SessionIDFromString(id string) (SessionID, error) {
	if id == "" {
		return "", domain_err.NewInvalidSessionIDError(id, "session ID cannot be empty")
	}
	// Validate UUID format
	if _, err := uuid.Parse(id); err != nil {
		return "", domain_err.NewInvalidSessionIDError(id, "invalid UUID format").WithParent(err)
	}
	return SessionID(id), nil
}

// ReconstructSessionID reconstructs a SessionID from a trusted data source without validation.
// This function is used when loading data from persistent storage (e.g., database)
// that has already been validated when it was first saved.
//
// Parameters:
//   - id string: The string representation of the session ID from a trusted source
//
// Returns:
//   - SessionID: The reconstructed session identifier
//
// Note:
//   - This function does not perform validation and assumes the input is valid
//   - Use SessionIDFromString() for untrusted input that requires validation
func ReconstructSessionID(id string) SessionID {
	return SessionID(id)
}

// String returns the string representation of the SessionID.
//
// Returns:
//   - string: The UUID string value
func (id SessionID) String() string {
	return string(id)
}

// IsEmpty checks if the SessionID is empty.
//
// Returns:
//   - bool: true if the ID is an empty string, false otherwise
func (id SessionID) IsEmpty() bool {
	return id == ""
}

// Equals checks equality with another SessionID.
//
// Parameters:
//   - other SessionID: The session ID to compare with
//
// Returns:
//   - bool: true if both IDs have the same value, false otherwise
func (id SessionID) Equals(other SessionID) bool {
	return id == other
}
