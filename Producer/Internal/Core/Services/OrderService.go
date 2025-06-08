package services

import models "github.com/Prompiriya084/go-mq/Models"

type OrderService interface {
	GetAll(filters *models.Order, preload []string) ([]models.Order, error)
	Get(filter *models.Order, preload []string) (*models.Order, error)
	Create(order *models.Order) error
}
