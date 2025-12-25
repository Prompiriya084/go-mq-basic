package models

import (
	"time"

	"gorm.io/gorm"
)

type Inventory struct {
	ProductID string         `json:"product_id" gorm:"primaryKey" validate:"required"`
	Qty       uint           `json:"qty" validate:"required,gt=0"`
	CreatedAt time.Time      `json:"created_at" swaggerignore:"true"`
	UpdatedAt time.Time      `json:"updated_at" swaggerignore:"true"`
	DeletedAt gorm.DeletedAt `json:"deleted_at" swaggerignore:"true" gorm:"index"`
}
