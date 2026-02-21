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

type mockMealPlanRepo struct {
	ListFn   func(ctx context.Context, userID int64) ([]model.MealPlan, error)
	GetFn    func(ctx context.Context, id int64) (*model.MealPlan, error)
	CreateFn func(ctx context.Context, userID int64, title string, mealType string, servings int, notes string) (*model.MealPlan, error)
	UpdateFn func(ctx context.Context, id int64, title string, mealType string, servings int, notes string) (*model.MealPlan, error)
	DeleteFn func(ctx context.Context, id int64) error
}

func (m *mockMealPlanRepo) List(ctx context.Context, userID int64) ([]model.MealPlan, error) {
	return m.ListFn(ctx, userID)
}
func (m *mockMealPlanRepo) Get(ctx context.Context, id int64) (*model.MealPlan, error) {
	return m.GetFn(ctx, id)
}
func (m *mockMealPlanRepo) Create(ctx context.Context, userID int64, title string, mealType string, servings int, notes string) (*model.MealPlan, error) {
	return m.CreateFn(ctx, userID, title, mealType, servings, notes)
}
func (m *mockMealPlanRepo) Update(ctx context.Context, id int64, title string, mealType string, servings int, notes string) (*model.MealPlan, error) {
	return m.UpdateFn(ctx, id, title, mealType, servings, notes)
}
func (m *mockMealPlanRepo) Delete(ctx context.Context, id int64) error {
	return m.DeleteFn(ctx, id)
}

// --- List ---

func TestMealPlanHandler_List(t *testing.T) {
	t.Run("empty", func(t *testing.T) {
		h := NewMealPlanHandler(&mockMealPlanRepo{
			ListFn: func(ctx context.Context, userID int64) ([]model.MealPlan, error) {
				return nil, nil
			},
		})
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rec := httptest.NewRecorder()
		h.List(rec, req)

		if rec.Code != http.StatusOK {
			t.Fatalf("got %d, want 200", rec.Code)
		}
		var result []model.MealPlan
		json.NewDecoder(rec.Body).Decode(&result)
		if len(result) != 0 {
			t.Fatalf("expected empty, got %d", len(result))
		}
	})

	t.Run("with data", func(t *testing.T) {
		now := time.Now()
		h := NewMealPlanHandler(&mockMealPlanRepo{
			ListFn: func(ctx context.Context, userID int64) ([]model.MealPlan, error) {
				return []model.MealPlan{
					{ID: 1, UserID: 1, Title: "BBQ", MealType: "dinner", Servings: 4, CreatedAt: now, UpdatedAt: now},
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

func TestMealPlanHandler_Get(t *testing.T) {
	now := time.Now()

	t.Run("success", func(t *testing.T) {
		h := NewMealPlanHandler(&mockMealPlanRepo{
			GetFn: func(ctx context.Context, id int64) (*model.MealPlan, error) {
				return &model.MealPlan{ID: id, UserID: 1, Title: "BBQ", MealType: "dinner", Servings: 4, CreatedAt: now, UpdatedAt: now}, nil
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
		h := NewMealPlanHandler(&mockMealPlanRepo{
			GetFn: func(ctx context.Context, id int64) (*model.MealPlan, error) {
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
		h := NewMealPlanHandler(&mockMealPlanRepo{})
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

func TestMealPlanHandler_Create(t *testing.T) {
	now := time.Now()

	t.Run("success", func(t *testing.T) {
		h := NewMealPlanHandler(&mockMealPlanRepo{
			CreateFn: func(ctx context.Context, userID int64, title string, mealType string, servings int, notes string) (*model.MealPlan, error) {
				return &model.MealPlan{ID: 1, UserID: userID, Title: title, MealType: mealType, Servings: servings, Notes: notes, CreatedAt: now, UpdatedAt: now}, nil
			},
		})
		body := `{"title":"BBQ Night","meal_type":"dinner","servings":4,"notes":"bring charcoal"}`
		req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBufferString(body))
		rec := httptest.NewRecorder()
		h.Create(rec, req)

		if rec.Code != http.StatusCreated {
			t.Fatalf("got %d, want 201", rec.Code)
		}
	})

	t.Run("defaults", func(t *testing.T) {
		h := NewMealPlanHandler(&mockMealPlanRepo{
			CreateFn: func(ctx context.Context, userID int64, title string, mealType string, servings int, notes string) (*model.MealPlan, error) {
				if mealType != "dinner" {
					t.Fatalf("expected default meal_type 'dinner', got %q", mealType)
				}
				if servings != 2 {
					t.Fatalf("expected default servings 2, got %d", servings)
				}
				return &model.MealPlan{ID: 1, UserID: userID, Title: title, MealType: mealType, Servings: servings, CreatedAt: now, UpdatedAt: now}, nil
			},
		})
		req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBufferString(`{"title":"Simple"}`))
		rec := httptest.NewRecorder()
		h.Create(rec, req)

		if rec.Code != http.StatusCreated {
			t.Fatalf("got %d, want 201", rec.Code)
		}
	})

	t.Run("empty title", func(t *testing.T) {
		h := NewMealPlanHandler(&mockMealPlanRepo{})
		req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBufferString(`{"title":""}`))
		rec := httptest.NewRecorder()
		h.Create(rec, req)

		if rec.Code != http.StatusBadRequest {
			t.Fatalf("got %d, want 400", rec.Code)
		}
	})

	t.Run("invalid body", func(t *testing.T) {
		h := NewMealPlanHandler(&mockMealPlanRepo{})
		req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBufferString("bad"))
		rec := httptest.NewRecorder()
		h.Create(rec, req)

		if rec.Code != http.StatusBadRequest {
			t.Fatalf("got %d, want 400", rec.Code)
		}
	})
}

// --- Update ---

func TestMealPlanHandler_Update(t *testing.T) {
	now := time.Now()

	t.Run("success", func(t *testing.T) {
		h := NewMealPlanHandler(&mockMealPlanRepo{
			UpdateFn: func(ctx context.Context, id int64, title string, mealType string, servings int, notes string) (*model.MealPlan, error) {
				return &model.MealPlan{ID: id, UserID: 1, Title: title, MealType: mealType, Servings: servings, Notes: notes, CreatedAt: now, UpdatedAt: now}, nil
			},
		})
		body := `{"title":"Updated","meal_type":"lunch","servings":3,"notes":"updated"}`
		req := httptest.NewRequest(http.MethodPut, "/1", bytes.NewBufferString(body))
		req = withChiParam(req, "id", "1")
		rec := httptest.NewRecorder()
		h.Update(rec, req)

		if rec.Code != http.StatusOK {
			t.Fatalf("got %d, want 200", rec.Code)
		}
	})

	t.Run("not found", func(t *testing.T) {
		h := NewMealPlanHandler(&mockMealPlanRepo{
			UpdateFn: func(ctx context.Context, id int64, title string, mealType string, servings int, notes string) (*model.MealPlan, error) {
				return nil, pgx.ErrNoRows
			},
		})
		req := httptest.NewRequest(http.MethodPut, "/999", bytes.NewBufferString(`{"title":"X","meal_type":"dinner","servings":2}`))
		req = withChiParam(req, "id", "999")
		rec := httptest.NewRecorder()
		h.Update(rec, req)

		if rec.Code != http.StatusNotFound {
			t.Fatalf("got %d, want 404", rec.Code)
		}
	})
}

// --- Delete ---

func TestMealPlanHandler_Delete(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		h := NewMealPlanHandler(&mockMealPlanRepo{
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
		h := NewMealPlanHandler(&mockMealPlanRepo{
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
