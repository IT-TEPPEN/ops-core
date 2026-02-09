package usecase

import (
	"context"
	"errors"

	"opscore/backend/internal/authentication/application/dto"
	application_err "opscore/backend/internal/authentication/application/err"
	"opscore/backend/internal/authentication/application/service"
	"opscore/backend/internal/authentication/domain/entity"
	domain_err "opscore/backend/internal/authentication/domain/err"
	"opscore/backend/internal/authentication/domain/repository"
	"opscore/backend/internal/authentication/domain/value_object"
)

// IssueSessionFromProviderUser issues an authentication session for a user authenticated by an external provider.
// This usecase orchestrates user registration/update and session creation based on provider authentication.
type IssueSessionFromProviderUser interface {
	// Execute issues an authentication session from external provider user information.
	//
	// This usecase:
	//   1. Validates provider user information
	//   2. Registers a new user or updates an existing user based on provider identity
	//   3. Issues an authentication session (access token and refresh token)
	//
	// Parameters:
	//   - ctx: Context for cancellation and timeout
	//   - providerUserInfo: User information from the external authentication provider
	//
	// Returns:
	//   - SessionResult containing access token, refresh token, and user information
	//   - Error if validation, user management, or session issuance fails
	//
	// Errors:
	//   - ErrInvalidInput: When provider information is invalid or domain rules are violated (wraps domain errors with WithParent)
	//   - ErrDataAccessFailure: When user data cannot be read from the database
	//   - ErrDataPersistFailure: When user data cannot be saved to the database
	//   - ErrDataConflict: When unique constraints are violated (e.g., email already exists)
	//   - ErrAuthenticationFailed: When session cannot be issued (JWT/RefreshToken generation failure)
	//   - ErrUnexpected: When an unexpected error occurs
	Execute(ctx context.Context, providerUserInfo dto.ProviderUserInfo) (*dto.SessionResult, error)
}

type issueSessionFromProviderUserImpl struct {
	userRepository repository.UserRepository
	sessionService service.SessionService
}

func NewIssueSessionFromProviderUser(
	userRepository repository.UserRepository,
	sessionService service.SessionService,
) IssueSessionFromProviderUser {
	return &issueSessionFromProviderUserImpl{
		userRepository: userRepository,
		sessionService: sessionService,
	}
}

func (u *issueSessionFromProviderUserImpl) Execute(ctx context.Context, providerUserInfo dto.ProviderUserInfo) (*dto.SessionResult, error) {
	// Step 1: Validate provider user information
	if err := validateProviderUserInfo(providerUserInfo); err != nil {
		return nil, application_err.NewInvalidInputError("providerUserInfo", err.Error())
	}

	// Step 2: Find or create user
	user, err := u.findOrCreateUser(ctx, providerUserInfo)
	if err != nil {
		return nil, err
	}

	// Step 3: Issue session
	session, err := u.sessionService.CreateSession(ctx, user.ID(), false)
	if err != nil {
		return nil, application_err.NewAuthenticationFailedError("session creation failed").WithParent(err)
	}

	return session, nil
}

func validateProviderUserInfo(info dto.ProviderUserInfo) error {
	if info.Provider == "" {
		return errors.New("provider is required")
	}
	if info.Subject == "" {
		return errors.New("subject is required")
	}
	if info.Email == "" {
		return errors.New("email is required")
	}
	return nil
}

func (u *issueSessionFromProviderUserImpl) findOrCreateUser(ctx context.Context, providerUserInfo dto.ProviderUserInfo) (*entity.User, error) {
	// Try to find existing user by provider identity
	user, err := u.userRepository.FindByProviderIdentity(ctx, providerUserInfo.Provider, providerUserInfo.Subject)

	if err != nil {
		// Check if it's a "not found" error
		if domainErr, ok := err.(*domain_err.DomainError); ok && domainErr.Kind() == domain_err.ErrNotFoundUser {
			return u.createNewUser(ctx, providerUserInfo)
		}
		// Other database errors
		return nil, application_err.NewUnexpectedError("failed to find user by provider identity").WithParent(err)
	}

	// User exists, update last login and identity information
	return u.updateExistingUser(ctx, user, providerUserInfo)
}

func (u *issueSessionFromProviderUserImpl) createNewUser(ctx context.Context, providerUserInfo dto.ProviderUserInfo) (*entity.User, error) {
	// Create new user entity
	userID := value_object.NewUserID()
	user := entity.NewUser(
		userID,
		providerUserInfo.Email,
		providerUserInfo.Name,
		providerUserInfo.Picture,
	)

	// Create identity
	identityID := value_object.NewIdentityID()
	identity, err := entity.NewIdentity(
		identityID,
		userID,
		providerUserInfo.Provider,
		providerUserInfo.Subject,
		providerUserInfo.Email,
		providerUserInfo.Name,
		providerUserInfo.Picture,
	)
	if err != nil {
		return nil, application_err.NewInvalidInputError("identity", "failed to create identity").WithParent(err)
	}

	// Add identity to user
	if err := user.AddIdentity(identity); err != nil {
		return nil, application_err.NewInvalidInputError("identity", "failed to add identity to user").WithParent(err)
	}

	// Save user to repository
	if err := u.userRepository.Save(ctx, user); err != nil {
		// Check error type
		if domainErr, ok := err.(*domain_err.DomainError); ok {
			switch domainErr.Kind() {
			case domain_err.ErrDataConflict:
				return nil, application_err.NewUnexpectedError("user already exists").WithParent(err)
			case domain_err.ErrDataPersistFailure:
				return nil, application_err.NewUnexpectedError("failed to save user").WithParent(err)
			}
		}
		return nil, application_err.NewUnexpectedError("failed to save user").WithParent(err)
	}

	return user, nil
}

func (u *issueSessionFromProviderUserImpl) updateExistingUser(ctx context.Context, user *entity.User, providerUserInfo dto.ProviderUserInfo) (*entity.User, error) {
	// Update last login time
	user.UpdateLastLogin()

	// Update identity information if changed
	identities := user.Identities()
	for _, identity := range identities {
		if identity.Provider() == providerUserInfo.Provider && identity.ProviderUserID() == providerUserInfo.Subject {
			identity.UpdateInfo(providerUserInfo.Email, providerUserInfo.Name, providerUserInfo.Picture)
			break
		}
	}

	// Save updated user
	if err := u.userRepository.Save(ctx, user); err != nil {
		// Check error type
		if domainErr, ok := err.(*domain_err.DomainError); ok {
			switch domainErr.Kind() {
			case domain_err.ErrDataConflict:
				return nil, application_err.NewUnexpectedError("failed to update user due to conflict").WithParent(err)
			case domain_err.ErrDataPersistFailure:
				return nil, application_err.NewUnexpectedError("failed to save user").WithParent(err)
			}
		}
		return nil, application_err.NewUnexpectedError("failed to save user").WithParent(err)
	}

	return user, nil
}
