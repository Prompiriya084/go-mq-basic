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
	existingOrder, err := s.Get(&models.Order{
		ID: order.ID,
	}, nil)
	if err != nil {
		return err
	}

	if existingOrder != nil {
		return errors.New("The order" + existingOrder.ID.String() + " already exists.")
	}
	order.ID = uuid.New()
	order.CreatedAt = time.Now()
	order.UpdatedAt = time.Now()
	byteMessage, err := json.Marshal(order)
	if err != nil {
		return err
	}
	if err := s.mqProducer.PublishMessage("order.create", byteMessage); err != nil {
		return err
	}
	return nil
}
func (s *orderServiceImpl) Update(order *models.Order) error {
	existingOrder, err := s.Get(&models.Order{
		ID: order.ID,
	}, nil)
	if err != nil {
		return err
	}
	if existingOrder == nil {
		return errors.New("The order " + order.ID.String() + " does not exists.")
	}
	existingOrder.ProductID = order.ProductID
	existingOrder.UpdatedAt = time.Now()

	byteMessage, err := json.Marshal(existingOrder)
	if err != nil {
		return err
	}
	if err := s.mqProducer.PublishMessage("order.update", byteMessage); err != nil {
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
		return errors.New("The order " + order.ID.String() + " does not exists.")
	}
	selectedOrder.DeleteAt = time.Now()

	byteMessage, err := json.Marshal(selectedOrder)
	if err != nil {
		return err
	}
	if err := s.mqProducer.PublishMessage("order.cancel", byteMessage); err != nil {
		return err
	}
	return nil
}
