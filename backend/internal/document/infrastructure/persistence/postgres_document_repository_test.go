package persistence

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"opscore/backend/internal/document/domain/entity"
	"opscore/backend/internal/document/domain/repository"
	"opscore/backend/internal/document/domain/value_object"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/suite"
)

type DocumentRepositoryTestSuite struct {
	suite.Suite
	db   *pgxpool.Pool
	repo repository.DocumentRepository
	ctx  context.Context
}

func TestPostgresDocumentRepository(t *testing.T) {
	dbURL := os.Getenv("TEST_DATABASE_URL")
	if dbURL == "" {
		t.Skip("TEST_DATABASE_URL is not set; skipping document repository tests")
	}

	suite.Run(t, new(DocumentRepositoryTestSuite))
}

func (s *DocumentRepositoryTestSuite) SetupSuite() {
	dbURL := os.Getenv("TEST_DATABASE_URL")
	var err error
	s.ctx = context.Background()
	s.db, err = pgxpool.New(s.ctx, dbURL)
	if err != nil {
		s.T().Fatalf("failed to connect to test database: %v", err)
	}

	s.repo = NewPostgresDocumentRepository(s.db)
}

func (s *DocumentRepositoryTestSuite) TearDownSuite() {
	if s.db != nil {
		s.db.Close()
	}
}

func (s *DocumentRepositoryTestSuite) SetupTest() {
	s.cleanupTables()
}

func (s *DocumentRepositoryTestSuite) cleanupTables() {
	_, _ = s.db.Exec(s.ctx, "DELETE FROM document_versions")
	_, _ = s.db.Exec(s.ctx, "DELETE FROM documents")
	_, _ = s.db.Exec(s.ctx, "DELETE FROM repositories")
}

func (s *DocumentRepositoryTestSuite) insertRepository() value_object.RepositoryID {
	id := uuid.New()
	name := fmt.Sprintf("repo-%s", id.String()[:8])
	url := fmt.Sprintf("https://example.com/%s.git", name)

	_, err := s.db.Exec(s.ctx, "INSERT INTO repositories (id, name, url, created_at, updated_at) VALUES ($1, $2, $3, $4, $5)", id, name, url, time.Now(), time.Now())
	s.Require().NoError(err)

	repoID, err := value_object.NewRepositoryID(id.String())
	s.Require().NoError(err)
	return repoID
}

func (s *DocumentRepositoryTestSuite) buildDocument(repoID value_object.RepositoryID, providerRepoID, owner, repo, title string) entity.Document {
	origin, err := value_object.NewRepositoryOrigin(providerRepoID, owner, repo)
	s.Require().NoError(err)

	docID := value_object.GenerateDocumentID()
	doc, err := entity.NewDocument(docID, repoID, &origin, value_object.AccessScopePublic)
	s.Require().NoError(err)

	filePath, err := value_object.NewFilePath("docs/sample.md")
	s.Require().NoError(err)
	commit, err := value_object.NewCommitHash("abcdef1")
	s.Require().NoError(err)
	source, err := value_object.NewDocumentSource(filePath, commit)
	s.Require().NoError(err)

	tag, err := value_object.NewTag("ops")
	s.Require().NoError(err)

	err = doc.Publish(source, title, value_object.DocumentTypeProcedure, []value_object.Tag{tag}, nil, "# content")
	s.Require().NoError(err)

	return doc
}

func (s *DocumentRepositoryTestSuite) TestSaveAndFindByID() {
	repoID := s.insertRepository()
	doc := s.buildDocument(repoID, "provider-1", "alice", "repo-one", "Doc One")

	err := s.repo.Save(s.ctx, doc)
	s.Require().NoError(err)

	found, err := s.repo.FindByID(s.ctx, doc.ID())
	s.Require().NoError(err)
	s.Require().NotNil(found)
	s.Equal(doc.ID(), found.ID())
	s.Equal(doc.Origin().ProviderRepositoryID(), found.Origin().ProviderRepositoryID())
	s.Equal(doc.Origin().Owner(), found.Origin().Owner())
	s.Equal(doc.Origin().Repository(), found.Origin().Repository())
	s.Equal(doc.CurrentVersion().VersionNumber().Int(), found.CurrentVersion().VersionNumber().Int())
}

func (s *DocumentRepositoryTestSuite) TestUpdateAddsNewVersion() {
	repoID := s.insertRepository()
	doc := s.buildDocument(repoID, "provider-2", "bob", "repo-two", "Doc Two")

	err := s.repo.Save(s.ctx, doc)
	s.Require().NoError(err)

	// Publish a new version
	filePath, _ := value_object.NewFilePath("docs/sample.md")
	commit, _ := value_object.NewCommitHash("abcdef2")
	source, _ := value_object.NewDocumentSource(filePath, commit)
	tag, _ := value_object.NewTag("ops")
	err = doc.Publish(source, "Doc Two v2", value_object.DocumentTypeProcedure, []value_object.Tag{tag}, nil, "# updated")
	s.Require().NoError(err)

	err = s.repo.Update(s.ctx, doc)
	s.Require().NoError(err)

	found, err := s.repo.FindByID(s.ctx, doc.ID())
	s.Require().NoError(err)
	s.Require().NotNil(found)
	s.Len(found.Versions(), 2)
	s.Equal(2, found.CurrentVersion().VersionNumber().Int())
}

func (s *DocumentRepositoryTestSuite) TestFindPublishedFilters() {
	repoID := s.insertRepository()
	docA := s.buildDocument(repoID, "provider-A", "carol", "repo-a", "Doc A")
	docB := s.buildDocument(repoID, "provider-B", "dave", "repo-b", "Doc B")

	s.Require().NoError(s.repo.Save(s.ctx, docA))
	s.Require().NoError(s.repo.Save(s.ctx, docB))

	filtered, err := s.repo.FindPublished(s.ctx, repository.NewProviderRepositoryIDFilter("provider-A"))
	s.Require().NoError(err)
	s.Len(filtered, 1)
	s.Equal(docA.ID(), filtered[0].ID())

	filtered, err = s.repo.FindPublished(s.ctx, repository.NewOwnerFilter("dave"))
	s.Require().NoError(err)
	s.Len(filtered, 1)
	s.Equal(docB.ID(), filtered[0].ID())

	filtered, err = s.repo.FindPublished(s.ctx, repository.NewRepositoryNameFilter("repo-a"))
	s.Require().NoError(err)
	s.Len(filtered, 1)
	s.Equal(docA.ID(), filtered[0].ID())
}
