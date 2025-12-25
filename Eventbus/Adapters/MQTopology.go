package adapters_eventbus

import "github.com/rabbitmq/amqp091-go"

const (
	OrderExchange = "order.events"
)

func DeclareTopology(ch *amqp091.Channel) error {
	// Exchange
	if err := ch.ExchangeDeclare(
		OrderExchange,
		"topic",
		true,
		false,
		false,
		false,
		nil,
	); err != nil {
		return err
	}

	// Main queue
	_, err := ch.QueueDeclare(
		"inventory.order.created",
		true,
		false,
		false,
		false,
		amqp091.Table{
			"x-dead-letter-exchange":    OrderExchange,
			"x-dead-letter-routing-key": "order.created.retry",
		},
	)
	if err != nil {
		return err
	}

	// Retry queue
	_, err = ch.QueueDeclare(
		"inventory.order.created.retry",
		true,
		false,
		false,
		false,
		amqp091.Table{
			"x-message-ttl":             int32(10000),
			"x-dead-letter-exchange":    OrderExchange,
			"x-dead-letter-routing-key": "order.created",
		},
	)

	// DLQ
	_, err = ch.QueueDeclare(
		"inventory.order.created.dlq",
		true,
		false,
		false,
		false,
		nil,
	)

	// Bindings
	ch.QueueBind("inventory.order.created", "order.created", OrderExchange, false, nil)
	ch.QueueBind("inventory.order.created.retry", "order.created.retry", OrderExchange, false, nil)
	ch.QueueBind("inventory.order.created.dlq", "order.created.dlq", OrderExchange, false, nil)

	return err
}
