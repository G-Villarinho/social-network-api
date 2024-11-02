package service

import (
	"fmt"

	"github.com/G-Villarinho/social-network/client"
	"github.com/G-Villarinho/social-network/config"
	"github.com/G-Villarinho/social-network/domain"
	"github.com/G-Villarinho/social-network/pkg"
)

type queueService struct {
	di             *pkg.Di
	rabbitMQClient client.RabbitMQClient
}

func NewQueueService(di *pkg.Di) (domain.QueueService, error) {
	rabbitMQClient, err := pkg.Invoke[client.RabbitMQClient](di)
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
