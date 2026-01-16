package adapters_eventbus

type rabbitTopology struct {
    ch *amqp.Channel
}

func NewRabbitTopology(ch *amqp.Channel) *RabbitTopology {
    return &RabbitTopology{ch: ch}
}

func (t *rabbitTopology) DeclareQueue(cfg QueueConfig)