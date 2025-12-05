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

// MockGroupUseCase is a mock implementation of the GroupUseCase interface
type MockGroupUseCase struct {
	mock.Mock
}

func (m *MockGroupUseCase) Create(ctx context.Context, req dto.CreateGroupRequest) (*dto.GroupResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.GroupResponse), args.Error(1)
}

func (m *MockGroupUseCase) GetByID(ctx context.Context, groupID string) (*dto.GroupResponse, error) {
	args := m.Called(ctx, groupID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.GroupResponse), args.Error(1)
}

func (m *MockGroupUseCase) GetAll(ctx context.Context) ([]dto.GroupResponse, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]dto.GroupResponse), args.Error(1)
}

func (m *MockGroupUseCase) Update(ctx context.Context, groupID string, req dto.UpdateGroupRequest) (*dto.GroupResponse, error) {
	args := m.Called(ctx, groupID, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.GroupResponse), args.Error(1)
}

func (m *MockGroupUseCase) Delete(ctx context.Context, groupID string) error {
	args := m.Called(ctx, groupID)
	return args.Error(0)
}

func (m *MockGroupUseCase) AddMember(ctx context.Context, groupID string, req dto.AddMemberRequest) (*dto.GroupResponse, error) {
	args := m.Called(ctx, groupID, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.GroupResponse), args.Error(1)
}

func (m *MockGroupUseCase) RemoveMember(ctx context.Context, groupID string, req dto.RemoveMemberRequest) (*dto.GroupResponse, error) {
	args := m.Called(ctx, groupID, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.GroupResponse), args.Error(1)
}

func (m *MockGroupUseCase) GetByMemberID(ctx context.Context, userID string) ([]dto.GroupResponse, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]dto.GroupResponse), args.Error(1)
}

func setupGroupRouter(handler *GroupHandler) *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.POST("/groups", handler.CreateGroup)
	router.GET("/groups", handler.ListGroups)
	router.GET("/groups/:groupId", handler.GetGroup)
	router.PUT("/groups/:groupId", handler.UpdateGroup)
	router.DELETE("/groups/:groupId", handler.DeleteGroup)
	router.POST("/groups/:groupId/members", handler.AddMember)
	router.DELETE("/groups/:groupId/members", handler.RemoveMember)
	router.GET("/users/:userId/groups", handler.GetUserGroups)
	return router
}

