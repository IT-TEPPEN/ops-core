package usecase

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"opscore/backend/internal/authentication/application/dto"
	"opscore/backend/internal/authentication/domain/entity"
	domain_err "opscore/backend/internal/authentication/domain/err"
	"opscore/backend/internal/authentication/domain/value_object"
)

// MockUserRepository is a mock implementation of UserRepository
type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) Save(ctx context.Context, user *entity.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *MockUserRepository) FindByID(ctx context.Context, userID value_object.UserID) (*entity.User, error) {
	args := m.Called(ctx, userID)
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

func (m *MockUserRepository) FindByProviderIdentity(ctx context.Context, provider, providerUserID string) (*entity.User, error) {
	args := m.Called(ctx, provider, providerUserID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.User), args.Error(1)
}

func (m *MockUserRepository) Delete(ctx context.Context, userID value_object.UserID) error {
	args := m.Called(ctx, userID)
	return args.Error(0)
}

// MockSessionService is a mock implementation of SessionService
type MockSessionService struct {
	mock.Mock
}

func (m *MockSessionService) CreateSession(ctx context.Context, userID value_object.UserID, rememberMe bool) (*dto.SessionResult, error) {
	args := m.Called(ctx, userID, rememberMe)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.SessionResult), args.Error(1)
}

func TestIssueSessionFromProviderUser_Execute_NewUser(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mockRepo := new(MockUserRepository)
	mockSessionService := new(MockSessionService)
	usecase := NewIssueSessionFromProviderUser(mockRepo, mockSessionService)

	providerInfo := dto.ProviderUserInfo{
		Provider:      "google",
		Subject:       "google-user-123",
		Email:         "newuser@example.com",
		EmailVerified: true,
		Name:          "New User",
		Picture:       "https://example.com/picture.jpg",
	}

	// User does not exist yet
	mockRepo.On("FindByProviderIdentity", ctx, "google", "google-user-123").
		Return(nil, domain_err.NewNotFoundUserError("google-user-123"))

	// Save will be called with new user
	mockRepo.On("Save", ctx, mock.AnythingOfType("*entity.User")).
		Return(nil)

	// SessionService will be called to create session
	expectedSession := &dto.SessionResult{
		AccessToken:  "mock-access-token",
		RefreshToken: "mock-refresh-token",
		User: dto.UserInfo{
			UserID:      "user-123",
			Email:       "newuser@example.com",
			DisplayName: "New User",
			PictureURL:  "https://example.com/picture.jpg",
		},
	}
	mockSessionService.On("CreateSession", ctx, mock.AnythingOfType("value_object.UserID"), false).
		Return(expectedSession, nil)

	// Act
	result, err := usecase.Execute(ctx, providerInfo)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "mock-access-token", result.AccessToken)
	assert.Equal(t, "mock-refresh-token", result.RefreshToken)
	assert.Equal(t, "newuser@example.com", result.User.Email)
	assert.Equal(t, "New User", result.User.DisplayName)
	mockRepo.AssertExpectations(t)
	mockSessionService.AssertExpectations(t)
}

