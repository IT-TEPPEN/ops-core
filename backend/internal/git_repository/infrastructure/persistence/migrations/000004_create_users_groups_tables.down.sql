-- Filepath: backend/internal/git_repository/infrastructure/persistence/migrations/000004_create_users_groups_tables.down.sql
-- Drop users, groups, and user_groups tables

DROP TABLE IF EXISTS user_groups;
DROP TABLE IF EXISTS groups;
DROP TABLE IF EXISTS users;
