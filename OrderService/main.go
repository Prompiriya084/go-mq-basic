package main

import (
	"os"

	eventbus "github.com/Prompiriya084/go-mq/EventBus"
	database "github.com/Prompiriya084/go-mq/Infrastructure/Database"
	models "github.com/Prompiriya084/go-mq/Models"
	adapters_handlers "github.com/Prompiriya084/go-mq/OrderService/Internal/Adapters/Handlers"
	adapters_repositories "github.com/Prompiriya084/go-mq/OrderService/Internal/Adapters/Repositories"
	services "github.com/Prompiriya084/go-mq/OrderService/Internal/Core/Services"
	utilities_validator "github.com/Prompiriya084/go-mq/OrderService/Internal/Core/Utilities/Validator"
	routes "github.com/Prompiriya084/go-mq/OrderService/Web/Routes"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/adaptor"
	httpSwagger "github.com/swaggo/http-swagger"
)

// @title Fiber v3 Swagger Example
// @version 1.0
// @description Swagger with Fiber v3
// @host localhost:8080
// @BasePath /api
func main() {
	app := fiber.New()

	db := database.InitDb()

	app.Get("/swagger/*", adaptor.HTTPHandlerFunc(httpSwagger.Handler(
		httpSwagger.URL("http://localhost:8080/swagger/doc.json"),
	)))
	mqOrderEventbus := eventbus.NewMQEventbus[models.Order](os.Getenv("RABBITMQ_URL"))
	orderRepo := adapters_repositories.NewOrderRepository(db)
	orderService := services.NewOrderService(orderRepo, mqOrderEventbus)

	validator := utilities_validator.NewValidator()

	orderHandler := adapters_handlers.NewOrderHandler(orderService, validator, mqOrderEventbus)
	routes.OrderSetupRouter(app, orderHandler)

	//For message queue
	go orderHandler.InventoryConfirmed()
	app.Listen(":8080")
}
