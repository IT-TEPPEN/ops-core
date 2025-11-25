package error

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidationError_Error(t *testing.T) {
	err := &ValidationError{
		Code:    CodeInvalidEmail,
		Field:   "email",
		Value:   "invalid-email",
		Message: "must be a valid email address",
	}

	expected := "[USD0005] validation failed for field 'email': must be a valid email address"
	assert.Equal(t, expected, err.Error())
}

func TestValidationError_Is(t *testing.T) {
	err := &ValidationError{Field: "name", Message: "required"}
	assert.True(t, errors.Is(err, ErrInvalidEntity))
}

func TestValidationError_ErrorCode(t *testing.T) {
	err := &ValidationError{
		Code:    CodeInvalidEmail,
		Field:   "email",
		Message: "invalid format",
	}

	assert.Equal(t, CodeInvalidEmail, err.ErrorCode())
}

func TestBusinessRuleViolationError_Error(t *testing.T) {
	err := &BusinessRuleViolationError{
		Code:    CodeDuplicateMember,
		Rule:    "unique_member",
		Entity:  "Group",
		Message: "user is already a member of this group",
	}

	expected := "[USD0010] business rule 'unique_member' violated for Group: user is already a member of this group"
	assert.Equal(t, expected, err.Error())
}

func TestBusinessRuleViolationError_Is(t *testing.T) {
	err := &BusinessRuleViolationError{
		Code:    CodeDuplicateMember,
		Rule:    "unique_member",
		Entity:  "Group",
		Message: "user is already a member",
	}
	assert.True(t, errors.Is(err, ErrInvariantViolation))
}

func TestInvalidStateTransitionError_Error(t *testing.T) {
	err := &InvalidStateTransitionError{
		Code:      CodeInvalidStateTransition,
		Entity:    "User",
		FromState: "active",
		ToState:   "deleted",
		Reason:    "cannot delete active user",
	}

	expected := "[USD0008] invalid state transition for User from 'active' to 'deleted': cannot delete active user"
	assert.Equal(t, expected, err.Error())
}

func TestInvalidStateTransitionError_Is(t *testing.T) {
	err := &InvalidStateTransitionError{
		Code:      CodeInvalidStateTransition,
		Entity:    "User",
		FromState: "active",
		ToState:   "deleted",
		Reason:    "cannot delete active user",
	}
	assert.True(t, errors.Is(err, ErrInvariantViolation))
}

func TestNewValidationError(t *testing.T) {
	err := NewValidationError("name", "", "cannot be empty")
	assert.Equal(t, CodeRequiredFieldMissing, err.Code)
	assert.Equal(t, "name", err.Field)
}

func TestNewEmailValidationError(t *testing.T) {
	err := NewEmailValidationError("email", "invalid", "invalid format")
	assert.Equal(t, CodeInvalidEmail, err.Code)
	assert.Equal(t, "email", err.Field)
}

func TestNewRoleValidationError(t *testing.T) {
	err := NewRoleValidationError("role", "superadmin", "invalid role")
	assert.Equal(t, CodeInvalidRole, err.Code)
	assert.Equal(t, "role", err.Field)
}

func TestNewDuplicateMemberError(t *testing.T) {
	err := NewDuplicateMemberError("Group", "user is already a member")
	assert.Equal(t, CodeDuplicateMember, err.Code)
	assert.Equal(t, "unique_member", err.Rule)
}

func TestNewMemberNotFoundError(t *testing.T) {
	err := NewMemberNotFoundError("Group", "user is not a member")
	assert.Equal(t, CodeMemberNotFound, err.Code)
	assert.Equal(t, "member_exists", err.Rule)
}
