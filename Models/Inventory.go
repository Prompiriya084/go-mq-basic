package models

import (
	"time"

	"github.com/google/uuid"
)

type Inventory struct {
	ID        uuid.UUID `json:"id" swaggerignore:"true"`
	ProductID string    `json:"product_id" validate:"required"`
	Qty       uint      `json:"qty" validate:"required"`
	CreatedAt time.Time `json:"created_at" swaggerignore:"true"`
	UpdatedAt time.Time `json:"updated_at" swaggerignore:"true"`
	DeleteAt  time.Time `json:"delete_at" swaggerignore:"true"`
}
