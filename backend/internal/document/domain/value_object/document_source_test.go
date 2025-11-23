package value_object

import "testing"

func TestNewDocumentSource(t *testing.T) {
	validPath, _ := NewFilePath("docs/test.md")
	validHash, _ := NewCommitHash("abc1234567")

	tests := []struct {
		name       string
		filePath   FilePath
		commitHash CommitHash
		wantErr    bool
	}{
		{
			name:       "valid document source",
			filePath:   validPath,
			commitHash: validHash,
			wantErr:    false,
		},
		{
			name:       "empty file path",
			filePath:   FilePath(""),
			commitHash: validHash,
			wantErr:    true,
		},
		{
			name:       "empty commit hash",
			filePath:   validPath,
			commitHash: CommitHash(""),
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewDocumentSource(tt.filePath, tt.commitHash)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewDocumentSource() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if !got.FilePath().Equals(tt.filePath) {
					t.Errorf("FilePath() = %v, want %v", got.FilePath(), tt.filePath)
				}
				if !got.CommitHash().Equals(tt.commitHash) {
					t.Errorf("CommitHash() = %v, want %v", got.CommitHash(), tt.commitHash)
				}
			}
		})
	}
}

func TestDocumentSource_Equals(t *testing.T) {
	path1, _ := NewFilePath("docs/test.md")
	path2, _ := NewFilePath("docs/other.md")
	hash1, _ := NewCommitHash("abc1234567")
	hash2, _ := NewCommitHash("def7890123")

	source1, _ := NewDocumentSource(path1, hash1)
	source2, _ := NewDocumentSource(path1, hash1)
	source3, _ := NewDocumentSource(path2, hash1)
	source4, _ := NewDocumentSource(path1, hash2)

	if !source1.Equals(source2) {
		t.Error("Equals() returned false for identical sources")
	}
	if source1.Equals(source3) {
		t.Error("Equals() returned true for different file paths")
	}
	if source1.Equals(source4) {
		t.Error("Equals() returned true for different commit hashes")
	}
}

func TestDocumentSource_String(t *testing.T) {
	path, _ := NewFilePath("docs/test.md")
	hash, _ := NewCommitHash("abc1234567")
	source, _ := NewDocumentSource(path, hash)

	expected := "docs/test.md@abc1234"
	if got := source.String(); got != expected {
		t.Errorf("String() = %v, want %v", got, expected)
	}
}
