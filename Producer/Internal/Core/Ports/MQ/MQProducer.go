package ports_mq

type MQProducer interface {
	PublishMessage(queueName string, body []byte) error
}
