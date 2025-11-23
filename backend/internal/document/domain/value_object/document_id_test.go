package value_object

import (
	"testing"

	"github.com/google/uuid"
)

func TestNewDocumentID(t *testing.T) {
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
		{
			name:    "valid UUID with uppercase",
			id:      "A1B2C3D4-E5F6-4789-ABCD-EF0123456789",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewDocumentID(tt.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewDocumentID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got.String() != tt.id {
				t.Errorf("NewDocumentID() = %v, want %v", got, tt.id)
			}
		})
	}
}

func TestGenerateDocumentID(t *testing.T) {
	id1 := GenerateDocumentID()
	id2 := GenerateDocumentID()

	if id1.IsEmpty() {
		t.Error("GenerateDocumentID() returned empty ID")
	}
	if id1.Equals(id2) {
		t.Error("GenerateDocumentID() generated duplicate IDs")
	}

	// Verify it's a valid UUID
	_, err := uuid.Parse(id1.String())
	if err != nil {
		t.Errorf("GenerateDocumentID() generated invalid UUID: %v", err)
	}
}

func TestDocumentID_IsEmpty(t *testing.T) {
	empty := DocumentID("")
	if !empty.IsEmpty() {
		t.Error("IsEmpty() returned false for empty DocumentID")
	}

	notEmpty := GenerateDocumentID()
	if notEmpty.IsEmpty() {
		t.Error("IsEmpty() returned true for non-empty DocumentID")
	}
}

func TestDocumentID_Equals(t *testing.T) {
	id := uuid.New().String()
	docID1, _ := NewDocumentID(id)
	docID2, _ := NewDocumentID(id)
	docID3 := GenerateDocumentID()

	if !docID1.Equals(docID2) {
		t.Error("Equals() returned false for identical DocumentIDs")
	}
	if docID1.Equals(docID3) {
		t.Error("Equals() returned true for different DocumentIDs")
	}
}
