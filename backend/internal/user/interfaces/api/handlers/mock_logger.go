package handlers

// MockLogger implements the Logger interface for testing
type MockLogger struct{}

func (m *MockLogger) Info(msg string, args ...any)  {}
func (m *MockLogger) Error(msg string, args ...any) {}
func (m *MockLogger) Debug(msg string, args ...any) {}
func (m *MockLogger) Warn(msg string, args ...any)  {}
