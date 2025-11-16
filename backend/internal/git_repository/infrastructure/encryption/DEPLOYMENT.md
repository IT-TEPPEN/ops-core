# Deployment Guide for Access Token Encryption

## Initial Deployment

### Prerequisites

1. PostgreSQL database running and accessible
2. Database migrations applied (up to migration 000003)
3. Secure key management solution in place (for production)

### Step 1: Generate Encryption Key

Generate a secure 32-byte encryption key:

```bash
# Using OpenSSL (recommended)
openssl rand -base64 32 > encryption_key.txt

# Store this key securely - never commit to version control!
```

### Step 2: Configure Environment

Set the encryption key as an environment variable:

```bash
# For local development
export ENCRYPTION_KEY=$(cat encryption_key.txt)

# For Docker
docker run -e ENCRYPTION_KEY="$(cat encryption_key.txt)" your-app

# For Kubernetes
kubectl create secret generic encryption-key \
  --from-literal=ENCRYPTION_KEY="$(cat encryption_key.txt)"
```

### Step 3: Deploy Application

Deploy the application with the encryption key configured. The application will:
- Automatically encrypt all new access tokens on save
- Automatically decrypt access tokens on retrieval

### Step 4: Verify Encryption

Verify that encryption is working:

```bash
# Check logs for successful startup (no encryption key errors)
# Test creating a new repository with an access token
# Query the database directly to verify tokens are encrypted

# Example SQL query to verify encryption:
SELECT id, name, 
       LENGTH(access_token) as encrypted_length,
       LEFT(access_token, 20) as encrypted_prefix
FROM repositories
LIMIT 5;

# The access_token should be base64-encoded and longer than plaintext
```

## Migration from Existing System

### WARNING: If You Have Existing Plaintext Tokens

If you already have access tokens stored in plaintext in your database, you need to handle migration carefully.

#### Option 1: Re-Enter Tokens (Recommended)

1. Deploy the application with encryption enabled
2. Notify users to re-enter their access tokens
3. Old plaintext tokens will fail to decrypt, requiring users to update

#### Option 2: Manual Migration Script

Create a migration script to encrypt existing tokens:

```go
package main

import (
    "context"
    "crypto/rand"
    "fmt"
    "log"
    "os"
    
    "opscore/backend/internal/git_repository/infrastructure/encryption"
    "github.com/jackc/pgx/v5/pgxpool"
)

func main() {
    // Get encryption key
    keyStr := os.Getenv("ENCRYPTION_KEY")
    if len(keyStr) != 32 {
        log.Fatal("ENCRYPTION_KEY must be 32 bytes")
    }
    
    encryptor, err := encryption.NewEncryptor([]byte(keyStr))
    if err != nil {
        log.Fatal(err)
    }
    
    // Connect to database
    dbURL := os.Getenv("DATABASE_URL")
    db, err := pgxpool.New(context.Background(), dbURL)
    if err != nil {
        log.Fatal(err)
    }
    defer db.Close()
    
    // Get all repositories
    rows, err := db.Query(context.Background(), 
        "SELECT id, access_token FROM repositories WHERE access_token IS NOT NULL AND access_token != ''")
    if err != nil {
        log.Fatal(err)
    }
    defer rows.Close()
    
    // Migrate each token
    for rows.Next() {
        var id, plainToken string
        if err := rows.Scan(&id, &plainToken); err != nil {
            log.Printf("Error scanning row: %v", err)
            continue
        }
        
        // Try to decrypt - if it fails, it's plaintext
        _, err := encryptor.Decrypt(plainToken)
        if err != nil {
            // Token is plaintext, encrypt it
            encrypted, err := encryptor.Encrypt(plainToken)
            if err != nil {
                log.Printf("Error encrypting token for repo %s: %v", id, err)
                continue
            }
            
            // Update database
            _, err = db.Exec(context.Background(),
                "UPDATE repositories SET access_token = $1 WHERE id = $2",
                encrypted, id)
            if err != nil {
                log.Printf("Error updating repo %s: %v", id, err)
                continue
            }
            
            log.Printf("Migrated token for repo %s", id)
        } else {
            log.Printf("Token for repo %s already encrypted", id)
        }
    }
    
    log.Println("Migration completed")
}
```

Run this script once during a maintenance window:

```bash
export ENCRYPTION_KEY="your-key-here"
export DATABASE_URL="your-db-url-here"
go run migrate_tokens.go
```

## Key Rotation

### When to Rotate Keys

- Periodically (e.g., every 90 days)
- When a key may have been compromised
- When changing environments
- During security audits

### Key Rotation Process

1. **Generate New Key**:
   ```bash
   openssl rand -base64 32 > new_encryption_key.txt
   ```

2. **Keep Old Key Available**:
   ```bash
   export OLD_ENCRYPTION_KEY=$(cat encryption_key.txt)
   export NEW_ENCRYPTION_KEY=$(cat new_encryption_key.txt)
   ```

