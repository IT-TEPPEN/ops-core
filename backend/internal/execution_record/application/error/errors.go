package error

import "errors"

var (
	// ErrNotFound is returned when a resource is not found.
	ErrNotFound = errors.New("resource not found")

	// ErrUnauthorized is returned when access is unauthorized.
	ErrUnauthorized = errors.New("unauthorized")

	// ErrForbidden is returned when access is forbidden.
	ErrForbidden = errors.New("forbidden")

	// ErrConflict is returned when there is a resource conflict.
	ErrConflict = errors.New("resource conflict")

	// ErrBadRequest is returned when the request is invalid.
	ErrBadRequest = errors.New("bad request")

	// ErrValidationFailed is returned when validation fails.
	ErrValidationFailed = errors.New("validation failed")
)

// NotFoundError represents a resource not found error.
type NotFoundError struct {
	ResourceType string
	ResourceID   string
	Cause        error
}

func (e *NotFoundError) Error() string {
	if e.Cause != nil {
		return e.ResourceType + " not found (ID: " + e.ResourceID + "): " + e.Cause.Error()
	}
	return e.ResourceType + " not found (ID: " + e.ResourceID + ")"
}

func (e *NotFoundError) Is(target error) bool {
	return target == ErrNotFound
}

func (e *NotFoundError) Unwrap() error {
	return e.Cause
}

// ValidationError represents a validation error.
type ValidationError struct {
	Field   string
	Message string
}

func (e *ValidationError) Error() string {
	return "validation failed for field '" + e.Field + "': " + e.Message
}

func (e *ValidationError) Is(target error) bool {
	return target == ErrValidationFailed
}

// ConflictError represents a resource conflict error.
type ConflictError struct {
	ResourceType string
	Identifier   string
	Reason       string
	Cause        error
}

func (e *ConflictError) Error() string {
	msg := "conflict: " + e.ResourceType + " with identifier '" + e.Identifier + "': " + e.Reason
	if e.Cause != nil {
		msg += ": " + e.Cause.Error()
	}
	return msg
}

func (e *ConflictError) Is(target error) bool {
	return target == ErrConflict
}

func (e *ConflictError) Unwrap() error {
	return e.Cause
}

// ForbiddenError represents a forbidden access error.
type ForbiddenError struct {
	Resource string
	Action   string
	UserID   string
}

func (e *ForbiddenError) Error() string {
	return "user " + e.UserID + " is forbidden to " + e.Action + " on resource " + e.Resource
}

func (e *ForbiddenError) Is(target error) bool {
	return target == ErrForbidden
}
