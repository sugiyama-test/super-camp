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
	"github.com/sugiyamadaiki/super-camp/backend/internal/repository"
)

// --- mock ---

type mockLayoutRepo struct {
	ListFn   func(ctx context.Context, userID int64) ([]repository.LayoutSummary, error)
	GetFn    func(ctx context.Context, id int64) (*model.Layout, error)
	CreateFn func(ctx context.Context, title string, data string, userID int64) (*model.Layout, error)
	UpdateFn func(ctx context.Context, id int64, title string, data string) (*model.Layout, error)
	DeleteFn func(ctx context.Context, id int64) error
}

func (m *mockLayoutRepo) List(ctx context.Context, userID int64) ([]repository.LayoutSummary, error) {
	return m.ListFn(ctx, userID)
}
func (m *mockLayoutRepo) Get(ctx context.Context, id int64) (*model.Layout, error) {
	return m.GetFn(ctx, id)
}
func (m *mockLayoutRepo) Create(ctx context.Context, title string, data string, userID int64) (*model.Layout, error) {
	return m.CreateFn(ctx, title, data, userID)
}
func (m *mockLayoutRepo) Update(ctx context.Context, id int64, title string, data string) (*model.Layout, error) {
	return m.UpdateFn(ctx, id, title, data)
}
func (m *mockLayoutRepo) Delete(ctx context.Context, id int64) error {
	return m.DeleteFn(ctx, id)
}

// --- List ---

func TestLayoutHandler_List(t *testing.T) {
	t.Run("empty", func(t *testing.T) {
		h := NewLayoutHandler(&mockLayoutRepo{
			ListFn: func(ctx context.Context, userID int64) ([]repository.LayoutSummary, error) {
				return nil, nil
			},
		})
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rec := httptest.NewRecorder()
		h.List(rec, req)

		if rec.Code != http.StatusOK {
			t.Fatalf("got %d, want 200", rec.Code)
		}
		var result []repository.LayoutSummary
		json.NewDecoder(rec.Body).Decode(&result)
		if len(result) != 0 {
			t.Fatalf("expected empty, got %d", len(result))
		}
	})

	t.Run("with data", func(t *testing.T) {
		h := NewLayoutHandler(&mockLayoutRepo{
			ListFn: func(ctx context.Context, userID int64) ([]repository.LayoutSummary, error) {
				return []repository.LayoutSummary{
					{ID: 1, UserID: 1, Title: "Camp A", ItemCount: 3},
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

func TestLayoutHandler_Get(t *testing.T) {
	now := time.Now()

	t.Run("success", func(t *testing.T) {
		h := NewLayoutHandler(&mockLayoutRepo{
			GetFn: func(ctx context.Context, id int64) (*model.Layout, error) {
				return &model.Layout{ID: id, UserID: 1, Title: "My Layout", Data: `{}`, CreatedAt: now, UpdatedAt: now}, nil
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
		h := NewLayoutHandler(&mockLayoutRepo{
			GetFn: func(ctx context.Context, id int64) (*model.Layout, error) {
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
		h := NewLayoutHandler(&mockLayoutRepo{})
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

func TestLayoutHandler_Create(t *testing.T) {
	now := time.Now()

	t.Run("success", func(t *testing.T) {
		h := NewLayoutHandler(&mockLayoutRepo{
			CreateFn: func(ctx context.Context, title string, data string, userID int64) (*model.Layout, error) {
				return &model.Layout{ID: 1, UserID: userID, Title: title, Data: data, CreatedAt: now, UpdatedAt: now}, nil
			},
		})
		body := `{"title":"New Layout","data":{"items":[]}}`
		req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBufferString(body))
		rec := httptest.NewRecorder()
		h.Create(rec, req)

		if rec.Code != http.StatusCreated {
			t.Fatalf("got %d, want 201", rec.Code)
		}
	})

	t.Run("empty title", func(t *testing.T) {
		h := NewLayoutHandler(&mockLayoutRepo{})
		req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBufferString(`{"title":"","data":{}}`))
		rec := httptest.NewRecorder()
		h.Create(rec, req)

		if rec.Code != http.StatusBadRequest {
			t.Fatalf("got %d, want 400", rec.Code)
		}
	})

	t.Run("invalid body", func(t *testing.T) {
		h := NewLayoutHandler(&mockLayoutRepo{})
		req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBufferString("bad"))
		rec := httptest.NewRecorder()
		h.Create(rec, req)

		if rec.Code != http.StatusBadRequest {
			t.Fatalf("got %d, want 400", rec.Code)
		}
	})
}

// --- Update ---

func TestLayoutHandler_Update(t *testing.T) {
	now := time.Now()

	t.Run("success", func(t *testing.T) {
		h := NewLayoutHandler(&mockLayoutRepo{
			UpdateFn: func(ctx context.Context, id int64, title string, data string) (*model.Layout, error) {
				return &model.Layout{ID: id, UserID: 1, Title: title, Data: data, CreatedAt: now, UpdatedAt: now}, nil
			},
		})
		body := `{"title":"Updated","data":{}}`
		req := httptest.NewRequest(http.MethodPut, "/1", bytes.NewBufferString(body))
		req = withChiParam(req, "id", "1")
		rec := httptest.NewRecorder()
		h.Update(rec, req)

		if rec.Code != http.StatusOK {
			t.Fatalf("got %d, want 200", rec.Code)
		}
	})

	t.Run("not found", func(t *testing.T) {
		h := NewLayoutHandler(&mockLayoutRepo{
			UpdateFn: func(ctx context.Context, id int64, title string, data string) (*model.Layout, error) {
				return nil, pgx.ErrNoRows
			},
		})
		req := httptest.NewRequest(http.MethodPut, "/999", bytes.NewBufferString(`{"title":"X","data":{}}`))
		req = withChiParam(req, "id", "999")
		rec := httptest.NewRecorder()
		h.Update(rec, req)

		if rec.Code != http.StatusNotFound {
			t.Fatalf("got %d, want 404", rec.Code)
		}
	})
}

// --- Delete ---

func TestLayoutHandler_Delete(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		h := NewLayoutHandler(&mockLayoutRepo{
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
		h := NewLayoutHandler(&mockLayoutRepo{
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
