package main

import (
	adapters_handlers "github.com/Prompiriya084/go-mq/InventoryService/Internal/Adapters/Handlers"
	database "github.com/Prompiriya084/go-mq/InventoryService/Internal/Infrastructure/Database"
	utilities_validator "github.com/Prompiriya084/go-mq/InventoryService/Internal/Utilities/Validator"
	routes "github.com/Prompiriya084/go-mq/InventoryService/Web/Routes"

	"github.com/gofiber/contrib/swagger"
	"github.com/gofiber/fiber/v2"

	adapters_repositories "github.com/Prompiriya084/go-mq/InventoryService/Internal/Adapters/Repositories"
	services "github.com/Prompiriya084/go-mq/InventoryService/Internal/Core/Services"
)

// @title Inventory Service API
// @version 1.0
// @description API Example
// @BasePath /
// @schemes http
func main() {
	app := fiber.New()
	db := database.InitDb()

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

	repo := adapters_repositories.NewInventoryRepository(db)
	// mqEventbus := eventbus.NewMQEventbus[models.Order](os.Getenv("RABBITMQ_URL"))

	inventoryAPIService := services.NewInventoryAPIService(repo, mqEventbus)

	validator := utilities_validator.NewValidator()
	inventoryAPIHandler := adapters_handlers.NewInventoryAPIHandler(inventoryAPIService, validator)

	routes.InventorySetupRouter(app, inventoryAPIHandler)

	// go inventoryHandler.CheckStock()
	// go inventoryHandler.ReverseStock()

	app.Listen(":8081")
}
