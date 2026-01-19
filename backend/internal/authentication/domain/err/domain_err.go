package domain_err

import (
	"fmt"
	"opscore/backend/internal/shared/errors"
)

type DomainErrorKind string

const (
	// ### ValueObject Errors (0000 - 0999) ####################################
	ErrInvalidUserID      DomainErrorKind = "ATD0001"
	ErrInvalidSessionID   DomainErrorKind = "ATD0002"
	ErrInvalidSessionType DomainErrorKind = "ATD0004"

	// ### Entity Errors (1000 - 1999) #########################################
	ErrInvalidProviderUserID       DomainErrorKind = "ATD1003"
	ErrCannotRemoveOnlyIdentity    DomainErrorKind = "ATD1005"
	ErrCannotUnlinkPrimaryIdentity DomainErrorKind = "ATD1006"
	ErrDuplicateIdentity           DomainErrorKind = "ATD1007"

	// ### Domain Service Errors (2000 - 2999) #################################
	// (No domain services defined yet)

	// ### Repository(IF) Errors (3000 - 3999) #################################
	ErrNotFoundUser       DomainErrorKind = "ATD3000"
	ErrNotFoundSession    DomainErrorKind = "ATD3001"
	ErrDataAccessFailure  DomainErrorKind = "ATD3002"
	ErrDataPersistFailure DomainErrorKind = "ATD3003"
	ErrDataConflict       DomainErrorKind = "ATD3004"

	// ### Unexpected Errors (9000 - 9999) #####################################
	ErrUnexpected DomainErrorKind = "ATD9999"
)

var factory = errors.NewFactory(errors.FeatureAuthentication, errors.LayerDomain)

var (
	// ### ValueObject Errors #################################################
	buildInvalidUserIDError      = factory.MustBuilder(string(ErrInvalidUserID))
	buildInvalidSessionIDError   = factory.MustBuilder(string(ErrInvalidSessionID))
	buildInvalidSessionTypeError = factory.MustBuilder(string(ErrInvalidSessionType))

	// ### Entity Errors ######################################################
	buildInvalidProviderUserIDError       = factory.MustBuilder(string(ErrInvalidProviderUserID))
	buildCannotRemoveOnlyIdentityError    = factory.MustBuilder(string(ErrCannotRemoveOnlyIdentity))
	buildCannotUnlinkPrimaryIdentityError = factory.MustBuilder(string(ErrCannotUnlinkPrimaryIdentity))
	buildDuplicateIdentityError           = factory.MustBuilder(string(ErrDuplicateIdentity))

	// ### Domain Service Errors ##############################################
	// (No domain services defined yet)

	// ### Repository(IF) Errors ##############################################
	buildNotFoundUserError       = factory.MustBuilder(string(ErrNotFoundUser))
	buildNotFoundSessionError    = factory.MustBuilder(string(ErrNotFoundSession))
	buildDataAccessFailureError  = factory.MustBuilder(string(ErrDataAccessFailure))
	buildDataPersistFailureError = factory.MustBuilder(string(ErrDataPersistFailure))
	buildDataConflictError       = factory.MustBuilder(string(ErrDataConflict))

	// ### Unexpected Errors ##################################################
	buildUnexpectedError = factory.MustBuilder(string(ErrUnexpected))
)

type DomainError struct {
	*errors.CustomError
}

func (e *DomainError) Kind() DomainErrorKind {
	return DomainErrorKind(e.CustomError.Kind())
}

func newDomainError(builder func(string) *errors.CustomError, message string) *DomainError {
	return &DomainError{
		CustomError: builder(message),
	}
}

func (e *DomainError) WithParent(err error) *DomainError {
	return &DomainError{
		CustomError: e.CustomError.WithParent(err),
	}
}

// ### ValueObject Errors #####################################################

func NewInvalidUserIDError(value string, reason string) *DomainError {
	return newDomainError(
		buildInvalidUserIDError,
		fmt.Sprintf("Invalid UserID '%s': %s", value, reason),
	)
}

func NewInvalidSessionIDError(value string, reason string) *DomainError {
	return newDomainError(
		buildInvalidSessionIDError,
		fmt.Sprintf("Invalid SessionID '%s': %s", value, reason),
	)
}

func NewInvalidSessionTypeError(value string) *DomainError {
	return newDomainError(
		buildInvalidSessionTypeError,
		fmt.Sprintf("Invalid SessionType: %s", value),
	)
}

// ### Entity Errors ##########################################################

func NewInvalidProviderUserIDError(providerUserID string) *DomainError {
	return newDomainError(
		buildInvalidProviderUserIDError,
		fmt.Sprintf("Invalid provider user ID: %s", providerUserID),
	)
}

func NewCannotRemoveOnlyIdentityError(userID string) *DomainError {
	return newDomainError(
		buildCannotRemoveOnlyIdentityError,
		fmt.Sprintf("Cannot remove the only identity for user: %s", userID),
	)
}

func NewCannotUnlinkPrimaryIdentityError(identityID string) *DomainError {
	return newDomainError(
		buildCannotUnlinkPrimaryIdentityError,
		fmt.Sprintf("Cannot unlink primary identity: %s", identityID),
	)
}

func NewDuplicateIdentityError(provider string, providerUserID string) *DomainError {
	return newDomainError(
		buildDuplicateIdentityError,
		fmt.Sprintf("Identity already exists: provider=%s, providerUserID=%s", provider, providerUserID),
	)
}

// ### Domain Service Errors ##################################################
// (No domain services defined yet)

// ### Repository(IF) Errors ##################################################

func NewNotFoundUserError(userID string) *DomainError {
	return newDomainError(
		buildNotFoundUserError,
		fmt.Sprintf("User not found: %s", userID),
	)
}

func NewNotFoundSessionError(sessionID string) *DomainError {
	return newDomainError(
		buildNotFoundSessionError,
		fmt.Sprintf("Session not found: %s", sessionID),
	)
}

func NewDataAccessFailureError(resource string) *DomainError {
	return newDomainError(
		buildDataAccessFailureError,
		fmt.Sprintf("Failed to access data: %s", resource),
	)
}

func NewDataPersistFailureError(resource string) *DomainError {
	return newDomainError(
		buildDataPersistFailureError,
		fmt.Sprintf("Failed to persist data: %s", resource),
	)
}

func NewDataConflictError(resource string, constraint string) *DomainError {
	return newDomainError(
		buildDataConflictError,
		fmt.Sprintf("Data conflict on %s: %s", resource, constraint),
	)
}

// ### Unexpected Errors ######################################################

func NewUnexpectedError(detail string) *DomainError {
	return newDomainError(
		buildUnexpectedError,
		fmt.Sprintf("Unexpected error occurred: %s", detail),
	)
}
