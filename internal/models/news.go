package models

import "time"

type News struct {
	ID          int       `json:"id" db:"id"`
	Title       string    `json:"title" db:"title"`
	Description string    `json:"description" db:"description"`
	FullText    string    `json:"full_text" db:"full_text"`
	ImagePath   string    `json:"image_path" db:"image_path"`
	Status      string    `json:"status" db:"status"` // published/not_published
	Link        string    `json:"link" db:"link"`     // автогенерация
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}
