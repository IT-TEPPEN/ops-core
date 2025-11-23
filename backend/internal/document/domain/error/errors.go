package error

import "errors"

var (
	// ErrDocumentNotFound is returned when a document is not found.
	ErrDocumentNotFound = errors.New("document not found")

	// ErrVersionNotFound is returned when a document version is not found.
	ErrVersionNotFound = errors.New("version not found")

	// ErrDocumentAlreadyExists is returned when attempting to create a duplicate document.
	ErrDocumentAlreadyExists = errors.New("document already exists")

	// ErrInvalidDocumentID is returned when a document ID is invalid.
	ErrInvalidDocumentID = errors.New("invalid document ID")

	// ErrInvalidVersionNumber is returned when a version number is invalid.
	ErrInvalidVersionNumber = errors.New("invalid version number")

	// ErrDocumentNotPublished is returned when attempting an operation on an unpublished document.
	ErrDocumentNotPublished = errors.New("document is not published")

	// ErrDocumentAlreadyPublished is returned when attempting to publish an already published document.
	ErrDocumentAlreadyPublished = errors.New("document is already published")

	// ErrVersionAlreadyUnpublished is returned when attempting to unpublish an already unpublished version.
	ErrVersionAlreadyUnpublished = errors.New("version is already unpublished")

	// ErrVersionMismatch is returned when a version does not belong to the document.
	ErrVersionMismatch = errors.New("version does not belong to this document")

	// ErrEmptyContent is returned when document content is empty.
	ErrEmptyContent = errors.New("content cannot be empty")

	// ErrInvalidMetadata is returned when document metadata is invalid.
	ErrInvalidMetadata = errors.New("invalid document metadata")
)
