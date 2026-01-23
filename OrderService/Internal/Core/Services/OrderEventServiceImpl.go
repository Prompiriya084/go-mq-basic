package services

import (
	"time"

	ports_repositories "github.com/Prompiriya084/go-mq/OrderService/Internal/Core/Ports/Repositories"
	models "github.com/Prompiriya084/go-mq/OrderService/Models"
	models_eventbus "github.com/Prompiriya084/go-mq/OrderService/Models/Eventbus"
	"github.com/google/uuid"
)

type orderEventServiceImpl struct {
	repo ports_repositories.OrderRepository
}

func NewOrderEventService(repo ports_repositories.OrderRepository) OrderEventService {
	return &orderEventServiceImpl{
		repo: repo,
	}
}

func (s *orderEventServiceImpl) OrderCreated(orderEvent *models_eventbus.OrderCreated) error {
	uuidEventOrderID, _ := uuid.Parse(orderEvent.ID)
	updatedOrder, err := s.repo.Get(&models.Order{
		ID: uuidEventOrderID,
	}, nil)

	if err != nil {
		return err
	}

	updatedOrder.Status = "Completed"
	updatedOrder.UpdatedAt = time.Now()

	if err := s.repo.Update(updatedOrder); err != nil {
		return err
	}

	return nil
}
