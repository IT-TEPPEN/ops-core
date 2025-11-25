package entity

import (
	"testing"
	"time"

	docvo "opscore/backend/internal/document/domain/value_object"
	"opscore/backend/internal/execution_record/domain/value_object"
)

func TestNewExecutionRecord(t *testing.T) {
	validID := value_object.GenerateExecutionRecordID()
	validDocID := docvo.GenerateDocumentID()
	validVersionID := docvo.GenerateVersionID()

	tests := []struct {
		name           string
		id             value_object.ExecutionRecordID
		documentID     docvo.DocumentID
		versionID      docvo.VersionID
		executorID     string
		title          string
		variableValues []value_object.VariableValue
		wantErr        bool
	}{
		{
			name:           "valid execution record",
			id:             validID,
			documentID:     validDocID,
			versionID:      validVersionID,
			executorID:     "user-123",
			title:          "Test Execution",
			variableValues: []value_object.VariableValue{},
			wantErr:        false,
		},
		{
			name:           "empty record ID",
			id:             value_object.ExecutionRecordID(""),
			documentID:     validDocID,
			versionID:      validVersionID,
			executorID:     "user-123",
			title:          "Test Execution",
			variableValues: []value_object.VariableValue{},
			wantErr:        true,
		},
		{
			name:           "empty document ID",
			id:             validID,
			documentID:     docvo.DocumentID(""),
			versionID:      validVersionID,
			executorID:     "user-123",
			title:          "Test Execution",
			variableValues: []value_object.VariableValue{},
			wantErr:        true,
		},
		{
			name:           "empty version ID",
			id:             validID,
			documentID:     validDocID,
			versionID:      docvo.VersionID(""),
			executorID:     "user-123",
			title:          "Test Execution",
			variableValues: []value_object.VariableValue{},
			wantErr:        true,
		},
		{
			name:           "empty executor ID",
			id:             validID,
			documentID:     validDocID,
			versionID:      validVersionID,
			executorID:     "",
			title:          "Test Execution",
			variableValues: []value_object.VariableValue{},
			wantErr:        true,
		},
		{
			name:           "empty title",
			id:             validID,
			documentID:     validDocID,
			versionID:      validVersionID,
			executorID:     "user-123",
			title:          "",
			variableValues: []value_object.VariableValue{},
			wantErr:        true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewExecutionRecord(
				tt.id,
				tt.documentID,
				tt.versionID,
				tt.executorID,
				tt.title,
				tt.variableValues,
			)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewExecutionRecord() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if got == nil {
					t.Error("NewExecutionRecord() returned nil")
					return
				}
				if got.Status() != value_object.ExecutionStatusInProgress {
					t.Errorf("Status() = %v, want %v", got.Status(), value_object.ExecutionStatusInProgress)
				}
				if got.AccessScope() != value_object.AccessScopePrivate {
					t.Errorf("AccessScope() = %v, want %v", got.AccessScope(), value_object.AccessScopePrivate)
				}
				if len(got.Steps()) != 0 {
					t.Errorf("Steps() length = %d, want 0", len(got.Steps()))
				}
				if got.CompletedAt() != nil {
					t.Error("CompletedAt() should be nil for new record")
				}
			}
		})
	}
}

func TestExecutionRecord_AddStep(t *testing.T) {
	record := createTestExecutionRecord(t)

	// Add first step
	err := record.AddStep(1, "First step")
	if err != nil {
		t.Errorf("AddStep() error = %v", err)
	}

	if len(record.Steps()) != 1 {
		t.Errorf("Steps() length = %d, want 1", len(record.Steps()))
	}

	// Add second step
	err = record.AddStep(2, "Second step")
	if err != nil {
		t.Errorf("AddStep() second step error = %v", err)
	}

	if len(record.Steps()) != 2 {
		t.Errorf("Steps() length = %d, want 2", len(record.Steps()))
	}

	// Try to add duplicate step number
	err = record.AddStep(1, "Duplicate step")
	if err == nil {
		t.Error("AddStep() should return error for duplicate step number")
	}
}

func TestExecutionRecord_AddStep_AfterCompletion(t *testing.T) {
	record := createTestExecutionRecord(t)

	// Complete the record
	err := record.Complete()
	if err != nil {
		t.Fatalf("Complete() error = %v", err)
	}

	// Try to add step after completion
	err = record.AddStep(1, "Step after completion")
	if err == nil {
		t.Error("AddStep() should return error for completed execution")
	}
}

func TestExecutionRecord_UpdateStepNotes(t *testing.T) {
	record := createTestExecutionRecord(t)

	// Add a step first
	err := record.AddStep(1, "Test step")
	if err != nil {
		t.Fatalf("AddStep() error = %v", err)
	}

	// Update notes
	notes := "Step completed successfully"
	err = record.UpdateStepNotes(1, notes)
	if err != nil {
		t.Errorf("UpdateStepNotes() error = %v", err)
	}

	if record.Steps()[0].Notes() != notes {
		t.Errorf("Step notes = %v, want %v", record.Steps()[0].Notes(), notes)
	}

	// Try to update notes for non-existent step
	err = record.UpdateStepNotes(99, "Notes")
	if err == nil {
		t.Error("UpdateStepNotes() should return error for non-existent step")
	}
}

func TestExecutionRecord_UpdateNotes(t *testing.T) {
	record := createTestExecutionRecord(t)

	notes := "Overall execution notes"
	record.UpdateNotes(notes)

	if record.Notes() != notes {
		t.Errorf("Notes() = %v, want %v", record.Notes(), notes)
	}
}

