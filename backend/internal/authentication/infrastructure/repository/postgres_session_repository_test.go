package repository

import (
	"context"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/suite"

	"opscore/backend/internal/authentication/domain/entity"
	domain_err "opscore/backend/internal/authentication/domain/err"
	"opscore/backend/internal/authentication/domain/value_object"
	"opscore/backend/internal/shared/testutil"
)

type SessionRepositoryTestSuite struct {
	suite.Suite
	db   *pgxpool.Pool
	repo *PostgresSessionRepository
	ctx  context.Context
}

func TestPostgresSessionRepository(t *testing.T) {
	testDB := testutil.SetupTestDatabase(t, 1)
	s := new(SessionRepositoryTestSuite)
	s.db = testDB.Pool
	suite.Run(t, s)
}

func (s *SessionRepositoryTestSuite) SetupSuite() {
	s.ctx = context.Background()
	s.repo = NewPostgresSessionRepository(s.db).(*PostgresSessionRepository)
}

func (s *SessionRepositoryTestSuite) SetupTest() {
	s.Require().NotNil(s.db, "database pool must be initialized")
	_, _ = s.db.Exec(s.ctx, "DELETE FROM refresh_tokens")
	_, _ = s.db.Exec(s.ctx, "DELETE FROM user_identities")
	_, _ = s.db.Exec(s.ctx, "DELETE FROM users")
}

// insertTestUser inserts a user directly for FK references.
func (s *SessionRepositoryTestSuite) insertTestUser() value_object.UserID {
	userID := value_object.NewUserID()
	_, err := s.db.Exec(s.ctx,
		"INSERT INTO users (id, email, display_name, created_at, updated_at, last_login_at) VALUES ($1, $2, $3, NOW(), NOW(), NOW())",
		userID.String(), "test@example.com", "Test User",
	)
	s.Require().NoError(err)
	return userID
}

func (s *SessionRepositoryTestSuite) TestCreate_Success() {
	userID := s.insertTestUser()
	session, _, err := entity.NewSession(userID, 7*24*time.Hour, false)
	s.Require().NoError(err)

	err = s.repo.Create(s.ctx, session)
	s.Require().NoError(err)
}

func (s *SessionRepositoryTestSuite) TestCreate_DuplicateID() {
	userID := s.insertTestUser()
	session, _, err := entity.NewSession(userID, 7*24*time.Hour, false)
	s.Require().NoError(err)

	err = s.repo.Create(s.ctx, session)
	s.Require().NoError(err)

	// Attempt to insert the same session again
	err = s.repo.Create(s.ctx, session)
	s.Require().Error(err)

	var domainErr *domain_err.DomainError
	s.Require().True(s.errorAs(err, &domainErr))
	s.Equal(domain_err.ErrDataConflict, domainErr.Kind())
}

func (s *SessionRepositoryTestSuite) TestFindByID_Found() {
	userID := s.insertTestUser()
	session, _, err := entity.NewSession(userID, 7*24*time.Hour, false)
	s.Require().NoError(err)
	s.Require().NoError(s.repo.Create(s.ctx, session))

	found, err := s.repo.FindByID(s.ctx, session.ID())
	s.Require().NoError(err)
	s.Equal(session.ID(), found.ID())
	s.Equal(session.UserID(), found.UserID())
	s.Equal(session.TokenHash(), found.TokenHash())
	s.Equal(session.JTI(), found.JTI())
	s.Equal(session.RememberMe(), found.RememberMe())
	s.False(found.IsRevoked())
}

func (s *SessionRepositoryTestSuite) TestFindByID_NotFound() {
	nonExistentID := value_object.NewSessionID()
	_, err := s.repo.FindByID(s.ctx, nonExistentID)
	s.Require().Error(err)

	var domainErr *domain_err.DomainError
	s.Require().True(s.errorAs(err, &domainErr))
	s.Equal(domain_err.ErrNotFoundSession, domainErr.Kind())
}

func (s *SessionRepositoryTestSuite) TestFindByTokenHash_Found() {
	userID := s.insertTestUser()
	session, _, err := entity.NewSession(userID, 7*24*time.Hour, true)
	s.Require().NoError(err)
	s.Require().NoError(s.repo.Create(s.ctx, session))

	found, err := s.repo.FindByTokenHash(s.ctx, session.TokenHash())
	s.Require().NoError(err)
	s.Equal(session.ID(), found.ID())
	s.True(found.RememberMe())
}

func (s *SessionRepositoryTestSuite) TestFindByTokenHash_NotFound() {
	_, err := s.repo.FindByTokenHash(s.ctx, "nonexistenthash")
	s.Require().Error(err)

	var domainErr *domain_err.DomainError
	s.Require().True(s.errorAs(err, &domainErr))
	s.Equal(domain_err.ErrNotFoundSession, domainErr.Kind())
}

