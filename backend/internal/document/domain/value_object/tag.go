package value_object

import (
	"errors"
	"strings"
)

// Tag represents a tag for categorizing documents.
type Tag struct {
	name string
}

// NewTag creates a new Tag from a string.
func NewTag(name string) (Tag, error) {
	trimmed := strings.TrimSpace(name)
	if trimmed == "" {
		return Tag{}, errors.New("tag name cannot be empty")
	}
	if len(trimmed) > 50 {
		return Tag{}, errors.New("tag name cannot exceed 50 characters")
	}
	return Tag{name: trimmed}, nil
}

// Name returns the tag name.
func (t Tag) Name() string {
	return t.name
}

// String returns the string representation of Tag.
func (t Tag) String() string {
	return t.name
}

// Equals checks if two Tags are equal.
func (t Tag) Equals(other Tag) bool {
	return t.name == other.name
}
