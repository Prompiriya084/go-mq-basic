package rabbitMQ

import (
	"os"

	"github.com/rabbitmq/amqp091-go"
)

func MQConnectionInit() (*amqp091.Channel, error) {
	// // Close previous connection/channel if they exist
	// if amqp091.channel != nil {
	// 	_ = c.channel.Close()
	// }
	// if c.conn != nil {
	// 	_ = c.conn.Close()
	// }
	// Dial new connection
	conn, err := amqp091.Dial(os.Getenv("RABBITMQ_URL"))
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

	order_core_exchange := os.Getenv("MQ_Inventory_Core_Exchange")
	order_core_routing := os.Getenv("MQ_Inventory_Core_Routing")
	order_core_queue := os.Getenv("MQ_Inventory_Core_Queue")
	order_retry_exchange := os.Getenv("MQ_Inventory_Retry_Exchange")
	order_retry_routing := os.Getenv("MQ_Inventory_Retry_Routing")
	order_retry_queue := os.Getenv("MQ_Inventory_Retry_Queue")
	order_dlq_exchange := os.Getenv("MQ_Inventory_DLQ_Exchange")
	order_dlq_routing := os.Getenv("MQ_Inventory_DLQ_Routing")
	order_dlq_queue := os.Getenv("MQ_Inventory_DLQ_Queue")

	if err := ch.ExchangeDeclare(
		order_core_exchange,
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
		order_retry_exchange,
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
		order_dlq_exchange,
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
		order_core_queue,
		true,
		false,
		false,
		false,
		amqp091.Table{
			"x-dead-letter-exchange":    order_retry_exchange,
			"x-dead-letter-routing-key": order_retry_routing,
		},
	)
	if err != nil {
		return err
	}

	// retry queue
	_, err = ch.QueueDeclare(
		order_retry_queue,
		true,
		false,
		false,
		false,
		amqp091.Table{
			"x-message-ttl":             5000,
			"x-dead-letter-exchange":    order_dlq_exchange,
			"x-dead-letter-routing-key": order_dlq_routing,
		},
	)
	if err != nil {
		return err
	}
	// dlq
	_, err = ch.QueueDeclare(
		order_dlq_queue,
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
		order_core_queue,    //queuename
		order_core_routing,  //routing key
		order_core_exchange, //exchange
		false,
		nil,
	); err != nil {
		return err
	}

	// bind retry process
	if err := ch.QueueBind(order_retry_queue, //queuename
		order_retry_routing,  //routing key
		order_retry_exchange, //exchange
		false,
		nil,
	); err != nil {
		return err
	}

	// bind deleted queue process
	if err := ch.QueueBind(
		order_dlq_queue,    //queuename
		order_dlq_routing,  //routing key
		order_dlq_exchange, //exchange
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
