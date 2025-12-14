-- Migration: Add performance optimization indexes
-- This migration adds indexes to improve query performance for frequently accessed data

-- Document indexes
CREATE INDEX IF NOT EXISTS idx_documents_repository_id ON documents (repository_id);
CREATE INDEX IF NOT EXISTS idx_documents_owner ON documents (owner);
CREATE INDEX IF NOT EXISTS idx_documents_is_published ON documents (is_published);
CREATE INDEX IF NOT EXISTS idx_documents_published ON documents (is_published, access_scope) WHERE is_published = true;
CREATE INDEX IF NOT EXISTS idx_documents_created_at ON documents (created_at);

-- Document version indexes
CREATE INDEX IF NOT EXISTS idx_document_versions_document_id ON document_versions (document_id);
CREATE UNIQUE INDEX IF NOT EXISTS idx_document_versions_document_version ON document_versions (document_id, version_number);
CREATE INDEX IF NOT EXISTS idx_document_versions_commit_hash ON document_versions (commit_hash);

-- Execution record indexes
CREATE INDEX IF NOT EXISTS idx_execution_records_document_id ON execution_records (document_id);
CREATE INDEX IF NOT EXISTS idx_execution_records_executor_id ON execution_records (executor_id);
CREATE INDEX IF NOT EXISTS idx_execution_records_status ON execution_records (status);
CREATE INDEX IF NOT EXISTS idx_execution_records_started_at ON execution_records (started_at);
CREATE INDEX IF NOT EXISTS idx_execution_records_completed_at ON execution_records (completed_at);
CREATE INDEX IF NOT EXISTS idx_execution_records_search ON execution_records (document_id, executor_id, status, started_at);

-- Execution record step indexes
CREATE INDEX IF NOT EXISTS idx_execution_steps_record_id ON execution_steps (execution_record_id);
CREATE UNIQUE INDEX IF NOT EXISTS idx_execution_steps_record_step ON execution_steps (execution_record_id, step_number);

-- Attachment indexes
CREATE INDEX IF NOT EXISTS idx_attachments_execution_record_id ON attachments (execution_record_id);
CREATE INDEX IF NOT EXISTS idx_attachments_step_id ON attachments (step_id);

-- View history indexes
CREATE INDEX IF NOT EXISTS idx_view_history_document_id ON view_history (document_id);
CREATE INDEX IF NOT EXISTS idx_view_history_user_id ON view_history (user_id);
CREATE INDEX IF NOT EXISTS idx_view_history_viewed_at ON view_history (viewed_at);
CREATE INDEX IF NOT EXISTS idx_view_history_aggregation ON view_history (document_id, viewed_at);

-- View statistics indexes
CREATE UNIQUE INDEX IF NOT EXISTS idx_view_statistics_document_id ON view_statistics (document_id);
CREATE INDEX IF NOT EXISTS idx_view_statistics_total_views ON view_statistics (total_views);
CREATE INDEX IF NOT EXISTS idx_view_statistics_last_viewed_at ON view_statistics (last_viewed_at);
CREATE INDEX IF NOT EXISTS idx_view_statistics_popularity ON view_statistics (total_views, last_viewed_at);

-- User indexes
CREATE UNIQUE INDEX IF NOT EXISTS idx_users_email ON users (email);
CREATE UNIQUE INDEX IF NOT EXISTS idx_users_username ON users (username);

-- Group indexes
CREATE INDEX IF NOT EXISTS idx_groups_name ON groups (name);

-- Group member indexes
CREATE INDEX IF NOT EXISTS idx_group_members_group_id ON group_members (group_id);
CREATE INDEX IF NOT EXISTS idx_group_members_user_id ON group_members (user_id);
CREATE UNIQUE INDEX IF NOT EXISTS idx_group_members_composite ON group_members (group_id, user_id);
