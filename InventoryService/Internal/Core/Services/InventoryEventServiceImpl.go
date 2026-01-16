package services

import (
	"fmt"
	"time"

	ports_eventbus "github.com/Prompiriya084/go-mq/InventoryService/Internal/Core/Ports/Eventbus"
	ports_repositories "github.com/Prompiriya084/go-mq/InventoryService/Internal/Core/Ports/Repositories"
	models "github.com/Prompiriya084/go-mq/InventoryService/Models"
	events "github.com/Prompiriya084/go-mq/InventoryService/Models/Events"
)

type inventoryEventServiceImpl struct {
	repo     ports_repositories.InventoryRepository
	eventbus ports_eventbus.EventBusPublisher
}

func NewInventoryEventService(repo ports_repositories.InventoryRepository, eventbus ports_eventbus.EventBusPublisher) InventoryEventService {
	return &inventoryEventServiceImpl{
		repo:     repo,
		eventbus: eventbus,
	}
}

func (s *inventoryEventServiceImpl) ReserveStock(order *events.OrderCreated) error {
	selectedStock, err := s.repo.Get(&models.Inventory{
		ProductID: order.ProductID,
	}, nil)
	if err != nil {
		return err
	}
	if selectedStock == nil {
		return fmt.Errorf("The product not found.")
	}

	if selectedStock.Qty < order.Qty {
		return fmt.Errorf("Not enough stock")
	}

	//Reserve stock
	selectedStock.Qty -= order.Qty
	selectedStock.UpdatedAt = time.Now()

	if err := s.repo.Update(selectedStock); err != nil {
		return err
	}

	//Publish event
	event := events.StockReserved{
		OrderID:   order.ID,
		ProductID: order.ProductID,
		Qty:       order.Qty,
		Timestamp: time.Now(),
	}

	return s.eventbus.Publish("stock.reserved", event)
}

func (s *inventoryEventServiceImpl) RollbackStock(rollbackData *models.Inventory) error {
	return nil
}
