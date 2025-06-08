package main

import (
	"os"

	adapters_handlers "github.com/Prompiriya084/go-mq/Customer/Adapters/Handlers"
	adapters_customers "github.com/Prompiriya084/go-mq/Customer/Adapters/MQ"
	adapters_repositories "github.com/Prompiriya084/go-mq/Customer/Adapters/Repositories"
	services "github.com/Prompiriya084/go-mq/Customer/Core/Services"
	database "github.com/Prompiriya084/go-mq/Infrastructure/Database"
	models "github.com/Prompiriya084/go-mq/Models"
)

func main() {
	db := database.InitDb()

	repo := adapters_repositories.NewOrderRepository(db)
	orderService := services.NewOrderService(repo)
	mqOrderCustomer := adapters_customers.NewMQCustomer[models.Order](os.Getenv("RABBITMQ_URL"))

	orderHandler := adapters_handlers.NewOrderHandler(orderService, mqOrderCustomer)
	orderHandler.Create()
}
