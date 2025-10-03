package eventbus

import (
	"encoding/json"
	"log"
	"sync"
	"time"

	"github.com/rabbitmq/amqp091-go"
)

type mqEventBusImpl[Tentity any] struct {
	connStr string
	conn    *amqp091.Connection
	channel *amqp091.Channel
	mu      sync.Mutex
}

func NewMQEventbus[Tentity any](connStr string) EventBus[Tentity] {
	return &mqEventBusImpl[Tentity]{
		connStr: connStr,
	}
}

// connect for publisher
func (c *mqEventBusImpl[Tentity]) connectPublisher() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	// close old
	if c.channel != nil {
		_ = c.channel.Close()
	}
	if c.conn != nil {
		_ = c.conn.Close()
	}

	conn, err := amqp091.Dial(c.connStr)
	if err != nil {
		return err
	}
	ch, err := conn.Channel()
	if err != nil {
		_ = conn.Close()
		return err
	}
	c.conn = conn
	c.channel = ch
	return nil
}

func (c *mqEventBusImpl[Tentity]) connectConsumer(queue string) (<-chan amqp091.Delivery, error) {
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
func (p *mqEventBusImpl[Tentity]) Publish(queue string, body []byte) error {
	if err := p.connectPublisher(); err != nil {
		log.Printf("connect failed: %v", err)
	}

	_, err := p.channel.QueueDeclare(
		queue,
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}
	err = p.channel.Publish("", queue, false, false, amqp091.Publishing{
		ContentType:  "application/json",
		Body:         body,
		DeliveryMode: amqp091.Persistent, // message survives restart
	})

	if err != nil {
		// Force reconnect and retry once
		log.Printf("❌ Publish failed, retrying: %v", err)
		_ = p.connectPublisher() // close and reconnect
		return p.Publish(queue, body)
	}

	return nil
}

func (c *mqEventBusImpl[Tentity]) Subscribe(queue string, handler func(data Tentity) error) error {
	for {
		msgs, err := c.connectConsumer(queue)
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
