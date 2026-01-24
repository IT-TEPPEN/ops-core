package dto

// LinkIdentityRequest represents the request to link a new identity to a user.
type LinkIdentityRequest struct {
	UserID       string           `json:"user_id"`
	ProviderInfo ProviderUserInfo `json:"provider_info"`
}
