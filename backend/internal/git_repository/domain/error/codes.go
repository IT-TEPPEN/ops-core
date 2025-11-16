package error

// ErrorCode represents a unique error identifier
// Format: <CONTEXT>_<LAYER>_<CATEGORY>_<NUMBER>
// Example: GITREPO_DOM_VAL_001
type ErrorCode string

const (
	// Validation errors (GITREPO_DOM_VAL_xxx)
	CodeInvalidEntityField   ErrorCode = "GITREPO_DOM_VAL_001"
	CodeRequiredFieldMissing ErrorCode = "GITREPO_DOM_VAL_002"
	CodeInvalidFieldFormat   ErrorCode = "GITREPO_DOM_VAL_003"
	CodeFieldValueOutOfRange ErrorCode = "GITREPO_DOM_VAL_004"
	CodeInvalidURL           ErrorCode = "GITREPO_DOM_VAL_005"
	CodeUnsupportedURLScheme ErrorCode = "GITREPO_DOM_VAL_006"

	// Business rule violations (GITREPO_DOM_BUS_xxx)
	CodeBusinessRuleViolation  ErrorCode = "GITREPO_DOM_BUS_001"
	CodeInvalidStateTransition ErrorCode = "GITREPO_DOM_BUS_002"
	CodeInvariantViolation     ErrorCode = "GITREPO_DOM_BUS_003"
)

// String returns the string representation of the error code
func (c ErrorCode) String() string {
	return string(c)
}
