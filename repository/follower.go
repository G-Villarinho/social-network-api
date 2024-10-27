package repository

import (
	"context"

	"github.com/G-Villarinho/social-network/domain"
	"github.com/G-Villarinho/social-network/pkg"
	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
)

type followerRepository struct {
	di          *pkg.Di
	db          *gorm.DB
	redisClient *redis.Client
}

func NewFollowerRepository(di *pkg.Di) (domain.FollowerRepository, error) {
	db, err := pkg.Invoke[*gorm.DB](di)
	if err != nil {
		return nil, err
	}

	redisClient, err := pkg.Invoke[*redis.Client](di)
	if err != nil {
		return nil, err
	}

	return &followerRepository{
		di:          di,
		db:          db,
		redisClient: redisClient,
	}, nil
}

func (f *followerRepository) CreateFollower(ctx context.Context, follower domain.Follower) error {
	if err := f.db.WithContext(ctx).Create(&follower).Error; err != nil {
		return err
	}

	return nil
}
