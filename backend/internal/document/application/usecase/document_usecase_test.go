package usecase

import (
	"context"
	"errors"
	"testing"

	"opscore/backend/internal/document/application/dto"
	apperror "opscore/backend/internal/document/application/error"
	"opscore/backend/internal/document/domain/entity"
	"opscore/backend/internal/document/domain/repository"
	"opscore/backend/internal/document/domain/value_object"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Helper function to create a test document
func createTestDocument(t *testing.T) entity.Document {
	docID := value_object.GenerateDocumentID()
	repoID, _ := value_object.NewRepositoryID("a1b2c3d4-e5f6-7890-1234-567890abcdef")
	accessScope, _ := value_object.NewAccessScope("public")

	doc, err := entity.NewDocument(docID, repoID, "test-owner", accessScope)
	assert.NoError(t, err)

	// Publish a version
	filePath, _ := value_object.NewFilePath("docs/test.md")
	commitHash, _ := value_object.NewCommitHash("abc1234567890")
	source, _ := value_object.NewDocumentSource(filePath, commitHash)
	docType, _ := value_object.NewDocumentType("procedure")

	err = doc.Publish(source, "Test Document", docType, nil, nil, "# Test Content")
	assert.NoError(t, err)

	return doc
}

func TestDocumentUseCase_CreateDocument(t *testing.T) {
	t.Run("正常にドキュメントを作成できる", func(t *testing.T) {
		mockRepo := new(repository.MockDocumentRepository)

		req := &dto.CreateDocumentRequest{
			RepositoryID: "a1b2c3d4-e5f6-7890-1234-567890abcdef",
			FilePath:     "docs/test.md",
			CommitHash:   "abc1234567890",
			Title:        "Test Document",
			DocType:      "procedure",
			Owner:        "test-owner",
			Tags:         []string{"test", "document"},
			Variables: []dto.VariableDefinitionDTO{
				{
					Name:         "server_name",
					Label:        "Server Name",
					Description:  "The target server",
					Type:         "string",
					Required:     true,
					DefaultValue: "localhost",
				},
			},
			Content:      "# Test Document\n\nThis is a test.",
			AccessScope:  "public",
			IsAutoUpdate: true,
		}

		// Mock Save to succeed
		mockRepo.On("Save", mock.Anything, mock.AnythingOfType("*entity.document")).Return(nil)

		uc := NewDocumentUseCase(mockRepo)
		result, err := uc.CreateDocument(context.Background(), req)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.NotEmpty(t, result.ID)
		assert.Equal(t, req.RepositoryID, result.RepositoryID)
		assert.Equal(t, req.Owner, result.Owner)
		assert.True(t, result.IsPublished)
		assert.True(t, result.IsAutoUpdate)
		assert.NotNil(t, result.CurrentVersion)
		assert.Equal(t, req.Title, result.CurrentVersion.Title)

		mockRepo.AssertExpectations(t)
	})

	t.Run("無効なリポジトリIDでエラーになる", func(t *testing.T) {
		mockRepo := new(repository.MockDocumentRepository)

		req := &dto.CreateDocumentRequest{
			RepositoryID: "invalid-uuid",
			FilePath:     "docs/test.md",
			CommitHash:   "abc1234567890",
			Title:        "Test Document",
			DocType:      "procedure",
			Owner:        "test-owner",
			Content:      "# Test",
			AccessScope:  "public",
		}

		uc := NewDocumentUseCase(mockRepo)
		result, err := uc.CreateDocument(context.Background(), req)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.True(t, errors.Is(err, apperror.ErrBadRequest))
	})

	t.Run("空のタイトルでエラーになる", func(t *testing.T) {
		mockRepo := new(repository.MockDocumentRepository)

		req := &dto.CreateDocumentRequest{
			RepositoryID: "a1b2c3d4-e5f6-7890-1234-567890abcdef",
			FilePath:     "docs/test.md",
			CommitHash:   "abc1234567890",
			Title:        "",
			DocType:      "procedure",
			Owner:        "test-owner",
			Content:      "# Test",
			AccessScope:  "public",
		}

		uc := NewDocumentUseCase(mockRepo)
		result, err := uc.CreateDocument(context.Background(), req)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.True(t, errors.Is(err, apperror.ErrBadRequest))
	})
}

func TestDocumentUseCase_GetDocument(t *testing.T) {
	t.Run("存在するドキュメントを取得できる", func(t *testing.T) {
		mockRepo := new(repository.MockDocumentRepository)
		testDoc := createTestDocument(t)

		docID := testDoc.ID()
		mockRepo.On("FindByID", mock.Anything, docID).Return(testDoc, nil)

		uc := NewDocumentUseCase(mockRepo)
		result, err := uc.GetDocument(context.Background(), docID.String())

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, docID.String(), result.ID)

		mockRepo.AssertExpectations(t)
	})

	t.Run("存在しないドキュメントでエラーになる", func(t *testing.T) {
		mockRepo := new(repository.MockDocumentRepository)

		docID, _ := value_object.NewDocumentID("a1b2c3d4-e5f6-7890-1234-567890abcdef")
		mockRepo.On("FindByID", mock.Anything, docID).Return(nil, nil)

		uc := NewDocumentUseCase(mockRepo)
		result, err := uc.GetDocument(context.Background(), docID.String())

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.True(t, errors.Is(err, apperror.ErrNotFound))

		mockRepo.AssertExpectations(t)
	})

	t.Run("無効なドキュメントIDでエラーになる", func(t *testing.T) {
		mockRepo := new(repository.MockDocumentRepository)

		uc := NewDocumentUseCase(mockRepo)
		result, err := uc.GetDocument(context.Background(), "invalid-uuid")

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.True(t, errors.Is(err, apperror.ErrBadRequest))
	})
}

