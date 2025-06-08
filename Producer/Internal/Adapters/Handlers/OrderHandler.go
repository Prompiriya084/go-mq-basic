package adapters_handlers

import (
	models "github.com/Prompiriya084/go-mq/Models"
	services "github.com/Prompiriya084/go-mq/Producer/Internal/Core/Services"
	utilities "github.com/Prompiriya084/go-mq/Producer/Internal/Core/Utilities"
	"github.com/gofiber/fiber/v3"
)

type OrderHandler struct {
	service   services.OrderService
	validator utilities.Validator
}

func NewOrderHandler(service services.OrderService, validator utilities.Validator) *OrderHandler {
	return &OrderHandler{
		service:   service,
		validator: validator,
	}
}

func (h *OrderHandler) GetAll(c fiber.Ctx) error {
	orders, err := h.service.GetAll(nil, nil)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}

	if len(orders) == 0 {
		return c.Status(fiber.StatusNotFound).SendString("Data not found.")
	}
	return c.JSON(fiber.Map{
		"data": orders,
	})
}

func (h *OrderHandler) Create(c fiber.Ctx) error {
	var order models.Order
	if err := c.Bind().Body(&order); err != nil {
		return c.SendStatus(fiber.StatusBadRequest)
	}

	if err := h.validator.ValidateStruct(order); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}
	if err := h.service.Create(&order); err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}

	return c.JSON(fiber.Map{
		"message": "create order successful.",
	})
}
