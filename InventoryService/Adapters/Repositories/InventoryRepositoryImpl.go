package adapters_repositories

import (
	ports_repositories "github.com/Prompiriya084/go-mq/InventoryService/Core/Ports/Repositories"
	models "github.com/Prompiriya084/go-mq/Models"
	"gorm.io/gorm"
)

type inventoryRepositoryImpl struct {
	*repositoryImpl[models.Inventory]
}

func NewInventoryRepository(db *gorm.DB) ports_repositories.InventoryRepository {
	return &inventoryRepositoryImpl{
		repositoryImpl: newRepository[models.Inventory](db),
	}
}
