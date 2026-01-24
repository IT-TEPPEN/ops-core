package service

import (
	"context"
	"time"

	"opscore/backend/internal/authentication/domain/value_object"
)

// TokenService defines the interface for JWT token management.
// This service is responsible for generating and validating access tokens (JWTs).
// Implementation should use industry-standard JWT libraries and follow OAuth 2.0 / OIDC practices.
type TokenService interface {
	// GenerateAccessToken generates a JWT access token for the authenticated user.
	// The token should include standard claims (sub, iss, aud, exp, iat, jti) and custom claims as needed.
	//
	// Parameters:
	//   - ctx context.Context: The context for cancellation and timeout
	//   - userID value_object.UserID: The user ID to include in the token (sub claim)
	//   - sessionJTI string: The session's JWT ID for token binding
	//   - expirationDuration time.Duration: How long the token should be valid
	//
	// Returns:
	//   - string: The generated JWT token (signed)
	//   - error: Application error if generation fails
	//
	// Errors:
	//   - application_err.ErrTokenGenerationError: When token generation fails
	GenerateAccessToken(ctx context.Context, userID value_object.UserID, sessionJTI string, expirationDuration time.Duration) (string, error)

	// ValidateAccessToken validates a JWT access token and extracts its claims.
	// This checks the signature, expiration, and other standard validations.
	//
	// Parameters:
	//   - ctx context.Context: The context for cancellation and timeout
	//   - token string: The JWT token to validate
	//
	// Returns:
	//   - *TokenClaims: Parsed token claims (userID, sessionJTI, etc.)
	//   - error: Application error if validation fails
	//
	// Errors:
	//   - application_err.ErrTokenValidationError: When token is invalid, expired, or malformed
	ValidateAccessToken(ctx context.Context, token string) (*TokenClaims, error)
}

// TokenClaims represents the parsed claims from a JWT token.
type TokenClaims struct {
	// UserID is the user identifier (sub claim)
	UserID value_object.UserID
	// SessionJTI is the session's JWT ID for token binding
	SessionJTI string
	// IssuedAt is when the token was issued
	IssuedAt time.Time
	// ExpiresAt is when the token expires
	ExpiresAt time.Time
}
