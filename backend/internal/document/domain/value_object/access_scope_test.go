package value_object

import "testing"

func TestNewAccessScope(t *testing.T) {
	tests := []struct {
		name    string
		scope   string
		want    AccessScope
		wantErr bool
	}{
		{
			name:    "public scope",
			scope:   "public",
			want:    AccessScopePublic,
			wantErr: false,
		},
		{
			name:    "private scope",
			scope:   "private",
			want:    AccessScopePrivate,
			wantErr: false,
		},
		{
			name:    "invalid scope",
			scope:   "invalid",
			want:    "",
			wantErr: true,
		},
		{
			name:    "group scope is invalid",
			scope:   "group",
			want:    "",
			wantErr: true,
		},
		{
			name:    "user scope is invalid",
			scope:   "user",
			want:    "",
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
			if got != tt.want {
				t.Errorf("NewAccessScope() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAccessScope_Checks(t *testing.T) {
	tests := []struct {
		name      string
		scope     AccessScope
		isPublic  bool
		isPrivate bool
	}{
		{
			name:      "public scope",
			scope:     AccessScopePublic,
			isPublic:  true,
			isPrivate: false,
		},
		{
			name:      "private scope",
			scope:     AccessScopePrivate,
			isPublic:  false,
			isPrivate: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.scope.IsPublic(); got != tt.isPublic {
				t.Errorf("IsPublic() = %v, want %v", got, tt.isPublic)
			}
			if got := tt.scope.IsPrivate(); got != tt.isPrivate {
				t.Errorf("IsPrivate() = %v, want %v", got, tt.isPrivate)
			}
		})
	}
}
