package repository

import (
	"context"
	"errors"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/suite"

	"opscore/backend/internal/authentication/domain/entity"
	domain_err "opscore/backend/internal/authentication/domain/err"
	"opscore/backend/internal/authentication/domain/value_object"
	"opscore/backend/internal/shared/testutil"
)

type UserRepositoryTestSuite struct {
	suite.Suite
	db   *pgxpool.Pool
	repo *PostgresUserRepository
	ctx  context.Context
}

func TestPostgresUserRepository(t *testing.T) {
	testDB := testutil.SetupTestDatabase(t, 1)
	s := new(UserRepositoryTestSuite)
	s.db = testDB.Pool
	suite.Run(t, s)
}

func (s *UserRepositoryTestSuite) SetupSuite() {
	s.ctx = context.Background()
	s.repo = NewPostgresUserRepository(s.db).(*PostgresUserRepository)
}

func (s *UserRepositoryTestSuite) SetupTest() {
	s.Require().NotNil(s.db, "database pool must be initialized")
	_, _ = s.db.Exec(s.ctx, "DELETE FROM refresh_tokens")
	_, _ = s.db.Exec(s.ctx, "DELETE FROM user_identities")
	_, _ = s.db.Exec(s.ctx, "DELETE FROM users")
}

// createTestUser creates a User aggregate with one identity and saves it.
func (s *UserRepositoryTestSuite) createTestUser(email, displayName, provider, providerUserID string) *entity.User {
	userID := value_object.NewUserID()
	user := entity.NewUser(userID, email, displayName, "https://example.com/pic.jpg")

	identityID := value_object.NewIdentityID()
	identity, err := entity.NewIdentity(identityID, userID, provider, providerUserID, email, displayName, "https://example.com/pic.jpg")
	s.Require().NoError(err)
	s.Require().NoError(user.AddIdentity(identity))

	err = s.repo.Save(s.ctx, user)
	s.Require().NoError(err)

	return user
}

func (s *UserRepositoryTestSuite) TestSave_NewUserWithIdentity() {
	user := s.createTestUser("alice@example.com", "Alice", "google", "google-123")

	// Verify user was persisted
	found, err := s.repo.FindByID(s.ctx, user.ID())
	s.Require().NoError(err)
	s.Equal(user.ID(), found.ID())
	s.Equal("alice@example.com", found.Email())
	s.Equal("Alice", found.DisplayName())
}

func (s *UserRepositoryTestSuite) TestSave_UserWithNoPictureURL() {
	userID := value_object.NewUserID()
	user := entity.NewUser(userID, "nopic@example.com", "NoPic", "")

	identityID := value_object.NewIdentityID()
	identity, err := entity.NewIdentity(identityID, userID, "github", "gh-nopic", "nopic@example.com", "NoPic", "")
	s.Require().NoError(err)
	s.Require().NoError(user.AddIdentity(identity))

	err = s.repo.Save(s.ctx, user)
	s.Require().NoError(err)

	found, err := s.repo.FindByID(s.ctx, user.ID())
	s.Require().NoError(err)
	s.Equal("", found.PictureURL())
}

func (s *UserRepositoryTestSuite) TestFindByID_NotFound() {
	nonExistentID := value_object.NewUserID()
	_, err := s.repo.FindByID(s.ctx, nonExistentID)
	s.Require().Error(err)

	var domainErr *domain_err.DomainError
	s.Require().True(s.errorAs(err, &domainErr))
	s.Equal(domain_err.ErrNotFoundUser, domainErr.Kind())
}

func (s *UserRepositoryTestSuite) TestFindByIDWithIdentities() {
	user := s.createTestUser("bob@example.com", "Bob", "github", "gh-bob")

	found, err := s.repo.FindByIDWithIdentities(s.ctx, user.ID())
	s.Require().NoError(err)
	s.True(found.HasIdentitiesLoaded())
	s.Len(found.Identities(), 1)

	identity := found.Identities()[0]
	s.Equal("github", identity.Provider())
	s.Equal("gh-bob", identity.ProviderUserID())
	s.Equal("bob@example.com", identity.Email())
	s.True(identity.IsPrimary()) // First identity is auto-primary
}

func (s *UserRepositoryTestSuite) TestFindByProviderIdentity_Found() {
	user := s.createTestUser("charlie@example.com", "Charlie", "gitlab", "gl-charlie")

	found, err := s.repo.FindByProviderIdentity(s.ctx, "gitlab", "gl-charlie")
	s.Require().NoError(err)
	s.Equal(user.ID(), found.ID())
	s.True(found.HasIdentitiesLoaded())
	s.Len(found.Identities(), 1)
}

func (s *UserRepositoryTestSuite) TestFindByProviderIdentity_NotFound() {
	_, err := s.repo.FindByProviderIdentity(s.ctx, "unknown", "unknown-id")
	s.Require().Error(err)

	var domainErr *domain_err.DomainError
	s.Require().True(s.errorAs(err, &domainErr))
	s.Equal(domain_err.ErrNotFoundUser, domainErr.Kind())
}

func (s *UserRepositoryTestSuite) TestSave_AddSecondIdentity() {
	user := s.createTestUser("dave@example.com", "Dave", "google", "g-dave")

	// Reload with identities
	user, err := s.repo.FindByIDWithIdentities(s.ctx, user.ID())
	s.Require().NoError(err)

	// Add a second identity
	newIdentityID := value_object.NewIdentityID()
	newIdentity, err := entity.NewIdentity(newIdentityID, user.ID(), "github", "gh-dave", "dave@github.com", "Dave GH", "")
	s.Require().NoError(err)
	s.Require().NoError(user.AddIdentity(newIdentity))

	err = s.repo.Save(s.ctx, user)
	s.Require().NoError(err)

	// Verify both identities exist
	found, err := s.repo.FindByIDWithIdentities(s.ctx, user.ID())
	s.Require().NoError(err)
	s.Len(found.Identities(), 2)
}

