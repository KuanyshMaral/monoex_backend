package repositories

import (
	"context"
	"database/sql"
	"monoex_backend/internal/models"
)

type LegislationRepository struct {
	db *sql.DB
}

func NewLegislationRepository(db *sql.DB) *LegislationRepository {
	return &LegislationRepository{db: db}
}

func (r *LegislationRepository) Create(ctx context.Context, l *models.Legislation) error {
	return r.db.QueryRowContext(ctx, `
        INSERT INTO legislations (title, description, file_path)
        VALUES ($1, $2, $3)
        RETURNING id, created_at, updated_at
    `, l.Title, l.Description, l.FilePath).Scan(&l.ID, &l.CreatedAt, &l.UpdatedAt)
}

func (r *LegislationRepository) GetByID(ctx context.Context, id int) (*models.Legislation, error) {
	var l models.Legislation
	err := r.db.QueryRowContext(ctx, `
        SELECT id, title, description, file_path, created_at, updated_at
        FROM legislations WHERE id = $1
    `, id).Scan(&l.ID, &l.Title, &l.Description, &l.FilePath, &l.CreatedAt, &l.UpdatedAt)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &l, nil
}

func (r *LegislationRepository) GetAll(ctx context.Context, limit, offset int) ([]*models.Legislation, error) {
	rows, err := r.db.QueryContext(ctx, `
        SELECT id, title, description, file_path, created_at, updated_at
        FROM legislations
        ORDER BY created_at DESC
        LIMIT $1 OFFSET $2
    `, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var legislations []*models.Legislation
	for rows.Next() {
		var l models.Legislation
		if err := rows.Scan(&l.ID, &l.Title, &l.Description, &l.FilePath, &l.CreatedAt, &l.UpdatedAt); err != nil {
			return nil, err
		}
		legislations = append(legislations, &l)
	}
	return legislations, nil
}

func (r *LegislationRepository) Update(ctx context.Context, l *models.Legislation) error {
	return r.db.QueryRowContext(ctx, `
        UPDATE legislations SET
            title = $1, description = $2, file_path = $3, updated_at = now()
        WHERE id = $4
        RETURNING updated_at
    `, l.Title, l.Description, l.FilePath, l.ID).Scan(&l.UpdatedAt)
}

func (r *LegislationRepository) Delete(ctx context.Context, id int) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM legislations WHERE id = $1`, id)
	return err
}

func (r *LegislationRepository) GetTotalCount(ctx context.Context) (int, error) {
	var count int
	err := r.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM legislations`).Scan(&count)
	return count, err
}
