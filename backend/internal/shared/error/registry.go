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
// Format: <CONTEXT><LAYER><NUMBER>
// - CONTEXT: 2 characters (e.g., GR for git_repository)
// - LAYER: 1 character (D for domain, A for application, I for infrastructure)
// - NUMBER: 4 digits (e.g., 0001)
var ErrorCodeRegistry = map[string]ErrorCodeInfo{
// Git Repository Context - Domain Layer
"GRD0001": {
Code:        "GRD0001",
Context:     "git_repository",
Layer:       "Domain",
Category:    "Validation",
Description: "Invalid entity field value",
Severity:    "MEDIUM",
Retryable:   false,
},
"GRD0002": {
Code:        "GRD0002",
Context:     "git_repository",
Layer:       "Domain",
Category:    "Validation",
Description: "Required field is missing",
Severity:    "MEDIUM",
Retryable:   false,
},
"GRD0003": {
Code:        "GRD0003",
Context:     "git_repository",
Layer:       "Domain",
Category:    "Validation",
Description: "Invalid field format",
Severity:    "MEDIUM",
Retryable:   false,
},
"GRD0004": {
Code:        "GRD0004",
Context:     "git_repository",
Layer:       "Domain",
Category:    "Validation",
Description: "Field value out of range",
Severity:    "MEDIUM",
Retryable:   false,
},
"GRD0005": {
Code:        "GRD0005",
Context:     "git_repository",
Layer:       "Domain",
Category:    "Validation",
Description: "Invalid URL format",
Severity:    "MEDIUM",
Retryable:   false,
},
"GRD0006": {
Code:        "GRD0006",
Context:     "git_repository",
Layer:       "Domain",
Category:    "Validation",
Description: "Unsupported URL scheme (only HTTPS is supported)",
Severity:    "MEDIUM",
Retryable:   false,
},
"GRD0007": {
Code:        "GRD0007",
Context:     "git_repository",
Layer:       "Domain",
Category:    "Business",
Description: "Business rule violation",
Severity:    "HIGH",
Retryable:   false,
},
"GRD0008": {
Code:        "GRD0008",
Context:     "git_repository",
Layer:       "Domain",
Category:    "Business",
Description: "Invalid state transition",
Severity:    "HIGH",
Retryable:   false,
},
"GRD0009": {
Code:        "GRD0009",
Context:     "git_repository",
Layer:       "Domain",
Category:    "Business",
Description: "Invariant violation",
Severity:    "HIGH",
Retryable:   false,
},

// Git Repository Context - Application Layer
"GRA0001": {
Code:        "GRA0001",
Context:     "git_repository",
Layer:       "Application",
Category:    "Resource",
Description: "Requested resource not found",
Severity:    "LOW",
Retryable:   false,
},
"GRA0002": {
Code:        "GRA0002",
Context:     "git_repository",
Layer:       "Application",
Category:    "Resource",
Description: "Resource conflict (duplicate)",
Severity:    "MEDIUM",
Retryable:   false,
},
"GRA0003": {
Code:        "GRA0003",
Context:     "git_repository",
Layer:       "Application",
Category:    "Authentication",
Description: "Unauthorized access",
Severity:    "HIGH",
Retryable:   false,
},
"GRA0004": {
Code:        "GRA0004",
Context:     "git_repository",
Layer:       "Application",
Category:    "Authorization",
Description: "Forbidden access",
Severity:    "HIGH",
Retryable:   false,
},
"GRA0005": {
Code:        "GRA0005",
Context:     "git_repository",
Layer:       "Application",
Category:    "Authentication",
Description: "Invalid credentials",
Severity:    "MEDIUM",
Retryable:   false,
},
"GRA0006": {
Code:        "GRA0006",
Context:     "git_repository",
Layer:       "Application",
Category:    "Validation",
Description: "Application-level validation failed",
Severity:    "MEDIUM",
Retryable:   false,
},
"GRA0007": {
Code:        "GRA0007",
Context:     "git_repository",
Layer:       "Application",
Category:    "Validation",
Description: "Invalid request format",
Severity:    "MEDIUM",
Retryable:   false,
},

// Git Repository Context - Infrastructure Layer
"GRI0001": {
Code:        "GRI0001",
Context:     "git_repository",
Layer:       "Infrastructure",
Category:    "Database",
Description: "Database connection error",
Severity:    "CRITICAL",
Retryable:   true,
},
"GRI0002": {
Code:        "GRI0002",
Context:     "git_repository",
Layer:       "Infrastructure",
Category:    "Database",
Description: "Database query error",
Severity:    "HIGH",
Retryable:   false,
},
"GRI0003": {
Code:        "GRI0003",
Context:     "git_repository",
Layer:       "Infrastructure",
Category:    "Database",
Description: "Database constraint violation",
Severity:    "MEDIUM",
Retryable:   false,
},
"GRI0004": {
Code:        "GRI0004",
Context:     "git_repository",
Layer:       "Infrastructure",
Category:    "Database",
Description: "Database query timeout",
Severity:    "HIGH",
Retryable:   true,
},
"GRI0005": {
Code:        "GRI0005",
Context:     "git_repository",
Layer:       "Infrastructure",
Category:    "ExternalAPI",
Description: "External API error",
Severity:    "HIGH",
Retryable:   false,
},
"GRI0006": {
Code:        "GRI0006",
Context:     "git_repository",
Layer:       "Infrastructure",
Category:    "ExternalAPI",
Description: "External API timeout",
Severity:    "HIGH",
Retryable:   true,
},
"GRI0007": {
Code:        "GRI0007",
Context:     "git_repository",
Layer:       "Infrastructure",
Category:    "ExternalAPI",
Description: "External API resource not found",
Severity:    "MEDIUM",
Retryable:   false,
},
"GRI0008": {
Code:        "GRI0008",
Context:     "git_repository",
Layer:       "Infrastructure",
Category:    "Connection",
Description: "Connection failed",
Severity:    "CRITICAL",
Retryable:   true,
},
"GRI0009": {
Code:        "GRI0009",
Context:     "git_repository",
Layer:       "Infrastructure",
Category:    "Connection",
Description: "Connection timeout",
Severity:    "HIGH",
Retryable:   true,
},
"GRI0010": {
Code:        "GRI0010",
Context:     "git_repository",
Layer:       "Infrastructure",
Category:    "Storage",
Description: "Storage operation failed",
Severity:    "HIGH",
Retryable:   false,
},
"GRI0011": {
Code:        "GRI0011",
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
