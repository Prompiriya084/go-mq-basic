package main

import (
	"os"

	eventbus "github.com/Prompiriya084/go-mq/Eventbus"

	adapters_handlers "github.com/Prompiriya084/go-mq/InventoryService/Internal/Adapters/Handlers"
	database "github.com/Prompiriya084/go-mq/InventoryService/Internal/Infrastructure/Database"
	utilities_validator "github.com/Prompiriya084/go-mq/InventoryService/Internal/Utilities/Validator"
	models "github.com/Prompiriya084/go-mq/InventoryService/Models"
	routes "github.com/Prompiriya084/go-mq/InventoryService/Web/Routes"

	"github.com/gofiber/fiber/v2"

	adapters_repositories "github.com/Prompiriya084/go-mq/InventoryService/Internal/Adapters/Repositories"
	services "github.com/Prompiriya084/go-mq/InventoryService/Internal/Core/Services"
)

func main() {
	app := fiber.New()
	db := database.InitDb()

	repo := adapters_repositories.NewInventoryRepository(db)
	mqEventbus := eventbus.NewMQEventbus[models.Order](os.Getenv("RABBITMQ_URL"))
	inventoryService := services.NewInventoryService(repo, mqEventbus)

	validator := utilities_validator.NewValidator()
	inventoryHandler := adapters_handlers.NewInventoryHandler(inventoryService, mqEventbus, validator)

	routes.InventorySetupRouter(app, inventoryHandler)

	go inventoryHandler.CheckStock()
	go inventoryHandler.ReverseStock()
	app.Listen(":8081")
}
