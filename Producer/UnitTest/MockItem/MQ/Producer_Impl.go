package unittest_mq

type MockProducer struct{}

func (m *MockProducer) PublishMessage(queue string, body []byte) error {
	return nil // simulate successful MQ send
}
