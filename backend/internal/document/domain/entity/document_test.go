package entity

import (
	"testing"

	"opscore/backend/internal/document/domain/value_object"
)

func TestNewDocument(t *testing.T) {
	validID := value_object.GenerateDocumentID()
	validRepoID, _ := value_object.NewRepositoryID("a1b2c3d4-e5f6-4789-abcd-ef0123456789")
	validScope, _ := value_object.NewAccessScope("public")

	tests := []struct {
		name         string
		id           value_object.DocumentID
		repositoryID value_object.RepositoryID
		owner        string
		accessScope  value_object.AccessScope
		wantErr      bool
	}{
		{
			name:         "valid document",
			id:           validID,
			repositoryID: validRepoID,
			owner:        "admin",
			accessScope:  validScope,
			wantErr:      false,
		},
		{
			name:         "empty document ID",
			id:           value_object.DocumentID(""),
			repositoryID: validRepoID,
			owner:        "admin",
			accessScope:  validScope,
			wantErr:      true,
		},
		{
			name:         "empty owner",
			id:           validID,
			repositoryID: validRepoID,
			owner:        "",
			accessScope:  validScope,
			wantErr:      true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewDocument(
				tt.id,
				tt.repositoryID,
				tt.owner,
				tt.accessScope,
			)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewDocument() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if got == nil {
					t.Error("NewDocument() returned nil")
					return
				}
				if got.IsPublished() {
					t.Error("IsPublished() = true for new document")
				}
				if got.IsAutoUpdate() {
					t.Error("IsAutoUpdate() = true for new document")
				}
				if len(got.Versions()) != 0 {
					t.Errorf("Versions() length = %d, want 0", len(got.Versions()))
				}
			}
		})
	}
}

func TestDocument_Publish(t *testing.T) {
	doc := createTestDocument(t)
	
	path, _ := value_object.NewFilePath("docs/test.md")
	hash, _ := value_object.NewCommitHash("abc1234567")
	source, _ := value_object.NewDocumentSource(path, hash)
	docType, _ := value_object.NewDocumentType("procedure")
	content := "# Test\n\nThis is a test document."

	err := doc.Publish(source, "Test Title", docType, []value_object.Tag{}, []value_object.VariableDefinition{}, content)
	if err != nil {
		t.Fatalf("Publish() error = %v", err)
	}

	if !doc.IsPublished() {
		t.Error("IsPublished() = false after publishing")
	}

	if len(doc.Versions()) != 1 {
		t.Errorf("Versions() length = %d, want 1", len(doc.Versions()))
	}

	currentVersion := doc.CurrentVersion()
	if currentVersion == nil {
		t.Fatal("CurrentVersion() returned nil after publishing")
	}

	if currentVersion.VersionNumber().Int() != 1 {
		t.Errorf("Version number = %d, want 1", currentVersion.VersionNumber().Int())
	}

	// Publish a second version
	path2, _ := value_object.NewFilePath("docs/test.md")
	hash2, _ := value_object.NewCommitHash("def4567890")
	source2, _ := value_object.NewDocumentSource(path2, hash2)
	content2 := "# Test v2\n\nUpdated content."

	err = doc.Publish(source2, "Test Title v2", docType, []value_object.Tag{}, []value_object.VariableDefinition{}, content2)
	if err != nil {
		t.Fatalf("Publish() second version error = %v", err)
	}

	if len(doc.Versions()) != 2 {
		t.Errorf("Versions() length = %d, want 2", len(doc.Versions()))
	}

	currentVersion = doc.CurrentVersion()
	if currentVersion.VersionNumber().Int() != 2 {
		t.Errorf("Version number = %d, want 2", currentVersion.VersionNumber().Int())
	}
}

