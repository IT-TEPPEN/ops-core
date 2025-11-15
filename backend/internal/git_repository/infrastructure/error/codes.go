package error

// ErrorCode represents a unique error identifier for infrastructure layer
type ErrorCode string

const (
	// Database errors (INF_DB_xxx)
	CodeDatabaseConnection ErrorCode = "INF_DB_001"
	CodeDatabaseQuery      ErrorCode = "INF_DB_002"
	CodeDatabaseConstraint ErrorCode = "INF_DB_003"
	CodeDatabaseTimeout    ErrorCode = "INF_DB_004"

	// External API errors (INF_EXT_xxx)
	CodeExternalAPIError    ErrorCode = "INF_EXT_001"
	CodeExternalAPITimeout  ErrorCode = "INF_EXT_002"
	CodeExternalAPINotFound ErrorCode = "INF_EXT_003"

	// Connection errors (INF_CONN_xxx)
	CodeConnectionFailed  ErrorCode = "INF_CONN_001"
	CodeConnectionTimeout ErrorCode = "INF_CONN_002"

	// Storage errors (INF_STOR_xxx)
	CodeStorageOperation ErrorCode = "INF_STOR_001"
	CodeStorageNotFound  ErrorCode = "INF_STOR_002"
)

// String returns the string representation of the error code
func (c ErrorCode) String() string {
	return string(c)
}
