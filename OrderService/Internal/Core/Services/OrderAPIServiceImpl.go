package services

import (
	"fmt"
	"time"

	ports_eventbus "github.com/Prompiriya084/go-mq/OrderService/Internal/Core/Ports/Eventbus"
	ports_repositories "github.com/Prompiriya084/go-mq/OrderService/Internal/Core/Ports/Repositories"
	models "github.com/Prompiriya084/go-mq/OrderService/Models"
	models_eventbus "github.com/Prompiriya084/go-mq/OrderService/Models/Eventbus"
	"github.com/google/uuid"
)

type orderAPIServiceImpl struct {
	repo     ports_repositories.OrderRepository
	eventbus ports_eventbus.EventBusPublisher
}

func NewOrderAPIService(repo ports_repositories.OrderRepository, eventbus ports_eventbus.EventBusPublisher) OrderAPIService {
	return &orderAPIServiceImpl{
		repo:     repo,
		eventbus: eventbus,
	}
}
func (s *orderAPIServiceImpl) GetAll(filters *models.Order, preload []string) ([]*models.Order, error) {
	return s.repo.GetAll(filters, preload)
}

func (s *orderAPIServiceImpl) Get(filters *models.Order, preload []string) (*models.Order, error) {
	return s.repo.Get(filters, preload)
}
func (s *orderAPIServiceImpl) Create(order *models.Order) (string, error) {
	errMessage := "Failed to create new order :"

	order.ID = uuid.New()
	order.Status = "PENDING"
	order.CreatedAt = time.Now()
	order.UpdatedAt = time.Now()
	if err := s.repo.Add(order); err != nil {
		return "", fmt.Errorf("%s %w", errMessage, err)
	}
	strOrderID := order.ID.String()
	sentEvent := &models_eventbus.OrderCreated{
		ID:        strOrderID,
		ProductID: order.ProductID,
		Qty:       order.Qty,
		Timestamp: time.Now(),
	}

	if err := s.eventbus.Publish("order.created", sentEvent); err != nil {
		return "", fmt.Errorf("%s %w", errMessage, err)
	}

	return strOrderID, nil
}
func (s *orderAPIServiceImpl) Update(order *models.Order) error {
	errMessage := "Failed to update order :"
	existingOrder, err := s.Get(&models.Order{
		ID: order.ID,
	}, nil)
	if err != nil {
		return fmt.Errorf("%s %w", errMessage, err)
	}
	if existingOrder == nil {
		return fmt.Errorf("%s The order %s does not exists", errMessage, order.ID.String())
	}
	existingOrder.ProductID = order.ProductID
	existingOrder.UpdatedAt = time.Now()

	if err := s.repo.Update(existingOrder); err != nil {
		return fmt.Errorf("%s %w", errMessage, err)
	}
	return nil
}
func (s *orderAPIServiceImpl) Delete(orderId string) error {
	errMessage := "Failed to delete order :"
	uuidOrder, err := uuid.Parse(orderId)
	if err != nil {
		return fmt.Errorf("%s %w", errMessage, err)
	}
	deletedOrder, err := s.Get(&models.Order{
		ID: uuidOrder,
	}, nil)
	if err != nil {
		return fmt.Errorf("%s %w", errMessage, err)
	}
	if deletedOrder == nil {
		return fmt.Errorf("%s The order %s does not exists", errMessage, orderId)
	}
	if err := s.repo.Delete(deletedOrder); err != nil {
		return fmt.Errorf("%s %w", errMessage, err)
	}

	return nil
}

// // For Message queue
// func (s *orderServiceImpl) InventoryConfirmed(evt *models.Order) error {
// 	evt.Status = "SUCCESS"
// 	evt.UpdatedAt = time.Now()
// 	if err := s.repo.Update(evt); err != nil {
// 		return err
// 	}

// 	return nil
// }
