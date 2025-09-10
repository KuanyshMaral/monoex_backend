package services

import (
	"context"
	"database/sql"
	"errors"
	"monoex_backend/internal/models"
	"monoex_backend/internal/repositories"
)

type NewsService struct {
	repo *repositories.NewsRepository
}

func NewNewsService(repo *repositories.NewsRepository) *NewsService {
	return &NewsService{repo: repo}
}

// Create new news
func (s *NewsService) Create(ctx context.Context, n *models.News) error {
	if n.Title == "" {
		return errors.New("title is required")
	}
	if n.Status == "" {
		n.Status = "draft"
	}
	return s.repo.Create(ctx, n)
}

// Get by ID
func (s *NewsService) GetByID(ctx context.Context, id int) (*models.News, error) {
	news, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return news, nil
}

// Get by link (only published)
func (s *NewsService) GetByLink(ctx context.Context, link string) (*models.News, error) {
	news, err := s.repo.GetByLink(ctx, link)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return news, nil
}

// Get all news with pagination
func (s *NewsService) GetAll(ctx context.Context, limit, offset int) ([]*models.News, error) {
	if limit <= 0 {
		limit = 10
	}
	return s.repo.GetAll(ctx, limit, offset)
}

// Get published news with pagination
func (s *NewsService) GetPublished(ctx context.Context, limit, offset int) ([]*models.News, error) {
	if limit <= 0 {
		limit = 10
	}
	return s.repo.GetPublished(ctx, limit, offset)
}

// Update existing news
func (s *NewsService) Update(ctx context.Context, n *models.News) error {
	if n.ID == 0 {
		return errors.New("id is required for update")
	}
	return s.repo.Update(ctx, n)
}

// Update only status
func (s *NewsService) UpdateStatus(ctx context.Context, id int, status string) error {
	if id == 0 {
		return errors.New("id is required")
	}
	if status == "" {
		return errors.New("status is required")
	}
	return s.repo.UpdateStatus(ctx, id, status)
}

// Delete by ID
func (s *NewsService) Delete(ctx context.Context, id int) error {
	if id == 0 {
		return errors.New("id is required")
	}
	return s.repo.Delete(ctx, id)
}

// Get total count
func (s *NewsService) GetTotalCount(ctx context.Context) (int, error) {
	return s.repo.GetTotalCount(ctx)
}

// Get published count
func (s *NewsService) GetPublishedCount(ctx context.Context) (int, error) {
	return s.repo.GetPublishedCount(ctx)
}
