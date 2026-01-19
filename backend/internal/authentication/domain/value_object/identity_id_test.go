package value_object
package value_object

import (
	"testing"

	"github.com/google/uuid"
)

func TestNewIdentityID(t *testing.T) {
	id := NewIdentityID()

	if id.IsEmpty() {
		t.Error("NewIdentityID should not return empty ID")
	}

	// Verify it's a valid UUID
	_, err := uuid.Parse(id.String())
	if err != nil {
		t.Errorf("NewIdentityID should generate valid UUID: %v", err)
	}
}

func TestReconstructIdentityID(t *testing.T) {
	input := "550e8400-e29b-41d4-a716-446655440000"
	id := ReconstructIdentityID(input)

	if id.String() != input {
		t.Errorf("expected %s, got %s", input, id.String())
	}

	// ReconstructIdentityID should not validate
	invalidID := ReconstructIdentityID("invalid")
	if invalidID.String() != "invalid" {
		t.Error("ReconstructIdentityID should not validate input")
	}
}

func TestIdentityID_IsEmpty(t *testing.T) {
	if !IdentityID("").IsEmpty() {
		t.Error("empty ID should return true")
	}

	if IdentityID("550e8400-e29b-41d4-a716-446655440000").IsEmpty() {
		t.Error("non-empty ID should return false")
	}
}

func TestIdentityID_Equals(t *testing.T) {
	id1 := IdentityID("550e8400-e29b-41d4-a716-446655440000")
	id2 := IdentityID("550e8400-e29b-41d4-a716-446655440000")
	id3 := IdentityID("660e8400-e29b-41d4-a716-446655440000")

	if !id1.Equals(id2) {
		t.Error("same IDs should be equal")
	}

	if id1.Equals(id3) {
		t.Error("different IDs should not be equal")
	}
}
