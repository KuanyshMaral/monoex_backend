package services

import (
	"context"
	"database/sql"
	"errors"
	"monoex_backend/internal/models"
	"monoex_backend/internal/repositories"
)

type LegislationService struct {
	repo *repositories.LegislationRepository
}

func NewLegislationService(repo *repositories.LegislationRepository) *LegislationService {
	return &LegislationService{repo: repo}
}

// Create new legislation
func (s *LegislationService) Create(ctx context.Context, l *models.Legislation) error {
	if l.Title == "" {
		return errors.New("title is required")
	}
	return s.repo.Create(ctx, l)
}

// Get by ID
func (s *LegislationService) GetByID(ctx context.Context, id int) (*models.Legislation, error) {
	legislation, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return legislation, nil
}

// Get all with pagination
func (s *LegislationService) GetAll(ctx context.Context, limit, offset int) ([]*models.Legislation, error) {
	if limit <= 0 {
		limit = 10
	}
	return s.repo.GetAll(ctx, limit, offset)
}

// Update existing legislation
func (s *LegislationService) Update(ctx context.Context, l *models.Legislation) error {
	if l.ID == 0 {
		return errors.New("id is required for update")
	}
	return s.repo.Update(ctx, l)
}

// Delete by ID
func (s *LegislationService) Delete(ctx context.Context, id int) error {
	if id == 0 {
		return errors.New("id is required")
	}
	return s.repo.Delete(ctx, id)
}

// Get total count
func (s *LegislationService) GetTotalCount(ctx context.Context) (int, error) {
	return s.repo.GetTotalCount(ctx)
}
