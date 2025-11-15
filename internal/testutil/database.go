package testutil

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"testing"
	"time"

	_ "github.com/lib/pq" // PostgreSQL driver
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"

	"github.com/raihanstark/trade-journal/internal/db"
)

// PostgresContainer holds the test database container and connection
type PostgresContainer struct {
	Container *postgres.PostgresContainer
	DB        *sql.DB
	Queries   *db.Queries
}

// SetupTestDatabase creates a PostgreSQL container and runs migrations
func SetupTestDatabase(t *testing.T) *PostgresContainer {
	t.Helper()

	ctx := context.Background()

	pgContainer, err := postgres.Run(ctx,
		"postgres:16-alpine", // Use PostgreSQL 16
		postgres.WithDatabase("testdb"),
		postgres.WithUsername("testuser"),
		postgres.WithPassword("testpass"),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(5*time.Second),
		),
	)
	if err != nil {
		t.Fatalf("failed to start postgres container: %v", err)
	}

	connStr, err := pgContainer.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		t.Fatalf("failed to get connection string: %v", err)
	}

	database, err := sql.Open("postgres", connStr)
	if err != nil {
		t.Fatalf("failed to connect to database: %v", err)
	}

	if err := database.Ping(); err != nil {
		t.Fatalf("failed to ping database: %v", err)
	}

	if err := runMigrations(database); err != nil {
		t.Fatalf("failed to run migrations: %v", err)
	}

	queries := db.New(database)

	t.Cleanup(func() {
		database.Close()
		if err := pgContainer.Terminate(ctx); err != nil {
			t.Logf("failed to terminate postgres container: %v", err)
		}
	})

	return &PostgresContainer{
		Container: pgContainer,
		DB:        database,
		Queries:   queries,
	}
}

// runMigrations reads and runs all migration files from db/migrations/
func runMigrations(db *sql.DB) error {
	projectRoot, err := findProjectRoot()
	if err != nil {
		return fmt.Errorf("failed to find project root: %w", err)
	}

	migrationsDir := filepath.Join(projectRoot, "db", "migrations")

	files, err := os.ReadDir(migrationsDir)
	if err != nil {
		return fmt.Errorf("failed to read migrations directory: %w", err)
	}

	var migrationFiles []string
	for _, file := range files {
		if !file.IsDir() && strings.HasSuffix(file.Name(), ".sql") {
			migrationFiles = append(migrationFiles, file.Name())
		}
	}

	sort.Strings(migrationFiles)

	for _, filename := range migrationFiles {
		filePath := filepath.Join(migrationsDir, filename)

		content, err := os.ReadFile(filePath)
		if err != nil {
			return fmt.Errorf("failed to read migration file %s: %w", filename, err)
		}

		upSQL := extractUpMigration(string(content))

		if _, err := db.Exec(upSQL); err != nil {
			return fmt.Errorf("failed to execute migration %s: %w", filename, err)
		}
	}

	return nil
}

func findProjectRoot() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}

	for {
		goModPath := filepath.Join(dir, "go.mod")
		if _, err := os.Stat(goModPath); err == nil {
			return dir, nil
		}

		parent := filepath.Dir(dir)
		if parent == dir {
			return "", fmt.Errorf("could not find go.mod in any parent directory")
		}
		dir = parent
	}
}

// extractUpMigration extracts the "migrate:up" section from a dbmate migration file
// dbmate format:
// -- migrate:up
// CREATE TABLE ...
// -- migrate:down
// DROP TABLE ...
func extractUpMigration(content string) string {
	// Split by the migrate:down marker
	re := regexp.MustCompile(`(?m)^--\s*migrate:down`)
	parts := re.Split(content, 2)

	if len(parts) == 0 {
		return content
	}

	// Take the first part (before migrate:down)
	upSection := parts[0]

	// Remove the migrate:up marker
	upRe := regexp.MustCompile(`(?m)^--\s*migrate:up\s*\n`)
	upSection = upRe.ReplaceAllString(upSection, "")

	return strings.TrimSpace(upSection)
}

// TruncateTables clears all data from tables
// Use this between tests to ensure isolation
func TruncateTables(t *testing.T, db *sql.DB) {
	t.Helper()

	tables := []string{
		"trade_strategies",
		"trades",
		"strategies",
		"accounts",
		"users",
	}

	for _, table := range tables {
		_, err := db.Exec(fmt.Sprintf("TRUNCATE TABLE %s RESTART IDENTITY CASCADE", table))
		if err != nil {
			t.Fatalf("failed to truncate table %s: %v", table, err)
		}
	}
}
