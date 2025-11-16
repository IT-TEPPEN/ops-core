package error

import (
	"errors"
	"testing"
)

func TestNotFoundError_Error(t *testing.T) {
	err := &NotFoundError{
		Code:         CodeResourceNotFound,
		ResourceType: "Repository",
		ResourceID:   "repo-123",
	}

	expected := "[GITREPO_APP_RES_001] Repository not found (ID: repo-123)"
	if err.Error() != expected {
		t.Errorf("Expected error message '%s', got '%s'", expected, err.Error())
	}
}

func TestNotFoundError_ErrorWithCause(t *testing.T) {
	cause := errors.New("database error")
	err := &NotFoundError{
		Code:         CodeResourceNotFound,
		ResourceType: "Repository",
		ResourceID:   "repo-123",
		Cause:        cause,
	}

	if !errors.Is(err, cause) {
		t.Error("Expected error to wrap cause")
	}
}

func TestNotFoundError_Is(t *testing.T) {
	err := &NotFoundError{
		Code:         CodeResourceNotFound,
		ResourceType: "Repository",
		ResourceID:   "repo-123",
	}

	if !errors.Is(err, ErrNotFound) {
		t.Error("NotFoundError should match ErrNotFound sentinel")
	}
}

func TestNotFoundError_Unwrap(t *testing.T) {
	cause := errors.New("underlying error")
	err := &NotFoundError{
		Code:         CodeResourceNotFound,
		ResourceType: "Repository",
		ResourceID:   "repo-123",
		Cause:        cause,
	}

	if err.Unwrap() != cause {
		t.Error("Unwrap should return the cause")
	}
}

func TestNotFoundError_ErrorCode(t *testing.T) {
	err := &NotFoundError{
		Code:         CodeResourceNotFound,
		ResourceType: "Repository",
		ResourceID:   "repo-123",
	}

	if err.ErrorCode() != CodeResourceNotFound {
		t.Errorf("Expected error code '%s', got '%s'", CodeResourceNotFound, err.ErrorCode())
	}
}

func TestUnauthorizedError_Error(t *testing.T) {
	err := &UnauthorizedError{
		Code:   CodeUnauthorized,
		Reason: "invalid token",
	}

	expected := "[GITREPO_APP_AUTH_001] unauthorized: invalid token"
	if err.Error() != expected {
		t.Errorf("Expected error message '%s', got '%s'", expected, err.Error())
	}
}

func TestUnauthorizedError_Is(t *testing.T) {
	err := &UnauthorizedError{
		Code:   CodeUnauthorized,
		Reason: "test",
	}

	if !errors.Is(err, ErrUnauthorized) {
		t.Error("UnauthorizedError should match ErrUnauthorized sentinel")
	}
}

func TestForbiddenError_Error(t *testing.T) {
	err := &ForbiddenError{
		Code:     CodeForbidden,
		Resource: "Repository",
		Action:   "delete",
		UserID:   "user-123",
	}

	expected := "[GITREPO_APP_AUTH_002] user user-123 is forbidden to delete on resource Repository"
	if err.Error() != expected {
		t.Errorf("Expected error message '%s', got '%s'", expected, err.Error())
	}
}

func TestForbiddenError_Is(t *testing.T) {
	err := &ForbiddenError{
		Code:     CodeForbidden,
		Resource: "Repository",
		Action:   "delete",
		UserID:   "user-123",
	}

	if !errors.Is(err, ErrForbidden) {
		t.Error("ForbiddenError should match ErrForbidden sentinel")
	}
}

func TestConflictError_Error(t *testing.T) {
	err := &ConflictError{
		Code:         CodeResourceConflict,
		ResourceType: "Repository",
		Identifier:   "https://github.com/test/repo",
		Reason:       "URL already registered",
	}

	expected := "[GITREPO_APP_RES_002] conflict: Repository with identifier 'https://github.com/test/repo' already exists: URL already registered"
	if err.Error() != expected {
		t.Errorf("Expected error message '%s', got '%s'", expected, err.Error())
	}
}

func TestConflictError_Is(t *testing.T) {
	err := &ConflictError{
		Code:         CodeResourceConflict,
		ResourceType: "Repository",
		Identifier:   "test",
		Reason:       "test",
	}

	if !errors.Is(err, ErrConflict) {
		t.Error("ConflictError should match ErrConflict sentinel")
	}
}

func TestValidationFailedError_Error(t *testing.T) {
	err := &ValidationFailedError{
		Code: CodeValidationFailed,
		Errors: []FieldError{
			{Field: "name", Message: "required", Code: "DOM_VAL_002"},
			{Field: "url", Message: "invalid format", Code: "DOM_VAL_005"},
		},
	}

	expected := "[GITREPO_APP_VAL_001] validation failed: 2 error(s)"
	if err.Error() != expected {
		t.Errorf("Expected error message '%s', got '%s'", expected, err.Error())
	}
}

func TestValidationFailedError_Is(t *testing.T) {
	err := &ValidationFailedError{
		Code: CodeValidationFailed,
		Errors: []FieldError{
			{Field: "test", Message: "test"},
		},
	}

	if !errors.Is(err, ErrBadRequest) {
		t.Error("ValidationFailedError should match ErrBadRequest sentinel")
	}
}

func TestValidationFailedError_ErrorCode(t *testing.T) {
	err := &ValidationFailedError{
		Code: CodeValidationFailed,
		Errors: []FieldError{
			{Field: "test", Message: "test"},
		},
	}

	if err.ErrorCode() != CodeValidationFailed {
		t.Errorf("Expected error code '%s', got '%s'", CodeValidationFailed, err.ErrorCode())
	}
}
