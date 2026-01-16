package events

type OrderCreated struct {
	ID        string `validate:"required"`
	ProductID string `validate:"required"`
	Qty       uint   `validate:"required"`
}
