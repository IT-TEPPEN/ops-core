package value_object

import (
	"testing"

	"github.com/google/uuid"
)

func TestNewExecutionRecordID(t *testing.T) {
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
			got, err := NewExecutionRecordID(tt.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewExecutionRecordID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got.String() != tt.id {
				t.Errorf("NewExecutionRecordID() = %v, want %v", got, tt.id)
			}
		})
	}
}

func TestGenerateExecutionRecordID(t *testing.T) {
	id := GenerateExecutionRecordID()
	if id.IsEmpty() {
		t.Error("GenerateExecutionRecordID() returned empty ID")
	}
	// Verify it's a valid UUID
	if _, err := uuid.Parse(id.String()); err != nil {
		t.Errorf("GenerateExecutionRecordID() returned invalid UUID: %v", err)
	}
}

func TestExecutionRecordID_IsEmpty(t *testing.T) {
	tests := []struct {
		name string
		id   ExecutionRecordID
		want bool
	}{
		{
			name: "empty ID",
			id:   ExecutionRecordID(""),
			want: true,
		},
		{
			name: "non-empty ID",
			id:   GenerateExecutionRecordID(),
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.id.IsEmpty(); got != tt.want {
				t.Errorf("ExecutionRecordID.IsEmpty() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestExecutionRecordID_Equals(t *testing.T) {
	id1 := GenerateExecutionRecordID()
	id2 := GenerateExecutionRecordID()
	id1Copy := ExecutionRecordID(id1.String())

	if !id1.Equals(id1Copy) {
		t.Error("Same IDs should be equal")
	}
	if id1.Equals(id2) {
		t.Error("Different IDs should not be equal")
	}
}
