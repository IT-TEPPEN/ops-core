-- Filepath: backend/migrate/000002_create_git_providers.up.sql
-- Create git_providers table for managing Git provider configurations

CREATE TABLE git_providers (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    provider_type VARCHAR(50) NOT NULL CHECK (provider_type IN ('github', 'gitlab', 'bitbucket', 'gitea', 'other')),
    provider_name VARCHAR(255) NOT NULL UNIQUE,
    base_url VARCHAR(255) NOT NULL,
    api_url VARCHAR(255) NOT NULL,
    is_saas BOOLEAN NOT NULL DEFAULT FALSE,
    client_id VARCHAR(255),
    client_secret_encrypted TEXT,
    enabled BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT chk_saas_credentials CHECK (
        (is_saas = TRUE AND client_id IS NULL AND client_secret_encrypted IS NULL) OR
        (is_saas = FALSE AND client_id IS NOT NULL AND client_secret_encrypted IS NOT NULL)
    )
);

CREATE INDEX idx_git_providers_provider_type ON git_providers(provider_type);
CREATE INDEX idx_git_providers_enabled ON git_providers(enabled);
CREATE INDEX idx_git_providers_is_saas ON git_providers(is_saas);

-- Comments
COMMENT ON TABLE git_providers IS 'Git provider configurations (GitHub, GitLab, etc.)';
COMMENT ON COLUMN git_providers.is_saas IS 'TRUE for SaaS providers (credentials from env vars), FALSE for self-hosted';
COMMENT ON COLUMN git_providers.client_id IS 'OAuth client ID (NULL for SaaS providers using env vars)';
COMMENT ON COLUMN git_providers.client_secret_encrypted IS 'Encrypted OAuth client secret (NULL for SaaS providers using env vars)';
