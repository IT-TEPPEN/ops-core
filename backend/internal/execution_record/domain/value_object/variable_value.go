package value_object

import "errors"

// VariableValue represents a variable name and its value during execution.
type VariableValue struct {
	name  string
	value interface{}
}

// NewVariableValue creates a new VariableValue.
func NewVariableValue(name string, value interface{}) (VariableValue, error) {
	if name == "" {
		return VariableValue{}, errors.New("variable name cannot be empty")
	}
	return VariableValue{
		name:  name,
		value: value,
	}, nil
}

// ReconstructVariableValue reconstructs a VariableValue from persistence data.
func ReconstructVariableValue(name string, value interface{}) VariableValue {
	return VariableValue{
		name:  name,
		value: value,
	}
}

// Name returns the variable name.
func (v VariableValue) Name() string {
	return v.name
}

// Value returns the variable value.
func (v VariableValue) Value() interface{} {
	return v.value
}

// Equals checks if two VariableValues are equal by name.
func (v VariableValue) Equals(other VariableValue) bool {
	return v.name == other.name
}
