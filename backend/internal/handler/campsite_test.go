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

type mockCampsiteRepo struct {
	ListFn   func(ctx context.Context) ([]model.Campsite, error)
	GetFn    func(ctx context.Context, id int64) (*model.Campsite, error)
	CreateFn func(ctx context.Context, name string, address string, latitude *float64, longitude *float64, notes string) (*model.Campsite, error)
	UpdateFn func(ctx context.Context, id int64, name string, address string, latitude *float64, longitude *float64, notes string) (*model.Campsite, error)
	DeleteFn func(ctx context.Context, id int64) error
}

func (m *mockCampsiteRepo) List(ctx context.Context) ([]model.Campsite, error) {
	return m.ListFn(ctx)
}
func (m *mockCampsiteRepo) Get(ctx context.Context, id int64) (*model.Campsite, error) {
	return m.GetFn(ctx, id)
}
func (m *mockCampsiteRepo) Create(ctx context.Context, name string, address string, latitude *float64, longitude *float64, notes string) (*model.Campsite, error) {
	return m.CreateFn(ctx, name, address, latitude, longitude, notes)
}
func (m *mockCampsiteRepo) Update(ctx context.Context, id int64, name string, address string, latitude *float64, longitude *float64, notes string) (*model.Campsite, error) {
	return m.UpdateFn(ctx, id, name, address, latitude, longitude, notes)
}
func (m *mockCampsiteRepo) Delete(ctx context.Context, id int64) error {
	return m.DeleteFn(ctx, id)
}

// --- List ---

func TestCampsiteHandler_List(t *testing.T) {
	t.Run("empty", func(t *testing.T) {
		h := NewCampsiteHandler(&mockCampsiteRepo{
			ListFn: func(ctx context.Context) ([]model.Campsite, error) {
				return nil, nil
			},
		})
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rec := httptest.NewRecorder()
		h.List(rec, req)

		if rec.Code != http.StatusOK {
			t.Fatalf("got %d, want 200", rec.Code)
		}
		var result []model.Campsite
		json.NewDecoder(rec.Body).Decode(&result)
		if len(result) != 0 {
			t.Fatalf("expected empty, got %d", len(result))
		}
	})

	t.Run("with data", func(t *testing.T) {
		now := time.Now()
		lat := 35.6762
		lng := 139.6503
		h := NewCampsiteHandler(&mockCampsiteRepo{
			ListFn: func(ctx context.Context) ([]model.Campsite, error) {
				return []model.Campsite{
					{ID: 1, Name: "ふもとっぱら", Address: "静岡県富士宮市", Latitude: &lat, Longitude: &lng, CreatedAt: now, UpdatedAt: now},
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

func TestCampsiteHandler_Get(t *testing.T) {
	now := time.Now()

	t.Run("success", func(t *testing.T) {
		h := NewCampsiteHandler(&mockCampsiteRepo{
			GetFn: func(ctx context.Context, id int64) (*model.Campsite, error) {
				return &model.Campsite{ID: id, Name: "ふもとっぱら", CreatedAt: now, UpdatedAt: now}, nil
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
		h := NewCampsiteHandler(&mockCampsiteRepo{
			GetFn: func(ctx context.Context, id int64) (*model.Campsite, error) {
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
		h := NewCampsiteHandler(&mockCampsiteRepo{})
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

func TestCampsiteHandler_Create(t *testing.T) {
	now := time.Now()

	t.Run("success", func(t *testing.T) {
		h := NewCampsiteHandler(&mockCampsiteRepo{
			CreateFn: func(ctx context.Context, name string, address string, latitude *float64, longitude *float64, notes string) (*model.Campsite, error) {
				return &model.Campsite{ID: 1, Name: name, Address: address, Latitude: latitude, Longitude: longitude, Notes: notes, CreatedAt: now, UpdatedAt: now}, nil
			},
		})
		body := `{"name":"ふもとっぱら","address":"静岡県富士宮市","latitude":35.6762,"longitude":139.6503,"notes":"富士山が見える"}`
		req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBufferString(body))
		rec := httptest.NewRecorder()
		h.Create(rec, req)

		if rec.Code != http.StatusCreated {
			t.Fatalf("got %d, want 201", rec.Code)
		}
	})

	t.Run("empty name", func(t *testing.T) {
		h := NewCampsiteHandler(&mockCampsiteRepo{})
		req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBufferString(`{"name":"","address":"test"}`))
		rec := httptest.NewRecorder()
		h.Create(rec, req)

		if rec.Code != http.StatusBadRequest {
			t.Fatalf("got %d, want 400", rec.Code)
		}
	})

	t.Run("invalid body", func(t *testing.T) {
		h := NewCampsiteHandler(&mockCampsiteRepo{})
		req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBufferString("bad"))
		rec := httptest.NewRecorder()
		h.Create(rec, req)

		if rec.Code != http.StatusBadRequest {
			t.Fatalf("got %d, want 400", rec.Code)
		}
	})
}

// --- Update ---

func TestCampsiteHandler_Update(t *testing.T) {
	now := time.Now()

	t.Run("success", func(t *testing.T) {
		h := NewCampsiteHandler(&mockCampsiteRepo{
			UpdateFn: func(ctx context.Context, id int64, name string, address string, latitude *float64, longitude *float64, notes string) (*model.Campsite, error) {
				return &model.Campsite{ID: id, Name: name, Address: address, CreatedAt: now, UpdatedAt: now}, nil
			},
		})
		body := `{"name":"道志の森","address":"山梨県道志村","notes":"川沿い"}`
		req := httptest.NewRequest(http.MethodPut, "/1", bytes.NewBufferString(body))
		req = withChiParam(req, "id", "1")
		rec := httptest.NewRecorder()
		h.Update(rec, req)

		if rec.Code != http.StatusOK {
			t.Fatalf("got %d, want 200", rec.Code)
		}
	})

	t.Run("not found", func(t *testing.T) {
		h := NewCampsiteHandler(&mockCampsiteRepo{
			UpdateFn: func(ctx context.Context, id int64, name string, address string, latitude *float64, longitude *float64, notes string) (*model.Campsite, error) {
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

func TestCampsiteHandler_Delete(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		h := NewCampsiteHandler(&mockCampsiteRepo{
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
		h := NewCampsiteHandler(&mockCampsiteRepo{
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
