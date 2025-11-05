-- Filepath: backend/infrastructure/persistence/migrations/000002_add_access_token_to_repositories.down.sql
-- Remove access_token column from repositories table
ALTER TABLE repositories
DROP COLUMN IF EXISTS access_token;