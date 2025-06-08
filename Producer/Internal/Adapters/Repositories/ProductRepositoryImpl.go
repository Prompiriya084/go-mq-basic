package adapters_repositories

import (
	models "github.com/Prompiriya084/go-mq/Models"
	ports_repositories "github.com/Prompiriya084/go-mq/Producer/Internal/Core/Ports/Repositories"
	"gorm.io/gorm"
)

type productRepositoryImpl struct {
	*repositoryImpl[models.Product]
}

func NewProductRepository(db *gorm.DB) ports_repositories.ProductRepository {
	return &productRepositoryImpl{
		repositoryImpl: newRepository[models.Product](db),
	}
}
