package routes

import (
	adapters_handlers "github.com/Prompiriya084/go-mq/InventoryService/Internal/Adapters/Handlers"
	"github.com/gofiber/fiber/v2"
)

func InventorySetupRouter(app *fiber.App, handler *adapters_handlers.InventoryHandler) {
	inventyApp := app.Group("/api/inventory")
	inventyApp.Post("", handler.Create)
	inventyApp.Get("", handler.GetAll)
	inventyApp.Get("/:id", handler.Get)
}
