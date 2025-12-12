package error

import (
	"errors"

	apperror "opscore/backend/internal/document/application/error"
)

// MapToHTTPError maps application errors to HTTP errors
func MapToHTTPError(err error, requestID string) *HTTPError {
	var httpErr *HTTPError

	switch {
	case errors.Is(err, apperror.ErrNotFound):
		httpErr = NotFound("Resource not found")
		var notFoundErr *apperror.NotFoundError
		if errors.As(err, &notFoundErr) {
			httpErr = httpErr.WithDetails(map[string]interface{}{
				"resource_type": notFoundErr.ResourceType,
				"resource_id":   notFoundErr.ResourceID,
			})
		}

	case errors.Is(err, apperror.ErrBadRequest):
		httpErr = BadRequest("Validation failed")
		var validationErr *apperror.ValidationFailedError
		if errors.As(err, &validationErr) {
			// Use the error code from ValidationFailedError
			httpErr.Code = string(validationErr.Code)
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

	case errors.Is(err, apperror.ErrConflict):
		httpErr = Conflict("Resource conflict")
		var conflictErr *apperror.ConflictError
		if errors.As(err, &conflictErr) {
			httpErr = httpErr.WithDetails(map[string]interface{}{
				"resource_type": conflictErr.ResourceType,
				"identifier":    conflictErr.Identifier,
				"reason":        conflictErr.Reason,
			})
		}

	case errors.Is(err, apperror.ErrUnauthorized):
		httpErr = Unauthorized("Authentication required")

	case errors.Is(err, apperror.ErrForbidden):
		httpErr = Forbidden("Access denied")

	default:
		httpErr = InternalServerError("An unexpected error occurred")
	}

	return httpErr.WithRequestID(requestID)
}
