package ports_repositories

import models "github.com/Prompiriya084/go-mq/Models"

type InventoryRepository interface {
	repository[models.Inventory]
}
