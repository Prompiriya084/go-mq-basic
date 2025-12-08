package services

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	eventbus "github.com/Prompiriya084/go-mq/Eventbus"
	ports_repositories "github.com/Prompiriya084/go-mq/InventoryService/Internal/Core/Ports/Repositories"
	models "github.com/Prompiriya084/go-mq/InventoryService/Models"
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

	if err := s.repo.Add(invent); err != nil {
		return fmt.Errorf("%s %w", errMessage, err)
	}
	return nil
}
func (s *inventoryServiceImpl) Update(updatedItem *models.Inventory) error {
	errMessage := "Failed to update stock :"
	existingItemInStock, err := s.repo.Get(&models.Inventory{
		ProductID: updatedItem.ProductID,
	}, nil)
	if err != nil {
		return fmt.Errorf("%s %w", errMessage, err)
	}
	if existingItemInStock == nil {
		return fmt.Errorf("%s The product %s doesn't exist", errMessage, updatedItem.ProductID)
	}
	updatedItem.UpdatedAt = time.Now()
	if err := s.repo.Update(updatedItem); err != nil {
		return fmt.Errorf("%s %w", errMessage, err)
	}
	return nil
}
func (s *inventoryServiceImpl) Delete(productId string) error {
	errMessage := "Failed to delete stock :"
	deletedItem, err := s.repo.Get(&models.Inventory{
		ProductID: productId,
	}, nil)
	if err != nil {
		return fmt.Errorf("%s %w", errMessage, err)
	}
	if deletedItem == nil {
		return fmt.Errorf("%s The product %s doesn't exist", errMessage, productId)
	}
	if err := s.repo.Delete(deletedItem); err != nil {
		return fmt.Errorf("%s %w", errMessage, err)
	}
	return nil
}

// For event driven oparation
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
func (s *inventoryServiceImpl) CheckStock(order *models.Order) error {
	fmt.Println("checstock: ", order.ProductID)
	testmodel := &models.Inventory{
		ProductID: order.ProductID,
	}
	fmt.Println("sending model: ", testmodel)
	existingItemInStock, err := s.repo.Get(testmodel, nil)
	if err != nil {
		return err
	}
	//If data not found
	if existingItemInStock == nil {
		log.Printf("The product Id : %s doesn't exists.", order.ProductID)
		s.checkStockFailed(order)
		return nil
	}
	//If the existing product qty is less than order qty
	if existingItemInStock.Qty < order.Qty {
		log.Printf("Check stock %v failed.", order)
		s.checkStockFailed(order)
		return nil
	}

	s.checkStockSuccessful(order)
	log.Printf("Check stock %v successful.", order)

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

	if err := s.repo.Update(existingItemInStock); err != nil {
		return err
	}

	if err := s.bus.Publish("inventory.checked", byteMessage); err != nil {
		return err
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

	if err := s.repo.Update(existingItemInStock); err != nil {
		return err
	}

	s.checkStockFailed(order)

	return nil
}