func TestDocumentUseCase_ListDocuments(t *testing.T) {
	t.Run("ドキュメント一覧を取得できる", func(t *testing.T) {
		mockRepo := new(repository.MockDocumentRepository)
		testDoc := createTestDocument(t)
		docs := []entity.Document{testDoc}

		mockRepo.On("FindPublished", mock.Anything, mock.Anything).Return(docs, nil)

		uc := NewDocumentUseCase(mockRepo)
		result, err := uc.ListDocuments(context.Background())

		assert.NoError(t, err)
		assert.Len(t, result, 1)

		mockRepo.AssertExpectations(t)
	})

	t.Run("ドキュメントが存在しない場合は空のリストを返す", func(t *testing.T) {
		mockRepo := new(repository.MockDocumentRepository)
		var emptyDocs []entity.Document

		mockRepo.On("FindPublished", mock.Anything, mock.Anything).Return(emptyDocs, nil)

		uc := NewDocumentUseCase(mockRepo)
		result, err := uc.ListDocuments(context.Background())

		assert.NoError(t, err)
		assert.Empty(t, result)

		mockRepo.AssertExpectations(t)
	})
}

func TestDocumentUseCase_UpdateDocument(t *testing.T) {
	t.Run("正常にドキュメントを更新できる", func(t *testing.T) {
		mockRepo := new(repository.MockDocumentRepository)
		testDoc := createTestDocument(t)

		docID := testDoc.ID()
		mockRepo.On("FindByID", mock.Anything, docID).Return(testDoc, nil)
		mockRepo.On("Update", mock.Anything, mock.AnythingOfType("*entity.document")).Return(nil)

		req := &dto.UpdateDocumentRequest{
			FilePath:   "docs/updated.md",
			CommitHash: "def4567890123",
			Title:      "Updated Document",
			DocType:    "procedure",
			Tags:       []string{"updated"},
			Content:    "# Updated Content",
		}

		uc := NewDocumentUseCase(mockRepo)
		result, err := uc.UpdateDocument(context.Background(), docID.String(), req)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		// The version count should increase
		assert.GreaterOrEqual(t, result.VersionCount, 2)

		mockRepo.AssertExpectations(t)
	})

	t.Run("存在しないドキュメントの更新でエラーになる", func(t *testing.T) {
		mockRepo := new(repository.MockDocumentRepository)

		docID, _ := value_object.NewDocumentID("a1b2c3d4-e5f6-7890-1234-567890abcdef")
		mockRepo.On("FindByID", mock.Anything, docID).Return(nil, nil)

		req := &dto.UpdateDocumentRequest{
			FilePath:   "docs/test.md",
			CommitHash: "abc1234567890",
			Title:      "Test",
			DocType:    "procedure",
			Content:    "# Test",
		}

		uc := NewDocumentUseCase(mockRepo)
		result, err := uc.UpdateDocument(context.Background(), docID.String(), req)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.True(t, errors.Is(err, apperror.ErrNotFound))

		mockRepo.AssertExpectations(t)
	})
}

