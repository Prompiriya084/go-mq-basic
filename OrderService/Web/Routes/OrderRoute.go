package routes

import (
	adapters_handlers "github.com/Prompiriya084/go-mq/OrderService/Internal/Adapters/Handlers"

	"github.com/gofiber/fiber/v2"
)

func OrderSetupRouter(app *fiber.App, handler *adapters_handlers.OrderHandler) {
	orderApp := app.Group("/api/orders")
	orderApp.Post("", handler.Create)
	orderApp.Get("", handler.GetAll)
	orderApp.Get("/:id", handler.Get)
	orderApp.Put("/:id", handler.Update)
	orderApp.Delete("/:id", handler.Delete)
}
