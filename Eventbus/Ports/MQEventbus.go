package ports_eventbus

type MQEventBus[Tentity any] interface {
	Publish(queue string, body []byte) error
	Subscribe(queue string, handler func([]byte) error) error
}
