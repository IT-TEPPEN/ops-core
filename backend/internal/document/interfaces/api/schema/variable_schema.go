package schema

// GetVariableDefinitionsResponse represents the response for getting variable definitions
type GetVariableDefinitionsResponse struct {
	DocumentID string                  `json:"documentId"`
	Variables  []VariableDefinitionDTO `json:"variables"`
}

// VariableDefinitionDTO represents a variable definition in the API schema
type VariableDefinitionDTO struct {
	Name         string      `json:"name"`
	Label        string      `json:"label"`
	Description  string      `json:"description,omitempty"`
	Type         string      `json:"type"` // "string", "number", "boolean", "date"
	Required     bool        `json:"required"`
	DefaultValue interface{} `json:"defaultValue,omitempty"`
}

// VariableValueDTO represents a variable value in the API schema
type VariableValueDTO struct {
	Name  string      `json:"name"`
	Value interface{} `json:"value"`
}

// ValidateVariableValuesRequest represents the request for validating variable values
type ValidateVariableValuesRequest struct {
	Values []VariableValueDTO `json:"values"`
}

// ValidateVariableValuesResponse represents the response for validating variable values
// 
// Design Note: This endpoint returns HTTP 200 OK for both successful and failed validations.
// The "valid" field indicates the validation result:
//   - valid: true  - All variable values passed validation
//   - valid: false - One or more validation errors occurred (see "errors" field)
// 
// This differs from typical REST conventions where validation failures return 4xx codes.
// The rationale is to distinguish between:
//   - Validation failures (HTTP 200, valid:false) - Expected client behavior, values don't meet requirements
//   - System errors (HTTP 4xx/5xx) - Malformed requests or server issues
type ValidateVariableValuesResponse struct {
	Valid  bool                `json:"valid"`
	Errors []ValidationErrorDTO `json:"errors,omitempty"`
}

// ValidationErrorDTO represents a validation error for a variable
// Example error message for a required field: "{Label} is required"
// where {Label} is the variable's label field (e.g., "Server Name is required")
type ValidationErrorDTO struct {
	Name    string `json:"name"`
	Message string `json:"message"`
}
