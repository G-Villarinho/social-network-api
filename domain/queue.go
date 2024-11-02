package domain

type QueueService interface {
	Publish(queueName string, message []byte) error
	Consume(queueName string) (<-chan []byte, error)
}
