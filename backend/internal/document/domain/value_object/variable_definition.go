package value_object

import (
	"errors"
	"regexp"
)

// VariableType represents the data type of a variable.
type VariableType string

const (
	// VariableTypeString represents a string variable type.
	VariableTypeString VariableType = "string"
	// VariableTypeNumber represents a number variable type.
	VariableTypeNumber VariableType = "number"
	// VariableTypeBoolean represents a boolean variable type.
	VariableTypeBoolean VariableType = "boolean"
	// VariableTypeDate represents a date variable type.
	VariableTypeDate VariableType = "date"
)

// NewVariableType creates a new VariableType from a string.
func NewVariableType(t string) (VariableType, error) {
	varType := VariableType(t)
	if !varType.IsValid() {
		return "", errors.New("invalid variable type: must be 'string', 'number', 'boolean', or 'date'")
	}
	return varType, nil
}

// IsValid checks if the VariableType is valid.
func (v VariableType) IsValid() bool {
	return v == VariableTypeString || v == VariableTypeNumber || v == VariableTypeBoolean || v == VariableTypeDate
}

// String returns the string representation of VariableType.
func (v VariableType) String() string {
	return string(v)
}

// VariableDefinition represents a variable definition in a document.
type VariableDefinition struct {
	name         string
	label        string
	description  string
	varType      VariableType
	required     bool
	defaultValue interface{}
}

var variableNameRegex = regexp.MustCompile(`^[a-zA-Z][a-zA-Z0-9_]*$`)

// NewVariableDefinition creates a new VariableDefinition.
func NewVariableDefinition(
	name string,
	label string,
	description string,
	varType VariableType,
	required bool,
	defaultValue interface{},
) (VariableDefinition, error) {
	if name == "" {
		return VariableDefinition{}, errors.New("variable name cannot be empty")
	}
	if !variableNameRegex.MatchString(name) {
		return VariableDefinition{}, errors.New("variable name must be alphanumeric with underscores, starting with a letter")
	}
	if label == "" {
		return VariableDefinition{}, errors.New("variable label cannot be empty")
	}
	if !varType.IsValid() {
		return VariableDefinition{}, errors.New("invalid variable type")
	}

	return VariableDefinition{
		name:         name,
		label:        label,
		description:  description,
		varType:      varType,
		required:     required,
		defaultValue: defaultValue,
	}, nil
}

// Name returns the variable name.
func (v VariableDefinition) Name() string {
	return v.name
}

// Label returns the variable label.
func (v VariableDefinition) Label() string {
	return v.label
}

// Description returns the variable description.
func (v VariableDefinition) Description() string {
	return v.description
}

// Type returns the variable type.
func (v VariableDefinition) Type() VariableType {
	return v.varType
}

// Required returns whether the variable is required.
func (v VariableDefinition) Required() bool {
	return v.required
}

// DefaultValue returns the default value.
func (v VariableDefinition) DefaultValue() interface{} {
	return v.defaultValue
}

// Equals checks if two VariableDefinitions are equal by name.
func (v VariableDefinition) Equals(other VariableDefinition) bool {
	return v.name == other.name
}
