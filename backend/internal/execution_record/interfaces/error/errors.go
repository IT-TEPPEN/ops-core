package error

import (
	"encoding/json"
	"net/http"
)

// HTTPError represents an HTTP error response.
type HTTPError struct {
	StatusCode int                    `json:"-"`
	Code       string                 `json:"code"`
	Message    string                 `json:"message"`
	Details    map[string]interface{} `json:"details,omitempty"`
	RequestID  string                 `json:"request_id,omitempty"`
}

// Error implements the error interface.
func (e *HTTPError) Error() string {
	return e.Message
}

// WriteJSON writes the error as JSON response.
func (e *HTTPError) WriteJSON(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(e.StatusCode)
	json.NewEncoder(w).Encode(e)
}

// NewHTTPError creates a new HTTPError.
func NewHTTPError(statusCode int, code, message string) *HTTPError {
	return &HTTPError{
		StatusCode: statusCode,
		Code:       code,
		Message:    message,
		Details:    make(map[string]interface{}),
	}
}

// WithDetails adds details to the error.
func (e *HTTPError) WithDetails(details map[string]interface{}) *HTTPError {
	e.Details = details
	return e
}

// WithRequestID adds request ID to the error.
func (e *HTTPError) WithRequestID(requestID string) *HTTPError {
	e.RequestID = requestID
	return e
}

// Predefined HTTP errors

// BadRequest creates a 400 Bad Request error.
func BadRequest(message string) *HTTPError {
	return NewHTTPError(http.StatusBadRequest, "BAD_REQUEST", message)
}

// Unauthorized creates a 401 Unauthorized error.
func Unauthorized(message string) *HTTPError {
	return NewHTTPError(http.StatusUnauthorized, "UNAUTHORIZED", message)
}

// Forbidden creates a 403 Forbidden error.
func Forbidden(message string) *HTTPError {
	return NewHTTPError(http.StatusForbidden, "FORBIDDEN", message)
}

// NotFound creates a 404 Not Found error.
func NotFound(message string) *HTTPError {
	return NewHTTPError(http.StatusNotFound, "NOT_FOUND", message)
}

// Conflict creates a 409 Conflict error.
func Conflict(message string) *HTTPError {
	return NewHTTPError(http.StatusConflict, "CONFLICT", message)
}

// InternalServerError creates a 500 Internal Server Error.
func InternalServerError(message string) *HTTPError {
	return NewHTTPError(http.StatusInternalServerError, "INTERNAL_SERVER_ERROR", message)
}

// ServiceUnavailable creates a 503 Service Unavailable error.
func ServiceUnavailable(message string) *HTTPError {
	return NewHTTPError(http.StatusServiceUnavailable, "SERVICE_UNAVAILABLE", message)
}
