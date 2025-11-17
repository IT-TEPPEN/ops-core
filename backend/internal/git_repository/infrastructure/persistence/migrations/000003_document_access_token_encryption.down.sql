-- Filepath: backend/infrastructure/persistence/migrations/000003_document_access_token_encryption.down.sql
-- Remove comment from access_token column
COMMENT ON COLUMN repositories.access_token IS NULL;
