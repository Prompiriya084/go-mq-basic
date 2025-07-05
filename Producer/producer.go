package main

import (
	"os"

	adapters_repositories "github.com/Prompiriya084/go-mq/Customer/Adapters/Repositories"
	database "github.com/Prompiriya084/go-mq/Infrastructure/Database"
	adapters_producers "github.com/Prompiriya084/go-mq/Producer/Internal/Adapters/MQ"
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
	repo := adapters_repositories.NewOrderRepository(db)
	routes.OrderSetupRouter(app, repo, mqProducer)
	app.Listen(":8080")
}
