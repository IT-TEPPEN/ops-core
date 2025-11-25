package error

// ErrorCode represents a unique error identifier
// Format: <CONTEXT><LAYER><NUMBER>
// - CONTEXT: 2 characters (e.g., US for user)
// - LAYER: 1 character (D for domain, A for application, I for infrastructure)
// - NUMBER: 4 digits (e.g., 0001)
// Example: USD0001
type ErrorCode string

const (
	// Validation errors
	CodeInvalidEntityField   ErrorCode = "USD0001"
	CodeRequiredFieldMissing ErrorCode = "USD0002"
	CodeInvalidFieldFormat   ErrorCode = "USD0003"
	CodeFieldValueOutOfRange ErrorCode = "USD0004"
	CodeInvalidEmail         ErrorCode = "USD0005"
	CodeInvalidRole          ErrorCode = "USD0006"

	// Business rule violations
	CodeBusinessRuleViolation  ErrorCode = "USD0007"
	CodeInvalidStateTransition ErrorCode = "USD0008"
	CodeInvariantViolation     ErrorCode = "USD0009"
	CodeDuplicateMember        ErrorCode = "USD0010"
	CodeMemberNotFound         ErrorCode = "USD0011"
)

// String returns the string representation of the error code
func (c ErrorCode) String() string {
	return string(c)
}
