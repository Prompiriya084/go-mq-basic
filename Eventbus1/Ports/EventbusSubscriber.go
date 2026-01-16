package ports_eventbus

type EventBusSubscriber[Tentity any] interface {
	Consume(eventName string, handler func(event string, body []byte) error) error
}
