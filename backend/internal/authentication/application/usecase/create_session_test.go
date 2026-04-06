package usecase

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"opscore/backend/internal/authentication/application/dto"
	application_err "opscore/backend/internal/authentication/application/err"
	"opscore/backend/internal/authentication/domain/entity"
	domain_err "opscore/backend/internal/authentication/domain/err"
	"opscore/backend/internal/authentication/domain/value_object"
)

func TestCreateSession_Execute_ShortSession(t *testing.T) {
	ctx := context.Background()
	mockUserRepo := new(MockUserRepository)
	mockSessionRepo := new(MockSessionRepository)
	uc := NewCreateSession(mockUserRepo, mockSessionRepo)

	userID := value_object.NewUserID()
	user := entity.NewUser(userID, "user@example.com", "Test User", "https://example.com/pic.jpg")

	mockUserRepo.On("FindByID", ctx, userID).Return(user, nil)
	mockSessionRepo.On("Create", ctx, mock.AnythingOfType("*entity.Session")).Return(nil)

	req := dto.CreateSessionRequest{
		UserID:      userID.String(),
		SessionType: dto.SessionTypeShort,
	}

	result, err := uc.Execute(ctx, req)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, userID.String(), result.UserID)
	assert.Equal(t, dto.SessionTypeShort, result.SessionType)
	assert.False(t, result.IsRevoked)
	assert.NotEmpty(t, result.SessionID)
	assert.True(t, result.ExpiresAt.After(result.CreatedAt))
	mockUserRepo.AssertExpectations(t)
	mockSessionRepo.AssertExpectations(t)
}

func TestCreateSession_Execute_LongSession(t *testing.T) {
	ctx := context.Background()
	mockUserRepo := new(MockUserRepository)
	mockSessionRepo := new(MockSessionRepository)
	uc := NewCreateSession(mockUserRepo, mockSessionRepo)

	userID := value_object.NewUserID()
	user := entity.NewUser(userID, "user@example.com", "Test User", "https://example.com/pic.jpg")

	mockUserRepo.On("FindByID", ctx, userID).Return(user, nil)
	mockSessionRepo.On("Create", ctx, mock.AnythingOfType("*entity.Session")).Return(nil)

	req := dto.CreateSessionRequest{
		UserID:      userID.String(),
		SessionType: dto.SessionTypeLong,
	}

	result, err := uc.Execute(ctx, req)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, dto.SessionTypeLong, result.SessionType)
	// Long session should have a later expiry than short session
	assert.True(t, result.ExpiresAt.After(result.CreatedAt))
	mockUserRepo.AssertExpectations(t)
	mockSessionRepo.AssertExpectations(t)
}

func TestCreateSession_Execute_InvalidUserID(t *testing.T) {
	ctx := context.Background()
	mockUserRepo := new(MockUserRepository)
	mockSessionRepo := new(MockSessionRepository)
	uc := NewCreateSession(mockUserRepo, mockSessionRepo)

	req := dto.CreateSessionRequest{
		UserID:      "not-a-uuid",
		SessionType: dto.SessionTypeShort,
	}

	result, err := uc.Execute(ctx, req)

	assert.Error(t, err)
	assert.Nil(t, result)
	var appErr *application_err.ApplicationError
	assert.ErrorAs(t, err, &appErr)
	assert.Equal(t, application_err.ErrInvalidInput, appErr.Kind())
	mockUserRepo.AssertNotCalled(t, "FindByID")
	mockSessionRepo.AssertNotCalled(t, "Create")
}

func TestCreateSession_Execute_EmptyUserID(t *testing.T) {
	ctx := context.Background()
	mockUserRepo := new(MockUserRepository)
	mockSessionRepo := new(MockSessionRepository)
	uc := NewCreateSession(mockUserRepo, mockSessionRepo)

	req := dto.CreateSessionRequest{
		UserID:      "",
		SessionType: dto.SessionTypeShort,
	}

	result, err := uc.Execute(ctx, req)

	assert.Error(t, err)
	assert.Nil(t, result)
	var appErr *application_err.ApplicationError
	assert.ErrorAs(t, err, &appErr)
	assert.Equal(t, application_err.ErrInvalidInput, appErr.Kind())
}

func TestCreateSession_Execute_InvalidSessionType(t *testing.T) {
	ctx := context.Background()
	mockUserRepo := new(MockUserRepository)
	mockSessionRepo := new(MockSessionRepository)
	uc := NewCreateSession(mockUserRepo, mockSessionRepo)

	userID := value_object.NewUserID()

	req := dto.CreateSessionRequest{
		UserID:      userID.String(),
		SessionType: "invalid",
	}

	result, err := uc.Execute(ctx, req)

	assert.Error(t, err)
	assert.Nil(t, result)
	var appErr *application_err.ApplicationError
	assert.ErrorAs(t, err, &appErr)
	assert.Equal(t, application_err.ErrInvalidInput, appErr.Kind())
	mockUserRepo.AssertNotCalled(t, "FindByID")
}

func TestCreateSession_Execute_UserNotFound(t *testing.T) {
	ctx := context.Background()
	mockUserRepo := new(MockUserRepository)
	mockSessionRepo := new(MockSessionRepository)
	uc := NewCreateSession(mockUserRepo, mockSessionRepo)

	userID := value_object.NewUserID()

	mockUserRepo.On("FindByID", ctx, userID).
		Return(nil, domain_err.NewNotFoundUserError(userID.String()))

	req := dto.CreateSessionRequest{
		UserID:      userID.String(),
		SessionType: dto.SessionTypeShort,
	}

	result, err := uc.Execute(ctx, req)

	assert.Error(t, err)
	assert.Nil(t, result)
	var appErr *application_err.ApplicationError
	assert.ErrorAs(t, err, &appErr)
	assert.Equal(t, application_err.ErrUserNotFound, appErr.Kind())
	mockSessionRepo.AssertNotCalled(t, "Create")
}

func TestCreateSession_Execute_SessionPersistFailure(t *testing.T) {
	ctx := context.Background()
	mockUserRepo := new(MockUserRepository)
	mockSessionRepo := new(MockSessionRepository)
	uc := NewCreateSession(mockUserRepo, mockSessionRepo)

	userID := value_object.NewUserID()
	user := entity.NewUser(userID, "user@example.com", "Test User", "https://example.com/pic.jpg")

	mockUserRepo.On("FindByID", ctx, userID).Return(user, nil)
	mockSessionRepo.On("Create", ctx, mock.AnythingOfType("*entity.Session")).
		Return(domain_err.NewDataPersistFailureError("sessions"))

	req := dto.CreateSessionRequest{
		UserID:      userID.String(),
		SessionType: dto.SessionTypeShort,
	}

	result, err := uc.Execute(ctx, req)

	assert.Error(t, err)
	assert.Nil(t, result)
	var appErr *application_err.ApplicationError
	assert.ErrorAs(t, err, &appErr)
	assert.Equal(t, application_err.ErrDataPersistFailure, appErr.Kind())
}
