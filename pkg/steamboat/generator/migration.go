package generator

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	_ "github.com/joho/godotenv/autoload"
)

// GenerateMigration creates a new SQL migration file pair
func GenerateMigration(name string) (string, error) {
	// Find the next migration number
	migrationNum, err := getNextMigrationNumber()
	if err != nil {
		return "", fmt.Errorf("failed to get next migration number: %w", err)
	}
	
	// Create migration file names
	migrationName := toSnakeCase(name)
	upFile := fmt.Sprintf("%06d_%s.up.sql", migrationNum, migrationName)
	downFile := fmt.Sprintf("%06d_%s.down.sql", migrationNum, migrationName)
	
	// Use hardcoded migrations path
	migrationsDir := "internal/database/migrations"
	
	// Ensure directory exists
	if err := os.MkdirAll(migrationsDir, 0755); err != nil {
		return "", fmt.Errorf("failed to create migrations directory: %w", err)
	}
	
	upPath := filepath.Join(migrationsDir, upFile)
	downPath := filepath.Join(migrationsDir, downFile)
	
	// Create up migration file
	upContent := fmt.Sprintf(`-- Migration: %s
-- Created: %s

-- Add your SQL here

`, toHumanReadableMigration(name), fmt.Sprintf("%06d", migrationNum))
	
	if err := os.WriteFile(upPath, []byte(upContent), 0644); err != nil {
		return "", fmt.Errorf("failed to create up migration file: %w", err)
	}
	
	// Create down migration file
	downContent := fmt.Sprintf(`-- Rollback: %s
-- Created: %s

-- Add your rollback SQL here

`, toHumanReadableMigration(name), fmt.Sprintf("%06d", migrationNum))
	
	if err := os.WriteFile(downPath, []byte(downContent), 0644); err != nil {
		return "", fmt.Errorf("failed to create down migration file: %w", err)
	}
	
	return fmt.Sprintf("Migration files created:\n  - %s\n  - %s", upPath, downPath), nil
}

func getNextMigrationNumber() (int, error) {
	migrationsDir := "internal/database/migrations"
	
	// Ensure directory exists
	if err := os.MkdirAll(migrationsDir, 0755); err != nil {
		return 0, err
	}
	
	entries, err := os.ReadDir(migrationsDir)
	if err != nil {
		return 0, err
	}
	
	maxNum := 0
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		
		name := entry.Name()
		if !strings.HasSuffix(name, ".sql") {
			continue
		}
		
		// Extract number from filename (e.g., "000001_create_users.up.sql" -> 1)
		parts := strings.Split(name, "_")
		if len(parts) < 2 {
			continue
		}
		
		if num, err := strconv.Atoi(parts[0]); err == nil {
			if num > maxNum {
				maxNum = num
			}
		}
	}
	
	return maxNum + 1, nil
}

func toHumanReadableMigration(s string) string {
	// Convert snake_case or camelCase to human readable
	s = strings.ReplaceAll(s, "_", " ")
	words := strings.Fields(s)
	for i, word := range words {
		words[i] = strings.Title(strings.ToLower(word))
	}
	return strings.Join(words, " ")
}