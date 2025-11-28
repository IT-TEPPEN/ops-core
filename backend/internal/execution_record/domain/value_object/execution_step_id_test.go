package value_object

import (
	"testing"

	"github.com/google/uuid"
)

func TestNewExecutionStepID(t *testing.T) {
	tests := []struct {
		name    string
		id      string
		wantErr bool
	}{
		{
			name:    "valid UUID",
			id:      uuid.New().String(),
			wantErr: false,
		},
		{
			name:    "empty string",
			id:      "",
			wantErr: true,
		},
		{
			name:    "invalid UUID",
			id:      "not-a-uuid",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewExecutionStepID(tt.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewExecutionStepID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got.String() != tt.id {
				t.Errorf("NewExecutionStepID() = %v, want %v", got, tt.id)
			}
		})
	}
}

func TestGenerateExecutionStepID(t *testing.T) {
	id := GenerateExecutionStepID()
	if id.IsEmpty() {
		t.Error("GenerateExecutionStepID() returned empty ID")
	}
	// Verify it's a valid UUID
	if _, err := uuid.Parse(id.String()); err != nil {
		t.Errorf("GenerateExecutionStepID() returned invalid UUID: %v", err)
	}
}

func TestExecutionStepID_IsEmpty(t *testing.T) {
	tests := []struct {
		name string
		id   ExecutionStepID
		want bool
	}{
		{
			name: "empty ID",
			id:   ExecutionStepID(""),
			want: true,
		},
		{
			name: "non-empty ID",
			id:   GenerateExecutionStepID(),
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.id.IsEmpty(); got != tt.want {
				t.Errorf("ExecutionStepID.IsEmpty() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestExecutionStepID_Equals(t *testing.T) {
	id1 := GenerateExecutionStepID()
	id2 := GenerateExecutionStepID()
	id1Copy := ExecutionStepID(id1.String())

	if !id1.Equals(id1Copy) {
		t.Error("Same IDs should be equal")
	}
	if id1.Equals(id2) {
		t.Error("Different IDs should not be equal")
	}
}
