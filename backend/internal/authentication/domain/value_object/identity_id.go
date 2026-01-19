package value_object

import (
	"github.com/google/uuid"
)

// IdentityID represents a unique identifier for an identity.
// This is a value object that wraps a UUID string to provide type safety
// and prevent mixing different ID types.
type IdentityID string

// NewIdentityID creates a new IdentityID with a generated UUID.
// This function never returns an error as UUID generation is guaranteed to succeed.
//
// Returns:
//   - IdentityID: A newly generated identity identifier
func NewIdentityID() IdentityID {
	return IdentityID(uuid.New().String())
}

// ReconstructIdentityID reconstructs an IdentityID from a trusted data source without validation.
// This function is used when loading data from persistent storage (e.g., database)
// that has already been validated when it was first saved.
//
// Parameters:
//   - id string: The string representation of the identity ID from a trusted source
//
// Returns:
//   - IdentityID: The reconstructed identity identifier
//
// Note:
//   - This function does not perform validation and assumes the input is valid
//   - Use IdentityIDFromString() for untrusted input that requires validation
func ReconstructIdentityID(id string) IdentityID {
	return IdentityID(id)
}

// String returns the string representation of the IdentityID.
//
// Returns:
//   - string: The UUID string value
func (id IdentityID) String() string {
	return string(id)
}

// IsEmpty checks if the IdentityID is empty.
//
// Returns:
//   - bool: true if the ID is an empty string, false otherwise
func (id IdentityID) IsEmpty() bool {
	return id == ""
}

// Equals checks equality with another IdentityID.
//
// Parameters:
//   - other IdentityID: The identity ID to compare with
//
// Returns:
//   - bool: true if both IDs have the same value, false otherwise
func (id IdentityID) Equals(other IdentityID) bool {
	return id == other
}
