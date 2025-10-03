package services

import models "github.com/Prompiriya084/go-mq/Models"

type OrderService interface {
	Create(order *models.Order) error
	Update(order *models.Order) error
	Delete(order *models.Order) error
}
