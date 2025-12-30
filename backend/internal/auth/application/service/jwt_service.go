package service

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"time"

	"opscore/backend/internal/auth/domain"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

// JWTService handles JWT token generation and validation
type JWTService struct {
	secret                   []byte
	issuer                   string
	accessTokenExpiration    time.Duration
	refreshTokenExpiration   time.Duration
	refreshTokenExpirationRM time.Duration // Remember Me
}

// NewJWTService creates a new JWTService
func NewJWTService(secret string, issuer string, accessTokenExpiration time.Duration) *JWTService {
	return &JWTService{
		secret:                   []byte(secret),
		issuer:                   issuer,
		accessTokenExpiration:    accessTokenExpiration,
		refreshTokenExpiration:   7 * 24 * time.Hour,  // 7 days
		refreshTokenExpirationRM: 30 * 24 * time.Hour, // 30 days (Remember Me)
	}
}

// GenerateToken generates a new JWT access token for the user
func (s *JWTService) GenerateToken(ctx context.Context, user *domain.User) (string, error) {
	now := time.Now()
	claims := domain.JWTClaims{
		UserID: user.ID,
		Email:  user.PrimaryEmail,
		Name:   user.DisplayName,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    s.issuer,
			Subject:   user.ID,
			ExpiresAt: jwt.NewNumericDate(now.Add(s.accessTokenExpiration)),
			IssuedAt:  jwt.NewNumericDate(now),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString(s.secret)
	if err != nil {
		return "", fmt.Errorf("failed to sign token: %w", err)
	}

	return signedToken, nil
}

// TokenPair represents an access token and refresh token pair
type TokenPair struct {
	AccessToken  string
	RefreshToken string
	ExpiresAt    time.Time
}

// GenerateTokenPair generates both access and refresh tokens
func (s *JWTService) GenerateTokenPair(ctx context.Context, user *domain.User, rememberMe bool) (*TokenPair, *domain.RefreshToken, error) {
	// Generate access token
	accessToken, err := s.GenerateToken(ctx, user)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to generate access token: %w", err)
	}

	// Generate refresh token (random string)
	refreshTokenString, err := generateRandomToken()
	if err != nil {
		return nil, nil, fmt.Errorf("failed to generate refresh token: %w", err)
	}

	// Hash the refresh token for storage
	tokenHash := hashToken(refreshTokenString)

	// Generate JTI (JWT ID)
	jti := uuid.New().String()

	// Calculate expiration based on rememberMe flag
	var expiresAt time.Time
	if rememberMe {
		expiresAt = time.Now().Add(s.refreshTokenExpirationRM)
	} else {
		expiresAt = time.Now().Add(s.refreshTokenExpiration)
	}

	// Create refresh token domain object
	refreshTokenID := uuid.New().String()
	refreshToken := domain.NewRefreshToken(refreshTokenID, user.ID, tokenHash, jti, expiresAt, rememberMe)

	return &TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshTokenString, // Return plain token to client
		ExpiresAt:    expiresAt,
	}, refreshToken, nil
}

// ValidateRefreshToken validates a refresh token string and returns the hash
func (s *JWTService) ValidateRefreshToken(tokenString string) (string, error) {
	if tokenString == "" {
		return "", fmt.Errorf("empty refresh token")
	}

	// Hash the token to compare with database
	tokenHash := hashToken(tokenString)
	return tokenHash, nil
}

// generateRandomToken generates a cryptographically secure random token
func generateRandomToken() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", fmt.Errorf("failed to generate random bytes: %w", err)
	}
	return base64.URLEncoding.EncodeToString(b), nil
}

// hashToken hashes a token using SHA-256
func hashToken(token string) string {
	hash := sha256.Sum256([]byte(token))
	return hex.EncodeToString(hash[:])
}

// ValidateToken validates a JWT token and returns the claims
func (s *JWTService) ValidateToken(ctx context.Context, tokenString string) (*domain.JWTClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &domain.JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		// Verify signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return s.secret, nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to parse token: %w", err)
	}

	claims, ok := token.Claims.(*domain.JWTClaims)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("invalid token claims")
	}

	return claims, nil
}
