package services

import models "github.com/Prompiriya084/go-mq/OrderService/Models"

type OrderService interface {
	GetAll(filters *models.Order, preload []string) ([]*models.Order, error)
	Get(filter *models.Order, preload []string) (*models.Order, error)
	Create(order *models.Order) (string, error)
	Update(order *models.Order) error
	Delete(orderId string) error
	InventoryConfirmed(evt *models.Order) error
	// InventoryConfirmedRetry(evt *models.Order) (string, error)
	// InventoryFailed(evt *models.Order) (string, error)
	// InventoryFailedRetry(evt *models.Order) (string, error)
}
