package services

import (
	"errors"
	"fmt"

	models "github.com/Prompiriya084/go-mq/Models"
	ports_repositories "github.com/Prompiriya084/go-mq/OrderConsumer/Core/Ports/Repositories"
)

type orderServiceImpl struct {
	repo ports_repositories.OrderRepository
}

func NewOrderService(repo ports_repositories.OrderRepository) OrderService {
	return &orderServiceImpl{
		repo: repo,
	}
}
func (s *orderServiceImpl) Create(order *models.Order) error {
	existingOrder, err := s.repo.Get(&models.Order{
		ID: order.ID,
	}, nil)
	if err != nil {
		return fmt.Errorf("failed to get order: %w", err)
	}
	if existingOrder != nil {
		return errors.New("The order" + order.ID.String() + " already exists.")
	}
	if err := s.repo.Add(order); err != nil {
		return fmt.Errorf("failed to add order: %w", err)
	}
	return nil
}
func (s *orderServiceImpl) Update(order *models.Order) error {
	existingOrder, err := s.repo.Get(&models.Order{
		ID: order.ID,
	}, nil)
	if err != nil {
		return fmt.Errorf("failed to get order: %w", err)
	}
	if existingOrder != nil {
		return errors.New("The order" + order.ID.String() + " already exists.")
	}
	if err := s.repo.Update(order); err != nil {
		return err
	}
	return nil
}
func (s *orderServiceImpl) Delete(order *models.Order) error {
	existingOrder, err := s.repo.Get(&models.Order{
		ID: order.ID,
	}, nil)
	if err != nil {
		return fmt.Errorf("failed to get order: %w", err)
	}
	if existingOrder == nil {
		return errors.New("The order" + order.ID.String() + " doesn't exist.")
	}
	if err := s.repo.Delete(order); err != nil {
		return err
	}
	return nil
}
