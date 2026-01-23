package services

import models_eventbus "github.com/Prompiriya084/go-mq/OrderService/Models/Eventbus"

type OrderEventService interface {
	OrderCreated(order *models_eventbus.OrderCreated) error
}
