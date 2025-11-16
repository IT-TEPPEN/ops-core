package error

// ErrorCode represents a unique error identifier
// Format: <CONTEXT><LAYER><NUMBER>
// - CONTEXT: 2 characters (e.g., GR for git_repository)
// - LAYER: 1 character (D for domain, A for application, I for infrastructure)
// - NUMBER: 4 digits (e.g., 0001)
// Example: GRD0001
type ErrorCode string

const (
	// Validation errors
	CodeInvalidEntityField   ErrorCode = "GRD0001"
	CodeRequiredFieldMissing ErrorCode = "GRD0002"
	CodeInvalidFieldFormat   ErrorCode = "GRD0003"
	CodeFieldValueOutOfRange ErrorCode = "GRD0004"
	CodeInvalidURL           ErrorCode = "GRD0005"
	CodeUnsupportedURLScheme ErrorCode = "GRD0006"

	// Business rule violations
	CodeBusinessRuleViolation  ErrorCode = "GRD0007"
	CodeInvalidStateTransition ErrorCode = "GRD0008"
	CodeInvariantViolation     ErrorCode = "GRD0009"
)

// String returns the string representation of the error code
func (c ErrorCode) String() string {
	return string(c)
}