3. **Run Migration Script**:
   ```go
   package main
   
   import (
       "context"
       "log"
       "os"
       
       "opscore/backend/internal/git_repository/infrastructure/encryption"
       "github.com/jackc/pgx/v5/pgxpool"
   )
   
   func main() {
       oldKey := os.Getenv("OLD_ENCRYPTION_KEY")
       newKey := os.Getenv("NEW_ENCRYPTION_KEY")
       
       if len(oldKey) != 32 || len(newKey) != 32 {
           log.Fatal("Keys must be 32 bytes")
       }
       
       oldEncryptor, _ := encryption.NewEncryptor([]byte(oldKey))
       newEncryptor, _ := encryption.NewEncryptor([]byte(newKey))
       
       db, err := pgxpool.New(context.Background(), os.Getenv("DATABASE_URL"))
       if err != nil {
           log.Fatal(err)
       }
       defer db.Close()
       
       rows, _ := db.Query(context.Background(), 
           "SELECT id, access_token FROM repositories WHERE access_token IS NOT NULL")
       
       for rows.Next() {
           var id, encryptedToken string
           rows.Scan(&id, &encryptedToken)
           
           // Decrypt with old key
           plaintext, err := oldEncryptor.Decrypt(encryptedToken)
           if err != nil {
               log.Printf("Failed to decrypt token for repo %s: %v", id, err)
               continue
           }
           
           // Re-encrypt with new key
           newEncrypted, _ := newEncryptor.Encrypt(plaintext)
           
           // Update database
           db.Exec(context.Background(),
               "UPDATE repositories SET access_token = $1 WHERE id = $2",
               newEncrypted, id)
           
           log.Printf("Rotated key for repo %s", id)
       }
   }
   ```

4. **Update Environment**:
   ```bash
   export ENCRYPTION_KEY=$(cat new_encryption_key.txt)
   ```

5. **Restart Application**

6. **Verify**:
   - Test retrieval of repositories
   - Verify tokens are decrypted correctly
   - Check application logs for errors

7. **Secure Old Key**:
   - Keep the old key in secure backup for 30 days
   - After verification period, securely delete old key

## Production Best Practices

### Key Management

1. **Use a Secret Management Service**:
   - AWS Secrets Manager
   - HashiCorp Vault
   - Google Secret Manager
   - Azure Key Vault

2. **Example with AWS Secrets Manager**:
   ```bash
   # Store the key
   aws secretsmanager create-secret \
     --name opscore/encryption-key \
     --secret-string "your-32-byte-key"
   
   # Retrieve in application startup
   ENCRYPTION_KEY=$(aws secretsmanager get-secret-value \
     --secret-id opscore/encryption-key \
     --query SecretString \
     --output text)
   ```

3. **Example with HashiCorp Vault**:
   ```bash
   # Store the key
   vault kv put secret/opscore encryption_key="your-32-byte-key"
   
   # Retrieve in application
   ENCRYPTION_KEY=$(vault kv get -field=encryption_key secret/opscore)
   ```

### Monitoring

1. **Log Encryption Errors** (but never log keys or tokens):
   ```go
   if err != nil {
       log.Error("Failed to encrypt access token", "error", err, "repo_id", repoID)
   }
   ```

2. **Set Up Alerts**:
   - Alert on encryption/decryption failures
   - Alert on missing encryption key at startup
   - Monitor key rotation completion

### Backup and Recovery

1. **Backup Encryption Key Separately**:
   - Store in secure, separate location from database backups
   - Use multiple secure storage locations
   - Document key recovery procedures

2. **Test Recovery Procedures**:
   - Regularly test restoring from backups
   - Verify encrypted data can be decrypted with backed-up keys
   - Document and test disaster recovery procedures

### Compliance

1. **Document Encryption Implementation**:
   - Algorithm: AES-256-GCM
   - Key size: 256 bits
   - Key rotation frequency
   - Access controls

2. **Audit Trail**:
   - Log key rotations (not the keys themselves)
   - Track access to encryption keys
   - Maintain compliance documentation

## Troubleshooting

### Issue: "invalid encryption key: must be 32 bytes"

**Solution**: Ensure ENCRYPTION_KEY is exactly 32 bytes:
```bash
echo -n "$ENCRYPTION_KEY" | wc -c  # Should output 32
```

### Issue: "failed to decrypt access token"

**Possible causes**:
1. Encryption key changed without migrating data
2. Plaintext token in database (not yet encrypted)
3. Database corruption

**Solution**:
1. Verify ENCRYPTION_KEY matches the key used to encrypt
2. Check if token is plaintext and needs migration
3. Restore from backup if corruption detected

### Issue: Application won't start

**Check**:
1. ENCRYPTION_KEY environment variable is set
2. Key is exactly 32 bytes
3. No special characters causing shell interpretation issues

## Security Checklist

- [ ] Encryption key is 32 bytes (256 bits)
- [ ] Key stored in secure key management service (production)
- [ ] Key never committed to version control
- [ ] Different keys for different environments
- [ ] Key rotation procedure documented and tested
- [ ] Backup keys stored securely
- [ ] TLS/SSL enabled for database connections
- [ ] Database backups are encrypted
- [ ] Access to encryption keys is restricted and audited
- [ ] Monitoring and alerting configured
- [ ] Compliance documentation completed
