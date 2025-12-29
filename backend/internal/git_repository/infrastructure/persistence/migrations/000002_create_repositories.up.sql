-- Filepath: backend/internal/git_repository/infrastructure/persistence/migrations/000002_create_repositories.up.sql
-- Create repositories table

CREATE TABLE repositories (
    id UUID PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    url VARCHAR(255) NOT NULL UNIQUE,
    access_token_encrypted TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_repositories_url ON repositories(url);

-- Create managed_files table
CREATE TABLE managed_files (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    repository_id UUID NOT NULL REFERENCES repositories(id) ON DELETE CASCADE,
    file_path TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    UNIQUE (repository_id, file_path)
);

CREATE INDEX idx_managed_files_repository_id ON managed_files(repository_id);

-- Comments
COMMENT ON TABLE repositories IS 'Git repositories registered in OpScore';
COMMENT ON COLUMN repositories.access_token_encrypted IS 'AES-256-GCM encrypted access token (for backward compatibility, prefer OAuth connections)';
COMMENT ON TABLE managed_files IS 'Files managed by OpScore in repositories';
