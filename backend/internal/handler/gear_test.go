package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/sugiyamadaiki/super-camp/backend/internal/model"
)

// --- mock ---

type mockGearRepo struct {
	ListFn   func(ctx context.Context, userID int64) ([]model.Gear, error)
	GetFn    func(ctx context.Context, id int64) (*model.Gear, error)
	CreateFn func(ctx context.Context, userID int64, name string, category string, brand string, weightGrams *float64, notes string) (*model.Gear, error)
	UpdateFn func(ctx context.Context, id int64, name string, category string, brand string, weightGrams *float64, notes string) (*model.Gear, error)
	DeleteFn func(ctx context.Context, id int64) error
}

func (m *mockGearRepo) List(ctx context.Context, userID int64) ([]model.Gear, error) {
	return m.ListFn(ctx, userID)
}
func (m *mockGearRepo) Get(ctx context.Context, id int64) (*model.Gear, error) {
	return m.GetFn(ctx, id)
}
func (m *mockGearRepo) Create(ctx context.Context, userID int64, name string, category string, brand string, weightGrams *float64, notes string) (*model.Gear, error) {
	return m.CreateFn(ctx, userID, name, category, brand, weightGrams, notes)
}
func (m *mockGearRepo) Update(ctx context.Context, id int64, name string, category string, brand string, weightGrams *float64, notes string) (*model.Gear, error) {
	return m.UpdateFn(ctx, id, name, category, brand, weightGrams, notes)
}
func (m *mockGearRepo) Delete(ctx context.Context, id int64) error {
	return m.DeleteFn(ctx, id)
}

// --- List ---

func TestGearHandler_List(t *testing.T) {
	t.Run("empty", func(t *testing.T) {
		h := NewGearHandler(&mockGearRepo{
			ListFn: func(ctx context.Context, userID int64) ([]model.Gear, error) {
				return nil, nil
			},
		})
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rec := httptest.NewRecorder()
		h.List(rec, req)

		if rec.Code != http.StatusOK {
			t.Fatalf("got %d, want 200", rec.Code)
		}
		var result []model.Gear
		json.NewDecoder(rec.Body).Decode(&result)
		if len(result) != 0 {
			t.Fatalf("expected empty, got %d", len(result))
		}
	})

	t.Run("with data", func(t *testing.T) {
		now := time.Now()
		w := 1500.0
		h := NewGearHandler(&mockGearRepo{
			ListFn: func(ctx context.Context, userID int64) ([]model.Gear, error) {
				return []model.Gear{
					{ID: 1, UserID: 1, Name: "テント", Category: "シェルター", Brand: "Snow Peak", WeightGrams: &w, CreatedAt: now, UpdatedAt: now},
				}, nil
			},
		})
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rec := httptest.NewRecorder()
		h.List(rec, req)

		if rec.Code != http.StatusOK {
			t.Fatalf("got %d, want 200", rec.Code)
		}
	})
}

// --- Get ---

func TestGearHandler_Get(t *testing.T) {
	now := time.Now()

	t.Run("success", func(t *testing.T) {
		h := NewGearHandler(&mockGearRepo{
			GetFn: func(ctx context.Context, id int64) (*model.Gear, error) {
				return &model.Gear{ID: id, UserID: 1, Name: "テント", CreatedAt: now, UpdatedAt: now}, nil
			},
		})
		req := httptest.NewRequest(http.MethodGet, "/1", nil)
		req = withChiParam(req, "id", "1")
		rec := httptest.NewRecorder()
		h.Get(rec, req)

		if rec.Code != http.StatusOK {
			t.Fatalf("got %d, want 200", rec.Code)
		}
	})

	t.Run("not found", func(t *testing.T) {
		h := NewGearHandler(&mockGearRepo{
			GetFn: func(ctx context.Context, id int64) (*model.Gear, error) {
				return nil, pgx.ErrNoRows
			},
		})
		req := httptest.NewRequest(http.MethodGet, "/999", nil)
		req = withChiParam(req, "id", "999")
		rec := httptest.NewRecorder()
		h.Get(rec, req)

		if rec.Code != http.StatusNotFound {
			t.Fatalf("got %d, want 404", rec.Code)
		}
	})

	t.Run("invalid id", func(t *testing.T) {
		h := NewGearHandler(&mockGearRepo{})
		req := httptest.NewRequest(http.MethodGet, "/abc", nil)
		req = withChiParam(req, "id", "abc")
		rec := httptest.NewRecorder()
		h.Get(rec, req)

		if rec.Code != http.StatusBadRequest {
			t.Fatalf("got %d, want 400", rec.Code)
		}
	})
}

