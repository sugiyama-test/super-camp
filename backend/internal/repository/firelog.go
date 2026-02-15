package repository

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sugiyamadaiki/super-camp/backend/internal/model"
)

type FireLogRepository struct {
	pool *pgxpool.Pool
}

func NewFireLogRepository(pool *pgxpool.Pool) *FireLogRepository {
	return &FireLogRepository{pool: pool}
}

func (r *FireLogRepository) List(ctx context.Context, userID int64) ([]model.FireLog, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT id, user_id, date, location, wood_type, duration_minutes, notes, temperature, campsite_id, created_at, updated_at
		FROM fire_logs WHERE user_id = $1
		ORDER BY date DESC, created_at DESC
	`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []model.FireLog
	for rows.Next() {
		var l model.FireLog
		if err := rows.Scan(&l.ID, &l.UserID, &l.Date, &l.Location, &l.WoodType, &l.DurationMinutes, &l.Notes, &l.Temperature, &l.CampsiteID, &l.CreatedAt, &l.UpdatedAt); err != nil {
			return nil, err
		}
		result = append(result, l)
	}
	return result, rows.Err()
}

func (r *FireLogRepository) Get(ctx context.Context, id int64) (*model.FireLog, error) {
	var l model.FireLog
	err := r.pool.QueryRow(ctx, `
		SELECT id, user_id, date, location, wood_type, duration_minutes, notes, temperature, campsite_id, created_at, updated_at
		FROM fire_logs WHERE id = $1
	`, id).Scan(&l.ID, &l.UserID, &l.Date, &l.Location, &l.WoodType, &l.DurationMinutes, &l.Notes, &l.Temperature, &l.CampsiteID, &l.CreatedAt, &l.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &l, nil
}

func (r *FireLogRepository) Create(ctx context.Context, userID int64, date string, location string, woodType string, durationMinutes int, notes string, temperature *float64) (*model.FireLog, error) {
	var l model.FireLog
	err := r.pool.QueryRow(ctx, `
		INSERT INTO fire_logs (user_id, date, location, wood_type, duration_minutes, notes, temperature)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id, user_id, date, location, wood_type, duration_minutes, notes, temperature, campsite_id, created_at, updated_at
	`, userID, date, location, woodType, durationMinutes, notes, temperature).Scan(
		&l.ID, &l.UserID, &l.Date, &l.Location, &l.WoodType, &l.DurationMinutes, &l.Notes, &l.Temperature, &l.CampsiteID, &l.CreatedAt, &l.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &l, nil
}

func (r *FireLogRepository) Update(ctx context.Context, id int64, date string, location string, woodType string, durationMinutes int, notes string, temperature *float64) (*model.FireLog, error) {
	var l model.FireLog
	err := r.pool.QueryRow(ctx, `
		UPDATE fire_logs SET date = $2, location = $3, wood_type = $4, duration_minutes = $5, notes = $6, temperature = $7
		WHERE id = $1
		RETURNING id, user_id, date, location, wood_type, duration_minutes, notes, temperature, campsite_id, created_at, updated_at
	`, id, date, location, woodType, durationMinutes, notes, temperature).Scan(
		&l.ID, &l.UserID, &l.Date, &l.Location, &l.WoodType, &l.DurationMinutes, &l.Notes, &l.Temperature, &l.CampsiteID, &l.CreatedAt, &l.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &l, nil
}

func (r *FireLogRepository) Delete(ctx context.Context, id int64) error {
	ct, err := r.pool.Exec(ctx, `DELETE FROM fire_logs WHERE id = $1`, id)
	if err != nil {
		return err
	}
	if ct.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}
	return nil
}
