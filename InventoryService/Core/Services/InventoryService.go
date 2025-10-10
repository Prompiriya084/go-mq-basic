package services

import models "github.com/Prompiriya084/go-mq/Models"

type InventoryService interface {
	Get(invent *models.Inventory) (*models.Inventory, error)
	CheckStock(order *models.Order) error
	ReverseStock(order *models.Order) error
	Create(invent *models.Inventory) error
	Update(invent *models.Inventory) error
	Delete(invent *models.Inventory) error
}
