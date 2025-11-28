package value_object

import (
	"testing"

	"github.com/google/uuid"
)

func TestNewAttachmentID(t *testing.T) {
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
			got, err := NewAttachmentID(tt.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewAttachmentID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got.String() != tt.id {
				t.Errorf("NewAttachmentID() = %v, want %v", got, tt.id)
			}
		})
	}
}

func TestGenerateAttachmentID(t *testing.T) {
	id := GenerateAttachmentID()
	if id.IsEmpty() {
		t.Error("GenerateAttachmentID() returned empty ID")
	}
	// Verify it's a valid UUID
	if _, err := uuid.Parse(id.String()); err != nil {
		t.Errorf("GenerateAttachmentID() returned invalid UUID: %v", err)
	}
}

func TestAttachmentID_IsEmpty(t *testing.T) {
	tests := []struct {
		name string
		id   AttachmentID
		want bool
	}{
		{
			name: "empty ID",
			id:   AttachmentID(""),
			want: true,
		},
		{
			name: "non-empty ID",
			id:   GenerateAttachmentID(),
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.id.IsEmpty(); got != tt.want {
				t.Errorf("AttachmentID.IsEmpty() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAttachmentID_Equals(t *testing.T) {
	id1 := GenerateAttachmentID()
	id2 := GenerateAttachmentID()
	id1Copy := AttachmentID(id1.String())

	if !id1.Equals(id1Copy) {
		t.Error("Same IDs should be equal")
	}
	if id1.Equals(id2) {
		t.Error("Different IDs should not be equal")
	}
}
