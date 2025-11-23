package value_object

import "testing"

func TestNewTag(t *testing.T) {
	tests := []struct {
		name    string
		tagName string
		wantErr bool
	}{
		{
			name:    "valid tag",
			tagName: "database",
			wantErr: false,
		},
		{
			name:    "tag with spaces trimmed",
			tagName: "  backup  ",
			wantErr: false,
		},
		{
			name:    "empty tag",
			tagName: "",
			wantErr: true,
		},
		{
			name:    "only spaces",
			tagName: "   ",
			wantErr: true,
		},
		{
			name:    "tag too long",
			tagName: "this-is-a-very-long-tag-name-that-exceeds-fifty-characters-limit",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewTag(tt.tagName)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewTag() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got.Name() == "" {
				t.Error("NewTag() returned tag with empty name")
			}
		})
	}
}

func TestTag_Equals(t *testing.T) {
	tag1, _ := NewTag("database")
	tag2, _ := NewTag("database")
	tag3, _ := NewTag("backup")

	if !tag1.Equals(tag2) {
		t.Error("Equals() returned false for identical tags")
	}
	if tag1.Equals(tag3) {
		t.Error("Equals() returned true for different tags")
	}
}
