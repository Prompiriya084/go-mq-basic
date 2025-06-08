package adapters_handlers

import (
	"log"

	ports_mq "github.com/Prompiriya084/go-mq/Customer/Core/Ports/MQ"
	services "github.com/Prompiriya084/go-mq/Customer/Core/Services"
	models "github.com/Prompiriya084/go-mq/Models"
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
	err := h.mq.ReceiveMessage("order", func(order models.Order) error {

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
