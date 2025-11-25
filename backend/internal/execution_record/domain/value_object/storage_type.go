package value_object

import "errors"

// StorageType represents the type of storage for attachments.
type StorageType string

const (
	// StorageTypeLocal represents local file system storage.
	StorageTypeLocal StorageType = "local"
	// StorageTypeS3 represents S3 storage.
	StorageTypeS3 StorageType = "s3"
	// StorageTypeMinio represents MinIO storage.
	StorageTypeMinio StorageType = "minio"
)

// NewStorageType creates a new StorageType from a string.
func NewStorageType(storageType string) (StorageType, error) {
	st := StorageType(storageType)
	if !st.IsValid() {
		return "", errors.New("invalid storage type: must be 'local', 's3', or 'minio'")
	}
	return st, nil
}

// IsValid checks if the StorageType is valid.
func (s StorageType) IsValid() bool {
	return s == StorageTypeLocal || s == StorageTypeS3 || s == StorageTypeMinio
}

// String returns the string representation of StorageType.
func (s StorageType) String() string {
	return string(s)
}

// IsLocal returns true if the storage type is local.
func (s StorageType) IsLocal() bool {
	return s == StorageTypeLocal
}

// IsS3Compatible returns true if the storage type is S3 or MinIO.
func (s StorageType) IsS3Compatible() bool {
	return s == StorageTypeS3 || s == StorageTypeMinio
}

// Equals checks if two StorageTypes are equal.
func (s StorageType) Equals(other StorageType) bool {
	return s == other
}
