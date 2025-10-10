package services

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"time"

	eventbus "github.com/Prompiriya084/go-mq/EventBus"
	ports_repositories "github.com/Prompiriya084/go-mq/InventoryService/Core/Ports/Repositories"
	models "github.com/Prompiriya084/go-mq/Models"
)

type inventoryServiceImpl struct {
	repo ports_repositories.InventoryRepository
	bus  eventbus.EventBus[models.Order]
}

func NewInventoryService(repo ports_repositories.InventoryRepository, bus eventbus.EventBus[models.Order]) InventoryService {
	return &inventoryServiceImpl{
		repo: repo,
		bus:  bus,
	}
}
func (s *inventoryServiceImpl) checkStockFailed(order *models.Order) error {
	byteMessage, err := json.Marshal(order)
	if err != nil {
		return err
	}
	if err := s.bus.Publish("inventory.failed", byteMessage); err != nil {
		return err
	}

	return nil
}
func (s *inventoryServiceImpl) checkStockSuccessful(order *models.Order) error {
	byteMessage, err := json.Marshal(order)
	if err != nil {
		return err
	}
	existingItemInStock, err := s.Get(&models.Inventory{
		ProductID: order.ProductID,
	})
	if err != nil {
		return err
	}
	existingItemInStock.Qty = existingItemInStock.Qty - order.Qty
	existingItemInStock.UpdatedAt = time.Now()

	if err := s.Update(existingItemInStock); err != nil {
		return err
	}

	if err := s.bus.Publish("inventory.checked", byteMessage); err != nil {
		return err
	}

	return nil
}

func (s *inventoryServiceImpl) CheckStock(order *models.Order) error {
	existingItemInStock, err := s.Get(&models.Inventory{
		ProductID: order.ProductID,
	})
	if err != nil {
		return err
	}
	if existingItemInStock.Qty >= order.Qty {
		log.Printf("Check stock %v successful.", order)
		s.checkStockSuccessful(order)

	} else {
		log.Printf("Check stock %v failed.", order)
		s.checkStockFailed(order)
	}

	return nil
}
func (s *inventoryServiceImpl) ReverseStock(order *models.Order) error {
	existingItemInStock, err := s.Get(&models.Inventory{
		ProductID: order.ProductID,
	})
	if err != nil {
		return err
	}
	existingItemInStock.Qty = existingItemInStock.Qty + order.Qty
	existingItemInStock.UpdatedAt = time.Now()

	if err := s.Update(existingItemInStock); err != nil {
		return err
	}

	s.checkStockFailed(order)

	return nil
}
func (s *inventoryServiceImpl) Get(invent *models.Inventory) (*models.Inventory, error) {
	existingItemInStock, err := s.repo.Get(invent, nil)
	if err != nil {
		return nil, err
	}

	return existingItemInStock, nil
}
func (s *inventoryServiceImpl) Create(invent *models.Inventory) error {
	existingItemInStock, err := s.Get(&models.Inventory{
		ProductID: invent.ProductID,
	})

	if err != nil {
		return fmt.Errorf("Creating new stock failed: %w", err)
	}

	if existingItemInStock != nil {
		return errors.New("The Product ID : " + existingItemInStock.ProductID + " already exists.")
	}

	if err := s.repo.Add(invent); err != nil {
		return fmt.Errorf("Failed to create the product ID : "+existingItemInStock.ProductID+" into inventory : %w", err)
	}
	return nil
}
func (s *inventoryServiceImpl) Update(invent *models.Inventory) error {
	existingItemInStock, err := s.Get(&models.Inventory{
		ID: invent.ID,
	})
	if err != nil {
		return fmt.Errorf("failed to get order: %w", err)
	}
	if existingItemInStock != nil {
		return errors.New("The order" + invent.ID.String() + " already exists.")
	}
	if err := s.repo.Update(invent); err != nil {
		return err
	}
	return nil
}
func (s *inventoryServiceImpl) Delete(invent *models.Inventory) error {
	existingItemInStock, err := s.Get(&models.Inventory{
		ID: invent.ID,
	})
	if err != nil {
		return fmt.Errorf("failed to get order: %w", err)
	}
	if existingItemInStock == nil {
		return errors.New("The order" + invent.ID.String() + " doesn't exist.")
	}
	if err := s.repo.Delete(invent); err != nil {
		return err
	}
	return nil
}
