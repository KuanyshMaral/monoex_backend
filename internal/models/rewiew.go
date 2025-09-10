package models

import "time"

type Review struct {
	ID          int       `json:"id" db:"id"`
	CompanyName string    `json:"company_name" db:"company_name"`
	ServiceType string    `json:"service_type" db:"service_type"`
	Description string    `json:"description" db:"description"`
	PDFPath     string    `json:"pdf_path" db:"pdf_path"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}
