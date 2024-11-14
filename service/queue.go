package service

import (
	"fmt"

	"github.com/G-Villarinho/social-network/client"
	"github.com/G-Villarinho/social-network/config"
	"github.com/G-Villarinho/social-network/internal"

	"github.com/G-Villarinho/social-network/domain"
)

type queueService struct {
	di             *internal.Di
	rabbitMQClient client.RabbitMQClient
}

func NewQueueService(di *internal.Di) (domain.QueueService, error) {
	rabbitMQClient, err := internal.Invoke[client.RabbitMQClient](di)
	if err != nil {
		return nil, err
	}

	return &queueService{
		di:             di,
		rabbitMQClient: rabbitMQClient,
	}, nil
}

func (q *queueService) Publish(queueName string, message []byte) error {
	if err := q.rabbitMQClient.Publish(config.QueueLikePost, message); err != nil {
		return fmt.Errorf("error publishing message to queue: %w", err)
	}

	return nil
}

func (q *queueService) Consume(queueName string) (<-chan []byte, error) {
	messages, err := q.rabbitMQClient.Consume(queueName)
	if err != nil {
		return nil, fmt.Errorf("error consuming message from queue: %w", err)
	}

	return messages, nil
}
