package error

// ErrorCode represents error codes for the document application layer
type ErrorCode string

const (
	CodeResourceNotFound ErrorCode = "RESOURCE_NOT_FOUND"
	CodeValidationFailed ErrorCode = "VALIDATION_FAILED"
	CodeResourceConflict ErrorCode = "RESOURCE_CONFLICT"
	CodeUnauthorized     ErrorCode = "UNAUTHORIZED"
	CodeForbidden        ErrorCode = "FORBIDDEN"
)
