package adapters_eventbus

import (
	"encoding/json"
	"fmt"
	"log"
	"sync"
	"time"

	ports_eventbus "github.com/Prompiriya084/go-mq/Eventbus/Ports"
	"github.com/rabbitmq/amqp091-go"
)

type mqEventBusSubScriberImpl[Tentity any] struct {
	connStr string
	conn    *amqp091.Connection
	channel *amqp091.Channel
	mu      sync.Mutex
}

func NewMQSubscriber[Tentity any](connStr string) ports_eventbus.EventBusSubscriber[Tentity] {
	return &mqEventBusSubScriberImpl[Tentity]{
		connStr: connStr,
	}
}

func (c *mqEventBusSubScriberImpl[Tentity]) connectConsumer(queue string) (<-chan amqp091.Delivery, error) {
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
	// ❌ Remove this — no create allowed
	// _, err := p.channel.QueueDeclare(...)
	// ✔ Declare ONLY as passive (must already exist)
	// _, err = b.channel.QueueDeclarePassive(queue, true, false, false, false, nil)
	// if err != nil {
	// 	return nil, err
	// }
	if err := c.ensureQueue(queue); err != nil {
		return nil, err
	}

	msgs, err := c.channel.Consume(queue, "", false, false, false, false, nil)
	if err != nil {
		return nil, err
	}

	return msgs, nil
}

func (c *mqEventBusSubScriberImpl[Tentity]) ensureQueue(queue string) error {
	// main
	errMessage := "Failed to declare queue " + queue + ": "
	_, err := c.channel.QueueDeclarePassive(queue, true, false, false, false, nil)
	if err == nil { // Queue exists
		return nil
	}

	_, err = c.channel.QueueDeclare(
		queue,
		true,  // durable
		false, // autoDelete
		false, // exclusive
		false, // noWait
		nil,
	)
	if err != nil {
		return fmt.Errorf("%s %w", errMessage, err)
	}

	// // failed (no retry)
	// _, err = b.channel.QueueDeclarePassive(
	// 	queue+".failed",
	// 	true,
	// 	false,
	// 	false,
	// 	false,
	// 	nil,
	// )
	// if err != nil { // Queue does NOT exist
	// 	_, err = b.channel.QueueDeclare(
	// 		queue+".failed",
	// 		true,
	// 		false,
	// 		false,
	// 		false,
	// 		nil,
	// 	)
	// 	if err != nil {
	// 		return fmt.Errorf("%s %w", errMessage, err)
	// 	}
	// }

	// // DLQ
	// _, err = b.channel.QueueDeclarePassive(
	// 	queue+".dlq",
	// 	true,
	// 	false,
	// 	false,
	// 	false,
	// 	nil,
	// )
	// if err != nil { // Queue does NOT exist
	// 	_, err = b.channel.QueueDeclare(
	// 		queue+".dlq",
	// 		true,
	// 		false,
	// 		false,
	// 		false,
	// 		nil,
	// 	)
	// 	if err != nil {
	// 		return fmt.Errorf("%s %w", errMessage, err)
	// 	}
	// }

	// // retry queue, dead-letter back to main after delay
	// _, err = b.channel.QueueDeclare(
	// 	queue+".retry",
	// 	true,
	// 	false,
	// 	false,
	// 	false,
	// 	amqp091.Table{
	// 		"x-message-ttl":             int32(10000), // retry after 10s
	// 		"x-dead-letter-exchange":    "testing",
	// 		"x-dead-letter-routing-key": queue,
	// 	},
	// )
	// if err != nil {
	// 	return fmt.Errorf("%s %w", errMessage, err)
	// }

	return nil
}
func (c *mqEventBusSubScriberImpl[Tentity]) Subscribe(queue string, handler func(data Tentity) error) error {
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
				// _ = msg.Ack(false)
				c.Publish(queue+".dlq", msg.Body) //Send the .dlq queue in case business logic error
				continue
			}

			if err := handler(model); err != nil {
				log.Printf("Model handler failed: %v", err)
				_ = msg.Nack(false, false)        // or requeue = false to prevent retry loop
				c.Publish(queue+".dlq", msg.Body) //Send the .dlq queue in case business logic error
				// _ = msg.Ack(false)
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
