package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/G-Villarinho/social-network/config"
	"github.com/G-Villarinho/social-network/domain"
	"github.com/G-Villarinho/social-network/internal"
	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
	jsoniter "github.com/json-iterator/go"
)

type memoryCacheRepository struct {
	di          *internal.Di
	redisClient *redis.Client
}

func NewMemoryCacheRepository(di *internal.Di) (domain.MemoryCacheRepository, error) {
	redisClient, err := internal.Invoke[*redis.Client](di)
	if err != nil {
		return nil, err
	}

	return &memoryCacheRepository{
		di:          di,
		redisClient: redisClient,
	}, nil
}

func (m *memoryCacheRepository) SetPostLike(ctx context.Context, postID uuid.UUID, userID uuid.UUID) error {
	if err := m.redisClient.Set(ctx, getLikeCacheKey(postID, userID), "liked", time.Duration(config.Env.Cache.CacheExp)*time.Minute).Err(); err != nil {
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

func (m *memoryCacheRepository) SetPost(ctx context.Context, userID uuid.UUID, posts *domain.Pagination[*domain.PostResponse], page, limit int) error {
	JSON, err := jsoniter.Marshal(posts)
	if err != nil {
		return err
	}

	if err := m.redisClient.
		Set(ctx, getPostCacheKey(userID, page, limit), JSON, time.Duration(config.Env.Cache.CacheExp)*time.Minute).
		Err(); err != nil {
		return err
	}

	return nil
}

func (m *memoryCacheRepository) GetPosts(ctx context.Context, userID uuid.UUID, page int, limit int) (*domain.Pagination[*domain.PostResponse], error) {
	posts := new(domain.Pagination[*domain.PostResponse])

	JSON, err := m.redisClient.Get(ctx, getPostCacheKey(userID, page, limit)).Result()
	if err != nil {
		if err == redis.Nil {
			return nil, nil
		}
		return nil, err
	}

	if err := jsoniter.UnmarshalFromString(JSON, posts); err != nil {
		return nil, err
	}

	return posts, nil
}

func (m *memoryCacheRepository) GetCachedLikes(ctx context.Context, userID uuid.UUID, postIDs []uuid.UUID) (*domain.LikeCache, error) {
	likeCache := new(domain.LikeCache)

	for _, postID := range postIDs {
		key := getLikeCacheKey(postID, userID)
		liked, err := m.redisClient.Get(ctx, key).Result()

		if err == redis.Nil {
			likeCache.MissingLikes = append(likeCache.MissingLikes, postID)
			continue
		}

		if err != nil {
			return nil, fmt.Errorf("error fetching from cache: %w", err)
		}

		if liked == "liked" {
			likeCache.CachedLikes = append(likeCache.CachedLikes, postID)
		}
	}

	return likeCache, nil
}

func (m *memoryCacheRepository) SetLikesByPostIDs(ctx context.Context, userID uuid.UUID, postIDs []uuid.UUID) error {
	for _, postID := range postIDs {
		key := getLikeCacheKey(postID, userID)
		if err := m.redisClient.Set(ctx, key, "liked", time.Duration(config.Env.Cache.CacheExp)*time.Minute).Err(); err != nil {
			return err
		}
	}
	return nil
}

func getLikeCacheKey(postID uuid.UUID, userID uuid.UUID) string {
	return fmt.Sprintf("like:%s:%s", postID.String(), userID.String())
}

func getPostCacheKey(userID uuid.UUID, page, limit int) string {
	return fmt.Sprintf("user:%s:feed:page:%d:limit:%d", userID, page, limit)
}
