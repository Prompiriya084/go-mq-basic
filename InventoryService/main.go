package main

import (
	"os"

	eventbus "github.com/Prompiriya084/go-mq/Eventbus"
	database "github.com/Prompiriya084/go-mq/Infrastructure/Database"
	adapters_handlers "github.com/Prompiriya084/go-mq/InventoryService/Adapters/Handlers"
	models "github.com/Prompiriya084/go-mq/Models"

	adapters_repositories "github.com/Prompiriya084/go-mq/InventoryService/Adapters/Repositories"
	services "github.com/Prompiriya084/go-mq/InventoryService/Core/Services"
)

func main() {
	db := database.InitDb()

	repo := adapters_repositories.NewInventoryRepository(db)
	mqEventbus := eventbus.NewMQEventbus[models.Order](os.Getenv("RABBITMQ_URL"))
	orderService := services.NewInventoryService(repo, mqEventbus)

	inventoryHandler := adapters_handlers.NewInventoryHandler(orderService, mqEventbus)
	inventoryHandler.CheckStock()
	inventoryHandler.ReverseStock()
}
