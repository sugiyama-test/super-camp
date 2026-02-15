package repository

import (
	"context"
	"encoding/json"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sugiyamadaiki/super-camp/backend/internal/model"
)

type LayoutRepository struct {
	pool *pgxpool.Pool
}

func NewLayoutRepository(pool *pgxpool.Pool) *LayoutRepository {
	return &LayoutRepository{pool: pool}
}

type LayoutSummary struct {
	ID        int64  `json:"id"`
	UserID    int64  `json:"user_id"`
	Title     string `json:"title"`
	ItemCount int    `json:"item_count"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

func (r *LayoutRepository) List(ctx context.Context, userID int64) ([]LayoutSummary, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT id, user_id, title, data, created_at, updated_at
		FROM layouts WHERE user_id = $1
		ORDER BY updated_at DESC
	`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []LayoutSummary
	for rows.Next() {
		var l model.Layout
		if err := rows.Scan(&l.ID, &l.UserID, &l.Title, &l.Data, &l.CreatedAt, &l.UpdatedAt); err != nil {
			return nil, err
		}
		itemCount := 0
		var data map[string]json.RawMessage
		if err := json.Unmarshal([]byte(l.Data), &data); err == nil {
			if items, ok := data["items"]; ok {
				var arr []json.RawMessage
				if err := json.Unmarshal(items, &arr); err == nil {
					itemCount = len(arr)
				}
			}
		}
		result = append(result, LayoutSummary{
			ID:        l.ID,
			UserID:    l.UserID,
			Title:     l.Title,
			ItemCount: itemCount,
			CreatedAt: l.CreatedAt.Format("2006-01-02T15:04:05Z"),
			UpdatedAt: l.UpdatedAt.Format("2006-01-02T15:04:05Z"),
		})
	}
	return result, rows.Err()
}

func (r *LayoutRepository) Get(ctx context.Context, id int64) (*model.Layout, error) {
	var l model.Layout
	err := r.pool.QueryRow(ctx, `
		SELECT id, user_id, title, data, created_at, updated_at
		FROM layouts WHERE id = $1
	`, id).Scan(&l.ID, &l.UserID, &l.Title, &l.Data, &l.CreatedAt, &l.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &l, nil
}

func (r *LayoutRepository) Create(ctx context.Context, title string, data string, userID int64) (*model.Layout, error) {
	var l model.Layout
	err := r.pool.QueryRow(ctx, `
		INSERT INTO layouts (title, data, user_id) VALUES ($1, $2, $3)
		RETURNING id, user_id, title, data, created_at, updated_at
	`, title, data, userID).Scan(&l.ID, &l.UserID, &l.Title, &l.Data, &l.CreatedAt, &l.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &l, nil
}

func (r *LayoutRepository) Update(ctx context.Context, id int64, title string, data string) (*model.Layout, error) {
	var l model.Layout
	err := r.pool.QueryRow(ctx, `
		UPDATE layouts SET title = $2, data = $3 WHERE id = $1
		RETURNING id, user_id, title, data, created_at, updated_at
	`, id, title, data).Scan(&l.ID, &l.UserID, &l.Title, &l.Data, &l.CreatedAt, &l.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &l, nil
}

func (r *LayoutRepository) Delete(ctx context.Context, id int64) error {
	ct, err := r.pool.Exec(ctx, `DELETE FROM layouts WHERE id = $1`, id)
	if err != nil {
		return err
	}
	if ct.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}
	return nil
}
