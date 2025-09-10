package repositories

import (
	"database/sql"
	"monoex_backend/internal/models"
)

type AdminRepository struct {
	db *sql.DB
}

func NewAdminRepository(db *sql.DB) *AdminRepository {
	return &AdminRepository{db: db}
}

func (r *AdminRepository) GetAdmin() (*models.Admin, error) {
	admin := &models.Admin{}
	err := r.db.QueryRow(`SELECT id, username, password FROM admins LIMIT 1`).Scan(&admin.ID, &admin.Username, &admin.Password)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return admin, nil
}

func (r *AdminRepository) CreateAdmin(username, password string) error {
	_, err := r.db.Exec(`INSERT INTO admins(username, password) VALUES($1, $2)`, username, password)
	return err
}
