package usecase

import (
	"context"

	"opscore/backend/internal/authentication/application/dto"
	application_err "opscore/backend/internal/authentication/application/err"
	"opscore/backend/internal/authentication/application/service"
	"opscore/backend/internal/authentication/domain/entity"
	domain_err "opscore/backend/internal/authentication/domain/err"
	"opscore/backend/internal/authentication/domain/repository"
	"opscore/backend/internal/authentication/domain/value_object"
)

// CreateSession creates a new authentication session for a user.
// This usecase generates a refresh token and stores the session in the database.
type CreateSession interface {
	// Execute creates a new session for the specified user.
	// The session duration is determined by the sessionType (short: 7 days, long: 30 days).
	//
	// Parameters:
	//   - ctx context.Context: The context for cancellation and timeout
	//   - dto dto.CreateSessionRequest: Request containing userID and sessionType
	//
	// Returns:
	//   - *dto.SessionInfo: Created session information with session ID and expiration
	//   - error: Error if session creation fails
	//
	// Errors:
	//   - application_err.ErrInvalidInput: When userID is invalid or sessionType is unknown
	//   - application_err.ErrUserNotFound: When specified user does not exist
	//   - application_err.ErrDataPersistFailure: When saving session data fails
	//   - application_err.ErrUnexpected: When session token generation fails
	Execute(ctx context.Context, dto dto.CreateSessionRequest) (*dto.SessionInfo, error)
}

type createSessionImpl struct {
	userRepository    repository.UserRepository
	sessionRepository repository.SessionRepository
}

func NewCreateSession(
	userRepository repository.UserRepository,
	sessionRepository repository.SessionRepository,
) CreateSession {
	return &createSessionImpl{
		userRepository:    userRepository,
		sessionRepository: sessionRepository,
	}
}

func (c *createSessionImpl) Execute(ctx context.Context, req dto.CreateSessionRequest) (*dto.SessionInfo, error) {
	// 1. Validate input
	userID, err := value_object.UserIDFromString(req.UserID)
	if err != nil {
		return nil, application_err.NewInvalidInputError("userID", "invalid user ID format").WithParent(err)
	}

	rememberMe, err := resolveSessionType(req.SessionType)
	if err != nil {
		return nil, err
	}

	// 2. Verify user exists
	_, err = c.userRepository.FindByID(ctx, userID)
	if err != nil {
		if domainErr, ok := err.(*domain_err.DomainError); ok {
			if domainErr.Kind() == domain_err.ErrNotFoundUser {
				return nil, application_err.NewUserNotFoundError(userID.String()).WithParent(err)
			}
		}
		return nil, application_err.NewUnexpectedError("failed to verify user existence").WithParent(err)
	}

	// 3. Determine session duration
	sessionDuration := service.ShortSessionDuration
	if rememberMe {
		sessionDuration = service.LongSessionDuration
	}

	// 4. Create session entity
	session, _, err := entity.NewSession(userID, sessionDuration, rememberMe)
	if err != nil {
		return nil, application_err.NewUnexpectedError("failed to generate session token").WithParent(err)
	}

	// 5. Persist session
	if err := c.sessionRepository.Create(ctx, session); err != nil {
		if domainErr, ok := err.(*domain_err.DomainError); ok {
			if domainErr.Kind() == domain_err.ErrDataPersistFailure {
				return nil, application_err.NewDataPersistFailureError("session").WithParent(err)
			}
		}
		return nil, application_err.NewUnexpectedError("failed to save session").WithParent(err)
	}

	// 6. Build result
	sessionType := dto.SessionTypeShort
	if rememberMe {
		sessionType = dto.SessionTypeLong
	}

	return &dto.SessionInfo{
		SessionID:   session.ID().String(),
		UserID:      session.UserID().String(),
		SessionType: sessionType,
		CreatedAt:   session.CreatedAt(),
		ExpiresAt:   session.ExpiresAt(),
		LastUsedAt:  session.LastUsedAt(),
		IsRevoked:   session.IsRevoked(),
	}, nil
}

func resolveSessionType(sessionType dto.SessionType) (bool, error) {
	switch sessionType {
	case dto.SessionTypeShort:
		return false, nil
	case dto.SessionTypeLong:
		return true, nil
	default:
		return false, application_err.NewInvalidInputError("sessionType", "must be 'short' or 'long'")
	}
}
