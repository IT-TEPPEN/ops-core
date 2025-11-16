package error

// ErrorCode represents a unique error identifier for infrastructure layer
// Format: <CONTEXT>_<LAYER>_<CATEGORY>_<NUMBER>
// Example: GITREPO_INF_DB_001
type ErrorCode string

const (
	// Database errors (GITREPO_INF_DB_xxx)
	CodeDatabaseConnection ErrorCode = "GITREPO_INF_DB_001"
	CodeDatabaseQuery      ErrorCode = "GITREPO_INF_DB_002"
	CodeDatabaseConstraint ErrorCode = "GITREPO_INF_DB_003"
	CodeDatabaseTimeout    ErrorCode = "GITREPO_INF_DB_004"

	// External API errors (GITREPO_INF_EXT_xxx)
	CodeExternalAPIError    ErrorCode = "GITREPO_INF_EXT_001"
	CodeExternalAPITimeout  ErrorCode = "GITREPO_INF_EXT_002"
	CodeExternalAPINotFound ErrorCode = "GITREPO_INF_EXT_003"

	// Connection errors (GITREPO_INF_CONN_xxx)
	CodeConnectionFailed  ErrorCode = "GITREPO_INF_CONN_001"
	CodeConnectionTimeout ErrorCode = "GITREPO_INF_CONN_002"

	// Storage errors (GITREPO_INF_STOR_xxx)
	CodeStorageOperation ErrorCode = "GITREPO_INF_STOR_001"
	CodeStorageNotFound  ErrorCode = "GITREPO_INF_STOR_002"
)

// String returns the string representation of the error code
func (c ErrorCode) String() string {
	return string(c)
}
