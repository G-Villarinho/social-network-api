package repository

import (
	"context"
	"time"

	"github.com/G-Villarinho/social-network/domain"
	"github.com/G-Villarinho/social-network/pkg"
	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
)

type memoryCacheRepository struct {
	di          *pkg.Di
	redisClient *redis.Client
}

func NewMemoryCacheRepository(di *pkg.Di) (domain.MemoryCacheRepository, error) {
	redisClient, err := pkg.Invoke[*redis.Client](di)
	if err != nil {
		return nil, err
	}

	return &memoryCacheRepository{
		di:          di,
		redisClient: redisClient,
	}, nil
}

func (m *memoryCacheRepository) SetPostLike(ctx context.Context, postID uuid.UUID, userID uuid.UUID) error {
	if err := m.redisClient.Set(ctx, getLikeCacheKey(postID, userID), "liked", 5*time.Minute).Err(); err != nil {
		return err
	}

	return nil
}

func (m *memoryCacheRepository) RemovePostLike(ctx context.Context, postID uuid.UUID, userID uuid.UUID) error {
	if err := m.redisClient.Del(ctx, getLikeCacheKey(postID, userID)).Err(); err != nil {
		return err
	}

	return nil
}

func getLikeCacheKey(postID uuid.UUID, userID uuid.UUID) string {
	return "like:" + postID.String() + ":" + userID.String()
}
