package error

// ErrorCodeInfo provides metadata about an error code
type ErrorCodeInfo struct {
Code        string
Context     string // Context/Bounded Context (e.g., "git_repository", "user", "project")
Layer       string
Category    string
Description string
Severity    string // "LOW", "MEDIUM", "HIGH", "CRITICAL"
Retryable   bool
}

// ErrorCodeRegistry is the global error code registry
// Format: <CONTEXT>_<LAYER>_<CATEGORY>_<NUMBER>
var ErrorCodeRegistry = map[string]ErrorCodeInfo{
// Git Repository Context - Domain Layer - Validation Errors
"GITREPO_DOM_VAL_001": {
Code:        "GITREPO_DOM_VAL_001",
Context:     "git_repository",
Layer:       "Domain",
Category:    "Validation",
Description: "Invalid entity field value",
Severity:    "MEDIUM",
Retryable:   false,
},
"GITREPO_DOM_VAL_002": {
Code:        "GITREPO_DOM_VAL_002",
Context:     "git_repository",
Layer:       "Domain",
Category:    "Validation",
Description: "Required field is missing",
Severity:    "MEDIUM",
Retryable:   false,
},
"GITREPO_DOM_VAL_003": {
Code:        "GITREPO_DOM_VAL_003",
Context:     "git_repository",
Layer:       "Domain",
Category:    "Validation",
Description: "Invalid field format",
Severity:    "MEDIUM",
Retryable:   false,
},
"GITREPO_DOM_VAL_004": {
Code:        "GITREPO_DOM_VAL_004",
Context:     "git_repository",
Layer:       "Domain",
Category:    "Validation",
Description: "Field value out of range",
Severity:    "MEDIUM",
Retryable:   false,
},
"GITREPO_DOM_VAL_005": {
Code:        "GITREPO_DOM_VAL_005",
Context:     "git_repository",
Layer:       "Domain",
Category:    "Validation",
Description: "Invalid URL format",
Severity:    "MEDIUM",
Retryable:   false,
},
"GITREPO_DOM_VAL_006": {
Code:        "GITREPO_DOM_VAL_006",
Context:     "git_repository",
Layer:       "Domain",
Category:    "Validation",
Description: "Unsupported URL scheme (only HTTPS is supported)",
Severity:    "MEDIUM",
Retryable:   false,
},

// Git Repository Context - Domain Layer - Business Rule Violations
"GITREPO_DOM_BUS_001": {
Code:        "GITREPO_DOM_BUS_001",
Context:     "git_repository",
Layer:       "Domain",
Category:    "Business",
Description: "Business rule violation",
Severity:    "HIGH",
Retryable:   false,
},
"GITREPO_DOM_BUS_002": {
Code:        "GITREPO_DOM_BUS_002",
Context:     "git_repository",
Layer:       "Domain",
Category:    "Business",
Description: "Invalid state transition",
Severity:    "HIGH",
Retryable:   false,
},
"GITREPO_DOM_BUS_003": {
Code:        "GITREPO_DOM_BUS_003",
Context:     "git_repository",
Layer:       "Domain",
Category:    "Business",
Description: "Invariant violation",
Severity:    "HIGH",
Retryable:   false,
},

// Git Repository Context - Application Layer - Resource Errors
"GITREPO_APP_RES_001": {
Code:        "GITREPO_APP_RES_001",
Context:     "git_repository",
Layer:       "Application",
Category:    "Resource",
Description: "Requested resource not found",
Severity:    "LOW",
Retryable:   false,
},
"GITREPO_APP_RES_002": {
Code:        "GITREPO_APP_RES_002",
Context:     "git_repository",
Layer:       "Application",
Category:    "Resource",
Description: "Resource conflict (duplicate)",
Severity:    "MEDIUM",
Retryable:   false,
},

// Git Repository Context - Application Layer - Authentication/Authorization
"GITREPO_APP_AUTH_001": {
Code:        "GITREPO_APP_AUTH_001",
Context:     "git_repository",
Layer:       "Application",
Category:    "Authentication",
Description: "Unauthorized access",
Severity:    "HIGH",
Retryable:   false,
},
"GITREPO_APP_AUTH_002": {
Code:        "GITREPO_APP_AUTH_002",
Context:     "git_repository",
Layer:       "Application",
Category:    "Authorization",
Description: "Forbidden access",
Severity:    "HIGH",
Retryable:   false,
},
"GITREPO_APP_AUTH_003": {
Code:        "GITREPO_APP_AUTH_003",
Context:     "git_repository",
Layer:       "Application",
Category:    "Authentication",
Description: "Invalid credentials",
Severity:    "MEDIUM",
Retryable:   false,
},

// Git Repository Context - Application Layer - Validation
"GITREPO_APP_VAL_001": {
Code:        "GITREPO_APP_VAL_001",
Context:     "git_repository",
Layer:       "Application",
Category:    "Validation",
Description: "Application-level validation failed",
Severity:    "MEDIUM",
Retryable:   false,
},
"GITREPO_APP_VAL_002": {
Code:        "GITREPO_APP_VAL_002",
Context:     "git_repository",
Layer:       "Application",
Category:    "Validation",
Description: "Invalid request format",
Severity:    "MEDIUM",
Retryable:   false,
},

// Git Repository Context - Infrastructure Layer - Database Errors
"GITREPO_INF_DB_001": {
Code:        "GITREPO_INF_DB_001",
Context:     "git_repository",
Layer:       "Infrastructure",
Category:    "Database",
Description: "Database connection error",
Severity:    "CRITICAL",
Retryable:   true,
},
"GITREPO_INF_DB_002": {
Code:        "GITREPO_INF_DB_002",
Context:     "git_repository",
Layer:       "Infrastructure",
Category:    "Database",
Description: "Database query error",
Severity:    "HIGH",
Retryable:   false,
},
"GITREPO_INF_DB_003": {
Code:        "GITREPO_INF_DB_003",
Context:     "git_repository",
Layer:       "Infrastructure",
Category:    "Database",
Description: "Database constraint violation",
Severity:    "MEDIUM",
Retryable:   false,
},
"GITREPO_INF_DB_004": {
Code:        "GITREPO_INF_DB_004",
Context:     "git_repository",
Layer:       "Infrastructure",
Category:    "Database",
Description: "Database query timeout",
Severity:    "HIGH",
Retryable:   true,
},

// Git Repository Context - Infrastructure Layer - External API Errors
"GITREPO_INF_EXT_001": {
Code:        "GITREPO_INF_EXT_001",
Context:     "git_repository",
Layer:       "Infrastructure",
Category:    "ExternalAPI",
Description: "External API error",
Severity:    "HIGH",
Retryable:   false,
},
"GITREPO_INF_EXT_002": {
Code:        "GITREPO_INF_EXT_002",
Context:     "git_repository",
Layer:       "Infrastructure",
Category:    "ExternalAPI",
Description: "External API timeout",
Severity:    "HIGH",
Retryable:   true,
},
"GITREPO_INF_EXT_003": {
Code:        "GITREPO_INF_EXT_003",
Context:     "git_repository",
Layer:       "Infrastructure",
Category:    "ExternalAPI",
Description: "External API resource not found",
Severity:    "MEDIUM",
Retryable:   false,
},

// Git Repository Context - Infrastructure Layer - Connection Errors
"GITREPO_INF_CONN_001": {
Code:        "GITREPO_INF_CONN_001",
Context:     "git_repository",
Layer:       "Infrastructure",
Category:    "Connection",
Description: "Connection failed",
Severity:    "CRITICAL",
Retryable:   true,
},
"GITREPO_INF_CONN_002": {
Code:        "GITREPO_INF_CONN_002",
Context:     "git_repository",
Layer:       "Infrastructure",
Category:    "Connection",
Description: "Connection timeout",
Severity:    "HIGH",
Retryable:   true,
},

// Git Repository Context - Infrastructure Layer - Storage Errors
"GITREPO_INF_STOR_001": {
Code:        "GITREPO_INF_STOR_001",
Context:     "git_repository",
Layer:       "Infrastructure",
Category:    "Storage",
Description: "Storage operation failed",
Severity:    "HIGH",
Retryable:   false,
},
"GITREPO_INF_STOR_002": {
Code:        "GITREPO_INF_STOR_002",
Context:     "git_repository",
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
