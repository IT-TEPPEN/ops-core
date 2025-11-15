package error

// ErrorCodeInfo provides metadata about an error code
type ErrorCodeInfo struct {
	Code        string
	Layer       string
	Category    string
	Description string
	Severity    string // "LOW", "MEDIUM", "HIGH", "CRITICAL"
	Retryable   bool
}

// ErrorCodeRegistry is the global error code registry
var ErrorCodeRegistry = map[string]ErrorCodeInfo{
	// Domain Layer - Validation Errors
	"DOM_VAL_001": {
		Code:        "DOM_VAL_001",
		Layer:       "Domain",
		Category:    "Validation",
		Description: "Invalid entity field value",
		Severity:    "MEDIUM",
		Retryable:   false,
	},
	"DOM_VAL_002": {
		Code:        "DOM_VAL_002",
		Layer:       "Domain",
		Category:    "Validation",
		Description: "Required field is missing",
		Severity:    "MEDIUM",
		Retryable:   false,
	},
	"DOM_VAL_003": {
		Code:        "DOM_VAL_003",
		Layer:       "Domain",
		Category:    "Validation",
		Description: "Invalid field format",
		Severity:    "MEDIUM",
		Retryable:   false,
	},
	"DOM_VAL_004": {
		Code:        "DOM_VAL_004",
		Layer:       "Domain",
		Category:    "Validation",
		Description: "Field value out of range",
		Severity:    "MEDIUM",
		Retryable:   false,
	},
	"DOM_VAL_005": {
		Code:        "DOM_VAL_005",
		Layer:       "Domain",
		Category:    "Validation",
		Description: "Invalid URL format",
		Severity:    "MEDIUM",
		Retryable:   false,
	},
	"DOM_VAL_006": {
		Code:        "DOM_VAL_006",
		Layer:       "Domain",
		Category:    "Validation",
		Description: "Unsupported URL scheme (only HTTPS is supported)",
		Severity:    "MEDIUM",
		Retryable:   false,
	},

	// Domain Layer - Business Rule Violations
	"DOM_BUS_001": {
		Code:        "DOM_BUS_001",
		Layer:       "Domain",
		Category:    "Business",
		Description: "Business rule violation",
		Severity:    "HIGH",
		Retryable:   false,
	},
	"DOM_BUS_002": {
		Code:        "DOM_BUS_002",
		Layer:       "Domain",
		Category:    "Business",
		Description: "Invalid state transition",
		Severity:    "HIGH",
		Retryable:   false,
	},
	"DOM_BUS_003": {
		Code:        "DOM_BUS_003",
		Layer:       "Domain",
		Category:    "Business",
		Description: "Invariant violation",
		Severity:    "HIGH",
		Retryable:   false,
	},

	// Application Layer - Resource Errors
	"APP_RES_001": {
		Code:        "APP_RES_001",
		Layer:       "Application",
		Category:    "Resource",
		Description: "Requested resource not found",
		Severity:    "LOW",
		Retryable:   false,
	},
	"APP_RES_002": {
		Code:        "APP_RES_002",
		Layer:       "Application",
		Category:    "Resource",
		Description: "Resource conflict (duplicate)",
		Severity:    "MEDIUM",
		Retryable:   false,
	},

	// Application Layer - Authentication/Authorization
	"APP_AUTH_001": {
		Code:        "APP_AUTH_001",
		Layer:       "Application",
		Category:    "Authentication",
		Description: "Unauthorized access",
		Severity:    "HIGH",
		Retryable:   false,
	},
	"APP_AUTH_002": {
		Code:        "APP_AUTH_002",
		Layer:       "Application",
		Category:    "Authorization",
		Description: "Forbidden access",
		Severity:    "HIGH",
		Retryable:   false,
	},
	"APP_AUTH_003": {
		Code:        "APP_AUTH_003",
		Layer:       "Application",
		Category:    "Authentication",
		Description: "Invalid credentials",
		Severity:    "MEDIUM",
		Retryable:   false,
	},

	// Application Layer - Validation
	"APP_VAL_001": {
		Code:        "APP_VAL_001",
		Layer:       "Application",
		Category:    "Validation",
		Description: "Application-level validation failed",
		Severity:    "MEDIUM",
		Retryable:   false,
	},
	"APP_VAL_002": {
		Code:        "APP_VAL_002",
		Layer:       "Application",
		Category:    "Validation",
		Description: "Invalid request format",
		Severity:    "MEDIUM",
		Retryable:   false,
	},

	// Infrastructure Layer - Database Errors
	"INF_DB_001": {
		Code:        "INF_DB_001",
		Layer:       "Infrastructure",
		Category:    "Database",
		Description: "Database connection error",
		Severity:    "CRITICAL",
		Retryable:   true,
	},
	"INF_DB_002": {
		Code:        "INF_DB_002",
		Layer:       "Infrastructure",
		Category:    "Database",
		Description: "Database query error",
		Severity:    "HIGH",
		Retryable:   false,
	},
	"INF_DB_003": {
		Code:        "INF_DB_003",
		Layer:       "Infrastructure",
		Category:    "Database",
		Description: "Database constraint violation",
		Severity:    "MEDIUM",
		Retryable:   false,
	},
	"INF_DB_004": {
		Code:        "INF_DB_004",
		Layer:       "Infrastructure",
		Category:    "Database",
		Description: "Database query timeout",
		Severity:    "HIGH",
		Retryable:   true,
	},

	// Infrastructure Layer - External API Errors
	"INF_EXT_001": {
		Code:        "INF_EXT_001",
		Layer:       "Infrastructure",
		Category:    "ExternalAPI",
		Description: "External API error",
		Severity:    "HIGH",
		Retryable:   false,
	},
	"INF_EXT_002": {
		Code:        "INF_EXT_002",
		Layer:       "Infrastructure",
		Category:    "ExternalAPI",
		Description: "External API timeout",
		Severity:    "HIGH",
		Retryable:   true,
	},
	"INF_EXT_003": {
		Code:        "INF_EXT_003",
		Layer:       "Infrastructure",
		Category:    "ExternalAPI",
		Description: "External API resource not found",
		Severity:    "MEDIUM",
		Retryable:   false,
	},

	// Infrastructure Layer - Connection Errors
	"INF_CONN_001": {
		Code:        "INF_CONN_001",
		Layer:       "Infrastructure",
		Category:    "Connection",
		Description: "Connection failed",
		Severity:    "CRITICAL",
		Retryable:   true,
	},
	"INF_CONN_002": {
		Code:        "INF_CONN_002",
		Layer:       "Infrastructure",
		Category:    "Connection",
		Description: "Connection timeout",
		Severity:    "HIGH",
		Retryable:   true,
	},

	// Infrastructure Layer - Storage Errors
	"INF_STOR_001": {
		Code:        "INF_STOR_001",
		Layer:       "Infrastructure",
		Category:    "Storage",
		Description: "Storage operation failed",
		Severity:    "HIGH",
		Retryable:   false,
	},
	"INF_STOR_002": {
		Code:        "INF_STOR_002",
		Layer:       "Infrastructure",
		Category:    "Storage",
		Description: "Storage resource not found",
		Severity:    "MEDIUM",
		Retryable:   false,
	},
}

// GetErrorCodeInfo retrieves metadata for an error code
func GetErrorCodeInfo(code string) (ErrorCodeInfo, bool) {
	info, exists := ErrorCodeRegistry[code]
	return info, exists
}
