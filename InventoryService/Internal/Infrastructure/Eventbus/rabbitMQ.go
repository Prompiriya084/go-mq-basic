package rabbitMQ

import (
	"os"

	"github.com/rabbitmq/amqp091-go"
)

func MQConnectionInit(connStr string) (*amqp091.Channel, error) {
	// // Close previous connection/channel if they exist
	// if amqp091.channel != nil {
	// 	_ = c.channel.Close()
	// }
	// if c.conn != nil {
	// 	_ = c.conn.Close()
	// }
	// Dial new connection
	conn, err := amqp091.Dial(connStr)
	if err != nil {
		return nil, err
	}

	channel, err := conn.Channel()
	if err != nil {
		return nil, err
	}

	return channel, nil
}

func DeclareTopology(ch *amqp091.Channel) error {

	inventory_core_exchange := os.Getenv("MQ_Inventory_Core_Exchange")
	inventory_core_routing := os.Getenv("MQ_Inventory_Core_Routing")
	inventory_core_queue := os.Getenv("MQ_Inventory_Core_Queue")
	inventory_retry_exchange := os.Getenv("MQ_Inventory_Retry_Exchange")
	inventory_retry_routing := os.Getenv("MQ_Inventory_Retry_Routing")
	inventory_retry_queue := os.Getenv("MQ_Inventory_Retry_Queue")
	inventory_dlq_exchange := os.Getenv("MQ_Inventory_DLQ_Exchange")
	inventory_dlq_routing := os.Getenv("MQ_Inventory_DLQ_Routing")
	inventory_dlq_queue := os.Getenv("MQ_Inventory_DLQ_Queue")

	if err := ch.ExchangeDeclare(
		inventory_core_exchange,
		"topic",
		true,
		false,
		false,
		false,
		nil,
	); err != nil {
		return err
	}
	// retry + dlq exchange
	if err := ch.ExchangeDeclare(
		inventory_retry_exchange,
		"direct",
		true,
		false,
		false,
		false,
		nil,
	); err != nil {
		return err
	}

	if err := ch.ExchangeDeclare(
		inventory_dlq_exchange,
		"direct",
		true,
		false,
		false,
		false,
		nil,
	); err != nil {
		return err
	}

	// Main Queue
	_, err := ch.QueueDeclare(
		inventory_core_queue,
		true,
		false,
		false,
		false,
		amqp091.Table{
			"x-dead-letter-exchange":    inventory_retry_exchange,
			"x-dead-letter-routing-key": inventory_retry_routing,
		},
	)
	if err != nil {
		return nil
	}

	// retry queue
	_, err = ch.QueueDeclare(
		inventory_retry_queue,
		true,
		false,
		false,
		false,
		amqp091.Table{
			"x-message-ttl":             5000,
			"x-dead-letter-exchange":    inventory_dlq_exchange,
			"x-dead-letter-routing-key": inventory_dlq_routing,
		},
	)
	if err != nil {
		return err
	}
	// dlq
	_, err = ch.QueueDeclare(
		inventory_dlq_queue,
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}

	// bind core process
	if err := ch.QueueBind(
		inventory_core_queue,    //queuename
		inventory_core_routing,  //routing key
		inventory_core_exchange, //exchange
		false,
		nil,
	); err != nil {
		return err
	}

	// bind retry process
	if err := ch.QueueBind(inventory_retry_queue, //queuename
		inventory_retry_routing,  //routing key
		inventory_retry_exchange, //exchange
		false,
		nil,
	); err != nil {
		return err
	}

	// bind deleted queue process
	if err := ch.QueueBind(
		inventory_dlq_queue,    //queuename
		inventory_dlq_routing,  //routing key
		inventory_dlq_exchange, //exchange
		false,
		nil,
	); err != nil {
		return err
	}

	// // Bind
	// return ch.QueueBind(
	// 	"inventory.queue",
	// 	"order.created",
	// 	"order.events",
	// 	false,
	// 	nil,
	// )
	return nil
}
