# Database Migrations

This directory contains PostgreSQL database migrations for the OpsCore application.

## Overview

The migrations are managed using [golang-migrate/migrate](https://github.com/golang-migrate/migrate) and implement the schema defined in ADR 0005 (Database Schema) and ADR 0006 (Database Migration Strategy).

## Migration Files

Each migration consists of two files:
- `{version}_{description}.up.sql` - SQL statements to apply the migration
- `{version}_{description}.down.sql` - SQL statements to rollback the migration

### Available Migrations

#### 000001: Create Initial Tables
Creates the base `repositories` and `managed_files` tables for repository configuration management.

#### 000002: Add Access Token to Repositories
Adds the `access_token` column to the `repositories` table for GitHub/GitLab authentication.

#### 000003: Document Access Token Encryption
Documents the encryption strategy for access tokens stored in the database.

#### 000004: Create Users and Groups Tables
Creates the user management system:
- `users` - User accounts with name, email, and role (admin/user)
- `groups` - User groups with name and description
- `user_groups` - Many-to-many relationship table for user-group membership

**Schema:**
```sql
users (id, name, email, role, created_at, updated_at)
groups (id, name, description, created_at, updated_at)
user_groups (user_id, group_id, joined_at)
```

**Indexes:**
- `idx_users_email` - Fast lookup by email
- `idx_users_role` - Filter by user role
- `idx_groups_name` - Fast lookup by group name
- `idx_user_groups_user_id` - Lookup users in groups
- `idx_user_groups_group_id` - Lookup groups for users

#### 000005: Create Documents Tables
Creates the document management system:
- `documents` - Published documents with metadata
- `document_versions` - Version history for each document

**Schema:**
```sql
documents (id, repository_id, owner, is_published, is_auto_update, 
           access_scope, current_version_id, created_at, updated_at)
document_versions (id, document_id, version_number, file_path, commit_hash,
                   title, doc_type, tags, variables, content, 
                   published_at, unpublished_at, created_at)
```

**Features:**
- Version control with sequential version numbers
- Support for procedure and knowledge document types
- Tag-based categorization (using PostgreSQL array type)
- Variable definitions stored as JSONB
- Circular reference handling for current_version_id

**Indexes:**
- `idx_documents_repository_id` - Filter documents by repository
- `idx_documents_is_published` - Filter published documents
- `idx_documents_access_scope` - Filter by access level
- `idx_document_versions_document_id` - Lookup document versions
- `idx_document_versions_doc_type` - Filter by document type
- `idx_document_versions_tags` - GIN index for tag searches
- `idx_document_versions_commit_hash` - Lookup by Git commit

#### 000006: Create Execution Records Tables
Creates the execution tracking system:
- `execution_records` - Record of procedure executions
- `execution_steps` - Individual steps within an execution
- `attachments` - Files attached to execution steps

**Schema:**
```sql
execution_records (id, document_id, document_version_id, executor_id,
                   title, variable_values, notes, status, access_scope,
                   started_at, completed_at, created_at, updated_at)
execution_steps (id, execution_record_id, step_number, description,
                 notes, executed_at)
attachments (id, execution_record_id, execution_step_id, file_name,
             file_size, mime_type, storage_type, storage_path,
             uploaded_by, uploaded_at)
```

**Features:**
- Track execution status (in_progress/completed/failed)
- Store variable values used in execution (JSONB)
- Support multiple storage backends (local/s3/minio)
- Associate attachments with execution steps

**Indexes:**
- `idx_execution_records_document_id` - Filter by document
- `idx_execution_records_executor_id` - Filter by executor
- `idx_execution_records_status` - Filter by status
- `idx_execution_records_started_at` - Sort by start time
- `idx_execution_records_completed_at` - Sort by completion time
- `idx_execution_records_variable_values` - GIN index for variable searches
- `idx_execution_steps_execution_record_id` - Lookup execution steps
- `idx_attachments_execution_record_id` - Lookup attachments by execution
- `idx_attachments_execution_step_id` - Lookup attachments by step
- `idx_attachments_uploaded_by` - Filter by uploader

#### 000007: Create View Tables
Creates the view tracking and statistics system:
- `view_history` - Record of document views
- `view_statistics` - Aggregated view statistics per document

**Schema:**
```sql
view_history (id, document_id, user_id, ip_address, user_agent, viewed_at)
view_statistics (document_id, total_views, unique_users, 
                 last_viewed_at, updated_at)
```

**Features:**
- Track anonymous and authenticated views
- Store IP address and user agent for analytics
- Maintain aggregated statistics for performance
- Support for NULL user_id (anonymous views)

**Indexes:**
- `idx_view_history_document_id` - Filter by document
- `idx_view_history_user_id` - Filter by user
- `idx_view_history_viewed_at` - Sort by view time

## Usage

### Using the Migration Tool

The migration tool is located at `backend/cmd/migrate/main.go`.

```bash
# Apply all pending migrations
go run ./cmd/migrate up

# Check migration status
go run ./cmd/migrate status

# Rollback the last migration
go run ./cmd/migrate down

# Rollback N migrations
go run ./cmd/migrate down 2

# Force set migration version (use with caution)
go run ./cmd/migrate force 5
```

### Environment Variables

Set the `DATABASE_URL` environment variable to specify the database connection:

```bash
export DATABASE_URL="postgres://username:password@host:port/database?sslmode=disable"
```

Default: `postgres://opscore_user:opscore_password@db:5432/opscore_db?sslmode=disable`

### Using golang-migrate CLI

If you have the `migrate` CLI installed:

```bash
# Apply all migrations
migrate -database "${DATABASE_URL}" -path ./migrations up

# Check status
migrate -database "${DATABASE_URL}" -path ./migrations version

# Rollback one migration
migrate -database "${DATABASE_URL}" -path ./migrations down 1
```

## Testing

Migration tests are located in `migrations_test.go`.

```bash
# Run migration tests
cd backend/internal/git_repository/infrastructure/persistence
go test -v -run TestMigrations

# Run constraint tests
go test -v -run TestMigrationConstraints
```

Tests verify:
- Migration up/down operations
- Schema integrity (tables exist)
- Index creation
- Foreign key constraints
- Check constraints
- Cascade deletion behavior
- Migration idempotency

## Best Practices

### Creating New Migrations

1. Use sequential version numbers:
   ```bash
   000008_description.up.sql
   000008_description.down.sql
   ```

2. Keep migrations focused on a single change

3. Always create both up and down migrations

4. Test migrations thoroughly before applying to production

5. Never modify applied migrations - create new ones instead

### Writing Migration SQL

1. Use explicit column types and constraints

2. Add indexes for frequently queried columns

3. Include comments for complex logic

4. Use `IF EXISTS` for down migrations to make them idempotent

5. Consider performance impact on large tables

### Rollback Safety

- Always test down migrations locally
- Some operations may cause data loss (document in comments)
- Use transactions when possible (most DDL is transactional in PostgreSQL)
- Consider backup before applying migrations in production

## Schema Diagram

```
repositories ─┐
              ├─> documents ─┐
              │              ├─> document_versions
              │              ├─> execution_records ─┐
              │              │                      ├─> execution_steps
              │              │                      └─> attachments
              │              └─> view_history
              └─> managed_files
                                view_statistics

users ─┐
       ├─> user_groups ─── groups
       ├─> execution_records
       └─> attachments
```

## References

- [ADR 0005: Database Schema](../../../../../adr/0005-database-schema.md)
- [ADR 0006: Database Migration Strategy](../../../../../adr/0006-database-migration.md)
- [golang-migrate Documentation](https://github.com/golang-migrate/migrate)
- [PostgreSQL Documentation](https://www.postgresql.org/docs/)
