package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/G-Villarinho/social-network/domain"
	"github.com/G-Villarinho/social-network/pkg"
	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
	jsoniter "github.com/json-iterator/go"
	"gorm.io/gorm"
)

type postRepository struct {
	di          *pkg.Di
	db          *gorm.DB
	redisClient *redis.Client
}

func NewPostRepository(di *pkg.Di) (domain.PostRepository, error) {
	db, err := pkg.Invoke[*gorm.DB](di)
	if err != nil {
		return nil, err
	}

	redisClient, err := pkg.Invoke[*redis.Client](di)
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

func (p *postRepository) LikePost(ctx context.Context, like domain.Like) error {
	tx := p.db.WithContext(ctx).Begin()
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

func (p *postRepository) HasUserLikedPost(ctx context.Context, ID uuid.UUID, userID uuid.UUID) (bool, error) {
	var like domain.Like

	if err := p.db.WithContext(ctx).
		Where("postId = ? AND userId = ?", ID, userID).
		First(&like).Error; err != nil {

		if err == gorm.ErrRecordNotFound {
			return false, nil
		}

		return false, err
	}

	return true, nil
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

func (p *postRepository) GetLikedPostIDs(ctx context.Context, userID uuid.UUID) (map[uuid.UUID]struct{}, error) {
	var likes []domain.Like
	if err := p.db.WithContext(ctx).
		Where("userID = ?", userID).
		Find(&likes).Error; err != nil {
		return nil, err
	}

	likedPostIDs := make(map[uuid.UUID]struct{})
	for _, like := range likes {
		likedPostIDs[like.PostID] = struct{}{}
	}

	return likedPostIDs, nil
}

func (p *postRepository) GetLikesByPostIDs(ctx context.Context, userID uuid.UUID, postIDs []uuid.UUID) ([]uuid.UUID, error) {
	var likes []domain.Like
	if err := p.db.WithContext(ctx).
		Where("userID = ? AND postID IN (?)", userID, postIDs).
		Find(&likes).Error; err != nil {
		return nil, err
	}

	var likedPostIDs []uuid.UUID
	for _, like := range likes {
		likedPostIDs = append(likedPostIDs, like.PostID)
	}

	return likedPostIDs, nil
}

func (p *postRepository) GetCachedPosts(ctx context.Context, cacheKey string) (*domain.Pagination[*domain.PostResponse], error) {
	var cachedFeed domain.Pagination[*domain.PostResponse]

	data, err := p.redisClient.Get(ctx, cacheKey).Result()
	if err != nil {
		if err == redis.Nil {
			return nil, nil
		}

		return nil, err
	}

	if err := jsoniter.UnmarshalFromString(data, &cachedFeed); err != nil {
		return nil, err
	}
	return &cachedFeed, nil
}

func (p *postRepository) CachePost(ctx context.Context, cacheKey string, feed *domain.Pagination[*domain.PostResponse]) error {
	data, err := jsoniter.MarshalToString(feed)
	if err != nil {
		return err
	}

	return p.redisClient.Set(ctx, cacheKey, data, 5*time.Minute).Err()
}
