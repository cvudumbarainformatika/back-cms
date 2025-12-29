package database

import (
	"fmt"
	"io/ioutil"
	"log"
	"path/filepath"
	"sort"
	"strings"

	"github.com/jmoiron/sqlx"
)

// Migrator handles database migrations
type Migrator struct {
	db            *sqlx.DB
	migrationsDir string
}

// NewMigrator creates a new migrator instance
func NewMigrator(db *sqlx.DB, migrationsDir string) *Migrator {
	return &Migrator{
		db:            db,
		migrationsDir: migrationsDir,
	}
}

// RunMigrations runs all migration files in order
func (m *Migrator) RunMigrations() error {
	// Create migrations table if not exists
	if err := m.createMigrationsTable(); err != nil {
		return fmt.Errorf("failed to create migrations table: %w", err)
	}

	// Get all migration files
	files, err := ioutil.ReadDir(m.migrationsDir)
	if err != nil {
		return fmt.Errorf("failed to read migrations directory: %w", err)
	}

	// Filter and sort SQL files
	var sqlFiles []string
	for _, file := range files {
		if !file.IsDir() && strings.HasSuffix(file.Name(), ".sql") {
			// Skip template files
			if !strings.Contains(file.Name(), "TEMPLATE") && !strings.Contains(file.Name(), "MIGRATION_TEMPLATE") {
				sqlFiles = append(sqlFiles, file.Name())
			}
		}
	}
	sort.Strings(sqlFiles)

	if len(sqlFiles) == 0 {
		log.Println("No migration files found")
		return nil
	}

	log.Printf("Found %d migration files\n", len(sqlFiles))

	// Run each migration
	for _, filename := range sqlFiles {
		// Check if migration already ran
		migrated, err := m.isMigrationRun(filename)
		if err != nil {
			return fmt.Errorf("failed to check migration status: %w", err)
		}

		if migrated {
			log.Printf("⊘ Skipping migration (already run): %s\n", filename)
			continue
		}

		// Read and execute migration file
		filePath := filepath.Join(m.migrationsDir, filename)
		content, err := ioutil.ReadFile(filePath)
		if err != nil {
			return fmt.Errorf("failed to read migration file %s: %w", filename, err)
		}

		// Split by semicolon to handle multiple statements
		queries := strings.Split(string(content), ";")
		for _, query := range queries {
			query = strings.TrimSpace(query)
			if query == "" {
				continue
			}

			if _, err := m.db.Exec(query); err != nil {
				return fmt.Errorf("failed to execute migration %s: %w", filename, err)
			}
		}

		// Record migration as run
		if err := m.recordMigration(filename); err != nil {
			return fmt.Errorf("failed to record migration: %w", err)
		}

		log.Printf("✓ Migration completed: %s\n", filename)
	}

	log.Println("✓ All migrations completed successfully!")
	return nil
}

// createMigrationsTable creates the migrations tracking table
func (m *Migrator) createMigrationsTable() error {
	query := `
	CREATE TABLE IF NOT EXISTS migrations (
		id INT AUTO_INCREMENT PRIMARY KEY,
		migration VARCHAR(255) NOT NULL UNIQUE,
		executed_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
	`
	_, err := m.db.Exec(query)
	return err
}

// isMigrationRun checks if a migration has already been run
func (m *Migrator) isMigrationRun(filename string) (bool, error) {
	var count int
	query := `SELECT COUNT(*) FROM migrations WHERE migration = ?`
	err := m.db.Get(&count, query, filename)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// recordMigration records a migration as run
func (m *Migrator) recordMigration(filename string) error {
	query := `INSERT INTO migrations (migration) VALUES (?)`
	_, err := m.db.Exec(query, filename)
	return err
}
