package adapters_handlers

import (
	"log"

	eventbus "github.com/Prompiriya084/go-mq/Eventbus"
	services "github.com/Prompiriya084/go-mq/InventoryService/Core/Services"
	models "github.com/Prompiriya084/go-mq/Models"
)

type InventoryHandler struct {
	service services.InventoryService
	bus     eventbus.EventBus[models.Order]
}

func NewInventoryHandler(service services.InventoryService, bus eventbus.EventBus[models.Order]) *InventoryHandler {
	return &InventoryHandler{
		service: service,
		bus:     bus,
	}
}
func (h *InventoryHandler) CheckStock() {
	err := h.bus.Subscribe("inventory.check.requested", func(param models.Order) error {

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
