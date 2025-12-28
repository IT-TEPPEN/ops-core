package service

import (
	"context"
	"fmt"
	"time"

	"opscore/backend/internal/auth/domain"

	"github.com/golang-jwt/jwt/v5"
)

// JWTService handles JWT token generation and validation
type JWTService struct {
	secret     []byte
	issuer     string
	expiration time.Duration
}

// NewJWTService creates a new JWTService
func NewJWTService(secret string, issuer string, expiration time.Duration) *JWTService {
	return &JWTService{
		secret:     []byte(secret),
		issuer:     issuer,
		expiration: expiration,
	}
}

// GenerateToken generates a new JWT token for the user
func (s *JWTService) GenerateToken(ctx context.Context, user *domain.User) (string, error) {
	now := time.Now()
	claims := domain.JWTClaims{
		UserID: user.ID,
		Email:  user.PrimaryEmail,
		Name:   user.DisplayName,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    s.issuer,
			Subject:   user.ID,
			ExpiresAt: jwt.NewNumericDate(now.Add(s.expiration)),
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
