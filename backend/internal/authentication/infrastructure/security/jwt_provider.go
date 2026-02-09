package security

import (
	"errors"
	"time"

	"opscore/backend/internal/authentication/application/port"

	"github.com/golang-jwt/jwt/v5"
)

// JWTProviderAdapter implements port.JWTProvider using golang-jwt library.
// This is an Adapter in hexagonal architecture - it adapts the external JWT library
// to the interface expected by the application layer.
type JWTProviderAdapter struct {
	secretKey []byte
	issuer    string
}

// NewJWTProviderAdapter creates a new JWT provider adapter instance.
//
// Parameters:
//   - secretKey: The secret key for signing JWT tokens
//   - issuer: The issuer claim to include in tokens (typically the application name)
//
// Returns:
//   - port.JWTProvider: The adapter implementing the JWTProvider interface
func NewJWTProviderAdapter(secretKey string, issuer string) port.JWTProvider {
	return &JWTProviderAdapter{
		secretKey: []byte(secretKey),
		issuer:    issuer,
	}
}

// GenerateAccessToken generates a JWT access token for a user
func (j *JWTProviderAdapter) GenerateAccessToken(userID string, jti string, expiresIn time.Duration) (string, error) {
	now := time.Now()
	claims := jwt.MapClaims{
		"sub": userID,                    // Subject: user ID
		"jti": jti,                       // JWT ID: for token invalidation
		"iss": j.issuer,                  // Issuer: application name
		"iat": now.Unix(),                // Issued At
		"exp": now.Add(expiresIn).Unix(), // Expiration Time
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString(j.secretKey)
	if err != nil {
		return "", err
	}

	return signedToken, nil
}

// ValidateAccessToken validates a JWT access token and returns the claims.
// This adapter extracts JWT standard claims and maps them to Application-defined JWTClaims.
func (j *JWTProviderAdapter) ValidateAccessToken(tokenString string) (*port.JWTClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Verify signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return j.secretKey, nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, errors.New("invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("invalid token claims")
	}

	// Extract and validate required claims as per Application's contract
	userID, ok := claims["sub"].(string)
	if !ok || userID == "" {
		return nil, errors.New("missing or invalid 'sub' claim")
	}

	jti, ok := claims["jti"].(string)
	if !ok || jti == "" {
		return nil, errors.New("missing or invalid 'jti' claim")
	}

	// Extract timestamps
	var issuedAt, expiresAt time.Time
	if iat, ok := claims["iat"].(float64); ok {
		issuedAt = time.Unix(int64(iat), 0)
	}
	if exp, ok := claims["exp"].(float64); ok {
		expiresAt = time.Unix(int64(exp), 0)
	}

	// Map Infrastructure JWT claims to Application-defined structure
	return &port.JWTClaims{
		UserID:    userID,
		JTI:       jti,
		IssuedAt:  issuedAt,
		ExpiresAt: expiresAt,
	}, nil
}
