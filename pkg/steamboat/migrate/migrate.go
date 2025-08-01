package migrate

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/sqlite3"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/joho/godotenv/autoload"
)

// Run executes all pending migrations
func Run() error {
	m, err := getMigrator()
	if err != nil {
		return err
	}
	defer m.Close()

	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return fmt.Errorf("failed to run migrations: %w", err)
	}

	return nil
}

// Rollback rolls back one migration
func Rollback() error {
	m, err := getMigrator()
	if err != nil {
		return err
	}
	defer m.Close()

	if err := m.Steps(-1); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return fmt.Errorf("failed to rollback migration: %w", err)
	}

	return nil
}

// Status returns the current migration version
func Status() (uint, bool, error) {
	m, err := getMigrator()
	if err != nil {
		return 0, false, err
	}
	defer m.Close()

	version, dirty, err := m.Version()
	if err != nil && !errors.Is(err, migrate.ErrNilVersion) {
		return 0, false, fmt.Errorf("failed to get migration status: %w", err)
	}

	return version, dirty, nil
}

func getMigrator() (*migrate.Migrate, error) {
	// Get database URL from environment
	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		return nil, fmt.Errorf("DB_URL environment variable is not set")
	}

	// Ensure the database directory exists
	dbDir := filepath.Dir(dbURL)
	if dbDir != "." && dbDir != "" {
		if err := os.MkdirAll(dbDir, 0755); err != nil {
			return nil, fmt.Errorf("failed to create database directory: %w", err)
		}
	}

	// Use hardcoded migrations path
	migrationsPath := "internal/database/migrations"

	// Ensure migrations directory exists
	if err := os.MkdirAll(migrationsPath, 0755); err != nil {
		return nil, fmt.Errorf("failed to create migrations directory: %w", err)
	}

	// Convert to absolute path for file:// URL
	absPath, err := filepath.Abs(migrationsPath)
	if err != nil {
		return nil, fmt.Errorf("failed to get absolute path: %w", err)
	}

	// Convert to SQLite URL format
	sqliteURL := fmt.Sprintf("sqlite3://%s", dbURL)
	migrationsURL := fmt.Sprintf("file://%s", absPath)

	// Create migrator
	m, err := migrate.New(migrationsURL, sqliteURL)
	if err != nil {
		return nil, fmt.Errorf("failed to create migrator: %w", err)
	}

	return m, nil
}