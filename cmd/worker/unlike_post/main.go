package main

import (
	"context"
	"log"
	"time"

	"github.com/G-Villarinho/social-network/client"
	"github.com/G-Villarinho/social-network/config"
	"github.com/G-Villarinho/social-network/database"
	"github.com/G-Villarinho/social-network/domain"
	"github.com/G-Villarinho/social-network/pkg" // ajuste conforme necessário
	"github.com/G-Villarinho/social-network/repository"
	"github.com/G-Villarinho/social-network/service"
	"github.com/go-redis/redis/v8"
	jsoniter "github.com/json-iterator/go"
	"gorm.io/gorm"
)

func main() {
	config.ConfigureLogger()
	config.LoadEnvironments()

	di := pkg.NewDi()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	db, err := database.NewMysqlConnection(ctx)
	if err != nil {
		log.Fatal("error to connect to mysql: ", err)
	}

	redisClient, err := database.NewRedisConnection(ctx)
	if err != nil {
		log.Fatal("error to connect to redis: ", err)
	}

	pkg.Provide(di, func(d *pkg.Di) (*gorm.DB, error) {
		return db, nil
	})

	pkg.Provide(di, func(d *pkg.Di) (*redis.Client, error) {
		return redisClient, nil
	})

	rabbitMQClient, err := client.NewRabbitMQClient(di)
	if err != nil {
		log.Fatal("error initializing RabbitMQ client: ", err)
	}
	if err := rabbitMQClient.Connect(); err != nil {
		log.Fatal("error connecting to RabbitMQ: ", err)
	}
	defer func() {
		if err := rabbitMQClient.Disconnect(); err != nil {
			log.Println("error disconnecting from RabbitMQ:", err)
		}
	}()

	pkg.Provide(di, func(d *pkg.Di) (client.RabbitMQClient, error) {
		return rabbitMQClient, nil
	})

	pkg.Provide(di, service.NewContextService)
	pkg.Provide(di, service.NewPostService)
	pkg.Provide(di, service.NewQueueService)
	pkg.Provide(di, service.NewSessionService)

	pkg.Provide(di, repository.NewMemoryCacheRepository)
	pkg.Provide(di, repository.NewPostRepository)
	pkg.Provide(di, repository.NewSessionRepository)

	postService, err := pkg.Invoke[domain.PostService](di)
	if err != nil {
		log.Fatal("error to create post service: ", err)
	}

	queueService, err := pkg.Invoke[domain.QueueService](di)
	if err != nil {
		log.Fatal("error to create queue service: ", err)
	}

	for {
		messages, err := queueService.Consume(config.QueueUnlikePost)
		if err != nil {
			log.Fatal("error to consume message from queue: ", err)
		}

		for message := range messages {
			var payload domain.LikePayload
			if err := jsoniter.Unmarshal(message, &payload); err != nil {
				log.Println("error unmarshalling like payload: ", err)
				continue
			}

			if err := postService.ProcessUnlikePost(context.Background(), payload); err != nil {
				log.Println("error processing like post: ", err)
			}

			log.Printf("like post processed: %s", payload.PostID)
		}
	}

}
