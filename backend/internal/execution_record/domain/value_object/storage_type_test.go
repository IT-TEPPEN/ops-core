package value_object

import "testing"

func TestNewStorageType(t *testing.T) {
	tests := []struct {
		name        string
		storageType string
		wantErr     bool
	}{
		{
			name:        "valid local",
			storageType: "local",
			wantErr:     false,
		},
		{
			name:        "valid s3",
			storageType: "s3",
			wantErr:     false,
		},
		{
			name:        "valid minio",
			storageType: "minio",
			wantErr:     false,
		},
		{
			name:        "invalid storage type",
			storageType: "invalid",
			wantErr:     true,
		},
		{
			name:        "empty storage type",
			storageType: "",
			wantErr:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewStorageType(tt.storageType)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewStorageType() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got.String() != tt.storageType {
				t.Errorf("NewStorageType() = %v, want %v", got, tt.storageType)
			}
		})
	}
}

func TestStorageType_IsValid(t *testing.T) {
	tests := []struct {
		name        string
		storageType StorageType
		want        bool
	}{
		{
			name:        "local is valid",
			storageType: StorageTypeLocal,
			want:        true,
		},
		{
			name:        "s3 is valid",
			storageType: StorageTypeS3,
			want:        true,
		},
		{
			name:        "minio is valid",
			storageType: StorageTypeMinio,
			want:        true,
		},
		{
			name:        "invalid storage type",
			storageType: StorageType("invalid"),
			want:        false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.storageType.IsValid(); got != tt.want {
				t.Errorf("StorageType.IsValid() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStorageType_IsLocal(t *testing.T) {
	tests := []struct {
		name        string
		storageType StorageType
		want        bool
	}{
		{
			name:        "local returns true",
			storageType: StorageTypeLocal,
			want:        true,
		},
		{
			name:        "s3 returns false",
			storageType: StorageTypeS3,
			want:        false,
		},
		{
			name:        "minio returns false",
			storageType: StorageTypeMinio,
			want:        false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.storageType.IsLocal(); got != tt.want {
				t.Errorf("StorageType.IsLocal() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStorageType_IsS3Compatible(t *testing.T) {
	tests := []struct {
		name        string
		storageType StorageType
		want        bool
	}{
		{
			name:        "local returns false",
			storageType: StorageTypeLocal,
			want:        false,
		},
		{
			name:        "s3 returns true",
			storageType: StorageTypeS3,
			want:        true,
		},
		{
			name:        "minio returns true",
			storageType: StorageTypeMinio,
			want:        true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.storageType.IsS3Compatible(); got != tt.want {
				t.Errorf("StorageType.IsS3Compatible() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStorageType_Equals(t *testing.T) {
	st1, _ := NewStorageType("local")
	st2, _ := NewStorageType("s3")
	st1Copy, _ := NewStorageType("local")

	if !st1.Equals(st1Copy) {
		t.Error("Same storage types should be equal")
	}
	if st1.Equals(st2) {
		t.Error("Different storage types should not be equal")
	}
}
