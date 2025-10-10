package unittest_eventbus

type MockEventbus[Tentity any] struct{}

func (m *MockEventbus[Tentity]) Publish(queue string, body []byte) error {
	return nil // simulate successful MQ send
}
func (m *MockEventbus[Tentity]) Subscribe(queue string, handler func(data Tentity) error) error {
	return nil // simulate successful MQ send
}

// type EventBus[Tentity any] interface {
// 	Publish(queue string, body []byte) error
// 	Subscribe(queue string, handler func(data Tentity) error) error
// }
