package error

import (
	"errors"
	"fmt"
)

// Sentinel errors for common application failures
var (
	ErrNotFound     = errors.New("resource not found")
	ErrUnauthorized = errors.New("unauthorized")
	ErrForbidden    = errors.New("forbidden")
	ErrConflict     = errors.New("resource conflict")
	ErrBadRequest   = errors.New("bad request")
)

// NotFoundError represents a resource not found error
type NotFoundError struct {
	Code         ErrorCode
	ResourceType string
	ResourceID   string
	Cause        error
}

func (e *NotFoundError) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("[%s] %s not found (ID: %s): %v", e.Code, e.ResourceType, e.ResourceID, e.Cause)
	}
	return fmt.Sprintf("[%s] %s not found (ID: %s)", e.Code, e.ResourceType, e.ResourceID)
}

func (e *NotFoundError) Is(target error) bool {
	return target == ErrNotFound
}

func (e *NotFoundError) Unwrap() error {
	return e.Cause
}

func (e *NotFoundError) ErrorCode() ErrorCode {
	return e.Code
}

// UnauthorizedError represents an authentication failure
type UnauthorizedError struct {
	Code   ErrorCode
	Reason string
	Cause  error
}

func (e *UnauthorizedError) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("[%s] unauthorized: %s: %v", e.Code, e.Reason, e.Cause)
	}
	return fmt.Sprintf("[%s] unauthorized: %s", e.Code, e.Reason)
}

func (e *UnauthorizedError) Is(target error) bool {
	return target == ErrUnauthorized
}

func (e *UnauthorizedError) Unwrap() error {
	return e.Cause
}

func (e *UnauthorizedError) ErrorCode() ErrorCode {
	return e.Code
}

// ForbiddenError represents an authorization failure
type ForbiddenError struct {
	Code     ErrorCode
	Resource string
	Action   string
	UserID   string
}

func (e *ForbiddenError) Error() string {
	return fmt.Sprintf("[%s] user %s is forbidden to %s on resource %s", e.Code, e.UserID, e.Action, e.Resource)
}

func (e *ForbiddenError) Is(target error) bool {
	return target == ErrForbidden
}

func (e *ForbiddenError) ErrorCode() ErrorCode {
	return e.Code
}

// ConflictError represents a resource conflict (e.g., duplicate key)
type ConflictError struct {
	Code         ErrorCode
	ResourceType string
	Identifier   string
	Reason       string
	Cause        error
}

func (e *ConflictError) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("[%s] conflict: %s with identifier '%s' already exists: %s: %v",
			e.Code, e.ResourceType, e.Identifier, e.Reason, e.Cause)
	}
	return fmt.Sprintf("[%s] conflict: %s with identifier '%s' already exists: %s",
		e.Code, e.ResourceType, e.Identifier, e.Reason)
}

func (e *ConflictError) Is(target error) bool {
	return target == ErrConflict
}

func (e *ConflictError) Unwrap() error {
	return e.Cause
}

func (e *ConflictError) ErrorCode() ErrorCode {
	return e.Code
}

// ValidationFailedError represents application-level validation failure
type ValidationFailedError struct {
	Code   ErrorCode
	Errors []FieldError
}

type FieldError struct {
	Field   string
	Message string
	Code    string
}

func (e *ValidationFailedError) Error() string {
	return fmt.Sprintf("[%s] validation failed: %d error(s)", e.Code, len(e.Errors))
}

func (e *ValidationFailedError) Is(target error) bool {
	return target == ErrBadRequest
}

func (e *ValidationFailedError) ErrorCode() ErrorCode {
	return e.Code
}

// Factory functions to prevent manual error struct/code pairing errors

// NewNotFoundError creates a new NotFoundError with the correct error code
func NewNotFoundError(resourceType, resourceID string, cause error) *NotFoundError {
	return &NotFoundError{
		Code:         CodeResourceNotFound,
		ResourceType: resourceType,
		ResourceID:   resourceID,
		Cause:        cause,
	}
}

// NewUnauthorizedError creates a new UnauthorizedError with the correct error code
func NewUnauthorizedError(reason string, cause error) *UnauthorizedError {
	return &UnauthorizedError{
		Code:   CodeUnauthorized,
		Reason: reason,
		Cause:  cause,
	}
}

// NewForbiddenError creates a new ForbiddenError with the correct error code
func NewForbiddenError(resource, action, userID string) *ForbiddenError {
	return &ForbiddenError{
		Code:     CodeForbidden,
		Resource: resource,
		Action:   action,
		UserID:   userID,
	}
}

// NewConflictError creates a new ConflictError with the correct error code
func NewConflictError(resourceType, identifier, reason string, cause error) *ConflictError {
	return &ConflictError{
		Code:         CodeResourceConflict,
		ResourceType: resourceType,
		Identifier:   identifier,
		Reason:       reason,
		Cause:        cause,
	}
}

// NewValidationFailedError creates a new ValidationFailedError with the correct error code
func NewValidationFailedError(errors []FieldError) *ValidationFailedError {
	return &ValidationFailedError{
		Code:   CodeValidationFailed,
		Errors: errors,
	}
}
