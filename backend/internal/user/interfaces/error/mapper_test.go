package error

import (
	"errors"
	"testing"

	apperror "opscore/backend/internal/user/application/error"
	domainerror "opscore/backend/internal/user/domain/error"
	infraerror "opscore/backend/internal/user/infrastructure/error"

	"github.com/stretchr/testify/assert"
)

func TestMapToHTTPError_NotFound(t *testing.T) {
	err := apperror.NewNotFoundError("User", "user-123", nil)

	httpErr := MapToHTTPError(err, "req-123")

	assert.Equal(t, 404, httpErr.StatusCode)
	assert.Equal(t, "NOT_FOUND", httpErr.Code)
	assert.Equal(t, "req-123", httpErr.RequestID)
	assert.Equal(t, "User", httpErr.Details["resource_type"])
}

func TestMapToHTTPError_Conflict(t *testing.T) {
	err := apperror.NewConflictError("User", "test@example.com", "email already registered", nil)

	httpErr := MapToHTTPError(err, "req-123")

	assert.Equal(t, 409, httpErr.StatusCode)
	assert.Equal(t, "CONFLICT", httpErr.Code)
	assert.Equal(t, "User", httpErr.Details["resource_type"])
}

func TestMapToHTTPError_BadRequest(t *testing.T) {
	err := apperror.NewValidationFailedError([]apperror.FieldError{
		{Field: "email", Message: "invalid format"},
	})

	httpErr := MapToHTTPError(err, "req-123")

	assert.Equal(t, 400, httpErr.StatusCode)
	assert.Equal(t, "BAD_REQUEST", httpErr.Code)
	assert.NotNil(t, httpErr.Details["validation_errors"])
}

func TestMapToHTTPError_Unauthorized(t *testing.T) {
	err := apperror.NewUnauthorizedError("invalid token", nil)

	httpErr := MapToHTTPError(err, "req-123")

	assert.Equal(t, 401, httpErr.StatusCode)
	assert.Equal(t, "UNAUTHORIZED", httpErr.Code)
}

func TestMapToHTTPError_Forbidden(t *testing.T) {
	err := apperror.NewForbiddenError("User", "delete", "user-456")

	httpErr := MapToHTTPError(err, "req-123")

	assert.Equal(t, 403, httpErr.StatusCode)
	assert.Equal(t, "FORBIDDEN", httpErr.Code)
}

func TestMapToHTTPError_DomainValidation(t *testing.T) {
	err := domainerror.NewEmailValidationError("email", "invalid", "invalid format")

	httpErr := MapToHTTPError(err, "req-123")

	assert.Equal(t, 400, httpErr.StatusCode)
	assert.Equal(t, "BAD_REQUEST", httpErr.Code)
}

func TestMapToHTTPError_DomainInvariant(t *testing.T) {
	err := domainerror.NewBusinessRuleViolationError("unique_email", "User", "email must be unique")

	httpErr := MapToHTTPError(err, "req-123")

	assert.Equal(t, 409, httpErr.StatusCode)
	assert.Equal(t, "CONFLICT", httpErr.Code)
}

func TestMapToHTTPError_DatabaseError(t *testing.T) {
	err := infraerror.NewDatabaseError("FindByID", "users", errors.New("connection failed"), false)

	httpErr := MapToHTTPError(err, "req-123")

	assert.Equal(t, 500, httpErr.StatusCode)
	assert.Equal(t, "INTERNAL_SERVER_ERROR", httpErr.Code)
}

func TestMapToHTTPError_ConnectionError(t *testing.T) {
	err := infraerror.NewConnectionError("database", errors.New("failed"), false)

	httpErr := MapToHTTPError(err, "req-123")

	assert.Equal(t, 503, httpErr.StatusCode)
	assert.Equal(t, "SERVICE_UNAVAILABLE", httpErr.Code)
}

func TestMapToHTTPError_UnknownError(t *testing.T) {
	err := errors.New("some unknown error")

	httpErr := MapToHTTPError(err, "req-123")

	assert.Equal(t, 500, httpErr.StatusCode)
	assert.Equal(t, "INTERNAL_SERVER_ERROR", httpErr.Code)
}
