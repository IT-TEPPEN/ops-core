package domain

import "github.com/golang-jwt/jwt/v5"

// JWTClaims represents the claims stored in JWT tokens
type JWTClaims struct {
	UserID string `json:"sub"`
	Email  string `json:"email"`
	Name   string `json:"name"`
	jwt.RegisteredClaims
}
