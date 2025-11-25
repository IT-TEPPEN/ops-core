package error

import (
	"errors"
	"fmt"
)

// Sentinel errors for infrastructure failures
var (
	ErrDatabase    = errors.New("database error")
	ErrExternalAPI = errors.New("external API error")
	ErrConnection  = errors.New("connection error")
	ErrTimeout     = errors.New("timeout error")
	ErrRetryable   = errors.New("retryable error")
)

// DatabaseError represents a database operation failure
type DatabaseError struct {
	Code      ErrorCode
	Operation string // "FindByID", "Save", "Delete", etc.
	Table     string
	Cause     error
	Retryable bool
}

func (e *DatabaseError) Error() string {
	return fmt.Sprintf("[%s] database error during %s on table %s: %v", e.Code, e.Operation, e.Table, e.Cause)
}

func (e *DatabaseError) Is(target error) bool {
	if target == ErrDatabase {
		return true
	}
	if target == ErrRetryable && e.Retryable {
		return true
	}
	return false
}

func (e *DatabaseError) Unwrap() error {
	return e.Cause
}

func (e *DatabaseError) ErrorCode() ErrorCode {
	return e.Code
}

// ExternalAPIError represents an external API call failure
type ExternalAPIError struct {
	Code       ErrorCode
	Service    string
	Endpoint   string
	StatusCode int
	Cause      error
	Retryable  bool
}

func (e *ExternalAPIError) Error() string {
	return fmt.Sprintf("[%s] external API error calling %s at %s (status %d): %v",
		e.Code, e.Service, e.Endpoint, e.StatusCode, e.Cause)
}

func (e *ExternalAPIError) Is(target error) bool {
	if target == ErrExternalAPI {
		return true
	}
	if target == ErrRetryable && e.Retryable {
		return true
	}
	return false
}

func (e *ExternalAPIError) Unwrap() error {
	return e.Cause
}

func (e *ExternalAPIError) ErrorCode() ErrorCode {
	return e.Code
}

// ConnectionError represents a connection failure
type ConnectionError struct {
	Code   ErrorCode
	Target string
	Cause  error
}

func (e *ConnectionError) Error() string {
	return fmt.Sprintf("[%s] connection error to %s: %v", e.Code, e.Target, e.Cause)
}

func (e *ConnectionError) Is(target error) bool {
	return target == ErrConnection
}

func (e *ConnectionError) Unwrap() error {
	return e.Cause
}

func (e *ConnectionError) ErrorCode() ErrorCode {
	return e.Code
}

// Factory functions to prevent manual error struct/code pairing errors

// NewDatabaseError creates a new DatabaseError with the appropriate error code
func NewDatabaseError(operation, table string, cause error, retryable bool) *DatabaseError {
	code := CodeDatabaseQuery
	if retryable {
		code = CodeDatabaseTimeout
	}

	return &DatabaseError{
		Code:      code,
		Operation: operation,
		Table:     table,
		Cause:     cause,
		Retryable: retryable,
	}
}

// NewDatabaseConnectionError creates a DatabaseError specifically for connection errors
func NewDatabaseConnectionError(operation, table string, cause error) *DatabaseError {
	return &DatabaseError{
		Code:      CodeDatabaseConnection,
		Operation: operation,
		Table:     table,
		Cause:     cause,
		Retryable: true,
	}
}

// NewDatabaseConstraintError creates a DatabaseError specifically for constraint violations
func NewDatabaseConstraintError(operation, table string, cause error) *DatabaseError {
	return &DatabaseError{
		Code:      CodeDatabaseConstraint,
		Operation: operation,
		Table:     table,
		Cause:     cause,
		Retryable: false,
	}
}

// NewConnectionError creates a new ConnectionError with the correct error code
func NewConnectionError(target string, cause error, isTimeout bool) *ConnectionError {
	code := CodeConnectionFailed
	if isTimeout {
		code = CodeConnectionTimeout
	}

	return &ConnectionError{
		Code:   code,
		Target: target,
		Cause:  cause,
	}
}
