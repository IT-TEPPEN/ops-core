package usecase

import (
	"context"

	"opscore/backend/internal/document/application/dto"

	"github.com/stretchr/testify/mock"
)

// MockVariableUseCase is a mock implementation of VariableUseCase for testing
type MockVariableUseCase struct {
	mock.Mock
}

// GetVariableDefinitions mocks the GetVariableDefinitions method
func (m *MockVariableUseCase) GetVariableDefinitions(ctx context.Context, documentID string) ([]dto.VariableDefinitionDTO, error) {
	args := m.Called(ctx, documentID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]dto.VariableDefinitionDTO), args.Error(1)
}

// ValidateVariableValues mocks the ValidateVariableValues method
func (m *MockVariableUseCase) ValidateVariableValues(ctx context.Context, documentID string, values []VariableValue) error {
	args := m.Called(ctx, documentID, values)
	return args.Error(0)
}

// SubstituteVariables mocks the SubstituteVariables method
func (m *MockVariableUseCase) SubstituteVariables(ctx context.Context, content string, values []VariableValue) (string, error) {
	args := m.Called(ctx, content, values)
	return args.String(0), args.Error(1)
}
