package storage

import (
"context"
"errors"
"fmt"
"io"
"time"

"github.com/aws/aws-sdk-go-v2/aws"
"github.com/aws/aws-sdk-go-v2/service/s3"
"github.com/aws/aws-sdk-go-v2/service/s3/types"
)

// S3Storage implements Storage interface for S3-compatible storage (AWS S3, MinIO, etc.).
type S3Storage struct {
client   *s3.Client
bucket   string
region   string
endpoint string
}

// S3Config holds configuration for S3Storage.
type S3Config struct {
Bucket   string
Region   string
Endpoint string // Optional: for MinIO or custom S3-compatible endpoints
}

// NewS3Storage creates a new S3Storage instance.
// For MinIO, set the endpoint in the AWS config's BaseEndpoint.
func NewS3Storage(awsConfig aws.Config, config S3Config) (*S3Storage, error) {
if config.Bucket == "" {
return nil, errors.New("bucket name is required")
}

// Create S3 client with optional custom endpoint for MinIO
var client *s3.Client
if config.Endpoint != "" {
client = s3.NewFromConfig(awsConfig, func(o *s3.Options) {
o.BaseEndpoint = aws.String(config.Endpoint)
o.UsePathStyle = true // Required for MinIO
})
} else {
client = s3.NewFromConfig(awsConfig)
}

return &S3Storage{
client:   client,
bucket:   config.Bucket,
region:   config.Region,
endpoint: config.Endpoint,
}, nil
}

// Save uploads a file to S3.
func (s *S3Storage) Save(ctx context.Context, key string, reader io.Reader, contentType string) error {
if key == "" {
return errors.New("key cannot be empty")
}
if reader == nil {
return errors.New("reader cannot be nil")
}

input := &s3.PutObjectInput{
Bucket: aws.String(s.bucket),
Key:    aws.String(key),
Body:   reader,
}

if contentType != "" {
input.ContentType = aws.String(contentType)
}

_, err := s.client.PutObject(ctx, input)
if err != nil {
return fmt.Errorf("failed to upload to S3: %w", err)
}

return nil
}

// Get retrieves a file from S3.
func (s *S3Storage) Get(ctx context.Context, key string) (io.ReadCloser, error) {
if key == "" {
return nil, errors.New("key cannot be empty")
}

input := &s3.GetObjectInput{
Bucket: aws.String(s.bucket),
Key:    aws.String(key),
}

result, err := s.client.GetObject(ctx, input)
if err != nil {
return nil, fmt.Errorf("failed to get object from S3: %w", err)
}

return result.Body, nil
}

// Delete removes a file from S3.
func (s *S3Storage) Delete(ctx context.Context, key string) error {
if key == "" {
return errors.New("key cannot be empty")
}

input := &s3.DeleteObjectInput{
Bucket: aws.String(s.bucket),
Key:    aws.String(key),
}

_, err := s.client.DeleteObject(ctx, input)
if err != nil {
return fmt.Errorf("failed to delete object from S3: %w", err)
}

// S3 DeleteObject is idempotent - it returns success even if the object doesn't exist
return nil
}

// Exists checks if a file exists in S3.
func (s *S3Storage) Exists(ctx context.Context, key string) (bool, error) {
if key == "" {
return false, errors.New("key cannot be empty")
}

input := &s3.HeadObjectInput{
Bucket: aws.String(s.bucket),
Key:    aws.String(key),
}

_, err := s.client.HeadObject(ctx, input)
if err != nil {
// Check if it's a "not found" error
var noSuchKey *types.NoSuchKey
if errors.As(err, &noSuchKey) {
return false, nil
}
// Also check for NotFound which may be returned in some cases
var notFound *types.NotFound
if errors.As(err, &notFound) {
return false, nil
}
return false, fmt.Errorf("failed to check object existence: %w", err)
}

return true, nil
}

// GetSignedURL generates a pre-signed URL for temporary access to the file.
func (s *S3Storage) GetSignedURL(ctx context.Context, key string, expiration time.Duration) (string, error) {
if key == "" {
return "", errors.New("key cannot be empty")
}

presignClient := s3.NewPresignClient(s.client)

input := &s3.GetObjectInput{
Bucket: aws.String(s.bucket),
Key:    aws.String(key),
}

presignedReq, err := presignClient.PresignGetObject(ctx, input, func(opts *s3.PresignOptions) {
opts.Expires = expiration
})
if err != nil {
return "", fmt.Errorf("failed to generate presigned URL: %w", err)
}

return presignedReq.URL, nil
}
