package domain

type QueueService interface {
	Publish(queueName string, message []byte) error
}