func TestExecutionRecord_UpdateTitle(t *testing.T) {
	record := createTestExecutionRecord(t)

	newTitle := "Updated Title"
	err := record.UpdateTitle(newTitle)
	if err != nil {
		t.Errorf("UpdateTitle() error = %v", err)
	}

	if record.Title() != newTitle {
		t.Errorf("Title() = %v, want %v", record.Title(), newTitle)
	}

	// Try to set empty title
	err = record.UpdateTitle("")
	if err == nil {
		t.Error("UpdateTitle() should return error for empty title")
	}
}

func TestExecutionRecord_Complete(t *testing.T) {
	record := createTestExecutionRecord(t)

	err := record.Complete()
	if err != nil {
		t.Errorf("Complete() error = %v", err)
	}

	if record.Status() != value_object.ExecutionStatusCompleted {
		t.Errorf("Status() = %v, want %v", record.Status(), value_object.ExecutionStatusCompleted)
	}

	if record.CompletedAt() == nil {
		t.Error("CompletedAt() should not be nil after completion")
	}

	// Try to complete again
	err = record.Complete()
	if err == nil {
		t.Error("Complete() should return error for already completed execution")
	}
}

func TestExecutionRecord_MarkAsFailed(t *testing.T) {
	record := createTestExecutionRecord(t)

	err := record.MarkAsFailed()
	if err != nil {
		t.Errorf("MarkAsFailed() error = %v", err)
	}

	if record.Status() != value_object.ExecutionStatusFailed {
		t.Errorf("Status() = %v, want %v", record.Status(), value_object.ExecutionStatusFailed)
	}

	if record.CompletedAt() == nil {
		t.Error("CompletedAt() should not be nil after marking as failed")
	}

	// Try to mark as failed again
	err = record.MarkAsFailed()
	if err == nil {
		t.Error("MarkAsFailed() should return error for already failed execution")
	}
}

func TestExecutionRecord_UpdateAccessScope(t *testing.T) {
	record := createTestExecutionRecord(t)

	// Default should be private
	if record.AccessScope() != value_object.AccessScopePrivate {
		t.Errorf("AccessScope() = %v, want %v", record.AccessScope(), value_object.AccessScopePrivate)
	}

	// Update to public
	record.UpdateAccessScope(value_object.AccessScopePublic)

	if record.AccessScope() != value_object.AccessScopePublic {
		t.Errorf("AccessScope() = %v, want %v", record.AccessScope(), value_object.AccessScopePublic)
	}
}

func TestReconstructExecutionRecord(t *testing.T) {
	id := value_object.GenerateExecutionRecordID()
	docID := docvo.GenerateDocumentID()
	versionID := docvo.GenerateVersionID()
	executorID := "user-456"
	title := "Reconstructed Execution"
	vv, _ := value_object.NewVariableValue("server", "localhost")
	variableValues := []value_object.VariableValue{vv}
	notes := "Test notes"
	status := value_object.ExecutionStatusCompleted
	accessScope := value_object.AccessScopePublic
	steps := []ExecutionStep{}
	startedAt := fixedTime()
	completedAt := fixedTime()
	createdAt := fixedTime()
	updatedAt := fixedTime()

	record := ReconstructExecutionRecord(
		id,
		docID,
		versionID,
		executorID,
		title,
		variableValues,
		notes,
		status,
		accessScope,
		steps,
		startedAt,
		&completedAt,
		createdAt,
		updatedAt,
	)

	if !record.ID().Equals(id) {
		t.Errorf("ID() = %v, want %v", record.ID(), id)
	}
	if !record.DocumentID().Equals(docID) {
		t.Errorf("DocumentID() = %v, want %v", record.DocumentID(), docID)
	}
	if !record.DocumentVersionID().Equals(versionID) {
		t.Errorf("DocumentVersionID() = %v, want %v", record.DocumentVersionID(), versionID)
	}
	if record.ExecutorID() != executorID {
		t.Errorf("ExecutorID() = %v, want %v", record.ExecutorID(), executorID)
	}
	if record.Title() != title {
		t.Errorf("Title() = %v, want %v", record.Title(), title)
	}
	if len(record.VariableValues()) != 1 {
		t.Errorf("VariableValues() length = %d, want 1", len(record.VariableValues()))
	}
	if record.Notes() != notes {
		t.Errorf("Notes() = %v, want %v", record.Notes(), notes)
	}
	if record.Status() != status {
		t.Errorf("Status() = %v, want %v", record.Status(), status)
	}
	if record.AccessScope() != accessScope {
		t.Errorf("AccessScope() = %v, want %v", record.AccessScope(), accessScope)
	}
}

// Helper function to create a test execution record
func createTestExecutionRecord(t *testing.T) ExecutionRecord {
	id := value_object.GenerateExecutionRecordID()
	docID := docvo.GenerateDocumentID()
	versionID := docvo.GenerateVersionID()

	record, err := NewExecutionRecord(
		id,
		docID,
		versionID,
		"user-123",
		"Test Execution",
		[]value_object.VariableValue{},
	)

	if err != nil {
		t.Fatalf("Failed to create test execution record: %v", err)
	}

	return record
}

// Helper function to get a fixed time for testing
func fixedTime() time.Time {
	return time.Date(2024, 1, 15, 10, 30, 0, 0, time.UTC)
}
