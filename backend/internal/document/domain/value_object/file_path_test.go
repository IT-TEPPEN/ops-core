package value_object

import "testing"

func TestNewFilePath(t *testing.T) {
	tests := []struct {
		name    string
		path    string
		wantErr bool
	}{
		{
			name:    "valid path",
			path:    "docs/test.md",
			wantErr: false,
		},
		{
			name:    "valid absolute path",
			path:    "/docs/procedures/backup.md",
			wantErr: false,
		},
		{
			name:    "empty path",
			path:    "",
			wantErr: true,
		},
		{
			name:    "path with parent directory escape",
			path:    "../../../etc/passwd",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewFilePath(tt.path)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewFilePath() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got.IsEmpty() {
				t.Error("NewFilePath() returned empty path")
			}
		})
	}
}

func TestFilePath_IsMarkdown(t *testing.T) {
	tests := []struct {
		name     string
		path     string
		expected bool
	}{
		{"markdown .md", "docs/test.md", true},
		{"markdown .markdown", "docs/test.markdown", true},
		{"uppercase .MD", "docs/test.MD", true},
		{"not markdown .txt", "docs/test.txt", false},
		{"not markdown .go", "src/main.go", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fp, _ := NewFilePath(tt.path)
			if got := fp.IsMarkdown(); got != tt.expected {
				t.Errorf("IsMarkdown() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestFilePath_Extension(t *testing.T) {
	fp, _ := NewFilePath("docs/test.md")
	if ext := fp.Extension(); ext != ".md" {
		t.Errorf("Extension() = %v, want .md", ext)
	}
}
