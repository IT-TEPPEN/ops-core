package usecase

import (
	"context"
	"io"

	docvo "opscore/backend/internal/document/domain/value_object"
	"opscore/backend/internal/execution_record/domain/entity"
	"opscore/backend/internal/execution_record/domain/repository"
	"opscore/backend/internal/execution_record/domain/value_object"
)

// MockExecutionRecordRepository is a mock implementation of ExecutionRecordRepository for testing.
type MockExecutionRecordRepository struct {
	SaveFunc            func(ctx context.Context, record entity.ExecutionRecord) error
	FindByIDFunc        func(ctx context.Context, id value_object.ExecutionRecordID) (entity.ExecutionRecord, error)
	FindByExecutorIDFunc func(ctx context.Context, executorID string) ([]entity.ExecutionRecord, error)
	FindByDocumentIDFunc func(ctx context.Context, documentID docvo.DocumentID) ([]entity.ExecutionRecord, error)
	SearchFunc          func(ctx context.Context, criteria repository.SearchCriteria) ([]entity.ExecutionRecord, error)
	UpdateFunc          func(ctx context.Context, record entity.ExecutionRecord) error
	DeleteFunc          func(ctx context.Context, id value_object.ExecutionRecordID) error
}

func (m *MockExecutionRecordRepository) Save(ctx context.Context, record entity.ExecutionRecord) error {
	if m.SaveFunc != nil {
		return m.SaveFunc(ctx, record)
	}
	return nil
}

func (m *MockExecutionRecordRepository) FindByID(ctx context.Context, id value_object.ExecutionRecordID) (entity.ExecutionRecord, error) {
	if m.FindByIDFunc != nil {
		return m.FindByIDFunc(ctx, id)
	}
	return nil, nil
}

func (m *MockExecutionRecordRepository) FindByExecutorID(ctx context.Context, executorID string) ([]entity.ExecutionRecord, error) {
	if m.FindByExecutorIDFunc != nil {
		return m.FindByExecutorIDFunc(ctx, executorID)
	}
	return nil, nil
}

func (m *MockExecutionRecordRepository) FindByDocumentID(ctx context.Context, documentID docvo.DocumentID) ([]entity.ExecutionRecord, error) {
	if m.FindByDocumentIDFunc != nil {
		return m.FindByDocumentIDFunc(ctx, documentID)
	}
	return nil, nil
}

func (m *MockExecutionRecordRepository) Search(ctx context.Context, criteria repository.SearchCriteria) ([]entity.ExecutionRecord, error) {
	if m.SearchFunc != nil {
		return m.SearchFunc(ctx, criteria)
	}
	return nil, nil
}

func (m *MockExecutionRecordRepository) Update(ctx context.Context, record entity.ExecutionRecord) error {
	if m.UpdateFunc != nil {
		return m.UpdateFunc(ctx, record)
	}
	return nil
}

func (m *MockExecutionRecordRepository) Delete(ctx context.Context, id value_object.ExecutionRecordID) error {
	if m.DeleteFunc != nil {
		return m.DeleteFunc(ctx, id)
	}
	return nil
}

// MockAttachmentRepository is a mock implementation of AttachmentRepository for testing.
type MockAttachmentRepository struct {
	SaveFunc                    func(ctx context.Context, attachment entity.Attachment, file io.Reader) error
	FindByIDFunc                func(ctx context.Context, id value_object.AttachmentID) (entity.Attachment, error)
	FindByExecutionRecordIDFunc func(ctx context.Context, recordID value_object.ExecutionRecordID) ([]entity.Attachment, error)
	FindByExecutionStepIDFunc   func(ctx context.Context, stepID value_object.ExecutionStepID) ([]entity.Attachment, error)
	GetFileFunc                 func(ctx context.Context, id value_object.AttachmentID) (io.ReadCloser, error)
	DeleteFunc                  func(ctx context.Context, id value_object.AttachmentID) error
}

func (m *MockAttachmentRepository) Save(ctx context.Context, attachment entity.Attachment, file io.Reader) error {
	if m.SaveFunc != nil {
		return m.SaveFunc(ctx, attachment, file)
	}
	return nil
}

func (m *MockAttachmentRepository) FindByID(ctx context.Context, id value_object.AttachmentID) (entity.Attachment, error) {
	if m.FindByIDFunc != nil {
		return m.FindByIDFunc(ctx, id)
	}
	return nil, nil
}

func (m *MockAttachmentRepository) FindByExecutionRecordID(ctx context.Context, recordID value_object.ExecutionRecordID) ([]entity.Attachment, error) {
	if m.FindByExecutionRecordIDFunc != nil {
		return m.FindByExecutionRecordIDFunc(ctx, recordID)
	}
	return nil, nil
}

func (m *MockAttachmentRepository) FindByExecutionStepID(ctx context.Context, stepID value_object.ExecutionStepID) ([]entity.Attachment, error) {
	if m.FindByExecutionStepIDFunc != nil {
		return m.FindByExecutionStepIDFunc(ctx, stepID)
	}
	return nil, nil
}

func (m *MockAttachmentRepository) GetFile(ctx context.Context, id value_object.AttachmentID) (io.ReadCloser, error) {
	if m.GetFileFunc != nil {
		return m.GetFileFunc(ctx, id)
	}
	return nil, nil
}

func (m *MockAttachmentRepository) Delete(ctx context.Context, id value_object.AttachmentID) error {
	if m.DeleteFunc != nil {
		return m.DeleteFunc(ctx, id)
	}
	return nil
}

// MockStorageManager is a mock implementation of StorageManager for testing.
type MockStorageManager struct {
	StoreFunc    func(ctx context.Context, path string, file io.Reader) (string, error)
	RetrieveFunc func(ctx context.Context, path string) (io.ReadCloser, error)
	DeleteFunc   func(ctx context.Context, path string) error
	TypeFunc     func() value_object.StorageType
}

func (m *MockStorageManager) Store(ctx context.Context, path string, file io.Reader) (string, error) {
	if m.StoreFunc != nil {
		return m.StoreFunc(ctx, path, file)
	}
	return path, nil
}

func (m *MockStorageManager) Retrieve(ctx context.Context, path string) (io.ReadCloser, error) {
	if m.RetrieveFunc != nil {
		return m.RetrieveFunc(ctx, path)
	}
	return nil, nil
}

func (m *MockStorageManager) Delete(ctx context.Context, path string) error {
	if m.DeleteFunc != nil {
		return m.DeleteFunc(ctx, path)
	}
	return nil
}

func (m *MockStorageManager) Type() value_object.StorageType {
	if m.TypeFunc != nil {
		return m.TypeFunc()
	}
	return value_object.StorageTypeLocal
}
