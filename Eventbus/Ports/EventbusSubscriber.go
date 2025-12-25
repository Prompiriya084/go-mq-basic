package ports_eventbus

type EventBusSubscriber[Tentity any] interface {
	Subscribe(queue string, handler func(data Tentity) error) error
}
