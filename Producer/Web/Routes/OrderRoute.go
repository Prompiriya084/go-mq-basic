package routes

import (
	adapters_handlers "github.com/Prompiriya084/go-mq/Producer/Internal/Adapters/Handlers"

	services "github.com/Prompiriya084/go-mq/Producer/Internal/Core/Services"
	utilities_validator "github.com/Prompiriya084/go-mq/Producer/Internal/Core/Utilities/Validator"
	"github.com/gofiber/fiber/v2"
)

func OrderSetupRouter(app *fiber.App, service services.OrderService) {

	validator := utilities_validator.NewValidator()
	handler := adapters_handlers.NewOrderHandler(service, validator)

	orderApp := app.Group("/api/orders")
	orderApp.Post("", handler.Create)
	orderApp.Get("", handler.GetAll)
	orderApp.Get("/:id", handler.Get)
	orderApp.Put("/:id", handler.Update)
	orderApp.Delete("/:id", handler.Delete)
}
