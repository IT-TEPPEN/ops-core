package security

import (
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
)

const (
	testSecretKey = "test-secret-key-for-jwt-provider-test"
	testIssuer    = "opscore-test"
)

func TestJWTProvider_GenerateAccessToken_Success(t *testing.T) {
	// Setup
	provider := NewJWTProviderAdapter(testSecretKey, testIssuer)
	userID := "user-123"
	jti := "session-456"
	expiresIn := 15 * time.Minute

	// Execute
	token, err := provider.GenerateAccessToken(userID, jti, expiresIn)

	// Assert
	assert.NoError(t, err)
	assert.NotEmpty(t, token)
}

func TestJWTProvider_ValidateAccessToken_Success(t *testing.T) {
	// Setup
	provider := NewJWTProviderAdapter(testSecretKey, testIssuer)
	userID := "user-123"
	jti := "session-456"
	expiresIn := 15 * time.Minute

	// Generate a token
	token, err := provider.GenerateAccessToken(userID, jti, expiresIn)
	assert.NoError(t, err)

	// Execute
	claims, err := provider.ValidateAccessToken(token)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, claims)
	assert.Equal(t, userID, claims.UserID)
	assert.Equal(t, jti, claims.JTI)
	assert.False(t, claims.IssuedAt.IsZero())
	assert.False(t, claims.ExpiresAt.IsZero())
	assert.True(t, claims.ExpiresAt.After(claims.IssuedAt))
}

func TestJWTProvider_ValidateAccessToken_InvalidToken(t *testing.T) {
	// Setup
	provider := NewJWTProviderAdapter(testSecretKey, testIssuer)
	invalidToken := "invalid.jwt.token"

	// Execute
	claims, err := provider.ValidateAccessToken(invalidToken)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, claims)
}

func TestJWTProvider_ValidateAccessToken_MissingRequiredClaims(t *testing.T) {
	// Setup - create a token manually without required claims
	provider := NewJWTProviderAdapter(testSecretKey, testIssuer).(*JWTProviderAdapter)

	// Create token with missing 'sub' claim
	claims := jwt.MapClaims{
		"jti": "session-123",
		"iss": testIssuer,
		"iat": time.Now().Unix(),
		"exp": time.Now().Add(15 * time.Minute).Unix(),
		// "sub" is missing!
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(provider.secretKey)
	assert.NoError(t, err)

	// Execute
	result, err := provider.ValidateAccessToken(tokenString)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "sub")
}

func TestJWTProvider_ValidateAccessToken_ExpiredToken(t *testing.T) {
	// Setup
	provider := NewJWTProviderAdapter(testSecretKey, testIssuer)
	userID := "user-123"
	jti := "session-456"
	expiresIn := 1 * time.Millisecond // Very short expiration

	// Generate a token
	token, err := provider.GenerateAccessToken(userID, jti, expiresIn)
	assert.NoError(t, err)

	// Wait for token to expire
	time.Sleep(10 * time.Millisecond)

	// Execute
	claims, err := provider.ValidateAccessToken(token)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, claims)
	assert.Contains(t, err.Error(), "expired")
}

func TestJWTProvider_ValidateAccessToken_WrongSecretKey(t *testing.T) {
	// Setup
	providerGenerate := NewJWTProviderAdapter(testSecretKey, testIssuer)
	providerValidate := NewJWTProviderAdapter("different-secret-key", testIssuer)
	userID := "user-123"
	jti := "session-456"
	expiresIn := 15 * time.Minute

	// Generate a token with first provider
	token, err := providerGenerate.GenerateAccessToken(userID, jti, expiresIn)
	assert.NoError(t, err)

	// Execute - try to validate with different secret key
	claims, err := providerValidate.ValidateAccessToken(token)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, claims)
}
