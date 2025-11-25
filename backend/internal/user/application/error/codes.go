package error

// ErrorCode represents a unique error identifier for application layer
// Format: <CONTEXT><LAYER><NUMBER>
// - CONTEXT: 2 characters (e.g., US for user)
// - LAYER: 1 character (D for domain, A for application, I for infrastructure)
// - NUMBER: 4 digits (e.g., 0001)
// Example: USA0001
type ErrorCode string

const (
	// Resource errors
	CodeResourceNotFound ErrorCode = "USA0001"
	CodeResourceConflict ErrorCode = "USA0002"

	// Authentication/Authorization
	CodeUnauthorized       ErrorCode = "USA0003"
	CodeForbidden          ErrorCode = "USA0004"
	CodeInvalidCredentials ErrorCode = "USA0005"

	// Validation errors
	CodeValidationFailed ErrorCode = "USA0006"
	CodeInvalidRequest   ErrorCode = "USA0007"
)

// String returns the string representation of the error code
func (c ErrorCode) String() string {
	return string(c)
}
