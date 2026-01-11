-- Filepath: backend/migrate/000005_create_document_versions_and_tags.up.sql
-- Create document_versions, document_current_version, tags, and document_tags tables

-- Document versions table (versioned snapshots; latest via version_no DESC)
CREATE TABLE document_versions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    document_id UUID NOT NULL REFERENCES documents(id) ON DELETE CASCADE,
    version_no INTEGER NOT NULL,
    commit_hash VARCHAR(64),
    connection_by UUID REFERENCES oauth_connections(id),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE (document_id, version_no)
);

CREATE INDEX idx_document_versions_latest ON document_versions(document_id, version_no DESC);
CREATE INDEX idx_document_versions_connection_id ON document_versions(connection_id);

-- Document current version table (stores latest version content)
CREATE TABLE document_current_version (
    document_id UUID PRIMARY KEY REFERENCES documents(id) ON DELETE CASCADE,
    document_version_id UUID NOT NULL REFERENCES document_versions(id) ON DELETE CASCADE,
    title TEXT,
    description TEXT,
    document_type VARCHAR(64),
    content TEXT NOT NULL
);

CREATE INDEX idx_document_current_version_document_version_id ON document_current_version(document_version_id);

-- Tags table
CREATE TABLE tags (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(), 
    name VARCHAR(64) NOT NULL UNIQUE
);

-- Document tags junction table
CREATE TABLE document_tags (
    document_version_id UUID NOT NULL REFERENCES document_versions(id) ON DELETE CASCADE,
    tag_id UUID NOT NULL REFERENCES tags(id) ON DELETE CASCADE,
    PRIMARY KEY (document_version_id, tag_id)
);

CREATE INDEX idx_document_tags_document_version_id ON document_tags(document_version_id);
CREATE INDEX idx_document_tags_tag_id ON document_tags(tag_id);

-- Comments
COMMENT ON TABLE document_versions IS 'Version history of documents (latest via version_no DESC)';
COMMENT ON TABLE document_current_version IS 'Current version content of each document';
COMMENT ON TABLE tags IS 'Tag master table for categorizing documents';
COMMENT ON TABLE document_tags IS 'Many-to-many relationship between document versions and tags';
