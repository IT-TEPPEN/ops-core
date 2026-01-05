-- Drop indexes added for origin metadata
DROP INDEX IF EXISTS idx_documents_owner_repository;
DROP INDEX IF EXISTS idx_documents_provider_repository_id;

-- Drop origin metadata columns
ALTER TABLE documents
    DROP COLUMN IF EXISTS repository,
    DROP COLUMN IF EXISTS provider_repository_id;
