package bootstrap

import (
	"os"

	adapters_eventbus "github.com/Prompiriya084/go-mq/OrderService/Internal/Adapters/Eventbus"
	adapters_handlers "github.com/Prompiriya084/go-mq/OrderService/Internal/Adapters/Handlers"
	adapters_repositories "github.com/Prompiriya084/go-mq/OrderService/Internal/Adapters/Repositories"
	services "github.com/Prompiriya084/go-mq/OrderService/Internal/Core/Services"
	utilities_validator "github.com/Prompiriya084/go-mq/OrderService/Internal/Core/Utilities/Validator"
	routes "github.com/Prompiriya084/go-mq/OrderService/Web/Routes"
	"github.com/gofiber/contrib/swagger"
	"github.com/gofiber/fiber/v2"
	"github.com/rabbitmq/amqp091-go"
	"gorm.io/gorm"
)

func StartOrderAPI(db *gorm.DB, ch *amqp091.Channel) {
	app := fiber.New()

	// Serve swagger.json
	app.Static("/docs", "./docs")
	// Correct swagger config (must NOT be nil)
	cfg := swagger.Config{
		Title:    "Order Service API",
		Path:     "swagger", // UI root
		FilePath: "./docs/swagger.json",
	}

	// Generate handler (this is where your panic happened)
	swaggerHandler := swagger.New(cfg)
	app.Get("/swagger/*", swaggerHandler)

	repo := adapters_repositories.NewOrderRepository(db)
	eventbus := adapters_eventbus.NewRabbitPublisher(ch, os.Getenv("MQ_Order_Core_Exchange"))
	service := services.NewOrderAPIService(repo, eventbus)

	validator := utilities_validator.NewValidator()

	handler := adapters_handlers.NewOrderHandler(service, validator)
	routes.OrderSetupRouter(app, handler)

	app.Listen(":8080")
}
