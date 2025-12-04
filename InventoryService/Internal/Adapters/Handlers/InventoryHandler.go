package adapters_handlers

import (
	"fmt"
	"log"

	eventbus "github.com/Prompiriya084/go-mq/Eventbus"
	services "github.com/Prompiriya084/go-mq/InventoryService/Internal/Core/Services"
	utilities_validator "github.com/Prompiriya084/go-mq/InventoryService/Internal/Utilities/Validator"
	models "github.com/Prompiriya084/go-mq/InventoryService/Models"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type InventoryHandler struct {
	service   services.InventoryService
	bus       eventbus.EventBus[models.Order]
	validator utilities_validator.Validator
}

func NewInventoryHandler(service services.InventoryService, bus eventbus.EventBus[models.Order], validator utilities_validator.Validator) *InventoryHandler {
	return &InventoryHandler{
		service: service,
		bus:     bus,
	}
}
func (h *InventoryHandler) CheckStock() {
	err := h.bus.Subscribe("order.created", func(param models.Order) error {

		log.Printf("✅ Processed Order: %v", param)
		if err := h.service.CheckStock(&param); err != nil {
			log.Printf("❌ Check stock failed: Order: %v, Exception: %v", param, err)
			return err
		}

		return nil //return null when wanting to acknowledge when process complete
	})
	if err != nil {
		panic(err)
	}
}
func (h *InventoryHandler) ReverseStock() {
	err := h.bus.Subscribe("payment.failed", func(param models.Order) error {

		log.Printf("✅ Processed Order: %v", param)
		if err := h.service.ReverseStock(&param); err != nil {
			log.Printf("❌ Reverse stock failed: Order: %v, Exception: %v", param, err)
			return err
		}

		return nil //return null when wanting to acknowledge when process complete
	})
	if err != nil {
		panic(err)
	}
}
func (h *InventoryHandler) Create(c *fiber.Ctx) error {
	var inventory models.Inventory

	if err := c.BodyParser(&inventory); err != nil {
		return c.SendStatus(fiber.StatusBadRequest)
	}

	if err := h.validator.ValidateStruct(inventory); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}
	if err := h.service.Create(&inventory); err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}

	return c.JSON(fiber.Map{
		"message": "Create stock successful.",
	})
}
func (h *InventoryHandler) GetAll(c *fiber.Ctx) error {
	orders, err := h.service.GetAll(nil, nil)
	fmt.Println("Stock : ", orders)
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

func (h *InventoryHandler) Get(c *fiber.Ctx) error {
	orderID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.SendStatus(fiber.StatusBadRequest)
	}

	orders, err := h.service.Get(&models.Inventory{
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

// func (h *InventoryHandler) Update() {
// 	err := h.bus.Subscribe("order.update", func(order models.Inventory) error {

// 		log.Printf("✅ Processed Order: ID=%s", order.ID)
// 		if err := h.service.Update(&order); err != nil {
// 			log.Printf("❌ DB update failed: %v", err)
// 			return err
// 		}

// 		return nil //return null when wanting to acknowledge when process complete
// 	})
// 	if err != nil {
// 		panic(err)
// 	}
// }
// func (h *InventoryHandler) Cancel() {
// 	err := h.bus.Subscribe("order.cancel", func(order models.Inventory) error {

// 		log.Printf("✅ Processed Order: ID=%s", order.ID)
// 		if err := h.service.Delete(&order); err != nil {
// 			log.Printf("❌ DB delete failed: %v", err)
// 			return err
// 		}

// 		return nil //return null when wanting to acknowledge when process complete
// 	})
// 	if err != nil {
// 		panic(err)
// 	}
// }
