package dto

// CreateSessionRequest represents the request to create a new session.
type CreateSessionRequest struct {
	UserID      string      `json:"user_id"`
	SessionType SessionType `json:"session_type"`
}
