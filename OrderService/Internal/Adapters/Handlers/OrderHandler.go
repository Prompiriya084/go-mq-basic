package adapters_handlers

import (
	"fmt"
	"log"

	eventbus "github.com/Prompiriya084/go-mq/EventBus"
	models "github.com/Prompiriya084/go-mq/Models"
	services "github.com/Prompiriya084/go-mq/OrderService/Internal/Core/Services"
	utilities_validator "github.com/Prompiriya084/go-mq/OrderService/Internal/Core/Utilities/Validator"
	"github.com/google/uuid"

	"github.com/gofiber/fiber/v2"
)

type OrderHandler struct {
	service   services.OrderService
	validator utilities_validator.Validator
	bus       eventbus.EventBus[models.Order]
}

func NewOrderHandler(service services.OrderService,
	validator utilities_validator.Validator,
	bus eventbus.EventBus[models.Order]) *OrderHandler {
	return &OrderHandler{
		service:   service,
		validator: validator,
		bus:       bus,
	}
}

func (h *OrderHandler) GetAll(c *fiber.Ctx) error {
	orders, err := h.service.GetAll(nil, nil)
	fmt.Println("Orders : ", orders)
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

func (h *OrderHandler) Get(c *fiber.Ctx) error {
	orderID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.SendStatus(fiber.StatusBadRequest)
	}

	orders, err := h.service.Get(&models.Order{
		ID: orderID,
	}, nil)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}

	if orders == nil {
		return c.Status(fiber.StatusNotFound).SendString("Data not found.")
	}

	return c.JSON(fiber.Map{
		"data": orders,
	})
}

func (h *OrderHandler) Create(c *fiber.Ctx) error {
	var order models.Order
	if err := c.BodyParser(&order); err != nil {
		return c.SendStatus(fiber.StatusBadRequest)
	}

	if err := h.validator.ValidateStruct(order); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}
	if err := h.service.Create(&order); err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}

	return c.JSON(fiber.Map{
		"message": "Create order successful.",
	})
}

func (h *OrderHandler) Update(c *fiber.Ctx) error {
	orderId, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}

	var order models.Order
	if err := c.BodyParser(&order); err != nil {
		return c.SendStatus(fiber.StatusBadRequest)
	}
	order.ID = orderId

	if err := h.validator.ValidateStruct(order); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}
	if err := h.service.Update(&order); err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}

	return c.JSON(fiber.Map{
		"message": "Update order successful.",
	})
}

func (h *OrderHandler) Delete(c *fiber.Ctx) error {
	orderID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.SendStatus(fiber.StatusBadRequest)
	}

	if err := h.service.Delete(&models.Order{
		ID: orderID,
	}); err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}

	return c.JSON(fiber.Map{
		"message": "Delete order successful.",
	})
}

func (h *OrderHandler) InventoryConfirmed() {
	err := h.bus.Subscribe("inventory.checked", func(param models.Order) error {

		log.Printf("✅ Processed Order: %v", param)
		// if err := h.service.ReverseStock(&param); err != nil {
		// 	log.Printf("❌ Reverse stock failed: Order: %v, Exception: %v", param, err)
		// 	return err
		// }

		return nil //return null when wanting to acknowledge when process complete
	})
	if err != nil {
		panic(err)
	}
}
