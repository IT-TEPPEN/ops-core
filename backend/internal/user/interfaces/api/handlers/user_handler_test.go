package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"opscore/backend/internal/user/application/dto"
	"opscore/backend/internal/user/interfaces/api/schema"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockUserUseCase is a mock implementation of the UserUseCase interface
type MockUserUseCase struct {
	mock.Mock
}

func (m *MockUserUseCase) Create(ctx context.Context, req dto.CreateUserRequest) (*dto.UserResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.UserResponse), args.Error(1)
}

func (m *MockUserUseCase) GetByID(ctx context.Context, userID string) (*dto.UserResponse, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.UserResponse), args.Error(1)
}

func (m *MockUserUseCase) GetAll(ctx context.Context) ([]dto.UserResponse, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]dto.UserResponse), args.Error(1)
}

func (m *MockUserUseCase) Update(ctx context.Context, userID string, req dto.UpdateUserRequest) (*dto.UserResponse, error) {
	args := m.Called(ctx, userID, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.UserResponse), args.Error(1)
}

func (m *MockUserUseCase) Delete(ctx context.Context, userID string) error {
	args := m.Called(ctx, userID)
	return args.Error(0)
}

func (m *MockUserUseCase) ChangeRole(ctx context.Context, userID string, req dto.ChangeRoleRequest) (*dto.UserResponse, error) {
	args := m.Called(ctx, userID, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.UserResponse), args.Error(1)
}

func (m *MockUserUseCase) JoinGroup(ctx context.Context, userID string, req dto.JoinGroupRequest) (*dto.UserResponse, error) {
	args := m.Called(ctx, userID, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.UserResponse), args.Error(1)
}

func (m *MockUserUseCase) LeaveGroup(ctx context.Context, userID string, req dto.LeaveGroupRequest) (*dto.UserResponse, error) {
	args := m.Called(ctx, userID, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.UserResponse), args.Error(1)
}

func setupUserRouter(handler *UserHandler) *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.POST("/users", handler.CreateUser)
	router.GET("/users", handler.ListUsers)
	router.GET("/users/:userId", handler.GetUser)
	router.PUT("/users/:userId", handler.UpdateUser)
	router.DELETE("/users/:userId", handler.DeleteUser)
	router.PUT("/users/:userId/role", handler.ChangeUserRole)
	router.POST("/users/:userId/groups", handler.JoinGroup)
	router.DELETE("/users/:userId/groups", handler.LeaveGroup)
	return router
}

