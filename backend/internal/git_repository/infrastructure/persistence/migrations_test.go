package persistence

import (
	"context"
	"fmt"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	// rollbackOneStep is the number of migration steps to rollback in tests
	rollbackOneStep = -1
)

// checkTableExists is a helper function to check if a table exists in the database
func checkTableExists(t *testing.T, ctx context.Context, conn *pgxpool.Pool, tableName string) bool {
	var exists bool
	err := conn.QueryRow(ctx, `
		SELECT EXISTS (
			SELECT FROM information_schema.tables 
			WHERE table_schema = 'public' 
			AND table_name = $1
		)
	`, tableName).Scan(&exists)
	require.NoError(t, err, "Failed to check if table %s exists", tableName)
	return exists
}

// TestMigrations tests the database migrations
func TestMigrations(t *testing.T) {
	// Skip test if PostgreSQL is not available
	if !checkDatabaseConnection(t) {
		t.Skip("Skipping migrations test - database is not available")
		return
	}

	config := getPostgreSQLConfig()

	// Generate unique test database name
	testDBName := fmt.Sprintf("%s_migrations_%s", config.DBName, uuid.New().String()[:8])

	// Connect to master database to create test database
	masterDSN := fmt.Sprintf("postgresql://%s:%s@%s:%s/postgres",
		config.User, config.Password, config.Host, config.Port)

	masterConn, err := pgxpool.New(context.Background(), masterDSN)
	require.NoError(t, err, "Failed to connect to PostgreSQL master database")
	defer masterConn.Close()

	// Create test database
	// Note: pgx.Identifier{}.Sanitize() properly quotes and escapes identifiers,
	// preventing SQL injection by treating the name as a literal identifier.
	// Example: "test; DROP TABLE users--" becomes "test; DROP TABLE users--" (quoted)
	createSQL := fmt.Sprintf("CREATE DATABASE %s", pgx.Identifier{testDBName}.Sanitize())
	_, err = masterConn.Exec(context.Background(), createSQL)
	require.NoError(t, err, "Failed to create test database")

	// Cleanup function
	defer func() {
		cleanupConn, err := pgxpool.New(context.Background(), masterDSN)
		if err == nil {
			dropSQL := fmt.Sprintf("DROP DATABASE IF EXISTS %s WITH (FORCE)", pgx.Identifier{testDBName}.Sanitize())
			_, _ = cleanupConn.Exec(context.Background(), dropSQL)
			cleanupConn.Close()
		}
	}()

	// Connect to test database
	testDSN := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s",
		config.User, config.Password, config.Host, config.Port, testDBName)

	testConn, err := pgxpool.New(context.Background(), testDSN)
	require.NoError(t, err, "Failed to connect to test database")
	defer testConn.Close()

	// Get migrations directory path
	_, filename, _, ok := runtime.Caller(0)
	require.True(t, ok, "Failed to get current file path")
	migrationsPath := filepath.Join(filepath.Dir(filename), "migrations")

	// Test migration up
	t.Run("Migration Up", func(t *testing.T) {
		sourceURL := fmt.Sprintf("file://%s", migrationsPath)
		m, err := migrate.New(sourceURL, testDSN)
		require.NoError(t, err, "Failed to create migrate instance")
		defer m.Close()

		err = m.Up()
		require.NoError(t, err, "Migration up failed")

		// Verify migration version
		version, dirty, err := m.Version()
		require.NoError(t, err, "Failed to get migration version")
		assert.False(t, dirty, "Database should not be in dirty state")
		assert.Equal(t, uint(7), version, "Expected migration version 7")
	})

	// Verify schema after migrations
	t.Run("Schema Verification", func(t *testing.T) {
		ctx := context.Background()

		// Check all required tables exist
		tables := []string{
			"repositories",
			"users",
			"groups",
			"user_groups",
			"documents",
			"document_versions",
			"execution_records",
			"execution_steps",
			"attachments",
			"view_history",
			"view_statistics",
		}

		for _, tableName := range tables {
			exists := checkTableExists(t, ctx, testConn, tableName)
			assert.True(t, exists, "%s table should exist", tableName)
		}
	})

	// Verify indexes
	t.Run("Index Verification", func(t *testing.T) {
		ctx := context.Background()

		// Check for important indexes
		indexes := []string{
			"idx_users_email",
			"idx_users_role",
			"idx_groups_name",
			"idx_user_groups_user_id",
			"idx_user_groups_group_id",
			"idx_documents_repository_id",
			"idx_documents_is_published",
			"idx_documents_access_scope",
			"idx_document_versions_document_id",
			"idx_document_versions_doc_type",
			"idx_document_versions_tags",
			"idx_document_versions_commit_hash",
			"idx_execution_records_document_id",
			"idx_execution_records_executor_id",
			"idx_execution_records_status",
			"idx_execution_steps_execution_record_id",
			"idx_attachments_execution_record_id",
			"idx_attachments_execution_step_id",
			"idx_view_history_document_id",
			"idx_view_history_user_id",
			"idx_view_history_viewed_at",
		}

		for _, indexName := range indexes {
			var indexExists bool
			err := testConn.QueryRow(ctx, `
				SELECT EXISTS (
					SELECT FROM pg_indexes 
					WHERE schemaname = 'public' 
					AND indexname = $1
				)
			`, indexName).Scan(&indexExists)
			require.NoError(t, err, "Failed to check index %s", indexName)
			assert.True(t, indexExists, "Index %s should exist", indexName)
		}
	})

	// Test migration down
	t.Run("Migration Down", func(t *testing.T) {
		sourceURL := fmt.Sprintf("file://%s", migrationsPath)
		m, err := migrate.New(sourceURL, testDSN)
		require.NoError(t, err, "Failed to create migrate instance")
		defer m.Close()

		// Rollback one migration
		err = m.Steps(rollbackOneStep)
		require.NoError(t, err, "Migration down failed")

		// Verify version decreased
		version, dirty, err := m.Version()
		require.NoError(t, err, "Failed to get migration version")
		assert.False(t, dirty, "Database should not be in dirty state")
		assert.Equal(t, uint(6), version, "Expected migration version 6 after rollback")

		// Verify view_history and view_statistics tables are dropped
		ctx := context.Background()
		var viewHistoryTableExists bool
		err = testConn.QueryRow(ctx, `
			SELECT EXISTS (
				SELECT FROM information_schema.tables 
				WHERE table_schema = 'public' 
				AND table_name = 'view_history'
			)
		`).Scan(&viewHistoryTableExists)
		require.NoError(t, err)
		assert.False(t, viewHistoryTableExists, "view_history table should be dropped")

		// Migrate back up
		err = m.Up()
		require.NoError(t, err, "Migration up failed after rollback")

		version, dirty, err = m.Version()
		require.NoError(t, err, "Failed to get migration version")
		assert.False(t, dirty, "Database should not be in dirty state")
		assert.Equal(t, uint(7), version, "Expected migration version 7 after re-applying")
	})

	// Test migration idempotency
	t.Run("Migration Idempotency", func(t *testing.T) {
		sourceURL := fmt.Sprintf("file://%s", migrationsPath)
		m, err := migrate.New(sourceURL, testDSN)
		require.NoError(t, err, "Failed to create migrate instance")
		defer m.Close()

		// Try to apply migrations again
		err = m.Up()
		assert.ErrorIs(t, err, migrate.ErrNoChange, "Should return ErrNoChange when no migrations to apply")
	})
}

