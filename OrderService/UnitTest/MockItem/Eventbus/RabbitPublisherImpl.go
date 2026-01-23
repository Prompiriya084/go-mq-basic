package unittest_eventbus

import ports_eventbus "github.com/Prompiriya084/go-mq/OrderService/Internal/Core/Ports/Eventbus"

type mockEventbus struct{}

func NewMockRabbitPublisher() ports_eventbus.EventBusPublisher {
	return &mockEventbus{}
}
func (m *mockEventbus) Publish(eventName string, payload any) error {
	return nil // simulate successful MQ send
}

// type EventBus[Tentity any] interface {
// 	Publish(queue string, body []byte) error
// 	Subscribe(queue string, handler func(data Tentity) error) error
// }
