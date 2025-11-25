package value_object

import "testing"

func TestNewAccessScope(t *testing.T) {
	tests := []struct {
		name    string
		scope   string
		wantErr bool
	}{
		{
			name:    "valid public",
			scope:   "public",
			wantErr: false,
		},
		{
			name:    "valid private",
			scope:   "private",
			wantErr: false,
		},
		{
			name:    "invalid scope",
			scope:   "invalid",
			wantErr: true,
		},
		{
			name:    "empty scope",
			scope:   "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewAccessScope(tt.scope)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewAccessScope() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got.String() != tt.scope {
				t.Errorf("NewAccessScope() = %v, want %v", got, tt.scope)
			}
		})
	}
}

func TestAccessScope_IsValid(t *testing.T) {
	tests := []struct {
		name  string
		scope AccessScope
		want  bool
	}{
		{
			name:  "public is valid",
			scope: AccessScopePublic,
			want:  true,
		},
		{
			name:  "private is valid",
			scope: AccessScopePrivate,
			want:  true,
		},
		{
			name:  "invalid scope",
			scope: AccessScope("invalid"),
			want:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.scope.IsValid(); got != tt.want {
				t.Errorf("AccessScope.IsValid() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAccessScope_IsPublic(t *testing.T) {
	tests := []struct {
		name  string
		scope AccessScope
		want  bool
	}{
		{
			name:  "public returns true",
			scope: AccessScopePublic,
			want:  true,
		},
		{
			name:  "private returns false",
			scope: AccessScopePrivate,
			want:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.scope.IsPublic(); got != tt.want {
				t.Errorf("AccessScope.IsPublic() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAccessScope_IsPrivate(t *testing.T) {
	tests := []struct {
		name  string
		scope AccessScope
		want  bool
	}{
		{
			name:  "private returns true",
			scope: AccessScopePrivate,
			want:  true,
		},
		{
			name:  "public returns false",
			scope: AccessScopePublic,
			want:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.scope.IsPrivate(); got != tt.want {
				t.Errorf("AccessScope.IsPrivate() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAccessScope_Equals(t *testing.T) {
	scope1, _ := NewAccessScope("public")
	scope2, _ := NewAccessScope("private")
	scope1Copy, _ := NewAccessScope("public")

	if !scope1.Equals(scope1Copy) {
		t.Error("Same access scopes should be equal")
	}
	if scope1.Equals(scope2) {
		t.Error("Different access scopes should not be equal")
	}
}
