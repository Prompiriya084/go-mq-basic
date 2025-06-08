package services

import models "github.com/Prompiriya084/go-mq/Models"

type ProductService interface {
	Create(product *models.Product) error
	GetAll() ([]models.Product, error)
}
