package eventbus

type EventBus[Tentity any] interface {
	Publish(queue string, body []byte) error
	Subscribe(queue string, handler func(data Tentity) error) error
}
