package dto

import (
	"testing"

	"opscore/backend/internal/document/domain/entity"
	"opscore/backend/internal/document/domain/value_object"
)

func TestToDocumentResponse_IncludesOrigin(t *testing.T) {
	repoID, _ := value_object.NewRepositoryID("a1b2c3d4-e5f6-4789-abcd-ef0123456789")
	origin, _ := value_object.NewRepositoryOrigin("12345", "org", "repo")
	accessScope, _ := value_object.NewAccessScope("public")
	docID := value_object.GenerateDocumentID()

	doc, err := entity.NewDocument(docID, repoID, &origin, accessScope)
	if err != nil {
		t.Fatalf("failed to create document: %v", err)
	}

	resp := ToDocumentResponse(doc)
	if resp.ProviderRepositoryID != "12345" {
		t.Fatalf("expected provider_repository_id to be 12345, got %s", resp.ProviderRepositoryID)
	}
	if resp.Owner != "org" {
		t.Fatalf("expected owner to be org, got %s", resp.Owner)
	}
	if resp.Repository != "repo" {
		t.Fatalf("expected repository to be repo, got %s", resp.Repository)
	}
}

func TestToDocumentListItemResponse_IncludesOrigin(t *testing.T) {
	repoID, _ := value_object.NewRepositoryID("a1b2c3d4-e5f6-4789-abcd-ef0123456789")
	origin, _ := value_object.NewRepositoryOrigin("12345", "org", "repo")
	accessScope, _ := value_object.NewAccessScope("public")
	docID := value_object.GenerateDocumentID()

	doc, err := entity.NewDocument(docID, repoID, &origin, accessScope)
	if err != nil {
		t.Fatalf("failed to create document: %v", err)
	}

	listItem := ToDocumentListItemResponse(doc)
	if listItem.ProviderRepositoryID != "12345" {
		t.Fatalf("expected provider_repository_id to be 12345, got %s", listItem.ProviderRepositoryID)
	}
	if listItem.Owner != "org" {
		t.Fatalf("expected owner to be org, got %s", listItem.Owner)
	}
	if listItem.Repository != "repo" {
		t.Fatalf("expected repository to be repo, got %s", listItem.Repository)
	}
}
