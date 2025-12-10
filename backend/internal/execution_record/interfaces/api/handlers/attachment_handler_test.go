package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"

	"opscore/backend/internal/execution_record/application/dto"
	"opscore/backend/internal/execution_record/domain/value_object"
)

// MockAttachmentUsecase is a mock implementation of AttachmentUsecase
type MockAttachmentUsecase struct {
	UploadAttachmentFunc          func(ctx context.Context, req *dto.UploadAttachmentRequest) (*dto.AttachmentResponse, error)
	GetAttachmentFunc             func(ctx context.Context, attachmentID string) (*dto.AttachmentResponse, error)
	GetAttachmentFileFunc         func(ctx context.Context, attachmentID string) (io.ReadCloser, *dto.AttachmentResponse, error)
	ListAttachmentsByRecordIDFunc func(ctx context.Context, recordID string) ([]*dto.AttachmentResponse, error)
	ListAttachmentsByStepIDFunc   func(ctx context.Context, stepID string) ([]*dto.AttachmentResponse, error)
	DeleteAttachmentFunc          func(ctx context.Context, attachmentID string) error
	GetAttachmentURLFunc          func(ctx context.Context, attachmentID string, expirationMinutes int) (string, error)
}

func (m *MockAttachmentUsecase) UploadAttachment(ctx context.Context, req *dto.UploadAttachmentRequest) (*dto.AttachmentResponse, error) {
	if m.UploadAttachmentFunc != nil {
		return m.UploadAttachmentFunc(ctx, req)
	}
	return nil, nil
}

func (m *MockAttachmentUsecase) GetAttachment(ctx context.Context, attachmentID string) (*dto.AttachmentResponse, error) {
	if m.GetAttachmentFunc != nil {
		return m.GetAttachmentFunc(ctx, attachmentID)
	}
	return nil, nil
}

func (m *MockAttachmentUsecase) GetAttachmentFile(ctx context.Context, attachmentID string) (io.ReadCloser, *dto.AttachmentResponse, error) {
	if m.GetAttachmentFileFunc != nil {
		return m.GetAttachmentFileFunc(ctx, attachmentID)
	}
	return nil, nil, nil
}

func (m *MockAttachmentUsecase) ListAttachmentsByRecordID(ctx context.Context, recordID string) ([]*dto.AttachmentResponse, error) {
	if m.ListAttachmentsByRecordIDFunc != nil {
		return m.ListAttachmentsByRecordIDFunc(ctx, recordID)
	}
	return nil, nil
}

func (m *MockAttachmentUsecase) ListAttachmentsByStepID(ctx context.Context, stepID string) ([]*dto.AttachmentResponse, error) {
	if m.ListAttachmentsByStepIDFunc != nil {
		return m.ListAttachmentsByStepIDFunc(ctx, stepID)
	}
	return nil, nil
}

func (m *MockAttachmentUsecase) DeleteAttachment(ctx context.Context, attachmentID string) error {
	if m.DeleteAttachmentFunc != nil {
		return m.DeleteAttachmentFunc(ctx, attachmentID)
	}
	return nil
}

func (m *MockAttachmentUsecase) GetAttachmentURL(ctx context.Context, attachmentID string, expirationMinutes int) (string, error) {
	if m.GetAttachmentURLFunc != nil {
		return m.GetAttachmentURLFunc(ctx, attachmentID, expirationMinutes)
	}
	return "", nil
}

