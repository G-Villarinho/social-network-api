package repository

import (
	"context"

	"github.com/G-Villarinho/social-network/domain"
	"github.com/G-Villarinho/social-network/internal"
	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type likeRepository struct {
	di          *internal.Di
	db          *gorm.DB
	redisClient *redis.Client
}

func NewLikeRepository(di *internal.Di) (domain.LikeRepository, error) {
	db, err := internal.Invoke[*gorm.DB](di)
	if err != nil {
		return nil, err
	}

	redisClient, err := internal.Invoke[*redis.Client](di)
	if err != nil {
		return nil, err
	}

	return &likeRepository{
		di:          di,
		db:          db,
		redisClient: redisClient,
	}, nil
}

func (l *likeRepository) CreateLike(ctx context.Context, like domain.Like) error {
	tx := l.db.WithContext(ctx).Begin()
	if tx.Error != nil {
		return tx.Error
	}

	if err := tx.Create(&like).Error; err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Model(&domain.Post{}).
		Where("id = ?", like.PostID).
		Updates(map[string]interface{}{"likes": gorm.Expr("likes + ?", 1)}).Error; err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}

func (l *likeRepository) UserLikedPost(ctx context.Context, ID uuid.UUID, userID uuid.UUID) (bool, error) {
	var like domain.Like

	if err := l.db.WithContext(ctx).
		Where("postId = ? AND userId = ?", ID, userID).
		First(&like).Error; err != nil {

		if err == gorm.ErrRecordNotFound {
			return false, nil
		}

		return false, err
	}

	return true, nil
}