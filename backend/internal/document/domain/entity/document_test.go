package entity

import (
	"testing"

	"opscore/backend/internal/document/domain/value_object"
)

func TestNewDocument(t *testing.T) {
	validID := value_object.GenerateDocumentID()
	validRepoID, _ := value_object.NewRepositoryID("a1b2c3d4-e5f6-4789-abcd-ef0123456789")
	validPath, _ := value_object.NewFilePath("docs/test.md")
	validType, _ := value_object.NewDocumentType("procedure")
	validCategory, _ := value_object.NewCategory("Database")
	validScope, _ := value_object.NewAccessScope("public")

	tests := []struct {
		name         string
		id           value_object.DocumentID
		repositoryID value_object.RepositoryID
		filePath     value_object.FilePath
		title        string
		owner        string
		docType      value_object.DocumentType
		tags         []value_object.Tag
		category     value_object.Category
		variables    []value_object.VariableDefinition
		accessScope  value_object.AccessScope
		wantErr      bool
	}{
		{
			name:         "valid document",
			id:           validID,
			repositoryID: validRepoID,
			filePath:     validPath,
			title:        "Test Document",
			owner:        "admin",
			docType:      validType,
			tags:         []value_object.Tag{},
			category:     validCategory,
			variables:    []value_object.VariableDefinition{},
			accessScope:  validScope,
			wantErr:      false,
		},
		{
			name:         "empty document ID",
			id:           value_object.DocumentID(""),
			repositoryID: validRepoID,
			filePath:     validPath,
			title:        "Test",
			owner:        "admin",
			docType:      validType,
			tags:         []value_object.Tag{},
			category:     validCategory,
			variables:    []value_object.VariableDefinition{},
			accessScope:  validScope,
			wantErr:      true,
		},
		{
			name:         "empty title",
			id:           validID,
			repositoryID: validRepoID,
			filePath:     validPath,
			title:        "",
			owner:        "admin",
			docType:      validType,
			tags:         []value_object.Tag{},
			category:     validCategory,
			variables:    []value_object.VariableDefinition{},
			accessScope:  validScope,
			wantErr:      true,
		},
		{
			name:         "empty owner",
			id:           validID,
			repositoryID: validRepoID,
			filePath:     validPath,
			title:        "Test",
			owner:        "",
			docType:      validType,
			tags:         []value_object.Tag{},
			category:     validCategory,
			variables:    []value_object.VariableDefinition{},
			accessScope:  validScope,
			wantErr:      true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewDocument(
				tt.id,
				tt.repositoryID,
				tt.filePath,
				tt.title,
				tt.owner,
				tt.docType,
				tt.tags,
				tt.category,
				tt.variables,
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
	commitHash, _ := value_object.NewCommitHash("abc123def456")
	content := "# Test\n\nThis is a test document."

	err := doc.Publish(commitHash, content)
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
	commitHash2, _ := value_object.NewCommitHash("def456ghi789")
	content2 := "# Test v2\n\nUpdated content."

	err = doc.Publish(commitHash2, content2)
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
	commitHash, _ := value_object.NewCommitHash("abc123def")
	err = doc.Publish(commitHash, "content")
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

func TestDocument_UpdateMetadata(t *testing.T) {
	doc := createTestDocument(t)

	newTag, _ := value_object.NewTag("new-tag")
	newCategory, _ := value_object.NewCategory("New Category")
	newVar, _ := value_object.NewVariableDefinition("test_var", "Test Var", "desc", value_object.VariableTypeString, false, "default")

	err := doc.UpdateMetadata(
		"New Title",
		"new-owner",
		[]value_object.Tag{newTag},
		newCategory,
		[]value_object.VariableDefinition{newVar},
	)

	if err != nil {
		t.Errorf("UpdateMetadata() error = %v", err)
	}

	if doc.Title() != "New Title" {
		t.Errorf("Title() = %v, want 'New Title'", doc.Title())
	}

	if doc.Owner() != "new-owner" {
		t.Errorf("Owner() = %v, want 'new-owner'", doc.Owner())
	}

	// Test empty title
	err = doc.UpdateMetadata("", "owner", []value_object.Tag{}, newCategory, []value_object.VariableDefinition{})
	if err == nil {
		t.Error("UpdateMetadata() should return error for empty title")
	}

	// Test empty owner
	err = doc.UpdateMetadata("Title", "", []value_object.Tag{}, newCategory, []value_object.VariableDefinition{})
	if err == nil {
		t.Error("UpdateMetadata() should return error for empty owner")
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
	commitHash1, _ := value_object.NewCommitHash("abc1234567")
	err = doc.Publish(commitHash1, "content v1")
	if err != nil {
		t.Fatalf("Publish() error = %v", err)
	}

	commitHash2, _ := value_object.NewCommitHash("def4567890")
	err = doc.Publish(commitHash2, "content v2")
	if err != nil {
		t.Fatalf("Publish() second version error = %v", err)
	}

	// Verify current version is 2
	if doc.CurrentVersion() == nil {
		t.Fatal("CurrentVersion() returned nil after publishing")
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

func TestDocument_AddVersion(t *testing.T) {
	doc := createTestDocument(t)

	versionID := value_object.GenerateVersionID()
	versionNum, _ := value_object.NewVersionNumber(1)
	commitHash, _ := value_object.NewCommitHash("abc1234567")

	version, _ := NewDocumentVersion(versionID, doc.ID(), versionNum, commitHash, "content")

	err := doc.AddVersion(version)
	if err != nil {
		t.Errorf("AddVersion() error = %v", err)
	}

	if len(doc.Versions()) != 1 {
		t.Errorf("Versions() length = %d, want 1", len(doc.Versions()))
	}

	// Test adding nil version
	err = doc.AddVersion(nil)
	if err == nil {
		t.Error("AddVersion() should return error for nil version")
	}

	// Test adding version with wrong document ID
	wrongDocID := value_object.GenerateDocumentID()
	wrongVersion, _ := NewDocumentVersion(value_object.GenerateVersionID(), wrongDocID, versionNum, commitHash, "content")
	err = doc.AddVersion(wrongVersion)
	if err == nil {
		t.Error("AddVersion() should return error for version with different document ID")
	}
}

// Helper function to create a test document
func createTestDocument(t *testing.T) Document {
	id := value_object.GenerateDocumentID()
	repoID, _ := value_object.NewRepositoryID("a1b2c3d4-e5f6-4789-abcd-ef0123456789")
	filePath, _ := value_object.NewFilePath("docs/test.md")
	docType, _ := value_object.NewDocumentType("procedure")
	category, _ := value_object.NewCategory("Test")
	accessScope, _ := value_object.NewAccessScope("public")

	doc, err := NewDocument(
		id,
		repoID,
		filePath,
		"Test Document",
		"admin",
		docType,
		[]value_object.Tag{},
		category,
		[]value_object.VariableDefinition{},
		accessScope,
	)

	if err != nil {
		t.Fatalf("Failed to create test document: %v", err)
	}

	return doc
}
