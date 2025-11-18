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

func TestNewCategory(t *testing.T) {
	tests := []struct {
		name         string
		categoryName string
		wantErr      bool
	}{
		{
			name:         "valid category",
			categoryName: "Database Operations",
			wantErr:      false,
		},
		{
			name:         "category with spaces trimmed",
			categoryName: "  System Maintenance  ",
			wantErr:      false,
		},
		{
			name:         "empty category",
			categoryName: "",
			wantErr:      true,
		},
		{
			name:         "only spaces",
			categoryName: "   ",
			wantErr:      true,
		},
		{
			name:         "category too long",
			categoryName: "this is a very long category name that definitely exceeds the one hundred character limit for categories in the system",
			wantErr:      true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewCategory(tt.categoryName)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewCategory() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got.Name() == "" {
				t.Error("NewCategory() returned category with empty name")
			}
		})
	}
}
