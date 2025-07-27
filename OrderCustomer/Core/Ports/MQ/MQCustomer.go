package ports_mq

type MQCustomer[Tentity any] interface {
	ReceiveMessage(queue string, handler func(data Tentity) error) error
}
