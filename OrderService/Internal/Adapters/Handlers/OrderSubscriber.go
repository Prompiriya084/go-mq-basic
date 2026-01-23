package adapters_handlers

import (
	"encoding/json"

	services "github.com/Prompiriya084/go-mq/OrderService/Internal/Core/Services"
	models_eventbus "github.com/Prompiriya084/go-mq/OrderService/Models/Eventbus"
	"github.com/rabbitmq/amqp091-go"
)

type OrderSubscriber struct {
	ch      *amqp091.Channel
	service services.OrderEventService
}

func NewOrderSubscriber(ch *amqp091.Channel, service services.OrderEventService) *OrderSubscriber {
	return &OrderSubscriber{
		ch:      ch,
		service: service,
	}
}

// Receive message from "inventory.queue"
func (s *OrderSubscriber) Start() error {
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

func (s *OrderSubscriber) handleMessage(msg amqp091.Delivery) {
	var event models_eventbus.OrderCreated
	if err := json.Unmarshal(msg.Body, &event); err != nil {
		msg.Nack(false, false) // → DLQ
		return
	}

	if err := s.service.OrderCreated(&event); err != nil {
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
