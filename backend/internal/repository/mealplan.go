package repository

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sugiyamadaiki/super-camp/backend/internal/model"
)

type MealPlanRepository struct {
	pool *pgxpool.Pool
}

func NewMealPlanRepository(pool *pgxpool.Pool) *MealPlanRepository {
	return &MealPlanRepository{pool: pool}
}

func (r *MealPlanRepository) List(ctx context.Context, userID int64) ([]model.MealPlan, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT id, user_id, title, meal_type, servings, notes, created_at, updated_at
		FROM meal_plans WHERE user_id = $1
		ORDER BY created_at DESC
	`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []model.MealPlan
	for rows.Next() {
		var m model.MealPlan
		if err := rows.Scan(&m.ID, &m.UserID, &m.Title, &m.MealType, &m.Servings, &m.Notes, &m.CreatedAt, &m.UpdatedAt); err != nil {
			return nil, err
		}
		result = append(result, m)
	}
	return result, rows.Err()
}

func (r *MealPlanRepository) Get(ctx context.Context, id int64) (*model.MealPlan, error) {
	var m model.MealPlan
	err := r.pool.QueryRow(ctx, `
		SELECT id, user_id, title, meal_type, servings, notes, created_at, updated_at
		FROM meal_plans WHERE id = $1
	`, id).Scan(&m.ID, &m.UserID, &m.Title, &m.MealType, &m.Servings, &m.Notes, &m.CreatedAt, &m.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &m, nil
}

func (r *MealPlanRepository) Create(ctx context.Context, userID int64, title string, mealType string, servings int, notes string) (*model.MealPlan, error) {
	var m model.MealPlan
	err := r.pool.QueryRow(ctx, `
		INSERT INTO meal_plans (user_id, title, meal_type, servings, notes)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, user_id, title, meal_type, servings, notes, created_at, updated_at
	`, userID, title, mealType, servings, notes).Scan(
		&m.ID, &m.UserID, &m.Title, &m.MealType, &m.Servings, &m.Notes, &m.CreatedAt, &m.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &m, nil
}

func (r *MealPlanRepository) Update(ctx context.Context, id int64, title string, mealType string, servings int, notes string) (*model.MealPlan, error) {
	var m model.MealPlan
	err := r.pool.QueryRow(ctx, `
		UPDATE meal_plans SET title = $2, meal_type = $3, servings = $4, notes = $5
		WHERE id = $1
		RETURNING id, user_id, title, meal_type, servings, notes, created_at, updated_at
	`, id, title, mealType, servings, notes).Scan(
		&m.ID, &m.UserID, &m.Title, &m.MealType, &m.Servings, &m.Notes, &m.CreatedAt, &m.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &m, nil
}

func (r *MealPlanRepository) Delete(ctx context.Context, id int64) error {
	ct, err := r.pool.Exec(ctx, `DELETE FROM meal_plans WHERE id = $1`, id)
	if err != nil {
		return err
	}
	if ct.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}
	return nil
}
