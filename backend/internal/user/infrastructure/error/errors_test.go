package error

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDatabaseError_Error(t *testing.T) {
	cause := errors.New("connection refused")
	err := &DatabaseError{
		Code:      CodeDatabaseQuery,
		Operation: "FindByID",
		Table:     "users",
		Cause:     cause,
		Retryable: false,
	}

	expected := "[USI0002] database error during FindByID on table users: connection refused"
	assert.Equal(t, expected, err.Error())
}

func TestDatabaseError_Is(t *testing.T) {
	err := NewDatabaseError("FindByID", "users", errors.New("error"), false)
	assert.True(t, errors.Is(err, ErrDatabase))
}

func TestDatabaseError_Is_Retryable(t *testing.T) {
	err := NewDatabaseError("FindByID", "users", errors.New("timeout"), true)
	assert.True(t, errors.Is(err, ErrRetryable))
}

func TestConnectionError_Error(t *testing.T) {
	cause := errors.New("network unreachable")
	err := &ConnectionError{
		Code:   CodeConnectionFailed,
		Target: "database",
		Cause:  cause,
	}

	expected := "[USI0008] connection error to database: network unreachable"
	assert.Equal(t, expected, err.Error())
}

func TestConnectionError_Is(t *testing.T) {
	err := NewConnectionError("database", errors.New("failed"), false)
	assert.True(t, errors.Is(err, ErrConnection))
}

func TestNewDatabaseConnectionError(t *testing.T) {
	err := NewDatabaseConnectionError("Save", "users", errors.New("connection failed"))
	assert.Equal(t, CodeDatabaseConnection, err.Code)
	assert.True(t, err.Retryable)
}

func TestNewDatabaseConstraintError(t *testing.T) {
	err := NewDatabaseConstraintError("Save", "users", errors.New("unique violation"))
	assert.Equal(t, CodeDatabaseConstraint, err.Code)
	assert.False(t, err.Retryable)
}
