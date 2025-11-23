package entity

import (
	"testing"
	"time"

	"opscore/backend/internal/document/domain/value_object"
)

func TestNewDocumentVersion(t *testing.T) {
	validID := value_object.GenerateVersionID()
	validDocID := value_object.GenerateDocumentID()
	validVersionNum, _ := value_object.NewVersionNumber(1)
	validPath, _ := value_object.NewFilePath("docs/test.md")
	validHash, _ := value_object.NewCommitHash("abc123def456")
	validSource, _ := value_object.NewDocumentSource(validPath, validHash)
	validDocType, _ := value_object.NewDocumentType("procedure")
	validContent := "# Test Document\n\nThis is test content."

	tests := []struct {
		name       string
		id         value_object.VersionID
		documentID value_object.DocumentID
		versionNum value_object.VersionNumber
		source     value_object.DocumentSource
		title      string
		docType    value_object.DocumentType
		content    string
		wantErr    bool
	}{
		{
			name:       "valid version",
			id:         validID,
			documentID: validDocID,
			versionNum: validVersionNum,
			source:     validSource,
			title:      "Test Document",
			docType:    validDocType,
			content:    validContent,
			wantErr:    false,
		},
		{
			name:       "empty version ID",
			id:         value_object.VersionID(""),
			documentID: validDocID,
			versionNum: validVersionNum,
			source:     validSource,
			title:      "Test",
			docType:    validDocType,
			content:    validContent,
			wantErr:    true,
		},
		{
			name:       "empty title",
			id:         validID,
			documentID: validDocID,
			versionNum: validVersionNum,
			source:     validSource,
			title:      "",
			docType:    validDocType,
			content:    validContent,
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewDocumentVersion(
				tt.id,
				tt.documentID,
				tt.versionNum,
				tt.source,
				tt.title,
				tt.docType,
				[]value_object.Tag{},
				[]value_object.VariableDefinition{},
				tt.content,
			)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewDocumentVersion() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if got == nil {
					t.Error("NewDocumentVersion() returned nil")
					return
				}
				if !got.ID().Equals(tt.id) {
					t.Errorf("ID() = %v, want %v", got.ID(), tt.id)
				}
				if got.Title() != tt.title {
					t.Errorf("Title() = %v, want %v", got.Title(), tt.title)
				}
				if !got.IsCurrentVersion() {
					t.Error("IsCurrentVersion() = false, want true for new version")
				}
				if !got.IsPublished() {
					t.Error("IsPublished() = false, want true for new version")
				}
			}
		})
	}
}

func TestDocumentVersion_MarkAsCurrent(t *testing.T) {
	validPath, _ := value_object.NewFilePath("docs/test.md")
	validHash, _ := value_object.NewCommitHash("abc123def456")
	validSource, _ := value_object.NewDocumentSource(validPath, validHash)
	validDocType, _ := value_object.NewDocumentType("procedure")

	now := time.Now()
	version := ReconstructDocumentVersion(
		value_object.GenerateVersionID(),
		value_object.GenerateDocumentID(),
		value_object.VersionNumber(1),
		validSource,
		"Test",
		validDocType,
		[]value_object.Tag{},
		[]value_object.VariableDefinition{},
		"content",
		now,
		nil,
		false, // not current
	)

	if version.IsCurrentVersion() {
		t.Error("IsCurrentVersion() = true before marking as current")
	}

	version.MarkAsCurrent()

	if !version.IsCurrentVersion() {
		t.Error("IsCurrentVersion() = false after marking as current")
	}
}

func TestDocumentVersion_Unpublish(t *testing.T) {
	validPath, _ := value_object.NewFilePath("docs/test.md")
	validHash, _ := value_object.NewCommitHash("abc123def456")
	validSource, _ := value_object.NewDocumentSource(validPath, validHash)
	validDocType, _ := value_object.NewDocumentType("procedure")

	version, _ := NewDocumentVersion(
		value_object.GenerateVersionID(),
		value_object.GenerateDocumentID(),
		value_object.VersionNumber(1),
		validSource,
		"Test",
		validDocType,
		[]value_object.Tag{},
		[]value_object.VariableDefinition{},
		"# Test",
	)

	if !version.IsPublished() {
		t.Fatal("IsPublished() = false, want true for new version")
	}

	err := version.Unpublish()
	if err != nil {
		t.Errorf("Unpublish() error = %v", err)
	}

	if version.IsPublished() {
		t.Error("IsPublished() = true after unpublishing")
	}

	if version.UnpublishedAt() == nil {
		t.Error("UnpublishedAt() = nil after unpublishing")
	}

	// Try to unpublish again
	err = version.Unpublish()
	if err == nil {
		t.Error("Unpublish() second time should return error")
	}
}
