package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5"
	"github.com/sugiyamadaiki/super-camp/backend/internal/model"
	"github.com/sugiyamadaiki/super-camp/backend/internal/repository"
)

// --- mock ---

type mockChecklistRepo struct {
	ListFn       func(ctx context.Context, userID int64) ([]repository.ChecklistSummary, error)
	GetFn        func(ctx context.Context, id int64) (*repository.ChecklistWithItems, error)
	CreateFn     func(ctx context.Context, title string, userID int64) (*model.Checklist, error)
	DeleteFn     func(ctx context.Context, id int64) error
	AddItemFn    func(ctx context.Context, checklistID int64, name string, quantity int) (*model.ChecklistItem, error)
	UpdateItemFn func(ctx context.Context, id int64, name string, isChecked bool, quantity int) (*model.ChecklistItem, error)
	DeleteItemFn func(ctx context.Context, id int64) error
}

func (m *mockChecklistRepo) List(ctx context.Context, userID int64) ([]repository.ChecklistSummary, error) {
	return m.ListFn(ctx, userID)
}
func (m *mockChecklistRepo) Get(ctx context.Context, id int64) (*repository.ChecklistWithItems, error) {
	return m.GetFn(ctx, id)
}
func (m *mockChecklistRepo) Create(ctx context.Context, title string, userID int64) (*model.Checklist, error) {
	return m.CreateFn(ctx, title, userID)
}
func (m *mockChecklistRepo) Delete(ctx context.Context, id int64) error {
	return m.DeleteFn(ctx, id)
}
func (m *mockChecklistRepo) AddItem(ctx context.Context, checklistID int64, name string, quantity int) (*model.ChecklistItem, error) {
	return m.AddItemFn(ctx, checklistID, name, quantity)
}
func (m *mockChecklistRepo) UpdateItem(ctx context.Context, id int64, name string, isChecked bool, quantity int) (*model.ChecklistItem, error) {
	return m.UpdateItemFn(ctx, id, name, isChecked, quantity)
}
func (m *mockChecklistRepo) DeleteItem(ctx context.Context, id int64) error {
	return m.DeleteItemFn(ctx, id)
}

// --- helpers ---

func withChiParam(r *http.Request, key, value string) *http.Request {
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add(key, value)
	return r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, rctx))
}

func withChiParams(r *http.Request, params map[string]string) *http.Request {
	rctx := chi.NewRouteContext()
	for k, v := range params {
		rctx.URLParams.Add(k, v)
	}
	return r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, rctx))
}

// --- List ---

