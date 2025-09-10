package services

import (
	"context"
	"database/sql"
	"errors"
	"monoex_backend/internal/models"
	"monoex_backend/internal/repositories"
)

type ReviewService struct {
	repo *repositories.ReviewRepository
}

func NewReviewService(repo *repositories.ReviewRepository) *ReviewService {
	return &ReviewService{repo: repo}
}

// Create new review
func (s *ReviewService) Create(ctx context.Context, r *models.Review) error {
	if r.CompanyName == "" || r.ServiceType == "" {
		return errors.New("company name and service type are required")
	}
	return s.repo.Create(ctx, r)
}

// Get by ID
func (s *ReviewService) GetByID(ctx context.Context, id int) (*models.Review, error) {
	review, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return review, nil
}

// Get all reviews with pagination
func (s *ReviewService) GetAll(ctx context.Context, limit, offset int) ([]*models.Review, error) {
	if limit <= 0 {
		limit = 10
	}
	return s.repo.GetAll(ctx, limit, offset)
}

// Get by service type with pagination
func (s *ReviewService) GetByServiceType(ctx context.Context, serviceType string, limit, offset int) ([]*models.Review, error) {
	if serviceType == "" {
		return nil, errors.New("service type is required")
	}
	return s.repo.GetByServiceType(ctx, serviceType, limit, offset)
}

// Update review by ID
func (s *ReviewService) Update(ctx context.Context, r *models.Review) error {
	if r.ID == 0 {
		return errors.New("id is required for update")
	}
	return s.repo.Update(ctx, r)
}

// Delete by ID
func (s *ReviewService) Delete(ctx context.Context, id int) error {
	if id == 0 {
		return errors.New("id is required")
	}
	return s.repo.Delete(ctx, id)
}

// Get total count of reviews
func (s *ReviewService) GetTotalCount(ctx context.Context) (int, error) {
	return s.repo.GetTotalCount(ctx)
}