func TestDocument_Unpublish(t *testing.T) {
	doc := createTestDocument(t)

	// Try to unpublish when not published
	err := doc.Unpublish()
	if err == nil {
		t.Error("Unpublish() should return error for unpublished document")
	}

	// Publish first
	path, _ := value_object.NewFilePath("docs/test.md")
	hash, _ := value_object.NewCommitHash("abc1234567")
	source, _ := value_object.NewDocumentSource(path, hash)
	docType, _ := value_object.NewDocumentType("procedure")
	
	err = doc.Publish(source, "Title", docType, []value_object.Tag{}, []value_object.VariableDefinition{}, "content")
	if err != nil {
		t.Fatalf("Publish() error = %v", err)
	}

	// Now unpublish
	err = doc.Unpublish()
	if err != nil {
		t.Errorf("Unpublish() error = %v", err)
	}

	if doc.IsPublished() {
		t.Error("IsPublished() = true after unpublishing")
	}
}

func TestDocument_UpdateAccessScope(t *testing.T) {
	doc := createTestDocument(t)

	newScope, _ := value_object.NewAccessScope("private")
	err := doc.UpdateAccessScope(newScope)
	if err != nil {
		t.Errorf("UpdateAccessScope() error = %v", err)
	}

	if !doc.AccessScope().Equals(newScope) {
		t.Errorf("AccessScope() = %v, want %v", doc.AccessScope(), newScope)
	}
}

func TestDocument_AutoUpdate(t *testing.T) {
	doc := createTestDocument(t)

	if doc.IsAutoUpdate() {
		t.Error("IsAutoUpdate() = true for new document")
	}

	doc.EnableAutoUpdate()
	if !doc.IsAutoUpdate() {
		t.Error("IsAutoUpdate() = false after EnableAutoUpdate()")
	}

	doc.DisableAutoUpdate()
	if doc.IsAutoUpdate() {
		t.Error("IsAutoUpdate() = true after DisableAutoUpdate()")
	}
}

func TestDocument_RollbackToVersion(t *testing.T) {
	doc := createTestDocument(t)

	// Try to rollback unpublished document
	v1, _ := value_object.NewVersionNumber(1)
	err := doc.RollbackToVersion(v1)
	if err == nil {
		t.Error("RollbackToVersion() should return error for unpublished document")
	}

	// Publish two versions
	path, _ := value_object.NewFilePath("docs/test.md")
	hash1, _ := value_object.NewCommitHash("abc1234567")
	source1, _ := value_object.NewDocumentSource(path, hash1)
	docType, _ := value_object.NewDocumentType("procedure")
	
	err = doc.Publish(source1, "Title v1", docType, []value_object.Tag{}, []value_object.VariableDefinition{}, "content v1")
	if err != nil {
		t.Fatalf("Publish() error = %v", err)
	}

	hash2, _ := value_object.NewCommitHash("def7890123")
	source2, _ := value_object.NewDocumentSource(path, hash2)
	err = doc.Publish(source2, "Title v2", docType, []value_object.Tag{}, []value_object.VariableDefinition{}, "content v2")
	if err != nil {
		t.Fatalf("Publish() second version error = %v", err)
	}

	// Rollback to version 1
	err = doc.RollbackToVersion(v1)
	if err != nil {
		t.Errorf("RollbackToVersion() error = %v", err)
	}

	if doc.CurrentVersion().VersionNumber().Int() != 1 {
		t.Errorf("CurrentVersion() = %d, want 1", doc.CurrentVersion().VersionNumber().Int())
	}

	// Try to rollback to non-existent version
	v99, _ := value_object.NewVersionNumber(99)
	err = doc.RollbackToVersion(v99)
	if err == nil {
		t.Error("RollbackToVersion() should return error for non-existent version")
	}
}

// Helper function to create a test document
func createTestDocument(t *testing.T) Document {
	id := value_object.GenerateDocumentID()
	repoID, _ := value_object.NewRepositoryID("a1b2c3d4-e5f6-4789-abcd-ef0123456789")
	accessScope, _ := value_object.NewAccessScope("public")

	doc, err := NewDocument(
		id,
		repoID,
		"admin",
		accessScope,
	)

	if err != nil {
		t.Fatalf("Failed to create test document: %v", err)
	}

	return doc
}
