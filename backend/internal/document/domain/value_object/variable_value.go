package value_object

import "errors"

// VariableValue represents a value assigned to a variable.
type VariableValue struct {
	name  string
	value interface{}
}

// NewVariableValue creates a new VariableValue.
func NewVariableValue(name string, value interface{}) (VariableValue, error) {
	if name == "" {
		return VariableValue{}, errors.New("variable name cannot be empty")
	}
	// Value can be nil for optional variables
	return VariableValue{
		name:  name,
		value: value,
	}, nil
}

// Name returns the variable name.
func (v VariableValue) Name() string {
	return v.name
}

// Value returns the variable value.
func (v VariableValue) Value() interface{} {
	return v.value
}

// Equals checks if two VariableValues have the same name.
func (v VariableValue) Equals(other VariableValue) bool {
	return v.name == other.name
}
