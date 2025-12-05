-- Filepath: backend/internal/git_repository/infrastructure/persistence/migrations/000007_create_view_tables.up.sql
-- Create view_history and view_statistics tables

-- view_history table
CREATE TABLE view_history (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    document_id UUID NOT NULL REFERENCES documents(id) ON DELETE CASCADE,
    user_id UUID REFERENCES users(id) ON DELETE SET NULL,
    ip_address INET,
    user_agent TEXT,
    viewed_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_view_history_document_id ON view_history(document_id);
CREATE INDEX idx_view_history_user_id ON view_history(user_id);
CREATE INDEX idx_view_history_viewed_at ON view_history(viewed_at);

-- view_statistics table
CREATE TABLE view_statistics (
    document_id UUID PRIMARY KEY REFERENCES documents(id) ON DELETE CASCADE,
    total_views BIGINT NOT NULL DEFAULT 0,
    unique_users BIGINT NOT NULL DEFAULT 0,
    last_viewed_at TIMESTAMPTZ,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);
