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

type mockFireLogRepo struct {
	ListFn   func(ctx context.Context, userID int64) ([]model.FireLog, error)
	GetFn    func(ctx context.Context, id int64) (*model.FireLog, error)
	CreateFn func(ctx context.Context, userID int64, date string, location string, woodType string, durationMinutes int, notes string, temperature *float64) (*model.FireLog, error)
	UpdateFn func(ctx context.Context, id int64, date string, location string, woodType string, durationMinutes int, notes string, temperature *float64) (*model.FireLog, error)
	DeleteFn func(ctx context.Context, id int64) error
}

func (m *mockFireLogRepo) List(ctx context.Context, userID int64) ([]model.FireLog, error) {
	return m.ListFn(ctx, userID)
}
func (m *mockFireLogRepo) Get(ctx context.Context, id int64) (*model.FireLog, error) {
	return m.GetFn(ctx, id)
}
func (m *mockFireLogRepo) Create(ctx context.Context, userID int64, date string, location string, woodType string, durationMinutes int, notes string, temperature *float64) (*model.FireLog, error) {
	return m.CreateFn(ctx, userID, date, location, woodType, durationMinutes, notes, temperature)
}
func (m *mockFireLogRepo) Update(ctx context.Context, id int64, date string, location string, woodType string, durationMinutes int, notes string, temperature *float64) (*model.FireLog, error) {
	return m.UpdateFn(ctx, id, date, location, woodType, durationMinutes, notes, temperature)
}
func (m *mockFireLogRepo) Delete(ctx context.Context, id int64) error {
	return m.DeleteFn(ctx, id)
}

// --- List ---

func TestFireLogHandler_List(t *testing.T) {
	t.Run("empty", func(t *testing.T) {
		h := NewFireLogHandler(&mockFireLogRepo{
			ListFn: func(ctx context.Context, userID int64) ([]model.FireLog, error) {
				return nil, nil
			},
		})
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rec := httptest.NewRecorder()
		h.List(rec, req)

		if rec.Code != http.StatusOK {
			t.Fatalf("got %d, want 200", rec.Code)
		}
		var result []model.FireLog
		json.NewDecoder(rec.Body).Decode(&result)
		if len(result) != 0 {
			t.Fatalf("expected empty, got %d", len(result))
		}
	})

	t.Run("with data", func(t *testing.T) {
		now := time.Now()
		h := NewFireLogHandler(&mockFireLogRepo{
			ListFn: func(ctx context.Context, userID int64) ([]model.FireLog, error) {
				return []model.FireLog{
					{ID: 1, UserID: 1, Date: now, Location: "River", WoodType: "Oak", DurationMinutes: 120, CreatedAt: now, UpdatedAt: now},
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

func TestFireLogHandler_Get(t *testing.T) {
	now := time.Now()

	t.Run("success", func(t *testing.T) {
		h := NewFireLogHandler(&mockFireLogRepo{
			GetFn: func(ctx context.Context, id int64) (*model.FireLog, error) {
				return &model.FireLog{ID: id, UserID: 1, Date: now, Location: "Lake", CreatedAt: now, UpdatedAt: now}, nil
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
		h := NewFireLogHandler(&mockFireLogRepo{
			GetFn: func(ctx context.Context, id int64) (*model.FireLog, error) {
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
		h := NewFireLogHandler(&mockFireLogRepo{})
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

func TestFireLogHandler_Create(t *testing.T) {
	now := time.Now()

	t.Run("success", func(t *testing.T) {
		h := NewFireLogHandler(&mockFireLogRepo{
			CreateFn: func(ctx context.Context, userID int64, date string, location string, woodType string, durationMinutes int, notes string, temperature *float64) (*model.FireLog, error) {
				return &model.FireLog{ID: 1, UserID: userID, Date: now, Location: location, WoodType: woodType, DurationMinutes: durationMinutes, CreatedAt: now, UpdatedAt: now}, nil
			},
		})
		body := `{"date":"2025-01-01","location":"River","wood_type":"Oak","duration_minutes":60,"notes":""}`
		req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBufferString(body))
		rec := httptest.NewRecorder()
		h.Create(rec, req)

		if rec.Code != http.StatusCreated {
			t.Fatalf("got %d, want 201", rec.Code)
		}
	})

	t.Run("empty date", func(t *testing.T) {
		h := NewFireLogHandler(&mockFireLogRepo{})
		req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBufferString(`{"date":"","location":"River"}`))
		rec := httptest.NewRecorder()
		h.Create(rec, req)

		if rec.Code != http.StatusBadRequest {
			t.Fatalf("got %d, want 400", rec.Code)
		}
	})

	t.Run("invalid body", func(t *testing.T) {
		h := NewFireLogHandler(&mockFireLogRepo{})
		req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBufferString("bad"))
		rec := httptest.NewRecorder()
		h.Create(rec, req)

		if rec.Code != http.StatusBadRequest {
			t.Fatalf("got %d, want 400", rec.Code)
		}
	})
}

// --- Update ---

func TestFireLogHandler_Update(t *testing.T) {
	now := time.Now()

	t.Run("success", func(t *testing.T) {
		h := NewFireLogHandler(&mockFireLogRepo{
			UpdateFn: func(ctx context.Context, id int64, date string, location string, woodType string, durationMinutes int, notes string, temperature *float64) (*model.FireLog, error) {
				return &model.FireLog{ID: id, UserID: 1, Date: now, Location: location, CreatedAt: now, UpdatedAt: now}, nil
			},
		})
		body := `{"date":"2025-02-01","location":"Mountain","wood_type":"Pine","duration_minutes":90,"notes":"nice"}`
		req := httptest.NewRequest(http.MethodPut, "/1", bytes.NewBufferString(body))
		req = withChiParam(req, "id", "1")
		rec := httptest.NewRecorder()
		h.Update(rec, req)

		if rec.Code != http.StatusOK {
			t.Fatalf("got %d, want 200", rec.Code)
		}
	})

	t.Run("not found", func(t *testing.T) {
		h := NewFireLogHandler(&mockFireLogRepo{
			UpdateFn: func(ctx context.Context, id int64, date string, location string, woodType string, durationMinutes int, notes string, temperature *float64) (*model.FireLog, error) {
				return nil, pgx.ErrNoRows
			},
		})
		req := httptest.NewRequest(http.MethodPut, "/999", bytes.NewBufferString(`{"date":"2025-01-01"}`))
		req = withChiParam(req, "id", "999")
		rec := httptest.NewRecorder()
		h.Update(rec, req)

		if rec.Code != http.StatusNotFound {
			t.Fatalf("got %d, want 404", rec.Code)
		}
	})
}

// --- Delete ---

func TestFireLogHandler_Delete(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		h := NewFireLogHandler(&mockFireLogRepo{
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
		h := NewFireLogHandler(&mockFireLogRepo{
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
