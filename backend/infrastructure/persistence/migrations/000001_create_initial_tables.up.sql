-- Filepath: backend/infrastructure/persistence/migrations/000001_create_initial_tables.up.sql
-- Create repositories table
CREATE TABLE repositories (
    id UUID PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    url VARCHAR(255) NOT NULL UNIQUE, -- Ensure URL is unique
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Create managed_files table
CREATE TABLE managed_files (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(), -- Use default UUID generation
    repository_id UUID NOT NULL,
    file_path TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (repository_id) REFERENCES repositories(id) ON DELETE CASCADE, -- Add foreign key constraint
    UNIQUE (repository_id, file_path) -- Ensure file path is unique within a repository
);

-- Optional: Create indexes for faster lookups if needed
-- CREATE INDEX idx_repositories_url ON repositories(url);
-- CREATE INDEX idx_managed_files_repository_id ON managed_files(repository_id);
