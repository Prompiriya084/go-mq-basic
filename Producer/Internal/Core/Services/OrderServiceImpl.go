package services

import (
	"encoding/json"
	"errors"
	"time"

	models "github.com/Prompiriya084/go-mq/Models"
	ports_mq "github.com/Prompiriya084/go-mq/Producer/Internal/Core/Ports/MQ"
	ports_repositories "github.com/Prompiriya084/go-mq/Producer/Internal/Core/Ports/Repositories"
	"github.com/google/uuid"
)

type orderServiceImpl struct {
	repo       ports_repositories.OrderRepository
	mqProducer ports_mq.MQProducer
}

func NewOrderService(repo ports_repositories.OrderRepository, mqProducer ports_mq.MQProducer) OrderService {
	return &orderServiceImpl{
		repo:       repo,
		mqProducer: mqProducer,
	}
}

func (s *orderServiceImpl) GetAll(filters *models.Order, preload []string) ([]*models.Order, error) {
	return s.repo.GetAll(filters, preload)
}

func (s *orderServiceImpl) Get(filters *models.Order, preload []string) (*models.Order, error) {
	return s.repo.Get(filters, preload)
}
func (s *orderServiceImpl) Create(order *models.Order) error {
	order.ID = uuid.New()
	order.CreatedAt = time.Now()
	order.UpdatedAt = time.Now()
	byteMessage, err := json.Marshal(order)
	if err != nil {
		return err
	}
	if err := s.mqProducer.PublishMessage("order_create", byteMessage); err != nil {
		return err
	}
	return nil
}
func (s *orderServiceImpl) Update(order *models.Order) error {
	selectedOrder, err := s.Get(&models.Order{
		ID: order.ID,
	}, nil)
	if err != nil {
		return err
	}
	if selectedOrder == nil {
		return errors.New("Order not found.")
	}
	selectedOrder.ProductID = order.ProductID
	selectedOrder.UpdatedAt = time.Now()

	byteMessage, err := json.Marshal(selectedOrder)
	if err != nil {
		return err
	}
	if err := s.mqProducer.PublishMessage("order_update", byteMessage); err != nil {
		return err
	}
	return nil
}
func (s *orderServiceImpl) Delete(order *models.Order) error {
	selectedOrder, err := s.Get(&models.Order{
		ID: order.ID,
	}, nil)
	if err != nil {
		return err
	}
	if selectedOrder == nil {
		return errors.New("Order not found.")
	}
	selectedOrder.DeleteAt = time.Now()

	byteMessage, err := json.Marshal(selectedOrder)
	if err != nil {
		return err
	}
	if err := s.mqProducer.PublishMessage("order_delete", byteMessage); err != nil {
		return err
	}
	return nil
}
