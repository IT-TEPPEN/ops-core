package usecase

import (
	"context"
	"fmt"
	"regexp"

	apperror "opscore/backend/internal/document/application/error"
	"opscore/backend/internal/document/application/dto"
	"opscore/backend/internal/document/domain/repository"
	"opscore/backend/internal/document/domain/value_object"
)

// VariableValue represents a variable name-value pair
type VariableValue struct {
	Name  string
	Value interface{}
}

// VariableUseCase defines the interface for variable-related use cases
type VariableUseCase interface {
	// GetVariableDefinitions retrieves variable definitions for a document
	GetVariableDefinitions(ctx context.Context, documentID string) ([]dto.VariableDefinitionDTO, error)

	// ValidateVariableValues validates variable values against their definitions
	ValidateVariableValues(ctx context.Context, documentID string, values []VariableValue) error

	// SubstituteVariables replaces variable placeholders in content with provided values
	SubstituteVariables(ctx context.Context, content string, values []VariableValue) (string, error)
}

// variableUseCase implements the VariableUseCase interface
type variableUseCase struct {
	docRepo repository.DocumentRepository
}

// NewVariableUseCase creates a new instance of variableUseCase
func NewVariableUseCase(docRepo repository.DocumentRepository) VariableUseCase {
	return &variableUseCase{
		docRepo: docRepo,
	}
}

// GetVariableDefinitions retrieves variable definitions for a document
func (uc *variableUseCase) GetVariableDefinitions(ctx context.Context, documentID string) ([]dto.VariableDefinitionDTO, error) {
	// Validate document ID
	docID, err := value_object.NewDocumentID(documentID)
	if err != nil {
		return nil, apperror.NewValidationFailedError([]apperror.FieldError{
			{Field: "document_id", Message: err.Error()},
		})
	}

	// Find the document
	doc, err := uc.docRepo.FindByID(ctx, docID)
	if err != nil {
		return nil, fmt.Errorf("failed to find document: %w", err)
	}
	if doc == nil {
		return nil, apperror.NewNotFoundError("Document", documentID, nil)
	}

	// Get current version
	currentVersion := doc.CurrentVersion()
	if currentVersion == nil {
		return []dto.VariableDefinitionDTO{}, nil
	}

	// Convert variables to DTOs
	variables := currentVersion.Variables()
	variableDTOs := make([]dto.VariableDefinitionDTO, len(variables))
	for i, v := range variables {
		variableDTOs[i] = dto.VariableDefinitionDTO{
			Name:         v.Name(),
			Label:        v.Label(),
			Description:  v.Description(),
			Type:         v.Type().String(),
			Required:     v.Required(),
			DefaultValue: v.DefaultValue(),
		}
	}

	return variableDTOs, nil
}

// ValidateVariableValues validates variable values against their definitions
func (uc *variableUseCase) ValidateVariableValues(ctx context.Context, documentID string, values []VariableValue) error {
	// Get variable definitions
	definitions, err := uc.GetVariableDefinitions(ctx, documentID)
	if err != nil {
		return err
	}

	// Create a map of definitions for quick lookup
	defMap := make(map[string]dto.VariableDefinitionDTO)
	for _, def := range definitions {
		defMap[def.Name] = def
	}

	// Create a map of values for quick lookup
	valueMap := make(map[string]interface{})
	for _, val := range values {
		valueMap[val.Name] = val.Value
	}

	// Validate required fields
	var fieldErrors []apperror.FieldError
	for _, def := range definitions {
		if def.Required {
			val, exists := valueMap[def.Name]
			empty := false
			if !exists || val == nil {
				empty = true
			} else {
				switch def.Type {
				case "string", "date":
					strVal, ok := val.(string)
					if ok && strVal == "" {
						empty = true
					}
				// For number and boolean, nil or missing is empty, 0/false is valid
				}
			}
			if empty {
				fieldErrors = append(fieldErrors, apperror.FieldError{
					Field:   def.Name,
					Message: fmt.Sprintf("%s is required", def.Label),
				})
			}
		}
	}

	// Validate variable types (basic type checking)
	for _, val := range values {
		def, exists := defMap[val.Name]
		if !exists {
			// Unknown variable - skip or report error
			continue
		}

		// Type validation
		if val.Value != nil && val.Value != "" {
			switch def.Type {
			case "number":
				// Check if value can be converted to number
				switch val.Value.(type) {
				case int, int32, int64, float32, float64:
					// Valid number type
				default:
					fieldErrors = append(fieldErrors, apperror.FieldError{
						Field:   def.Name,
						Message: fmt.Sprintf("%s must be a number", def.Label),
					})
				}
			case "boolean":
				// Check if value is boolean
				switch val.Value.(type) {
				case bool:
					// Valid boolean type
				default:
					fieldErrors = append(fieldErrors, apperror.FieldError{
						Field:   def.Name,
						Message: fmt.Sprintf("%s must be a boolean", def.Label),
					})
				}
			}
		}
	}

	if len(fieldErrors) > 0 {
		return apperror.NewValidationFailedError(fieldErrors)
	}

	return nil
}

// SubstituteVariables replaces variable placeholders in content with provided values
func (uc *variableUseCase) SubstituteVariables(ctx context.Context, content string, values []VariableValue) (string, error) {
	result := content

	// Replace each variable in the content
	// Precompile regexes for all variable names
	regexMap := make(map[string]*regexp.Regexp, len(values))
	for _, val := range values {
		pattern := fmt.Sprintf(`\{\{%s\}\}`, regexp.QuoteMeta(val.Name))
		re, err := regexp.Compile(pattern)
		if err != nil {
			return "", fmt.Errorf("failed to compile regex for variable %s: %w", val.Name, err)
		}
		regexMap[val.Name] = re
	}

	// Replace each variable in the content using precompiled regexes
	for _, val := range values {
		strValue := ""
		if val.Value != nil {
			strValue = fmt.Sprintf("%v", val.Value)
		}
		result = regexMap[val.Name].ReplaceAllString(result, strValue)
	}

	return result, nil
}
