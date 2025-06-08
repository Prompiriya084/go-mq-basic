package adapters_customers

import (
	"encoding/json"
	"log"
	"time"

	ports_customers "github.com/Prompiriya084/go-mq/Customer/Core/Ports/MQ"
	"github.com/rabbitmq/amqp091-go"
)

type mqCustomerImpl[Tentity any] struct {
	connStr string
	conn    *amqp091.Connection
	channel *amqp091.Channel
}

func NewMQCustomer[Tentity any](connStr string) ports_customers.MQCustomer[Tentity] {
	return &mqCustomerImpl[Tentity]{
		connStr: connStr,
	}
}

func (c *mqCustomerImpl[Tentity]) connect(queue string) (<-chan amqp091.Delivery, error) {
	var err error
	// Close previous connection/channel if they exist
	if c.channel != nil {
		_ = c.channel.Close()
	}
	if c.conn != nil {
		_ = c.conn.Close()
	}
	// Dial new connection
	c.conn, err = amqp091.Dial(c.connStr)
	if err != nil {
		return nil, err
	}

	c.channel, err = c.conn.Channel()
	if err != nil {
		return nil, err
	}
	// Optional: Set prefetch
	if err := c.channel.Qos(1, 0, false); err != nil {
		return nil, err
	}
	// Declare queue
	_, err = c.channel.QueueDeclare(queue, true, false, false, false, nil)
	if err != nil {
		return nil, err
	}

	msgs, err := c.channel.Consume(queue, "", false, false, false, false, nil)
	if err != nil {
		return nil, err
	}

	return msgs, nil
}

func (c *mqCustomerImpl[Tentity]) ReceiveMessage(queue string, handler func(data Tentity) error) error {
	for {
		msgs, err := c.connect(queue)
		if err != nil {
			log.Printf("Failed to connect to RabbitMQ: %v. Retrying in 5s...", err)
			time.Sleep(5 * time.Second)
			continue
		}
		log.Println("Connected to RabbitMQ. Waiting for messages...")
		for msg := range msgs {
			var model Tentity
			if err := json.Unmarshal(msg.Body, &model); err != nil {
				log.Printf("Failed to decode message: %v", err)
				_ = msg.Nack(false, false)
				continue
			}

			if err := handler(model); err != nil {
				log.Printf("Model handler failed: %v", err)
				_ = msg.Nack(false, true) // or requeue = false to prevent retry loop
				continue
			}

			if err := msg.Ack(false); err != nil {
				log.Printf("Failed to ack message: %v", err)
			}

		}

		log.Println("Message channel closed. Reconnecting in 5s...")
		time.Sleep(5 * time.Second)
	}
}
