package repository

import (
	"context"

	"github.com/G-Villarinho/social-network/domain"
	"github.com/G-Villarinho/social-network/pkg"
	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
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

func (f *followerRepository) DeleteFollower(ctx context.Context, followerId uuid.UUID) error {
	if err := f.db.WithContext(ctx).Where("id = ?", followerId).Delete(&domain.Follower{}).Error; err != nil {
		return err
	}

	return nil
}

func (f *followerRepository) GetFollower(ctx context.Context, userID uuid.UUID, followerID uuid.UUID) (*domain.Follower, error) {
	var follower *domain.Follower

	if err := f.db.WithContext(ctx).Where("userId = ? AND followerId = ?", userID, followerID).First(&follower).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}

	return follower, nil
}

func (f *followerRepository) GetFollowers(ctx context.Context, userID uuid.UUID) ([]*domain.Follower, error) {
	var followers []*domain.Follower

	if err := f.db.WithContext(ctx).Preload("Follower").Where("userId = ?", userID).Find(&followers).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}

		return nil, err
	}

	return followers, nil
}

func (f *followerRepository) GetFollowings(ctx context.Context, userID uuid.UUID) ([]*domain.Follower, error) {
	var followings []*domain.Follower

	if err := f.db.WithContext(ctx).Preload("User").Where("followerId = ?", userID).Find(&followings).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}

		return nil, err
	}

	return followings, nil
}
