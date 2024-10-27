package database

import (
	"context"

	"github.com/G-Villarinho/social-network/config"
	"github.com/go-redis/redis/v8"
)

func NewRedisConnection(ctx context.Context) (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     config.Env.Redis.Address,
		Password: config.Env.Redis.Password,
		DB:       config.Env.Redis.DB,
	})

	_, err := client.Ping(ctx).Result()
	if err != nil {
		return nil, err
	}

	return client, nil
}
