package services

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	eventbus "github.com/Prompiriya084/go-mq/Eventbus"
	ports_repositories "github.com/Prompiriya084/go-mq/InventoryService/Internal/Core/Ports/Repositories"
	models "github.com/Prompiriya084/go-mq/Models"
	"github.com/google/uuid"
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
	existingItemInStock, err := s.repo.Get(&models.Inventory{
		ProductID: order.ProductID,
	}, nil)
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
	existingItemInStock, err := s.repo.Get(&models.Inventory{
		ProductID: order.ProductID,
	}, nil)
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
	existingItemInStock, err := s.repo.Get(&models.Inventory{
		ProductID: order.ProductID,
	}, nil)
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
func (s *inventoryServiceImpl) GetAll(filters *models.Inventory, preload []string) ([]*models.Inventory, error) {
	return s.repo.GetAll(filters, preload)
}

func (s *inventoryServiceImpl) Get(filters *models.Inventory, preload []string) (*models.Inventory, error) {
	return s.repo.Get(filters, preload)
}
func (s *inventoryServiceImpl) Create(invent *models.Inventory) error {
	errMessage := "Failed to create new stock :"
	existingItemInStock, err := s.repo.Get(&models.Inventory{
		ProductID: invent.ProductID,
	}, nil)

	if err != nil {
		return fmt.Errorf("%s %w", errMessage, err)
	}

	if existingItemInStock != nil {
		return fmt.Errorf("%s The product %s already exists", errMessage, existingItemInStock.ProductID)
	}
	invent.ID = uuid.New()

	if err := s.repo.Add(invent); err != nil {
		return fmt.Errorf("%s %w", errMessage, err)
	}
	return nil
}
func (s *inventoryServiceImpl) Update(invent *models.Inventory) error {
	errMessage := "Failed to update stock :"
	existingItemInStock, err := s.repo.Get(&models.Inventory{
		ID: invent.ID,
	}, nil)
	if err != nil {
		return fmt.Errorf("%s %w", errMessage, err)
	}
	if existingItemInStock != nil {
		return fmt.Errorf("%s The order %s already exists", errMessage, invent.ID.String())
	}
	if err := s.repo.Update(invent); err != nil {
		return fmt.Errorf("%s %w", errMessage, err)
	}
	return nil
}
func (s *inventoryServiceImpl) Delete(invent *models.Inventory) error {
	errMessage := "Failed to delete stock :"
	existingItemInStock, err := s.repo.Get(&models.Inventory{
		ID: invent.ID,
	}, nil)
	if err != nil {
		return fmt.Errorf("%s %w", errMessage, err)
	}
	if existingItemInStock == nil {
		return fmt.Errorf("%s The order %s doesn't exist", errMessage, invent.ID.String())
	}
	if err := s.repo.Delete(invent); err != nil {
		return fmt.Errorf("%s %w", errMessage, err)
	}
	return nil
}
