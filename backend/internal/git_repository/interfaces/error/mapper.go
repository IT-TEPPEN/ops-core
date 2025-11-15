package error

import (
	"errors"

	apperror "opscore/backend/internal/git_repository/application/error"
	domainerror "opscore/backend/internal/git_repository/domain/error"
	infraerror "opscore/backend/internal/git_repository/infrastructure/error"
)

// MapToHTTPError maps application/domain/infrastructure errors to HTTP errors
func MapToHTTPError(err error, requestID string) *HTTPError {
	var httpErr *HTTPError

	switch {
	// Application layer errors
	case errors.Is(err, apperror.ErrNotFound):
		httpErr = NotFound("Resource not found")
		// Extract details if it's a NotFoundError
		var notFoundErr *apperror.NotFoundError
		if errors.As(err, &notFoundErr) {
			httpErr = httpErr.WithDetails(map[string]interface{}{
				"resource_type": notFoundErr.ResourceType,
				"resource_id":   notFoundErr.ResourceID,
			})
		}

	case errors.Is(err, apperror.ErrUnauthorized):
		httpErr = Unauthorized("Authentication required")

	case errors.Is(err, apperror.ErrForbidden):
		httpErr = Forbidden("Access denied")

	case errors.Is(err, apperror.ErrConflict):
		httpErr = Conflict("Resource already exists")
		// Extract details if it's a ConflictError
		var conflictErr *apperror.ConflictError
		if errors.As(err, &conflictErr) {
			httpErr = httpErr.WithDetails(map[string]interface{}{
				"resource_type": conflictErr.ResourceType,
				"identifier":    conflictErr.Identifier,
				"reason":        conflictErr.Reason,
			})
		}

	case errors.Is(err, apperror.ErrBadRequest):
		httpErr = BadRequest("Invalid request")
		// Extract details if it's a ValidationFailedError
		var validationErr *apperror.ValidationFailedError
		if errors.As(err, &validationErr) {
			fieldErrors := make([]map[string]interface{}, len(validationErr.Errors))
			for i, fieldErr := range validationErr.Errors {
				fieldErrors[i] = map[string]interface{}{
					"field":   fieldErr.Field,
					"message": fieldErr.Message,
					"code":    fieldErr.Code,
				}
			}
			httpErr = httpErr.WithDetails(map[string]interface{}{
				"validation_errors": fieldErrors,
			})
		}

	// Domain layer errors (when not wrapped by application layer)
	case errors.Is(err, domainerror.ErrInvalidEntity):
		httpErr = BadRequest("Invalid entity")
		// Extract details if it's a ValidationError
		var domainValidationErr *domainerror.ValidationError
		if errors.As(err, &domainValidationErr) {
			httpErr = httpErr.WithDetails(map[string]interface{}{
				"field":   domainValidationErr.Field,
				"message": domainValidationErr.Message,
			})
		}

	case errors.Is(err, domainerror.ErrInvariantViolation):
		httpErr = Conflict("Business rule violation")

	// Infrastructure layer errors
	case errors.Is(err, infraerror.ErrConnection):
		httpErr = ServiceUnavailable("Service temporarily unavailable")

	case errors.Is(err, infraerror.ErrTimeout):
		httpErr = ServiceUnavailable("Request timeout")

	case errors.Is(err, infraerror.ErrRetryable):
		httpErr = ServiceUnavailable("Temporary error, please retry")

	case errors.Is(err, infraerror.ErrDatabase):
		httpErr = InternalServerError("Database error occurred")

	case errors.Is(err, infraerror.ErrExternalAPI):
		httpErr = InternalServerError("External service error")

	// Default to internal server error
	default:
		httpErr = InternalServerError("An unexpected error occurred")
	}

	return httpErr.WithRequestID(requestID)
}
