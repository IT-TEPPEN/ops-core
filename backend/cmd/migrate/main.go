package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	command := os.Args[1]

	// Get database URL from environment variable
	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		databaseURL = "postgres://opscore_user:opscore_password@db:5432/opscore_db?sslmode=disable"
		fmt.Println("Warning: DATABASE_URL environment variable not set, using default.")
	}

	// Verify database connection
	ctx := context.Background()
	dbpool, err := pgxpool.New(ctx, databaseURL)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to create connection pool: %v\n", err)
		os.Exit(1)
	}
	defer dbpool.Close()

	if err := dbpool.Ping(ctx); err != nil {
		fmt.Fprintf(os.Stderr, "Unable to ping database: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("Successfully connected to database")

	// Get migrations directory path
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		fmt.Fprintf(os.Stderr, "Failed to get current file path\n")
		os.Exit(1)
	}
	migrationsPath := filepath.Join(filepath.Dir(filename), "..", "..", "internal", "git_repository", "infrastructure", "persistence", "migrations")
	absPath, err := filepath.Abs(migrationsPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to get absolute path: %v\n", err)
		os.Exit(1)
	}

	// Create migrate instance
	sourceURL := fmt.Sprintf("file://%s", absPath)
	m, err := migrate.New(sourceURL, databaseURL)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create migrate instance: %v\n", err)
		os.Exit(1)
	}
	defer m.Close()

	// Execute command
	switch command {
	case "up":
		if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
			fmt.Fprintf(os.Stderr, "Migration up failed: %v\n", err)
			os.Exit(1)
		}
		if errors.Is(err, migrate.ErrNoChange) {
			fmt.Println("No migrations to apply")
		} else {
			fmt.Println("Migrations applied successfully")
		}

	case "down":
		steps := 1
		if len(os.Args) > 2 {
			fmt.Sscanf(os.Args[2], "%d", &steps)
		}
		if err := m.Steps(-steps); err != nil && !errors.Is(err, migrate.ErrNoChange) {
			fmt.Fprintf(os.Stderr, "Migration down failed: %v\n", err)
			os.Exit(1)
		}
		if errors.Is(err, migrate.ErrNoChange) {
			fmt.Println("No migrations to rollback")
		} else {
			fmt.Printf("Rolled back %d migration(s) successfully\n", steps)
		}

	case "status":
		version, dirty, err := m.Version()
		if err != nil && !errors.Is(err, migrate.ErrNilVersion) {
			fmt.Fprintf(os.Stderr, "Failed to get migration status: %v\n", err)
			os.Exit(1)
		}
		if errors.Is(err, migrate.ErrNilVersion) {
			fmt.Println("No migrations have been applied yet")
		} else {
			fmt.Printf("Current migration version: %d\n", version)
			if dirty {
				fmt.Println("Warning: Database is in a dirty state")
			} else {
				fmt.Println("Database is in a clean state")
			}
		}

	case "force":
		if len(os.Args) < 3 {
			fmt.Fprintf(os.Stderr, "Force command requires a version number\n")
			printUsage()
			os.Exit(1)
		}
		var version int
		fmt.Sscanf(os.Args[2], "%d", &version)
		if err := m.Force(version); err != nil {
			fmt.Fprintf(os.Stderr, "Failed to force version: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("Forced migration version to %d\n", version)

	default:
		fmt.Fprintf(os.Stderr, "Unknown command: %s\n", command)
		printUsage()
		os.Exit(1)
	}
}

func printUsage() {
	fmt.Println("Usage: migrate <command> [options]")
	fmt.Println("\nCommands:")
	fmt.Println("  up              - Apply all pending migrations")
	fmt.Println("  down [N]        - Rollback N migrations (default: 1)")
	fmt.Println("  status          - Show current migration status")
	fmt.Println("  force <version> - Force set migration version (use with caution)")
	fmt.Println("\nEnvironment Variables:")
	fmt.Println("  DATABASE_URL    - PostgreSQL connection string")
}
