package models

import (
	"time"
)

type Inventory struct {
	ProductID string    `json:"product_id" gorm:"primaryKey" validate:"required"`
	Qty       uint      `json:"qty" validate:"required,gt=0"`
	CreatedAt time.Time `json:"created_at" swaggerignore:"true"`
	UpdatedAt time.Time `json:"updated_at" swaggerignore:"true"`
	DeleteAt  time.Time `json:"delete_at" swaggerignore:"true"`
}
