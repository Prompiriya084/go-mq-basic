package services

import models "github.com/Prompiriya084/go-mq/InventoryService/Models"

type InventoryAPIService interface {
	Get(filters *models.Inventory, preload []string) (*models.Inventory, error)
	GetAll(filters *models.Inventory, preload []string) ([]*models.Inventory, error)
	Create(invent *models.Inventory) error
	Update(updatedItem *models.Inventory) error
	Delete(productId string) error
	// CheckStock(order *models.Order) error
	// ReverseStock(order *models.Order) error
}
