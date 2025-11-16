package error

// ErrorCode represents a unique error identifier for infrastructure layer
// Format: <CONTEXT><LAYER><NUMBER>
// - CONTEXT: 2 characters (e.g., GR for git_repository)
// - LAYER: 1 character (D for domain, A for application, I for infrastructure)
// - NUMBER: 4 digits (e.g., 0001)
// Example: GRI0001
type ErrorCode string

const (
	// Database errors
	CodeDatabaseConnection ErrorCode = "GRI0001"
	CodeDatabaseQuery      ErrorCode = "GRI0002"
	CodeDatabaseConstraint ErrorCode = "GRI0003"
	CodeDatabaseTimeout    ErrorCode = "GRI0004"

	// External API errors
	CodeExternalAPIError    ErrorCode = "GRI0005"
	CodeExternalAPITimeout  ErrorCode = "GRI0006"
	CodeExternalAPINotFound ErrorCode = "GRI0007"

	// Connection errors
	CodeConnectionFailed  ErrorCode = "GRI0008"
	CodeConnectionTimeout ErrorCode = "GRI0009"

	// Storage errors
	CodeStorageOperation ErrorCode = "GRI0010"
	CodeStorageNotFound  ErrorCode = "GRI0011"
)

// String returns the string representation of the error code
func (c ErrorCode) String() string {
	return string(c)
}
