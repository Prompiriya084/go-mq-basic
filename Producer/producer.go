package main

import (
	database "github.com/Prompiriya084/go-mq/Infrastructure/Database"
	routes "github.com/Prompiriya084/go-mq/Producer/Web/Routes"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/adaptor"
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

	routes.OrderSetupRouter(db, app)
	app.Listen(":8080")
}
