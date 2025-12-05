package repository

import (
	"context"

	"opscore/backend/internal/user/domain/entity"
	"opscore/backend/internal/user/domain/value_object"

	"github.com/stretchr/testify/mock"
)

// MockGroupRepository is a mock implementation of GroupRepository
type MockGroupRepository struct {
	mock.Mock
}

// Save mocks the Save method
func (m *MockGroupRepository) Save(ctx context.Context, group entity.Group) error {
	args := m.Called(ctx, group)
	return args.Error(0)
}

// FindByID mocks the FindByID method
func (m *MockGroupRepository) FindByID(ctx context.Context, id value_object.GroupID) (entity.Group, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(entity.Group), args.Error(1)
}

// FindByMemberID mocks the FindByMemberID method
func (m *MockGroupRepository) FindByMemberID(ctx context.Context, userID value_object.UserID) ([]entity.Group, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]entity.Group), args.Error(1)
}

// FindAll mocks the FindAll method
func (m *MockGroupRepository) FindAll(ctx context.Context) ([]entity.Group, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]entity.Group), args.Error(1)
}

// Update mocks the Update method
func (m *MockGroupRepository) Update(ctx context.Context, group entity.Group) error {
	args := m.Called(ctx, group)
	return args.Error(0)
}

// Delete mocks the Delete method
func (m *MockGroupRepository) Delete(ctx context.Context, id value_object.GroupID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}