func (s *SessionRepositoryTestSuite) TestFindByUserID_ReturnsActiveSessions() {
	userID := s.insertTestUser()

	// Create 2 active sessions
	session1, _, err := entity.NewSession(userID, 7*24*time.Hour, false)
	s.Require().NoError(err)
	s.Require().NoError(s.repo.Create(s.ctx, session1))

	session2, _, err := entity.NewSession(userID, 30*24*time.Hour, true)
	s.Require().NoError(err)
	s.Require().NoError(s.repo.Create(s.ctx, session2))

	// Create a revoked session
	session3, _, err := entity.NewSession(userID, 7*24*time.Hour, false)
	s.Require().NoError(err)
	session3.Revoke()
	s.Require().NoError(s.repo.Create(s.ctx, session3))

	sessions, err := s.repo.FindByUserID(s.ctx, userID)
	s.Require().NoError(err)
	// Only active (not revoked, not expired) sessions should be returned
	s.Len(sessions, 2)
}

func (s *SessionRepositoryTestSuite) TestUpdate_UpdateFields() {
	userID := s.insertTestUser()
	session, _, err := entity.NewSession(userID, 7*24*time.Hour, false)
	s.Require().NoError(err)
	s.Require().NoError(s.repo.Create(s.ctx, session))

	// Update last used
	session.UpdateLastUsed()
	err = s.repo.Update(s.ctx, session)
	s.Require().NoError(err)

	found, err := s.repo.FindByID(s.ctx, session.ID())
	s.Require().NoError(err)
	s.NotNil(found.LastUsedAt())
}

func (s *SessionRepositoryTestSuite) TestUpdate_Revoke() {
	userID := s.insertTestUser()
	session, _, err := entity.NewSession(userID, 7*24*time.Hour, false)
	s.Require().NoError(err)
	s.Require().NoError(s.repo.Create(s.ctx, session))

	session.Revoke()
	err = s.repo.Update(s.ctx, session)
	s.Require().NoError(err)

	found, err := s.repo.FindByID(s.ctx, session.ID())
	s.Require().NoError(err)
	s.True(found.IsRevoked())
	s.NotNil(found.RevokedAt())
}

func (s *SessionRepositoryTestSuite) TestUpdate_NotFound() {
	nonExistent, _, err := entity.NewSession(value_object.NewUserID(), 7*24*time.Hour, false)
	s.Require().NoError(err)

	err = s.repo.Update(s.ctx, nonExistent)
	s.Require().Error(err)

	var domainErr *domain_err.DomainError
	s.Require().True(s.errorAs(err, &domainErr))
	s.Equal(domain_err.ErrNotFoundSession, domainErr.Kind())
}

func (s *SessionRepositoryTestSuite) TestRevokeAllByUserID() {
	userID := s.insertTestUser()

	// Create 3 sessions
	for i := 0; i < 3; i++ {
		session, _, err := entity.NewSession(userID, 7*24*time.Hour, false)
		s.Require().NoError(err)
		s.Require().NoError(s.repo.Create(s.ctx, session))
	}

	count, err := s.repo.RevokeAllByUserID(s.ctx, userID)
	s.Require().NoError(err)
	s.Equal(3, count)

	// Active sessions should be 0
	sessions, err := s.repo.FindByUserID(s.ctx, userID)
	s.Require().NoError(err)
	s.Len(sessions, 0)
}

func (s *SessionRepositoryTestSuite) TestRevokeAllByUserID_NoSessions() {
	userID := s.insertTestUser()

	count, err := s.repo.RevokeAllByUserID(s.ctx, userID)
	s.Require().NoError(err)
	s.Equal(0, count)
}

func (s *SessionRepositoryTestSuite) TestDeleteExpired() {
	userID := s.insertTestUser()

	// Insert an already-expired session directly
	expiredID := value_object.NewSessionID()
	_, err := s.db.Exec(s.ctx,
		`INSERT INTO refresh_tokens (id, user_id, token_hash, jti, expires_at, created_at, is_revoked, remember_me)
		 VALUES ($1, $2, $3, $4, $5, NOW(), FALSE, FALSE)`,
		expiredID.String(), userID.String(), "expiredhash", "expired-jti",
		time.Now().Add(-1*time.Hour),
	)
	s.Require().NoError(err)

	// Insert a valid session
	session, _, err := entity.NewSession(userID, 7*24*time.Hour, false)
	s.Require().NoError(err)
	s.Require().NoError(s.repo.Create(s.ctx, session))

	err = s.repo.DeleteExpired(s.ctx)
	s.Require().NoError(err)

	// Expired session should be gone
	_, err = s.repo.FindByID(s.ctx, expiredID)
	s.Require().Error(err)

	// Valid session should remain
	found, err := s.repo.FindByID(s.ctx, session.ID())
	s.Require().NoError(err)
	s.NotNil(found)
}

// errorAs is a helper that wraps errors.As for DomainError.
func (s *SessionRepositoryTestSuite) errorAs(err error, target interface{}) bool {
	return errorAs(err, target)
}
