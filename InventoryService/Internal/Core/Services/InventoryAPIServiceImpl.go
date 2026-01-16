package services

import (
	"fmt"
	"time"

	ports_eventbus "github.com/Prompiriya084/go-mq/InventoryService/Internal/Core/Ports/Eventbus"
	ports_repositories "github.com/Prompiriya084/go-mq/InventoryService/Internal/Core/Ports/Repositories"
	models "github.com/Prompiriya084/go-mq/InventoryService/Models"
)

type inventoryAPIServiceImpl struct {
	repo     ports_repositories.InventoryRepository
	eventbus ports_eventbus.EventBusPublisher
}

func NewInventoryAPIService(repo ports_repositories.InventoryRepository, eventbus ports_eventbus.EventBusPublisher) InventoryAPIService {
	return &inventoryAPIServiceImpl{
		repo:     repo,
		eventbus: eventbus,
	}
}

func (s *inventoryAPIServiceImpl) GetAll(filters *models.Inventory, preload []string) ([]*models.Inventory, error) {
	return s.repo.GetAll(filters, preload)
}

func (s *inventoryAPIServiceImpl) Get(filters *models.Inventory, preload []string) (*models.Inventory, error) {
	return s.repo.Get(filters, preload)
}
func (s *inventoryAPIServiceImpl) Create(invent *models.Inventory) error {
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
func (s *inventoryAPIServiceImpl) Update(updatedItem *models.Inventory) error {
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
func (s *inventoryAPIServiceImpl) Delete(productId string) error {
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

// // For event driven oparation
// func (s *inventoryAPIServiceImpl) checkStockFailed(order *models.Order) error {
// 	byteMessage, err := json.Marshal(order)
// 	if err != nil {
// 		return err
// 	}
// 	if err := s.eventbus.Publish("inventory.failed", byteMessage); err != nil {
// 		return err
// 	}

// 	return nil
// }

// func (s *inventoryAPIServiceImpl) CheckStock(order *models.Order) error {
// 	fmt.Println("checstock: ", order.ProductID)
// 	testmodel := &models.Inventory{
// 		ProductID: order.ProductID,
// 	}
// 	fmt.Println("sending model: ", testmodel)
// 	existingItemInStock, err := s.repo.Get(testmodel, nil)
// 	if err != nil {
// 		return err
// 	}
// 	//If data not found
// 	if existingItemInStock == nil {
// 		log.Printf("The product Id : %s doesn't exists.", order.ProductID)
// 		s.checkStockFailed(order)
// 		return nil
// 	}
// 	//If the existing product qty is less than order qty
// 	if existingItemInStock.Qty < order.Qty {
// 		log.Printf("Check stock %v failed.", order)
// 		s.checkStockFailed(order)
// 		return nil
// 	}

// 	s.checkStockSuccessful(order)
// 	log.Printf("Check stock %v successful.", order)

// 	return nil
// }
// func (s *inventoryAPIServiceImpl) checkStockSuccessful(order *models.Order) error {
// 	byteMessage, err := json.Marshal(order)
// 	if err != nil {
// 		return err
// 	}
// 	existingItemInStock, err := s.repo.Get(&models.Inventory{
// 		ProductID: order.ProductID,
// 	}, nil)
// 	if err != nil {
// 		return err
// 	}
// 	existingItemInStock.Qty = existingItemInStock.Qty - order.Qty
// 	existingItemInStock.UpdatedAt = time.Now()

// 	if err := s.repo.Update(existingItemInStock); err != nil {
// 		return err
// 	}

// 	if err := s.eventbus.Publish("inventory.checked", byteMessage); err != nil {
// 		return err
// 	}

// 	return nil
// }
// func (s *inventoryAPIServiceImpl) ReverseStock(order *models.Order) error {
// 	existingItemInStock, err := s.repo.Get(&models.Inventory{
// 		ProductID: order.ProductID,
// 	}, nil)
// 	if err != nil {
// 		return err
// 	}
// 	existingItemInStock.Qty = existingItemInStock.Qty + order.Qty
// 	existingItemInStock.UpdatedAt = time.Now()

// 	if err := s.repo.Update(existingItemInStock); err != nil {
// 		return err
// 	}

// 	s.checkStockFailed(order)

// 	return nil
// }
