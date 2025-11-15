package error

// ErrorCode represents a unique error identifier
type ErrorCode string

const (
	// Validation errors (DOM_VAL_xxx)
	CodeInvalidEntityField   ErrorCode = "DOM_VAL_001"
	CodeRequiredFieldMissing ErrorCode = "DOM_VAL_002"
	CodeInvalidFieldFormat   ErrorCode = "DOM_VAL_003"
	CodeFieldValueOutOfRange ErrorCode = "DOM_VAL_004"
	CodeInvalidURL           ErrorCode = "DOM_VAL_005"
	CodeUnsupportedURLScheme ErrorCode = "DOM_VAL_006"

	// Business rule violations (DOM_BUS_xxx)
	CodeBusinessRuleViolation  ErrorCode = "DOM_BUS_001"
	CodeInvalidStateTransition ErrorCode = "DOM_BUS_002"
	CodeInvariantViolation     ErrorCode = "DOM_BUS_003"
)

// String returns the string representation of the error code
func (c ErrorCode) String() string {
	return string(c)
}
