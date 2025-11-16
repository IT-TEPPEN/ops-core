# Access Token Encryption

## Overview

This package provides AES-256-GCM encryption for sensitive access tokens stored in the database. As specified in ADR 0005, credentials must be encrypted to ensure security.

## Encryption Details

- **Algorithm**: AES-256-GCM (Galois/Counter Mode)
- **Key Size**: 32 bytes (256 bits)
- **Storage Format**: Base64-encoded ciphertext with embedded nonce
- **Nonce**: Randomly generated for each encryption (12 bytes)

## Setup

### Environment Variable

The encryption key must be provided via the `ENCRYPTION_KEY` environment variable:

```bash
export ENCRYPTION_KEY="your-32-byte-encryption-key-here"
```

**Important**: The key must be exactly 32 bytes (256 bits) for AES-256 encryption.

### Generating a Secure Key

To generate a secure random key, use one of these methods:

#### Using OpenSSL:
```bash
openssl rand -base64 32
```

#### Using Go:
```go
package main

import (
    "crypto/rand"
    "encoding/base64"
    "fmt"
)

func main() {
    key := make([]byte, 32)
    rand.Read(key)
    fmt.Println(base64.StdEncoding.EncodeToString(key))
}
```

### Key Management Best Practices

1. **Never commit the encryption key to version control**
2. **Use a secure key management service in production**:
   - AWS Secrets Manager
   - HashiCorp Vault
   - Google Secret Manager
   - Azure Key Vault

3. **Rotate encryption keys periodically**
4. **Store keys separately from the application code**
5. **Use different keys for different environments** (development, staging, production)

## Usage

The encryption is transparent to application code. The `PostgresRepository` automatically:
- Encrypts access tokens when saving to the database
- Decrypts access tokens when retrieving from the database

```go
// Example usage (encryption happens automatically)
repo := entity.NewRepository(id, name, url, "my-access-token")
err := repository.Save(ctx, repo)  // Token is encrypted before storage

// Retrieval (decryption happens automatically)
retrieved, err := repository.FindByID(ctx, id)
token := retrieved.AccessToken()  // Returns decrypted token
```

## Security Considerations

1. **Key Storage**: The encryption key should never be stored in the codebase. Use environment variables or secure key management services.

2. **Key Rotation**: When rotating keys, you'll need to:
   - Keep the old key available for decryption
   - Decrypt all tokens with the old key
   - Re-encrypt with the new key
   - Remove the old key

3. **Backup Security**: Database backups contain encrypted tokens, but they're only as secure as the encryption key. Ensure backups are stored securely.

4. **Network Security**: Use TLS/SSL for database connections to prevent token interception during transmission.

## Testing

Tests use randomly generated keys for each test run to ensure isolation:

```go
key := make([]byte, 32)
rand.Read(key)
encryptor, _ := encryption.NewEncryptor(key)
```

## Error Handling

The encryptor returns specific errors:
- `ErrInvalidKey`: The provided key is not 32 bytes
- `ErrInvalidCiphertext`: The ciphertext format is invalid or corrupted

## Migration from Plaintext

If you have existing plaintext tokens in the database:

1. **CRITICAL**: Back up your database before proceeding
2. The application will attempt to decrypt tokens on read
3. If decryption fails (plaintext tokens), the error will be returned
4. You'll need to manually re-encrypt or re-enter tokens

**Note**: It's recommended to start with encryption from the beginning, or perform a controlled migration with a maintenance window.

## Compliance

This implementation helps meet compliance requirements:
- SOC 2: Encryption of sensitive data at rest
- PCI DSS: Encryption of cardholder data
- GDPR: Appropriate security measures for personal data
- HIPAA: Encryption of protected health information
