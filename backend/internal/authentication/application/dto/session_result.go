package dto

// SessionResult represents the result of session issuance from provider authentication
type SessionResult struct {
	// AccessToken is the JWT access token for API authentication
	AccessToken string
	// RefreshToken is the refresh token for obtaining new access tokens
	RefreshToken string
	// User is the authenticated user's basic information
	User UserInfo
}
