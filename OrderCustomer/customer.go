package main

import (
	"os"

	database "github.com/Prompiriya084/go-mq/Infrastructure/Database"
	models "github.com/Prompiriya084/go-mq/Models"
	adapters_handlers "github.com/Prompiriya084/go-mq/OrderCustomer/Adapters/Handlers"
	adapters_customers "github.com/Prompiriya084/go-mq/OrderCustomer/Adapters/MQ"
	adapters_repositories "github.com/Prompiriya084/go-mq/OrderCustomer/Adapters/Repositories"
	services "github.com/Prompiriya084/go-mq/OrderCustomer/Core/Services"
)

func main() {
	db := database.InitDb()

	repo := adapters_repositories.NewOrderRepository(db)
	orderService := services.NewOrderService(repo)
	mqOrderCustomer := adapters_customers.NewMQCustomer[models.Order](os.Getenv("RABBITMQ_URL"))

	orderHandler := adapters_handlers.NewOrderHandler(orderService, mqOrderCustomer)
	orderHandler.Create()
	orderHandler.Update()
	orderHandler.Delete()
}
