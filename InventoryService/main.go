package main

import (
	"os"

	eventbus "github.com/Prompiriya084/go-mq/Eventbus"

	adapters_handlers "github.com/Prompiriya084/go-mq/InventoryService/Internal/Adapters/Handlers"
	database "github.com/Prompiriya084/go-mq/InventoryService/Internal/Infrastructure/Database"
	utilities_validator "github.com/Prompiriya084/go-mq/InventoryService/Internal/Utilities/Validator"
	models "github.com/Prompiriya084/go-mq/InventoryService/Models"
	routes "github.com/Prompiriya084/go-mq/InventoryService/Web/Routes"

	"github.com/gofiber/contrib/swagger"
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

	// Serve swagger.json
	app.Static("/docs", "./docs")
	// Correct swagger config (must NOT be nil)
	cfg := swagger.Config{
		Title: "Order Service API",
		Path:  "swagger", // UI root
		// BasePath: "/",       // Mount path
		FilePath: "./docs/swagger.json",
	}

	// Generate handler (this is where your panic happened)
	swaggerHandler := swagger.New(cfg)
	app.Get("/swagger/*", swaggerHandler)

	routes.InventorySetupRouter(app, inventoryHandler)

	// go inventoryHandler.CheckStock()
	// go inventoryHandler.ReverseStock()
	app.Listen(":8081")
}
