package services

import (
	"encoding/json"
	"fmt"
	"time"

	eventbus "github.com/Prompiriya084/go-mq/EventBus"

	ports_repositories "github.com/Prompiriya084/go-mq/OrderService/Internal/Core/Ports/Repositories"
	models "github.com/Prompiriya084/go-mq/OrderService/Models"
	"github.com/google/uuid"
)

type orderServiceImpl struct {
	repo ports_repositories.OrderRepository
	bus  eventbus.EventBus[models.Order]
}

func NewOrderService(repo ports_repositories.OrderRepository, bus eventbus.EventBus[models.Order]) OrderService {
	return &orderServiceImpl{
		repo: repo,
		bus:  bus,
	}
}
func (s *orderServiceImpl) GetAll(filters *models.Order, preload []string) ([]*models.Order, error) {
	return s.repo.GetAll(filters, preload)
}

func (s *orderServiceImpl) Get(filters *models.Order, preload []string) (*models.Order, error) {
	return s.repo.Get(filters, preload)
}
func (s *orderServiceImpl) Create(order *models.Order) (string, error) {
	errMessage := "Failed to create new order :"

	order.ID = uuid.New()
	order.Status = "PENDING"
	order.CreatedAt = time.Now()
	order.UpdatedAt = time.Now()
	if err := s.repo.Add(order); err != nil {
		return "", fmt.Errorf("%s %w", errMessage, err)
	}

	byteMessage, err := json.Marshal(&models.Order{
		ID:        order.ID,
		ProductID: order.ProductID,
		Qty:       order.Qty,
	})
	if err != nil {
		return "", fmt.Errorf("%s %w", errMessage, err)
	}
	if err := s.bus.Publish("order.created", byteMessage); err != nil {
		return "", fmt.Errorf("%s %w", errMessage, err)
	}

	return order.ID.String(), nil
}
func (s *orderServiceImpl) Update(order *models.Order) error {
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
func (s *orderServiceImpl) Delete(orderId string) error {
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

// For Message queue
func (s *orderServiceImpl) InventoryConfirmed(evt *models.Order) error {
	evt.Status = "SUCCESS"
	evt.UpdatedAt = time.Now()
	if err := s.repo.Update(evt); err != nil {
		return err
	}

	return nil
}
