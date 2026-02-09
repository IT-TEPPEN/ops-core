package service

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"opscore/backend/internal/authentication/application/port"
	"opscore/backend/internal/authentication/domain/entity"
	domain_err "opscore/backend/internal/authentication/domain/err"
	"opscore/backend/internal/authentication/domain/value_object"
)

// MockUserRepository is a mock implementation of repository.UserRepository
type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) FindByID(ctx context.Context, userID value_object.UserID) (*entity.User, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.User), args.Error(1)
}

func (m *MockUserRepository) FindByProviderIdentity(ctx context.Context, provider string, subject string) (*entity.User, error) {
	args := m.Called(ctx, provider, subject)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.User), args.Error(1)
}

func (m *MockUserRepository) FindByIDWithIdentities(ctx context.Context, userID value_object.UserID) (*entity.User, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.User), args.Error(1)
}

func (m *MockUserRepository) FindByEmail(ctx context.Context, email string) (*entity.User, error) {
	args := m.Called(ctx, email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.User), args.Error(1)
}

func (m *MockUserRepository) Save(ctx context.Context, user *entity.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *MockUserRepository) Delete(ctx context.Context, userID value_object.UserID) error {
	args := m.Called(ctx, userID)
	return args.Error(0)
}

// MockSessionRepository is a mock implementation of repository.SessionRepository
type MockSessionRepository struct {
	mock.Mock
}

func (m *MockSessionRepository) Create(ctx context.Context, session *entity.Session) error {
	args := m.Called(ctx, session)
	return args.Error(0)
}

func (m *MockSessionRepository) FindByID(ctx context.Context, sessionID value_object.SessionID) (*entity.Session, error) {
	args := m.Called(ctx, sessionID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.Session), args.Error(1)
}

func (m *MockSessionRepository) FindByTokenHash(ctx context.Context, tokenHash string) (*entity.Session, error) {
	args := m.Called(ctx, tokenHash)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.Session), args.Error(1)
}

func (m *MockSessionRepository) FindActiveByUserID(ctx context.Context, userID value_object.UserID) ([]*entity.Session, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*entity.Session), args.Error(1)
}

func (m *MockSessionRepository) Update(ctx context.Context, session *entity.Session) error {
	args := m.Called(ctx, session)
	return args.Error(0)
}

func (m *MockSessionRepository) FindByUserID(ctx context.Context, userID value_object.UserID) ([]*entity.Session, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*entity.Session), args.Error(1)
}

func (m *MockSessionRepository) RevokeAllByUserID(ctx context.Context, userID value_object.UserID) (int, error) {
	args := m.Called(ctx, userID)
	return args.Int(0), args.Error(1)
}

func (m *MockSessionRepository) Delete(ctx context.Context, sessionID value_object.SessionID) error {
	args := m.Called(ctx, sessionID)
	return args.Error(0)
}

func (m *MockSessionRepository) DeleteExpired(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

// MockJWTProvider is a mock implementation of port.JWTProvider
type MockJWTProvider struct {
	mock.Mock
}

// Ensure MockJWTProvider implements port.JWTProvider
var _ port.JWTProvider = (*MockJWTProvider)(nil)

func (m *MockJWTProvider) GenerateAccessToken(userID string, jti string, expiresIn time.Duration) (string, error) {
	args := m.Called(userID, jti, expiresIn)
	return args.String(0), args.Error(1)
}

func (m *MockJWTProvider) ValidateAccessToken(token string) (*port.JWTClaims, error) {
	args := m.Called(token)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*port.JWTClaims), args.Error(1)
}

// Test: Successful session creation with short-lived session (rememberMe = false)
func TestSessionService_CreateSession_Success_ShortLived(t *testing.T) {
	// Setup
	mockUserRepo := new(MockUserRepository)
	mockSessionRepo := new(MockSessionRepository)
	mockJWTProvider := new(MockJWTProvider)
	service := NewSessionService(mockUserRepo, mockSessionRepo, mockJWTProvider)

	ctx := context.Background()
	userID := value_object.NewUserID()

	// Create mock user
	user := entity.NewUser(userID, "test@example.com", "Test User", "")
	mockUserRepo.On("FindByID", ctx, userID).Return(user, nil)

	// Mock session repository
	mockSessionRepo.On("Create", ctx, mock.AnythingOfType("*entity.Session")).Return(nil)

	// Mock JWT provider
	mockJWTProvider.On("GenerateAccessToken", userID.String(), mock.AnythingOfType("string"), AccessTokenDuration).Return("mock-access-token", nil)

	// Execute
	result, err := service.CreateSession(ctx, userID, false)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "mock-access-token", result.AccessToken)
	assert.NotEmpty(t, result.RefreshToken)
	assert.Equal(t, user.ID().String(), result.User.UserID)
	assert.Equal(t, user.Email(), result.User.Email)
	assert.Equal(t, user.DisplayName(), result.User.DisplayName)

	mockUserRepo.AssertExpectations(t)
	mockSessionRepo.AssertExpectations(t)
	mockJWTProvider.AssertExpectations(t)
}

