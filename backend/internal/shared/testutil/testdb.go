package testutil

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// TestDatabase holds the resources for an isolated test database.
// Create via SetupTestDatabase and clean up via TearDown.
type TestDatabase struct {
	Pool   *pgxpool.Pool
	DBName string
	dbURL  string // admin connection URL (pointing to "postgres" DB)
}

// SetupTestDatabase creates a new temporary database for testing.
// It connects to the "postgres" admin database, creates a uniquely named DB,
// runs migrations, and returns a connection pool to the new DB.
//
// migrationSteps controls how many migrations to apply.
// If 0 or omitted, all migrations are applied.
//
// The caller must call TearDown() (typically via t.Cleanup or TearDownSuite)
// to drop the database and release resources.
func SetupTestDatabase(t testing.TB, migrationSteps ...int) *TestDatabase {
	t.Helper()

	adminURL := adminDatabaseURL(t)

	ctx := context.Background()

	// Connect to the "postgres" admin database to create the test DB.
	adminConn, err := pgx.Connect(ctx, adminURL)
	if err != nil {
		t.Fatalf("testutil: failed to connect to admin database: %v", err)
	}
	defer adminConn.Close(ctx)

	// Generate a unique database name.
	suffix := strings.ReplaceAll(uuid.New().String(), "-", "")[:12]
	dbName := fmt.Sprintf("opscore_test_%s", suffix)

	// CREATE DATABASE does not support parameters; dbName is generated internally (UUID-based), not from user input.
	_, err = adminConn.Exec(ctx, fmt.Sprintf("CREATE DATABASE %s", pgx.Identifier{dbName}.Sanitize()))
	if err != nil {
		t.Fatalf("testutil: failed to create test database %s: %v", dbName, err)
	}

	// Build connection URL for the new test database.
	testDBURL := replaceDBName(adminURL, dbName)

	// Run migrations.
	migrationsPath := resolveMigrationsPath(t)
	m, err := migrate.New(fmt.Sprintf("file://%s", migrationsPath), testDBURL)
	if err != nil {
		// Clean up: drop the database if migration setup fails.
		dropDatabase(t, adminURL, dbName)
		t.Fatalf("testutil: failed to create migrate instance: %v", err)
	}
	steps := 0
	if len(migrationSteps) > 0 {
		steps = migrationSteps[0]
	}

	var migErr error
	if steps > 0 {
		migErr = m.Steps(steps)
	} else {
		migErr = m.Up()
	}
	if migErr != nil && migErr != migrate.ErrNoChange {
		srcErr, dbErr := m.Close()
		_ = srcErr
		_ = dbErr
		dropDatabase(t, adminURL, dbName)
		t.Fatalf("testutil: failed to run migrations: %v", migErr)
	}
	srcErr, dbErr := m.Close()
	if srcErr != nil {
		t.Logf("testutil: warning closing migration source: %v", srcErr)
	}
	if dbErr != nil {
		t.Logf("testutil: warning closing migration db: %v", dbErr)
	}

	// Create connection pool for the test database.
	pool, err := pgxpool.New(ctx, testDBURL)
	if err != nil {
		dropDatabase(t, adminURL, dbName)
		t.Fatalf("testutil: failed to create connection pool for test database: %v", err)
	}

	td := &TestDatabase{
		Pool:   pool,
		DBName: dbName,
		dbURL:  adminURL,
	}

	t.Cleanup(func() {
		td.TearDown(t)
	})

	return td
}

// TearDown closes the connection pool and drops the test database.
func (td *TestDatabase) TearDown(t testing.TB) {
	t.Helper()

	if td.Pool != nil {
		td.Pool.Close()
		td.Pool = nil
	}

	dropDatabase(t, td.dbURL, td.DBName)
}

// adminDatabaseURL returns a connection URL to the "postgres" admin database
// on the same host as DATABASE_URL. It replaces the database name with "postgres".
func adminDatabaseURL(t testing.TB) string {
	t.Helper()

	// We use DATABASE_URL (the dev connection) only to derive the host/port/credentials.
	// We never connect to opscore_db itself.
	dbURL := getDatabaseURL(t)
	return replaceDBName(dbURL, "postgres")
}

// getDatabaseURL retrieves the database connection URL from environment.
// It checks TEST_DATABASE_URL first, then falls back to DATABASE_URL.
func getDatabaseURL(t testing.TB) string {
	t.Helper()

	// Allow override via TEST_DATABASE_URL for CI or custom setups.
	for _, key := range []string{"TEST_DATABASE_URL", "DATABASE_URL"} {
		if val := lookupEnv(key); val != "" {
			return val
		}
	}

	t.Skip("Neither TEST_DATABASE_URL nor DATABASE_URL is set; skipping integration tests")
	return "" // unreachable, but satisfies compiler
}

// lookupEnv retrieves an environment variable value.
func lookupEnv(key string) string {
	return os.Getenv(key)
}

// resolveMigrationsPath finds the backend/migrate/ directory using runtime.Caller.
func resolveMigrationsPath(t testing.TB) string {
	t.Helper()

	// This file is at backend/internal/shared/testutil/testdb.go
	// We need backend/migrate/
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("testutil: failed to resolve current file path via runtime.Caller")
	}

	// Navigate from testutil/ -> shared/ -> internal/ -> backend/migrate/
	migrationsPath := filepath.Join(filepath.Dir(filename), "..", "..", "..", "migrate")
	absPath, err := filepath.Abs(migrationsPath)
	if err != nil {
		t.Fatalf("testutil: failed to resolve migrations path: %v", err)
	}
	return absPath
}

// replaceDBName replaces the database name in a PostgreSQL connection URL.
// Supports the format: postgres://user:pass@host:port/dbname?params
func replaceDBName(connURL, newDB string) string {
	// Find the last '/' before '?' (query params)
	queryIdx := strings.Index(connURL, "?")
	base := connURL
	query := ""
	if queryIdx >= 0 {
		base = connURL[:queryIdx]
		query = connURL[queryIdx:]
	}

	lastSlash := strings.LastIndex(base, "/")
	if lastSlash < 0 {
		return connURL // Unexpected format; return as-is.
	}
	return base[:lastSlash+1] + newDB + query
}

// dropDatabase connects to the admin DB and drops the specified database.
func dropDatabase(t testing.TB, adminURL, dbName string) {
	t.Helper()

	ctx := context.Background()
	conn, err := pgx.Connect(ctx, adminURL)
	if err != nil {
		t.Logf("testutil: warning: failed to connect for dropping database %s: %v", dbName, err)
		return
	}
	defer conn.Close(ctx)

	// Terminate any remaining connections to the test database.
	_, _ = conn.Exec(ctx, "SELECT pg_terminate_backend(pid) FROM pg_stat_activity WHERE datname = $1 AND pid <> pg_backend_pid()", dbName)

	// DROP DATABASE does not support parameters; dbName is generated internally (UUID-based), not from user input.
	_, err = conn.Exec(ctx, fmt.Sprintf("DROP DATABASE IF EXISTS %s", pgx.Identifier{dbName}.Sanitize()))
	if err != nil {
		t.Logf("testutil: warning: failed to drop test database %s: %v", dbName, err)
	}
}
