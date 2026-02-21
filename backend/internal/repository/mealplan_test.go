//go:build integration

package repository

import (
	"context"
	"testing"

	"github.com/jackc/pgx/v5"
	"github.com/sugiyamadaiki/super-camp/backend/internal/testutil"
)

func TestMealPlanRepository_CRUD(t *testing.T) {
	pool := testutil.SetupTestDB(t)
	repo := NewMealPlanRepository(pool)
	ctx := context.Background()
	var userID int64 = 1

	// Create
	plan, err := repo.Create(ctx, userID, "BBQナイト", "dinner", 4, "炭を忘れずに")
	if err != nil {
		t.Fatalf("Create: %v", err)
	}
	if plan.Title != "BBQナイト" {
		t.Fatalf("expected title 'BBQナイト', got %q", plan.Title)
	}
	if plan.MealType != "dinner" {
		t.Fatalf("expected meal_type 'dinner', got %q", plan.MealType)
	}
	if plan.Servings != 4 {
		t.Fatalf("expected servings 4, got %d", plan.Servings)
	}

	// Get
	got, err := repo.Get(ctx, plan.ID)
	if err != nil {
		t.Fatalf("Get: %v", err)
	}
	if got.Notes != "炭を忘れずに" {
		t.Fatalf("Get notes mismatch: %q", got.Notes)
	}

	// Update
	updated, err := repo.Update(ctx, plan.ID, "BBQ&焼きそば", "lunch", 6, "多めに準備")
	if err != nil {
		t.Fatalf("Update: %v", err)
	}
	if updated.Title != "BBQ&焼きそば" {
		t.Fatalf("expected updated title, got %q", updated.Title)
	}
	if updated.MealType != "lunch" {
		t.Fatalf("expected 'lunch', got %q", updated.MealType)
	}
	if updated.Servings != 6 {
		t.Fatalf("expected 6, got %d", updated.Servings)
	}

	// Create another
	plan2, err := repo.Create(ctx, userID, "朝食", "breakfast", 2, "")
	if err != nil {
		t.Fatalf("Create 2: %v", err)
	}

	// List
	list, err := repo.List(ctx, userID)
	if err != nil {
		t.Fatalf("List: %v", err)
	}
	if len(list) != 2 {
		t.Fatalf("expected 2 plans, got %d", len(list))
	}

	// Delete
	if err := repo.Delete(ctx, plan.ID); err != nil {
		t.Fatalf("Delete: %v", err)
	}
	if err := repo.Delete(ctx, plan2.ID); err != nil {
		t.Fatalf("Delete 2: %v", err)
	}

	_, err = repo.Get(ctx, plan.ID)
	if err != pgx.ErrNoRows {
		t.Fatalf("expected ErrNoRows, got %v", err)
	}
}

func TestMealPlanRepository_NotFound(t *testing.T) {
	pool := testutil.SetupTestDB(t)
	repo := NewMealPlanRepository(pool)
	ctx := context.Background()

	_, err := repo.Get(ctx, 99999)
	if err != pgx.ErrNoRows {
		t.Fatalf("expected ErrNoRows, got %v", err)
	}

	_, err = repo.Update(ctx, 99999, "x", "dinner", 2, "")
	if err != pgx.ErrNoRows {
		t.Fatalf("expected ErrNoRows, got %v", err)
	}

	err = repo.Delete(ctx, 99999)
	if err != pgx.ErrNoRows {
		t.Fatalf("expected ErrNoRows, got %v", err)
	}
}
