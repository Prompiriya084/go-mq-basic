package services

import (
	models "github.com/Prompiriya084/go-mq/InventoryService/Models"
	events "github.com/Prompiriya084/go-mq/InventoryService/Models/Events"
)

type InventoryEventService interface {
	ReserveStock(order *events.OrderCreated) error
	RollbackStock(rollbackData *models.Inventory) error
}
