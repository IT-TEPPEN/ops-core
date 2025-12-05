package error

import "net/http"

// HTTPError represents an HTTP error response
type HTTPError struct {
	StatusCode int
	Code       string
	Message    string
	Details    map[string]interface{}
}

// WithDetails adds details to the HTTP error
func (e *HTTPError) WithDetails(details map[string]interface{}) *HTTPError {
	e.Details = details
	return e
}

// WithRequestID adds request ID to the HTTP error details
func (e *HTTPError) WithRequestID(requestID string) *HTTPError {
	if e.Details == nil {
		e.Details = make(map[string]interface{})
	}
	e.Details["request_id"] = requestID
	return e
}

// BadRequest returns a 400 Bad Request error
func BadRequest(message string) *HTTPError {
	return &HTTPError{
		StatusCode: http.StatusBadRequest,
		Code:       "BAD_REQUEST",
		Message:    message,
	}
}

// NotFound returns a 404 Not Found error
func NotFound(message string) *HTTPError {
	return &HTTPError{
		StatusCode: http.StatusNotFound,
		Code:       "NOT_FOUND",
		Message:    message,
	}
}

// Conflict returns a 409 Conflict error
func Conflict(message string) *HTTPError {
	return &HTTPError{
		StatusCode: http.StatusConflict,
		Code:       "CONFLICT",
		Message:    message,
	}
}

// InternalServerError returns a 500 Internal Server Error
func InternalServerError(message string) *HTTPError {
	return &HTTPError{
		StatusCode: http.StatusInternalServerError,
		Code:       "INTERNAL_ERROR",
		Message:    message,
	}
}

// Unauthorized returns a 401 Unauthorized error
func Unauthorized(message string) *HTTPError {
	return &HTTPError{
		StatusCode: http.StatusUnauthorized,
		Code:       "UNAUTHORIZED",
		Message:    message,
	}
}

// Forbidden returns a 403 Forbidden error
func Forbidden(message string) *HTTPError {
	return &HTTPError{
		StatusCode: http.StatusForbidden,
		Code:       "FORBIDDEN",
		Message:    message,
	}
}
