package domain

//go:generate mockery --name=QueueService --output=../mocks --outpkg=mocks

const (
	QueueSendEmail  = "send_email_queue"
	QueueLikePost   = "like_post_queue"
	QueueUnlikePost = "unlike_post_queue"
)

type QueueService interface {
	Publish(queueName string, message []byte) error
	Consume(queueName string) (<-chan []byte, error)
}