func TestGroupHandler_CreateGroup(t *testing.T) {
	t.Run("有効なリクエストでグループが作成される", func(t *testing.T) {
		mockUseCase := new(MockGroupUseCase)
		handler := NewGroupHandler(mockUseCase, &MockLogger{})
		router := setupGroupRouter(handler)

		expectedResponse := &dto.GroupResponse{
			ID:          "group-123",
			Name:        "Test Group",
			Description: "Test Description",
			MemberIDs:   []string{},
		}

		mockUseCase.On("Create", mock.Anything, mock.MatchedBy(func(req dto.CreateGroupRequest) bool {
			return req.Name == "Test Group" && req.Description == "Test Description"
		})).Return(expectedResponse, nil)

		body := schema.CreateGroupRequest{
			Name:        "Test Group",
			Description: "Test Description",
		}
		jsonBody, _ := json.Marshal(body)

		req, _ := http.NewRequest("POST", "/groups", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)

		var response schema.GroupResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "group-123", response.ID)
		assert.Equal(t, "Test Group", response.Name)
		mockUseCase.AssertExpectations(t)
	})

	t.Run("無効なリクエストボディで400が返される", func(t *testing.T) {
		mockUseCase := new(MockGroupUseCase)
		handler := NewGroupHandler(mockUseCase, &MockLogger{})
		router := setupGroupRouter(handler)

		req, _ := http.NewRequest("POST", "/groups", bytes.NewBuffer([]byte("invalid json")))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

func TestGroupHandler_GetGroup(t *testing.T) {
	t.Run("存在するグループが取得できる", func(t *testing.T) {
		mockUseCase := new(MockGroupUseCase)
		handler := NewGroupHandler(mockUseCase, &MockLogger{})
		router := setupGroupRouter(handler)

		expectedResponse := &dto.GroupResponse{
			ID:          "group-123",
			Name:        "Test Group",
			Description: "Test Description",
			MemberIDs:   []string{},
		}

		mockUseCase.On("GetByID", mock.Anything, "group-123").Return(expectedResponse, nil)

		req, _ := http.NewRequest("GET", "/groups/group-123", nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response schema.GroupResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "group-123", response.ID)
		mockUseCase.AssertExpectations(t)
	})
}

func TestGroupHandler_ListGroups(t *testing.T) {
	t.Run("グループ一覧が取得できる", func(t *testing.T) {
		mockUseCase := new(MockGroupUseCase)
		handler := NewGroupHandler(mockUseCase, &MockLogger{})
		router := setupGroupRouter(handler)

		expectedResponses := []dto.GroupResponse{
			{ID: "group-1", Name: "Group 1", Description: "Desc 1", MemberIDs: []string{}},
			{ID: "group-2", Name: "Group 2", Description: "Desc 2", MemberIDs: []string{}},
		}

		mockUseCase.On("GetAll", mock.Anything).Return(expectedResponses, nil)

		req, _ := http.NewRequest("GET", "/groups", nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response schema.ListGroupsResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Len(t, response.Groups, 2)
		mockUseCase.AssertExpectations(t)
	})
}

func TestGroupHandler_UpdateGroup(t *testing.T) {
	t.Run("有効なリクエストでグループが更新される", func(t *testing.T) {
		mockUseCase := new(MockGroupUseCase)
		handler := NewGroupHandler(mockUseCase, &MockLogger{})
		router := setupGroupRouter(handler)

		expectedResponse := &dto.GroupResponse{
			ID:          "group-123",
			Name:        "New Name",
			Description: "New Description",
			MemberIDs:   []string{},
		}

		mockUseCase.On("Update", mock.Anything, "group-123", mock.MatchedBy(func(req dto.UpdateGroupRequest) bool {
			return req.Name == "New Name" && req.Description == "New Description"
		})).Return(expectedResponse, nil)

		body := schema.UpdateGroupRequest{
			Name:        "New Name",
			Description: "New Description",
		}
		jsonBody, _ := json.Marshal(body)

		req, _ := http.NewRequest("PUT", "/groups/group-123", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		mockUseCase.AssertExpectations(t)
	})
}

func TestGroupHandler_DeleteGroup(t *testing.T) {
	t.Run("存在するグループが削除できる", func(t *testing.T) {
		mockUseCase := new(MockGroupUseCase)
		handler := NewGroupHandler(mockUseCase, &MockLogger{})
		router := setupGroupRouter(handler)

		mockUseCase.On("Delete", mock.Anything, "group-123").Return(nil)

		req, _ := http.NewRequest("DELETE", "/groups/group-123", nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		mockUseCase.AssertExpectations(t)
	})
}

func TestGroupHandler_AddMember(t *testing.T) {
	t.Run("グループにメンバーを追加できる", func(t *testing.T) {
		mockUseCase := new(MockGroupUseCase)
		handler := NewGroupHandler(mockUseCase, &MockLogger{})
		router := setupGroupRouter(handler)

		expectedResponse := &dto.GroupResponse{
			ID:          "group-123",
			Name:        "Test Group",
			Description: "Test Description",
			MemberIDs:   []string{"user-123"},
		}

		mockUseCase.On("AddMember", mock.Anything, "group-123", mock.MatchedBy(func(req dto.AddMemberRequest) bool {
			return req.UserID == "user-123"
		})).Return(expectedResponse, nil)

		body := schema.AddMemberRequest{
			UserID: "user-123",
		}
		jsonBody, _ := json.Marshal(body)

		req, _ := http.NewRequest("POST", "/groups/group-123/members", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		mockUseCase.AssertExpectations(t)
	})
}

func TestGroupHandler_RemoveMember(t *testing.T) {
	t.Run("グループからメンバーを削除できる", func(t *testing.T) {
		mockUseCase := new(MockGroupUseCase)
		handler := NewGroupHandler(mockUseCase, &MockLogger{})
		router := setupGroupRouter(handler)

		expectedResponse := &dto.GroupResponse{
			ID:          "group-123",
			Name:        "Test Group",
			Description: "Test Description",
			MemberIDs:   []string{},
		}

		mockUseCase.On("RemoveMember", mock.Anything, "group-123", mock.MatchedBy(func(req dto.RemoveMemberRequest) bool {
			return req.UserID == "user-123"
		})).Return(expectedResponse, nil)

		body := schema.RemoveMemberRequest{
			UserID: "user-123",
		}
		jsonBody, _ := json.Marshal(body)

		req, _ := http.NewRequest("DELETE", "/groups/group-123/members", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		mockUseCase.AssertExpectations(t)
	})
}

func TestGroupHandler_GetUserGroups(t *testing.T) {
	t.Run("ユーザーのグループ一覧が取得できる", func(t *testing.T) {
		mockUseCase := new(MockGroupUseCase)
		handler := NewGroupHandler(mockUseCase, &MockLogger{})
		router := setupGroupRouter(handler)

		expectedResponses := []dto.GroupResponse{
			{ID: "group-1", Name: "Group 1", Description: "Desc 1", MemberIDs: []string{}},
		}

		mockUseCase.On("GetByMemberID", mock.Anything, "user-123").Return(expectedResponses, nil)

		req, _ := http.NewRequest("GET", "/users/user-123/groups", nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response schema.ListGroupsResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Len(t, response.Groups, 1)
		mockUseCase.AssertExpectations(t)
	})
}
