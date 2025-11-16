package main

import (
	"fmt"
	"os"
	"sort"

	sharedError "opscore/backend/internal/shared/error"
)

func main() {
	// Create docs directory if it doesn't exist
	docsDir := "docs"
	if err := os.MkdirAll(docsDir, 0755); err != nil {
		fmt.Fprintf(os.Stderr, "Error creating docs directory: %v\n", err)
		os.Exit(1)
	}

	// Open output file
	f, err := os.Create("docs/error-codes.md")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating error-codes.md: %v\n", err)
		os.Exit(1)
	}
	defer f.Close()

	// Write header
	fmt.Fprintln(f, "# Error Code Reference")
	fmt.Fprintln(f, "")
	fmt.Fprintln(f, "This document provides a comprehensive reference for all error codes used in the OpsCore backend application.")
	fmt.Fprintln(f, "")
	fmt.Fprintln(f, "## Error Code Format")
	fmt.Fprintln(f, "")
	fmt.Fprintln(f, "Error codes follow the format: `<LAYER>_<CATEGORY>_<NUMBER>`")
	fmt.Fprintln(f, "")
	fmt.Fprintln(f, "**Layer Prefixes:**")
	fmt.Fprintln(f, "- `DOM`: Domain Layer")
	fmt.Fprintln(f, "- `APP`: Application Layer")
	fmt.Fprintln(f, "- `INF`: Infrastructure Layer")
	fmt.Fprintln(f, "")
	fmt.Fprintln(f, "**Severity Levels:**")
	fmt.Fprintln(f, "- `LOW`: Minor issues that don't significantly impact functionality")
	fmt.Fprintln(f, "- `MEDIUM`: Issues that may impact some functionality")
	fmt.Fprintln(f, "- `HIGH`: Serious issues that impact core functionality")
	fmt.Fprintln(f, "- `CRITICAL`: Critical failures that prevent system operation")
	fmt.Fprintln(f, "")
	fmt.Fprintln(f, "## Error Codes")
	fmt.Fprintln(f, "")
	fmt.Fprintln(f, "| Code | Layer | Category | Description | Severity | Retryable |")
	fmt.Fprintln(f, "|------|-------|----------|-------------|----------|-----------|")

	// Sort error codes for consistent output
	codes := make([]string, 0, len(sharedError.ErrorCodeRegistry))
	for code := range sharedError.ErrorCodeRegistry {
		codes = append(codes, code)
	}
	sort.Strings(codes)

	// Write error codes
	for _, code := range codes {
		info := sharedError.ErrorCodeRegistry[code]
		retryable := "No"
		if info.Retryable {
			retryable = "Yes"
		}
		fmt.Fprintf(f, "| %s | %s | %s | %s | %s | %s |\n",
			info.Code, info.Layer, info.Category, info.Description, info.Severity, retryable)
	}

	// Write usage section
	fmt.Fprintln(f, "")
	fmt.Fprintln(f, "## Usage")
	fmt.Fprintln(f, "")
	fmt.Fprintln(f, "### In Code")
	fmt.Fprintln(f, "")
	fmt.Fprintln(f, "```go")
	fmt.Fprintln(f, "// Domain Layer")
	fmt.Fprintln(f, "return &error.ValidationError{")
	fmt.Fprintln(f, "    Code:    error.CodeInvalidURL,")
	fmt.Fprintln(f, "    Field:   \"url\",")
	fmt.Fprintln(f, "    Message: \"must be a valid HTTPS URL\",")
	fmt.Fprintln(f, "}")
	fmt.Fprintln(f, "")
	fmt.Fprintln(f, "// Application Layer")
	fmt.Fprintln(f, "return &error.NotFoundError{")
	fmt.Fprintln(f, "    Code:         error.CodeResourceNotFound,")
	fmt.Fprintln(f, "    ResourceType: \"Repository\",")
	fmt.Fprintln(f, "    ResourceID:   id,")
	fmt.Fprintln(f, "}")
	fmt.Fprintln(f, "```")
	fmt.Fprintln(f, "")
	fmt.Fprintln(f, "### In Logs")
	fmt.Fprintln(f, "")
	fmt.Fprintln(f, "Error codes are automatically included in log messages:")
	fmt.Fprintln(f, "")
	fmt.Fprintln(f, "```")
	fmt.Fprintln(f, "[DOM_VAL_005] validation failed for field 'url': must be a valid HTTPS URL")
	fmt.Fprintln(f, "```")
	fmt.Fprintln(f, "")
	fmt.Fprintln(f, "### In HTTP Responses")
	fmt.Fprintln(f, "")
	fmt.Fprintln(f, "Error codes are mapped to appropriate HTTP status codes and included in responses:")
	fmt.Fprintln(f, "")
	fmt.Fprintln(f, "```json")
	fmt.Fprintln(f, "{")
	fmt.Fprintln(f, "  \"code\": \"APP_RES_001\",")
	fmt.Fprintln(f, "  \"message\": \"Resource not found\",")
	fmt.Fprintln(f, "  \"details\": {")
	fmt.Fprintln(f, "    \"resource_type\": \"Repository\",")
	fmt.Fprintln(f, "    \"resource_id\": \"repo-123\"")
	fmt.Fprintln(f, "  },")
	fmt.Fprintln(f, "  \"request_id\": \"550e8400-e29b-41d4-a716-446655440000\",")
	fmt.Fprintln(f, "  \"timestamp\": \"2025-11-15T00:10:18Z\"")
	fmt.Fprintln(f, "}")
	fmt.Fprintln(f, "```")
	fmt.Fprintln(f, "")
	fmt.Fprintln(f, "## Related ADRs")
	fmt.Fprintln(f, "")
	fmt.Fprintln(f, "- [ADR 0007: Backend Architecture - Onion Architecture](../adr/0007-backend-architecture-onion.md)")
	fmt.Fprintln(f, "- [ADR 0008: Backend Logging Strategy](../adr/0008-backend-logging-strategy.md)")
	fmt.Fprintln(f, "- [ADR 0015: Backend Custom Error Design](../adr/0015-backend-custom-error-design.md)")

	fmt.Println("✓ Successfully generated docs/error-codes.md")
}
