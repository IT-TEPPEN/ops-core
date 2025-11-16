package error

// ErrorCode represents a unique error identifier for application layer
// Format: <CONTEXT>_<LAYER>_<CATEGORY>_<NUMBER>
// Example: GITREPO_APP_RES_001
type ErrorCode string

const (
	// Resource errors (GITREPO_APP_RES_xxx)
	CodeResourceNotFound ErrorCode = "GITREPO_APP_RES_001"
	CodeResourceConflict ErrorCode = "GITREPO_APP_RES_002"

	// Authentication/Authorization (GITREPO_APP_AUTH_xxx)
	CodeUnauthorized       ErrorCode = "GITREPO_APP_AUTH_001"
	CodeForbidden          ErrorCode = "GITREPO_APP_AUTH_002"
	CodeInvalidCredentials ErrorCode = "GITREPO_APP_AUTH_003"

	// Validation errors (GITREPO_APP_VAL_xxx)
	CodeValidationFailed ErrorCode = "GITREPO_APP_VAL_001"
	CodeInvalidRequest   ErrorCode = "GITREPO_APP_VAL_002"
)

// String returns the string representation of the error code
func (c ErrorCode) String() string {
	return string(c)
}
