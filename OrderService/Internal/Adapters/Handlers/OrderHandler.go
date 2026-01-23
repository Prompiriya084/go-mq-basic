package adapters_handlers

import (
	"fmt"

	services "github.com/Prompiriya084/go-mq/OrderService/Internal/Core/Services"
	utilities_validator "github.com/Prompiriya084/go-mq/OrderService/Internal/Core/Utilities/Validator"
	models "github.com/Prompiriya084/go-mq/OrderService/Models"
	"github.com/google/uuid"

	"github.com/gofiber/fiber/v2"
)

type OrderHandler struct {
	service   services.OrderAPIService
	validator utilities_validator.Validator
}

func NewOrderHandler(service services.OrderAPIService,
	validator utilities_validator.Validator) *OrderHandler {
	return &OrderHandler{
		service:   service,
		validator: validator,
	}
}

// Order godoc
// @Summary Get all orders
// @Description Get all orders
// @Tags Orders
// @Accept json
// @Produce json
// @Success 200 {array} models.Order
// @Router /api/orders [get]
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

// Order godoc
// @Summary Get Order By Id
// @Description Get order by Id
// @Tags Orders
// @Accept json
// @Produce json
// @Param   id   path     string  true  "The ID of the resource"
// @Success 200 {array} models.Order
// @Router /api/orders/{id} [get]
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

// Order godoc
// @Summary Create order
// @Description Create a new order
// @Tags Orders
// @Accept json
// @Produce json
// @Param request body models.Order true "Order info"
// @Success 202 {object} dto.MessageResponse
// @Failure 400 {string} string "Bad Request"
// @Failure 500 {string} string "Internal Server Error"
// @Router /api/orders [post]
func (h *OrderHandler) Create(c *fiber.Ctx) error {
	var order models.Order
	if err := c.BodyParser(&order); err != nil {
		return c.SendStatus(fiber.StatusBadRequest)
	}

	if err := h.validator.ValidateStruct(order); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}
	orderId, err := h.service.Create(&order)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}

	return c.Status(202).JSON(fiber.Map{
		"order_id": orderId,
		"status":   "PENDING",
		"message":  "Order received. Waiting for inventory confirmation.",
	})
}

// Order godoc
// @Summary Update order
// @Description Update an order
// @Tags Orders
// @Accept json
// @Produce json
// @Param request body models.Order true "Order info"
// @Param   id   path     string  true  "The ID of the resource"
// @Success 201 {object} dto.MessageResponse
// @Failure 400 {string} string "Bad Request"
// @Failure 500 {string} string "Internal Server Error"
// @Router /api/orders/{id} [put]
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
		"order_id": orderId,
		"message":  "Update order successful.",
	})
}

// Order godoc
// @Summary Create order
// @Description Create a new order
// @Tags Orders
// @Accept json
// @Produce json
// @Param   id   path     string  true  "The ID of the resource"
// @Success 202 {object} dto.MessageResponse
// @Failure 400 {string} string "Bad Request"
// @Failure 500 {string} string "Internal Server Error"
// @Router /api/orders [delete]
func (h *OrderHandler) Delete(c *fiber.Ctx) error {
	orderId := c.Params("id")
	if orderId == "" {
		return c.SendStatus(fiber.StatusBadRequest)
	}

	if err := h.service.Delete(orderId); err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}

	return c.JSON(fiber.Map{
		"order_id": orderId,
		"message":  "Deleted order successful.",
	})
}

// func (h *OrderHandler) InventoryConfirmed() {
// 	err := h.bus.Subscribe("inventory.checked", func(param models.Order) error {

// 		log.Printf("✅ Processed Order: %v", param)
// 		param.Status = "COMPLETED"
// 		if err := h.service.Update(&param); err != nil {
// 			log.Println(err.Error())
// 			return err
// 		}

// 		return nil //return null when wanting to acknowledge when process complete
// 	})
// 	if err != nil {
// 		panic(err)
// 	}
// }
// func (h *OrderHandler) InventoryFailed() {
// 	err := h.bus.Subscribe("inventory.failed", func(evt models.Order) error {

// 		log.Printf("✅ Processed Order: %v", evt)
// 		evt.Status = "FAILED"
// 		if err := h.service.Update(&evt); err != nil {
// 			log.Println(err.Error())
// 			return err
// 		}

// 		return nil //return null when wanting to acknowledge when process complete
// 	})
// 	if err != nil {
// 		panic(err)
// 	}
// }
