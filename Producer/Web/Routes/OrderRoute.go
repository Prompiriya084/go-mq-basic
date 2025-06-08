package routes

import (
	"os"

	adapters_handlers "github.com/Prompiriya084/go-mq/Producer/Internal/Adapters/Handlers"
	adapters_producers "github.com/Prompiriya084/go-mq/Producer/Internal/Adapters/MQ"
	adapters_repositories "github.com/Prompiriya084/go-mq/Producer/Internal/Adapters/Repositories"
	adapters_utilities "github.com/Prompiriya084/go-mq/Producer/Internal/Adapters/Utilities"
	services "github.com/Prompiriya084/go-mq/Producer/Internal/Core/Services"
	"github.com/gofiber/fiber/v3"
	"gorm.io/gorm"
)

func OrderSetupRouter(db *gorm.DB, app *fiber.App) {

	mqProducer := adapters_producers.NewMQProducer(os.Getenv("RABBITMQ_URL"))
	repo := adapters_repositories.NewOrderRepository(db)
	service := services.NewOrderService(repo, mqProducer)
	validator := adapters_utilities.NewValidator()
	handler := adapters_handlers.NewOrderHandler(service, validator)

	orderApp := app.Group("/api/orders")
	orderApp.Post("", handler.Create)
	orderApp.Get("", handler.GetAll)
}
