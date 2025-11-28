package value_object

import (
	"errors"

	"github.com/google/uuid"
)

// AttachmentID represents a unique identifier for an attachment.
type AttachmentID string

// NewAttachmentID creates a new AttachmentID from a string.
func NewAttachmentID(id string) (AttachmentID, error) {
	if id == "" {
		return "", errors.New("attachment ID cannot be empty")
	}
	// Validate that it's a valid UUID
	if _, err := uuid.Parse(id); err != nil {
		return "", errors.New("attachment ID must be a valid UUID")
	}
	return AttachmentID(id), nil
}

// GenerateAttachmentID generates a new AttachmentID using UUID v4.
func GenerateAttachmentID() AttachmentID {
	return AttachmentID(uuid.New().String())
}

// String returns the string representation of AttachmentID.
func (a AttachmentID) String() string {
	return string(a)
}

// IsEmpty returns true if the AttachmentID is empty.
func (a AttachmentID) IsEmpty() bool {
	return string(a) == ""
}

// Equals checks if two AttachmentIDs are equal.
func (a AttachmentID) Equals(other AttachmentID) bool {
	return a == other
}
