package application_err

import (
	"fmt"
	"opscore/backend/internal/shared/errors"
)

type ApplicationErrorKind string

const (
	// ### Input Validation Errors (0000 - 0999) ###############################
	ErrInvalidInput ApplicationErrorKind = "ATA0001"

	// ### Authentication Errors (1000 - 1999) #################################
	ErrAuthenticationFailed       ApplicationErrorKind = "ATA1000"
	ErrSessionExpired             ApplicationErrorKind = "ATA1001"
	ErrSessionRevoked             ApplicationErrorKind = "ATA1002"
	ErrSessionNotFound            ApplicationErrorKind = "ATA1003"
	ErrTokenGenerationError       ApplicationErrorKind = "ATA1004"
	ErrTokenValidationError       ApplicationErrorKind = "ATA1005"
	ErrUnsupportedProvider        ApplicationErrorKind = "ATA1006"
	ErrProviderConfigurationError ApplicationErrorKind = "ATA1007"

	// ### Authorization Errors (2000 - 2999) ##################################
	ErrResourceNotOwned ApplicationErrorKind = "ATA2001"

	// ### Business Logic Errors (3000 - 3999) #################################
	ErrIdentityAlreadyLinked ApplicationErrorKind = "ATA3000"
	ErrUserNotFound          ApplicationErrorKind = "ATA3001"

	// ### Data Access Errors (4000 - 4999) ####################################
	ErrDataAccessFailure  ApplicationErrorKind = "ATA4000"
	ErrDataPersistFailure ApplicationErrorKind = "ATA4001"

	// ### Unexpected Errors (9000 - 9999) #####################################
	ErrUnexpected ApplicationErrorKind = "ATA9999"
)

var factory = errors.NewFactory(errors.FeatureAuthentication, errors.LayerApplication)

var (
	// ### Input Validation Errors ############################################
	buildInvalidInputError = factory.MustBuilder(string(ErrInvalidInput))

	// ### Authentication Errors ##############################################
	buildAuthenticationFailedError  = factory.MustBuilder(string(ErrAuthenticationFailed))
	buildSessionExpiredError        = factory.MustBuilder(string(ErrSessionExpired))
	buildSessionRevokedError        = factory.MustBuilder(string(ErrSessionRevoked))
	buildSessionNotFoundError       = factory.MustBuilder(string(ErrSessionNotFound))
	buildTokenGenerationError       = factory.MustBuilder(string(ErrTokenGenerationError))
	buildTokenValidationError       = factory.MustBuilder(string(ErrTokenValidationError))
	buildUnsupportedProviderError   = factory.MustBuilder(string(ErrUnsupportedProvider))
	buildProviderConfigurationError = factory.MustBuilder(string(ErrProviderConfigurationError))

	// ### Authorization Errors ###############################################
	buildResourceNotOwnedError = factory.MustBuilder(string(ErrResourceNotOwned))

	// ### Business Logic Errors ##############################################
	buildIdentityAlreadyLinkedError = factory.MustBuilder(string(ErrIdentityAlreadyLinked))
	buildUserNotFoundError          = factory.MustBuilder(string(ErrUserNotFound))

	// ### Data Access Errors #################################################
	buildDataAccessFailureError  = factory.MustBuilder(string(ErrDataAccessFailure))
	buildDataPersistFailureError = factory.MustBuilder(string(ErrDataPersistFailure))

	// ### Unexpected Errors ##################################################
	buildUnexpectedError = factory.MustBuilder(string(ErrUnexpected))
)

type ApplicationError struct {
	*errors.CustomError
}

func (e *ApplicationError) Kind() ApplicationErrorKind {
	return ApplicationErrorKind(e.CustomError.Kind())
}

func newApplicationError(builder func(string) *errors.CustomError, message string) *ApplicationError {
	return &ApplicationError{
		CustomError: builder(message),
	}
}

func (e *ApplicationError) WithParent(err error) *ApplicationError {
	return &ApplicationError{
		CustomError: e.CustomError.WithParent(err),
	}
}

// ### Input Validation Errors ################################################

func NewInvalidInputError(field string, reason string) *ApplicationError {
	return newApplicationError(
		buildInvalidInputError,
		fmt.Sprintf("Invalid input for field '%s': %s", field, reason),
	)
}

// ### Authentication Errors ##################################################

func NewAuthenticationFailedError(reason string) *ApplicationError {
	return newApplicationError(
		buildAuthenticationFailedError,
		fmt.Sprintf("Authentication failed: %s", reason),
	)
}

func NewSessionExpiredError(sessionID string) *ApplicationError {
	return newApplicationError(
		buildSessionExpiredError,
		fmt.Sprintf("Session expired: %s", sessionID),
	)
}

func NewSessionRevokedError(sessionID string) *ApplicationError {
	return newApplicationError(
		buildSessionRevokedError,
		fmt.Sprintf("Session has been revoked: %s", sessionID),
	)
}

func NewSessionNotFoundError(sessionID string) *ApplicationError {
	return newApplicationError(
		buildSessionNotFoundError,
		fmt.Sprintf("Session not found: %s", sessionID),
	)
}

func NewTokenGenerationError(tokenType string) *ApplicationError {
	return newApplicationError(
		buildTokenGenerationError,
		fmt.Sprintf("Failed to generate %s token", tokenType),
	)
}

func NewTokenValidationError(reason string) *ApplicationError {
	return newApplicationError(
		buildTokenValidationError,
		fmt.Sprintf("Token validation failed: %s", reason),
	)
}

func NewUnsupportedProviderError(provider string) *ApplicationError {
	return newApplicationError(
		buildUnsupportedProviderError,
		fmt.Sprintf("Provider '%s' is not supported", provider),
	)
}

func NewProviderConfigurationError(provider string, detail string) *ApplicationError {
	return newApplicationError(
		buildProviderConfigurationError,
		fmt.Sprintf("Provider '%s' configuration error: %s", provider, detail),
	)
}

// ### Authorization Errors ###################################################

func NewResourceNotOwnedError(resourceType string, resourceID string, userID string) *ApplicationError {
	return newApplicationError(
		buildResourceNotOwnedError,
		fmt.Sprintf("User %s does not own %s: %s", userID, resourceType, resourceID),
	)
}

// ### Business Logic Errors ##################################################

func NewIdentityAlreadyLinkedError(provider string, userID string) *ApplicationError {
	return newApplicationError(
		buildIdentityAlreadyLinkedError,
		fmt.Sprintf("Identity from provider '%s' is already linked to user %s", provider, userID),
	)
}

func NewUserNotFoundError(userID string) *ApplicationError {
	return newApplicationError(
		buildUserNotFoundError,
		fmt.Sprintf("User not found: %s", userID),
	)
}

// ### Data Access Errors #####################################################

func NewDataAccessFailureError(resource string) *ApplicationError {
	return newApplicationError(
		buildDataAccessFailureError,
		fmt.Sprintf("Failed to access data: %s", resource),
	)
}

func NewDataPersistFailureError(resource string) *ApplicationError {
	return newApplicationError(
		buildDataPersistFailureError,
		fmt.Sprintf("Failed to persist data: %s", resource),
	)
}

// ### Unexpected Errors ######################################################

func NewUnexpectedError(detail string) *ApplicationError {
	return newApplicationError(
		buildUnexpectedError,
		fmt.Sprintf("Unexpected error occurred: %s", detail),
	)
}
