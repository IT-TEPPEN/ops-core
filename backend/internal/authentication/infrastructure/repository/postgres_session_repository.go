package repository

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"

	"opscore/backend/internal/authentication/domain/entity"
	domain_err "opscore/backend/internal/authentication/domain/err"
	"opscore/backend/internal/authentication/domain/repository"
	"opscore/backend/internal/authentication/domain/value_object"
)

// PostgresSessionRepository is the PostgreSQL implementation of SessionRepository.
type PostgresSessionRepository struct {
	db *pgxpool.Pool
}

// NewPostgresSessionRepository creates a new PostgresSessionRepository.
func NewPostgresSessionRepository(db *pgxpool.Pool) repository.SessionRepository {
	return &PostgresSessionRepository{db: db}
}

func (r *PostgresSessionRepository) Create(ctx context.Context, session *entity.Session) error {
	query := `
		INSERT INTO refresh_tokens (id, user_id, token_hash, jti, expires_at, created_at, last_used_at, revoked_at, is_revoked, remember_me)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`

	_, err := r.db.Exec(ctx, query,
		session.ID().String(),
		session.UserID().String(),
		session.TokenHash(),
		session.JTI(),
		session.ExpiresAt(),
		session.CreatedAt(),
		session.LastUsedAt(),
		session.RevokedAt(),
		session.IsRevoked(),
		session.RememberMe(),
	)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return domain_err.NewDataConflictError("refresh_tokens", pgErr.ConstraintName).WithParent(err)
		}
		return domain_err.NewDataPersistFailureError("session").WithParent(err)
	}

	return nil
}

func (r *PostgresSessionRepository) FindByID(ctx context.Context, sessionID value_object.SessionID) (*entity.Session, error) {
	query := `
		SELECT id, user_id, token_hash, jti, expires_at, created_at, last_used_at, revoked_at, is_revoked, remember_me
		FROM refresh_tokens
		WHERE id = $1`

	return r.scanSession(ctx, query, sessionID.String())
}

func (r *PostgresSessionRepository) FindByTokenHash(ctx context.Context, tokenHash string) (*entity.Session, error) {
	query := `
		SELECT id, user_id, token_hash, jti, expires_at, created_at, last_used_at, revoked_at, is_revoked, remember_me
		FROM refresh_tokens
		WHERE token_hash = $1`

	return r.scanSession(ctx, query, tokenHash)
}

func (r *PostgresSessionRepository) FindByUserID(ctx context.Context, userID value_object.UserID) ([]*entity.Session, error) {
	query := `
		SELECT id, user_id, token_hash, jti, expires_at, created_at, last_used_at, revoked_at, is_revoked, remember_me
		FROM refresh_tokens
		WHERE user_id = $1 AND is_revoked = FALSE AND expires_at > NOW()
		ORDER BY created_at DESC`

	rows, err := r.db.Query(ctx, query, userID.String())
	if err != nil {
		return nil, domain_err.NewDataAccessFailureError("sessions").WithParent(err)
	}
	defer rows.Close()

	var sessions []*entity.Session
	for rows.Next() {
		session, err := r.scanSessionFromRow(rows)
		if err != nil {
			return nil, domain_err.NewDataAccessFailureError("sessions").WithParent(err)
		}
		sessions = append(sessions, session)
	}
	if err := rows.Err(); err != nil {
		return nil, domain_err.NewDataAccessFailureError("sessions").WithParent(err)
	}

	return sessions, nil
}

func (r *PostgresSessionRepository) Update(ctx context.Context, session *entity.Session) error {
	query := `
		UPDATE refresh_tokens
		SET last_used_at = $1, revoked_at = $2, is_revoked = $3
		WHERE id = $4`

	result, err := r.db.Exec(ctx, query,
		session.LastUsedAt(),
		session.RevokedAt(),
		session.IsRevoked(),
		session.ID().String(),
	)
	if err != nil {
		return domain_err.NewDataPersistFailureError("session").WithParent(err)
	}

	if result.RowsAffected() == 0 {
		return domain_err.NewNotFoundSessionError(session.ID().String())
	}

	return nil
}

func (r *PostgresSessionRepository) RevokeAllByUserID(ctx context.Context, userID value_object.UserID) (int, error) {
	query := `
		UPDATE refresh_tokens
		SET is_revoked = TRUE, revoked_at = $1
		WHERE user_id = $2 AND is_revoked = FALSE`

	now := time.Now()
	result, err := r.db.Exec(ctx, query, now, userID.String())
	if err != nil {
		return 0, domain_err.NewDataPersistFailureError("sessions").WithParent(err)
	}

	return int(result.RowsAffected()), nil
}

func (r *PostgresSessionRepository) DeleteExpired(ctx context.Context) error {
	query := `DELETE FROM refresh_tokens WHERE expires_at < NOW()`

	_, err := r.db.Exec(ctx, query)
	if err != nil {
		return domain_err.NewDataPersistFailureError("expired sessions").WithParent(err)
	}

	return nil
}

// scanSession executes a single-row query and returns a Session entity.
func (r *PostgresSessionRepository) scanSession(ctx context.Context, query string, arg interface{}) (*entity.Session, error) {
	var (
		id         string
		userID     string
		tokenHash  string
		jti        string
		expiresAt  time.Time
		createdAt  time.Time
		lastUsedAt *time.Time
		revokedAt  *time.Time
		isRevoked  bool
		rememberMe bool
	)

	err := r.db.QueryRow(ctx, query, arg).Scan(
		&id, &userID, &tokenHash, &jti, &expiresAt, &createdAt,
		&lastUsedAt, &revokedAt, &isRevoked, &rememberMe,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain_err.NewNotFoundSessionError(formatArg(arg))
		}
		return nil, domain_err.NewDataAccessFailureError("session").WithParent(err)
	}

	return reconstructSession(id, userID, tokenHash, jti, expiresAt, createdAt, lastUsedAt, revokedAt, isRevoked, rememberMe)
}

// scanSessionFromRow scans a session from an already-iterated row.
func (r *PostgresSessionRepository) scanSessionFromRow(rows pgx.Rows) (*entity.Session, error) {
	var (
		id         string
		userID     string
		tokenHash  string
		jti        string
		expiresAt  time.Time
		createdAt  time.Time
		lastUsedAt *time.Time
		revokedAt  *time.Time
		isRevoked  bool
		rememberMe bool
	)

	err := rows.Scan(
		&id, &userID, &tokenHash, &jti, &expiresAt, &createdAt,
		&lastUsedAt, &revokedAt, &isRevoked, &rememberMe,
	)
	if err != nil {
		return nil, err
	}

	return reconstructSession(id, userID, tokenHash, jti, expiresAt, createdAt, lastUsedAt, revokedAt, isRevoked, rememberMe)
}

// reconstructSession creates a Session entity from raw database fields.
// Database data is trusted, so we use Reconstruct* functions (no validation).
func reconstructSession(id, userID, tokenHash, jti string, expiresAt, createdAt time.Time, lastUsedAt, revokedAt *time.Time, isRevoked, rememberMe bool) (*entity.Session, error) {
	sessionID := value_object.ReconstructSessionID(id)
	uid := value_object.ReconstructUserID(userID)

	return entity.ReconstructSession(sessionID, uid, tokenHash, jti, expiresAt, createdAt, lastUsedAt, revokedAt, isRevoked, rememberMe), nil
}

// formatArg converts an argument to string for error messages.
func formatArg(arg interface{}) string {
	if s, ok := arg.(string); ok {
		return s
	}
	return "unknown"
}
