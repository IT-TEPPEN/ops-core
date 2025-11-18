package value_object

import (
	"errors"
	"path/filepath"
	"strings"
)

// FilePath represents a file path within a repository.
type FilePath string

// NewFilePath creates a new FilePath from a string.
func NewFilePath(path string) (FilePath, error) {
	if path == "" {
		return "", errors.New("file path cannot be empty")
	}
	
	// Clean the path to normalize it
	cleaned := filepath.Clean(path)
	
	// Ensure it doesn't try to escape the repository root
	if strings.HasPrefix(cleaned, "..") {
		return "", errors.New("file path cannot escape repository root")
	}
	
	return FilePath(cleaned), nil
}

// String returns the string representation of FilePath.
func (f FilePath) String() string {
	return string(f)
}

// IsEmpty returns true if the FilePath is empty.
func (f FilePath) IsEmpty() bool {
	return string(f) == ""
}

// Equals checks if two FilePaths are equal.
func (f FilePath) Equals(other FilePath) bool {
	return f == other
}

// Extension returns the file extension.
func (f FilePath) Extension() string {
	return filepath.Ext(string(f))
}

// IsMarkdown returns true if the file is a Markdown file.
func (f FilePath) IsMarkdown() bool {
	ext := strings.ToLower(f.Extension())
	return ext == ".md" || ext == ".markdown"
}
