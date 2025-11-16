package error

// ErrorCode represents a unique error identifier for application layer
// Format: <CONTEXT><LAYER><NUMBER>
// - CONTEXT: 2 characters (e.g., GR for git_repository)
// - LAYER: 1 character (D for domain, A for application, I for infrastructure)
// - NUMBER: 4 digits (e.g., 0001)
// Example: GRA0001
type ErrorCode string

const (
	// Resource errors
	CodeResourceNotFound ErrorCode = "GRA0001"
	CodeResourceConflict ErrorCode = "GRA0002"

	// Authentication/Authorization
	CodeUnauthorized       ErrorCode = "GRA0003"
	CodeForbidden          ErrorCode = "GRA0004"
	CodeInvalidCredentials ErrorCode = "GRA0005"

	// Validation errors
	CodeValidationFailed ErrorCode = "GRA0006"
	CodeInvalidRequest   ErrorCode = "GRA0007"
)

// String returns the string representation of the error code
func (c ErrorCode) String() string {
	return string(c)
}
