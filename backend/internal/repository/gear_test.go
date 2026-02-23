//go:build integration

package repository

import (
	"context"
	"testing"

	"github.com/jackc/pgx/v5"
	"github.com/sugiyamadaiki/super-camp/backend/internal/testutil"
)

func TestGearRepository_CRUD(t *testing.T) {
	pool := testutil.SetupTestDB(t)
	repo := NewGearRepository(pool)
	ctx := context.Background()
	var userID int64 = 1

	weight := 1500.0

	// Create
	gear, err := repo.Create(ctx, userID, "テント", "シェルター", "Snow Peak", &weight, "軽量モデル")
	if err != nil {
		t.Fatalf("Create: %v", err)
	}
	if gear.Name != "テント" {
		t.Fatalf("expected name 'テント', got %q", gear.Name)
	}
	if gear.WeightGrams == nil || *gear.WeightGrams != 1500.0 {
		t.Fatalf("expected weight 1500, got %v", gear.WeightGrams)
	}

	// Create without weight
	gear2, err := repo.Create(ctx, userID, "ナイフ", "ツール", "", nil, "")
	if err != nil {
		t.Fatalf("Create without weight: %v", err)
	}
	if gear2.WeightGrams != nil {
		t.Fatalf("expected nil weight, got %v", gear2.WeightGrams)
	}

	// Get
	got, err := repo.Get(ctx, gear.ID)
	if err != nil {
		t.Fatalf("Get: %v", err)
	}
	if got.Brand != "Snow Peak" {
		t.Fatalf("Get brand mismatch: %q", got.Brand)
	}

	// Update
	newWeight := 2000.0
	updated, err := repo.Update(ctx, gear.ID, "大型テント", "シェルター", "MSR", &newWeight, "ファミリー用")
	if err != nil {
		t.Fatalf("Update: %v", err)
	}
	if updated.Name != "大型テント" {
		t.Fatalf("expected updated name, got %q", updated.Name)
	}
	if updated.Brand != "MSR" {
		t.Fatalf("expected MSR, got %q", updated.Brand)
	}

	// List
	list, err := repo.List(ctx, userID)
	if err != nil {
		t.Fatalf("List: %v", err)
	}
	if len(list) != 2 {
		t.Fatalf("expected 2 gears, got %d", len(list))
	}

	// Delete
	if err := repo.Delete(ctx, gear.ID); err != nil {
		t.Fatalf("Delete: %v", err)
	}
	if err := repo.Delete(ctx, gear2.ID); err != nil {
		t.Fatalf("Delete 2: %v", err)
	}

	_, err = repo.Get(ctx, gear.ID)
	if err != pgx.ErrNoRows {
		t.Fatalf("expected ErrNoRows, got %v", err)
	}
}

func TestGearRepository_NotFound(t *testing.T) {
	pool := testutil.SetupTestDB(t)
	repo := NewGearRepository(pool)
	ctx := context.Background()

	_, err := repo.Get(ctx, 99999)
	if err != pgx.ErrNoRows {
		t.Fatalf("expected ErrNoRows, got %v", err)
	}

	_, err = repo.Update(ctx, 99999, "x", "", "", nil, "")
	if err != pgx.ErrNoRows {
		t.Fatalf("expected ErrNoRows, got %v", err)
	}

	err = repo.Delete(ctx, 99999)
	if err != pgx.ErrNoRows {
		t.Fatalf("expected ErrNoRows, got %v", err)
	}
}
