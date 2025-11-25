package error

import (
	"errors"
	"fmt"
)

// Sentinel errors for infrastructure failures
var (
	ErrDatabase    = errors.New("database error")
	ErrStorage     = errors.New("storage error")
	ErrConnection  = errors.New("connection error")
	ErrTimeout     = errors.New("timeout error")
	ErrRetryable   = errors.New("retryable error")
)

// DatabaseError represents a database operation failure.
type DatabaseError struct {
	Operation string
	Table     string
	Cause     error
	Retryable bool
}

func (e *DatabaseError) Error() string {
	return fmt.Sprintf("database error during %s on table %s: %v", e.Operation, e.Table, e.Cause)
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

// StorageError represents a storage operation failure.
type StorageError struct {
	Operation string
	Path      string
	Cause     error
}

func (e *StorageError) Error() string {
	return fmt.Sprintf("storage error during %s at path %s: %v", e.Operation, e.Path, e.Cause)
}

func (e *StorageError) Is(target error) bool {
	return target == ErrStorage
}

func (e *StorageError) Unwrap() error {
	return e.Cause
}

// ConnectionError represents a connection failure.
type ConnectionError struct {
	Target string
	Cause  error
}

func (e *ConnectionError) Error() string {
	return fmt.Sprintf("connection error to %s: %v", e.Target, e.Cause)
}

func (e *ConnectionError) Is(target error) bool {
	return target == ErrConnection
}

func (e *ConnectionError) Unwrap() error {
	return e.Cause
}
