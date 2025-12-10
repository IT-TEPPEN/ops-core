package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"opscore/backend/internal/document/application/dto"
	"opscore/backend/internal/document/application/usecase"
	"opscore/backend/internal/document/interfaces/api/schema"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestVariableHandler_GetVariableDefinitions(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("正常に変数定義を取得できる", func(t *testing.T) {
		mockUseCase := new(usecase.MockVariableUseCase)
		mockLogger := new(MockLogger)
		handler := NewVariableHandler(mockUseCase, mockLogger)

		docID := "a1b2c3d4-e5f6-7890-1234-567890abcdef"
		expectedVars := []dto.VariableDefinitionDTO{
			{
				Name:         "server_name",
				Label:        "Server Name",
				Description:  "Target server name",
				Type:         "string",
				Required:     true,
				DefaultValue: "localhost",
			},
		}

		// Mock logger calls
		mockLogger.On("Info", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Maybe()
		mockLogger.On("Info", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Maybe()

		// Mock usecase call
		mockUseCase.On("GetVariableDefinitions", mock.Anything, docID).Return(expectedVars, nil)

		// Create request
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = gin.Params{{Key: "docId", Value: docID}}
		c.Request = httptest.NewRequest("GET", "/api/v1/documents/"+docID+"/variables", nil)

		// Call handler
		handler.GetVariableDefinitions(c)

		// Assert response
		assert.Equal(t, http.StatusOK, w.Code)

		var response schema.GetVariableDefinitionsResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, docID, response.DocumentID)
		assert.Equal(t, 1, len(response.Variables))
		assert.Equal(t, "server_name", response.Variables[0].Name)

		mockUseCase.AssertExpectations(t)
	})

	t.Run("ドキュメントIDが空の場合はエラーを返す", func(t *testing.T) {
		mockUseCase := new(usecase.MockVariableUseCase)
		mockLogger := new(MockLogger)
		handler := NewVariableHandler(mockUseCase, mockLogger)

		// Mock logger calls
		mockLogger.On("Warn", mock.Anything, mock.Anything, mock.Anything).Maybe()

		// Create request without docId
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest("GET", "/api/v1/documents//variables", nil)

		// Call handler
		handler.GetVariableDefinitions(c)

		// Assert response
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

func TestVariableHandler_ValidateVariableValues(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("有効な変数値の場合は成功を返す", func(t *testing.T) {
		mockUseCase := new(usecase.MockVariableUseCase)
		mockLogger := new(MockLogger)
		handler := NewVariableHandler(mockUseCase, mockLogger)

		docID := "a1b2c3d4-e5f6-7890-1234-567890abcdef"
		reqBody := schema.ValidateVariableValuesRequest{
			Values: []schema.VariableValueDTO{
				{Name: "server_name", Value: "prod-server"},
			},
		}

		// Mock logger calls
		mockLogger.On("Info", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Maybe()

		// Mock usecase call - validation succeeds (returns nil)
		mockUseCase.On("ValidateVariableValues", mock.Anything, docID, mock.AnythingOfType("[]usecase.VariableValue")).Return(nil)

		// Create request
		body, _ := json.Marshal(reqBody)
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = gin.Params{{Key: "docId", Value: docID}}
		c.Request = httptest.NewRequest("POST", "/api/v1/documents/"+docID+"/validate-variables", bytes.NewBuffer(body))
		c.Request.Header.Set("Content-Type", "application/json")

		// Call handler
		handler.ValidateVariableValues(c)

		// Assert response
		assert.Equal(t, http.StatusOK, w.Code)

		var response schema.ValidateVariableValuesResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.True(t, response.Valid)
		assert.Equal(t, 0, len(response.Errors))

		mockUseCase.AssertExpectations(t)
	})

	t.Run("無効なリクエストボディの場合はエラーを返す", func(t *testing.T) {
		mockUseCase := new(usecase.MockVariableUseCase)
		mockLogger := new(MockLogger)
		handler := NewVariableHandler(mockUseCase, mockLogger)

		docID := "a1b2c3d4-e5f6-7890-1234-567890abcdef"

		// Mock logger calls
		mockLogger.On("Error", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Maybe()

		// Create request with invalid JSON
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = gin.Params{{Key: "docId", Value: docID}}
		c.Request = httptest.NewRequest("POST", "/api/v1/documents/"+docID+"/validate-variables", bytes.NewBufferString("invalid json"))
		c.Request.Header.Set("Content-Type", "application/json")

		// Call handler
		handler.ValidateVariableValues(c)

		// Assert response
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("バリデーションエラーの場合は適切なレスポンスを返す", func(t *testing.T) {
		mockUseCase := new(usecase.MockVariableUseCase)
		mockLogger := new(MockLogger)
		handler := NewVariableHandler(mockUseCase, mockLogger)

		docID := "a1b2c3d4-e5f6-7890-1234-567890abcdef"
		reqBody := schema.ValidateVariableValuesRequest{
			Values: []schema.VariableValueDTO{
				{Name: "server_name", Value: ""},
			},
		}

		// Mock logger calls
		mockLogger.On("Info", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Maybe()

		// Mock validation error
		validationErr := apperror.NewValidationFailedError([]apperror.FieldError{
			{Field: "server_name", Message: "Server Name is required"},
		})
		mockUseCase.On("ValidateVariableValues", mock.Anything, docID, mock.AnythingOfType("[]usecase.VariableValue")).Return(validationErr)

		// Create request
		body, _ := json.Marshal(reqBody)
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = gin.Params{{Key: "docId", Value: docID}}
		c.Request = httptest.NewRequest("POST", "/api/v1/documents/"+docID+"/validate-variables", bytes.NewBuffer(body))
		c.Request.Header.Set("Content-Type", "application/json")

		// Call handler
		handler.ValidateVariableValues(c)

		// Assert response - validation failures return 200 OK with valid: false
		assert.Equal(t, http.StatusOK, w.Code)

		var response schema.ValidateVariableValuesResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.False(t, response.Valid)
		assert.Greater(t, len(response.Errors), 0)
		assert.Equal(t, "server_name", response.Errors[0].Name)

		mockUseCase.AssertExpectations(t)
	})
}
