package bootstrap

import (
	"os"

	adapters_eventbus "github.com/Prompiriya084/go-mq/InventoryService/Internal/Adapters/Eventbus"
	adapters_handlers "github.com/Prompiriya084/go-mq/InventoryService/Internal/Adapters/Handlers"
	adapters_repositories "github.com/Prompiriya084/go-mq/InventoryService/Internal/Adapters/Repositories"
	services "github.com/Prompiriya084/go-mq/InventoryService/Internal/Core/Services"
	utilities_validator "github.com/Prompiriya084/go-mq/InventoryService/Internal/Utilities/Validator"
	routes "github.com/Prompiriya084/go-mq/InventoryService/Web/Routes"
	"github.com/gofiber/contrib/swagger"
	"github.com/gofiber/fiber/v2"
	"github.com/rabbitmq/amqp091-go"
	"gorm.io/gorm"
)

func StartInventoryAPI(db *gorm.DB, ch *amqp091.Channel) error {
	app := fiber.New()
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
	eventPublisher := adapters_eventbus.NewRabbitPublisher(ch, os.Getenv("MQ_Inventory_Core_Exchange"))

	validator := utilities_validator.NewValidator()
	service := services.NewInventoryAPIService(repo, eventPublisher)
	handler := adapters_handlers.NewInventoryAPIHandler(service, validator)

	routes.InventorySetupRouter(app, handler)

	app.Listen(":8081")
	return nil
}
