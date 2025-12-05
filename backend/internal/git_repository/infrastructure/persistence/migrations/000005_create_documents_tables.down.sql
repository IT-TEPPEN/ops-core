-- Filepath: backend/internal/git_repository/infrastructure/persistence/migrations/000005_create_documents_tables.down.sql
-- Drop documents and document_versions tables

ALTER TABLE documents DROP CONSTRAINT IF EXISTS fk_documents_current_version;
DROP TABLE IF EXISTS document_versions;
DROP TABLE IF EXISTS documents;
