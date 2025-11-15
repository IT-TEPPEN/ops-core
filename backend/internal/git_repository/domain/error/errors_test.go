package error

import (
	"errors"
	"testing"
)

func TestValidationError_Error(t *testing.T) {
	err := &ValidationError{
		Code:    CodeInvalidFieldFormat,
		Field:   "email",
		Value:   "invalid-email",
		Message: "must be a valid email address",
	}

	expected := "[DOM_VAL_003] validation failed for field 'email': must be a valid email address"
	if err.Error() != expected {
		t.Errorf("Expected error message '%s', got '%s'", expected, err.Error())
	}
}

func TestValidationError_Is(t *testing.T) {
	err := &ValidationError{
		Code:    CodeRequiredFieldMissing,
		Field:   "name",
		Message: "required",
	}

	if !errors.Is(err, ErrInvalidEntity) {
		t.Error("ValidationError should match ErrInvalidEntity sentinel")
	}
}

func TestValidationError_ErrorCode(t *testing.T) {
	err := &ValidationError{
		Code:    CodeInvalidURL,
		Field:   "url",
		Message: "invalid format",
	}

	if err.ErrorCode() != CodeInvalidURL {
		t.Errorf("Expected error code '%s', got '%s'", CodeInvalidURL, err.ErrorCode())
	}
}

func TestBusinessRuleViolationError_Error(t *testing.T) {
	err := &BusinessRuleViolationError{
		Code:    CodeBusinessRuleViolation,
		Rule:    "UniqueURL",
		Entity:  "Repository",
		Message: "repository with this URL already exists",
	}

	expected := "[DOM_BUS_001] business rule 'UniqueURL' violated for Repository: repository with this URL already exists"
	if err.Error() != expected {
		t.Errorf("Expected error message '%s', got '%s'", expected, err.Error())
	}
}

func TestBusinessRuleViolationError_Is(t *testing.T) {
	err := &BusinessRuleViolationError{
		Code:    CodeBusinessRuleViolation,
		Rule:    "TestRule",
		Entity:  "TestEntity",
		Message: "test",
	}

	if !errors.Is(err, ErrInvariantViolation) {
		t.Error("BusinessRuleViolationError should match ErrInvariantViolation sentinel")
	}
}

func TestBusinessRuleViolationError_ErrorCode(t *testing.T) {
	err := &BusinessRuleViolationError{
		Code:    CodeBusinessRuleViolation,
		Rule:    "TestRule",
		Entity:  "TestEntity",
		Message: "test",
	}

	if err.ErrorCode() != CodeBusinessRuleViolation {
		t.Errorf("Expected error code '%s', got '%s'", CodeBusinessRuleViolation, err.ErrorCode())
	}
}

func TestInvalidStateTransitionError_Error(t *testing.T) {
	err := &InvalidStateTransitionError{
		Code:      CodeInvalidStateTransition,
		Entity:    "Repository",
		FromState: "Active",
		ToState:   "Deleted",
		Reason:    "cannot delete active repository with pending operations",
	}

	expected := "[DOM_BUS_002] invalid state transition for Repository from 'Active' to 'Deleted': cannot delete active repository with pending operations"
	if err.Error() != expected {
		t.Errorf("Expected error message '%s', got '%s'", expected, err.Error())
	}
}

func TestInvalidStateTransitionError_Is(t *testing.T) {
	err := &InvalidStateTransitionError{
		Code:      CodeInvalidStateTransition,
		Entity:    "TestEntity",
		FromState: "A",
		ToState:   "B",
		Reason:    "test",
	}

	if !errors.Is(err, ErrInvariantViolation) {
		t.Error("InvalidStateTransitionError should match ErrInvariantViolation sentinel")
	}
}

func TestInvalidStateTransitionError_ErrorCode(t *testing.T) {
	err := &InvalidStateTransitionError{
		Code:      CodeInvalidStateTransition,
		Entity:    "TestEntity",
		FromState: "A",
		ToState:   "B",
		Reason:    "test",
	}

	if err.ErrorCode() != CodeInvalidStateTransition {
		t.Errorf("Expected error code '%s', got '%s'", CodeInvalidStateTransition, err.ErrorCode())
	}
}
