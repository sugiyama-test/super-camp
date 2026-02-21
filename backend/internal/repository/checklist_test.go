//go:build integration

package repository

import (
	"context"
	"testing"

	"github.com/jackc/pgx/v5"
	"github.com/sugiyamadaiki/super-camp/backend/internal/testutil"
)

func TestChecklistRepository_CRUD(t *testing.T) {
	pool := testutil.SetupTestDB(t)
	repo := NewChecklistRepository(pool)
	ctx := context.Background()
	var userID int64 = 1

	// Create
	cl, err := repo.Create(ctx, "キャンプ持ち物", userID)
	if err != nil {
		t.Fatalf("Create: %v", err)
	}
	if cl.Title != "キャンプ持ち物" {
		t.Fatalf("expected title 'キャンプ持ち物', got %q", cl.Title)
	}
	if cl.ID == 0 {
		t.Fatal("expected non-zero ID")
	}

	// Get
	got, err := repo.Get(ctx, cl.ID)
	if err != nil {
		t.Fatalf("Get: %v", err)
	}
	if got.Title != "キャンプ持ち物" {
		t.Fatalf("Get title mismatch: %q", got.Title)
	}
	if len(got.Items) != 0 {
		t.Fatalf("expected 0 items, got %d", len(got.Items))
	}

	// AddItem
	item, err := repo.AddItem(ctx, cl.ID, "テント", 1)
	if err != nil {
		t.Fatalf("AddItem: %v", err)
	}
	if item.Name != "テント" {
		t.Fatalf("expected name 'テント', got %q", item.Name)
	}
	if item.ChecklistID != cl.ID {
		t.Fatalf("expected checklist_id %d, got %d", cl.ID, item.ChecklistID)
	}

	// Add second item
	item2, err := repo.AddItem(ctx, cl.ID, "寝袋", 2)
	if err != nil {
		t.Fatalf("AddItem 2: %v", err)
	}

	// Get with items
	got, err = repo.Get(ctx, cl.ID)
	if err != nil {
		t.Fatalf("Get with items: %v", err)
	}
	if len(got.Items) != 2 {
		t.Fatalf("expected 2 items, got %d", len(got.Items))
	}

	// UpdateItem
	updated, err := repo.UpdateItem(ctx, item.ID, "テント(2人用)", true, 1)
	if err != nil {
		t.Fatalf("UpdateItem: %v", err)
	}
	if updated.Name != "テント(2人用)" {
		t.Fatalf("expected updated name, got %q", updated.Name)
	}
	if !updated.IsChecked {
		t.Fatal("expected is_checked=true")
	}

	// List
	list, err := repo.List(ctx, userID)
	if err != nil {
		t.Fatalf("List: %v", err)
	}
	if len(list) != 1 {
		t.Fatalf("expected 1 checklist, got %d", len(list))
	}
	if list[0].ItemCount != 2 {
		t.Fatalf("expected item_count=2, got %d", list[0].ItemCount)
	}
	if list[0].CheckedCount != 1 {
		t.Fatalf("expected checked_count=1, got %d", list[0].CheckedCount)
	}

	// DeleteItem
	if err := repo.DeleteItem(ctx, item2.ID); err != nil {
		t.Fatalf("DeleteItem: %v", err)
	}

	// Delete
	if err := repo.Delete(ctx, cl.ID); err != nil {
		t.Fatalf("Delete: %v", err)
	}

	// Get after delete => not found
	_, err = repo.Get(ctx, cl.ID)
	if err != pgx.ErrNoRows {
		t.Fatalf("expected ErrNoRows after delete, got %v", err)
	}
}

func TestChecklistRepository_NotFound(t *testing.T) {
	pool := testutil.SetupTestDB(t)
	repo := NewChecklistRepository(pool)
	ctx := context.Background()

	_, err := repo.Get(ctx, 99999)
	if err != pgx.ErrNoRows {
		t.Fatalf("expected ErrNoRows, got %v", err)
	}

	err = repo.Delete(ctx, 99999)
	if err != pgx.ErrNoRows {
		t.Fatalf("expected ErrNoRows, got %v", err)
	}

	_, err = repo.UpdateItem(ctx, 99999, "x", false, 1)
	if err != pgx.ErrNoRows {
		t.Fatalf("expected ErrNoRows, got %v", err)
	}

	err = repo.DeleteItem(ctx, 99999)
	if err != pgx.ErrNoRows {
		t.Fatalf("expected ErrNoRows, got %v", err)
	}
}
