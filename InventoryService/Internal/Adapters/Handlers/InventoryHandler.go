package adapters_handlers

import (
	"fmt"
	"log"

	eventbus "github.com/Prompiriya084/go-mq/Eventbus"
	services "github.com/Prompiriya084/go-mq/InventoryService/Internal/Core/Services"
	utilities_validator "github.com/Prompiriya084/go-mq/InventoryService/Internal/Utilities/Validator"
	models "github.com/Prompiriya084/go-mq/InventoryService/Models"
	"github.com/gofiber/fiber/v2"
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

// Inventory godoc
// @Summary Create stocks
// @Description Create a new stock
// @Tags Inventory
// @Accept json
// @Produce json
// @Param request body models.Inventory true "Inventory info"
// @Success 202 {object} dto.MessageResponse
// @Failure 400 {string} string "Bad Request"
// @Failure 500 {string} string "Internal Server Error"
// @Router /api/inventory [post]
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
		"product_id": inventory.ProductID,
		"message":    "Created stock successful.",
	})
}

// Inventory godoc
// @Summary Get all stocks
// @Description Get all stocks
// @Tags Inventory
// @Accept json
// @Produce json
// @Success 200 {array} models.Inventory
// @Router /api/inventory [get]
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

// Inventory godoc
// @Summary Get stocks By Id
// @Description Get stocks by product Id
// @Tags Inventory
// @Accept json
// @Produce json
// @Param   id   path     string  true  "The ID of the resource"
// @Success 200 {array} models.Inventory
// @Router /api/inventory/{id} [get]
func (h *InventoryHandler) Get(c *fiber.Ctx) error {
	// orderID, err := uuid.Parse(c.Params("id"))
	productId := c.Params("id")
	if productId == "" {
		return c.SendStatus(fiber.StatusBadRequest)
	}
	orders, err := h.service.Get(&models.Inventory{
		ProductID: productId,
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

// UpdateInventory godoc
// @Summary Update Inventory
// @Description Update stocks
// @Tags Inventory
// @Accept json
// @Produce json
// @Param request body models.Inventory true "Inventory info"
// @Param   id   path     string  true  "The ID of the resource"
// @Success 200 {object} dto.MessageResponse
// @Failure 400 {string} string "Bad Request"
// @Failure 500 {string} string "Internal Server Error"
// @Router /api/inventory/{id} [put]
func (h *InventoryHandler) Update(c *fiber.Ctx) error {
	productId := c.Params("id")
	if productId == "" {
		return c.SendStatus(fiber.StatusBadRequest)
	}

	var updatedItem models.Inventory
	if err := c.BodyParser(&updatedItem); err != nil {
		return c.SendStatus(fiber.StatusBadRequest)
	}
	updatedItem.ProductID = productId

	if err := h.validator.ValidateStruct(updatedItem); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}
	if err := h.service.Update(&updatedItem); err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}

	return c.JSON(fiber.Map{
		"product_id": updatedItem.ProductID,
		"message":    "Updated stock successful.",
	})
}

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