func TestChecklistHandler_List(t *testing.T) {
	t.Run("empty", func(t *testing.T) {
		h := NewChecklistHandler(&mockChecklistRepo{
			ListFn: func(ctx context.Context, userID int64) ([]repository.ChecklistSummary, error) {
				return nil, nil
			},
		})
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rec := httptest.NewRecorder()
		h.List(rec, req)

		if rec.Code != http.StatusOK {
			t.Fatalf("got %d, want 200", rec.Code)
		}
		var result []repository.ChecklistSummary
		json.NewDecoder(rec.Body).Decode(&result)
		if len(result) != 0 {
			t.Fatalf("expected empty slice, got %d", len(result))
		}
	})

	t.Run("with data", func(t *testing.T) {
		now := time.Now()
		h := NewChecklistHandler(&mockChecklistRepo{
			ListFn: func(ctx context.Context, userID int64) ([]repository.ChecklistSummary, error) {
				return []repository.ChecklistSummary{
					{Checklist: model.Checklist{ID: 1, UserID: 1, Title: "A", CreatedAt: now, UpdatedAt: now}, ItemCount: 3, CheckedCount: 1},
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

func TestChecklistHandler_Get(t *testing.T) {
	now := time.Now()

	t.Run("success", func(t *testing.T) {
		h := NewChecklistHandler(&mockChecklistRepo{
			GetFn: func(ctx context.Context, id int64) (*repository.ChecklistWithItems, error) {
				return &repository.ChecklistWithItems{
					Checklist: model.Checklist{ID: id, UserID: 1, Title: "Camp", CreatedAt: now, UpdatedAt: now},
					Items:     []model.ChecklistItem{},
				}, nil
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
		h := NewChecklistHandler(&mockChecklistRepo{
			GetFn: func(ctx context.Context, id int64) (*repository.ChecklistWithItems, error) {
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
		h := NewChecklistHandler(&mockChecklistRepo{})
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

func TestChecklistHandler_Create(t *testing.T) {
	now := time.Now()

	t.Run("success", func(t *testing.T) {
		h := NewChecklistHandler(&mockChecklistRepo{
			CreateFn: func(ctx context.Context, title string, userID int64) (*model.Checklist, error) {
				return &model.Checklist{ID: 1, UserID: userID, Title: title, CreatedAt: now, UpdatedAt: now}, nil
			},
		})
		body := `{"title":"New List"}`
		req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBufferString(body))
		rec := httptest.NewRecorder()
		h.Create(rec, req)

		if rec.Code != http.StatusCreated {
			t.Fatalf("got %d, want 201", rec.Code)
		}
	})

	t.Run("empty title", func(t *testing.T) {
		h := NewChecklistHandler(&mockChecklistRepo{})
		req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBufferString(`{"title":""}`))
		rec := httptest.NewRecorder()
		h.Create(rec, req)

		if rec.Code != http.StatusBadRequest {
			t.Fatalf("got %d, want 400", rec.Code)
		}
	})

	t.Run("invalid body", func(t *testing.T) {
		h := NewChecklistHandler(&mockChecklistRepo{})
		req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBufferString("bad"))
		rec := httptest.NewRecorder()
		h.Create(rec, req)

		if rec.Code != http.StatusBadRequest {
			t.Fatalf("got %d, want 400", rec.Code)
		}
	})
}

// --- Delete ---

func TestChecklistHandler_Delete(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		h := NewChecklistHandler(&mockChecklistRepo{
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
		h := NewChecklistHandler(&mockChecklistRepo{
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

// --- AddItem ---

func TestChecklistHandler_AddItem(t *testing.T) {
	now := time.Now()

	t.Run("success", func(t *testing.T) {
		h := NewChecklistHandler(&mockChecklistRepo{
			AddItemFn: func(ctx context.Context, checklistID int64, name string, quantity int) (*model.ChecklistItem, error) {
				return &model.ChecklistItem{ID: 1, ChecklistID: checklistID, Name: name, Quantity: quantity, CreatedAt: now, UpdatedAt: now}, nil
			},
		})
		req := httptest.NewRequest(http.MethodPost, "/1/items", bytes.NewBufferString(`{"name":"Tent","quantity":1}`))
		req = withChiParam(req, "id", "1")
		rec := httptest.NewRecorder()
		h.AddItem(rec, req)

		if rec.Code != http.StatusCreated {
			t.Fatalf("got %d, want 201", rec.Code)
		}
	})

	t.Run("empty name", func(t *testing.T) {
		h := NewChecklistHandler(&mockChecklistRepo{})
		req := httptest.NewRequest(http.MethodPost, "/1/items", bytes.NewBufferString(`{"name":"","quantity":1}`))
		req = withChiParam(req, "id", "1")
		rec := httptest.NewRecorder()
		h.AddItem(rec, req)

		if rec.Code != http.StatusBadRequest {
			t.Fatalf("got %d, want 400", rec.Code)
		}
	})

	t.Run("default quantity", func(t *testing.T) {
		h := NewChecklistHandler(&mockChecklistRepo{
			AddItemFn: func(ctx context.Context, checklistID int64, name string, quantity int) (*model.ChecklistItem, error) {
				if quantity != 1 {
					t.Fatalf("expected default quantity 1, got %d", quantity)
				}
				return &model.ChecklistItem{ID: 1, ChecklistID: checklistID, Name: name, Quantity: quantity, CreatedAt: now, UpdatedAt: now}, nil
			},
		})
		req := httptest.NewRequest(http.MethodPost, "/1/items", bytes.NewBufferString(`{"name":"Chair"}`))
		req = withChiParam(req, "id", "1")
		rec := httptest.NewRecorder()
		h.AddItem(rec, req)

		if rec.Code != http.StatusCreated {
			t.Fatalf("got %d, want 201", rec.Code)
		}
	})
}

// --- UpdateItem ---

func TestChecklistHandler_UpdateItem(t *testing.T) {
	now := time.Now()

	t.Run("success", func(t *testing.T) {
		h := NewChecklistHandler(&mockChecklistRepo{
			UpdateItemFn: func(ctx context.Context, id int64, name string, isChecked bool, quantity int) (*model.ChecklistItem, error) {
				return &model.ChecklistItem{ID: id, ChecklistID: 1, Name: name, IsChecked: isChecked, Quantity: quantity, CreatedAt: now, UpdatedAt: now}, nil
			},
		})
		req := httptest.NewRequest(http.MethodPut, "/1/items/1", bytes.NewBufferString(`{"name":"Tent","is_checked":true,"quantity":2}`))
		req = withChiParams(req, map[string]string{"id": "1", "itemID": "1"})
		rec := httptest.NewRecorder()
		h.UpdateItem(rec, req)

		if rec.Code != http.StatusOK {
			t.Fatalf("got %d, want 200", rec.Code)
		}
	})

	t.Run("not found", func(t *testing.T) {
		h := NewChecklistHandler(&mockChecklistRepo{
			UpdateItemFn: func(ctx context.Context, id int64, name string, isChecked bool, quantity int) (*model.ChecklistItem, error) {
				return nil, pgx.ErrNoRows
			},
		})
		req := httptest.NewRequest(http.MethodPut, "/1/items/999", bytes.NewBufferString(`{"name":"X","quantity":1}`))
		req = withChiParams(req, map[string]string{"id": "1", "itemID": "999"})
		rec := httptest.NewRecorder()
		h.UpdateItem(rec, req)

		if rec.Code != http.StatusNotFound {
			t.Fatalf("got %d, want 404", rec.Code)
		}
	})
}

// --- DeleteItem ---

func TestChecklistHandler_DeleteItem(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		h := NewChecklistHandler(&mockChecklistRepo{
			DeleteItemFn: func(ctx context.Context, id int64) error { return nil },
		})
		req := httptest.NewRequest(http.MethodDelete, "/1/items/1", nil)
		req = withChiParams(req, map[string]string{"id": "1", "itemID": "1"})
		rec := httptest.NewRecorder()
		h.DeleteItem(rec, req)

		if rec.Code != http.StatusNoContent {
			t.Fatalf("got %d, want 204", rec.Code)
		}
	})

	t.Run("not found", func(t *testing.T) {
		h := NewChecklistHandler(&mockChecklistRepo{
			DeleteItemFn: func(ctx context.Context, id int64) error { return pgx.ErrNoRows },
		})
		req := httptest.NewRequest(http.MethodDelete, "/1/items/999", nil)
		req = withChiParams(req, map[string]string{"id": "1", "itemID": "999"})
		rec := httptest.NewRecorder()
		h.DeleteItem(rec, req)

		if rec.Code != http.StatusNotFound {
			t.Fatalf("got %d, want 404", rec.Code)
		}
	})
}
