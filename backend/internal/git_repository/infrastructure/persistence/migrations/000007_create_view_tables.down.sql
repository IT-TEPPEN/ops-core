-- Filepath: backend/internal/git_repository/infrastructure/persistence/migrations/000007_create_view_tables.down.sql
-- Drop view_history and view_statistics tables

DROP TABLE IF EXISTS view_statistics;
DROP TABLE IF EXISTS view_history;
