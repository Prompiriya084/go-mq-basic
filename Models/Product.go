package models

import (
	"time"

	"github.com/google/uuid"
)

type Product struct {
	ID        uuid.UUID `json:"id" swaggerignore:"true"`
	Name      string    `json:"name" validate:"required"`
	Price     uint      `json:"price" validate:"required"`
	CreatedAt time.Time `json:"created_at" swaggerignore:"true"`
	UpdatedAt time.Time `json:"updated_at" swaggerignore:"true"`
	DeleteAt  time.Time `json:"delete_at" swaggerignore:"true"`
}
