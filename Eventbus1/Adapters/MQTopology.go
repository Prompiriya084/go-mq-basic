package adapters_eventbus

import (
	domain_eventbus "github.com/Prompiriya084/go-mq/Eventbus/Domain"
	"github.com/rabbitmq/amqp091-go"
)

type RabbitTopology struct {
	ch *amqp091.Channel
}

func NewRabbitTopology(ch *amqp091.Channel) *RabbitTopology {
	return &RabbitTopology{ch: ch}
}

func (t *RabbitTopology) DeclareQueue(cfg domain_eventbus.QueueConfig) error {

	// DLQ
	if _, err := t.ch.QueueDeclare(
		cfg.DLQName,
		true,
		false,
		false,
		false,
		nil,
	); err != nil {
		return err
	}

	// Retry Queue
	if _, err := t.ch.QueueDeclare(
		cfg.RetryName,
		true,
		false,
		false,
		false,
		amqp091.Table{
			"x-message-ttl":             cfg.RetryTTL,
			"x-dead-letter-exchange":    "",
			"x-dead-letter-routing-key": cfg.Name,
		},
	); err != nil {
		return err
	}

	// Main Queue
	if _, err := t.ch.QueueDeclare(
		cfg.Name,
		true,
		false,
		false,
		false,
		amqp091.Table{
			"x-dead-letter-exchange":    "",
			"x-dead-letter-routing-key": cfg.RetryName,
		},
	); err != nil {
		return err
	}

	return nil
}
