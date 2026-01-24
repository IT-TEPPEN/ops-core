package dto

// SetPrimaryIdentityRequest represents the request to set a primary identity.
type SetPrimaryIdentityRequest struct {
	UserID     string `json:"user_id"`
	IdentityID string `json:"identity_id"`
}
