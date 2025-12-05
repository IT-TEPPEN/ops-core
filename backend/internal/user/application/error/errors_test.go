package error

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNotFoundError_Error(t *testing.T) {
	err := &NotFoundError{
		Code:         CodeResourceNotFound,
		ResourceType: "User",
		ResourceID:   "user-123",
	}

	expected := "[USA0001] User not found (ID: user-123)"
	assert.Equal(t, expected, err.Error())
}

func TestNotFoundError_ErrorWithCause(t *testing.T) {
	cause := errors.New("database error")
	err := &NotFoundError{
		Code:         CodeResourceNotFound,
		ResourceType: "User",
		ResourceID:   "user-123",
		Cause:        cause,
	}

	expected := "[USA0001] User not found (ID: user-123): database error"
	assert.Equal(t, expected, err.Error())
}

func TestNotFoundError_Is(t *testing.T) {
	err := NewNotFoundError("User", "user-123", nil)
	assert.True(t, errors.Is(err, ErrNotFound))
}

func TestNotFoundError_Unwrap(t *testing.T) {
	cause := errors.New("underlying error")
	err := NewNotFoundError("User", "user-123", cause)
	assert.Equal(t, cause, errors.Unwrap(err))
}

func TestUnauthorizedError_Error(t *testing.T) {
	err := &UnauthorizedError{
		Code:   CodeUnauthorized,
		Reason: "invalid token",
	}

	expected := "[USA0003] unauthorized: invalid token"
	assert.Equal(t, expected, err.Error())
}

func TestUnauthorizedError_Is(t *testing.T) {
	err := NewUnauthorizedError("invalid token", nil)
	assert.True(t, errors.Is(err, ErrUnauthorized))
}

func TestForbiddenError_Error(t *testing.T) {
	err := &ForbiddenError{
		Code:     CodeForbidden,
		Resource: "User",
		Action:   "delete",
		UserID:   "user-456",
	}

	expected := "[USA0004] user user-456 is forbidden to delete on resource User"
	assert.Equal(t, expected, err.Error())
}

func TestForbiddenError_Is(t *testing.T) {
	err := NewForbiddenError("User", "delete", "user-456")
	assert.True(t, errors.Is(err, ErrForbidden))
}

func TestConflictError_Error(t *testing.T) {
	err := &ConflictError{
		Code:         CodeResourceConflict,
		ResourceType: "User",
		Identifier:   "test@example.com",
		Reason:       "email already registered",
	}

	expected := "[USA0002] conflict: User with identifier 'test@example.com' already exists: email already registered"
	assert.Equal(t, expected, err.Error())
}

func TestConflictError_Is(t *testing.T) {
	err := NewConflictError("User", "test@example.com", "email already registered", nil)
	assert.True(t, errors.Is(err, ErrConflict))
}

func TestValidationFailedError_Error(t *testing.T) {
	err := &ValidationFailedError{
		Code: CodeValidationFailed,
		Errors: []FieldError{
			{Field: "email", Message: "invalid format"},
		},
	}

	expected := "[USA0006] validation failed: 1 error(s)"
	assert.Equal(t, expected, err.Error())
}

func TestValidationFailedError_Is(t *testing.T) {
	err := NewValidationFailedError([]FieldError{
		{Field: "name", Message: "required"},
	})
	assert.True(t, errors.Is(err, ErrBadRequest))
}
