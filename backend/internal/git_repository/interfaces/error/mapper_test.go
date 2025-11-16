package error

import (
	"errors"
	"net/http"
	"testing"

	apperror "opscore/backend/internal/git_repository/application/error"
	domainerror "opscore/backend/internal/git_repository/domain/error"
	infraerror "opscore/backend/internal/git_repository/infrastructure/error"
)

func TestMapToHTTPError_NotFound(t *testing.T) {
	err := &apperror.NotFoundError{
		Code:         apperror.CodeResourceNotFound,
		ResourceType: "Repository",
		ResourceID:   "repo-123",
	}

	httpErr := MapToHTTPError(err, "req-123")

	if httpErr.StatusCode != http.StatusNotFound {
		t.Errorf("Expected status code %d, got %d", http.StatusNotFound, httpErr.StatusCode)
	}

	if httpErr.Code != "NOT_FOUND" {
		t.Errorf("Expected code 'NOT_FOUND', got '%s'", httpErr.Code)
	}

	if httpErr.RequestID != "req-123" {
		t.Errorf("Expected request ID 'req-123', got '%s'", httpErr.RequestID)
	}

	if httpErr.Details["resource_type"] != "Repository" {
		t.Errorf("Expected resource_type 'Repository', got '%v'", httpErr.Details["resource_type"])
	}
}

func TestMapToHTTPError_ValidationError(t *testing.T) {
	err := &apperror.ValidationFailedError{
		Code: apperror.CodeValidationFailed,
		Errors: []apperror.FieldError{
			{Field: "url", Message: "invalid format", Code: "DOM_VAL_005"},
		},
	}

	httpErr := MapToHTTPError(err, "req-456")

	if httpErr.StatusCode != http.StatusBadRequest {
		t.Errorf("Expected status code %d, got %d", http.StatusBadRequest, httpErr.StatusCode)
	}

	if httpErr.Code != "BAD_REQUEST" {
		t.Errorf("Expected code 'BAD_REQUEST', got '%s'", httpErr.Code)
	}

	validationErrors, ok := httpErr.Details["validation_errors"].([]map[string]interface{})
	if !ok {
		t.Fatal("Expected validation_errors in details")
	}

	if len(validationErrors) != 1 {
		t.Errorf("Expected 1 validation error, got %d", len(validationErrors))
	}
}

func TestMapToHTTPError_ConflictError(t *testing.T) {
	err := &apperror.ConflictError{
		Code:         apperror.CodeResourceConflict,
		ResourceType: "Repository",
		Identifier:   "https://github.com/test/repo",
		Reason:       "URL already exists",
	}

	httpErr := MapToHTTPError(err, "req-789")

	if httpErr.StatusCode != http.StatusConflict {
		t.Errorf("Expected status code %d, got %d", http.StatusConflict, httpErr.StatusCode)
	}

	if httpErr.Code != "CONFLICT" {
		t.Errorf("Expected code 'CONFLICT', got '%s'", httpErr.Code)
	}

	if httpErr.Details["resource_type"] != "Repository" {
		t.Errorf("Expected resource_type 'Repository', got '%v'", httpErr.Details["resource_type"])
	}
}

func TestMapToHTTPError_DomainValidationError(t *testing.T) {
	err := &domainerror.ValidationError{
		Code:    domainerror.CodeInvalidURL,
		Field:   "url",
		Value:   "invalid-url",
		Message: "must be a valid HTTPS URL",
	}

	httpErr := MapToHTTPError(err, "req-111")

	if httpErr.StatusCode != http.StatusBadRequest {
		t.Errorf("Expected status code %d, got %d", http.StatusBadRequest, httpErr.StatusCode)
	}

	if httpErr.Details["field"] != "url" {
		t.Errorf("Expected field 'url', got '%v'", httpErr.Details["field"])
	}
}

func TestMapToHTTPError_DatabaseError(t *testing.T) {
	err := &infraerror.DatabaseError{
		Code:      infraerror.CodeDatabaseQuery,
		Operation: "FindByID",
		Table:     "repositories",
		Cause:     errors.New("query error"),
		Retryable: false,
	}

	httpErr := MapToHTTPError(err, "req-222")

	if httpErr.StatusCode != http.StatusInternalServerError {
		t.Errorf("Expected status code %d, got %d", http.StatusInternalServerError, httpErr.StatusCode)
	}

	if httpErr.Code != "INTERNAL_SERVER_ERROR" {
		t.Errorf("Expected code 'INTERNAL_SERVER_ERROR', got '%s'", httpErr.Code)
	}
}

func TestMapToHTTPError_RetryableError(t *testing.T) {
	err := &infraerror.DatabaseError{
		Code:      infraerror.CodeDatabaseTimeout,
		Operation: "Query",
		Table:     "repositories",
		Cause:     errors.New("timeout"),
		Retryable: true,
	}

	httpErr := MapToHTTPError(err, "req-333")

	if httpErr.StatusCode != http.StatusServiceUnavailable {
		t.Errorf("Expected status code %d, got %d", http.StatusServiceUnavailable, httpErr.StatusCode)
	}

	if httpErr.Code != "SERVICE_UNAVAILABLE" {
		t.Errorf("Expected code 'SERVICE_UNAVAILABLE', got '%s'", httpErr.Code)
	}
}

func TestMapToHTTPError_UnknownError(t *testing.T) {
	err := errors.New("some unknown error")

	httpErr := MapToHTTPError(err, "req-444")

	if httpErr.StatusCode != http.StatusInternalServerError {
		t.Errorf("Expected status code %d, got %d", http.StatusInternalServerError, httpErr.StatusCode)
	}

	if httpErr.Code != "INTERNAL_SERVER_ERROR" {
		t.Errorf("Expected code 'INTERNAL_SERVER_ERROR', got '%s'", httpErr.Code)
	}

	if httpErr.Message != "An unexpected error occurred" {
		t.Errorf("Expected generic error message, got '%s'", httpErr.Message)
	}
}

func TestMapToHTTPError_Unauthorized(t *testing.T) {
	err := &apperror.UnauthorizedError{
		Code:   apperror.CodeUnauthorized,
		Reason: "invalid token",
	}

	httpErr := MapToHTTPError(err, "req-555")

	if httpErr.StatusCode != http.StatusUnauthorized {
		t.Errorf("Expected status code %d, got %d", http.StatusUnauthorized, httpErr.StatusCode)
	}

	if httpErr.Code != "UNAUTHORIZED" {
		t.Errorf("Expected code 'UNAUTHORIZED', got '%s'", httpErr.Code)
	}
}

func TestMapToHTTPError_Forbidden(t *testing.T) {
	err := &apperror.ForbiddenError{
		Code:     apperror.CodeForbidden,
		Resource: "Repository",
		Action:   "delete",
		UserID:   "user-123",
	}

	httpErr := MapToHTTPError(err, "req-666")

	if httpErr.StatusCode != http.StatusForbidden {
		t.Errorf("Expected status code %d, got %d", http.StatusForbidden, httpErr.StatusCode)
	}

	if httpErr.Code != "FORBIDDEN" {
		t.Errorf("Expected code 'FORBIDDEN', got '%s'", httpErr.Code)
	}
}
