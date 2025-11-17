-- Filepath: backend/infrastructure/persistence/migrations/000003_document_access_token_encryption.up.sql
-- Add comment to document that access_token is encrypted at the application level
COMMENT ON COLUMN repositories.access_token IS 'Encrypted using AES-256-GCM at application level. Stored as base64-encoded ciphertext.';
