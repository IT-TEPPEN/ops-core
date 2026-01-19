package value_object

import (
	"testing"

	"github.com/google/uuid"
)

func TestNewSessionID(t *testing.T) {
	id := NewSessionID()

	if id.IsEmpty() {
		t.Error("NewSessionID should not return empty ID")
	}

	// Verify it's a valid UUID
	_, err := uuid.Parse(id.String())
	if err != nil {
		t.Errorf("NewSessionID should generate valid UUID: %v", err)
	}
}

func TestSessionIDFromString(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		wantError bool
	}{
		{
			name:      "valid UUID",
			input:     "550e8400-e29b-41d4-a716-446655440000",
			wantError: false,
		},
		{
			name:      "empty string",
			input:     "",
			wantError: true,
		},
		{
			name:      "invalid UUID format",
			input:     "invalid-uuid",
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			id, err := SessionIDFromString(tt.input)

			if tt.wantError {
				if err == nil {
					t.Error("expected error but got nil")
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
				if id.String() != tt.input {
					t.Errorf("expected %s, got %s", tt.input, id.String())
				}
			}
		})
	}
}

func TestReconstructSessionID(t *testing.T) {
	input := "550e8400-e29b-41d4-a716-446655440000"
	id := ReconstructSessionID(input)

	if id.String() != input {
		t.Errorf("expected %s, got %s", input, id.String())
	}
}

func TestSessionID_IsEmpty(t *testing.T) {
	if !SessionID("").IsEmpty() {
		t.Error("empty ID should return true")
	}

	if SessionID("550e8400-e29b-41d4-a716-446655440000").IsEmpty() {
		t.Error("non-empty ID should return false")
	}
}

func TestSessionID_Equals(t *testing.T) {
	id1 := SessionID("550e8400-e29b-41d4-a716-446655440000")
	id2 := SessionID("550e8400-e29b-41d4-a716-446655440000")
	id3 := SessionID("660e8400-e29b-41d4-a716-446655440000")

	if !id1.Equals(id2) {
		t.Error("same IDs should be equal")
	}

	if id1.Equals(id3) {
		t.Error("different IDs should not be equal")
	}
}
