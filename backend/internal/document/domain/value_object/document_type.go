package value_object

import "errors"

// DocumentType represents the type of a document.
type DocumentType string

const (
	// DocumentTypeProcedure represents a procedure document.
	DocumentTypeProcedure DocumentType = "procedure"
	// DocumentTypeKnowledge represents a knowledge document.
	DocumentTypeKnowledge DocumentType = "knowledge"
)

// NewDocumentType creates a new DocumentType from a string.
func NewDocumentType(t string) (DocumentType, error) {
	docType := DocumentType(t)
	if !docType.IsValid() {
		return "", errors.New("invalid document type: must be 'procedure' or 'knowledge'")
	}
	return docType, nil
}

// IsValid checks if the DocumentType is valid.
func (d DocumentType) IsValid() bool {
	return d == DocumentTypeProcedure || d == DocumentTypeKnowledge
}

// String returns the string representation of DocumentType.
func (d DocumentType) String() string {
	return string(d)
}

// Equals checks if two DocumentTypes are equal.
func (d DocumentType) Equals(other DocumentType) bool {
	return d == other
}

// IsProcedure returns true if the document type is procedure.
func (d DocumentType) IsProcedure() bool {
	return d == DocumentTypeProcedure
}

// IsKnowledge returns true if the document type is knowledge.
func (d DocumentType) IsKnowledge() bool {
	return d == DocumentTypeKnowledge
}
