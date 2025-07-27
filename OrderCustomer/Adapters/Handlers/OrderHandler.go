package adapters_handlers

import (
	"log"

	models "github.com/Prompiriya084/go-mq/Models"
	ports_mq "github.com/Prompiriya084/go-mq/OrderCustomer/Core/Ports/MQ"
	services "github.com/Prompiriya084/go-mq/OrderCustomer/Core/Services"
)

type OrderHandler struct {
	service services.OrderService
	mq      ports_mq.MQCustomer[models.Order]
}

func NewOrderHandler(service services.OrderService, mq ports_mq.MQCustomer[models.Order]) *OrderHandler {
	return &OrderHandler{
		service: service,
		mq:      mq,
	}
}
func (h *OrderHandler) Create() {
	err := h.mq.ReceiveMessage("order_create", func(order models.Order) error {

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
	err := h.mq.ReceiveMessage("order_update", func(order models.Order) error {

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
func (h *OrderHandler) Delete() {
	err := h.mq.ReceiveMessage("order_delete", func(order models.Order) error {

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
