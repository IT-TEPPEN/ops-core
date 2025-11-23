package entity

import (
	"errors"
	"time"

	"opscore/backend/internal/document/domain/value_object"
)

// document represents a document in the system (aggregate root).
type document struct {
	id             value_object.DocumentID
	repositoryID   value_object.RepositoryID
	owner          string
	isPublished    bool
	isAutoUpdate   bool
	accessScope    value_object.AccessScope
	currentVersion DocumentVersion
	versions       []DocumentVersion
	createdAt      time.Time
	updatedAt      time.Time
}

// Document is the interface for a document (aggregate root).
type Document interface {
	ID() value_object.DocumentID
	RepositoryID() value_object.RepositoryID
	Owner() string
	IsPublished() bool
	IsAutoUpdate() bool
	AccessScope() value_object.AccessScope
	CurrentVersion() DocumentVersion
	Versions() []DocumentVersion
	CreatedAt() time.Time
	UpdatedAt() time.Time

	// Behaviors
	Publish(
		source value_object.DocumentSource,
		title string,
		docType value_object.DocumentType,
		tags []value_object.Tag,
		variables []value_object.VariableDefinition,
		content string,
	) error
	Unpublish() error
	UpdateAccessScope(scope value_object.AccessScope) error
	EnableAutoUpdate()
	DisableAutoUpdate()
	RollbackToVersion(versionNumber value_object.VersionNumber) error
	AddVersion(version DocumentVersion) error
}

// NewDocument creates a new Document instance.
func NewDocument(
	id value_object.DocumentID,
	repositoryID value_object.RepositoryID,
	owner string,
	accessScope value_object.AccessScope,
) (Document, error) {
	if id.IsEmpty() {
		return nil, errors.New("document ID cannot be empty")
	}
	if repositoryID.IsEmpty() {
		return nil, errors.New("repository ID cannot be empty")
	}
	if owner == "" {
		return nil, errors.New("owner cannot be empty")
	}
	if !accessScope.IsValid() {
		return nil, errors.New("invalid access scope")
	}

	now := time.Now()
	return &document{
		id:           id,
		repositoryID: repositoryID,
		owner:        owner,
		isPublished:  false,
		isAutoUpdate: false,
		accessScope:  accessScope,
		versions:     []DocumentVersion{},
		createdAt:    now,
		updatedAt:    now,
	}, nil
}

// ReconstructDocument reconstructs a Document from persistence data.
func ReconstructDocument(
	id value_object.DocumentID,
	repositoryID value_object.RepositoryID,
	owner string,
	isPublished bool,
	isAutoUpdate bool,
	accessScope value_object.AccessScope,
	currentVersion DocumentVersion,
	versions []DocumentVersion,
	createdAt time.Time,
	updatedAt time.Time,
) Document {
	return &document{
		id:             id,
		repositoryID:   repositoryID,
		owner:          owner,
		isPublished:    isPublished,
		isAutoUpdate:   isAutoUpdate,
		accessScope:    accessScope,
		currentVersion: currentVersion,
		versions:       versions,
		createdAt:      createdAt,
		updatedAt:      updatedAt,
	}
}

// Getter methods
func (d *document) ID() value_object.DocumentID {
	return d.id
}

func (d *document) RepositoryID() value_object.RepositoryID {
	return d.repositoryID
}

func (d *document) Owner() string {
	return d.owner
}

func (d *document) IsPublished() bool {
	return d.isPublished
}

func (d *document) IsAutoUpdate() bool {
	return d.isAutoUpdate
}

func (d *document) AccessScope() value_object.AccessScope {
	return d.accessScope
}

func (d *document) CurrentVersion() DocumentVersion {
	return d.currentVersion
}

func (d *document) Versions() []DocumentVersion {
	return d.versions
}

func (d *document) CreatedAt() time.Time {
	return d.createdAt
}

func (d *document) UpdatedAt() time.Time {
	return d.updatedAt
}

