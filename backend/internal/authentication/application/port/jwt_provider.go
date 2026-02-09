package port

import "time"

// JWTClaims represents the claims extracted from a validated JWT access token.
// This defines what the Application layer expects from JWT validation,
// following the Dependency Inversion Principle - Application layer defines
// the contract, Infrastructure layer implements it.
type JWTClaims struct {
	// UserID is the authenticated user's identifier
	UserID string
	// JTI is the JWT ID used for token invalidation
	JTI string
	// IssuedAt is when the token was issued
	IssuedAt time.Time
	// ExpiresAt is when the token expires
	ExpiresAt time.Time
}

// JWTProvider defines the interface for JWT token generation and validation.
// This is a Port in hexagonal architecture - an interface that the application layer
// expects external adapters (infrastructure) to implement.
//
// The application layer depends on this interface, while infrastructure provides
// the concrete implementation, following the Dependency Inversion Principle.
type JWTProvider interface {
	// GenerateAccessToken generates a JWT access token for a user.
	//
	// Parameters:
	//   - userID: The user ID to include in the token
	//   - jti: JWT ID to include in the token (for token invalidation)
	//   - expiresIn: Duration until the token expires
	//
	// Returns:
	//   - string: The encoded JWT token
	//   - error: Error if token generation fails
	GenerateAccessToken(userID string, jti string, expiresIn time.Duration) (string, error)

	// ValidateAccessToken validates a JWT access token and returns the claims.
	//
	// This method defines the contract: Infrastructure must extract these specific
	// claims and return them in the Application-defined structure, not a generic map.
	//
	// Parameters:
	//   - token: The JWT token to validate
	//
	// Returns:
	//   - *JWTClaims: The validated claims in Application-defined structure
	//   - error: Error if token is invalid, expired, or missing required claims
	ValidateAccessToken(token string) (*JWTClaims, error)
}
