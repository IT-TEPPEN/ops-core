package value_object

import "errors"

// VersionNumber represents a version number (sequential starting from 1).
type VersionNumber int

// NewVersionNumber creates a new VersionNumber from an int.
func NewVersionNumber(num int) (VersionNumber, error) {
	if num < 1 {
		return 0, errors.New("version number must be at least 1")
	}
	return VersionNumber(num), nil
}

// Int returns the int representation of VersionNumber.
func (v VersionNumber) Int() int {
	return int(v)
}

// IsZero returns true if the VersionNumber is zero (uninitialized).
func (v VersionNumber) IsZero() bool {
	return v == 0
}

// Equals checks if two VersionNumbers are equal.
func (v VersionNumber) Equals(other VersionNumber) bool {
	return v == other
}

// Next returns the next version number.
func (v VersionNumber) Next() VersionNumber {
	return VersionNumber(v.Int() + 1)
}

// Previous returns the previous version number if valid.
func (v VersionNumber) Previous() (VersionNumber, error) {
	if v <= 1 {
		return 0, errors.New("cannot get previous version of version 1")
	}
	return VersionNumber(v.Int() - 1), nil
}
