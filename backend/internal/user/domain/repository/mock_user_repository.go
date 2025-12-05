package repository

import (
	"context"

	"opscore/backend/internal/user/domain/entity"
	"opscore/backend/internal/user/domain/value_object"

	"github.com/stretchr/testify/mock"
)

// MockUserRepository is a mock implementation of UserRepository
type MockUserRepository struct {
	mock.Mock
}

// Save mocks the Save method
func (m *MockUserRepository) Save(ctx context.Context, user entity.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

// FindByID mocks the FindByID method
func (m *MockUserRepository) FindByID(ctx context.Context, id value_object.UserID) (entity.User, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(entity.User), args.Error(1)
}

// FindByEmail mocks the FindByEmail method
func (m *MockUserRepository) FindByEmail(ctx context.Context, email value_object.Email) (entity.User, error) {
	args := m.Called(ctx, email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(entity.User), args.Error(1)
}

// FindAll mocks the FindAll method
func (m *MockUserRepository) FindAll(ctx context.Context) ([]entity.User, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]entity.User), args.Error(1)
}

// Update mocks the Update method
func (m *MockUserRepository) Update(ctx context.Context, user entity.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

// Delete mocks the Delete method
func (m *MockUserRepository) Delete(ctx context.Context, id value_object.UserID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}
