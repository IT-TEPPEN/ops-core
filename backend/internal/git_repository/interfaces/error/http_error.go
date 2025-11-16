package error

import (
	"encoding/json"
	"net/http"
	"time"
)

// HTTPError represents an HTTP error response
type HTTPError struct {
	StatusCode int                    `json:"-"`
	Code       string                 `json:"code"`
	Message    string                 `json:"message"`
	Details    map[string]interface{} `json:"details,omitempty"`
	RequestID  string                 `json:"request_id,omitempty"`
	Timestamp  time.Time              `json:"timestamp"`
}

func (e *HTTPError) Error() string {
	return e.Message
}

// WriteJSON writes the error as JSON response
func (e *HTTPError) WriteJSON(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(e.StatusCode)
	json.NewEncoder(w).Encode(e)
}

// NewHTTPError creates a new HTTPError
func NewHTTPError(statusCode int, code, message string) *HTTPError {
	return &HTTPError{
		StatusCode: statusCode,
		Code:       code,
		Message:    message,
		Details:    make(map[string]interface{}),
		Timestamp:  time.Now(),
	}
}

// WithDetails adds details to the error
func (e *HTTPError) WithDetails(details map[string]interface{}) *HTTPError {
	e.Details = details
	return e
}

// WithRequestID adds request ID to the error
func (e *HTTPError) WithRequestID(requestID string) *HTTPError {
	e.RequestID = requestID
	return e
}

// Predefined HTTP errors
func BadRequest(message string) *HTTPError {
	return NewHTTPError(http.StatusBadRequest, "BAD_REQUEST", message)
}

func Unauthorized(message string) *HTTPError {
	return NewHTTPError(http.StatusUnauthorized, "UNAUTHORIZED", message)
}

func Forbidden(message string) *HTTPError {
	return NewHTTPError(http.StatusForbidden, "FORBIDDEN", message)
}

func NotFound(message string) *HTTPError {
	return NewHTTPError(http.StatusNotFound, "NOT_FOUND", message)
}

func Conflict(message string) *HTTPError {
	return NewHTTPError(http.StatusConflict, "CONFLICT", message)
}

func InternalServerError(message string) *HTTPError {
	return NewHTTPError(http.StatusInternalServerError, "INTERNAL_SERVER_ERROR", message)
}

func ServiceUnavailable(message string) *HTTPError {
	return NewHTTPError(http.StatusServiceUnavailable, "SERVICE_UNAVAILABLE", message)
}
