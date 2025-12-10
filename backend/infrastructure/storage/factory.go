package storage

import (
	"errors"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/credentials"
)

// StorageType represents the type of storage backend.
type StorageType string

const (
	// StorageTypeLocal represents local file system storage.
	StorageTypeLocal StorageType = "local"
	// StorageTypeS3 represents AWS S3 storage.
	StorageTypeS3 StorageType = "s3"
	// StorageTypeMinio represents MinIO storage.
	StorageTypeMinio StorageType = "minio"
)

// StorageConfig holds configuration for creating a Storage instance.
type StorageConfig struct {
	// Type specifies the storage backend type: "local", "s3", or "minio"
	Type StorageType

	// LocalPath is the base path for local file system storage (required for local)
	LocalPath string

	// S3Bucket is the S3 bucket name (required for s3/minio)
	S3Bucket string

	// S3Region is the AWS region (required for s3, optional for minio)
	S3Region string

	// S3Endpoint is the custom endpoint URL (required for minio, optional for s3)
	S3Endpoint string

	// S3AccessKeyID is the AWS access key ID (required for s3/minio)
	S3AccessKeyID string

	// S3SecretAccessKey is the AWS secret access key (required for s3/minio)
	S3SecretAccessKey string
}

// NewStorage creates a new Storage instance based on the configuration.
func NewStorage(config StorageConfig) (Storage, error) {
	switch config.Type {
	case StorageTypeLocal:
		return newLocalStorage(config)
	case StorageTypeS3:
		return newS3Storage(config, false)
	case StorageTypeMinio:
		return newS3Storage(config, true)
	default:
		return nil, fmt.Errorf("unsupported storage type: %s (must be 'local', 's3', or 'minio')", config.Type)
	}
}

// newLocalStorage creates a local file system storage.
func newLocalStorage(config StorageConfig) (Storage, error) {
	if config.LocalPath == "" {
		return nil, errors.New("local path is required for local storage")
	}

	return NewLocalStorage(config.LocalPath)
}

// newS3Storage creates an S3 or MinIO storage.
func newS3Storage(config StorageConfig, isMinio bool) (Storage, error) {
	if config.S3Bucket == "" {
		return nil, errors.New("S3 bucket is required for S3/MinIO storage")
	}

	if config.S3AccessKeyID == "" || config.S3SecretAccessKey == "" {
		return nil, errors.New("S3 credentials (access key and secret key) are required")
	}

	// For MinIO, endpoint is required
	if isMinio && config.S3Endpoint == "" {
		return nil, errors.New("S3 endpoint is required for MinIO storage")
	}

	// Create AWS credentials
	creds := credentials.NewStaticCredentialsProvider(
		config.S3AccessKeyID,
		config.S3SecretAccessKey,
		"", // session token (empty for static credentials)
	)

	// Create AWS config
	awsConfig := aws.Config{
		Region:      config.S3Region,
		Credentials: creds,
	}

	// If region is not specified, use a default
	if awsConfig.Region == "" {
		if isMinio {
			awsConfig.Region = "us-east-1" // MinIO default
		} else {
			return nil, errors.New("S3 region is required for AWS S3")
		}
	}

	// Create S3 storage configuration
	s3Config := S3Config{
		Bucket:   config.S3Bucket,
		Region:   awsConfig.Region,
		Endpoint: config.S3Endpoint,
	}

	return NewS3Storage(awsConfig, s3Config)
}
