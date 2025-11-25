package usecase

import (
	"context"
	"errors"
	"testing"

	docvo "opscore/backend/internal/document/domain/value_object"
	apperror "opscore/backend/internal/execution_record/application/error"
	"opscore/backend/internal/execution_record/application/dto"
	"opscore/backend/internal/execution_record/domain/entity"
	"opscore/backend/internal/execution_record/domain/value_object"
)

func TestExecutionRecordUsecase_CreateExecutionRecord(t *testing.T) {
	mockRepo := &MockExecutionRecordRepository{}
	uc := NewExecutionRecordUsecase(mockRepo)

	ctx := context.Background()
	docID := docvo.GenerateDocumentID()
	versionID := docvo.GenerateVersionID()

	req := &dto.CreateExecutionRecordRequest{
		DocumentID:        docID.String(),
		DocumentVersionID: versionID.String(),
		ExecutorID:        "user-123",
		Title:             "Test Execution",
		VariableValues:    []dto.VariableValueDTO{},
	}

	resp, err := uc.CreateExecutionRecord(ctx, req)
	if err != nil {
		t.Fatalf("CreateExecutionRecord() error = %v", err)
	}

	if resp == nil {
		t.Fatal("CreateExecutionRecord() returned nil response")
	}

	if resp.Title != req.Title {
		t.Errorf("Title = %v, want %v", resp.Title, req.Title)
	}

	if resp.ExecutorID != req.ExecutorID {
		t.Errorf("ExecutorID = %v, want %v", resp.ExecutorID, req.ExecutorID)
	}

	if resp.Status != "in_progress" {
		t.Errorf("Status = %v, want in_progress", resp.Status)
	}
}

func TestExecutionRecordUsecase_CreateExecutionRecord_InvalidDocumentID(t *testing.T) {
	mockRepo := &MockExecutionRecordRepository{}
	uc := NewExecutionRecordUsecase(mockRepo)

	ctx := context.Background()

	req := &dto.CreateExecutionRecordRequest{
		DocumentID:        "invalid-id",
		DocumentVersionID: docvo.GenerateVersionID().String(),
		ExecutorID:        "user-123",
		Title:             "Test Execution",
	}

	_, err := uc.CreateExecutionRecord(ctx, req)
	if err == nil {
		t.Fatal("CreateExecutionRecord() should return error for invalid document ID")
	}

	if !errors.Is(err, apperror.ErrValidationFailed) {
		t.Errorf("Error should be validation error, got %v", err)
	}
}

func TestExecutionRecordUsecase_GetExecutionRecord(t *testing.T) {
	docID := docvo.GenerateDocumentID()
	versionID := docvo.GenerateVersionID()
	recordID := value_object.GenerateExecutionRecordID()

	record, _ := entity.NewExecutionRecord(
		recordID,
		docID,
		versionID,
		"user-123",
		"Test Execution",
		[]value_object.VariableValue{},
	)

	mockRepo := &MockExecutionRecordRepository{
		FindByIDFunc: func(ctx context.Context, id value_object.ExecutionRecordID) (entity.ExecutionRecord, error) {
			if id.Equals(recordID) {
				return record, nil
			}
			return nil, nil
		},
	}

	uc := NewExecutionRecordUsecase(mockRepo)
	ctx := context.Background()

	resp, err := uc.GetExecutionRecord(ctx, recordID.String())
	if err != nil {
		t.Fatalf("GetExecutionRecord() error = %v", err)
	}

	if resp == nil {
		t.Fatal("GetExecutionRecord() returned nil response")
	}

	if resp.ID != recordID.String() {
		t.Errorf("ID = %v, want %v", resp.ID, recordID.String())
	}
}

func TestExecutionRecordUsecase_GetExecutionRecord_NotFound(t *testing.T) {
	mockRepo := &MockExecutionRecordRepository{
		FindByIDFunc: func(ctx context.Context, id value_object.ExecutionRecordID) (entity.ExecutionRecord, error) {
			return nil, nil
		},
	}

	uc := NewExecutionRecordUsecase(mockRepo)
	ctx := context.Background()

	recordID := value_object.GenerateExecutionRecordID()
	_, err := uc.GetExecutionRecord(ctx, recordID.String())
	if err == nil {
		t.Fatal("GetExecutionRecord() should return error for not found")
	}

	if !errors.Is(err, apperror.ErrNotFound) {
		t.Errorf("Error should be not found error, got %v", err)
	}
}

func TestExecutionRecordUsecase_AddStep(t *testing.T) {
	docID := docvo.GenerateDocumentID()
	versionID := docvo.GenerateVersionID()
	recordID := value_object.GenerateExecutionRecordID()

	record, _ := entity.NewExecutionRecord(
		recordID,
		docID,
		versionID,
		"user-123",
		"Test Execution",
		[]value_object.VariableValue{},
	)

	mockRepo := &MockExecutionRecordRepository{
		FindByIDFunc: func(ctx context.Context, id value_object.ExecutionRecordID) (entity.ExecutionRecord, error) {
			return record, nil
		},
	}

	uc := NewExecutionRecordUsecase(mockRepo)
	ctx := context.Background()

	req := &dto.AddStepRequest{
		ExecutionRecordID: recordID.String(),
		StepNumber:        1,
		Description:       "First step",
	}

	resp, err := uc.AddStep(ctx, req)
	if err != nil {
		t.Fatalf("AddStep() error = %v", err)
	}

	if len(resp.Steps) != 1 {
		t.Errorf("Steps length = %d, want 1", len(resp.Steps))
	}

	if resp.Steps[0].Description != "First step" {
		t.Errorf("Step description = %v, want 'First step'", resp.Steps[0].Description)
	}
}

