package domain

//go:generate mockery --name=QueueService --output=../mocks --outpkg=mocks

type QueueService interface {
	Publish(queueName string, message []byte) error
	Consume(queueName string) (<-chan []byte, error)
}
