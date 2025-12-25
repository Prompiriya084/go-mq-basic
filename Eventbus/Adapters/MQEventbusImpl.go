package adapters_eventbus

import (
	"fmt"
	"log"
	"sync"
	"time"

	ports_eventbus "github.com/Prompiriya084/go-mq/Eventbus/Ports"
	"github.com/rabbitmq/amqp091-go"
)

type mqEventBusImpl[Tentity any] struct {
	connStr    string
	conn       *amqp091.Connection
	pubChannel *amqp091.Channel
	subChannel *amqp091.Channel
	mu         sync.Mutex
}

func NewMQPublisher[Tentity any](connStr string) ports_eventbus.MQEventBus[Tentity] {
	return &mqEventBusImpl[Tentity]{
		connStr: connStr,
	}
}
func (b *mqEventBusImpl[Tentity]) DeclareQueue(queue string) error {
	ch, err := b.publisherChannel()
	if err != nil {
		return err
	}

	_, err = ch.QueueDeclare(
		queue,
		true,
		false,
		false,
		false,
		nil,
	)

	return err
}

func (b *mqEventBusImpl[Tentity]) connect() error {
	b.mu.Lock()
	defer b.mu.Unlock()

	if b.conn != nil && !b.conn.IsClosed() {
		return nil
	}

	conn, err := amqp091.Dial(b.connStr)
	if err != nil {
		return err
	}

	b.conn = conn
	log.Println("✅ RabbitMQ connected")

	go b.watchConnection()
	return nil
}
func (b *mqEventBusImpl[Tentity]) watchConnection() {
	errChan := b.conn.NotifyClose(make(chan *amqp091.Error))

	err := <-errChan
	if err != nil {
		log.Printf("❌ RabbitMQ connection closed: %v", err)
	}

	for {
		log.Println("🔁 Reconnecting to RabbitMQ...")
		if err := b.connect(); err == nil {
			return
		}
		time.Sleep(5 * time.Second)
	}
}
func (b *mqEventBusImpl[Tentity]) publisherChannel() (*amqp091.Channel, error) {
	if err := b.connect(); err != nil {
		return nil, err
	}

	if b.pubChannel != nil && !b.pubChannel.IsClosed() {
		return b.pubChannel, nil
	}

	ch, err := b.conn.Channel()
	if err != nil {
		return nil, err
	}

	b.pubChannel = ch
	return ch, nil
}

func (p *mqEventBusImpl[Tentity]) Publish(queue string, body []byte) error {
	if err := p.connect(); err != nil {
		log.Printf("connect failed: %v", err)
	}

	err := p.pubChannel.Publish("", queue, false, false, amqp091.Publishing{
		ContentType:  "application/json",
		Body:         body,
		DeliveryMode: amqp091.Persistent, // message survives restart
	})

	if err != nil {
		// Force reconnect and retry once
		log.Printf("❌ Publish failed, retrying: %v", err)
		// _ = p.InitalizePublisherConnection() // close and reconnect
		return p.Publish(queue, body)
	}

	fmt.Println("Publish successful.")

	return nil
}

func (b *mqEventBusImpl[Tentity]) consumerChannel() (*amqp091.Channel, error) {
	if err := b.connect(); err != nil {
		return nil, err
	}

	if b.subChannel != nil && !b.subChannel.IsClosed() {
		return b.subChannel, nil
	}

	ch, err := b.conn.Channel()
	if err != nil {
		return nil, err
	}

	if err := ch.Qos(1, 0, false); err != nil {
		return nil, err
	}

	b.subChannel = ch
	return ch, nil
}

func (b *mqEventBusImpl[Tentity]) Subscribe(
	queue string,
	handler func([]byte) error,
) error {

	ch, err := b.consumerChannel()
	if err != nil {
		return err
	}

	msgs, err := ch.Consume(
		queue,
		"",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}

	go func() {
		for msg := range msgs {
			if err := handler(msg.Body); err != nil {
				_ = msg.Nack(false, true)
				continue
			}
			_ = msg.Ack(false)
		}
	}()

	return nil
}
