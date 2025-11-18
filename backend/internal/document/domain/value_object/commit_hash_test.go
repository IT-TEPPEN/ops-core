package value_object

import "testing"

func TestNewCommitHash(t *testing.T) {
	tests := []struct {
		name    string
		hash    string
		wantErr bool
	}{
		{
			name:    "valid full SHA-1",
			hash:    "abc123def456789012345678901234567890abcd",
			wantErr: false,
		},
		{
			name:    "valid short hash",
			hash:    "abc1234",
			wantErr: false,
		},
		{
			name:    "minimum length hash",
			hash:    "abc1234",
			wantErr: false,
		},
		{
			name:    "too short hash",
			hash:    "abc123",
			wantErr: true,
		},
		{
			name:    "empty hash",
			hash:    "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewCommitHash(tt.hash)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewCommitHash() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got.String() != tt.hash {
				t.Errorf("NewCommitHash() = %v, want %v", got.String(), tt.hash)
			}
		})
	}
}

func TestCommitHash_Short(t *testing.T) {
	tests := []struct {
		name     string
		hash     string
		expected string
	}{
		{
			name:     "long hash",
			hash:     "abc123def456789012345678901234567890abcd",
			expected: "abc123d",
		},
		{
			name:     "already short hash",
			hash:     "abc1234",
			expected: "abc1234",
		},
		{
			name:     "exact 7 characters",
			hash:     "abcdefg",
			expected: "abcdefg",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ch, _ := NewCommitHash(tt.hash)
			if got := ch.Short(); got != tt.expected {
				t.Errorf("Short() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestCommitHash_IsEmpty(t *testing.T) {
	empty := CommitHash("")
	if !empty.IsEmpty() {
		t.Error("IsEmpty() returned false for empty commit hash")
	}

	nonEmpty, _ := NewCommitHash("abc1234")
	if nonEmpty.IsEmpty() {
		t.Error("IsEmpty() returned true for non-empty commit hash")
	}
}
