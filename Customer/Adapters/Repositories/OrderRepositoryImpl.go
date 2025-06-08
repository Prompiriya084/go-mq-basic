package adapters_repositories

import (
	ports_repositories "github.com/Prompiriya084/go-mq/Customer/Core/Ports/Repositories"
	models "github.com/Prompiriya084/go-mq/Models"
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
