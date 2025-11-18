package entity

import (
	"errors"
	"time"

	"opscore/backend/internal/document/domain/value_object"
)

// documentVersion represents a specific version of a document.
type documentVersion struct {
	id               value_object.VersionID
	documentID       value_object.DocumentID
	versionNumber    value_object.VersionNumber
	commitHash       value_object.CommitHash
	content          string
	publishedAt      time.Time
	unpublishedAt    *time.Time
	isCurrentVersion bool
}

// DocumentVersion is the interface for a document version.
type DocumentVersion interface {
	ID() value_object.VersionID
	DocumentID() value_object.DocumentID
	VersionNumber() value_object.VersionNumber
	CommitHash() value_object.CommitHash
	Content() string
	PublishedAt() time.Time
	UnpublishedAt() *time.Time
	IsCurrentVersion() bool
	MarkAsCurrent()
	Unpublish() error
	IsPublished() bool
}

// NewDocumentVersion creates a new DocumentVersion instance.
func NewDocumentVersion(
	id value_object.VersionID,
	documentID value_object.DocumentID,
	versionNumber value_object.VersionNumber,
	commitHash value_object.CommitHash,
	content string,
) (DocumentVersion, error) {
	if id.IsEmpty() {
		return nil, errors.New("version ID cannot be empty")
	}
	if documentID.IsEmpty() {
		return nil, errors.New("document ID cannot be empty")
	}
	if versionNumber.IsZero() {
		return nil, errors.New("version number cannot be zero")
	}
	if commitHash.IsEmpty() {
		return nil, errors.New("commit hash cannot be empty")
	}
	if content == "" {
		return nil, errors.New("content cannot be empty")
	}

	now := time.Now()
	return &documentVersion{
		id:               id,
		documentID:       documentID,
		versionNumber:    versionNumber,
		commitHash:       commitHash,
		content:          content,
		publishedAt:      now,
		unpublishedAt:    nil,
		isCurrentVersion: true,
	}, nil
}

// ReconstructDocumentVersion reconstructs a DocumentVersion from persistence data.
func ReconstructDocumentVersion(
	id value_object.VersionID,
	documentID value_object.DocumentID,
	versionNumber value_object.VersionNumber,
	commitHash value_object.CommitHash,
	content string,
	publishedAt time.Time,
	unpublishedAt *time.Time,
	isCurrentVersion bool,
) DocumentVersion {
	return &documentVersion{
		id:               id,
		documentID:       documentID,
		versionNumber:    versionNumber,
		commitHash:       commitHash,
		content:          content,
		publishedAt:      publishedAt,
		unpublishedAt:    unpublishedAt,
		isCurrentVersion: isCurrentVersion,
	}
}

// ID returns the version ID.
func (v *documentVersion) ID() value_object.VersionID {
	return v.id
}

// DocumentID returns the document ID.
func (v *documentVersion) DocumentID() value_object.DocumentID {
	return v.documentID
}

// VersionNumber returns the version number.
func (v *documentVersion) VersionNumber() value_object.VersionNumber {
	return v.versionNumber
}

// CommitHash returns the commit hash.
func (v *documentVersion) CommitHash() value_object.CommitHash {
	return v.commitHash
}

// Content returns the content.
func (v *documentVersion) Content() string {
	return v.content
}

// PublishedAt returns the published timestamp.
func (v *documentVersion) PublishedAt() time.Time {
	return v.publishedAt
}

// UnpublishedAt returns the unpublished timestamp.
func (v *documentVersion) UnpublishedAt() *time.Time {
	return v.unpublishedAt
}

// IsCurrentVersion returns whether this is the current version.
func (v *documentVersion) IsCurrentVersion() bool {
	return v.isCurrentVersion
}

// MarkAsCurrent marks this version as the current version.
func (v *documentVersion) MarkAsCurrent() {
	v.isCurrentVersion = true
}

// Unpublish unpublishes this version.
func (v *documentVersion) Unpublish() error {
	if v.unpublishedAt != nil {
		return errors.New("version is already unpublished")
	}
	now := time.Now()
	v.unpublishedAt = &now
	v.isCurrentVersion = false
	return nil
}

// IsPublished returns whether this version is currently published.
func (v *documentVersion) IsPublished() bool {
	return v.unpublishedAt == nil
}
