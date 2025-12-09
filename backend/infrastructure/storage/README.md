# Storage Abstraction Layer

This package provides a unified interface for file storage operations, supporting multiple backends including local file system, AWS S3, and MinIO.

## Features

- **Multiple Storage Backends**: Local file system, AWS S3, and MinIO support
- **Security**: Path traversal protection for local storage
- **Pre-signed URLs**: Support for temporary file access URLs (S3/MinIO)
- **Factory Pattern**: Easy configuration-based storage creation
- **Comprehensive Testing**: Full test coverage with security tests

## Installation

The package requires AWS SDK v2:

```bash
go get github.com/aws/aws-sdk-go-v2/service/s3
go get github.com/aws/aws-sdk-go-v2/config
```

## Usage

### Creating Storage Instances

#### Using the Factory (Recommended)

```go
import "opscore/backend/infrastructure/storage"

// Local file system storage
config := storage.StorageConfig{
    Type:      storage.StorageTypeLocal,
    LocalPath: "/data/attachments",
}
store, err := storage.NewStorage(config)

// AWS S3 storage
config := storage.StorageConfig{
    Type:              storage.StorageTypeS3,
    S3Bucket:          "my-bucket",
    S3Region:          "us-east-1",
    S3AccessKeyID:     "AKIAIOSFODNN7EXAMPLE",
    S3SecretAccessKey: "wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY",
}
store, err := storage.NewStorage(config)

// MinIO storage
config := storage.StorageConfig{
    Type:              storage.StorageTypeMinio,
    S3Bucket:          "my-bucket",
    S3Endpoint:        "http://localhost:9000",
    S3AccessKeyID:     "minioadmin",
    S3SecretAccessKey: "minioadmin",
    S3Region:          "us-east-1", // Optional for MinIO
}
store, err := storage.NewStorage(config)
```

#### Direct Instantiation

```go
// Local storage
localStorage, err := storage.NewLocalStorage("/data/attachments")

// S3 storage (requires AWS config)
import (
    "github.com/aws/aws-sdk-go-v2/aws"
    "github.com/aws/aws-sdk-go-v2/credentials"
)

awsConfig := aws.Config{
    Region:      "us-east-1",
    Credentials: credentials.NewStaticCredentialsProvider(accessKey, secretKey, ""),
}

s3Config := storage.S3Config{
    Bucket:   "my-bucket",
    Region:   "us-east-1",
    Endpoint: "", // Empty for AWS S3, or custom endpoint for MinIO
}

s3Storage, err := storage.NewS3Storage(awsConfig, s3Config)
```

### Basic Operations

```go
import (
    "context"
    "bytes"
    "time"
)

ctx := context.Background()

// Save a file
content := bytes.NewReader([]byte("file content"))
err := store.Save(ctx, "path/to/file.txt", content, "text/plain")

// Get a file
reader, err := store.Get(ctx, "path/to/file.txt")
if err == nil {
    defer reader.Close()
    // Read content from reader
}

// Check if file exists
exists, err := store.Exists(ctx, "path/to/file.txt")

// Delete a file (idempotent)
err := store.Delete(ctx, "path/to/file.txt")

// Get a pre-signed URL (S3/MinIO only, returns empty string for local)
url, err := store.GetSignedURL(ctx, "path/to/file.txt", 1*time.Hour)
```

## Environment Variables

For configuration via environment variables:

```bash
# Storage type
ATTACHMENT_STORAGE_TYPE=local # or s3, minio

# Local storage settings
ATTACHMENT_STORAGE_PATH=/data/attachments

# S3/MinIO settings
S3_BUCKET=opscore-attachments
S3_REGION=us-east-1
S3_ENDPOINT=https://s3.amazonaws.com # MinIO: http://localhost:9000
S3_ACCESS_KEY_ID=your-access-key
S3_SECRET_ACCESS_KEY=your-secret-key
```

Example configuration loading:

```go
config := storage.StorageConfig{
    Type:              storage.StorageType(os.Getenv("ATTACHMENT_STORAGE_TYPE")),
    LocalPath:         os.Getenv("ATTACHMENT_STORAGE_PATH"),
    S3Bucket:          os.Getenv("S3_BUCKET"),
    S3Region:          os.Getenv("S3_REGION"),
    S3Endpoint:        os.Getenv("S3_ENDPOINT"),
    S3AccessKeyID:     os.Getenv("S3_ACCESS_KEY_ID"),
    S3SecretAccessKey: os.Getenv("S3_SECRET_ACCESS_KEY"),
}

store, err := storage.NewStorage(config)
```

## Security Considerations

### Local Storage

- **Path Traversal Protection**: All file paths are validated to prevent directory escape attacks
- **Base Directory Isolation**: Files are confined to the configured base directory
- **Robust Validation**: Uses `filepath.Rel()` for proper path validation

### S3 Storage

- **Credential Management**: Use AWS credentials best practices
- **Pre-signed URL Expiration**: Always set appropriate expiration times
- **IAM Policies**: Configure proper S3 bucket policies and IAM roles

## Testing

Run tests:

```bash
go test ./infrastructure/storage/...
```

Run tests with coverage:

```bash
go test -cover ./infrastructure/storage/...
```

## Examples

### Using with AttachmentRepository

```go
type AttachmentRepository struct {
    db      *sql.DB
    storage storage.Storage
}

func (r *AttachmentRepository) Save(ctx context.Context, attachment *model.Attachment, file io.Reader) error {
    // Generate storage key
    key := fmt.Sprintf("attachments/%s/%s", 
        attachment.ExecutionRecordID(), 
        attachment.ID())
    
    // Save to storage
    if err := r.storage.Save(ctx, key, file, attachment.MimeType()); err != nil {
        return fmt.Errorf("failed to save file: %w", err)
    }
    
    // Save metadata to database
    // ... database operations
    
    return nil
}
```

### Handling Large Files

```go
// For large files, use streaming
file, err := os.Open("large-file.bin")
if err != nil {
    return err
}
defer file.Close()

err = store.Save(ctx, "large/file.bin", file, "application/octet-stream")
```

## Interface

```go
type Storage interface {
    Save(ctx context.Context, key string, reader io.Reader, contentType string) error
    Get(ctx context.Context, key string) (io.ReadCloser, error)
    Delete(ctx context.Context, key string) error
    Exists(ctx context.Context, key string) (bool, error)
    GetSignedURL(ctx context.Context, key string, expiration time.Duration) (string, error)
}
```

## License

See project root LICENSE file.
