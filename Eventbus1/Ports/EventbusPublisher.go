package ports_eventbus

type EventBusPublisher interface {
	Publish(eventName string, payload any) error
}
