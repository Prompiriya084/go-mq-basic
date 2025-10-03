package main

import (
	"os"

	eventbus "github.com/Prompiriya084/go-mq/Eventbus"
	database "github.com/Prompiriya084/go-mq/Infrastructure/Database"
	models "github.com/Prompiriya084/go-mq/Models"
	adapters_handlers "github.com/Prompiriya084/go-mq/OrderConsumer/Adapters/Handlers"

	adapters_repositories "github.com/Prompiriya084/go-mq/OrderConsumer/Adapters/Repositories"
	services "github.com/Prompiriya084/go-mq/OrderConsumer/Core/Services"
)

func main() {
	db := database.InitDb()

	repo := adapters_repositories.NewOrderRepository(db)
	orderService := services.NewOrderService(repo)
	mqEventbus := eventbus.NewMQEventbus[models.Order](os.Getenv("RABBITMQ_URL"))

	orderHandler := adapters_handlers.NewOrderHandler(orderService, mqEventbus)
	orderHandler.Create()
	orderHandler.Update()
	orderHandler.Cancel()
}
