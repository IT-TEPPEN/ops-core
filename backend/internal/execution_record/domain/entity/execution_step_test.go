package entity

import (
	"testing"

	"opscore/backend/internal/execution_record/domain/value_object"
)

func TestNewExecutionStep(t *testing.T) {
	validID := value_object.GenerateExecutionStepID()
	validRecordID := value_object.GenerateExecutionRecordID()

	tests := []struct {
		name        string
		id          value_object.ExecutionStepID
		recordID    value_object.ExecutionRecordID
		stepNumber  int
		description string
		wantErr     bool
	}{
		{
			name:        "valid execution step",
			id:          validID,
			recordID:    validRecordID,
			stepNumber:  1,
			description: "First step",
			wantErr:     false,
		},
		{
			name:        "empty step ID",
			id:          value_object.ExecutionStepID(""),
			recordID:    validRecordID,
			stepNumber:  1,
			description: "First step",
			wantErr:     true,
		},
		{
			name:        "empty record ID",
			id:          validID,
			recordID:    value_object.ExecutionRecordID(""),
			stepNumber:  1,
			description: "First step",
			wantErr:     true,
		},
		{
			name:        "zero step number",
			id:          validID,
			recordID:    validRecordID,
			stepNumber:  0,
			description: "First step",
			wantErr:     true,
		},
		{
			name:        "negative step number",
			id:          validID,
			recordID:    validRecordID,
			stepNumber:  -1,
			description: "First step",
			wantErr:     true,
		},
		{
			name:        "empty description",
			id:          validID,
			recordID:    validRecordID,
			stepNumber:  1,
			description: "",
			wantErr:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewExecutionStep(tt.id, tt.recordID, tt.stepNumber, tt.description)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewExecutionStep() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if got == nil {
					t.Error("NewExecutionStep() returned nil")
					return
				}
				if !got.ID().Equals(tt.id) {
					t.Errorf("ID() = %v, want %v", got.ID(), tt.id)
				}
				if !got.ExecutionRecordID().Equals(tt.recordID) {
					t.Errorf("ExecutionRecordID() = %v, want %v", got.ExecutionRecordID(), tt.recordID)
				}
				if got.StepNumber() != tt.stepNumber {
					t.Errorf("StepNumber() = %v, want %v", got.StepNumber(), tt.stepNumber)
				}
				if got.Description() != tt.description {
					t.Errorf("Description() = %v, want %v", got.Description(), tt.description)
				}
				if got.Notes() != "" {
					t.Error("Notes() should be empty for new step")
				}
			}
		})
	}
}

func TestExecutionStep_UpdateNotes(t *testing.T) {
	step := createTestExecutionStep(t)

	notes := "This step was executed successfully"
	step.UpdateNotes(notes)

	if step.Notes() != notes {
		t.Errorf("Notes() = %v, want %v", step.Notes(), notes)
	}
}

func TestReconstructExecutionStep(t *testing.T) {
	id := value_object.GenerateExecutionStepID()
	recordID := value_object.GenerateExecutionRecordID()
	stepNumber := 2
	description := "Test step"
	notes := "Some notes"
	executedAt := fixedTime()

	step := ReconstructExecutionStep(id, recordID, stepNumber, description, notes, executedAt)

	if !step.ID().Equals(id) {
		t.Errorf("ID() = %v, want %v", step.ID(), id)
	}
	if !step.ExecutionRecordID().Equals(recordID) {
		t.Errorf("ExecutionRecordID() = %v, want %v", step.ExecutionRecordID(), recordID)
	}
	if step.StepNumber() != stepNumber {
		t.Errorf("StepNumber() = %v, want %v", step.StepNumber(), stepNumber)
	}
	if step.Description() != description {
		t.Errorf("Description() = %v, want %v", step.Description(), description)
	}
	if step.Notes() != notes {
		t.Errorf("Notes() = %v, want %v", step.Notes(), notes)
	}
	if !step.ExecutedAt().Equal(executedAt) {
		t.Errorf("ExecutedAt() = %v, want %v", step.ExecutedAt(), executedAt)
	}
}

// Helper function to create a test execution step
func createTestExecutionStep(t *testing.T) ExecutionStep {
	id := value_object.GenerateExecutionStepID()
	recordID := value_object.GenerateExecutionRecordID()

	step, err := NewExecutionStep(id, recordID, 1, "Test step description")
	if err != nil {
		t.Fatalf("Failed to create test execution step: %v", err)
	}

	return step
}
