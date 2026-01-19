package value_object

import (
	"testing"

	"github.com/google/uuid"
)

func TestNewUserID(t *testing.T) {
	id := NewUserID()

	if id.IsEmpty() {
		t.Error("NewUserID should not return empty ID")
	}

	// Verify it's a valid UUID
	_, err := uuid.Parse(id.String())
	if err != nil {
		t.Errorf("NewUserID should generate valid UUID: %v", err)
	}
}

func TestUserIDFromString(t *testing.T) {
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
		{
			name:      "not a UUID",
			input:     "12345",
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			id, err := UserIDFromString(tt.input)

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

func TestReconstructUserID(t *testing.T) {
	input := "550e8400-e29b-41d4-a716-446655440000"
	id := ReconstructUserID(input)

	if id.String() != input {
		t.Errorf("expected %s, got %s", input, id.String())
	}

	// ReconstructUserID should not validate
	invalidID := ReconstructUserID("invalid")
	if invalidID.String() != "invalid" {
		t.Error("ReconstructUserID should not validate input")
	}
}

func TestUserID_IsEmpty(t *testing.T) {
	tests := []struct {
		name     string
		id       UserID
		expected bool
	}{
		{
			name:     "empty ID",
			id:       UserID(""),
			expected: true,
		},
		{
			name:     "non-empty ID",
			id:       UserID("550e8400-e29b-41d4-a716-446655440000"),
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if result := tt.id.IsEmpty(); result != tt.expected {
				t.Errorf("expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestUserID_Equals(t *testing.T) {
	id1 := UserID("550e8400-e29b-41d4-a716-446655440000")
	id2 := UserID("550e8400-e29b-41d4-a716-446655440000")
	id3 := UserID("660e8400-e29b-41d4-a716-446655440000")

	if !id1.Equals(id2) {
		t.Error("same IDs should be equal")
	}

	if id1.Equals(id3) {
		t.Error("different IDs should not be equal")
	}
}

func TestUserID_String(t *testing.T) {
	expected := "550e8400-e29b-41d4-a716-446655440000"
	id := UserID(expected)

	if id.String() != expected {
		t.Errorf("expected %s, got %s", expected, id.String())
	}
}
