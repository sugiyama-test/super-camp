//go:build integration

package repository

import (
	"context"
	"testing"

	"github.com/jackc/pgx/v5"
	"github.com/sugiyamadaiki/super-camp/backend/internal/testutil"
)

func TestLayoutRepository_CRUD(t *testing.T) {
	pool := testutil.SetupTestDB(t)
	repo := NewLayoutRepository(pool)
	ctx := context.Background()
	var userID int64 = 1

	// Create
	layout, err := repo.Create(ctx, "Camp Site A", `{"items":[{"name":"tent","x":100,"y":200}]}`, userID)
	if err != nil {
		t.Fatalf("Create: %v", err)
	}
	if layout.Title != "Camp Site A" {
		t.Fatalf("expected title 'Camp Site A', got %q", layout.Title)
	}
	if layout.ID == 0 {
		t.Fatal("expected non-zero ID")
	}

	// Get
	got, err := repo.Get(ctx, layout.ID)
	if err != nil {
		t.Fatalf("Get: %v", err)
	}
	if got.Title != "Camp Site A" {
		t.Fatalf("Get title mismatch: %q", got.Title)
	}

	// Update
	updated, err := repo.Update(ctx, layout.ID, "Camp Site B", `{"items":[]}`)
	if err != nil {
		t.Fatalf("Update: %v", err)
	}
	if updated.Title != "Camp Site B" {
		t.Fatalf("expected updated title, got %q", updated.Title)
	}

	// List
	list, err := repo.List(ctx, userID)
	if err != nil {
		t.Fatalf("List: %v", err)
	}
	if len(list) != 1 {
		t.Fatalf("expected 1 layout, got %d", len(list))
	}
	if list[0].Title != "Camp Site B" {
		t.Fatalf("expected 'Camp Site B', got %q", list[0].Title)
	}

	// Delete
	if err := repo.Delete(ctx, layout.ID); err != nil {
		t.Fatalf("Delete: %v", err)
	}

	_, err = repo.Get(ctx, layout.ID)
	if err != pgx.ErrNoRows {
		t.Fatalf("expected ErrNoRows after delete, got %v", err)
	}
}

func TestLayoutRepository_NotFound(t *testing.T) {
	pool := testutil.SetupTestDB(t)
	repo := NewLayoutRepository(pool)
	ctx := context.Background()

	_, err := repo.Get(ctx, 99999)
	if err != pgx.ErrNoRows {
		t.Fatalf("expected ErrNoRows, got %v", err)
	}

	_, err = repo.Update(ctx, 99999, "x", "{}")
	if err != pgx.ErrNoRows {
		t.Fatalf("expected ErrNoRows, got %v", err)
	}

	err = repo.Delete(ctx, 99999)
	if err != pgx.ErrNoRows {
		t.Fatalf("expected ErrNoRows, got %v", err)
	}
}
