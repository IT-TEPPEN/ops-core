package usecase

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"opscore/backend/internal/authentication/application/dto"
	application_err "opscore/backend/internal/authentication/application/err"
	"opscore/backend/internal/authentication/domain/entity"
	domain_err "opscore/backend/internal/authentication/domain/err"
	"opscore/backend/internal/authentication/domain/value_object"
)

// MockSessionRepository is a mock implementation of SessionRepository
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

func (m *MockSessionRepository) FindByUserID(ctx context.Context, userID value_object.UserID) ([]*entity.Session, error) {
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

func (m *MockSessionRepository) RevokeAllByUserID(ctx context.Context, userID value_object.UserID) (int, error) {
	args := m.Called(ctx, userID)
	return args.Int(0), args.Error(1)
}

func (m *MockSessionRepository) DeleteExpired(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

func TestGetAuthenticatedUser_Execute_ValidSession(t *testing.T) {
	ctx := context.Background()
	mockSessionRepo := new(MockSessionRepository)
	mockUserRepo := new(MockUserRepository)
	uc := NewGetAuthenticatedUser(mockSessionRepo, mockUserRepo)

	userID := value_object.NewUserID()
	sessionID := value_object.NewSessionID()
	session := entity.ReconstructSession(
		sessionID, userID, "tokenhash", "jti-123",
		time.Now().Add(7*24*time.Hour), time.Now(),
		nil, nil, false, false,
	)
	user := entity.NewUser(userID, "user@example.com", "Test User", "https://example.com/pic.jpg")

	mockSessionRepo.On("FindByID", ctx, sessionID).Return(session, nil)
	mockUserRepo.On("FindByID", ctx, userID).Return(user, nil)

	req := dto.SessionIDRequest{SessionID: sessionID.String()}
	result, err := uc.Execute(ctx, req)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, userID.String(), result.UserID)
	assert.Equal(t, "user@example.com", result.Email)
	assert.Equal(t, "Test User", result.DisplayName)
	assert.Equal(t, "https://example.com/pic.jpg", result.PictureURL)
	mockSessionRepo.AssertExpectations(t)
	mockUserRepo.AssertExpectations(t)
}

func TestGetAuthenticatedUser_Execute_InvalidSessionID(t *testing.T) {
	ctx := context.Background()
	mockSessionRepo := new(MockSessionRepository)
	mockUserRepo := new(MockUserRepository)
	uc := NewGetAuthenticatedUser(mockSessionRepo, mockUserRepo)

	req := dto.SessionIDRequest{SessionID: "not-a-uuid"}
	result, err := uc.Execute(ctx, req)

	assert.Error(t, err)
	assert.Nil(t, result)
	var appErr *application_err.ApplicationError
	assert.ErrorAs(t, err, &appErr)
	assert.Equal(t, application_err.ErrInvalidInput, appErr.Kind())
	mockSessionRepo.AssertNotCalled(t, "FindByID")
}

func TestGetAuthenticatedUser_Execute_EmptySessionID(t *testing.T) {
	ctx := context.Background()
	mockSessionRepo := new(MockSessionRepository)
	mockUserRepo := new(MockUserRepository)
	uc := NewGetAuthenticatedUser(mockSessionRepo, mockUserRepo)

	req := dto.SessionIDRequest{SessionID: ""}
	result, err := uc.Execute(ctx, req)

	assert.Error(t, err)
	assert.Nil(t, result)
	var appErr *application_err.ApplicationError
	assert.ErrorAs(t, err, &appErr)
	assert.Equal(t, application_err.ErrInvalidInput, appErr.Kind())
}

func TestGetAuthenticatedUser_Execute_SessionNotFound(t *testing.T) {
	ctx := context.Background()
	mockSessionRepo := new(MockSessionRepository)
	mockUserRepo := new(MockUserRepository)
	uc := NewGetAuthenticatedUser(mockSessionRepo, mockUserRepo)

	sessionID := value_object.NewSessionID()
	mockSessionRepo.On("FindByID", ctx, sessionID).
		Return(nil, domain_err.NewNotFoundSessionError(sessionID.String()))

	req := dto.SessionIDRequest{SessionID: sessionID.String()}
	result, err := uc.Execute(ctx, req)

	assert.Error(t, err)
	assert.Nil(t, result)
	var appErr *application_err.ApplicationError
	assert.ErrorAs(t, err, &appErr)
	assert.Equal(t, application_err.ErrSessionNotFound, appErr.Kind())
	mockUserRepo.AssertNotCalled(t, "FindByID")
}

func TestGetAuthenticatedUser_Execute_SessionExpired(t *testing.T) {
	ctx := context.Background()
	mockSessionRepo := new(MockSessionRepository)
	mockUserRepo := new(MockUserRepository)
	uc := NewGetAuthenticatedUser(mockSessionRepo, mockUserRepo)

	userID := value_object.NewUserID()
	sessionID := value_object.NewSessionID()
	// Session expired 1 hour ago
	session := entity.ReconstructSession(
		sessionID, userID, "tokenhash", "jti-123",
		time.Now().Add(-1*time.Hour), time.Now().Add(-8*24*time.Hour),
		nil, nil, false, false,
	)

	mockSessionRepo.On("FindByID", ctx, sessionID).Return(session, nil)

	req := dto.SessionIDRequest{SessionID: sessionID.String()}
	result, err := uc.Execute(ctx, req)

	assert.Error(t, err)
	assert.Nil(t, result)
	var appErr *application_err.ApplicationError
	assert.ErrorAs(t, err, &appErr)
	assert.Equal(t, application_err.ErrSessionExpired, appErr.Kind())
	mockUserRepo.AssertNotCalled(t, "FindByID")
}

func TestGetAuthenticatedUser_Execute_SessionRevoked(t *testing.T) {
	ctx := context.Background()
	mockSessionRepo := new(MockSessionRepository)
	mockUserRepo := new(MockUserRepository)
	uc := NewGetAuthenticatedUser(mockSessionRepo, mockUserRepo)

	userID := value_object.NewUserID()
	sessionID := value_object.NewSessionID()
	revokedAt := time.Now().Add(-1 * time.Hour)
	// Session is revoked but not expired
	session := entity.ReconstructSession(
		sessionID, userID, "tokenhash", "jti-123",
		time.Now().Add(7*24*time.Hour), time.Now(),
		nil, &revokedAt, true, false,
	)

	mockSessionRepo.On("FindByID", ctx, sessionID).Return(session, nil)

	req := dto.SessionIDRequest{SessionID: sessionID.String()}
	result, err := uc.Execute(ctx, req)

	assert.Error(t, err)
	assert.Nil(t, result)
	var appErr *application_err.ApplicationError
	assert.ErrorAs(t, err, &appErr)
	assert.Equal(t, application_err.ErrSessionRevoked, appErr.Kind())
	mockUserRepo.AssertNotCalled(t, "FindByID")
}

func TestGetAuthenticatedUser_Execute_UserNotFound(t *testing.T) {
	ctx := context.Background()
	mockSessionRepo := new(MockSessionRepository)
	mockUserRepo := new(MockUserRepository)
	uc := NewGetAuthenticatedUser(mockSessionRepo, mockUserRepo)

	userID := value_object.NewUserID()
	sessionID := value_object.NewSessionID()
	session := entity.ReconstructSession(
		sessionID, userID, "tokenhash", "jti-123",
		time.Now().Add(7*24*time.Hour), time.Now(),
		nil, nil, false, false,
	)

	mockSessionRepo.On("FindByID", ctx, sessionID).Return(session, nil)
	mockUserRepo.On("FindByID", ctx, userID).
		Return(nil, domain_err.NewNotFoundUserError(userID.String()))

	req := dto.SessionIDRequest{SessionID: sessionID.String()}
	result, err := uc.Execute(ctx, req)

	assert.Error(t, err)
	assert.Nil(t, result)
	var appErr *application_err.ApplicationError
	assert.ErrorAs(t, err, &appErr)
	assert.Equal(t, application_err.ErrUserNotFound, appErr.Kind())
}

func TestGetAuthenticatedUser_Execute_DataAccessFailure(t *testing.T) {
	ctx := context.Background()
	mockSessionRepo := new(MockSessionRepository)
	mockUserRepo := new(MockUserRepository)
	uc := NewGetAuthenticatedUser(mockSessionRepo, mockUserRepo)

	sessionID := value_object.NewSessionID()
	mockSessionRepo.On("FindByID", ctx, sessionID).
		Return(nil, domain_err.NewDataAccessFailureError("sessions"))

	req := dto.SessionIDRequest{SessionID: sessionID.String()}
	result, err := uc.Execute(ctx, req)

	assert.Error(t, err)
	assert.Nil(t, result)
	var appErr *application_err.ApplicationError
	assert.ErrorAs(t, err, &appErr)
	assert.Equal(t, application_err.ErrUnexpected, appErr.Kind())
}
