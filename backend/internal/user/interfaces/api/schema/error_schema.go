package schema

// ErrorResponse represents the API error response format
type ErrorResponse struct {
	Code    string `json:"code" example:"INVALID_REQUEST"`
	Message string `json:"message" example:"Invalid request body"`
}
