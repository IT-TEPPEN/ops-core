package value_object

import (
	domain_err "opscore/backend/internal/authentication/domain/err"

	"github.com/google/uuid"
)

// UserID represents a unique identifier for a user.
// This is a value object that wraps a UUID string to provide type safety
// and prevent mixing different ID types.
type UserID string

// NewUserID creates a new UserID with a generated UUID.
// This function never returns an error as UUID generation is guaranteed to succeed.
//
// Returns:
//   - UserID: A newly generated user identifier
func NewUserID() UserID {
	return UserID(uuid.New().String())
}

// UserIDFromString creates a UserID from a string with validation.
// This function is used when receiving user IDs from external sources (e.g., API parameters).
//
// Parameters:
//   - id string: The string representation of the user ID
//
// Returns:
//   - UserID: The validated user identifier
//   - error: Domain error if validation fails
//
// Errors:
//   - domain_err.ErrInvalidUserID: When id is empty or not a valid UUID format
func UserIDFromString(id string) (UserID, error) {
	if id == "" {
		return "", domain_err.NewInvalidUserIDError(id, "user ID cannot be empty")
	}
	// Validate UUID format
	if _, err := uuid.Parse(id); err != nil {
		return "", domain_err.NewInvalidUserIDError(id, "invalid UUID format").WithParent(err)
	}
	return UserID(id), nil
}

// ReconstructUserID reconstructs a UserID from a trusted data source without validation.
// This function is used when loading data from persistent storage (e.g., database)
// that has already been validated when it was first saved.
//
// Parameters:
//   - id string: The string representation of the user ID from a trusted source
//
// Returns:
//   - UserID: The reconstructed user identifier
//
// Note:
//   - This function does not perform validation and assumes the input is valid
//   - Use UserIDFromString() for untrusted input that requires validation
func ReconstructUserID(id string) UserID {
	return UserID(id)
}

// String returns the string representation of the UserID.
//
// Returns:
//   - string: The UUID string value
func (id UserID) String() string {
	return string(id)
}

// IsEmpty checks if the UserID is empty.
//
// Returns:
//   - bool: true if the ID is an empty string, false otherwise
func (id UserID) IsEmpty() bool {
	return id == ""
}

// Equals checks equality with another UserID.
//
// Parameters:
//   - other UserID: The user ID to compare with
//
// Returns:
//   - bool: true if both IDs have the same value, false otherwise
func (id UserID) Equals(other UserID) bool {
	return id == other
}
