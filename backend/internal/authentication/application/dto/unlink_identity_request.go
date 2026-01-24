package dto

// UnlinkIdentityRequest represents the request to unlink an identity from a user.
type UnlinkIdentityRequest struct {
	UserID     string `json:"user_id"`
	IdentityID string `json:"identity_id"`
}
