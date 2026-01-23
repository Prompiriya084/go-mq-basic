package models_eventbus

import "time"

type OrderCreated struct {
	ID        string `validate:"required"`
	ProductID string `validate:"required"`
	Qty       uint   `validate:"required"`
	Timestamp time.Time
}
