package value_object

import "testing"

func TestNewDocumentType(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    DocumentType
		wantErr bool
	}{
		{
			name:    "valid procedure type",
			input:   "procedure",
			want:    DocumentTypeProcedure,
			wantErr: false,
		},
		{
			name:    "valid knowledge type",
			input:   "knowledge",
			want:    DocumentTypeKnowledge,
			wantErr: false,
		},
		{
			name:    "invalid type",
			input:   "invalid",
			want:    "",
			wantErr: true,
		},
		{
			name:    "empty string",
			input:   "",
			want:    "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewDocumentType(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewDocumentType() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("NewDocumentType() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDocumentType_IsValid(t *testing.T) {
	tests := []struct {
		name     string
		docType  DocumentType
		expected bool
	}{
		{"procedure is valid", DocumentTypeProcedure, true},
		{"knowledge is valid", DocumentTypeKnowledge, true},
		{"invalid type", DocumentType("invalid"), false},
		{"empty type", DocumentType(""), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.docType.IsValid(); got != tt.expected {
				t.Errorf("IsValid() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestDocumentType_IsProcedure(t *testing.T) {
	if !DocumentTypeProcedure.IsProcedure() {
		t.Error("IsProcedure() returned false for procedure type")
	}
	if DocumentTypeKnowledge.IsProcedure() {
		t.Error("IsProcedure() returned true for knowledge type")
	}
}

func TestDocumentType_IsKnowledge(t *testing.T) {
	if !DocumentTypeKnowledge.IsKnowledge() {
		t.Error("IsKnowledge() returned false for knowledge type")
	}
	if DocumentTypeProcedure.IsKnowledge() {
		t.Error("IsKnowledge() returned true for procedure type")
	}
}