func TestUserHandler_CreateUser(t *testing.T) {
	t.Run("有効なリクエストでユーザーが作成される", func(t *testing.T) {
		mockUseCase := new(MockUserUseCase)
		handler := NewUserHandler(mockUseCase, &MockLogger{})
		router := setupUserRouter(handler)

		expectedResponse := &dto.UserResponse{
			ID:       "user-123",
			Name:     "Test User",
			Email:    "test@example.com",
			Role:     "user",
			GroupIDs: []string{},
		}

		mockUseCase.On("Create", mock.Anything, mock.MatchedBy(func(req dto.CreateUserRequest) bool {
			return req.Name == "Test User" && req.Email == "test@example.com"
		})).Return(expectedResponse, nil)

		body := schema.CreateUserRequest{
			Name:  "Test User",
			Email: "test@example.com",
			Role:  "user",
		}
		jsonBody, _ := json.Marshal(body)

		req, _ := http.NewRequest("POST", "/users", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)

		var response schema.UserResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "user-123", response.ID)
		assert.Equal(t, "Test User", response.Name)
		mockUseCase.AssertExpectations(t)
	})

	t.Run("無効なリクエストボディで400が返される", func(t *testing.T) {
		mockUseCase := new(MockUserUseCase)
		handler := NewUserHandler(mockUseCase, &MockLogger{})
		router := setupUserRouter(handler)

		req, _ := http.NewRequest("POST", "/users", bytes.NewBuffer([]byte("invalid json")))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

func TestUserHandler_GetUser(t *testing.T) {
	t.Run("存在するユーザーが取得できる", func(t *testing.T) {
		mockUseCase := new(MockUserUseCase)
		handler := NewUserHandler(mockUseCase, &MockLogger{})
		router := setupUserRouter(handler)

		expectedResponse := &dto.UserResponse{
			ID:       "user-123",
			Name:     "Test User",
			Email:    "test@example.com",
			Role:     "user",
			GroupIDs: []string{},
		}

		mockUseCase.On("GetByID", mock.Anything, "user-123").Return(expectedResponse, nil)

		req, _ := http.NewRequest("GET", "/users/user-123", nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response schema.UserResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "user-123", response.ID)
		mockUseCase.AssertExpectations(t)
	})
}

func TestUserHandler_ListUsers(t *testing.T) {
	t.Run("ユーザー一覧が取得できる", func(t *testing.T) {
		mockUseCase := new(MockUserUseCase)
		handler := NewUserHandler(mockUseCase, &MockLogger{})
		router := setupUserRouter(handler)

		expectedResponses := []dto.UserResponse{
			{ID: "user-1", Name: "User 1", Email: "user1@example.com", Role: "user", GroupIDs: []string{}},
			{ID: "user-2", Name: "User 2", Email: "user2@example.com", Role: "admin", GroupIDs: []string{}},
		}

		mockUseCase.On("GetAll", mock.Anything).Return(expectedResponses, nil)

		req, _ := http.NewRequest("GET", "/users", nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response schema.ListUsersResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Len(t, response.Users, 2)
		mockUseCase.AssertExpectations(t)
	})
}

func TestUserHandler_UpdateUser(t *testing.T) {
	t.Run("有効なリクエストでユーザーが更新される", func(t *testing.T) {
		mockUseCase := new(MockUserUseCase)
		handler := NewUserHandler(mockUseCase, &MockLogger{})
		router := setupUserRouter(handler)

		expectedResponse := &dto.UserResponse{
			ID:       "user-123",
			Name:     "New Name",
			Email:    "new@example.com",
			Role:     "user",
			GroupIDs: []string{},
		}

		mockUseCase.On("Update", mock.Anything, "user-123", mock.MatchedBy(func(req dto.UpdateUserRequest) bool {
			return req.Name == "New Name" && req.Email == "new@example.com"
		})).Return(expectedResponse, nil)

		body := schema.UpdateUserRequest{
			Name:  "New Name",
			Email: "new@example.com",
		}
		jsonBody, _ := json.Marshal(body)

		req, _ := http.NewRequest("PUT", "/users/user-123", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		mockUseCase.AssertExpectations(t)
	})
}

func TestUserHandler_DeleteUser(t *testing.T) {
	t.Run("存在するユーザーが削除できる", func(t *testing.T) {
		mockUseCase := new(MockUserUseCase)
		handler := NewUserHandler(mockUseCase, &MockLogger{})
		router := setupUserRouter(handler)

		mockUseCase.On("Delete", mock.Anything, "user-123").Return(nil)

		req, _ := http.NewRequest("DELETE", "/users/user-123", nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		mockUseCase.AssertExpectations(t)
	})
}

func TestUserHandler_ChangeUserRole(t *testing.T) {
	t.Run("有効なリクエストでロールが変更される", func(t *testing.T) {
		mockUseCase := new(MockUserUseCase)
		handler := NewUserHandler(mockUseCase, &MockLogger{})
		router := setupUserRouter(handler)

		expectedResponse := &dto.UserResponse{
			ID:       "user-123",
			Name:     "Test User",
			Email:    "test@example.com",
			Role:     "admin",
			GroupIDs: []string{},
		}

		mockUseCase.On("ChangeRole", mock.Anything, "user-123", mock.MatchedBy(func(req dto.ChangeRoleRequest) bool {
			return req.Role == "admin"
		})).Return(expectedResponse, nil)

		body := schema.ChangeRoleRequest{
			Role: "admin",
		}
		jsonBody, _ := json.Marshal(body)

		req, _ := http.NewRequest("PUT", "/users/user-123/role", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		mockUseCase.AssertExpectations(t)
	})
}

func TestUserHandler_JoinGroup(t *testing.T) {
	t.Run("ユーザーがグループに参加できる", func(t *testing.T) {
		mockUseCase := new(MockUserUseCase)
		handler := NewUserHandler(mockUseCase, &MockLogger{})
		router := setupUserRouter(handler)

		expectedResponse := &dto.UserResponse{
			ID:       "user-123",
			Name:     "Test User",
			Email:    "test@example.com",
			Role:     "user",
			GroupIDs: []string{"group-123"},
		}

		mockUseCase.On("JoinGroup", mock.Anything, "user-123", mock.MatchedBy(func(req dto.JoinGroupRequest) bool {
			return req.GroupID == "group-123"
		})).Return(expectedResponse, nil)

		body := schema.JoinGroupRequest{
			GroupID: "group-123",
		}
		jsonBody, _ := json.Marshal(body)

		req, _ := http.NewRequest("POST", "/users/user-123/groups", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		mockUseCase.AssertExpectations(t)
	})
}
