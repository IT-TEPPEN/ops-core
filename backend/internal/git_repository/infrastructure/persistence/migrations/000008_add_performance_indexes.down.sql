-- Migration rollback: Remove performance optimization indexes

-- Document indexes
DROP INDEX IF EXISTS idx_documents_repository_id;
DROP INDEX IF EXISTS idx_documents_owner;
DROP INDEX IF EXISTS idx_documents_is_published;
DROP INDEX IF EXISTS idx_documents_published;
DROP INDEX IF EXISTS idx_documents_created_at;

-- Document version indexes
DROP INDEX IF EXISTS idx_document_versions_document_id;
DROP INDEX IF EXISTS idx_document_versions_document_version;
DROP INDEX IF EXISTS idx_document_versions_commit_hash;

-- Execution record indexes
DROP INDEX IF EXISTS idx_execution_records_document_id;
DROP INDEX IF EXISTS idx_execution_records_executor_id;
DROP INDEX IF EXISTS idx_execution_records_status;
DROP INDEX IF EXISTS idx_execution_records_started_at;
DROP INDEX IF EXISTS idx_execution_records_completed_at;
DROP INDEX IF EXISTS idx_execution_records_search;

-- Execution record step indexes
DROP INDEX IF EXISTS idx_execution_steps_record_id;
DROP INDEX IF EXISTS idx_execution_steps_record_step;

-- Attachment indexes
DROP INDEX IF EXISTS idx_attachments_execution_record_id;
DROP INDEX IF EXISTS idx_attachments_step_id;

-- View history indexes
DROP INDEX IF EXISTS idx_view_history_document_id;
DROP INDEX IF EXISTS idx_view_history_user_id;
DROP INDEX IF EXISTS idx_view_history_viewed_at;
DROP INDEX IF EXISTS idx_view_history_aggregation;

-- View statistics indexes
DROP INDEX IF EXISTS idx_view_statistics_document_id;
DROP INDEX IF EXISTS idx_view_statistics_total_views;
DROP INDEX IF EXISTS idx_view_statistics_last_viewed_at;
DROP INDEX IF EXISTS idx_view_statistics_popularity;

-- User indexes
DROP INDEX IF EXISTS idx_users_email;
DROP INDEX IF EXISTS idx_users_username;

-- Group indexes
DROP INDEX IF EXISTS idx_groups_name;

-- Group member indexes
DROP INDEX IF EXISTS idx_group_members_group_id;
DROP INDEX IF EXISTS idx_group_members_user_id;
DROP INDEX IF EXISTS idx_group_members_composite;
