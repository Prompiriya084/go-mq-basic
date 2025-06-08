package adapters_producers

import (
	"log"
	"sync"

	ports_mq "github.com/Prompiriya084/go-mq/Producer/Internal/Core/Ports/MQ"
	"github.com/rabbitmq/amqp091-go"
)

type mqProducerImpl struct {
	connStr string
	conn    *amqp091.Connection
	chanel  *amqp091.Channel
	mu      sync.Mutex
}

func NewMQProducer(connStr string) ports_mq.MQProducer {
	return &mqProducerImpl{
		connStr: connStr,
	}
}
func (p *mqProducerImpl) connect() error {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.conn != nil && !p.conn.IsClosed() {
		return nil
	}
	//Close old
	if p.chanel != nil {
		_ = p.chanel.Close()
	}
	if p.conn != nil {
		_ = p.conn.Close()
	}
	var err error

	p.conn, err = amqp091.Dial(p.connStr)
	if err != nil {
		return err
	}

	p.chanel, err = p.conn.Channel()
	if err != nil {
		return nil
	}

	return nil
}
func (p *mqProducerImpl) PublishMessage(queueName string, body []byte) error {
	if err := p.connect(); err != nil {
		log.Printf("connect failed: %v", err)
	}

	_, err := p.chanel.QueueDeclare(
		queueName,
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}
	err = p.chanel.Publish("", queueName, false, false, amqp091.Publishing{
		ContentType:  "application/json",
		Body:         body,
		DeliveryMode: amqp091.Persistent, // message survives restart
	})

	if err != nil {
		// Force reconnect and retry once
		log.Printf("❌ Publish failed, retrying: %v", err)
		_ = p.connect() // close and reconnect
		return p.PublishMessage(queueName, body)
	}

	return nil
}
