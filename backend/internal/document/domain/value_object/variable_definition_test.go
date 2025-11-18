package value_object

import "testing"

func TestNewVariableDefinition(t *testing.T) {
	tests := []struct {
		name         string
		varName      string
		label        string
		description  string
		varType      VariableType
		required     bool
		defaultValue interface{}
		wantErr      bool
	}{
		{
			name:         "valid string variable",
			varName:      "server_name",
			label:        "Server Name",
			description:  "Name of the server",
			varType:      VariableTypeString,
			required:     true,
			defaultValue: "localhost",
			wantErr:      false,
		},
		{
			name:         "valid number variable",
			varName:      "port",
			label:        "Port",
			description:  "Server port",
			varType:      VariableTypeNumber,
			required:     false,
			defaultValue: 8080,
			wantErr:      false,
		},
		{
			name:         "empty name",
			varName:      "",
			label:        "Label",
			description:  "Description",
			varType:      VariableTypeString,
			required:     false,
			defaultValue: nil,
			wantErr:      true,
		},
		{
			name:         "invalid name format",
			varName:      "123invalid",
			label:        "Label",
			description:  "Description",
			varType:      VariableTypeString,
			required:     false,
			defaultValue: nil,
			wantErr:      true,
		},
		{
			name:         "name with spaces",
			varName:      "invalid name",
			label:        "Label",
			description:  "Description",
			varType:      VariableTypeString,
			required:     false,
			defaultValue: nil,
			wantErr:      true,
		},
		{
			name:         "empty label",
			varName:      "valid_name",
			label:        "",
			description:  "Description",
			varType:      VariableTypeString,
			required:     false,
			defaultValue: nil,
			wantErr:      true,
		},
		{
			name:         "invalid type",
			varName:      "valid_name",
			label:        "Label",
			description:  "Description",
			varType:      VariableType("invalid"),
			required:     false,
			defaultValue: nil,
			wantErr:      true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewVariableDefinition(tt.varName, tt.label, tt.description, tt.varType, tt.required, tt.defaultValue)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewVariableDefinition() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if got.Name() != tt.varName {
					t.Errorf("Name() = %v, want %v", got.Name(), tt.varName)
				}
				if got.Label() != tt.label {
					t.Errorf("Label() = %v, want %v", got.Label(), tt.label)
				}
				if got.Type() != tt.varType {
					t.Errorf("Type() = %v, want %v", got.Type(), tt.varType)
				}
				if got.Required() != tt.required {
					t.Errorf("Required() = %v, want %v", got.Required(), tt.required)
				}
			}
		})
	}
}

func TestNewVariableType(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    VariableType
		wantErr bool
	}{
		{
			name:    "string type",
			input:   "string",
			want:    VariableTypeString,
			wantErr: false,
		},
		{
			name:    "number type",
			input:   "number",
			want:    VariableTypeNumber,
			wantErr: false,
		},
		{
			name:    "boolean type",
			input:   "boolean",
			want:    VariableTypeBoolean,
			wantErr: false,
		},
		{
			name:    "date type",
			input:   "date",
			want:    VariableTypeDate,
			wantErr: false,
		},
		{
			name:    "invalid type",
			input:   "invalid",
			want:    "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewVariableType(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewVariableType() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("NewVariableType() = %v, want %v", got, tt.want)
			}
		})
	}
}
