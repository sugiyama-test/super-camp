package repository

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sugiyamadaiki/super-camp/backend/internal/model"
)

type CampsiteRepository struct {
	pool *pgxpool.Pool
}

func NewCampsiteRepository(pool *pgxpool.Pool) *CampsiteRepository {
	return &CampsiteRepository{pool: pool}
}

func (r *CampsiteRepository) List(ctx context.Context) ([]model.Campsite, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT id, name, address, latitude, longitude, notes, created_at, updated_at
		FROM campsites
		ORDER BY updated_at DESC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []model.Campsite
	for rows.Next() {
		var c model.Campsite
		if err := rows.Scan(&c.ID, &c.Name, &c.Address, &c.Latitude, &c.Longitude, &c.Notes, &c.CreatedAt, &c.UpdatedAt); err != nil {
			return nil, err
		}
		result = append(result, c)
	}
	return result, rows.Err()
}

func (r *CampsiteRepository) Get(ctx context.Context, id int64) (*model.Campsite, error) {
	var c model.Campsite
	err := r.pool.QueryRow(ctx, `
		SELECT id, name, address, latitude, longitude, notes, created_at, updated_at
		FROM campsites WHERE id = $1
	`, id).Scan(&c.ID, &c.Name, &c.Address, &c.Latitude, &c.Longitude, &c.Notes, &c.CreatedAt, &c.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &c, nil
}

func (r *CampsiteRepository) Create(ctx context.Context, name string, address string, latitude *float64, longitude *float64, notes string) (*model.Campsite, error) {
	var c model.Campsite
	err := r.pool.QueryRow(ctx, `
		INSERT INTO campsites (name, address, latitude, longitude, notes)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, name, address, latitude, longitude, notes, created_at, updated_at
	`, name, address, latitude, longitude, notes).Scan(
		&c.ID, &c.Name, &c.Address, &c.Latitude, &c.Longitude, &c.Notes, &c.CreatedAt, &c.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &c, nil
}

func (r *CampsiteRepository) Update(ctx context.Context, id int64, name string, address string, latitude *float64, longitude *float64, notes string) (*model.Campsite, error) {
	var c model.Campsite
	err := r.pool.QueryRow(ctx, `
		UPDATE campsites SET name = $2, address = $3, latitude = $4, longitude = $5, notes = $6
		WHERE id = $1
		RETURNING id, name, address, latitude, longitude, notes, created_at, updated_at
	`, id, name, address, latitude, longitude, notes).Scan(
		&c.ID, &c.Name, &c.Address, &c.Latitude, &c.Longitude, &c.Notes, &c.CreatedAt, &c.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &c, nil
}

func (r *CampsiteRepository) Delete(ctx context.Context, id int64) error {
	ct, err := r.pool.Exec(ctx, `DELETE FROM campsites WHERE id = $1`, id)
	if err != nil {
		return err
	}
	if ct.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}
	return nil
}
