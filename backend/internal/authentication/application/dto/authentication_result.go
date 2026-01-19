package dto

// AuthenticationResult represents the result of user authentication
type AuthenticationResult struct {
	// User is the authenticated user's information
	User UserInfo
	// Identity is the identity used for authentication
	Identity IdentityInfo
	// IsNewUser indicates if this is a newly created user
	IsNewUser bool
}
