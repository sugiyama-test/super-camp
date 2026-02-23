package repository

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sugiyamadaiki/super-camp/backend/internal/model"
)

type GearRepository struct {
	pool *pgxpool.Pool
}

func NewGearRepository(pool *pgxpool.Pool) *GearRepository {
	return &GearRepository{pool: pool}
}

func (r *GearRepository) List(ctx context.Context, userID int64) ([]model.Gear, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT id, user_id, name, category, brand, weight_grams, notes, created_at, updated_at
		FROM gears WHERE user_id = $1
		ORDER BY updated_at DESC
	`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []model.Gear
	for rows.Next() {
		var g model.Gear
		if err := rows.Scan(&g.ID, &g.UserID, &g.Name, &g.Category, &g.Brand, &g.WeightGrams, &g.Notes, &g.CreatedAt, &g.UpdatedAt); err != nil {
			return nil, err
		}
		result = append(result, g)
	}
	return result, rows.Err()
}

func (r *GearRepository) Get(ctx context.Context, id int64) (*model.Gear, error) {
	var g model.Gear
	err := r.pool.QueryRow(ctx, `
		SELECT id, user_id, name, category, brand, weight_grams, notes, created_at, updated_at
		FROM gears WHERE id = $1
	`, id).Scan(&g.ID, &g.UserID, &g.Name, &g.Category, &g.Brand, &g.WeightGrams, &g.Notes, &g.CreatedAt, &g.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &g, nil
}

func (r *GearRepository) Create(ctx context.Context, userID int64, name string, category string, brand string, weightGrams *float64, notes string) (*model.Gear, error) {
	var g model.Gear
	err := r.pool.QueryRow(ctx, `
		INSERT INTO gears (user_id, name, category, brand, weight_grams, notes)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, user_id, name, category, brand, weight_grams, notes, created_at, updated_at
	`, userID, name, category, brand, weightGrams, notes).Scan(
		&g.ID, &g.UserID, &g.Name, &g.Category, &g.Brand, &g.WeightGrams, &g.Notes, &g.CreatedAt, &g.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &g, nil
}

func (r *GearRepository) Update(ctx context.Context, id int64, name string, category string, brand string, weightGrams *float64, notes string) (*model.Gear, error) {
	var g model.Gear
	err := r.pool.QueryRow(ctx, `
		UPDATE gears SET name = $2, category = $3, brand = $4, weight_grams = $5, notes = $6
		WHERE id = $1
		RETURNING id, user_id, name, category, brand, weight_grams, notes, created_at, updated_at
	`, id, name, category, brand, weightGrams, notes).Scan(
		&g.ID, &g.UserID, &g.Name, &g.Category, &g.Brand, &g.WeightGrams, &g.Notes, &g.CreatedAt, &g.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &g, nil
}

func (r *GearRepository) Delete(ctx context.Context, id int64) error {
	ct, err := r.pool.Exec(ctx, `DELETE FROM gears WHERE id = $1`, id)
	if err != nil {
		return err
	}
	if ct.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}
	return nil
}
