package testutil

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
)

const defaultTestDBURL = "postgres://supercamp:supercamp@localhost:5435/supercamp_test?sslmode=disable"

// TestDBURL returns the test database URL from env or default.
func TestDBURL() string {
	if url := os.Getenv("TEST_DATABASE_URL"); url != "" {
		return url
	}
	return defaultTestDBURL
}

// SetupTestDB creates a connection pool and applies migrations.
// It returns the pool and a cleanup function.
func SetupTestDB(t *testing.T) *pgxpool.Pool {
	t.Helper()

	ctx := context.Background()
	pool, err := pgxpool.New(ctx, TestDBURL())
	if err != nil {
		t.Fatalf("connect to test db: %v", err)
	}

	if err := pool.Ping(ctx); err != nil {
		t.Fatalf("ping test db: %v", err)
	}

	// Apply schema
	if err := applyMigrations(ctx, pool); err != nil {
		t.Fatalf("apply migrations: %v", err)
	}

	t.Cleanup(func() {
		cleanTables(ctx, pool)
		pool.Close()
	})

	return pool
}

func applyMigrations(ctx context.Context, pool *pgxpool.Pool) error {
	// Read and execute migration files
	migrations := []string{
		"migrations/000001_init_schema.up.sql",
		"migrations/000002_seed_user.up.sql",
	}

	for _, path := range migrations {
		sql, err := os.ReadFile(path)
		if err != nil {
			// Try from backend/ prefix (when running from project root)
			sql, err = os.ReadFile("../../" + path)
			if err != nil {
				return fmt.Errorf("read migration %s: %w", path, err)
			}
		}
		if _, err := pool.Exec(ctx, string(sql)); err != nil {
			// Ignore errors from IF NOT EXISTS / ON CONFLICT
			// This allows re-running migrations
			continue
		}
	}
	return nil
}

// CleanTables truncates all application tables.
func cleanTables(ctx context.Context, pool *pgxpool.Pool) {
	tables := []string{
		"checklist_items",
		"checklists",
		"layouts",
		"fire_logs",
		"meal_plans",
	}
	for _, table := range tables {
		pool.Exec(ctx, fmt.Sprintf("DELETE FROM %s", table))
	}
}

// CleanTablesForTest can be called between subtests to reset state.
func CleanTablesForTest(t *testing.T, pool *pgxpool.Pool) {
	t.Helper()
	cleanTables(context.Background(), pool)
}
