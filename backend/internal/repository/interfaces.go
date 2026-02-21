package repository

import (
	"context"

	"github.com/sugiyamadaiki/super-camp/backend/internal/model"
)

type ChecklistRepo interface {
	List(ctx context.Context, userID int64) ([]ChecklistSummary, error)
	Get(ctx context.Context, id int64) (*ChecklistWithItems, error)
	Create(ctx context.Context, title string, userID int64) (*model.Checklist, error)
	Delete(ctx context.Context, id int64) error
	AddItem(ctx context.Context, checklistID int64, name string, quantity int) (*model.ChecklistItem, error)
	UpdateItem(ctx context.Context, id int64, name string, isChecked bool, quantity int) (*model.ChecklistItem, error)
	DeleteItem(ctx context.Context, id int64) error
}

type LayoutRepo interface {
	List(ctx context.Context, userID int64) ([]LayoutSummary, error)
	Get(ctx context.Context, id int64) (*model.Layout, error)
	Create(ctx context.Context, title string, data string, userID int64) (*model.Layout, error)
	Update(ctx context.Context, id int64, title string, data string) (*model.Layout, error)
	Delete(ctx context.Context, id int64) error
}

type FireLogRepo interface {
	List(ctx context.Context, userID int64) ([]model.FireLog, error)
	Get(ctx context.Context, id int64) (*model.FireLog, error)
	Create(ctx context.Context, userID int64, date string, location string, woodType string, durationMinutes int, notes string, temperature *float64) (*model.FireLog, error)
	Update(ctx context.Context, id int64, date string, location string, woodType string, durationMinutes int, notes string, temperature *float64) (*model.FireLog, error)
	Delete(ctx context.Context, id int64) error
}

type MealPlanRepo interface {
	List(ctx context.Context, userID int64) ([]model.MealPlan, error)
	Get(ctx context.Context, id int64) (*model.MealPlan, error)
	Create(ctx context.Context, userID int64, title string, mealType string, servings int, notes string) (*model.MealPlan, error)
	Update(ctx context.Context, id int64, title string, mealType string, servings int, notes string) (*model.MealPlan, error)
	Delete(ctx context.Context, id int64) error
}
