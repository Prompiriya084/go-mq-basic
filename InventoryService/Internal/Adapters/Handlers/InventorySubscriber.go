package adapters_handlers

import (
	"encoding/json"

	services "github.com/Prompiriya084/go-mq/InventoryService/Internal/Core/Services"
	events "github.com/Prompiriya084/go-mq/InventoryService/Models/Events"
	"github.com/rabbitmq/amqp091-go"
)

type InventorySubscriber struct {
	ch       *amqp091.Channel
	service  services.InventoryEventService
	maxRetry int
}

func NewInventorySubscriber(ch *amqp091.Channel, service services.InventoryEventService, maxRetry int) *InventorySubscriber {
	return &InventorySubscriber{
		ch:       ch,
		service:  service,
		maxRetry: maxRetry,
	}
}

// Receive message from "inventory.queue"
func (s *InventorySubscriber) Start() error {
	msgs, _ := s.ch.Consume(
		"inventory.queue",
		"",
		false, // manual ack
		false,
		false,
		false,
		nil,
	)

	go func() {
		for msg := range msgs {
			s.handleMessage(msg)
		}
	}()

	return nil
}

func (s *InventorySubscriber) handleMessage(msg amqp091.Delivery) {
	var event events.OrderCreated
	if err := json.Unmarshal(msg.Body, &event); err != nil {
		msg.Nack(false, false) // → DLQ
		return
	}

	if err := s.service.ReserveStock(&event); err != nil {
		// retry := s.getRetryCount(msg.Headers)
		// if retry >= s.maxRetry {
		// 	msg.Nack(false, false) // → DLQ
		// 	return
		// }

		// headers := msg.Headers
		// headers["x-retry"] = retry + 1

		// // ส่งไป retry exchange
		// s.ch.Publish(
		// 	"inventory.retry.exchange",
		// 	"inventory.retry",
		// 	false,
		// 	false,
		// 	amqp091.Publishing{
		// 		Headers: headers,
		// 		Body:    msg.Body,
		// 	},
		// )

		// msg.Ack(false)
		// return

		msg.Nack(false, false) // → retry
		return
	}

	msg.Ack(false)
}

func (s *InventorySubscriber) getRetryCount(headers amqp091.Table) int {
	if headers == nil {
		return 0
	}
	if v, ok := headers["x-retry"]; ok {
		return int(v.(int32))
	}
	return 0
}
