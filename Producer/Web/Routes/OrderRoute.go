package routes

import (
	ports_repositories "github.com/Prompiriya084/go-mq/Customer/Core/Ports/Repositories"
	adapters_handlers "github.com/Prompiriya084/go-mq/Producer/Internal/Adapters/Handlers"

	ports_mq "github.com/Prompiriya084/go-mq/Producer/Internal/Core/Ports/MQ"
	services "github.com/Prompiriya084/go-mq/Producer/Internal/Core/Services"
	utilities_validator "github.com/Prompiriya084/go-mq/Producer/Internal/Core/Utilities/Validator"
	"github.com/gofiber/fiber/v3"
)

func OrderSetupRouter(app *fiber.App, repo ports_repositories.OrderRepository, mqProducer ports_mq.MQProducer) {

	service := services.NewOrderService(repo, mqProducer)
	validator := utilities_validator.NewValidator()
	handler := adapters_handlers.NewOrderHandler(service, validator)

	orderApp := app.Group("/api/orders")
	orderApp.Post("", handler.Create)
	orderApp.Get("", handler.GetAll)
}
