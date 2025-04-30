-- Filepath: backend/infrastructure/persistence/migrations/000002_add_access_token_to_repositories.up.sql
-- Add access_token column to repositories table
ALTER TABLE repositories
ADD COLUMN access_token TEXT;