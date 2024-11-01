package service

import (
	"context"
	"fmt"
	"log/slog"
	"sync"
	"time"

	"github.com/G-Villarinho/social-network/client"
	"github.com/G-Villarinho/social-network/config"
	"github.com/G-Villarinho/social-network/domain"
	"github.com/G-Villarinho/social-network/pkg"
	jsoniter "github.com/json-iterator/go"
)

type likeQueueService struct {
	di             *pkg.Di
	mu             sync.Mutex
	likes          []domain.LikePayload
	rabbitMQClient client.RabbitMQClient
	interval       time.Duration
	maxBatch       int
}

func NewLikeQueueService(di *pkg.Di) (domain.LikeQueueService, error) {
	rabbitMQClient, err := pkg.Invoke[client.RabbitMQClient](di)
	if err != nil {
		return nil, fmt.Errorf("error to initialize RabbitMQ client: %w", err)
	}

	service := &likeQueueService{
		di:             di,
		rabbitMQClient: rabbitMQClient,
		interval:       10 * time.Second,
		maxBatch:       100,
	}

	go service.startProcessor()
	return service, nil
}

func (l *likeQueueService) AddLike(ctx context.Context, payload domain.LikePayload) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.likes = append(l.likes, payload)

	if len(l.likes) >= l.maxBatch {
		l.flush()
	}
}

func (l *likeQueueService) startProcessor() {
	ticker := time.NewTicker(l.interval)
	defer ticker.Stop()

	for range ticker.C {
		l.mu.Lock()
		if len(l.likes) > 0 {
			l.flush()
		}
		l.mu.Unlock()
	}
}

func (l *likeQueueService) flush() {
	if len(l.likes) == 0 {
		return
	}

	data, err := jsoniter.Marshal(l.likes)
	if err != nil {
		slog.Error("error marshalling likes batch", slog.String("error", err.Error()))
		return
	}

	if err := l.rabbitMQClient.Publish(config.QueueLikePost, data); err != nil {
		slog.Error("error publishing likes batch to RabbitMQ", slog.String("error", err.Error()))
	} else {
		slog.Info("published batch of likes", slog.Int("count", len(l.likes)))
	}

	l.likes = nil
}
