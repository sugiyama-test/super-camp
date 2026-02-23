//go:build integration

package repository

import (
	"context"
	"testing"

	"github.com/jackc/pgx/v5"
	"github.com/sugiyamadaiki/super-camp/backend/internal/testutil"
)

func TestCampsiteRepository_CRUD(t *testing.T) {
	pool := testutil.SetupTestDB(t)
	repo := NewCampsiteRepository(pool)
	ctx := context.Background()

	lat := 35.6762
	lng := 139.6503

	// Create
	campsite, err := repo.Create(ctx, "ふもとっぱら", "静岡県富士宮市", &lat, &lng, "富士山が見える")
	if err != nil {
		t.Fatalf("Create: %v", err)
	}
	if campsite.Name != "ふもとっぱら" {
		t.Fatalf("expected name 'ふもとっぱら', got %q", campsite.Name)
	}
	if campsite.Latitude == nil || *campsite.Latitude != 35.6762 {
		t.Fatalf("expected latitude 35.6762, got %v", campsite.Latitude)
	}

	// Create without coordinates
	campsite2, err := repo.Create(ctx, "道志の森", "山梨県道志村", nil, nil, "")
	if err != nil {
		t.Fatalf("Create without coords: %v", err)
	}
	if campsite2.Latitude != nil {
		t.Fatalf("expected nil latitude, got %v", campsite2.Latitude)
	}

	// Get
	got, err := repo.Get(ctx, campsite.ID)
	if err != nil {
		t.Fatalf("Get: %v", err)
	}
	if got.Address != "静岡県富士宮市" {
		t.Fatalf("Get address mismatch: %q", got.Address)
	}

	// Update
	newLat := 35.45
	newLng := 138.77
	updated, err := repo.Update(ctx, campsite.ID, "ふもとっぱらキャンプ場", "静岡県富士宮市麓156", &newLat, &newLng, "広大なフリーサイト")
	if err != nil {
		t.Fatalf("Update: %v", err)
	}
	if updated.Name != "ふもとっぱらキャンプ場" {
		t.Fatalf("expected updated name, got %q", updated.Name)
	}
	if updated.Address != "静岡県富士宮市麓156" {
		t.Fatalf("expected updated address, got %q", updated.Address)
	}

	// List (no userID filter)
	list, err := repo.List(ctx)
	if err != nil {
		t.Fatalf("List: %v", err)
	}
	if len(list) != 2 {
		t.Fatalf("expected 2 campsites, got %d", len(list))
	}

	// Delete
	if err := repo.Delete(ctx, campsite.ID); err != nil {
		t.Fatalf("Delete: %v", err)
	}
	if err := repo.Delete(ctx, campsite2.ID); err != nil {
		t.Fatalf("Delete 2: %v", err)
	}

	_, err = repo.Get(ctx, campsite.ID)
	if err != pgx.ErrNoRows {
		t.Fatalf("expected ErrNoRows, got %v", err)
	}
}

func TestCampsiteRepository_NotFound(t *testing.T) {
	pool := testutil.SetupTestDB(t)
	repo := NewCampsiteRepository(pool)
	ctx := context.Background()

	_, err := repo.Get(ctx, 99999)
	if err != pgx.ErrNoRows {
		t.Fatalf("expected ErrNoRows, got %v", err)
	}

	_, err = repo.Update(ctx, 99999, "x", "", nil, nil, "")
	if err != pgx.ErrNoRows {
		t.Fatalf("expected ErrNoRows, got %v", err)
	}

	err = repo.Delete(ctx, 99999)
	if err != pgx.ErrNoRows {
		t.Fatalf("expected ErrNoRows, got %v", err)
	}
}
