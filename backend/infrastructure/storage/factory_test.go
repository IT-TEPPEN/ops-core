package storage

import (
"testing"
)

func TestNewStorage_Local(t *testing.T) {
t.Run("creates local storage successfully", func(t *testing.T) {
tmpDir := t.TempDir()

config := StorageConfig{
Type:      StorageTypeLocal,
LocalPath: tmpDir,
}

storage, err := NewStorage(config)

if err != nil {
t.Fatalf("NewStorage() error = %v, want nil", err)
}
if storage == nil {
t.Fatal("NewStorage() returned nil")
}

// Verify it's a LocalStorage
_, ok := storage.(*LocalStorage)
if !ok {
t.Error("NewStorage() did not return LocalStorage instance")
}
})

t.Run("returns error when local path is missing", func(t *testing.T) {
config := StorageConfig{
Type: StorageTypeLocal,
}

_, err := NewStorage(config)

if err == nil {
t.Error("NewStorage() with missing LocalPath should return error")
}
})
}

func TestNewStorage_S3(t *testing.T) {
t.Run("creates S3 storage with valid config", func(t *testing.T) {
config := StorageConfig{
Type:              StorageTypeS3,
S3Bucket:          "test-bucket",
S3Region:          "us-east-1",
S3AccessKeyID:     "test-access-key",
S3SecretAccessKey: "test-secret-key",
}

storage, err := NewStorage(config)

if err != nil {
t.Fatalf("NewStorage() error = %v, want nil", err)
}
if storage == nil {
t.Fatal("NewStorage() returned nil")
}

// Verify it's an S3Storage
_, ok := storage.(*S3Storage)
if !ok {
t.Error("NewStorage() did not return S3Storage instance")
}
})

t.Run("returns error when bucket is missing", func(t *testing.T) {
config := StorageConfig{
Type:              StorageTypeS3,
S3Region:          "us-east-1",
S3AccessKeyID:     "test-access-key",
S3SecretAccessKey: "test-secret-key",
}

_, err := NewStorage(config)

if err == nil {
t.Error("NewStorage() with missing bucket should return error")
}
})

t.Run("returns error when region is missing", func(t *testing.T) {
config := StorageConfig{
Type:              StorageTypeS3,
S3Bucket:          "test-bucket",
S3AccessKeyID:     "test-access-key",
S3SecretAccessKey: "test-secret-key",
}

_, err := NewStorage(config)

if err == nil {
t.Error("NewStorage() with missing region should return error")
}
})

t.Run("returns error when credentials are missing", func(t *testing.T) {
config := StorageConfig{
Type:     StorageTypeS3,
S3Bucket: "test-bucket",
S3Region: "us-east-1",
}

_, err := NewStorage(config)

if err == nil {
t.Error("NewStorage() with missing credentials should return error")
}
})
}

func TestNewStorage_MinIO(t *testing.T) {
t.Run("creates MinIO storage with valid config", func(t *testing.T) {
config := StorageConfig{
Type:              StorageTypeMinio,
S3Bucket:          "test-bucket",
S3Region:          "us-east-1",
S3Endpoint:        "http://localhost:9000",
S3AccessKeyID:     "test-access-key",
S3SecretAccessKey: "test-secret-key",
}

storage, err := NewStorage(config)

if err != nil {
t.Fatalf("NewStorage() error = %v, want nil", err)
}
if storage == nil {
t.Fatal("NewStorage() returned nil")
}

// Verify it's an S3Storage
_, ok := storage.(*S3Storage)
if !ok {
t.Error("NewStorage() did not return S3Storage instance")
}
})

t.Run("uses default region if not specified", func(t *testing.T) {
config := StorageConfig{
Type:              StorageTypeMinio,
S3Bucket:          "test-bucket",
S3Endpoint:        "http://localhost:9000",
S3AccessKeyID:     "test-access-key",
S3SecretAccessKey: "test-secret-key",
}

storage, err := NewStorage(config)

if err != nil {
t.Fatalf("NewStorage() error = %v, want nil", err)
}
if storage == nil {
t.Fatal("NewStorage() returned nil")
}
})

t.Run("returns error when endpoint is missing", func(t *testing.T) {
config := StorageConfig{
Type:              StorageTypeMinio,
S3Bucket:          "test-bucket",
S3Region:          "us-east-1",
S3AccessKeyID:     "test-access-key",
S3SecretAccessKey: "test-secret-key",
}

_, err := NewStorage(config)

if err == nil {
t.Error("NewStorage() MinIO with missing endpoint should return error")
}
})

t.Run("returns error when credentials are missing", func(t *testing.T) {
config := StorageConfig{
Type:       StorageTypeMinio,
S3Bucket:   "test-bucket",
S3Endpoint: "http://localhost:9000",
}

_, err := NewStorage(config)

if err == nil {
t.Error("NewStorage() with missing credentials should return error")
}
})
}

func TestNewStorage_UnsupportedType(t *testing.T) {
config := StorageConfig{
Type: StorageType("unsupported"),
}

_, err := NewStorage(config)

if err == nil {
t.Error("NewStorage() with unsupported type should return error")
}
}
