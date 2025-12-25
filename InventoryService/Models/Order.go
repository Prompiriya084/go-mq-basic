package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Order struct {
	ID        uuid.UUID      `json:"id" swaggerignore:"true"`
	ProductID string         `json:"product_id" validate:"required"`
	Qty       uint           `json:"qty" validate:"required"`
	Status    string         `json:"status" swaggerignore:"true"`
	CreatedAt time.Time      `json:"created_at" swaggerignore:"true"`
	UpdatedAt time.Time      `json:"updated_at" swaggerignore:"true"`
	DeletedAt gorm.DeletedAt `json:"deleted_at" swaggerignore:"true" gorm:"index"`
}
