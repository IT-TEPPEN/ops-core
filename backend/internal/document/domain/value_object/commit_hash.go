package value_object

import "errors"

// CommitHash represents a Git commit hash.
type CommitHash string

// NewCommitHash creates a new CommitHash from a string.
func NewCommitHash(hash string) (CommitHash, error) {
	if hash == "" {
		return "", errors.New("commit hash cannot be empty")
	}
	// Git commit hashes are typically 40 characters (SHA-1) or 64 characters (SHA-256)
	// but we also accept short hashes (at least 7 characters)
	if len(hash) < 7 {
		return "", errors.New("commit hash must be at least 7 characters")
	}
	return CommitHash(hash), nil
}

// String returns the string representation of CommitHash.
func (c CommitHash) String() string {
	return string(c)
}

// IsEmpty returns true if the CommitHash is empty.
func (c CommitHash) IsEmpty() bool {
	return string(c) == ""
}

// Equals checks if two CommitHashes are equal.
func (c CommitHash) Equals(other CommitHash) bool {
	return c == other
}

// Short returns the first 7 characters of the commit hash.
func (c CommitHash) Short() string {
	if len(c) <= 7 {
		return string(c)
	}
	return string(c)[:7]
}