// Publish publishes a new version of the document.
func (d *document) Publish(
	source value_object.DocumentSource,
	title string,
	docType value_object.DocumentType,
	tags []value_object.Tag,
	variables []value_object.VariableDefinition,
	content string,
) error {
	if source.FilePath().IsEmpty() {
		return errors.New("source file path cannot be empty")
	}
	if source.CommitHash().IsEmpty() {
		return errors.New("source commit hash cannot be empty")
	}
	if title == "" {
		return errors.New("title cannot be empty")
	}
	if !docType.IsValid() {
		return errors.New("invalid document type")
	}
	if content == "" {
		return errors.New("content cannot be empty")
	}

	// Determine the next version number
	var nextVersionNumber value_object.VersionNumber
	if len(d.versions) == 0 {
		var err error
		nextVersionNumber, err = value_object.NewVersionNumber(1)
		if err != nil {
			return err
		}
	} else {
		// Find the highest version number
		maxVersion := d.versions[0].VersionNumber()
		for _, v := range d.versions {
			if v.VersionNumber().Int() > maxVersion.Int() {
				maxVersion = v.VersionNumber()
			}
		}
		nextVersionNumber = maxVersion.Next()
	}

	// Create a new version
	versionID := value_object.GenerateVersionID()
	newVersion, err := NewDocumentVersion(
		versionID,
		d.id,
		nextVersionNumber,
		source,
		title,
		docType,
		tags,
		variables,
		content,
	)
	if err != nil {
		return err
	}

	// Add the new version and set as current
	d.versions = append(d.versions, newVersion)
	d.currentVersion = newVersion
	d.isPublished = true
	d.updatedAt = time.Now()

	return nil
}

// Unpublish unpublishes the document.
func (d *document) Unpublish() error {
	if !d.isPublished {
		return errors.New("document is not published")
	}

	if d.currentVersion != nil {
		if err := d.currentVersion.Unpublish(); err != nil {
			return err
		}
	}

	d.isPublished = false
	d.updatedAt = time.Now()
	return nil
}

// UpdateAccessScope updates the document's access scope.
func (d *document) UpdateAccessScope(scope value_object.AccessScope) error {
	if !scope.IsValid() {
		return errors.New("invalid access scope")
	}

	d.accessScope = scope
	d.updatedAt = time.Now()
	return nil
}

// EnableAutoUpdate enables automatic updates for the document.
func (d *document) EnableAutoUpdate() {
	d.isAutoUpdate = true
	d.updatedAt = time.Now()
}

// DisableAutoUpdate disables automatic updates for the document.
func (d *document) DisableAutoUpdate() {
	d.isAutoUpdate = false
	d.updatedAt = time.Now()
}

// RollbackToVersion rolls back the document to a specific version.
func (d *document) RollbackToVersion(versionNumber value_object.VersionNumber) error {
	if !d.isPublished {
		return errors.New("cannot rollback unpublished document")
	}

	// Find the version
	var targetVersion DocumentVersion
	for _, v := range d.versions {
		if v.VersionNumber().Equals(versionNumber) {
			targetVersion = v
			break
		}
	}

	if targetVersion == nil {
		return errors.New("version not found")
	}

	if !targetVersion.IsPublished() {
		return errors.New("cannot rollback to unpublished version")
	}

	// Mark the target version as current
	targetVersion.MarkAsCurrent()
	d.currentVersion = targetVersion
	d.updatedAt = time.Now()

	return nil
}

// AddVersion adds a version to the document (used when reconstructing from persistence).
func (d *document) AddVersion(version DocumentVersion) error {
	if version == nil {
		return errors.New("version cannot be nil")
	}
	if !version.DocumentID().Equals(d.id) {
		return errors.New("version does not belong to this document")
	}

	d.versions = append(d.versions, version)
	if version.IsCurrentVersion() {
		d.currentVersion = version
	}

	return nil
}
