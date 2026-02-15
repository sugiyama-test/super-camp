package repository

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sugiyamadaiki/super-camp/backend/internal/model"
)

type ChecklistRepository struct {
	pool *pgxpool.Pool
}

func NewChecklistRepository(pool *pgxpool.Pool) *ChecklistRepository {
	return &ChecklistRepository{pool: pool}
}

type ChecklistSummary struct {
	model.Checklist
	ItemCount    int `json:"item_count"`
	CheckedCount int `json:"checked_count"`
}

type ChecklistWithItems struct {
	model.Checklist
	Items []model.ChecklistItem `json:"items"`
}

func (r *ChecklistRepository) List(ctx context.Context, userID int64) ([]ChecklistSummary, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT c.id, c.user_id, c.title, c.created_at, c.updated_at,
			COUNT(ci.id) AS item_count,
			COUNT(CASE WHEN ci.is_checked THEN 1 END) AS checked_count
		FROM checklists c
		LEFT JOIN checklist_items ci ON ci.checklist_id = c.id
		WHERE c.user_id = $1
		GROUP BY c.id
		ORDER BY c.updated_at DESC
	`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []ChecklistSummary
	for rows.Next() {
		var s ChecklistSummary
		if err := rows.Scan(&s.ID, &s.UserID, &s.Title, &s.CreatedAt, &s.UpdatedAt, &s.ItemCount, &s.CheckedCount); err != nil {
			return nil, err
		}
		result = append(result, s)
	}
	return result, rows.Err()
}

func (r *ChecklistRepository) Get(ctx context.Context, id int64) (*ChecklistWithItems, error) {
	var c ChecklistWithItems
	err := r.pool.QueryRow(ctx, `
		SELECT id, user_id, title, created_at, updated_at
		FROM checklists WHERE id = $1
	`, id).Scan(&c.ID, &c.UserID, &c.Title, &c.CreatedAt, &c.UpdatedAt)
	if err != nil {
		return nil, err
	}

	rows, err := r.pool.Query(ctx, `
		SELECT id, checklist_id, name, is_checked, quantity, sort_order, created_at, updated_at
		FROM checklist_items WHERE checklist_id = $1
		ORDER BY sort_order, id
	`, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	c.Items = []model.ChecklistItem{}
	for rows.Next() {
		var item model.ChecklistItem
		if err := rows.Scan(&item.ID, &item.ChecklistID, &item.Name, &item.IsChecked, &item.Quantity, &item.SortOrder, &item.CreatedAt, &item.UpdatedAt); err != nil {
			return nil, err
		}
		c.Items = append(c.Items, item)
	}
	return &c, rows.Err()
}

func (r *ChecklistRepository) Create(ctx context.Context, title string, userID int64) (*model.Checklist, error) {
	var c model.Checklist
	err := r.pool.QueryRow(ctx, `
		INSERT INTO checklists (title, user_id) VALUES ($1, $2)
		RETURNING id, user_id, title, created_at, updated_at
	`, title, userID).Scan(&c.ID, &c.UserID, &c.Title, &c.CreatedAt, &c.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &c, nil
}

func (r *ChecklistRepository) Delete(ctx context.Context, id int64) error {
	ct, err := r.pool.Exec(ctx, `DELETE FROM checklists WHERE id = $1`, id)
	if err != nil {
		return err
	}
	if ct.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}
	return nil
}

func (r *ChecklistRepository) AddItem(ctx context.Context, checklistID int64, name string, quantity int) (*model.ChecklistItem, error) {
	var item model.ChecklistItem
	err := r.pool.QueryRow(ctx, `
		INSERT INTO checklist_items (checklist_id, name, quantity, sort_order)
		VALUES ($1, $2, $3, COALESCE((SELECT MAX(sort_order) + 1 FROM checklist_items WHERE checklist_id = $1), 0))
		RETURNING id, checklist_id, name, is_checked, quantity, sort_order, created_at, updated_at
	`, checklistID, name, quantity).Scan(&item.ID, &item.ChecklistID, &item.Name, &item.IsChecked, &item.Quantity, &item.SortOrder, &item.CreatedAt, &item.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &item, nil
}

func (r *ChecklistRepository) UpdateItem(ctx context.Context, id int64, name string, isChecked bool, quantity int) (*model.ChecklistItem, error) {
	var item model.ChecklistItem
	err := r.pool.QueryRow(ctx, `
		UPDATE checklist_items SET name = $2, is_checked = $3, quantity = $4
		WHERE id = $1
		RETURNING id, checklist_id, name, is_checked, quantity, sort_order, created_at, updated_at
	`, id, name, isChecked, quantity).Scan(&item.ID, &item.ChecklistID, &item.Name, &item.IsChecked, &item.Quantity, &item.SortOrder, &item.CreatedAt, &item.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &item, nil
}

func (r *ChecklistRepository) DeleteItem(ctx context.Context, id int64) error {
	ct, err := r.pool.Exec(ctx, `DELETE FROM checklist_items WHERE id = $1`, id)
	if err != nil {
		return err
	}
	if ct.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}
	return nil
}
