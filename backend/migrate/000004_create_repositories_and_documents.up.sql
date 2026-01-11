-- Filepath: backend/migrate/000004_create_repositories_and_documents.up.sql
-- Create repositories, managed_files, and documents tables

-- Repositories table
CREATE TABLE repositories (
    id UUID PRIMARY KEY,
    provider_repository_id VARCHAR(255) NOT NULL UNIQUE,
    provider_repository_owner VARCHAR(255) NOT NULL,
    provider_repository_name VARCHAR(255) NOT NULL,
    url VARCHAR(255) NOT NULL UNIQUE,
    created_by UUID REFERENCES users(id),
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_repositories_url ON repositories(url);
CREATE INDEX idx_repositories_created_by ON repositories(created_by);

-- Managed files table
CREATE TABLE managed_files (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    repository_id UUID NOT NULL REFERENCES repositories(id) ON DELETE CASCADE,
    file_path TEXT NOT NULL,
    branch VARCHAR(255),
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(repository_id, file_path, branch)
);

CREATE INDEX idx_managed_files_repository_id ON managed_files(repository_id);

-- Documents table (logical document bound to a managed file)
CREATE TABLE documents (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    managed_file_id UUID NOT NULL REFERENCES managed_files(id) ON DELETE CASCADE,
    connection_id UUID REFERENCES oauth_connections(id),
    is_auto_update BOOLEAN NOT NULL DEFAULT FALSE,
    enabled BOOLEAN NOT NULL DEFAULT TRUE,
    created_by UUID REFERENCES users(id),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_documents_managed_file_id ON documents(managed_file_id);
CREATE INDEX idx_documents_connection_id ON documents(connection_id);
CREATE INDEX idx_documents_created_by ON documents(created_by);
CREATE INDEX idx_documents_enabled ON documents(enabled);

-- Comments
COMMENT ON TABLE repositories IS 'Git repositories registered in OpScore';
COMMENT ON TABLE managed_files IS 'Files managed by OpScore in repositories';
COMMENT ON TABLE documents IS 'Logical documents managed per repository file';
COMMENT ON COLUMN documents.is_auto_update IS 'Whether to automatically fetch updates from the repository';