func TestDocumentUseCase_GetDocumentVersions(t *testing.T) {
	t.Run("バージョン履歴を取得できる", func(t *testing.T) {
		mockRepo := new(repository.MockDocumentRepository)
		testDoc := createTestDocument(t)

		docID := testDoc.ID()
		versions := testDoc.Versions()

		mockRepo.On("FindByID", mock.Anything, docID).Return(testDoc, nil)
		mockRepo.On("FindVersionsByDocumentID", mock.Anything, docID).Return(versions, nil)

		uc := NewDocumentUseCase(mockRepo)
		result, err := uc.GetDocumentVersions(context.Background(), docID.String())

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, docID.String(), result.DocumentID)
		assert.Len(t, result.Versions, len(versions))

		mockRepo.AssertExpectations(t)
	})

	t.Run("存在しないドキュメントでエラーになる", func(t *testing.T) {
		mockRepo := new(repository.MockDocumentRepository)

		docID, _ := value_object.NewDocumentID("a1b2c3d4-e5f6-7890-1234-567890abcdef")
		mockRepo.On("FindByID", mock.Anything, docID).Return(nil, nil)

		uc := NewDocumentUseCase(mockRepo)
		result, err := uc.GetDocumentVersions(context.Background(), docID.String())

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.True(t, errors.Is(err, apperror.ErrNotFound))

		mockRepo.AssertExpectations(t)
	})
}

func TestDocumentUseCase_RollbackDocumentVersion(t *testing.T) {
	t.Run("正常にロールバックできる", func(t *testing.T) {
		mockRepo := new(repository.MockDocumentRepository)

		// Create a document with multiple versions
		docID := value_object.GenerateDocumentID()
		repoID, _ := value_object.NewRepositoryID("a1b2c3d4-e5f6-7890-1234-567890abcdef")
		accessScope, _ := value_object.NewAccessScope("public")

		doc, _ := entity.NewDocument(docID, repoID, "test-owner", accessScope)

		filePath1, _ := value_object.NewFilePath("docs/v1.md")
		commitHash1, _ := value_object.NewCommitHash("abc1234567890")
		source1, _ := value_object.NewDocumentSource(filePath1, commitHash1)
		docType, _ := value_object.NewDocumentType("procedure")
		doc.Publish(source1, "Version 1", docType, nil, nil, "# Version 1")

		filePath2, _ := value_object.NewFilePath("docs/v2.md")
		commitHash2, _ := value_object.NewCommitHash("def4567890123")
		source2, _ := value_object.NewDocumentSource(filePath2, commitHash2)
		doc.Publish(source2, "Version 2", docType, nil, nil, "# Version 2")

		mockRepo.On("FindByID", mock.Anything, docID).Return(doc, nil)
		mockRepo.On("Update", mock.Anything, mock.AnythingOfType("*entity.document")).Return(nil)

		uc := NewDocumentUseCase(mockRepo)
		result, err := uc.RollbackDocumentVersion(context.Background(), docID.String(), 1)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, 1, result.CurrentVersion.VersionNumber)

		mockRepo.AssertExpectations(t)
	})
}

func TestDocumentUseCase_UpdateDocumentMetadata(t *testing.T) {
	t.Run("正常にメタデータを更新できる", func(t *testing.T) {
		mockRepo := new(repository.MockDocumentRepository)
		testDoc := createTestDocument(t)

		docID := testDoc.ID()
		mockRepo.On("FindByID", mock.Anything, docID).Return(testDoc, nil)
		mockRepo.On("Update", mock.Anything, mock.AnythingOfType("*entity.document")).Return(nil)

		newScope := "private"
		isAutoUpdate := false
		req := &dto.UpdateDocumentMetadataRequest{
			AccessScope:  &newScope,
			IsAutoUpdate: &isAutoUpdate,
		}

		uc := NewDocumentUseCase(mockRepo)
		result, err := uc.UpdateDocumentMetadata(context.Background(), docID.String(), req)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, "private", result.AccessScope)
		assert.False(t, result.IsAutoUpdate)

		mockRepo.AssertExpectations(t)
	})
}
