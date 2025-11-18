package value_object

import (
	"errors"
	"strings"
)

// Category represents a category for organizing documents.
type Category struct {
	name string
}

// NewCategory creates a new Category from a string.
func NewCategory(name string) (Category, error) {
	trimmed := strings.TrimSpace(name)
	if trimmed == "" {
		return Category{}, errors.New("category name cannot be empty")
	}
	if len(trimmed) > 100 {
		return Category{}, errors.New("category name cannot exceed 100 characters")
	}
	return Category{name: trimmed}, nil
}

// Name returns the category name.
func (c Category) Name() string {
	return c.name
}

// String returns the string representation of Category.
func (c Category) String() string {
	return c.name
}

// Equals checks if two Categories are equal.
func (c Category) Equals(other Category) bool {
	return c.name == other.name
}
