package services

import (
	ports_repositories "github.com/Prompiriya084/go-mq/Customer/Core/Ports/Repositories"
	models "github.com/Prompiriya084/go-mq/Models"
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

	if err := s.repo.Add(order); err != nil {
		return err
	}
	return nil
}
