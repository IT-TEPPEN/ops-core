-- Add origin metadata columns to documents for provider-based filtering
ALTER TABLE documents
    ADD COLUMN provider_repository_id VARCHAR(255),
    ADD COLUMN repository VARCHAR(255);

-- Backfill existing rows to satisfy NOT NULL constraint
UPDATE documents
SET
    provider_repository_id = COALESCE(provider_repository_id, ''),
    repository = COALESCE(repository, '');

-- Enforce presence of origin metadata
ALTER TABLE documents
    ALTER COLUMN provider_repository_id SET NOT NULL,
    ALTER COLUMN repository SET NOT NULL;

-- Indexes to support list filters
CREATE INDEX idx_documents_provider_repository_id ON documents(provider_repository_id);
CREATE INDEX idx_documents_owner_repository ON documents(owner, repository);