func TestIssueSessionFromProviderUser_Execute_ExistingUser(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mockRepo := new(MockUserRepository)
	mockSessionService := new(MockSessionService)
	usecase := NewIssueSessionFromProviderUser(mockRepo, mockSessionService)

	providerInfo := dto.ProviderUserInfo{
		Provider:      "github",
		Subject:       "github-user-456",
		Email:         "existinguser@example.com",
		EmailVerified: true,
		Name:          "Updated Name",
		Picture:       "https://example.com/new-picture.jpg",
	}

	// Create existing user with identity
	userID := value_object.NewUserID()
	existingUser := entity.NewUser(
		userID,
		"existinguser@example.com",
		"Existing User",
		"https://example.com/old-picture.jpg",
	)

	identityID := value_object.NewIdentityID()
	identity, _ := entity.NewIdentity(
		identityID,
		userID,
		"github",
		"github-user-456",
		"existinguser@example.com",
		"Existing User",
		"https://example.com/old-picture.jpg",
	)
	_ = existingUser.AddIdentity(identity)

	// Repository will find the existing user
	mockRepo.On("FindByProviderIdentity", ctx, "github", "github-user-456").
		Return(existingUser, nil)

	// Save will be called to update LastLogin and identity info
	mockRepo.On("Save", ctx, mock.AnythingOfType("*entity.User")).
		Return(nil)

	// SessionService will be called to create session
	expectedSession := &dto.SessionResult{
		AccessToken:  "mock-access-token",
		RefreshToken: "mock-refresh-token",
		User: dto.UserInfo{
			UserID:      userID.String(),
			Email:       "existinguser@example.com",
			DisplayName: "Updated Name",
			PictureURL:  "https://example.com/new-picture.jpg",
		},
	}
	mockSessionService.On("CreateSession", ctx, userID, false).
		Return(expectedSession, nil)

	// Act
	result, err := usecase.Execute(ctx, providerInfo)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "mock-access-token", result.AccessToken)
	assert.Equal(t, "mock-refresh-token", result.RefreshToken)
	assert.Equal(t, userID.String(), result.User.UserID)
	assert.Equal(t, "existinguser@example.com", result.User.Email)
	mockRepo.AssertExpectations(t)
	mockSessionService.AssertExpectations(t)
}

func TestIssueSessionFromProviderUser_Execute_InvalidProviderInfo_EmptyProvider(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mockRepo := new(MockUserRepository)
	mockSessionService := new(MockSessionService)
	usecase := NewIssueSessionFromProviderUser(mockRepo, mockSessionService)

	providerInfo := dto.ProviderUserInfo{
		Provider:      "", // Invalid: empty provider
		Subject:       "user-123",
		Email:         "user@example.com",
		EmailVerified: true,
		Name:          "User Name",
		Picture:       "https://example.com/picture.jpg",
	}

	// Act
	result, err := usecase.Execute(ctx, providerInfo)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)

	// Once implementation is complete, this assertion should pass:
	// var appErr *application_err.ApplicationError
	// assert.ErrorAs(t, err, &appErr)
	// assert.Equal(t, application_err.ErrInvalidInput, appErr.Kind())

	// No repository calls should be made
	mockRepo.AssertNotCalled(t, "FindByProviderIdentity")
	mockRepo.AssertNotCalled(t, "Save")
}

func TestIssueSessionFromProviderUser_Execute_InvalidProviderInfo_EmptySubject(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mockRepo := new(MockUserRepository)
	mockSessionService := new(MockSessionService)
	usecase := NewIssueSessionFromProviderUser(mockRepo, mockSessionService)

	providerInfo := dto.ProviderUserInfo{
		Provider:      "google",
		Subject:       "", // Invalid: empty subject
		Email:         "user@example.com",
		EmailVerified: true,
		Name:          "User Name",
		Picture:       "https://example.com/picture.jpg",
	}

	// Act
	result, err := usecase.Execute(ctx, providerInfo)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)

	// No repository calls should be made
	mockRepo.AssertNotCalled(t, "FindByProviderIdentity")
	mockRepo.AssertNotCalled(t, "Save")
}

func TestIssueSessionFromProviderUser_Execute_InvalidProviderInfo_EmptyEmail(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mockRepo := new(MockUserRepository)
	mockSessionService := new(MockSessionService)
	usecase := NewIssueSessionFromProviderUser(mockRepo, mockSessionService)

	providerInfo := dto.ProviderUserInfo{
		Provider:      "google",
		Subject:       "user-123",
		Email:         "", // Invalid: empty email
		EmailVerified: true,
		Name:          "User Name",
		Picture:       "https://example.com/picture.jpg",
	}

	// Act
	result, err := usecase.Execute(ctx, providerInfo)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)

	// No repository calls should be made
	mockRepo.AssertNotCalled(t, "FindByProviderIdentity")
	mockRepo.AssertNotCalled(t, "Save")
}

