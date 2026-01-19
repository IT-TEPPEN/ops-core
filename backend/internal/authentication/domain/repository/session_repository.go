package repository

import (
	"context"

	"opscore/backend/internal/authentication/domain/entity"
	"opscore/backend/internal/authentication/domain/value_object"
)

// SessionRepository defines the interface for Session persistence.
// This repository manages authentication sessions (refresh tokens) for users.
type SessionRepository interface {
	// Create creates a new session.
	// This is called when a user successfully authenticates.
	//
	// Parameters:
	//   - ctx context.Context: The context for cancellation and timeout
	//   - session *entity.Session: The session entity to create
	//
	// Returns:
	//   - error: Domain error if creation fails
	//
	// Errors:
	//   - domain_err.ErrDataPersistFailure: When saving data fails
	//   - domain_err.ErrDataConflict: When session with same ID already exists
	Create(ctx context.Context, session *entity.Session) error

	// FindByID finds a session by ID.
	//
	// Parameters:
	//   - ctx context.Context: The context for cancellation and timeout
	//   - sessionID value_object.SessionID: The session ID to find
	//
	// Returns:
	//   - *entity.Session: The session entity
	//   - error: Domain error if operation fails
	//
	// Errors:
	//   - domain_err.ErrNotFoundSession: When session with given ID does not exist
	//   - domain_err.ErrDataAccessFailure: When loading data fails
	FindByID(ctx context.Context, sessionID value_object.SessionID) (*entity.Session, error)

	// FindByTokenHash finds a session by token hash.
	// This is used to validate refresh tokens.
	//
	// Parameters:
	//   - ctx context.Context: The context for cancellation and timeout
	//   - tokenHash string: SHA256 hash of the session token
	//
	// Returns:
	//   - *entity.Session: The session entity
	//   - error: Domain error if operation fails
	//
	// Errors:
	//   - domain_err.ErrNotFoundSession: When session with given token hash does not exist
	//   - domain_err.ErrDataAccessFailure: When loading data fails
	FindByTokenHash(ctx context.Context, tokenHash string) (*entity.Session, error)

	// FindByUserID finds all active sessions for a user.
	// This is used to list user's current sessions.
	//
	// Parameters:
	//   - ctx context.Context: The context for cancellation and timeout
	//   - userID value_object.UserID: The user ID to find sessions for
	//
	// Returns:
	//   - []*entity.Session: The list of sessions (may be empty)
	//   - error: Domain error if operation fails
	//
	// Errors:
	//   - domain_err.ErrDataAccessFailure: When loading data fails
	FindByUserID(ctx context.Context, userID value_object.UserID) ([]*entity.Session, error)

	// Update updates an existing session.
	// This is called when session metadata changes (e.g., last used timestamp).
	//
	// Parameters:
	//   - ctx context.Context: The context for cancellation and timeout
	//   - session *entity.Session: The session entity to update
	//
	// Returns:
	//   - error: Domain error if update fails
	//
	// Errors:
	//   - domain_err.ErrNotFoundSession: When session does not exist
	//   - domain_err.ErrDataPersistFailure: When saving data fails
	Update(ctx context.Context, session *entity.Session) error

	// RevokeAllByUserID revokes all sessions for a user.
	// This is called during logout or security events.
	//
	// Parameters:
	//   - ctx context.Context: The context for cancellation and timeout
	//   - userID value_object.UserID: The user ID to revoke sessions for
	//
	// Returns:
	//   - int: Number of sessions revoked
	//   - error: Domain error if operation fails
	//
	// Errors:
	//   - domain_err.ErrDataPersistFailure: When saving data fails
	RevokeAllByUserID(ctx context.Context, userID value_object.UserID) (int, error)

	// DeleteExpired deletes expired sessions (cleanup).
	// This is typically called by a background job to clean up old sessions.
	//
	// Parameters:
	//   - ctx context.Context: The context for cancellation and timeout
	//
	// Returns:
	//   - error: Domain error if operation fails
	//
	// Errors:
	//   - domain_err.ErrDataPersistFailure: When deleting data fails
	DeleteExpired(ctx context.Context) error
}
