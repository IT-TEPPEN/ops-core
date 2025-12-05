package error

// ErrorCode represents a unique error identifier for infrastructure layer
// Format: <CONTEXT><LAYER><NUMBER>
// - CONTEXT: 2 characters (e.g., US for user)
// - LAYER: 1 character (D for domain, A for application, I for infrastructure)
// - NUMBER: 4 digits (e.g., 0001)
// Example: USI0001
type ErrorCode string

const (
	// Database errors
	CodeDatabaseConnection ErrorCode = "USI0001"
	CodeDatabaseQuery      ErrorCode = "USI0002"
	CodeDatabaseConstraint ErrorCode = "USI0003"
	CodeDatabaseTimeout    ErrorCode = "USI0004"

	// External API errors
	CodeExternalAPIError    ErrorCode = "USI0005"
	CodeExternalAPITimeout  ErrorCode = "USI0006"
	CodeExternalAPINotFound ErrorCode = "USI0007"

	// Connection errors
	CodeConnectionFailed  ErrorCode = "USI0008"
	CodeConnectionTimeout ErrorCode = "USI0009"
)

// String returns the string representation of the error code
func (c ErrorCode) String() string {
	return string(c)
}
