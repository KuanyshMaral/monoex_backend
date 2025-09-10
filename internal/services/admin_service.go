package services

import (
	"context"
	"database/sql"
	"errors"

	"golang.org/x/crypto/bcrypt"
)

type AdminService struct {
	DB *sql.DB
}

func NewAdminService(db *sql.DB) *AdminService {
	return &AdminService{DB: db}
}

// Проверка, создан ли админ
func (s *AdminService) IsAdminCreated(ctx context.Context) (bool, error) {
	var count int
	err := s.DB.QueryRowContext(ctx, "SELECT COUNT(*) FROM admins").Scan(&count)
	return count > 0, err
}

// Создание единственного админа
func (s *AdminService) CreateAdmin(ctx context.Context, username, password string) error {
	exists, err := s.IsAdminCreated(ctx)
	if err != nil {
		return err
	}
	if exists {
		return errors.New("admin already exists")
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	_, err = s.DB.ExecContext(ctx,
		"INSERT INTO admins (username, password) VALUES ($1, $2)",
		username, string(hash))
	return err
}

// Проверка логина/пароля админа для middleware
func (s *AdminService) ValidateAdmin(username, password string) bool {
	var hash string
	err := s.DB.QueryRow("SELECT password FROM admins WHERE username=$1", username).Scan(&hash)
	if err != nil {
		return false
	}
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)) == nil
}
