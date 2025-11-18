package value_object

import "testing"

func TestNewVersionNumber(t *testing.T) {
	tests := []struct {
		name    string
		num     int
		want    int
		wantErr bool
	}{
		{
			name:    "valid version 1",
			num:     1,
			want:    1,
			wantErr: false,
		},
		{
			name:    "valid version 10",
			num:     10,
			want:    10,
			wantErr: false,
		},
		{
			name:    "zero version",
			num:     0,
			want:    0,
			wantErr: true,
		},
		{
			name:    "negative version",
			num:     -1,
			want:    0,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewVersionNumber(tt.num)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewVersionNumber() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got.Int() != tt.want {
				t.Errorf("NewVersionNumber() = %v, want %v", got.Int(), tt.want)
			}
		})
	}
}

func TestVersionNumber_Next(t *testing.T) {
	v1, _ := NewVersionNumber(1)
	v2 := v1.Next()

	if v2.Int() != 2 {
		t.Errorf("Next() = %v, want 2", v2.Int())
	}

	v10, _ := NewVersionNumber(10)
	v11 := v10.Next()

	if v11.Int() != 11 {
		t.Errorf("Next() = %v, want 11", v11.Int())
	}
}

func TestVersionNumber_Previous(t *testing.T) {
	tests := []struct {
		name    string
		num     int
		want    int
		wantErr bool
	}{
		{
			name:    "version 2 to 1",
			num:     2,
			want:    1,
			wantErr: false,
		},
		{
			name:    "version 10 to 9",
			num:     10,
			want:    9,
			wantErr: false,
		},
		{
			name:    "version 1 has no previous",
			num:     1,
			want:    0,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v, _ := NewVersionNumber(tt.num)
			got, err := v.Previous()
			if (err != nil) != tt.wantErr {
				t.Errorf("Previous() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got.Int() != tt.want {
				t.Errorf("Previous() = %v, want %v", got.Int(), tt.want)
			}
		})
	}
}

func TestVersionNumber_IsZero(t *testing.T) {
	zero := VersionNumber(0)
	if !zero.IsZero() {
		t.Error("IsZero() returned false for zero version")
	}

	v1, _ := NewVersionNumber(1)
	if v1.IsZero() {
		t.Error("IsZero() returned true for non-zero version")
	}
}

func TestVersionNumber_Equals(t *testing.T) {
	v1a, _ := NewVersionNumber(1)
	v1b, _ := NewVersionNumber(1)
	v2, _ := NewVersionNumber(2)

	if !v1a.Equals(v1b) {
		t.Error("Equals() returned false for identical version numbers")
	}
	if v1a.Equals(v2) {
		t.Error("Equals() returned true for different version numbers")
	}
}
