package error

import (
	"errors"
	"testing"
)

func TestDatabaseError_Error(t *testing.T) {
	cause := errors.New("connection refused")
	err := &DatabaseError{
		Code:      CodeDatabaseConnection,
		Operation: "FindByID",
		Table:     "repositories",
		Cause:     cause,
		Retryable: true,
	}

	if !errors.Is(err, cause) {
		t.Error("Expected error to wrap cause")
	}
}

func TestDatabaseError_Is(t *testing.T) {
	err := &DatabaseError{
		Code:      CodeDatabaseQuery,
		Operation: "Save",
		Table:     "repositories",
		Cause:     errors.New("query error"),
		Retryable: false,
	}

	if !errors.Is(err, ErrDatabase) {
		t.Error("DatabaseError should match ErrDatabase sentinel")
	}
}

func TestDatabaseError_Is_Retryable(t *testing.T) {
	err := &DatabaseError{
		Code:      CodeDatabaseTimeout,
		Operation: "Query",
		Table:     "repositories",
		Cause:     errors.New("timeout"),
		Retryable: true,
	}

	if !errors.Is(err, ErrRetryable) {
		t.Error("Retryable DatabaseError should match ErrRetryable sentinel")
	}
}

func TestDatabaseError_Unwrap(t *testing.T) {
	cause := errors.New("underlying error")
	err := &DatabaseError{
		Code:      CodeDatabaseQuery,
		Operation: "Save",
		Table:     "repositories",
		Cause:     cause,
		Retryable: false,
	}

	if err.Unwrap() != cause {
		t.Error("Unwrap should return the cause")
	}
}

func TestDatabaseError_ErrorCode(t *testing.T) {
	err := &DatabaseError{
		Code:      CodeDatabaseConnection,
		Operation: "Connect",
		Table:     "repositories",
		Cause:     errors.New("error"),
		Retryable: true,
	}

	if err.ErrorCode() != CodeDatabaseConnection {
		t.Errorf("Expected error code '%s', got '%s'", CodeDatabaseConnection, err.ErrorCode())
	}
}

func TestExternalAPIError_Error(t *testing.T) {
	err := &ExternalAPIError{
		Code:       CodeExternalAPIError,
		Service:    "GitHub",
		Endpoint:   "/repos/owner/repo",
		StatusCode: 500,
		Cause:      errors.New("internal server error"),
		Retryable:  true,
	}

	expected := "[GITREPO_INF_EXT_001] external API error calling GitHub at /repos/owner/repo (status 500): internal server error"
	if err.Error() != expected {
		t.Errorf("Expected error message '%s', got '%s'", expected, err.Error())
	}
}

func TestExternalAPIError_Is(t *testing.T) {
	err := &ExternalAPIError{
		Code:       CodeExternalAPIError,
		Service:    "GitHub",
		Endpoint:   "/test",
		StatusCode: 404,
		Cause:      errors.New("not found"),
		Retryable:  false,
	}

	if !errors.Is(err, ErrExternalAPI) {
		t.Error("ExternalAPIError should match ErrExternalAPI sentinel")
	}
}

func TestExternalAPIError_Is_Retryable(t *testing.T) {
	err := &ExternalAPIError{
		Code:       CodeExternalAPITimeout,
		Service:    "GitHub",
		Endpoint:   "/test",
		StatusCode: 503,
		Cause:      errors.New("timeout"),
		Retryable:  true,
	}

	if !errors.Is(err, ErrRetryable) {
		t.Error("Retryable ExternalAPIError should match ErrRetryable sentinel")
	}
}

func TestConnectionError_Error(t *testing.T) {
	err := &ConnectionError{
		Code:   CodeConnectionFailed,
		Target: "database",
		Cause:  errors.New("connection refused"),
	}

	expected := "[GITREPO_INF_CONN_001] connection error to database: connection refused"
	if err.Error() != expected {
		t.Errorf("Expected error message '%s', got '%s'", expected, err.Error())
	}
}

func TestConnectionError_Is(t *testing.T) {
	err := &ConnectionError{
		Code:   CodeConnectionFailed,
		Target: "database",
		Cause:  errors.New("error"),
	}

	if !errors.Is(err, ErrConnection) {
		t.Error("ConnectionError should match ErrConnection sentinel")
	}
}

func TestStorageError_Error(t *testing.T) {
	err := &StorageError{
		Code:      CodeStorageOperation,
		Operation: "Upload",
		Path:      "/tmp/file.txt",
		Cause:     errors.New("disk full"),
	}

	expected := "[GITREPO_INF_STOR_001] storage error during Upload at path /tmp/file.txt: disk full"
	if err.Error() != expected {
		t.Errorf("Expected error message '%s', got '%s'", expected, err.Error())
	}
}

func TestStorageError_Unwrap(t *testing.T) {
	cause := errors.New("underlying error")
	err := &StorageError{
		Code:      CodeStorageOperation,
		Operation: "Delete",
		Path:      "/test",
		Cause:     cause,
	}

	if err.Unwrap() != cause {
		t.Error("Unwrap should return the cause")
	}
}

func TestStorageError_ErrorCode(t *testing.T) {
	err := &StorageError{
		Code:      CodeStorageNotFound,
		Operation: "Download",
		Path:      "/test",
		Cause:     errors.New("not found"),
	}

	if err.ErrorCode() != CodeStorageNotFound {
		t.Errorf("Expected error code '%s', got '%s'", CodeStorageNotFound, err.ErrorCode())
	}
}