func TestIssueSessionFromProviderUser_Execute_DataAccessFailure(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mockRepo := new(MockUserRepository)
	mockSessionService := new(MockSessionService)
	usecase := NewIssueSessionFromProviderUser(mockRepo, mockSessionService)

	providerInfo := dto.ProviderUserInfo{
		Provider:      "google",
		Subject:       "google-user-123",
		Email:         "user@example.com",
		EmailVerified: true,
		Name:          "User Name",
		Picture:       "https://example.com/picture.jpg",
	}

	// FindByProviderIdentity fails with database error
	mockRepo.On("FindByProviderIdentity", ctx, "google", "google-user-123").
		Return(nil, domain_err.NewDataAccessFailureError("database connection error"))

	// Act
	result, err := usecase.Execute(ctx, providerInfo)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)

	// No other calls should be made
	mockRepo.AssertNotCalled(t, "Save")
	mockSessionService.AssertNotCalled(t, "CreateSession")
}

func TestIssueSessionFromProviderUser_Execute_RepositorySaveError(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mockRepo := new(MockUserRepository)
	mockSessionService := new(MockSessionService)
	usecase := NewIssueSessionFromProviderUser(mockRepo, mockSessionService)

	providerInfo := dto.ProviderUserInfo{
		Provider:      "google",
		Subject:       "google-user-999",
		Email:         "user@example.com",
		EmailVerified: true,
		Name:          "User Name",
		Picture:       "https://example.com/picture.jpg",
	}

	// User does not exist
	mockRepo.On("FindByProviderIdentity", ctx, "google", "google-user-999").
		Return(nil, domain_err.NewNotFoundUserError("google-user-999"))

	// Save fails with a persistence error
	mockRepo.On("Save", ctx, mock.AnythingOfType("*entity.User")).
		Return(domain_err.NewDataPersistFailureError("database error"))

	// Act
	result, err := usecase.Execute(ctx, providerInfo)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)

	// SessionService should not be called if user save fails
	mockSessionService.AssertNotCalled(t, "CreateSession")
	mockRepo.AssertExpectations(t)
}

func TestIssueSessionFromProviderUser_Execute_UniqueConstraintViolation(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mockRepo := new(MockUserRepository)
	mockSessionService := new(MockSessionService)
	usecase := NewIssueSessionFromProviderUser(mockRepo, mockSessionService)

	providerInfo := dto.ProviderUserInfo{
		Provider:      "google",
		Subject:       "google-user-conflict",
		Email:         "conflict@example.com",
		EmailVerified: true,
		Name:          "Conflict User",
		Picture:       "https://example.com/picture.jpg",
	}

	// User does not exist
	mockRepo.On("FindByProviderIdentity", ctx, "google", "google-user-conflict").
		Return(nil, domain_err.NewNotFoundUserError("google-user-conflict"))

	// Save fails with a conflict error (e.g., email already exists)
	mockRepo.On("Save", ctx, mock.AnythingOfType("*entity.User")).
		Return(domain_err.NewDataConflictError("email", "conflict@example.com"))

	// Act
	result, err := usecase.Execute(ctx, providerInfo)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)

	// SessionService should not be called if user save fails
	mockSessionService.AssertNotCalled(t, "CreateSession")
	mockRepo.AssertExpectations(t)
}

func TestIssueSessionFromProviderUser_Execute_SessionCreationFailure(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mockRepo := new(MockUserRepository)
	mockSessionService := new(MockSessionService)
	usecase := NewIssueSessionFromProviderUser(mockRepo, mockSessionService)

	providerInfo := dto.ProviderUserInfo{
		Provider:      "google",
		Subject:       "google-user-123",
		Email:         "user@example.com",
		EmailVerified: true,
		Name:          "User Name",
		Picture:       "https://example.com/picture.jpg",
	}

	// User does not exist
	mockRepo.On("FindByProviderIdentity", ctx, "google", "google-user-123").
		Return(nil, domain_err.NewNotFoundUserError("google-user-123"))

	// User save succeeds
	mockRepo.On("Save", ctx, mock.AnythingOfType("*entity.User")).
		Return(nil)

	// SessionService fails to create session (e.g., JWT generation failure)
	// SessionService returns a generic error that the usecase will wrap as ErrAuthenticationFailed
	mockSessionService.On("CreateSession", ctx, mock.AnythingOfType("value_object.UserID"), false).
		Return(nil, assert.AnError)

	// Act
	result, err := usecase.Execute(ctx, providerInfo)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)

	mockRepo.AssertExpectations(t)
	mockSessionService.AssertExpectations(t)
}
