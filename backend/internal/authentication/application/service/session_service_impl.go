package service

import (
	"context"
	"time"

	"opscore/backend/internal/authentication/application/dto"
	application_err "opscore/backend/internal/authentication/application/err"
	"opscore/backend/internal/authentication/application/port"
	"opscore/backend/internal/authentication/domain/entity"
	domain_err "opscore/backend/internal/authentication/domain/err"
	"opscore/backend/internal/authentication/domain/repository"
	"opscore/backend/internal/authentication/domain/value_object"
)

const (
	// AccessTokenDuration defines the lifetime of JWT access tokens (15 minutes)
	AccessTokenDuration = 15 * time.Minute
	// ShortSessionDuration defines the lifetime of short-lived refresh tokens (7 days)
	ShortSessionDuration = 7 * 24 * time.Hour
	// LongSessionDuration defines the lifetime of long-lived refresh tokens (30 days)
	LongSessionDuration = 30 * 24 * time.Hour
)

// sessionServiceImpl implements SessionService
type sessionServiceImpl struct {
	userRepository    repository.UserRepository
	sessionRepository repository.SessionRepository
	jwtProvider       port.JWTProvider
}

// NewSessionService creates a new instance of SessionService
func NewSessionService(
	userRepository repository.UserRepository,
	sessionRepository repository.SessionRepository,
	jwtProvider port.JWTProvider,
) SessionService {
	return &sessionServiceImpl{
		userRepository:    userRepository,
		sessionRepository: sessionRepository,
		jwtProvider:       jwtProvider,
	}
}

// CreateSession creates a new authentication session for the user
func (s *sessionServiceImpl) CreateSession(ctx context.Context, userID value_object.UserID, rememberMe bool) (*dto.SessionResult, error) {
	// 1. Load user information
	user, err := s.userRepository.FindByID(ctx, userID)
	if err != nil {
		if domainErr, ok := err.(*domain_err.DomainError); ok {
			switch domainErr.Kind() {
			case domain_err.ErrNotFoundUser:
				return nil, application_err.NewUnexpectedError("user not found").WithParent(err)
			case domain_err.ErrDataAccessFailure:
				return nil, application_err.NewUnexpectedError("failed to load user").WithParent(err)
			default:
				return nil, application_err.NewUnexpectedError("unexpected error while loading user").WithParent(err)
			}
		}
		return nil, application_err.NewUnexpectedError("unexpected error while loading user").WithParent(err)
	}

	// 2. Determine session duration based on rememberMe
	sessionDuration := ShortSessionDuration
	if rememberMe {
		sessionDuration = LongSessionDuration
	}

	// 3. Create refresh token (Session entity)
	session, rawToken, err := entity.NewSession(userID, sessionDuration, rememberMe)
	if err != nil {
		return nil, application_err.NewAuthenticationFailedError("failed to create refresh token").WithParent(err)
	}

	// 4. Persist the session
	err = s.sessionRepository.Create(ctx, session)
	if err != nil {
		if domainErr, ok := err.(*domain_err.DomainError); ok {
			switch domainErr.Kind() {
			case domain_err.ErrDataPersistFailure:
				return nil, application_err.NewUnexpectedError("failed to save session").WithParent(err)
			case domain_err.ErrDataConflict:
				return nil, application_err.NewUnexpectedError("session conflict").WithParent(err)
			default:
				return nil, application_err.NewUnexpectedError("unexpected error while saving session").WithParent(err)
			}
		}
		return nil, application_err.NewUnexpectedError("unexpected error while saving session").WithParent(err)
	}

	// 5. Generate JWT access token
	accessToken, err := s.jwtProvider.GenerateAccessToken(
		userID.String(),
		session.JTI(),
		AccessTokenDuration,
	)
	if err != nil {
		return nil, application_err.NewAuthenticationFailedError("failed to generate access token").WithParent(err)
	}

	// 6. Build and return result
	userInfo := dto.UserInfo{
		UserID:      user.ID().String(),
		Email:       user.Email(),
		DisplayName: user.DisplayName(),
		PictureURL:  user.PictureURL(),
		CreatedAt:   user.CreatedAt(),
		LastLoginAt: user.LastLoginAt(),
	}

	return &dto.SessionResult{
		AccessToken:  accessToken,
		RefreshToken: rawToken,
		User:         userInfo,
	}, nil
}
