package usecase

import (
	"context"
	"testing"

	"opscore/backend/internal/document/domain/entity"
	"opscore/backend/internal/document/domain/repository"
	"opscore/backend/internal/document/domain/value_object"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestVariableUseCase_GetVariableDefinitions(t *testing.T) {
	t.Run("正常に変数定義を取得できる", func(t *testing.T) {
		mockRepo := new(repository.MockDocumentRepository)

		// Create test document with variables
		doc := createTestDocument(t)
		docID := doc.ID().String()

		// Mock FindByID to return the document
		mockRepo.On("FindByID", mock.Anything, mock.AnythingOfType("value_object.DocumentID")).Return(&doc, nil)

		uc := NewVariableUseCase(mockRepo)
		result, err := uc.GetVariableDefinitions(context.Background(), docID)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		// The test document doesn't have variables by default, so expect empty list
		assert.Equal(t, 0, len(result))

		mockRepo.AssertExpectations(t)
	})

	t.Run("ドキュメントが存在しない場合はエラーを返す", func(t *testing.T) {
		mockRepo := new(repository.MockDocumentRepository)

		docID := value_object.GenerateDocumentID().String()

		// Mock FindByID to return nil (not found)
		mockRepo.On("FindByID", mock.Anything, mock.AnythingOfType("value_object.DocumentID")).Return(nil, nil)

		uc := NewVariableUseCase(mockRepo)
		result, err := uc.GetVariableDefinitions(context.Background(), docID)

		assert.Error(t, err)
		assert.Nil(t, result)

		mockRepo.AssertExpectations(t)
	})

	t.Run("無効なドキュメントIDの場合はエラーを返す", func(t *testing.T) {
		mockRepo := new(repository.MockDocumentRepository)

		uc := NewVariableUseCase(mockRepo)
		result, err := uc.GetVariableDefinitions(context.Background(), "invalid-id")

		assert.Error(t, err)
		assert.Nil(t, result)
	})
}

func TestVariableUseCase_ValidateVariableValues(t *testing.T) {
	t.Run("有効な変数値の場合はエラーを返さない", func(t *testing.T) {
		mockRepo := new(repository.MockDocumentRepository)

		// Create test document with variables
		doc := createTestDocumentWithVariables(t)
		docID := doc.ID().String()

		// Mock FindByID to return the document
		mockRepo.On("FindByID", mock.Anything, mock.AnythingOfType("value_object.DocumentID")).Return(&doc, nil)

		uc := NewVariableUseCase(mockRepo)

		values := []VariableValue{
			{Name: "server_name", Value: "prod-server"},
			{Name: "backup_path", Value: "/backup/db"},
		}

		err := uc.ValidateVariableValues(context.Background(), docID, values)

		assert.NoError(t, err)

		mockRepo.AssertExpectations(t)
	})

	t.Run("必須変数が欠けている場合はエラーを返す", func(t *testing.T) {
		mockRepo := new(repository.MockDocumentRepository)

		// Create test document with variables
		doc := createTestDocumentWithVariables(t)
		docID := doc.ID().String()

		// Mock FindByID to return the document
		mockRepo.On("FindByID", mock.Anything, mock.AnythingOfType("value_object.DocumentID")).Return(&doc, nil)

		uc := NewVariableUseCase(mockRepo)

		// Missing required variable 'server_name'
		values := []VariableValue{
			{Name: "backup_path", Value: "/backup/db"},
		}

		err := uc.ValidateVariableValues(context.Background(), docID, values)

		assert.Error(t, err)

		mockRepo.AssertExpectations(t)
	})
}

func TestVariableUseCase_SubstituteVariables(t *testing.T) {
	t.Run("変数を正しく置換できる", func(t *testing.T) {
		mockRepo := new(repository.MockDocumentRepository)

		uc := NewVariableUseCase(mockRepo)

		content := "Connect to {{server_name}} and backup to {{backup_path}}"
		values := []VariableValue{
			{Name: "server_name", Value: "prod-server"},
			{Name: "backup_path", Value: "/backup/db"},
		}

		result, err := uc.SubstituteVariables(context.Background(), content, values)

		assert.NoError(t, err)
		assert.Equal(t, "Connect to prod-server and backup to /backup/db", result)
	})

	t.Run("複数箇所の同じ変数を置換できる", func(t *testing.T) {
		mockRepo := new(repository.MockDocumentRepository)

		uc := NewVariableUseCase(mockRepo)

		content := "Server: {{server_name}}, Connect to {{server_name}}"
		values := []VariableValue{
			{Name: "server_name", Value: "prod-server"},
		}

		result, err := uc.SubstituteVariables(context.Background(), content, values)

		assert.NoError(t, err)
		assert.Equal(t, "Server: prod-server, Connect to prod-server", result)
	})

	t.Run("変数が存在しない場合はそのまま残す", func(t *testing.T) {
		mockRepo := new(repository.MockDocumentRepository)

		uc := NewVariableUseCase(mockRepo)

		content := "Connect to {{server_name}} and backup to {{backup_path}}"
		values := []VariableValue{
			{Name: "server_name", Value: "prod-server"},
		}

		result, err := uc.SubstituteVariables(context.Background(), content, values)

		assert.NoError(t, err)
		assert.Equal(t, "Connect to prod-server and backup to {{backup_path}}", result)
	})

	t.Run("変数値にnilを指定した場合は空文字列に置換される", func(t *testing.T) {
		mockRepo := new(repository.MockDocumentRepository)

		uc := NewVariableUseCase(mockRepo)

		content := "Server: {{server_name}}, Path: {{backup_path}}"
		values := []VariableValue{
			{Name: "server_name", Value: nil},
			{Name: "backup_path", Value: "/backup"},
		}

		result, err := uc.SubstituteVariables(context.Background(), content, values)

		assert.NoError(t, err)
		assert.Equal(t, "Server: , Path: /backup", result)
	})

	t.Run("変数値に特殊文字が含まれていても正しく置換される", func(t *testing.T) {
		mockRepo := new(repository.MockDocumentRepository)

		uc := NewVariableUseCase(mockRepo)

		content := "Path: {{path}}, Pattern: {{pattern}}"
		values := []VariableValue{
			{Name: "path", Value: "C:\\Program Files\\App"},
			{Name: "pattern", Value: "*.txt"},
		}

		result, err := uc.SubstituteVariables(context.Background(), content, values)

		assert.NoError(t, err)
		assert.Equal(t, "Path: C:\\Program Files\\App, Pattern: *.txt", result)
	})
}

// Helper function to create a test document with variables
func createTestDocumentWithVariables(t *testing.T) entity.Document {
	docID := value_object.GenerateDocumentID()
	repoID, _ := value_object.NewRepositoryID("a1b2c3d4-e5f6-7890-1234-567890abcdef")
	accessScope, _ := value_object.NewAccessScope("public")

	doc, err := entity.NewDocument(docID, repoID, "test-owner", accessScope)
	assert.NoError(t, err)

	// Create variables
	serverNameVar, _ := value_object.NewVariableDefinition(
		"server_name",
		"Server Name",
		"The target server name",
		value_object.VariableTypeString,
		true,
		"localhost",
	)

	backupPathVar, _ := value_object.NewVariableDefinition(
		"backup_path",
		"Backup Path",
		"Path to backup directory",
		value_object.VariableTypeString,
		true,
		"/backup",
	)

	variables := []value_object.VariableDefinition{serverNameVar, backupPathVar}

	// Publish a version with variables
	filePath, _ := value_object.NewFilePath("docs/test.md")
	commitHash, _ := value_object.NewCommitHash("abc1234567890")
	source, _ := value_object.NewDocumentSource(filePath, commitHash)
	docType, _ := value_object.NewDocumentType("procedure")

	err = doc.Publish(source, "Test Document", docType, nil, variables, "# Test Content")
	assert.NoError(t, err)

	return doc
}