func TestExecutionRecordUsecase_Complete(t *testing.T) {
	docID := docvo.GenerateDocumentID()
	versionID := docvo.GenerateVersionID()
	recordID := value_object.GenerateExecutionRecordID()

	record, _ := entity.NewExecutionRecord(
		recordID,
		docID,
		versionID,
		"user-123",
		"Test Execution",
		[]value_object.VariableValue{},
	)

	mockRepo := &MockExecutionRecordRepository{
		FindByIDFunc: func(ctx context.Context, id value_object.ExecutionRecordID) (entity.ExecutionRecord, error) {
			return record, nil
		},
	}

	uc := NewExecutionRecordUsecase(mockRepo)
	ctx := context.Background()

	req := &dto.CompleteExecutionRequest{
		ExecutionRecordID: recordID.String(),
	}

	resp, err := uc.Complete(ctx, req)
	if err != nil {
		t.Fatalf("Complete() error = %v", err)
	}

	if resp.Status != "completed" {
		t.Errorf("Status = %v, want 'completed'", resp.Status)
	}

	if resp.CompletedAt == nil {
		t.Error("CompletedAt should not be nil after completion")
	}
}

func TestExecutionRecordUsecase_MarkAsFailed(t *testing.T) {
	docID := docvo.GenerateDocumentID()
	versionID := docvo.GenerateVersionID()
	recordID := value_object.GenerateExecutionRecordID()

	record, _ := entity.NewExecutionRecord(
		recordID,
		docID,
		versionID,
		"user-123",
		"Test Execution",
		[]value_object.VariableValue{},
	)

	mockRepo := &MockExecutionRecordRepository{
		FindByIDFunc: func(ctx context.Context, id value_object.ExecutionRecordID) (entity.ExecutionRecord, error) {
			return record, nil
		},
	}

	uc := NewExecutionRecordUsecase(mockRepo)
	ctx := context.Background()

	req := &dto.MarkAsFailedRequest{
		ExecutionRecordID: recordID.String(),
	}

	resp, err := uc.MarkAsFailed(ctx, req)
	if err != nil {
		t.Fatalf("MarkAsFailed() error = %v", err)
	}

	if resp.Status != "failed" {
		t.Errorf("Status = %v, want 'failed'", resp.Status)
	}
}

func TestExecutionRecordUsecase_UpdateAccessScope(t *testing.T) {
	docID := docvo.GenerateDocumentID()
	versionID := docvo.GenerateVersionID()
	recordID := value_object.GenerateExecutionRecordID()

	record, _ := entity.NewExecutionRecord(
		recordID,
		docID,
		versionID,
		"user-123",
		"Test Execution",
		[]value_object.VariableValue{},
	)

	mockRepo := &MockExecutionRecordRepository{
		FindByIDFunc: func(ctx context.Context, id value_object.ExecutionRecordID) (entity.ExecutionRecord, error) {
			return record, nil
		},
	}

	uc := NewExecutionRecordUsecase(mockRepo)
	ctx := context.Background()

	req := &dto.UpdateAccessScopeRequest{
		ExecutionRecordID: recordID.String(),
		AccessScope:       "public",
	}

	resp, err := uc.UpdateAccessScope(ctx, req)
	if err != nil {
		t.Fatalf("UpdateAccessScope() error = %v", err)
	}

	if resp.AccessScope != "public" {
		t.Errorf("AccessScope = %v, want 'public'", resp.AccessScope)
	}
}

func TestExecutionRecordUsecase_DeleteExecutionRecord(t *testing.T) {
	docID := docvo.GenerateDocumentID()
	versionID := docvo.GenerateVersionID()
	recordID := value_object.GenerateExecutionRecordID()

	record, _ := entity.NewExecutionRecord(
		recordID,
		docID,
		versionID,
		"user-123",
		"Test Execution",
		[]value_object.VariableValue{},
	)

	deleted := false
	mockRepo := &MockExecutionRecordRepository{
		FindByIDFunc: func(ctx context.Context, id value_object.ExecutionRecordID) (entity.ExecutionRecord, error) {
			return record, nil
		},
		DeleteFunc: func(ctx context.Context, id value_object.ExecutionRecordID) error {
			deleted = true
			return nil
		},
	}

	uc := NewExecutionRecordUsecase(mockRepo)
	ctx := context.Background()

	err := uc.DeleteExecutionRecord(ctx, recordID.String())
	if err != nil {
		t.Fatalf("DeleteExecutionRecord() error = %v", err)
	}

	if !deleted {
		t.Error("Delete was not called")
	}
}
