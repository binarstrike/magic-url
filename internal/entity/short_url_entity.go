package entity

import (
	"time"

	"github.com/google/uuid"
)

type ShortURL struct {
	Id          uuid.UUID `json:"id" db:"id"`
	OriginalURL string    `json:"original_url,omitempty" db:"original_url" validate:"http_url"`
	ShortCode   string    `json:"short_code,omitempty" db:"short_code" validate:"lte=16"`
	CreatedAt   time.Time `json:"created_at,omitempty" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at,omitempty" db:"updated_at"`
}
