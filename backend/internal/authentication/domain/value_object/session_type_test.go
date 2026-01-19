package value_object
package value_object

import (
	"testing"
	"time"
)

func TestSessionTypeFromString(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		expected  SessionType
		wantError bool
	}{
		{
			name:      "short session",
			input:     "short",
			expected:  SessionTypeShort,
			wantError: false,
		},
		{
			name:      "long session",
			input:     "long",
			expected:  SessionTypeLong,
			wantError: false,
		},
		{
			name:      "invalid type",
			input:     "invalid",
			wantError: true,
		},
		{
			name:      "empty string",
			input:     "",
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			st, err := SessionTypeFromString(tt.input)

			if tt.wantError {
				if err == nil {
					t.Error("expected error but got nil")
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
				if st != tt.expected {
					t.Errorf("expected %s, got %s", tt.expected, st)
				}
			}
		})
	}
}

func TestReconstructSessionType(t *testing.T) {
	st := ReconstructSessionType("short")
	if st != SessionTypeShort {
		t.Errorf("expected %s, got %s", SessionTypeShort, st)
	}

	// ReconstructSessionType should not validate
	invalidST := ReconstructSessionType("invalid")
	if invalidST != SessionType("invalid") {
		t.Error("ReconstructSessionType should not validate input")
	}
}

func TestSessionType_ExpirationDuration(t *testing.T) {
	tests := []struct {
		name     string
		st       SessionType
		expected time.Duration
	}{
		{
			name:     "short session",
			st:       SessionTypeShort,
			expected: 7 * 24 * time.Hour,
		},
		{
			name:     "long session",
			st:       SessionTypeLong,
			expected: 30 * 24 * time.Hour,
		},
		{
			name:     "invalid type defaults to short",
			st:       SessionType("invalid"),
			expected: 7 * 24 * time.Hour,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if duration := tt.st.ExpirationDuration(); duration != tt.expected {
				t.Errorf("expected %v, got %v", tt.expected, duration)
			}
		})
	}
}

func TestSessionType_RememberMe(t *testing.T) {
	if SessionTypeShort.RememberMe() {
		t.Error("short session should not be remember me")
	}

	if !SessionTypeLong.RememberMe() {
		t.Error("long session should be remember me")
	}
}

func TestSessionType_String(t *testing.T) {
	if SessionTypeShort.String() != "short" {
		t.Errorf("expected 'short', got %s", SessionTypeShort.String())
	}

	if SessionTypeLong.String() != "long" {
		t.Errorf("expected 'long', got %s", SessionTypeLong.String())
	}
}
