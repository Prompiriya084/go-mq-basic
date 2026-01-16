package domain_eventbus

type QueueConfig struct {
	Name      string
	RetryName string
	DLQName   string
	RetryTTL  int32
}