// --- Create ---

func TestGearHandler_Create(t *testing.T) {
	now := time.Now()

	t.Run("success", func(t *testing.T) {
		h := NewGearHandler(&mockGearRepo{
			CreateFn: func(ctx context.Context, userID int64, name string, category string, brand string, weightGrams *float64, notes string) (*model.Gear, error) {
				return &model.Gear{ID: 1, UserID: userID, Name: name, Category: category, Brand: brand, WeightGrams: weightGrams, Notes: notes, CreatedAt: now, UpdatedAt: now}, nil
			},
		})
		body := `{"name":"テント","category":"シェルター","brand":"Snow Peak","weight_grams":1500,"notes":"軽量"}`
		req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBufferString(body))
		rec := httptest.NewRecorder()
		h.Create(rec, req)

		if rec.Code != http.StatusCreated {
			t.Fatalf("got %d, want 201", rec.Code)
		}
	})

	t.Run("empty name", func(t *testing.T) {
		h := NewGearHandler(&mockGearRepo{})
		req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBufferString(`{"name":"","category":"test"}`))
		rec := httptest.NewRecorder()
		h.Create(rec, req)

		if rec.Code != http.StatusBadRequest {
			t.Fatalf("got %d, want 400", rec.Code)
		}
	})

	t.Run("invalid body", func(t *testing.T) {
		h := NewGearHandler(&mockGearRepo{})
		req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBufferString("bad"))
		rec := httptest.NewRecorder()
		h.Create(rec, req)

		if rec.Code != http.StatusBadRequest {
			t.Fatalf("got %d, want 400", rec.Code)
		}
	})
}

// --- Update ---

func TestGearHandler_Update(t *testing.T) {
	now := time.Now()

	t.Run("success", func(t *testing.T) {
		h := NewGearHandler(&mockGearRepo{
			UpdateFn: func(ctx context.Context, id int64, name string, category string, brand string, weightGrams *float64, notes string) (*model.Gear, error) {
				return &model.Gear{ID: id, UserID: 1, Name: name, Category: category, Brand: brand, CreatedAt: now, UpdatedAt: now}, nil
			},
		})
		body := `{"name":"タープ","category":"シェルター","brand":"MSR","notes":"大型"}`
		req := httptest.NewRequest(http.MethodPut, "/1", bytes.NewBufferString(body))
		req = withChiParam(req, "id", "1")
		rec := httptest.NewRecorder()
		h.Update(rec, req)

		if rec.Code != http.StatusOK {
			t.Fatalf("got %d, want 200", rec.Code)
		}
	})

	t.Run("not found", func(t *testing.T) {
		h := NewGearHandler(&mockGearRepo{
			UpdateFn: func(ctx context.Context, id int64, name string, category string, brand string, weightGrams *float64, notes string) (*model.Gear, error) {
				return nil, pgx.ErrNoRows
			},
		})
		req := httptest.NewRequest(http.MethodPut, "/999", bytes.NewBufferString(`{"name":"x"}`))
		req = withChiParam(req, "id", "999")
		rec := httptest.NewRecorder()
		h.Update(rec, req)

		if rec.Code != http.StatusNotFound {
			t.Fatalf("got %d, want 404", rec.Code)
		}
	})
}

// --- Delete ---

func TestGearHandler_Delete(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		h := NewGearHandler(&mockGearRepo{
			DeleteFn: func(ctx context.Context, id int64) error { return nil },
		})
		req := httptest.NewRequest(http.MethodDelete, "/1", nil)
		req = withChiParam(req, "id", "1")
		rec := httptest.NewRecorder()
		h.Delete(rec, req)

		if rec.Code != http.StatusNoContent {
			t.Fatalf("got %d, want 204", rec.Code)
		}
	})

	t.Run("not found", func(t *testing.T) {
		h := NewGearHandler(&mockGearRepo{
			DeleteFn: func(ctx context.Context, id int64) error { return pgx.ErrNoRows },
		})
		req := httptest.NewRequest(http.MethodDelete, "/999", nil)
		req = withChiParam(req, "id", "999")
		rec := httptest.NewRecorder()
		h.Delete(rec, req)

		if rec.Code != http.StatusNotFound {
			t.Fatalf("got %d, want 404", rec.Code)
		}
	})
}
