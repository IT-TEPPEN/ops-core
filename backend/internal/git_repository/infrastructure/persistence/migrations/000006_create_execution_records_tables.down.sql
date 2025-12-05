-- Filepath: backend/internal/git_repository/infrastructure/persistence/migrations/000006_create_execution_records_tables.down.sql
-- Drop execution_records, execution_steps, and attachments tables

DROP TABLE IF EXISTS attachments;
DROP TABLE IF EXISTS execution_steps;
DROP TABLE IF EXISTS execution_records;
