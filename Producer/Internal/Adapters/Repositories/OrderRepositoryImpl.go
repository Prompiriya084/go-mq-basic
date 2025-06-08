package adapters_repositories

import (
	models "github.com/Prompiriya084/go-mq/Models"
	ports_repositories "github.com/Prompiriya084/go-mq/Producer/Internal/Core/Ports/Repositories"
	"gorm.io/gorm"
)

type orderRepositoryImpl struct {
	*repositoryImpl[models.Order]
}

func NewOrderRepository(db *gorm.DB) ports_repositories.OrderRepository {
	return &orderRepositoryImpl{
		repositoryImpl: newRepository[models.Order](db),
	}
}
