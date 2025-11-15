package error

// ErrorCode represents a unique error identifier for application layer
type ErrorCode string

const (
	// Resource errors (APP_RES_xxx)
	CodeResourceNotFound ErrorCode = "APP_RES_001"
	CodeResourceConflict ErrorCode = "APP_RES_002"

	// Authentication/Authorization (APP_AUTH_xxx)
	CodeUnauthorized       ErrorCode = "APP_AUTH_001"
	CodeForbidden          ErrorCode = "APP_AUTH_002"
	CodeInvalidCredentials ErrorCode = "APP_AUTH_003"

	// Validation errors (APP_VAL_xxx)
	CodeValidationFailed ErrorCode = "APP_VAL_001"
	CodeInvalidRequest   ErrorCode = "APP_VAL_002"
)

// String returns the string representation of the error code
func (c ErrorCode) String() string {
	return string(c)
}