// Test: Successful session creation with long-lived session (rememberMe = true)
func TestSessionService_CreateSession_Success_LongLived(t *testing.T) {
	// Setup
	mockUserRepo := new(MockUserRepository)
	mockSessionRepo := new(MockSessionRepository)
	mockJWTProvider := new(MockJWTProvider)
	service := NewSessionService(mockUserRepo, mockSessionRepo, mockJWTProvider)

	ctx := context.Background()
	userID := value_object.NewUserID()

	// Create mock user
	user := entity.NewUser(userID, "test@example.com", "Test User", "")
	mockUserRepo.On("FindByID", ctx, userID).Return(user, nil)

	// Mock session repository
	mockSessionRepo.On("Create", ctx, mock.AnythingOfType("*entity.Session")).Return(nil)

	// Mock JWT provider
	mockJWTProvider.On("GenerateAccessToken", userID.String(), mock.AnythingOfType("string"), AccessTokenDuration).Return("mock-access-token", nil)

	// Execute
	result, err := service.CreateSession(ctx, userID, true)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "mock-access-token", result.AccessToken)
	assert.NotEmpty(t, result.RefreshToken)

	mockUserRepo.AssertExpectations(t)
	mockSessionRepo.AssertExpectations(t)
	mockJWTProvider.AssertExpectations(t)
}

// Test: User not found
func TestSessionService_CreateSession_UserNotFound(t *testing.T) {
	// Setup
	mockUserRepo := new(MockUserRepository)
	mockSessionRepo := new(MockSessionRepository)
	mockJWTProvider := new(MockJWTProvider)
	service := NewSessionService(mockUserRepo, mockSessionRepo, mockJWTProvider)

	ctx := context.Background()
	userID := value_object.NewUserID()

	mockUserRepo.On("FindByID", ctx, userID).Return(nil, domain_err.NewNotFoundUserError(userID.String()))

	// Execute
	result, err := service.CreateSession(ctx, userID, false)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "user not found")

	mockUserRepo.AssertExpectations(t)
}

// Test: Session repository failure
func TestSessionService_CreateSession_SessionRepositoryFailure(t *testing.T) {
	// Setup
	mockUserRepo := new(MockUserRepository)
	mockSessionRepo := new(MockSessionRepository)
	mockJWTProvider := new(MockJWTProvider)
	service := NewSessionService(mockUserRepo, mockSessionRepo, mockJWTProvider)

	ctx := context.Background()
	userID := value_object.NewUserID()

	user := entity.NewUser(userID, "test@example.com", "Test User", "")
	mockUserRepo.On("FindByID", ctx, userID).Return(user, nil)

	mockSessionRepo.On("Create", ctx, mock.AnythingOfType("*entity.Session")).Return(domain_err.NewDataPersistFailureError("database error"))

	// Execute
	result, err := service.CreateSession(ctx, userID, false)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "failed to save session")

	mockUserRepo.AssertExpectations(t)
	mockSessionRepo.AssertExpectations(t)
}

// Test: JWT generation failure
func TestSessionService_CreateSession_JWTGenerationFailure(t *testing.T) {
	// Setup
	mockUserRepo := new(MockUserRepository)
	mockSessionRepo := new(MockSessionRepository)
	mockJWTProvider := new(MockJWTProvider)
	service := NewSessionService(mockUserRepo, mockSessionRepo, mockJWTProvider)

	ctx := context.Background()
	userID := value_object.NewUserID()

	user := entity.NewUser(userID, "test@example.com", "Test User", "")
	mockUserRepo.On("FindByID", ctx, userID).Return(user, nil)

	mockSessionRepo.On("Create", ctx, mock.AnythingOfType("*entity.Session")).Return(nil)

	mockJWTProvider.On("GenerateAccessToken", userID.String(), mock.AnythingOfType("string"), AccessTokenDuration).Return("", assert.AnError)

	// Execute
	result, err := service.CreateSession(ctx, userID, false)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "failed to generate access token")

	mockUserRepo.AssertExpectations(t)
	mockSessionRepo.AssertExpectations(t)
	mockJWTProvider.AssertExpectations(t)
}
