-- Filepath: backend/internal/git_repository/infrastructure/persistence/migrations/000005_create_documents_tables.up.sql
-- Create documents and document_versions tables

-- documents table
CREATE TABLE documents (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    repository_id UUID NOT NULL REFERENCES repositories(id) ON DELETE CASCADE,
    owner VARCHAR(255) NOT NULL,
    is_published BOOLEAN NOT NULL DEFAULT false,
    is_auto_update BOOLEAN NOT NULL DEFAULT false,
    access_scope VARCHAR(50) NOT NULL CHECK (access_scope IN ('public', 'private')),
    current_version_id UUID,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_documents_repository_id ON documents(repository_id);
CREATE INDEX idx_documents_is_published ON documents(is_published);
CREATE INDEX idx_documents_access_scope ON documents(access_scope);

-- document_versions table
CREATE TABLE document_versions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    document_id UUID NOT NULL REFERENCES documents(id) ON DELETE CASCADE,
    version_number INTEGER NOT NULL,
    file_path TEXT NOT NULL,
    commit_hash VARCHAR(64) NOT NULL,
    title VARCHAR(255) NOT NULL,
    doc_type VARCHAR(50) NOT NULL CHECK (doc_type IN ('procedure', 'knowledge')),
    tags TEXT[] NOT NULL DEFAULT '{}',
    variables JSONB,
    content TEXT NOT NULL,
    published_at TIMESTAMPTZ NOT NULL,
    unpublished_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    UNIQUE (document_id, version_number),
    UNIQUE (document_id, file_path, commit_hash)
);

CREATE INDEX idx_document_versions_document_id ON document_versions(document_id);
CREATE INDEX idx_document_versions_doc_type ON document_versions(doc_type);
CREATE INDEX idx_document_versions_tags ON document_versions USING GIN(tags);
CREATE INDEX idx_document_versions_commit_hash ON document_versions(commit_hash);

-- Add foreign key constraint after document_versions table is created to avoid circular dependency
ALTER TABLE documents
ADD CONSTRAINT fk_documents_current_version
FOREIGN KEY (current_version_id) REFERENCES document_versions(id);
