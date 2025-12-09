package handlers

import (
	"github.com/stretchr/testify/mock"
)

// MockLogger is a mock implementation of Logger for testing.
type MockLogger struct {
	mock.Mock
}

// Info mocks the Info method.
func (m *MockLogger) Info(msg string, args ...any) {
	m.Called(msg, args)
}

// Error mocks the Error method.
func (m *MockLogger) Error(msg string, args ...any) {
	m.Called(msg, args)
}

// Debug mocks the Debug method.
func (m *MockLogger) Debug(msg string, args ...any) {
	m.Called(msg, args)
}

// Warn mocks the Warn method.
func (m *MockLogger) Warn(msg string, args ...any) {
	m.Called(msg, args)
}
