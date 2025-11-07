package dto

// ErrorResponse represents a generic error response
type ErrorResponse struct {
	Code    string `json:"code" example:"INVALID_REQUEST"`
	Message string `json:"message" example:"Invalid request body"`
}
