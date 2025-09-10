package repositories

import (
	"context"
	"database/sql"
	"monoex_backend/internal/models"
)

type NewsRepository struct {
	db *sql.DB
}

func NewNewsRepository(db *sql.DB) *NewsRepository {
	return &NewsRepository{db: db}
}

func (r *NewsRepository) Create(ctx context.Context, n *models.News) error {
	return r.db.QueryRowContext(ctx, `
		INSERT INTO news (title, description, full_text, image_path, status, link)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, created_at, updated_at
	`, n.Title, n.Description, n.FullText, n.ImagePath, n.Status, n.Link).Scan(&n.ID, &n.CreatedAt, &n.UpdatedAt)
}

func (r *NewsRepository) GetByID(ctx context.Context, id int) (*models.News, error) {
	var n models.News

	err := r.db.QueryRowContext(ctx, `
		SELECT id, title, description, full_text, image_path, status, link, created_at, updated_at
		FROM news WHERE id = $1
	`, id).Scan(&n.ID, &n.Title, &n.Description, &n.FullText, &n.ImagePath, &n.Status, &n.Link, &n.CreatedAt, &n.UpdatedAt)

	if err != nil {
		return nil, err
	}

	return &n, nil
}

func (r *NewsRepository) GetByLink(ctx context.Context, link string) (*models.News, error) {
	var n models.News

	err := r.db.QueryRowContext(ctx, `
		SELECT id, title, description, full_text, image_path, status, link, created_at, updated_at
		FROM news WHERE link = $1 AND status = 'published'
	`, link).Scan(&n.ID, &n.Title, &n.Description, &n.FullText, &n.ImagePath, &n.Status, &n.Link, &n.CreatedAt, &n.UpdatedAt)

	if err != nil {
		return nil, err
	}

	return &n, nil
}

func (r *NewsRepository) GetAll(ctx context.Context, limit, offset int) ([]*models.News, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT id, title, description, full_text, image_path, status, link, created_at, updated_at
		FROM news
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var news []*models.News

	for rows.Next() {
		var n models.News
		err := rows.Scan(&n.ID, &n.Title, &n.Description, &n.FullText, &n.ImagePath, &n.Status, &n.Link, &n.CreatedAt, &n.UpdatedAt)
		if err != nil {
			continue
		}
		news = append(news, &n)
	}

	return news, nil
}

func (r *NewsRepository) GetPublished(ctx context.Context, limit, offset int) ([]*models.News, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT id, title, description, full_text, image_path, status, link, created_at, updated_at
		FROM news
		WHERE status = 'published'
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var news []*models.News

	for rows.Next() {
		var n models.News
		err := rows.Scan(&n.ID, &n.Title, &n.Description, &n.FullText, &n.ImagePath, &n.Status, &n.Link, &n.CreatedAt, &n.UpdatedAt)
		if err != nil {
			continue
		}
		news = append(news, &n)
	}

	return news, nil
}

func (r *NewsRepository) Update(ctx context.Context, n *models.News) error {
	_, err := r.db.ExecContext(ctx, `
		UPDATE news SET
			title = $1, description = $2, full_text = $3, image_path = $4, 
			status = $5, link = $6, updated_at = now()
		WHERE id = $7
	`, n.Title, n.Description, n.FullText, n.ImagePath, n.Status, n.Link, n.ID)
	return err
}

func (r *NewsRepository) UpdateStatus(ctx context.Context, id int, status string) error {
	_, err := r.db.ExecContext(ctx, `
		UPDATE news SET status = $1, updated_at = now() WHERE id = $2
	`, status, id)
	return err
}

func (r *NewsRepository) Delete(ctx context.Context, id int) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM news WHERE id = $1`, id)
	return err
}

func (r *NewsRepository) GetTotalCount(ctx context.Context) (int, error) {
	var count int
	err := r.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM news`).Scan(&count)
	return count, err
}

func (r *NewsRepository) GetPublishedCount(ctx context.Context) (int, error) {
	var count int
	err := r.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM news WHERE status = 'published'`).Scan(&count)
	return count, err
}
