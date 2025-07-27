package main

import (
	"os"

	database "github.com/Prompiriya084/go-mq/Infrastructure/Database"
	adapters_producers "github.com/Prompiriya084/go-mq/Producer/Internal/Adapters/MQ"
	adapters_repositories "github.com/Prompiriya084/go-mq/Producer/Internal/Adapters/Repositories"
	services "github.com/Prompiriya084/go-mq/Producer/Internal/Core/Services"
	routes "github.com/Prompiriya084/go-mq/Producer/Web/Routes"
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
	mqProducer := adapters_producers.NewMQProducer(os.Getenv("RABBITMQ_URL"))
	orderRepo := adapters_repositories.NewOrderRepository(db)
	orderService := services.NewOrderService(orderRepo, mqProducer)

	routes.OrderSetupRouter(app, orderService)
	app.Listen(":8080")
}