func TestAttachmentHandler_UploadAttachment(t *testing.T) {
	gin.SetMode(gin.TestMode)

	recordID := value_object.GenerateExecutionRecordID()
	stepID := value_object.GenerateExecutionStepID()
	attachmentID := value_object.GenerateAttachmentID()

	mockUsecase := &MockAttachmentUsecase{
		UploadAttachmentFunc: func(ctx context.Context, req *dto.UploadAttachmentRequest) (*dto.AttachmentResponse, error) {
			return &dto.AttachmentResponse{
				ID:                attachmentID.String(),
				ExecutionRecordID: recordID.String(),
				ExecutionStepID:   stepID.String(),
				FileName:          req.FileName,
				FileSize:          req.FileSize,
				MimeType:          req.MimeType,
				StorageType:       "local",
				UploadedBy:        req.UploadedBy,
			}, nil
		},
	}

	handler := NewAttachmentHandler(mockUsecase)

	// Create multipart form
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	writer.WriteField("execution_step_id", stepID.String())
	part, _ := writer.CreateFormFile("file", "test.png")
	part.Write([]byte("test file content"))
	writer.Close()

	req, _ := http.NewRequest("POST", "/api/v1/execution-records/"+recordID.String()+"/attachments", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	c.Params = gin.Params{{Key: "id", Value: recordID.String()}}
	c.Set("user_id", "user-123")

	handler.UploadAttachment(c)

	assert.Equal(t, http.StatusCreated, w.Code)
}

func TestAttachmentHandler_GetAttachment(t *testing.T) {
	gin.SetMode(gin.TestMode)

	attachmentID := value_object.GenerateAttachmentID()
	recordID := value_object.GenerateExecutionRecordID()
	stepID := value_object.GenerateExecutionStepID()

	mockUsecase := &MockAttachmentUsecase{
		GetAttachmentFunc: func(ctx context.Context, id string) (*dto.AttachmentResponse, error) {
			return &dto.AttachmentResponse{
				ID:                attachmentID.String(),
				ExecutionRecordID: recordID.String(),
				ExecutionStepID:   stepID.String(),
				FileName:          "test.png",
				FileSize:          1024,
				MimeType:          "image/png",
				StorageType:       "local",
				UploadedBy:        "user-123",
			}, nil
		},
	}

	handler := NewAttachmentHandler(mockUsecase)

	req, _ := http.NewRequest("GET", "/api/v1/attachments/"+attachmentID.String(), nil)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	c.Params = gin.Params{{Key: "id", Value: attachmentID.String()}}

	handler.GetAttachment(c)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, attachmentID.String(), response["id"])
	assert.Equal(t, "test.png", response["file_name"])
}

func TestAttachmentHandler_DeleteAttachment(t *testing.T) {
	gin.SetMode(gin.TestMode)

	attachmentID := value_object.GenerateAttachmentID()

	mockUsecase := &MockAttachmentUsecase{
		DeleteAttachmentFunc: func(ctx context.Context, id string) error {
			return nil
		},
	}

	handler := NewAttachmentHandler(mockUsecase)

	router := gin.New()
	router.DELETE("/attachments/:id", handler.DeleteAttachment)

	req, _ := http.NewRequest("DELETE", "/attachments/"+attachmentID.String(), nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNoContent, w.Code)
}

func TestAttachmentHandler_ListAttachments(t *testing.T) {
	gin.SetMode(gin.TestMode)

	recordID := value_object.GenerateExecutionRecordID()
	stepID := value_object.GenerateExecutionStepID()

	mockUsecase := &MockAttachmentUsecase{
		ListAttachmentsByRecordIDFunc: func(ctx context.Context, id string) ([]*dto.AttachmentResponse, error) {
			return []*dto.AttachmentResponse{
				{
					ID:                value_object.GenerateAttachmentID().String(),
					ExecutionRecordID: recordID.String(),
					ExecutionStepID:   stepID.String(),
					FileName:          "test1.png",
					FileSize:          1024,
					MimeType:          "image/png",
					StorageType:       "local",
					UploadedBy:        "user-123",
				},
				{
					ID:                value_object.GenerateAttachmentID().String(),
					ExecutionRecordID: recordID.String(),
					ExecutionStepID:   stepID.String(),
					FileName:          "test2.png",
					FileSize:          2048,
					MimeType:          "image/png",
					StorageType:       "local",
					UploadedBy:        "user-123",
				},
			}, nil
		},
	}

	handler := NewAttachmentHandler(mockUsecase)

	req, _ := http.NewRequest("GET", "/api/v1/execution-records/"+recordID.String()+"/attachments", nil)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	c.Params = gin.Params{{Key: "id", Value: recordID.String()}}

	handler.ListAttachments(c)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	attachments := response["attachments"].([]interface{})
	assert.Len(t, attachments, 2)
}

func TestAttachmentHandler_GetAttachmentURL(t *testing.T) {
	gin.SetMode(gin.TestMode)

	attachmentID := value_object.GenerateAttachmentID()

	mockUsecase := &MockAttachmentUsecase{
		GetAttachmentURLFunc: func(ctx context.Context, id string, expirationMinutes int) (string, error) {
			return "", nil // Local storage returns empty
		},
	}

	handler := NewAttachmentHandler(mockUsecase)

	req, _ := http.NewRequest("GET", "/api/v1/attachments/"+attachmentID.String()+"/url", nil)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	c.Params = gin.Params{{Key: "id", Value: attachmentID.String()}}

	handler.GetAttachmentURL(c)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Contains(t, response["url"], "/attachments/") // Should return download URL
}
