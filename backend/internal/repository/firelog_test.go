//go:build integration

package repository

import (
	"context"
	"testing"

	"github.com/jackc/pgx/v5"
	"github.com/sugiyamadaiki/super-camp/backend/internal/testutil"
)

func TestFireLogRepository_CRUD(t *testing.T) {
	pool := testutil.SetupTestDB(t)
	repo := NewFireLogRepository(pool)
	ctx := context.Background()
	var userID int64 = 1

	temp := 25.5

	// Create
	log, err := repo.Create(ctx, userID, "2025-03-15", "河原サイト", "薪(広葉樹)", 90, "良い火", &temp)
	if err != nil {
		t.Fatalf("Create: %v", err)
	}
	if log.Location != "河原サイト" {
		t.Fatalf("expected location '河原サイト', got %q", log.Location)
	}
	if log.Temperature == nil || *log.Temperature != 25.5 {
		t.Fatalf("expected temperature 25.5, got %v", log.Temperature)
	}

	// Create without temperature
	log2, err := repo.Create(ctx, userID, "2025-03-16", "山サイト", "針葉樹", 60, "", nil)
	if err != nil {
		t.Fatalf("Create without temp: %v", err)
	}
	if log2.Temperature != nil {
		t.Fatalf("expected nil temperature, got %v", log2.Temperature)
	}

	// Get
	got, err := repo.Get(ctx, log.ID)
	if err != nil {
		t.Fatalf("Get: %v", err)
	}
	if got.WoodType != "薪(広葉樹)" {
		t.Fatalf("Get wood_type mismatch: %q", got.WoodType)
	}

	// Update
	newTemp := 30.0
	updated, err := repo.Update(ctx, log.ID, "2025-03-15", "湖畔サイト", "薪(広葉樹)", 120, "最高", &newTemp)
	if err != nil {
		t.Fatalf("Update: %v", err)
	}
	if updated.Location != "湖畔サイト" {
		t.Fatalf("expected updated location, got %q", updated.Location)
	}
	if updated.DurationMinutes != 120 {
		t.Fatalf("expected 120 min, got %d", updated.DurationMinutes)
	}

	// List
	list, err := repo.List(ctx, userID)
	if err != nil {
		t.Fatalf("List: %v", err)
	}
	if len(list) != 2 {
		t.Fatalf("expected 2 logs, got %d", len(list))
	}

	// Delete
	if err := repo.Delete(ctx, log.ID); err != nil {
		t.Fatalf("Delete: %v", err)
	}
	if err := repo.Delete(ctx, log2.ID); err != nil {
		t.Fatalf("Delete 2: %v", err)
	}

	_, err = repo.Get(ctx, log.ID)
	if err != pgx.ErrNoRows {
		t.Fatalf("expected ErrNoRows, got %v", err)
	}
}

func TestFireLogRepository_NotFound(t *testing.T) {
	pool := testutil.SetupTestDB(t)
	repo := NewFireLogRepository(pool)
	ctx := context.Background()

	_, err := repo.Get(ctx, 99999)
	if err != pgx.ErrNoRows {
		t.Fatalf("expected ErrNoRows, got %v", err)
	}

	_, err = repo.Update(ctx, 99999, "2025-01-01", "", "", 0, "", nil)
	if err != pgx.ErrNoRows {
		t.Fatalf("expected ErrNoRows, got %v", err)
	}

	err = repo.Delete(ctx, 99999)
	if err != pgx.ErrNoRows {
		t.Fatalf("expected ErrNoRows, got %v", err)
	}
}
