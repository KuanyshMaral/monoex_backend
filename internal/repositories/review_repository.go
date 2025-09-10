package repositories

import (
	"context"
	"database/sql"
	"monoex_backend/internal/models"
)

type ReviewRepository struct {
	db *sql.DB
}

func NewReviewRepository(db *sql.DB) *ReviewRepository {
	return &ReviewRepository{db: db}
}

func (r *ReviewRepository) Create(ctx context.Context, review *models.Review) error {
	return r.db.QueryRowContext(ctx, `
		INSERT INTO reviews (company_name, service_type, description, pdf_path, created_at, updated_at)
		VALUES ($1, $2, $3, $4, now(), now())
		RETURNING id, created_at, updated_at
	`, review.CompanyName, review.ServiceType, review.Description, review.PDFPath).
		Scan(&review.ID, &review.CreatedAt, &review.UpdatedAt)
}

func (r *ReviewRepository) GetByID(ctx context.Context, id int) (*models.Review, error) {
	var review models.Review
	err := r.db.QueryRowContext(ctx, `
		SELECT id, company_name, service_type, description, pdf_path, created_at, updated_at
		FROM reviews
		WHERE id = $1
	`, id).Scan(
		&review.ID,
		&review.CompanyName,
		&review.ServiceType,
		&review.Description,
		&review.PDFPath,
		&review.CreatedAt,
		&review.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &review, nil
}

func (r *ReviewRepository) GetAll(ctx context.Context, limit, offset int) ([]*models.Review, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT id, company_name, service_type, description, pdf_path, created_at, updated_at
		FROM reviews
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var reviews []*models.Review
	for rows.Next() {
		var review models.Review
		if err := rows.Scan(
			&review.ID,
			&review.CompanyName,
			&review.ServiceType,
			&review.Description,
			&review.PDFPath,
			&review.CreatedAt,
			&review.UpdatedAt,
		); err != nil {
			continue
		}
		reviews = append(reviews, &review)
	}
	return reviews, nil
}

func (r *ReviewRepository) GetByServiceType(ctx context.Context, serviceType string, limit, offset int) ([]*models.Review, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT id, company_name, service_type, description, pdf_path, created_at, updated_at
		FROM reviews
		WHERE service_type = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`, serviceType, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var reviews []*models.Review
	for rows.Next() {
		var review models.Review
		if err := rows.Scan(
			&review.ID,
			&review.CompanyName,
			&review.ServiceType,
			&review.Description,
			&review.PDFPath,
			&review.CreatedAt,
			&review.UpdatedAt,
		); err != nil {
			continue
		}
		reviews = append(reviews, &review)
	}
	return reviews, nil
}

func (r *ReviewRepository) Update(ctx context.Context, review *models.Review) error {
	_, err := r.db.ExecContext(ctx, `
		UPDATE reviews SET
			company_name = $1,
			service_type = $2,
			description = $3,
			pdf_path = $4,
			updated_at = now()
		WHERE id = $5
	`, review.CompanyName, review.ServiceType, review.Description, review.PDFPath, review.ID)
	return err
}

func (r *ReviewRepository) Delete(ctx context.Context, id int) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM reviews WHERE id = $1`, id)
	return err
}

func (r *ReviewRepository) GetTotalCount(ctx context.Context) (int, error) {
	var count int
	err := r.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM reviews`).Scan(&count)
	return count, err
}