// TestMigrationConstraints tests foreign key constraints and check constraints
func TestMigrationConstraints(t *testing.T) {
	// Skip test if PostgreSQL is not available
	if !checkDatabaseConnection(t) {
		t.Skip("Skipping constraints test - database is not available")
		return
	}

	config := getPostgreSQLConfig()
	testDBName := fmt.Sprintf("%s_constraints_%s", config.DBName, uuid.New().String()[:8])

	// Setup test database
	masterDSN := fmt.Sprintf("postgresql://%s:%s@%s:%s/postgres",
		config.User, config.Password, config.Host, config.Port)

	masterConn, err := pgxpool.New(context.Background(), masterDSN)
	require.NoError(t, err)
	defer masterConn.Close()

	// Create test database (using pgx.Identifier for safe quoting)
	createSQL := fmt.Sprintf("CREATE DATABASE %s", pgx.Identifier{testDBName}.Sanitize())
	_, err = masterConn.Exec(context.Background(), createSQL)
	require.NoError(t, err)

	defer func() {
		cleanupConn, err := pgxpool.New(context.Background(), masterDSN)
		if err == nil {
			dropSQL := fmt.Sprintf("DROP DATABASE IF EXISTS %s WITH (FORCE)", pgx.Identifier{testDBName}.Sanitize())
			_, _ = cleanupConn.Exec(context.Background(), dropSQL)
			cleanupConn.Close()
		}
	}()

	testDSN := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s",
		config.User, config.Password, config.Host, config.Port, testDBName)

	testConn, err := pgxpool.New(context.Background(), testDSN)
	require.NoError(t, err)
	defer testConn.Close()

	// Apply migrations
	_, filename, _, ok := runtime.Caller(0)
	require.True(t, ok)
	migrationsPath := filepath.Join(filepath.Dir(filename), "migrations")
	sourceURL := fmt.Sprintf("file://%s", migrationsPath)

	m, err := migrate.New(sourceURL, testDSN)
	require.NoError(t, err)
	defer m.Close()

	err = m.Up()
	require.NoError(t, err)

	ctx := context.Background()

	// Test users table constraints
	t.Run("Users Role Constraint", func(t *testing.T) {
		// Try to insert user with invalid role
		_, err := testConn.Exec(ctx, `
			INSERT INTO users (name, email, role)
			VALUES ($1, $2, $3)
		`, "Test User", "test@example.com", "invalid_role")
		assert.Error(t, err, "Should fail with invalid role")
	})

	// Test documents table constraints
	t.Run("Documents Access Scope Constraint", func(t *testing.T) {
		// First create a repository
		var repoID string
		err := testConn.QueryRow(ctx, `
			INSERT INTO repositories (name, url)
			VALUES ($1, $2)
			RETURNING id
		`, "Test Repo", "https://github.com/test/repo").Scan(&repoID)
		require.NoError(t, err)

		// Try to insert document with invalid access_scope
		_, err = testConn.Exec(ctx, `
			INSERT INTO documents (repository_id, owner, access_scope)
			VALUES ($1, $2, $3)
		`, repoID, "test@example.com", "invalid_scope")
		assert.Error(t, err, "Should fail with invalid access_scope")
	})

	// Test execution_records table constraints
	t.Run("ExecutionRecords Status Constraint", func(t *testing.T) {
		// Create necessary parent records
		var userID string
		err := testConn.QueryRow(ctx, `
			INSERT INTO users (name, email)
			VALUES ($1, $2)
			RETURNING id
		`, "Executor", "executor@example.com").Scan(&userID)
		require.NoError(t, err)

		var repoID string
		err = testConn.QueryRow(ctx, `
			INSERT INTO repositories (name, url)
			VALUES ($1, $2)
			RETURNING id
		`, "Exec Repo", "https://github.com/test/exec").Scan(&repoID)
		require.NoError(t, err)

		var docID string
		err = testConn.QueryRow(ctx, `
			INSERT INTO documents (repository_id, owner, access_scope)
			VALUES ($1, $2, $3)
			RETURNING id
		`, repoID, "test@example.com", "public").Scan(&docID)
		require.NoError(t, err)

		var versionID string
		err = testConn.QueryRow(ctx, `
			INSERT INTO document_versions (document_id, version_number, file_path, commit_hash, title, doc_type, content, published_at)
			VALUES ($1, $2, $3, $4, $5, $6, $7, NOW())
			RETURNING id
		`, docID, 1, "/test.md", "abc123", "Test Doc", "procedure", "Test content").Scan(&versionID)
		require.NoError(t, err)

		// Try to insert execution record with invalid status
		_, err = testConn.Exec(ctx, `
			INSERT INTO execution_records (document_id, document_version_id, executor_id, title, status, access_scope)
			VALUES ($1, $2, $3, $4, $5, $6)
		`, docID, versionID, userID, "Test Execution", "invalid_status", "public")
		assert.Error(t, err, "Should fail with invalid status")
	})

	// Test foreign key constraints
	t.Run("Foreign Key Constraints", func(t *testing.T) {
		// Try to insert document with non-existent repository_id
		fakeRepoID := uuid.New().String()
		_, err := testConn.Exec(ctx, `
			INSERT INTO documents (repository_id, owner, access_scope)
			VALUES ($1, $2, $3)
		`, fakeRepoID, "test@example.com", "public")
		assert.Error(t, err, "Should fail with non-existent repository_id")
	})

	// Test cascade deletion
	t.Run("Cascade Deletion", func(t *testing.T) {
		// Create user and group
		var userID string
		err := testConn.QueryRow(ctx, `
			INSERT INTO users (name, email)
			VALUES ($1, $2)
			RETURNING id
		`, "Cascade User", "cascade@example.com").Scan(&userID)
		require.NoError(t, err)

		var groupID string
		err = testConn.QueryRow(ctx, `
			INSERT INTO groups (name)
			VALUES ($1)
			RETURNING id
		`, "Cascade Group").Scan(&groupID)
		require.NoError(t, err)

		// Add user to group
		_, err = testConn.Exec(ctx, `
			INSERT INTO user_groups (user_id, group_id)
			VALUES ($1, $2)
		`, userID, groupID)
		require.NoError(t, err)

		// Delete user and verify user_groups record is also deleted
		_, err = testConn.Exec(ctx, "DELETE FROM users WHERE id = $1", userID)
		require.NoError(t, err)

		var count int
		err = testConn.QueryRow(ctx, "SELECT COUNT(*) FROM user_groups WHERE user_id = $1", userID).Scan(&count)
		require.NoError(t, err)
		assert.Equal(t, 0, count, "user_groups record should be deleted via CASCADE")
	})
}