func (s *UserRepositoryTestSuite) TestSave_RemoveIdentity() {
	user := s.createTestUser("eve@example.com", "Eve", "google", "g-eve")

	// Reload and add second identity
	user, err := s.repo.FindByIDWithIdentities(s.ctx, user.ID())
	s.Require().NoError(err)

	secondID := value_object.NewIdentityID()
	secondIdentity, err := entity.NewIdentity(secondID, user.ID(), "github", "gh-eve", "eve@github.com", "Eve GH", "")
	s.Require().NoError(err)
	s.Require().NoError(user.AddIdentity(secondIdentity))
	err = s.repo.Save(s.ctx, user)
	s.Require().NoError(err)

	// Reload and remove the non-primary identity
	user, err = s.repo.FindByIDWithIdentities(s.ctx, user.ID())
	s.Require().NoError(err)
	s.Len(user.Identities(), 2)

	err = user.RemoveIdentity(secondID)
	s.Require().NoError(err)

	err = s.repo.Save(s.ctx, user)
	s.Require().NoError(err)

	// Verify only one identity remains
	found, err := s.repo.FindByIDWithIdentities(s.ctx, user.ID())
	s.Require().NoError(err)
	s.Len(found.Identities(), 1)
}

func (s *UserRepositoryTestSuite) TestSave_ChangePrimaryIdentity() {
	user := s.createTestUser("frank@example.com", "Frank", "google", "g-frank")

	// Reload and add second identity
	user, err := s.repo.FindByIDWithIdentities(s.ctx, user.ID())
	s.Require().NoError(err)

	secondID := value_object.NewIdentityID()
	secondIdentity, err := entity.NewIdentity(secondID, user.ID(), "github", "gh-frank", "frank@github.com", "Frank GH", "")
	s.Require().NoError(err)
	s.Require().NoError(user.AddIdentity(secondIdentity))
	err = s.repo.Save(s.ctx, user)
	s.Require().NoError(err)

	// Reload and change primary
	user, err = s.repo.FindByIDWithIdentities(s.ctx, user.ID())
	s.Require().NoError(err)
	err = user.SetPrimaryIdentity(secondID)
	s.Require().NoError(err)

	err = s.repo.Save(s.ctx, user)
	s.Require().NoError(err)

	// Verify primary changed
	found, err := s.repo.FindByIDWithIdentities(s.ctx, user.ID())
	s.Require().NoError(err)
	primary := found.GetPrimaryIdentity()
	s.Require().NotNil(primary)
	s.Equal(secondID, primary.ID())
}

func (s *UserRepositoryTestSuite) TestSave_UpdateLastLogin() {
	user := s.createTestUser("grace@example.com", "Grace", "google", "g-grace")

	// Reload
	user, err := s.repo.FindByIDWithIdentities(s.ctx, user.ID())
	s.Require().NoError(err)

	originalLoginAt := user.LastLoginAt()
	user.UpdateLastLogin()

	err = s.repo.Save(s.ctx, user)
	s.Require().NoError(err)

	found, err := s.repo.FindByID(s.ctx, user.ID())
	s.Require().NoError(err)
	s.True(found.LastLoginAt().After(originalLoginAt) || found.LastLoginAt().Equal(originalLoginAt))
}

func (s *UserRepositoryTestSuite) TestDelete_Success() {
	user := s.createTestUser("henry@example.com", "Henry", "google", "g-henry")

	err := s.repo.Delete(s.ctx, user.ID())
	s.Require().NoError(err)

	// User should not exist anymore
	_, err = s.repo.FindByID(s.ctx, user.ID())
	s.Require().Error(err)

	var domainErr *domain_err.DomainError
	s.Require().True(s.errorAs(err, &domainErr))
	s.Equal(domain_err.ErrNotFoundUser, domainErr.Kind())
}

func (s *UserRepositoryTestSuite) TestDelete_CascadesIdentities() {
	user := s.createTestUser("iris@example.com", "Iris", "github", "gh-iris")

	err := s.repo.Delete(s.ctx, user.ID())
	s.Require().NoError(err)

	// Identity should also be gone (ON DELETE CASCADE)
	var count int
	err = s.db.QueryRow(s.ctx, "SELECT COUNT(*) FROM user_identities WHERE user_id = $1", user.ID().String()).Scan(&count)
	s.Require().NoError(err)
	s.Equal(0, count)
}

func (s *UserRepositoryTestSuite) TestDelete_NotFound() {
	nonExistentID := value_object.NewUserID()
	err := s.repo.Delete(s.ctx, nonExistentID)
	s.Require().Error(err)

	var domainErr *domain_err.DomainError
	s.Require().True(s.errorAs(err, &domainErr))
	s.Equal(domain_err.ErrNotFoundUser, domainErr.Kind())
}

// errorAs is a helper that wraps errors.As for DomainError.
func (s *UserRepositoryTestSuite) errorAs(err error, target interface{}) bool {
	return errorAs(err, target)
}

// errorAs is a package-level helper for errors.As with DomainError.
func errorAs(err error, target interface{}) bool {
	return errors.As(err, target.(**domain_err.DomainError))
}
