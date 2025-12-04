package services

import models "github.com/Prompiriya084/go-mq/OrderService/Models"

type OrderService interface {
	GetAll(filters *models.Order, preload []string) ([]*models.Order, error)
	Get(filter *models.Order, preload []string) (*models.Order, error)
	Create(order *models.Order) error
	Update(order *models.Order) error
	Delete(order *models.Order) error
}
