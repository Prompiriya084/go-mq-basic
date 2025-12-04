package unittest_repositories

import models "github.com/Prompiriya084/go-mq/OrderService/Models"

type MockOrderRepo struct {
	*MockRepositoryImpl[models.Order]
}
