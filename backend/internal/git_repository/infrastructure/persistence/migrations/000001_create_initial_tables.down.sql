-- Filepath: backend/infrastructure/persistence/migrations/000001_create_initial_tables.down.sql
-- Drop managed_files table first due to foreign key constraint
DROP TABLE IF EXISTS managed_files;

-- Drop repositories table
DROP TABLE IF EXISTS repositories;
