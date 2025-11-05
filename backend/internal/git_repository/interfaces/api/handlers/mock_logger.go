package handlers

import "github.com/stretchr/testify/mock"

// MockLogger is a mock implementation of the Logger interface for testing
type MockLogger struct {
	mock.Mock
}

// Info is a mock implementation of the Logger.Info method
func (m *MockLogger) Info(msg string, args ...any) {
	callArgs := []interface{}{msg}
	for _, arg := range args {
		callArgs = append(callArgs, arg)
	}
	m.Called(callArgs...)
}

// Error is a mock implementation of the Logger.Error method
func (m *MockLogger) Error(msg string, args ...any) {
	callArgs := []interface{}{msg}
	for _, arg := range args {
		callArgs = append(callArgs, arg)
	}
	m.Called(callArgs...)
}

// Debug is a mock implementation of the Logger.Debug method
func (m *MockLogger) Debug(msg string, args ...any) {
	callArgs := []interface{}{msg}
	for _, arg := range args {
		callArgs = append(callArgs, arg)
	}
	m.Called(callArgs...)
}

// Warn is a mock implementation of the Logger.Warn method
func (m *MockLogger) Warn(msg string, args ...any) {
	callArgs := []interface{}{msg}
	for _, arg := range args {
		callArgs = append(callArgs, arg)
	}
	m.Called(callArgs...)
}
