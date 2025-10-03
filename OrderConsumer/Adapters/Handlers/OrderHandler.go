package adapters_handlers

import (
	"log"

	eventbus "github.com/Prompiriya084/go-mq/Eventbus"
	models "github.com/Prompiriya084/go-mq/Models"
	services "github.com/Prompiriya084/go-mq/OrderConsumer/Core/Services"
)

type OrderHandler struct {
	service services.OrderService
	bus     eventbus.EventBus[models.Order]
}

func NewOrderHandler(service services.OrderService, bus eventbus.EventBus[models.Order]) *OrderHandler {
	return &OrderHandler{
		service: service,
		bus:     bus,
	}
}
func (h *OrderHandler) Create() {
	err := h.bus.Subscribe("order.create", func(order models.Order) error {

		log.Printf("✅ Processed Order: ID=%s", order.ID)
		if err := h.service.Create(&order); err != nil {
			log.Printf("❌ DB insert failed: %v", err)
			return err
		}

		return nil //return null when wanting to acknowledge when process complete
	})
	if err != nil {
		panic(err)
	}
}
func (h *OrderHandler) Update() {
	err := h.bus.Subscribe("order.update", func(order models.Order) error {

		log.Printf("✅ Processed Order: ID=%s", order.ID)
		if err := h.service.Update(&order); err != nil {
			log.Printf("❌ DB update failed: %v", err)
			return err
		}

		return nil //return null when wanting to acknowledge when process complete
	})
	if err != nil {
		panic(err)
	}
}
func (h *OrderHandler) Cancel() {
	err := h.bus.Subscribe("order.cancel", func(order models.Order) error {

		log.Printf("✅ Processed Order: ID=%s", order.ID)
		if err := h.service.Delete(&order); err != nil {
			log.Printf("❌ DB delete failed: %v", err)
			return err
		}

		return nil //return null when wanting to acknowledge when process complete
	})
	if err != nil {
		panic(err)
	}
}
