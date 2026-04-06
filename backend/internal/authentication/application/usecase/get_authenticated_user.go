package usecase

import (
	"context"

	"opscore/backend/internal/authentication/application/dto"
	application_err "opscore/backend/internal/authentication/application/err"
	domain_err "opscore/backend/internal/authentication/domain/err"
	"opscore/backend/internal/authentication/domain/repository"
	"opscore/backend/internal/authentication/domain/value_object"
)

// GetAuthenticatedUser retrieves the currently authenticated user from a session.
// This usecase validates the session and returns the user information.
type GetAuthenticatedUser interface {
	// Execute validates the session and returns the authenticated user's information.
	// The session must be valid, not expired, and not revoked.
	//
	// Parameters:
	//   - ctx context.Context: The context for cancellation and timeout
	//   - dto dto.SessionIDRequest: Request containing sessionID to validate
	//
	// Returns:
	//   - *dto.UserInfo: Authenticated user's information
	//   - error: Error if session is invalid or user retrieval fails
	//
	// Errors:
	//   - application_err.ErrInvalidInput: When sessionID format is invalid
	//   - application_err.ErrSessionNotFound: When session with given ID does not exist
	//   - application_err.ErrSessionExpired: When session has expired
	//   - application_err.ErrSessionRevoked: When session has been revoked
	//   - application_err.ErrUserNotFound: When user associated with session does not exist
	//   - application_err.ErrUnexpected: When loading session or user data fails
	Execute(ctx context.Context, dto dto.SessionIDRequest) (*dto.UserInfo, error)
}

type getAuthenticatedUserImpl struct {
	sessionRepository repository.SessionRepository
	userRepository    repository.UserRepository
}

func NewGetAuthenticatedUser(
	sessionRepository repository.SessionRepository,
	userRepository repository.UserRepository,
) GetAuthenticatedUser {
	return &getAuthenticatedUserImpl{
		sessionRepository: sessionRepository,
		userRepository:    userRepository,
	}
}

func (g *getAuthenticatedUserImpl) Execute(ctx context.Context, req dto.SessionIDRequest) (*dto.UserInfo, error) {
	// 1. Validate session ID format
	sessionID, err := value_object.SessionIDFromString(req.SessionID)
	if err != nil {
		return nil, application_err.NewInvalidInputError("sessionID", "invalid session ID format").WithParent(err)
	}

	// 2. Find session
	session, err := g.sessionRepository.FindByID(ctx, sessionID)
	if err != nil {
		if domainErr, ok := err.(*domain_err.DomainError); ok {
			if domainErr.Kind() == domain_err.ErrNotFoundSession {
				return nil, application_err.NewSessionNotFoundError(sessionID.String()).WithParent(err)
			}
		}
		return nil, application_err.NewUnexpectedError("failed to load session").WithParent(err)
	}

	// 3. Check if session is revoked
	if session.IsRevoked() {
		return nil, application_err.NewSessionRevokedError(sessionID.String())
	}

	// 4. Check if session is expired (check expiry separately to return specific error)
	if !session.IsValid() {
		return nil, application_err.NewSessionExpiredError(sessionID.String())
	}

	// 5. Find user
	user, err := g.userRepository.FindByID(ctx, session.UserID())
	if err != nil {
		if domainErr, ok := err.(*domain_err.DomainError); ok {
			if domainErr.Kind() == domain_err.ErrNotFoundUser {
				return nil, application_err.NewUserNotFoundError(session.UserID().String()).WithParent(err)
			}
		}
		return nil, application_err.NewUnexpectedError("failed to load user").WithParent(err)
	}

	// 6. Return user info
	return &dto.UserInfo{
		UserID:      user.ID().String(),
		Email:       user.Email(),
		DisplayName: user.DisplayName(),
		PictureURL:  user.PictureURL(),
		CreatedAt:   user.CreatedAt(),
		LastLoginAt: user.LastLoginAt(),
	}, nil
}
