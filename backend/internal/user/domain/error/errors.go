package error

import (
	"errors"
	"fmt"
)

// Sentinel errors for common domain violations
// These are used for error comparison with errors.Is() across layer boundaries
var (
	ErrInvalidEntity      = errors.New("invalid entity")
	ErrInvariantViolation = errors.New("invariant violation")
)

// ValidationError represents a domain validation failure
type ValidationError struct {
	Code    ErrorCode
	Field   string
	Value   interface{}
	Message string
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("[%s] validation failed for field '%s': %s", e.Code, e.Field, e.Message)
}

// Is enables errors.Is() comparison with sentinel error
func (e *ValidationError) Is(target error) bool {
	return target == ErrInvalidEntity
}

// ErrorCode returns the error code
func (e *ValidationError) ErrorCode() ErrorCode {
	return e.Code
}

// BusinessRuleViolationError represents a business rule violation
type BusinessRuleViolationError struct {
	Code    ErrorCode
	Rule    string
	Entity  string
	Message string
}

func (e *BusinessRuleViolationError) Error() string {
	return fmt.Sprintf("[%s] business rule '%s' violated for %s: %s", e.Code, e.Rule, e.Entity, e.Message)
}

func (e *BusinessRuleViolationError) Is(target error) bool {
	return target == ErrInvariantViolation
}

func (e *BusinessRuleViolationError) ErrorCode() ErrorCode {
	return e.Code
}

// InvalidStateTransitionError represents an invalid state transition
type InvalidStateTransitionError struct {
	Code      ErrorCode
	Entity    string
	FromState string
	ToState   string
	Reason    string
}

func (e *InvalidStateTransitionError) Error() string {
	return fmt.Sprintf("[%s] invalid state transition for %s from '%s' to '%s': %s",
		e.Code, e.Entity, e.FromState, e.ToState, e.Reason)
}

func (e *InvalidStateTransitionError) Is(target error) bool {
	return target == ErrInvariantViolation
}

func (e *InvalidStateTransitionError) ErrorCode() ErrorCode {
	return e.Code
}

// Factory functions to prevent manual error struct/code pairing errors

// NewValidationError creates a new ValidationError with the appropriate error code
func NewValidationError(field string, value interface{}, message string) *ValidationError {
	code := CodeInvalidFieldFormat // Default
	if message == "required" || message == "cannot be empty" {
		code = CodeRequiredFieldMissing
	}

	return &ValidationError{
		Code:    code,
		Field:   field,
		Value:   value,
		Message: message,
	}
}

// NewEmailValidationError creates a new ValidationError specifically for email validation
func NewEmailValidationError(field string, value interface{}, message string) *ValidationError {
	return &ValidationError{
		Code:    CodeInvalidEmail,
		Field:   field,
		Value:   value,
		Message: message,
	}
}

// NewRoleValidationError creates a new ValidationError specifically for role validation
func NewRoleValidationError(field string, value interface{}, message string) *ValidationError {
	return &ValidationError{
		Code:    CodeInvalidRole,
		Field:   field,
		Value:   value,
		Message: message,
	}
}

// NewBusinessRuleViolationError creates a new BusinessRuleViolationError with the correct error code
func NewBusinessRuleViolationError(rule, entity, message string) *BusinessRuleViolationError {
	return &BusinessRuleViolationError{
		Code:    CodeBusinessRuleViolation,
		Rule:    rule,
		Entity:  entity,
		Message: message,
	}
}

// NewDuplicateMemberError creates a new BusinessRuleViolationError for duplicate member
func NewDuplicateMemberError(entity, message string) *BusinessRuleViolationError {
	return &BusinessRuleViolationError{
		Code:    CodeDuplicateMember,
		Rule:    "unique_member",
		Entity:  entity,
		Message: message,
	}
}

// NewMemberNotFoundError creates a new BusinessRuleViolationError for member not found
func NewMemberNotFoundError(entity, message string) *BusinessRuleViolationError {
	return &BusinessRuleViolationError{
		Code:    CodeMemberNotFound,
		Rule:    "member_exists",
		Entity:  entity,
		Message: message,
	}
}

// NewInvalidStateTransitionError creates a new InvalidStateTransitionError with the correct error code
func NewInvalidStateTransitionError(entity, fromState, toState, reason string) *InvalidStateTransitionError {
	return &InvalidStateTransitionError{
		Code:      CodeInvalidStateTransition,
		Entity:    entity,
		FromState: fromState,
		ToState:   toState,
		Reason:    reason,
	}
}
