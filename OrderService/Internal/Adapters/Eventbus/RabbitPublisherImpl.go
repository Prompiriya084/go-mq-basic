package adapters_eventbus

import (
	"encoding/json"

	ports_eventbus "github.com/Prompiriya084/go-mq/InventoryService/Internal/Core/Ports/Eventbus"
	"github.com/rabbitmq/amqp091-go"
)

type eventbusPublisherImpl struct {
	ch       *amqp091.Channel
	exchange string
}

func NewRabbitPublisher(ch *amqp091.Channel, exchange string) ports_eventbus.EventBusPublisher {
	return &eventbusPublisherImpl{
		ch:       ch,
		exchange: exchange,
	}
}

func (p *eventbusPublisherImpl) Publish(eventName string, payload any) error {
	body, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	return p.ch.Publish(
		p.exchange,
		eventName,
		false,
		false,
		amqp091.Publishing{
			ContentType: "application/json",
			Body:        body,
		},
	)
}
