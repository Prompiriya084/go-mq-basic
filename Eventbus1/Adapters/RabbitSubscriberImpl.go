package adapters_eventbus

import (
	ports_eventbus "github.com/Prompiriya084/go-mq/Eventbus/Ports"
	"github.com/rabbitmq/amqp091-go"
)

type rabbitSubscriber struct {
	ch       *amqp091.Channel
	exchange string
}

func NewRabbitSubscriber(ch *amqp091.Channel, exchange string) ports_eventbus.EventBusSubscriber {
	return &rabbitSubscriber{
		ch:       ch,
		exchange: exchange,
	}
}

func (s *rabbitSubscriber) Consume(eventName string, handler func(event string, body []byte) error) error {
	msgs, _ := s.ch.Consume(eventName, "", false, false, false, false, nil)

	for msg := range msgs {

		err := handler(msg.RoutingKey, msg.Body)

		if err == nil {
			msg.Ack(false)
			continue
		}

		retry := getRetryCount(msg)

		if retry >= maxRetry {
			// ❌ send to DLQ
			publishWithHeader(
				s.ch,
				"order.dlx",
				msg.RoutingKey,
				msg.Body,
				retry,
			)
			msg.Ack(false)
			continue
		}

		// 🔁 retry
		publishWithHeader(
			s.ch,
			"order.dlx",
			msg.RoutingKey,
			msg.Body,
			retry+1,
		)
		msg.Ack(false)
	}
}

func (s *rabbitSubscriber) publishWithHeader(
	ch *amqp091.Channel,
	exchange string,
	routingKey string,
	body []byte,
	retry int,
) error {

	return ch.Publish(
		exchange,
		routingKey,
		false,
		false,
		amqp091.Publishing{
			Headers: amqp091.Table{
				"x-retry-count": int32(retry),
			},
			Body: body,
		},
	)
}
