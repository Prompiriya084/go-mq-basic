package models

import (
	"time"

	"github.com/google/uuid"
)

// Order represents a order in the system
// @Description Represents the Order entity in the system
// @type Order
type Order struct {
	ID        uuid.UUID `json:"id" swaggerignore:"true"`
	ProductID string    `json:"product_id"`
	CreatedAt time.Time `json:"created_at" swaggerignore:"true"`
	UpdatedAt time.Time `json:"updated_at" swaggerignore:"true"`
	DeleteAt  time.Time `json:"delete_at" swaggerignore:"true"`
}
