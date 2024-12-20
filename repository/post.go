package repository

import (
	"context"
	"fmt"

	"github.com/G-Villarinho/social-network/domain"
	"github.com/G-Villarinho/social-network/internal"
	"github.com/go-redis/redis/v8"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type postRepository struct {
	di          *internal.Di
	db          *gorm.DB
	redisClient *redis.Client
}

func NewPostRepository(di *internal.Di) (domain.PostRepository, error) {
	db, err := internal.Invoke[*gorm.DB](di)
	if err != nil {
		return nil, err
	}

	redisClient, err := internal.Invoke[*redis.Client](di)
	if err != nil {
		return nil, err
	}

	return &postRepository{
		di:          di,
		db:          db,
		redisClient: redisClient,
	}, nil
}

func (p *postRepository) CreatePost(ctx context.Context, post domain.Post) error {
	if err := p.db.WithContext(ctx).
		Create(&post).Error; err != nil {
		return err
	}

	return nil
}

func (p *postRepository) GetPaginatedPosts(ctx context.Context, userID uuid.UUID, page int, limit int) (*domain.Pagination[*domain.Post], error) {
	pagination := &domain.Pagination[*domain.Post]{
		Limit: limit,
		Page:  page,
		Sort:  "createdAt desc",
	}

	subQuery := p.db.Table("Follower").Select("userId").Where("followerId = ?", userID)

	paginatedPosts, err := paginate(pagination,
		p.db.WithContext(ctx).
			Preload("Author").
			Where("authorId = ? OR authorId IN (?)", userID, subQuery))
	if err != nil {
		return nil, fmt.Errorf("error to get paginated feed in repository: %w", err)
	}

	return paginatedPosts, nil
}

func (p *postRepository) GetPostById(ctx context.Context, ID uuid.UUID, preload bool) (*domain.Post, error) {
	var post domain.Post

	query := p.db.WithContext(ctx)

	if preload {
		query = query.Preload("Author")
	}

	if err := query.
		Where("id = ?", ID).
		First(&post).Error; err != nil {

		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}

		return nil, err
	}

	return &post, nil
}

func (p *postRepository) UpdatePost(ctx context.Context, ID uuid.UUID, post domain.Post) error {
	if err := p.db.WithContext(ctx).
		Model(&post).
		Where("id = ?", ID).
		Updates(&post).Error; err != nil {
		return err
	}

	return nil
}

func (p *postRepository) DeletePost(ctx context.Context, ID uuid.UUID) error {
	if err := p.db.WithContext(ctx).
		Where("id = ?", ID).
		Delete(&domain.Post{}).Error; err != nil {
		return err
	}

	return nil
}

func (p *postRepository) GetByUserID(ctx context.Context, userID uuid.UUID) ([]*domain.Post, error) {
	var posts []*domain.Post

	if err := p.db.WithContext(ctx).
		Preload("Author").
		Where("authorId = ?", userID).
		Find(&posts).Error; err != nil {

		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}

		return nil, err
	}

	return posts, nil
}

func (p *postRepository) UnlikePost(ctx context.Context, ID uuid.UUID, userID uuid.UUID) error {
	tx := p.db.WithContext(ctx).Begin()
	if tx.Error != nil {
		return tx.Error
	}

	if err := tx.Where("postId = ? AND userId = ?", ID, userID).
		Delete(&domain.Like{}).Error; err != nil {

		tx.Rollback()
		return err
	}

	if err := tx.Model(&domain.Post{}).
		Where("id = ?", ID).
		Updates(map[string]interface{}{"likes": gorm.Expr("likes - ?", 1)}).Error; err != nil {

		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}
