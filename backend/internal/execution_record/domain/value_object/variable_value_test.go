package value_object

import "testing"

func TestNewVariableValue(t *testing.T) {
	tests := []struct {
		name    string
		varName string
		value   interface{}
		wantErr bool
	}{
		{
			name:    "valid string value",
			varName: "server_name",
			value:   "production-server",
			wantErr: false,
		},
		{
			name:    "valid number value",
			varName: "port",
			value:   8080,
			wantErr: false,
		},
		{
			name:    "valid boolean value",
			varName: "enabled",
			value:   true,
			wantErr: false,
		},
		{
			name:    "valid nil value",
			varName: "optional",
			value:   nil,
			wantErr: false,
		},
		{
			name:    "empty variable name",
			varName: "",
			value:   "value",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewVariableValue(tt.varName, tt.value)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewVariableValue() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if got.Name() != tt.varName {
					t.Errorf("VariableValue.Name() = %v, want %v", got.Name(), tt.varName)
				}
				if got.Value() != tt.value {
					t.Errorf("VariableValue.Value() = %v, want %v", got.Value(), tt.value)
				}
			}
		})
	}
}

func TestReconstructVariableValue(t *testing.T) {
	name := "server"
	value := "localhost"

	vv := ReconstructVariableValue(name, value)

	if vv.Name() != name {
		t.Errorf("ReconstructVariableValue().Name() = %v, want %v", vv.Name(), name)
	}
	if vv.Value() != value {
		t.Errorf("ReconstructVariableValue().Value() = %v, want %v", vv.Value(), value)
	}
}

func TestVariableValue_Equals(t *testing.T) {
	vv1, _ := NewVariableValue("name1", "value1")
	vv2, _ := NewVariableValue("name2", "value2")
	vv1Copy, _ := NewVariableValue("name1", "different_value")

	if !vv1.Equals(vv1Copy) {
		t.Error("VariableValues with same name should be equal")
	}
	if vv1.Equals(vv2) {
		t.Error("VariableValues with different names should not be equal")
	}
}
