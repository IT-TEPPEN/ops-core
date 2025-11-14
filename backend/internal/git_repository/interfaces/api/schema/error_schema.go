package schema

// ErrorResponse represents a generic API error response
type ErrorResponse struct {
	Code    string `json:"code" example:"INVALID_REQUEST"`
	Message string `json:"message" example:"Invalid request body"`
}
