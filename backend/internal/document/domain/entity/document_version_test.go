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
	validCommitHash, _ := value_object.NewCommitHash("abc123def456")
	validContent := "# Test Document\n\nThis is test content."

	tests := []struct {
		name       string
		id         value_object.VersionID
		documentID value_object.DocumentID
		versionNum value_object.VersionNumber
		commitHash value_object.CommitHash
		content    string
		wantErr    bool
	}{
		{
			name:       "valid version",
			id:         validID,
			documentID: validDocID,
			versionNum: validVersionNum,
			commitHash: validCommitHash,
			content:    validContent,
			wantErr:    false,
		},
		{
			name:       "empty version ID",
			id:         value_object.VersionID(""),
			documentID: validDocID,
			versionNum: validVersionNum,
			commitHash: validCommitHash,
			content:    validContent,
			wantErr:    true,
		},
		{
			name:       "empty document ID",
			id:         validID,
			documentID: value_object.DocumentID(""),
			versionNum: validVersionNum,
			commitHash: validCommitHash,
			content:    validContent,
			wantErr:    true,
		},
		{
			name:       "zero version number",
			id:         validID,
			documentID: validDocID,
			versionNum: value_object.VersionNumber(0),
			commitHash: validCommitHash,
			content:    validContent,
			wantErr:    true,
		},
		{
			name:       "empty commit hash",
			id:         validID,
			documentID: validDocID,
			versionNum: validVersionNum,
			commitHash: value_object.CommitHash(""),
			content:    validContent,
			wantErr:    true,
		},
		{
			name:       "empty content",
			id:         validID,
			documentID: validDocID,
			versionNum: validVersionNum,
			commitHash: validCommitHash,
			content:    "",
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewDocumentVersion(tt.id, tt.documentID, tt.versionNum, tt.commitHash, tt.content)
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
	validID := value_object.GenerateVersionID()
	validDocID := value_object.GenerateDocumentID()
	validVersionNum, _ := value_object.NewVersionNumber(1)
	validCommitHash, _ := value_object.NewCommitHash("abc123def456")

	// Create a version using reconstructor so it's not current by default
	now := time.Now()
	version := ReconstructDocumentVersion(
		validID,
		validDocID,
		validVersionNum,
		validCommitHash,
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
	validID := value_object.GenerateVersionID()
	validDocID := value_object.GenerateDocumentID()
	validVersionNum, _ := value_object.NewVersionNumber(1)
	validCommitHash, _ := value_object.NewCommitHash("abc123def456")
	validContent := "# Test"

	version, _ := NewDocumentVersion(validID, validDocID, validVersionNum, validCommitHash, validContent)

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

	if version.IsCurrentVersion() {
		t.Error("IsCurrentVersion() = true after unpublishing")
	}

	// Try to unpublish again
	err = version.Unpublish()
	if err == nil {
		t.Error("Unpublish() second time should return error")
	}
}

func TestReconstructDocumentVersion(t *testing.T) {
	id := value_object.GenerateVersionID()
	docID := value_object.GenerateDocumentID()
	versionNum, _ := value_object.NewVersionNumber(2)
	commitHash, _ := value_object.NewCommitHash("abc123")
	content := "test content"
	publishedAt := time.Now().Add(-24 * time.Hour)
	unpublishedAt := time.Now()

	version := ReconstructDocumentVersion(
		id,
		docID,
		versionNum,
		commitHash,
		content,
		publishedAt,
		&unpublishedAt,
		false,
	)

	if version == nil {
		t.Fatal("ReconstructDocumentVersion() returned nil")
	}

	if !version.ID().Equals(id) {
		t.Errorf("ID() = %v, want %v", version.ID(), id)
	}

	if !version.DocumentID().Equals(docID) {
		t.Errorf("DocumentID() = %v, want %v", version.DocumentID(), docID)
	}

	if version.IsPublished() {
		t.Error("IsPublished() = true for unpublished version")
	}

	if version.IsCurrentVersion() {
		t.Error("IsCurrentVersion() = true when reconstructed as false")
	}
}
